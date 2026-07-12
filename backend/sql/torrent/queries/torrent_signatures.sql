-- name: CreateTorrentSignature :execresult
INSERT INTO torrent_signatures (torrent_hash, user_id, signature_blob, payload_hash)
VALUES (?, ?, ?, ?);

-- name: GetTorrentSignaturesByHash :many
SELECT *
FROM torrent_signatures
WHERE torrent_hash = ?;

-- name: GetTorrentSignatureByUserAndHash :one
SELECT *
FROM torrent_signatures
WHERE torrent_hash = ? AND user_id = ?
    LIMIT 1;

-- name: DeleteTorrentSignaturesByHash :execresult
DELETE FROM torrent_signatures WHERE torrent_hash = ?;