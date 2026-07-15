-- name: CreateTorrent :execresult
INSERT INTO torrents (info_hash,info_bytes, creator_public_key, creator_user_id, name, size_bytes, piece_length, is_private)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: CreateUserTorrent :execresult
INSERT INTO user_torrents (user_id, torrent_info_hash, creator_public_key, status)
VALUES (?, ?, ?, ?);

-- name: GetTorrentByIdentity :one
SELECT * FROM torrents
WHERE info_hash = ? AND creator_public_key = ?
    LIMIT 1;

-- name: GetUserTorrents :many
SELECT
    v.info_hash,
    v.creator_public_key,
    v.creator_user_id,
    v.name,
    v.size_bytes,
    v.info_bytes,
    v.piece_length,
    v.is_private,
    v.added_at,
    ut.status,
    ut.progress,
    v.signature_id,
    v.signature_blob,
    v.payload_hash,
    v.signer_id,
    v.signer_public_key,
    v.signature_created_at
FROM v_torrent_with_signatures v
         JOIN user_torrents ut
              ON v.info_hash = ut.torrent_info_hash
                  AND v.creator_public_key = ut.creator_public_key
WHERE ut.user_id = ?;

-- name: UpdateUserTorrentStatus :exec
UPDATE user_torrents
SET status = ?
WHERE user_id = ? AND torrent_info_hash = ? AND creator_public_key = ?;

-- name: UpdateUserTorrentProgress :exec
UPDATE user_torrents
SET progress = ?
WHERE user_id = ? AND torrent_info_hash = ? AND creator_public_key = ?;

-- name: DeleteUserTorrent :execresult
DELETE FROM user_torrents
WHERE user_id = ? AND torrent_info_hash = ? AND creator_public_key = ?;