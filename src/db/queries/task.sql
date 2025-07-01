-- name: CreateTask :one
INSERT INTO tasks (
    to_id, title, content, due_date, root_id
) VALUES (
    ?, ?, ?, ?, ?
)
RETURNING *;

-- name: GetCurrentUserTasks :many
SELECT * from tasks where to_id = ?;

-- name: GetRootUserCreatedTasks :many
SELECT * from tasks where root_id = ? OR to_id = ?;

-- name: MarkTaskAsDone :exec
UPDATE tasks
SET 
is_completed = 1,
updated_at = CURRENT_TIMESTAMP
WHERE id = ? AND (to_id = ? OR root_id = ?)