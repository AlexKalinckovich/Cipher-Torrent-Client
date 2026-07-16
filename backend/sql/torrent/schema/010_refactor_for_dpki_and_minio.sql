-- +goose Up

-- 1. Drop the dependent view first to avoid locking issues
DROP VIEW IF EXISTS v_torrent_with_signatures;

-- ============================================
-- PHASE 1: Update parent table (torrents) FIRST
-- ============================================
ALTER TABLE torrents
    ADD COLUMN creator_public_key BINARY(32) NOT NULL DEFAULT 0x0000000000000000000000000000000000000000000000000000000000000000 AFTER info_hash,
    ADD COLUMN creator_user_id BIGINT NOT NULL DEFAULT 1 AFTER creator_public_key;

ALTER TABLE torrents
DROP PRIMARY KEY,
    ADD PRIMARY KEY (info_hash, creator_public_key);

ALTER TABLE torrents
    ADD INDEX idx_creator_user (creator_user_id);

ALTER TABLE torrents
    ADD CONSTRAINT fk_torrent_creator
        FOREIGN KEY (creator_user_id) REFERENCES users(id)
            ON DELETE CASCADE;

-- ============================================
-- PHASE 2: Update Child Tables
-- ============================================

ALTER TABLE user_torrents
DROP FOREIGN KEY fk_ut_torrent,
    DROP FOREIGN KEY fk_ut_user,
    DROP PRIMARY KEY;

ALTER TABLE user_torrents
    ADD COLUMN creator_public_key BINARY(32) NOT NULL DEFAULT 0x0000000000000000000000000000000000000000000000000000000000000000 AFTER torrent_info_hash;

ALTER TABLE user_torrents
    ADD PRIMARY KEY (user_id, torrent_info_hash, creator_public_key);

ALTER TABLE user_torrents
    ADD CONSTRAINT fk_ut_user
        FOREIGN KEY (user_id) REFERENCES users(id)
            ON DELETE CASCADE;

ALTER TABLE user_torrents
    ADD CONSTRAINT fk_ut_torrent
        FOREIGN KEY (torrent_info_hash, creator_public_key)
            REFERENCES torrents(info_hash, creator_public_key)
            ON DELETE CASCADE;

ALTER TABLE torrent_signature_map
DROP FOREIGN KEY fk_tsm_torrent,
    DROP FOREIGN KEY fk_tsm_sig,
    DROP FOREIGN KEY fk_tsm_signer,
    DROP PRIMARY KEY,
DROP INDEX idx_torrent_signer,
DROP INDEX idx_signer;

ALTER TABLE torrent_signature_map
    ADD COLUMN creator_public_key BINARY(32) NOT NULL DEFAULT 0x0000000000000000000000000000000000000000000000000000000000000000 AFTER torrent_info_hash;

ALTER TABLE torrent_signature_map
    ADD PRIMARY KEY (torrent_info_hash, creator_public_key, signature_id);

ALTER TABLE torrent_signature_map
    ADD INDEX idx_torrent_signer (torrent_info_hash, creator_public_key, signer_id);

ALTER TABLE torrent_signature_map
    ADD INDEX idx_signer (signer_id);

ALTER TABLE torrent_signature_map
    ADD CONSTRAINT fk_tsm_torrent
        FOREIGN KEY (torrent_info_hash, creator_public_key)
            REFERENCES torrents(info_hash, creator_public_key)
            ON DELETE CASCADE;

ALTER TABLE torrent_signature_map
    ADD CONSTRAINT fk_tsm_sig
        FOREIGN KEY (signature_id) REFERENCES torrent_signatures(id)
            ON DELETE CASCADE;

ALTER TABLE torrent_signature_map
    ADD CONSTRAINT fk_tsm_signer
        FOREIGN KEY (signer_id) REFERENCES users(id)
            ON DELETE CASCADE;

