package domain

import (
	"time"

	"github.com/google/uuid"
)

// Finance represents financial information of a corporation
// Based on gBizINFO Finance CSV specification
type Finance struct {
	ID                           uuid.UUID `json:"id" db:"id"`
	CorporateNumber              string    `json:"corporate_number" db:"corporate_number"`
	CorporateNameFromNumber      string    `json:"corporate_name_from_number" db:"corporate_name_from_number"`
	HeadOfficeLocationFromNumber string    `json:"head_office_location_from_number" db:"head_office_location_from_number"`
	CorporateName                string    `json:"corporate_name" db:"corporate_name"`
	HeadOfficeLocation           string    `json:"head_office_location" db:"head_office_location"`
	AccountingStandards          string    `json:"accounting_standards" db:"accounting_standards"`
	BusinessYear                 string    `json:"business_year" db:"business_year"`
	PeriodNumber                 string    `json:"period_number" db:"period_number"`
	SalesRevenue                 string    `json:"sales_revenue" db:"sales_revenue"`
	SalesRevenueUnit             string    `json:"sales_revenue_unit" db:"sales_revenue_unit"`
	OperatingRevenue1            string    `json:"operating_revenue1" db:"operating_revenue1"`
	OperatingRevenue1Unit        string    `json:"operating_revenue1_unit" db:"operating_revenue1_unit"`
	OperatingRevenue2            string    `json:"operating_revenue2" db:"operating_revenue2"`
	OperatingRevenue2Unit        string    `json:"operating_revenue2_unit" db:"operating_revenue2_unit"`
	GrossOperatingRevenue        string    `json:"gross_operating_revenue" db:"gross_operating_revenue"`
	GrossOperatingRevenueUnit    string    `json:"gross_operating_revenue_unit" db:"gross_operating_revenue_unit"`
	OrdinaryRevenue              string    `json:"ordinary_revenue" db:"ordinary_revenue"`
	OrdinaryRevenueUnit          string    `json:"ordinary_revenue_unit" db:"ordinary_revenue_unit"`
	NetPremiumsWritten           string    `json:"net_premiums_written" db:"net_premiums_written"`
	NetPremiumsWrittenUnit       string    `json:"net_premiums_written_unit" db:"net_premiums_written_unit"`
	OrdinaryIncome               string    `json:"ordinary_income" db:"ordinary_income"`
	OrdinaryIncomeUnit           string    `json:"ordinary_income_unit" db:"ordinary_income_unit"`
	NetIncome                    string    `json:"net_income" db:"net_income"`
	NetIncomeUnit                string    `json:"net_income_unit" db:"net_income_unit"`
	CapitalStock                 string    `json:"capital_stock" db:"capital_stock"`
	CapitalStockUnit             string    `json:"capital_stock_unit" db:"capital_stock_unit"`
	NetAssets                    string    `json:"net_assets" db:"net_assets"`
	NetAssetsUnit                string    `json:"net_assets_unit" db:"net_assets_unit"`
	TotalAssets                  string    `json:"total_assets" db:"total_assets"`
	TotalAssetsUnit              string    `json:"total_assets_unit" db:"total_assets_unit"`
	NumberOfEmployees            string    `json:"number_of_employees" db:"number_of_employees"`
	NumberOfEmployeesUnit        string    `json:"number_of_employees_unit" db:"number_of_employees_unit"`
	MajorShareholder1            string    `json:"major_shareholder1" db:"major_shareholder1"`
	ShareholdingRatio1           string    `json:"shareholding_ratio1" db:"shareholding_ratio1"`
	MajorShareholder2            string    `json:"major_shareholder2" db:"major_shareholder2"`
	ShareholdingRatio2           string    `json:"shareholding_ratio2" db:"shareholding_ratio2"`
	MajorShareholder3            string    `json:"major_shareholder3" db:"major_shareholder3"`
	ShareholdingRatio3           string    `json:"shareholding_ratio3" db:"shareholding_ratio3"`
	MajorShareholder4            string    `json:"major_shareholder4" db:"major_shareholder4"`
	ShareholdingRatio4           string    `json:"shareholding_ratio4" db:"shareholding_ratio4"`
	MajorShareholder5            string    `json:"major_shareholder5" db:"major_shareholder5"`
	ShareholdingRatio5           string    `json:"shareholding_ratio5" db:"shareholding_ratio5"`
	CreatedAt                    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                    time.Time `json:"updated_at" db:"updated_at"`
}

// FinanceRepository defines the interface for finance data persistence
type FinanceRepository interface {
	// Create a new finance record
	Create(finance *Finance) error

	// CreateBatch creates multiple finance records efficiently
	CreateBatch(finances []*Finance) error

	// GetByCorporateNumber retrieves all finance records for a specific corporate number
	GetByCorporateNumber(corporateNumber string) ([]*Finance, error)

	// GetLatest retrieves the latest finance record for a specific corporate number
	GetLatest(corporateNumber string) (*Finance, error)

	// Count returns the total number of finance records
	Count() (int64, error)

	// DeleteAll removes all finance records
	DeleteAll() error

	// DeleteByCorporateNumber removes all finance records for a specific corporate number
	DeleteByCorporateNumber(corporateNumber string) error
}

// NewFinance creates a new Finance with a UUIDv7 primary key
func NewFinance(corporateNumber string) *Finance {
	now := time.Now()

	return &Finance{
		ID:              MustNewUUIDv7(), // Generate UUIDv7 for better DB performance
		CorporateNumber: corporateNumber,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}
