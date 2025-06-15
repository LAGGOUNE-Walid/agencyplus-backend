-- name: CreateCalendar :one
INSERT INTO calendar_events (
    user_id, title, content, for_date
) VALUES (
    ?, ?, ?, ?
)
RETURNING *;

-- name: UpdateCalendar :exec
UPDATE calendar_events SET
title = ?,
content = ?,
for_date = ?,
updated_at = CURRENT_TIMESTAMP
WHERE id = ? and user_id = ?;

-- name: DeleteCalendar :exec
UPDATE calendar_events set deleted_at = CURRENT_TIMESTAMP where id = ? and user_id = ?;

-- name: GetUserCalendars :many
SELECT * from calendar_events where user_id = ? and deleted_at is null;