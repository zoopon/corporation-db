-- Updated for gBizINFO API compliance

-- name: GetCorporations :many
SELECT id, corporate_number, name, kana, name_en, postal_code, location, 
       status, close_date, close_cause, representative_name, representative_position,
       date_of_establishment, founding_year, capital_stock, employee_number,
       company_size_male, company_size_female, business_items, business_summary,
       company_url, qualification_grade, number_of_activity, update_date,
       created_at, updated_at 
FROM corporations 
ORDER BY created_at DESC;

-- name: GetCorporationByID :one
SELECT id, corporate_number, name, kana, name_en, postal_code, location, 
       status, close_date, close_cause, representative_name, representative_position,
       date_of_establishment, founding_year, capital_stock, employee_number,
       company_size_male, company_size_female, business_items, business_summary,
       company_url, qualification_grade, number_of_activity, update_date,
       created_at, updated_at 
FROM corporations 
WHERE id = $1 LIMIT 1;

-- name: GetCorporationByCorporateNumber :one
SELECT id, corporate_number, name, kana, name_en, postal_code, location, 
       status, close_date, close_cause, representative_name, representative_position,
       date_of_establishment, founding_year, capital_stock, employee_number,
       company_size_male, company_size_female, business_items, business_summary,
       company_url, qualification_grade, number_of_activity, update_date,
       created_at, updated_at 
FROM corporations 
WHERE corporate_number = $1 LIMIT 1;

-- name: GetCorporationsWithFilter :many
SELECT id, corporate_number, name, kana, name_en, postal_code, location, 
       status, close_date, close_cause, representative_name, representative_position,
       date_of_establishment, founding_year, capital_stock, employee_number,
       company_size_male, company_size_female, business_items, business_summary,
       company_url, qualification_grade, number_of_activity, update_date,
       created_at, updated_at 
FROM corporations 
WHERE ($1 = '' OR corporate_number = $1)
  AND ($2 = '' OR name ILIKE '%' || $2 || '%')
  AND ($3 = '' OR location ILIKE '%' || $3 || '%')
  AND ($4 = '' OR status = $4)
ORDER BY created_at DESC
LIMIT $5 OFFSET $6;

-- name: CountCorporationsWithFilter :one
SELECT COUNT(*) 
FROM corporations 
WHERE ($1 = '' OR corporate_number = $1)
  AND ($2 = '' OR name ILIKE '%' || $2 || '%')
  AND ($3 = '' OR location ILIKE '%' || $3 || '%')
  AND ($4 = '' OR status = $4);

-- name: CreateCorporation :one
INSERT INTO corporations (
    corporate_number, name, kana, name_en, postal_code, location,
    status, close_date, close_cause, representative_name, representative_position,
    date_of_establishment, founding_year, capital_stock, employee_number,
    company_size_male, company_size_female, business_items, business_summary,
    company_url, qualification_grade, number_of_activity, update_date
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23
)
RETURNING id, corporate_number, name, kana, name_en, postal_code, location, 
          status, close_date, close_cause, representative_name, representative_position,
          date_of_establishment, founding_year, capital_stock, employee_number,
          company_size_male, company_size_female, business_items, business_summary,
          company_url, qualification_grade, number_of_activity, update_date,
          created_at, updated_at;

-- name: UpdateCorporation :one
UPDATE corporations
SET name = $2, kana = $3, name_en = $4, postal_code = $5, location = $6,
    status = $7, close_date = $8, close_cause = $9, representative_name = $10,
    representative_position = $11, date_of_establishment = $12, founding_year = $13,
    capital_stock = $14, employee_number = $15, company_size_male = $16,
    company_size_female = $17, business_items = $18, business_summary = $19,
    company_url = $20, qualification_grade = $21, number_of_activity = $22,
    update_date = $23, updated_at = NOW()
WHERE id = $1
RETURNING id, corporate_number, name, kana, name_en, postal_code, location, 
          status, close_date, close_cause, representative_name, representative_position,
          date_of_establishment, founding_year, capital_stock, employee_number,
          company_size_male, company_size_female, business_items, business_summary,
          company_url, qualification_grade, number_of_activity, update_date,
          created_at, updated_at;

-- name: DeleteCorporation :exec
DELETE FROM corporations
WHERE id = $1;

-- name: DeleteCorporationByCorporateNumber :exec
DELETE FROM corporations
WHERE corporate_number = $1;

-- name: UpsertCorporation :one
INSERT INTO corporations (
    corporate_number, name, kana, name_en, postal_code, location,
    status, close_date, close_cause, representative_name, representative_position,
    date_of_establishment, founding_year, capital_stock, employee_number,
    company_size_male, company_size_female, business_items, business_summary,
    company_url, qualification_grade, number_of_activity, update_date
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23
)
ON CONFLICT (corporate_number)
DO UPDATE SET
    name = EXCLUDED.name,
    kana = EXCLUDED.kana,
    name_en = EXCLUDED.name_en,
    postal_code = EXCLUDED.postal_code,
    location = EXCLUDED.location,
    status = EXCLUDED.status,
    close_date = EXCLUDED.close_date,
    close_cause = EXCLUDED.close_cause,
    representative_name = EXCLUDED.representative_name,
    representative_position = EXCLUDED.representative_position,
    date_of_establishment = EXCLUDED.date_of_establishment,
    founding_year = EXCLUDED.founding_year,
    capital_stock = EXCLUDED.capital_stock,
    employee_number = EXCLUDED.employee_number,
    company_size_male = EXCLUDED.company_size_male,
    company_size_female = EXCLUDED.company_size_female,
    business_items = EXCLUDED.business_items,
    business_summary = EXCLUDED.business_summary,
    company_url = EXCLUDED.company_url,
    qualification_grade = EXCLUDED.qualification_grade,
    number_of_activity = EXCLUDED.number_of_activity,
    update_date = EXCLUDED.update_date,
    updated_at = NOW()
RETURNING id, corporate_number, name, kana, name_en, postal_code, location, 
          status, close_date, close_cause, representative_name, representative_position,
          date_of_establishment, founding_year, capital_stock, employee_number,
          company_size_male, company_size_female, business_items, business_summary,
          company_url, qualification_grade, number_of_activity, update_date,
          created_at, updated_at;
