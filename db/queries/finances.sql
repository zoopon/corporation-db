-- name: CreateFinance :exec
INSERT INTO finances (
    corporate_number,
    corporate_name_from_number,
    head_office_location_from_number,
    corporate_name,
    head_office_location,
    accounting_standards,
    business_year,
    period_number,
    sales_revenue,
    sales_revenue_unit,
    operating_revenue1,
    operating_revenue1_unit,
    operating_revenue2,
    operating_revenue2_unit,
    gross_operating_revenue,
    gross_operating_revenue_unit,
    ordinary_revenue,
    ordinary_revenue_unit,
    net_premiums_written,
    net_premiums_written_unit,
    ordinary_income,
    ordinary_income_unit,
    net_income,
    net_income_unit,
    capital_stock,
    capital_stock_unit,
    net_assets,
    net_assets_unit,
    total_assets,
    total_assets_unit,
    number_of_employees,
    number_of_employees_unit,
    major_shareholder1,
    shareholding_ratio1,
    major_shareholder2,
    shareholding_ratio2,
    major_shareholder3,
    shareholding_ratio3,
    major_shareholder4,
    shareholding_ratio4,
    major_shareholder5,
    shareholding_ratio5
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
    $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42
);

-- name: GetFinancesByCorporateNumber :many
SELECT * FROM finances 
WHERE corporate_number = $1
ORDER BY business_year DESC, period_number DESC;

-- name: GetLatestFinanceByCorporateNumber :one
SELECT * FROM finances 
WHERE corporate_number = $1 
ORDER BY business_year DESC, period_number DESC 
LIMIT 1;

-- name: CountAllFinances :one
SELECT COUNT(*) FROM finances;

-- name: DeleteAllFinances :exec
DELETE FROM finances;

-- name: GetFinancesByBusinessYear :many
SELECT * FROM finances 
WHERE business_year = $1
ORDER BY corporate_number;

-- name: DeleteFinancesByCorporateNumber :exec
DELETE FROM finances 
WHERE corporate_number = $1;

-- name: CountFinances :one
SELECT COUNT(*) FROM finances;
