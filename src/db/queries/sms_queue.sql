-- name: CreateSMSQueue :one
INSERT INTO sms_queues (
    user_id, title, content, from_number, total_recipients, scheduled_at
) VALUES (
    ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: AddSMSQueueContact :one
INSERT INTO sms_queue_contacts (
    sms_queue_id, phone_number
) VALUES (
    ?, ?
)
RETURNING *;

-- name: GetSmsQueue :one
SELECT * from sms_queues where id = ? LIMIT 1;

-- name: GetSmsContacts :many
SELECT * from sms_queue_contacts where sms_queue_id = ?