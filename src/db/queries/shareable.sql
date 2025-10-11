-- name: CreateShareable :one
INSERT INTO shareables (token, model_type, model_id, user_id)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: GetShareable :one
SELECT * FROM shareables 
WHERE token = ?;

-- name: GetShareableByModelId :one
SELECT * from shareables where model_type = ? and model_id = ? LIMIT 1;

-- name: GetShareableWithModel :one
SELECT s.token, s.model_type, s.model_id, 
    b.*,
    d.*
FROM shareables s
LEFT JOIN buildings b ON s.model_type = 'building' AND s.model_id = b.id
LEFT JOIN building_documents d ON s.model_type = 'document' AND s.model_id = d.id
WHERE s.token = ?;