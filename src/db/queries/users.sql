-- name: CreateUser :execresult
INSERT INTO users (
  fullname, role, root_id, email, phone, agency_name, agency_address,
  agency_logo, wilaya, daira, password
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: UpdateUser :exec
UPDATE users
SET
  fullname = ?,
  email = ?,
  phone = ?,
  agency_name = ?,
  agency_address = ?,
  wilaya = ?,
  daira = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE id = ? AND deleted_at is NULL;

-- name: UpdatePassword :exec
UPDATE users
SET
  password = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE id = ? AND deleted_at is NULL;

-- name: UpdateLogo :exec
UPDATE users
SET
  agency_logo = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE id = ? AND deleted_at is NULL;

-- name: CountUsersByEmail :one
SELECT COUNT(*) FROM users WHERE email = ?;

-- name: CountUsersByPhone :one
SELECT COUNT(*) FROM users WHERE phone = ?;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ? AND deleted_at is NULL LIMIT 1;

-- name: GetUserById :one
SELECT id, fullname, role, root_id, email, phone, agency_name, agency_address,
  agency_logo, wilaya, daira FROM users WHERE id = ? AND deleted_at is NULL LIMIT 1;


-- name: CountUsersByEmailExcludingID :one
SELECT COUNT(*) FROM users
WHERE email = ? AND id != ?;

-- name: GetUserAgents :many
SELECT id, fullname, role, root_id, email, phone, agency_name, agency_address,
  agency_logo, wilaya, daira from users where root_id = ?;

-- name: GetAgencyUsers :many
SELECT id, fullname, role, root_id, email, phone, agency_name, agency_address,
  agency_logo, wilaya, daira from users where root_id = ?;

-- name: GetUserSubscriptions :many
SELECT * from user_subscriptions where user_id = ? ORDER BY id DESC;

-- name: GetCurrentUserSubscription :one
SELECT * FROM user_subscriptions where user_id = ? ORDER by id desc limit 1;

-- name: CreateUsersubscription :exec
INSERT INTO user_subscriptions (
  user_id, plan_id, status, current_period_start, current_period_end, next_billing_date,
  trial_start, trial_end, amount
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: UpdateUserSubscription :exec
UPDATE user_subscriptions 
SET 
    plan_id = ?,
    status = ?,
    current_period_start = ?,
    current_period_end = ?,
    next_billing_date = ?,
    trial_start = ?,
    trial_end = ?,
    amount = ?,
    updated_at = datetime('now')
WHERE id = ?;

-- name: GetUser :one
SELECT * FROM users where id = ?;

-- name: GetUsers :many
SELECT * FROM users where id IN(sqlc.slice('users_id'));

-- name: ForceDelete :exec
DELETE FROM users where id = ?;

-- name: UpdateAgencyUsersSubscriptionStatus :exec
UPDATE user_subscriptions SET status = ? WHERE user_id IN(sqlc.slice('users_id'));