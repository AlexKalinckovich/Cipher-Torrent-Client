-- name: CreateTorrent :execresult
INSERT INTO torrents (info_hash, info_bytes, name, size_bytes, piece_length, is_private, storage_path, added_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetTorrentByInfoHash :one
SELECT * FROM torrents WHERE info_hash = ? LIMIT 1;

-- name: CreateUserTorrent :execresult
INSERT INTO user_torrents (user_id, torrent_info_hash, status, progress)
VALUES (?, ?, ?, ?);

-- name: GetUserTorrent :one
SELECT * FROM user_torrents WHERE user_id = ? AND torrent_info_hash = ? LIMIT 1;

-- name: GetUserTorrents :many
SELECT t.*, ut.status, ut.progress
FROM torrents t
         JOIN user_torrents ut ON t.info_hash = ut.torrent_info_hash
WHERE ut.user_id = ?;

-- name: UpdateUserTorrentStatus :exec
UPDATE user_torrents SET status = ? WHERE user_id = ? AND torrent_info_hash = ?;

-- name: UpdateUserTorrentProgress :exec
UPDATE user_torrents SET progress = ? WHERE user_id = ? AND torrent_info_hash = ?;

-- name: DeleteTorrent :execresult
DELETE FROM torrents WHERE info_hash = ?;