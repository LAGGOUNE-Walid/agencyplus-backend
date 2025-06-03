-- name: CreateContact :execresult
INSERT INTO contacts (
  user_id,
  fullname,
  phone,
  email,
  wilaya,
  daira,
  client_type,
  searching_for,
  preferred_location_type,
  house_finishing,
  renting_floor_looking_for,
  is_married,
  min_budget,
  max_budget,
  created_at,
  updated_at
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
);

-- name: CountContactsByPhone :one
SELECT COUNT(*) FROM contacts WHERE phone = ?;

-- name: CountContactsByEmail :one
SELECT COUNT(*) FROM contacts WHERE email = ?;

-- name: GetContact :one
SELECT
  id,
  user_id,
  fullname,
  phone,
  email,
  wilaya,
  daira,
  client_type,
  searching_for,
  preferred_location_type,
  house_finishing,
  renting_floor_looking_for,
  is_married,
  min_budget,
  max_budget,
  created_at,
  updated_at
FROM contacts
WHERE id = ? AND user_id = ?;

-- name: GetAllContacts :many
SELECT
  id,
  user_id,
  fullname,
  phone,
  email,
  wilaya,
  daira,
  client_type,
  searching_for,
  preferred_location_type,
  house_finishing,
  renting_floor_looking_for,
  is_married,
  min_budget,
  max_budget,
  created_at,
  updated_at
FROM contacts
WHERE user_id = ?
ORDER BY id DESC;

-- name: DeleteContact :exec
DELETE FROM contacts
WHERE id = ? AND user_id = ?;

-- name: GetContactsById :many
SELECT
  id,
  user_id,
  phone
FROM contacts
WHERE user_id = ? AND id IN (sqlc.slice('ids'))
ORDER BY id DESC;