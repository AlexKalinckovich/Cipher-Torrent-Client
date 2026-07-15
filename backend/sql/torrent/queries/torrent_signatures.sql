-- name: CreateSignature :execresult
INSERT INTO torrent_signatures (torrent_hash, user_id, signature_blob, payload_hash)
VALUES (?, ?, ?, ?);

-- name: CreateSignatureMap :execresult
INSERT INTO torrent_signature_map (torrent_info_hash, creator_public_key, signature_id, signer_id, trust_level)
VALUES (?, ?, ?, ?, ?);

-- name: GetSignatureMapByTorrentAndSigner :one
SELECT * FROM torrent_signature_map
WHERE torrent_info_hash = ? AND creator_public_key = ? AND signer_id = ?
    LIMIT 1;

-- name: GetSignaturesWithSignerKey :many
SELECT
    ts.id AS signature_id,
    ts.signature_blob,
    ts.torrent_hash,
    ts.payload_hash,
    ts.created_at,
    tsm.trust_level,
    tsm.signer_id,
    u.public_key AS signer_public_key
FROM torrent_signature_map tsm
         JOIN torrent_signatures ts ON tsm.signature_id = ts.id
         JOIN users u ON tsm.signer_id = u.id
WHERE tsm.torrent_info_hash = ? AND tsm.creator_public_key = ?;

-- name: DeleteSignatureMapByTorrent :execresult
DELETE FROM torrent_signature_map
WHERE torrent_info_hash = ? AND creator_public_key = ?;

-- name: DeleteSignature :execresult
DELETE FROM torrent_signatures WHERE id = ?;