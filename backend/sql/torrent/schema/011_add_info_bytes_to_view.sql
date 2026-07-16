-- +goose Up

CREATE OR REPLACE VIEW v_torrent_with_signatures AS
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
    t.info_bytes,
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

CREATE OR REPLACE VIEW v_torrent_with_signatures AS
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