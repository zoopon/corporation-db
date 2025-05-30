-- name: GetCorporations :many
SELECT id, corporate_number, name, name_kana, english_name, postal_code, address, 
       prefecture_code, city_code, founding_date, status, corporate_type, 
       capital_stock, employee_number, representative, business_description, 
       industry, website, phone, email, last_updated, created_at, updated_at 
FROM corporations 
ORDER BY created_at DESC;

-- name: GetCorporationByID :one
SELECT id, corporate_number, name, name_kana, english_name, postal_code, address, 
       prefecture_code, city_code, founding_date, status, corporate_type, 
       capital_stock, employee_number, representative, business_description, 
       industry, website, phone, email, last_updated, created_at, updated_at 
FROM corporations 
WHERE id = $1 LIMIT 1;

-- name: GetCorporationByCorporateNumber :one
SELECT id, corporate_number, name, name_kana, english_name, postal_code, address, 
       prefecture_code, city_code, founding_date, status, corporate_type, 
       capital_stock, employee_number, representative, business_description, 
       industry, website, phone, email, last_updated, created_at, updated_at 
FROM corporations 
WHERE corporate_number = $1 LIMIT 1;

-- name: GetCorporationsWithFilter :many
SELECT id, corporate_number, name, name_kana, english_name, postal_code, address, 
       prefecture_code, city_code, founding_date, status, corporate_type, 
       capital_stock, employee_number, representative, business_description, 
       industry, website, phone, email, last_updated, created_at, updated_at 
FROM corporations 
WHERE ($1 = '' OR corporate_number = $1)
  AND ($2 = '' OR name ILIKE '%' || $2 || '%')
  AND ($3 = '' OR prefecture_code = $3)
  AND ($4 = '' OR status = $4)
  AND ($5 = '' OR corporate_type = $5)
ORDER BY created_at DESC
LIMIT $6 OFFSET $7;

-- name: CountCorporationsWithFilter :one
SELECT COUNT(*) 
FROM corporations 
WHERE ($1 = '' OR corporate_number = $1)
  AND ($2 = '' OR name ILIKE '%' || $2 || '%')
  AND ($3 = '' OR prefecture_code = $3)
  AND ($4 = '' OR status = $4)
  AND ($5 = '' OR corporate_type = $5);

-- name: CreateCorporation :one
INSERT INTO corporations (
    corporate_number, name, name_kana, english_name, postal_code, address,
    prefecture_code, city_code, founding_date, status, corporate_type,
    capital_stock, employee_number, representative, business_description,
    industry, website, phone, email, last_updated
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20
)
RETURNING id, corporate_number, name, name_kana, english_name, postal_code, address, 
          prefecture_code, city_code, founding_date, status, corporate_type, 
          capital_stock, employee_number, representative, business_description, 
          industry, website, phone, email, last_updated, created_at, updated_at;

-- name: UpdateCorporation :one
UPDATE corporations
SET name = $2, name_kana = $3, english_name = $4, postal_code = $5, address = $6,
    prefecture_code = $7, city_code = $8, founding_date = $9, status = $10, 
    corporate_type = $11, capital_stock = $12, employee_number = $13, 
    representative = $14, business_description = $15, industry = $16, 
    website = $17, phone = $18, email = $19, last_updated = $20, updated_at = NOW()
WHERE id = $1
RETURNING id, corporate_number, name, name_kana, english_name, postal_code, address, 
          prefecture_code, city_code, founding_date, status, corporate_type, 
          capital_stock, employee_number, representative, business_description, 
          industry, website, phone, email, last_updated, created_at, updated_at;

-- name: DeleteCorporation :exec
DELETE FROM corporations
WHERE id = $1;

-- name: DeleteCorporationByCorporateNumber :exec
DELETE FROM corporations
WHERE corporate_number = $1;

-- name: UpsertCorporation :one
INSERT INTO corporations (
    corporate_number, name, name_kana, english_name, postal_code, address,
    prefecture_code, city_code, founding_date, status, corporate_type,
    capital_stock, employee_number, representative, business_description,
    industry, website, phone, email, last_updated
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20
)
ON CONFLICT (corporate_number)
DO UPDATE SET
    name = EXCLUDED.name,
    name_kana = EXCLUDED.name_kana,
    english_name = EXCLUDED.english_name,
    postal_code = EXCLUDED.postal_code,
    address = EXCLUDED.address,
    prefecture_code = EXCLUDED.prefecture_code,
    city_code = EXCLUDED.city_code,
    founding_date = EXCLUDED.founding_date,
    status = EXCLUDED.status,
    corporate_type = EXCLUDED.corporate_type,
    capital_stock = EXCLUDED.capital_stock,
    employee_number = EXCLUDED.employee_number,
    representative = EXCLUDED.representative,
    business_description = EXCLUDED.business_description,
    industry = EXCLUDED.industry,
    website = EXCLUDED.website,
    phone = EXCLUDED.phone,
    email = EXCLUDED.email,
    last_updated = EXCLUDED.last_updated,
    updated_at = NOW()
RETURNING id, corporate_number, name, name_kana, english_name, postal_code, address, 
          prefecture_code, city_code, founding_date, status, corporate_type, 
          capital_stock, employee_number, representative, business_description, 
          industry, website, phone, email, last_updated, created_at, updated_at;
