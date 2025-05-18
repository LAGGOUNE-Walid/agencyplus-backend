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
