-- Updated for gBizINFO API compliance

-- name: GetCorporations :many
SELECT id, corporate_number, name, kana, name_en, search_name, postal_code, location, 
       prefecture_code, status, close_date, close_cause, representative_name, representative_position,
       date_of_establishment, founding_year, capital_stock, employee_number,
       company_size_male, company_size_female, business_items, business_summary,
       company_url, qualification_grade, number_of_activity, update_date,
       created_at, updated_at 
FROM corporations 
ORDER BY created_at DESC;

-- name: GetCorporationByID :one
SELECT id, corporate_number, name, kana, name_en, search_name, postal_code, location, 
       prefecture_code, status, close_date, close_cause, representative_name, representative_position,
       date_of_establishment, founding_year, capital_stock, employee_number,
       company_size_male, company_size_female, business_items, business_summary,
       company_url, qualification_grade, number_of_activity, update_date,
       created_at, updated_at 
FROM corporations 
WHERE id = $1 LIMIT 1;

-- name: GetCorporationByCorporateNumber :one
SELECT id, corporate_number, name, kana, name_en, search_name, postal_code, location, 
       prefecture_code, status, close_date, close_cause, representative_name, representative_position,
       date_of_establishment, founding_year, capital_stock, employee_number,
       company_size_male, company_size_female, business_items, business_summary,
       company_url, qualification_grade, number_of_activity, update_date,
       created_at, updated_at 
FROM corporations 
WHERE corporate_number = $1 LIMIT 1;

-- name: GetCorporationsWithFilter :many
SELECT id, corporate_number, name, kana, name_en, search_name, postal_code, location, 
       prefecture_code, status, close_date, close_cause, representative_name, representative_position,
       date_of_establishment, founding_year, capital_stock, employee_number,
       company_size_male, company_size_female, business_items, business_summary,
       company_url, qualification_grade, number_of_activity, update_date,
       created_at, updated_at 
FROM corporations 
WHERE ($1 = '' OR corporate_number = $1)
  AND ($2 = '' OR (
    search_name ILIKE '%' || $2 || '%' OR 
    kana ILIKE '%' || $2 || '%' OR 
    name_en ILIKE '%' || $2 || '%'
  ))
  AND ($3 = '' OR location ILIKE '%' || $3 || '%')
  AND ($4 = '' OR status = $4)
  AND ($5 = '' OR prefecture_code = $5)
ORDER BY created_at DESC
LIMIT $6 OFFSET $7;

-- name: CountCorporationsWithFilter :one
SELECT COUNT(*) 
FROM corporations 
WHERE ($1 = '' OR corporate_number = $1)
  AND ($2 = '' OR (
    search_name ILIKE '%' || $2 || '%' OR 
    kana ILIKE '%' || $2 || '%' OR 
    name_en ILIKE '%' || $2 || '%'
  ))
  AND ($3 = '' OR location ILIKE '%' || $3 || '%')
  AND ($4 = '' OR status = $4)
  AND ($5 = '' OR prefecture_code = $5);

-- name: CreateCorporation :one
INSERT INTO corporations (
    corporate_number, name, kana, name_en, search_name, postal_code, location, prefecture_code,
    status, close_date, close_cause, representative_name, representative_position,
    date_of_establishment, founding_year, capital_stock, employee_number,
    company_size_male, company_size_female, business_items, business_summary,
    company_url, qualification_grade, number_of_activity, update_date
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25
)
RETURNING id, corporate_number, name, kana, name_en, search_name, postal_code, location, 
          prefecture_code, status, close_date, close_cause, representative_name, representative_position,
          date_of_establishment, founding_year, capital_stock, employee_number,
          company_size_male, company_size_female, business_items, business_summary,
          company_url, qualification_grade, number_of_activity, update_date,
          created_at, updated_at;

-- name: UpdateCorporation :one
UPDATE corporations
SET name = $2, kana = $3, name_en = $4, search_name = $5, postal_code = $6, location = $7,
    prefecture_code = $8, status = $9, close_date = $10, close_cause = $11, representative_name = $12,
    representative_position = $13, date_of_establishment = $14, founding_year = $15,
    capital_stock = $16, employee_number = $17, company_size_male = $18,
    company_size_female = $19, business_items = $20, business_summary = $21,
    company_url = $22, qualification_grade = $23, number_of_activity = $24,
    update_date = $25, updated_at = NOW()
WHERE id = $1
RETURNING id, corporate_number, name, kana, name_en, search_name, postal_code, location, 
          prefecture_code, status, close_date, close_cause, representative_name, representative_position,
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
    corporate_number, name, kana, name_en, search_name, postal_code, location, prefecture_code,
    status, close_date, close_cause, representative_name, representative_position,
    date_of_establishment, founding_year, capital_stock, employee_number,
    company_size_male, company_size_female, business_items, business_summary,
    company_url, qualification_grade, number_of_activity, update_date
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25
)
ON CONFLICT (corporate_number)
DO UPDATE SET
    name = EXCLUDED.name,
    kana = EXCLUDED.kana,
    name_en = EXCLUDED.name_en,
    search_name = EXCLUDED.search_name,
    postal_code = EXCLUDED.postal_code,
    location = EXCLUDED.location,
    prefecture_code = EXCLUDED.prefecture_code,
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
RETURNING id, corporate_number, name, kana, name_en, search_name, postal_code, location, 
          prefecture_code, status, close_date, close_cause, representative_name, representative_position,
          date_of_establishment, founding_year, capital_stock, employee_number,
          company_size_male, company_size_female, business_items, business_summary,
          company_url, qualification_grade, number_of_activity, update_date,
          created_at, updated_at;
