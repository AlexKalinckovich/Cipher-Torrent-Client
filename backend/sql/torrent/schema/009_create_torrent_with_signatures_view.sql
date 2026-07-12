-- +goose Up
CREATE OR REPLACE VIEW v_torrent_with_signatures AS
SELECT
    t.info_hash,
    t.info_bytes,
    t.name,
    t.size_bytes,
    t.piece_length,
    t.is_private,
    t.storage_path,
    t.added_at,
    ts.id AS signature_id,
    ts.signature_blob,
    ts.payload_hash,
    ts.user_id,
    u.public_key AS signer_public_key,
    ts.created_at AS signature_created_at
FROM torrents t
         LEFT JOIN torrent_signatures ts ON t.info_hash = ts.torrent_hash
         LEFT JOIN users u ON ts.user_id = u.id;

-- +goose Down
DROP VIEW IF EXISTS v_torrent_with_signatures;