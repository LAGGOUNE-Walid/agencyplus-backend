-- name: CreatePayment :exec
INSERT INTO payments (
    user_id, payload
) VALUES (
    ?, ?
);