-- db/queries/bases.sql
-- Base/Branch management queries

-- name: CreateBase :one
INSERT INTO bases (
    id,
    corporation_id,
    corporate_number,
    base_name,
    country_code,
    postal_code,
    location,
    phone_number,
    fax_number,
    data_obtained_at,
    data_source_url,
    is_head_office
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
)
RETURNING *;

-- name: CreateBaseBatch :exec
INSERT INTO bases (
    id,
    corporation_id,
    corporate_number,
    base_name,
    country_code,
    postal_code,
    location,
    phone_number,
    fax_number,
    data_obtained_at,
    data_source_url,
    is_head_office
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);

-- name: GetBaseByID :one
SELECT * FROM bases
WHERE id = $1;

-- name: GetBasesByCorporationID :many
SELECT * FROM bases
WHERE corporation_id = $1
ORDER BY is_head_office DESC, base_name ASC;

-- name: GetBasesByCorporateNumber :many
SELECT * FROM bases
WHERE corporate_number = $1
ORDER BY is_head_office DESC, base_name ASC;

-- name: UpdateBase :one
UPDATE bases
SET 
    base_name = $2,
    country_code = $3,
    postal_code = $4,
    location = $5,
    phone_number = $6,
    fax_number = $7,
    data_obtained_at = $8,
    data_source_url = $9,
    is_head_office = $10,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteBase :exec
DELETE FROM bases
WHERE id = $1;

-- name: DeleteBasesByCorporationID :exec
DELETE FROM bases
WHERE corporation_id = $1;

-- name: GetHeadOfficesByCorporateNumbers :many
SELECT * FROM bases
WHERE corporate_number = ANY($1::char(13)[])
AND is_head_office = true
ORDER BY corporate_number;

-- name: CountBasesByCorporateNumber :one
SELECT COUNT(*) FROM bases
WHERE corporate_number = $1;

-- name: ListAllBases :many
SELECT * FROM bases
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: SearchBasesByName :many
SELECT * FROM bases
WHERE base_name ILIKE '%' || $1 || '%'
ORDER BY base_name ASC
LIMIT $2 OFFSET $3;

-- name: GetBasesByCountry :many
SELECT * FROM bases
WHERE country_code = $1
ORDER BY base_name ASC
LIMIT $2 OFFSET $3;

-- name: GetHeadOfficeByCorporateNumber :one
SELECT * FROM bases
WHERE corporate_number = $1 AND is_head_office = true
LIMIT 1;
