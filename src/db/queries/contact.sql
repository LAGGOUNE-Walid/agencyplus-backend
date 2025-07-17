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
  preferred_building_types,
  preferred_features,
  min_rooms,
  max_rooms,
  min_surface,
  max_surface,
  furnished,
  acceptable_payment_type,
  max_year_built,
  created_at,
  updated_at
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
);

-- name: CountContactsByPhone :one
SELECT COUNT(*) FROM contacts WHERE phone = ?;

-- name: CountContactsByEmail :one
SELECT COUNT(*) FROM contacts WHERE email = ?;

-- name: CountUserContacts :one
SELECT COUNT(*) FROM contacts WHERE user_id = ?;

-- name: GetContact :one
SELECT
  *
FROM contacts
WHERE id = ? AND user_id = ?;

-- name: GetContactsList :many
SELECT
  *
FROM contacts
WHERE id IN (sqlc.slice('contact_ids')) AND user_id = ?;

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

-- name: InsertContactEmbeddings :exec
INSERT INTO contact_embeddings(contact_id, embedding, created_at) VALUES (?, ?, CURRENT_TIMESTAMP);

-- name: GetContactEmbeddings :one
SELECT * FROM contact_embeddings where contact_id = ?;

-- name: GetContactsWithEmbeddings :many
SELECT
  contacts.id,
  contact_embeddings.embedding
FROM contacts
RIGHT JOIN contact_embeddings ON contact_embeddings.contact_id = contacts.id
WHERE user_id = ?
ORDER BY contacts.id DESC;