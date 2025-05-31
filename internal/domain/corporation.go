package domain

import (
	"errors"
	"time"
)

// Corporation errors
var (
	ErrCorporationNotFound      = errors.New("corporation not found")
	ErrInvalidCorporateNumber   = errors.New("invalid corporate number")
	ErrCorporationAlreadyExists = errors.New("corporation already exists")
)

// Corporation represents a corporation entity based on gBizINFO REST API specification
type Corporation struct {
	ID int64 `json:"id"`

	// Basic Information (基本情報) - matches gBizINFO API response
	CorporateNumber string  `json:"corporate_number"`      // corporate_number (法人番号)
	Name            string  `json:"name"`                  // name (法人名)
	Kana            *string `json:"kana,omitempty"`        // kana (法人名ふりがな)
	NameEn          *string `json:"name_en,omitempty"`     // name_en (英語法人名)
	PostalCode      *string `json:"postal_code,omitempty"` // postal_code (郵便番号)
	Location        *string `json:"location,omitempty"`    // location (所在地)
	Status          string  `json:"status"`                // status (法人状態)

	// Registration Information (登記情報)
	CloseDate  *time.Time `json:"close_date,omitempty"`  // close_date (登記記録の閉鎖等年月日)
	CloseCause *string    `json:"close_cause,omitempty"` // close_cause (登記記録の閉鎖等の事由)

	// Representative Information (代表者情報)
	RepresentativeName     *string `json:"representative_name,omitempty"`     // representative_name
	RepresentativePosition *string `json:"representative_position,omitempty"` // representative_position

	// Company Details (企業詳細)
	DateOfEstablishment *time.Time `json:"date_of_establishment,omitempty"` // date_of_establishment
	FoundingYear        *int32     `json:"founding_year,omitempty"`         // founding_year
	CapitalStock        *int64     `json:"capital_stock,omitempty"`         // capital_stock
	EmployeeNumber      *int32     `json:"employee_number,omitempty"`       // employee_number
	CompanySizeMale     *int32     `json:"company_size_male,omitempty"`     // company_size_male
	CompanySizeFemale   *int32     `json:"company_size_female,omitempty"`   // company_size_female

	// Business Information (事業情報)
	BusinessItems      *string `json:"business_items,omitempty"`      // business_items (JSON string)
	BusinessSummary    *string `json:"business_summary,omitempty"`    // business_summary
	CompanyUrl         *string `json:"company_url,omitempty"`         // company_url
	QualificationGrade *string `json:"qualification_grade,omitempty"` // qualification_grade
	NumberOfActivity   *string `json:"number_of_activity,omitempty"`  // number_of_activity

	// gBizINFO Metadata
	UpdateDate *time.Time `json:"update_date,omitempty"` // update_date

	// Database Metadata
	CreatedAt time.Time `json:"created_at"` // DB登録日時
	UpdatedAt time.Time `json:"updated_at"` // DB更新日時
}

// CreateCorporationRequest represents request to create a corporation (based on gBizINFO API)
type CreateCorporationRequest struct {
	CorporateNumber        string     `json:"corporate_number" validate:"required,len=13"`
	Name                   string     `json:"name" validate:"required"`
	Kana                   *string    `json:"kana,omitempty"`
	NameEn                 *string    `json:"name_en,omitempty"`
	PostalCode             *string    `json:"postal_code,omitempty"`
	Location               *string    `json:"location,omitempty"`
	Status                 string     `json:"status" validate:"required"`
	CloseDate              *time.Time `json:"close_date,omitempty"`
	CloseCause             *string    `json:"close_cause,omitempty"`
	RepresentativeName     *string    `json:"representative_name,omitempty"`
	RepresentativePosition *string    `json:"representative_position,omitempty"`
	DateOfEstablishment    *time.Time `json:"date_of_establishment,omitempty"`
	FoundingYear           *int32     `json:"founding_year,omitempty"`
	CapitalStock           *int64     `json:"capital_stock,omitempty"`
	EmployeeNumber         *int32     `json:"employee_number,omitempty"`
	CompanySizeMale        *int32     `json:"company_size_male,omitempty"`
	CompanySizeFemale      *int32     `json:"company_size_female,omitempty"`
	BusinessItems          *string    `json:"business_items,omitempty"`
	BusinessSummary        *string    `json:"business_summary,omitempty"`
	CompanyUrl             *string    `json:"company_url,omitempty"`
	QualificationGrade     *string    `json:"qualification_grade,omitempty"`
	NumberOfActivity       *string    `json:"number_of_activity,omitempty"`
	UpdateDate             *time.Time `json:"update_date,omitempty"`
}

// CorporationFilter represents filtering options for corporation queries
type CorporationFilter struct {
	CorporateNumber *string `json:"corporate_number,omitempty"`
	Name            *string `json:"name,omitempty"`
	Prefecture      *string `json:"prefecture,omitempty"`
	Status          *string `json:"status,omitempty"`
	CorporateType   *string `json:"corporate_type,omitempty"`
	Limit           int     `json:"limit"`
	Offset          int     `json:"offset"`
}
