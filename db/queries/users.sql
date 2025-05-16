-- name: CreateUser :execresult
INSERT INTO users (
  fullname, role, email, phone, agency_name, agency_address,
  agency_logo, wilaya, daira, password
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
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
WHERE id = ?;

-- name: UpdatePassword :exec
UPDATE users
SET
  password = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: UpdateLogo :exec
UPDATE users
SET
  agency_logo = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: CountUsersByEmail :one
SELECT COUNT(*) FROM users WHERE email = ?;

-- name: CountUsersByPhone :one
SELECT COUNT(*) FROM users WHERE phone = ?;
