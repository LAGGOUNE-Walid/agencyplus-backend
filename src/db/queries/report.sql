-- name: CreateReport :one
INSERT INTO reports (
    user_id, title, content
) VALUES (
    ?, ?, ?
)
RETURNING *;

-- name: GetUserReportById :one
SELECT * from reports where user_id = ? and id = ?;

-- name: GetUserReportForMasterById :one
SELECT * from reports where user_id IN (sqlc.slice('users_id')) and id = ?;

-- name: UpdateReport :exec
UPDATE reports SET
title = ?,
content = ?,
updated_at = CURRENT_TIMESTAMP
WHERE id = ? and user_id = ?;

-- name: UpdateReportByMaster :exec
UPDATE reports SET
title = ?,
content = ?,
updated_at = CURRENT_TIMESTAMP
WHERE id = ? and user_id IN (sqlc.slice('users_id'));

-- name: DeleteReport :exec
UPDATE reports set deleted_at = CURRENT_TIMESTAMP where id = ? and user_id = ?;

-- name: DeleteReportByMaster :exec
UPDATE reports set deleted_at = CURRENT_TIMESTAMP where id = ? and user_id IN (sqlc.slice('users_id'));

-- name: GetUserReports :many
SELECT * from reports where user_id = ? and deleted_at is null ORDER BY id desc;

-- name: GetUserMasterReports :many
SELECT * from reports where user_id IN (sqlc.slice('users_id')) and deleted_at is null ORDER BY id desc;