-- ============================================
-- PHASE 3: Cleanup Defaults
-- ============================================
ALTER TABLE torrents ALTER COLUMN creator_public_key DROP DEFAULT;
ALTER TABLE torrents ALTER COLUMN creator_user_id DROP DEFAULT;
ALTER TABLE user_torrents ALTER COLUMN creator_public_key DROP DEFAULT;
ALTER TABLE torrent_signature_map ALTER COLUMN creator_public_key DROP DEFAULT;

-- ============================================
-- PHASE 4: Recreate View
-- ============================================
CREATE VIEW v_torrent_with_signatures AS
SELECT
    t.info_hash,
    t.creator_public_key,
    t.creator_user_id,
    t.name,
    t.size_bytes,
    t.piece_length,
    t.is_private,
    t.storage_path,
    t.added_at,
    tsm.signature_id,
    ts.signature_blob,
    ts.payload_hash,
    tsm.signer_id,
    u.public_key AS signer_public_key,
    ts.created_at AS signature_created_at
FROM torrents t
         LEFT JOIN torrent_signature_map tsm
                   ON t.info_hash = tsm.torrent_info_hash
                       AND t.creator_public_key = tsm.creator_public_key
         LEFT JOIN torrent_signatures ts
                   ON tsm.signature_id = ts.id
         LEFT JOIN users u
                   ON tsm.signer_id = u.id;

-- +goose Down

DROP VIEW IF EXISTS v_torrent_with_signatures;

ALTER TABLE torrent_signature_map
DROP FOREIGN KEY fk_tsm_torrent,
    DROP FOREIGN KEY fk_tsm_sig,
    DROP FOREIGN KEY fk_tsm_signer,
    DROP PRIMARY KEY,
DROP INDEX idx_torrent_signer,
DROP INDEX idx_signer,
DROP COLUMN creator_public_key;

ALTER TABLE torrent_signature_map
    ADD PRIMARY KEY (torrent_info_hash, signature_id),
    ADD INDEX idx_torrent_signer (torrent_info_hash, signer_id),
    ADD INDEX idx_signer (signer_id);

ALTER TABLE torrent_signature_map
    ADD CONSTRAINT fk_tsm_torrent
        FOREIGN KEY (torrent_info_hash) REFERENCES torrents(info_hash) ON DELETE CASCADE;

ALTER TABLE torrent_signature_map
    ADD CONSTRAINT fk_tsm_sig
        FOREIGN KEY (signature_id) REFERENCES torrent_signatures(id) ON DELETE CASCADE;

ALTER TABLE torrent_signature_map
    ADD CONSTRAINT fk_tsm_signer
        FOREIGN KEY (signer_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE user_torrents
DROP FOREIGN KEY fk_ut_torrent,
    DROP FOREIGN KEY fk_ut_user,
    DROP PRIMARY KEY,
    DROP COLUMN creator_public_key;

ALTER TABLE user_torrents
    ADD PRIMARY KEY (user_id, torrent_info_hash);

ALTER TABLE user_torrents
    ADD CONSTRAINT fk_ut_user
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE user_torrents
    ADD CONSTRAINT fk_ut_torrent
        FOREIGN KEY (torrent_info_hash) REFERENCES torrents(info_hash) ON DELETE CASCADE;

ALTER TABLE torrents
    DROP FOREIGN KEY fk_torrent_creator,
    DROP INDEX idx_creator_user,
    DROP PRIMARY KEY,
        DROP COLUMN creator_user_id,
        DROP COLUMN creator_public_key;

ALTER TABLE torrents
    ADD PRIMARY KEY (info_hash);

CREATE VIEW v_torrent_with_signatures AS
SELECT
    t.info_hash,
    t.name,
    t.size_bytes,
    t.piece_length,
    t.is_private,
    t.storage_path,
    t.added_at,
    tsm.signature_id,
    ts.signature_blob,
    ts.payload_hash,
    tsm.signer_id,
    u.public_key AS signer_public_key,
    ts.created_at AS signature_created_at
FROM torrents t
         LEFT JOIN torrent_signature_map tsm ON t.info_hash = tsm.torrent_info_hash
         LEFT JOIN torrent_signatures ts ON tsm.signature_id = ts.id
         LEFT JOIN users u ON tsm.signer_id = u.id;