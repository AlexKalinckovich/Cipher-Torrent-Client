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

-- name: ListAllTorrents :many
SELECT
    t.info_hash,
    t.creator_public_key,
    t.creator_user_id,
    t.name,
    t.size_bytes,
    t.info_bytes,
    t.piece_length,
    t.is_private,
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