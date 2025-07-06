-- name: CreateBuilding :execresult
INSERT INTO buildings (
  user_id, location, title, wilaya, daira, building_type, is_promotion_building, is_residency,
  status, price, surface_total, surface_built, rooms, bathrooms, floors_total, parking_spaces,
  is_by_the_sea, has_water, has_electricity, has_gas, has_internet, has_garden, has_pool,
  has_elevator, has_central_heating, has_water_tank, has_air_conditioner, has_equipped_kitchen,
  has_terrace, has_notarial_deed, has_land_booklet, has_act_in_joint_ownership,
  has_certificate_of_conformity, has_decision, has_concession, has_stamped_paper,
  has_building_permit, has_off_plan_sales_contract, building_finished_type,
  acceptable_payment_type, furnished, year_built, description, shareable_link
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: CreateBuildingImage :exec
INSERT INTO building_images (user_id, building_id, path, mimetype, size)
VALUES (?, ?, ?, ?, ?);

-- name: DeleteBuilding :exec
UPDATE buildings SET deleted_at = CURRENT_TIMESTAMP
WHERE user_id = ? AND id = ? AND deleted_at is NULL;

-- name: DeleteBuildingImage :exec
UPDATE building_images SET deleted_at = CURRENT_TIMESTAMP
WHERE building_id = ? AND user_id = ? AND id = ? AND deleted_at is NULL;

-- name: DeleteBuildingImages :exec
UPDATE building_images SET deleted_at = CURRENT_TIMESTAMP
WHERE building_id = ? AND user_id = ? AND deleted_at is NULL;

-- name: CreateBuildingDocument :exec
INSERT INTO building_documents (user_id, building_id, path, mimetype, size, thumbnail)
VALUES (?, ?, ?, ?, ?, ?);

-- name: DeleteBuildingDocument :exec
UPDATE building_documents SET deleted_at = CURRENT_TIMESTAMP
WHERE building_id = ? AND user_id = ? AND id = ? AND deleted_at is NULL;

-- name: DeleteBuildingDocuments :exec
UPDATE building_documents SET deleted_at = CURRENT_TIMESTAMP
WHERE building_id = ? AND user_id = ? AND deleted_at is NULL;

-- name: ListPaginatedBuildings :many
SELECT * FROM buildings
WHERE user_id = ? AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: ListImagesForBuildingIDs :many
SELECT * FROM building_images
WHERE building_id IN (sqlc.slice('building_ids')) AND deleted_at is NULL;

-- name: ListDocumentsForBuildingIDs :many
SELECT * FROM building_documents
WHERE building_id IN (sqlc.slice('building_ids')) AND deleted_at is NULL;

-- name: GetBuilding :one
SELECT * FROM buildings
WHERE user_id = ? AND id = ? AND deleted_at is NULL;

-- name: UpdateBuilding :exec
UPDATE buildings
SET
  title = ?,
  status = ?,
  location = ?,
  wilaya = ?,
  daira = ?,
  building_type = ?,
  is_promotion_building = ?,
  is_residency = ?,
  price = ?,
  surface_total = ?,
  surface_built = ?,
  rooms = ?,
  bathrooms = ?,
  floors_total = ?,
  parking_spaces = ?,
  is_by_the_sea = ?,
  has_water = ?,
  has_electricity = ?,
  has_gas = ?,
  has_internet = ?,
  has_garden = ?,
  has_pool = ?,
  has_elevator = ?,
  has_central_heating = ?,
  has_water_tank = ?,
  has_air_conditioner = ?,
  has_equipped_kitchen = ?,
  has_terrace = ?,
  has_notarial_deed = ?,
  has_land_booklet = ?,
  has_act_in_joint_ownership = ?,
  has_certificate_of_conformity = ?,
  has_decision = ?,
  has_concession = ?,
  has_stamped_paper = ?,
  has_building_permit = ?,
  has_off_plan_sales_contract = ?,
  building_finished_type = ?,
  acceptable_payment_type = ?,
  furnished = ?,
  year_built = ?,
  description = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE id = ? AND user_id = ? AND deleted_at is NULL;

-- name: CreateBuildingVue :exec
INSERT INTO building_vues(building_id, ip_address, user_agent)
VALUES(?, ?, ?);

-- name: InsertEmbeddings :exec
INSERT INTO building_embeddings(building_id, embedding, created_at) VALUES (?, ?, CURRENT_TIMESTAMP);

-- name: GetBuildingEmbeddings :one
SELECT * FROM building_embeddings where building_id = ?