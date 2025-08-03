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
