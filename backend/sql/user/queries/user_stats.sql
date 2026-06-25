-- name: CreateUserStats :exec
INSERT INTO user_stats (user_id, total_uploaded_bytes, total_downloaded_bytes, reputation_score, signed_torrents_count, active_torrents_count, peers_trusted_count)
VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: GetUserStatsByUserID :one
SELECT user_id, total_uploaded_bytes, total_downloaded_bytes, reputation_score, signed_torrents_count, active_torrents_count, peers_trusted_count
FROM user_stats
WHERE user_id = ? LIMIT 1;

-- name: UpdateUserStats :exec
UPDATE user_stats
SET total_uploaded_bytes = ?, total_downloaded_bytes = ?, reputation_score = ?, signed_torrents_count = ?, active_torrents_count = ?, peers_trusted_count = ?
WHERE user_id = ?;

-- name: DeleteUserStats :exec
DELETE FROM user_stats
WHERE user_id = ?;