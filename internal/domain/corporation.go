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

// Corporation represents a corporation entity from gBizINFO
type Corporation struct {
	ID                  int64      `json:"id"`
	CorporateNumber     string     `json:"corporate_number"`               // 法人番号（13桁）
	Name                string     `json:"name"`                           // 法人名
	NameKana            *string    `json:"name_kana,omitempty"`            // 法人名（カナ）
	EnglishName         *string    `json:"english_name,omitempty"`         // 英語法人名
	PostalCode          *string    `json:"postal_code,omitempty"`          // 郵便番号
	Address             *string    `json:"address,omitempty"`              // 所在地
	PrefectureCode      *string    `json:"prefecture_code,omitempty"`      // 都道府県コード
	CityCode            *string    `json:"city_code,omitempty"`            // 市区町村コード
	FoundingDate        *time.Time `json:"founding_date,omitempty"`        // 設立年月日
	Status              string     `json:"status"`                         // 法人状態（活動中/解散等）
	CorporateType       *string    `json:"corporate_type,omitempty"`       // 法人種別
	CapitalStock        *int64     `json:"capital_stock,omitempty"`        // 資本金
	EmployeeNumber      *int32     `json:"employee_number,omitempty"`      // 従業員数
	Representative      *string    `json:"representative,omitempty"`       // 代表者名
	BusinessDescription *string    `json:"business_description,omitempty"` // 事業内容
	Industry            *string    `json:"industry,omitempty"`             // 業種
	Website             *string    `json:"website,omitempty"`              // ウェブサイト
	Phone               *string    `json:"phone,omitempty"`                // 電話番号
	Email               *string    `json:"email,omitempty"`                // メールアドレス
	LastUpdated         *time.Time `json:"last_updated,omitempty"`         // 最終更新日時（gBizINFO側）
	CreatedAt           time.Time  `json:"created_at"`                     // DB登録日時
	UpdatedAt           time.Time  `json:"updated_at"`                     // DB更新日時
}

// CreateCorporationRequest represents request to create a corporation
type CreateCorporationRequest struct {
	CorporateNumber     string     `json:"corporate_number" validate:"required,len=13"`
	Name                string     `json:"name" validate:"required"`
	NameKana            *string    `json:"name_kana,omitempty"`
	EnglishName         *string    `json:"english_name,omitempty"`
	PostalCode          *string    `json:"postal_code,omitempty"`
	Address             *string    `json:"address,omitempty"`
	PrefectureCode      *string    `json:"prefecture_code,omitempty"`
	CityCode            *string    `json:"city_code,omitempty"`
	FoundingDate        *time.Time `json:"founding_date,omitempty"`
	Status              string     `json:"status" validate:"required"`
	CorporateType       *string    `json:"corporate_type,omitempty"`
	CapitalStock        *int64     `json:"capital_stock,omitempty"`
	EmployeeNumber      *int32     `json:"employee_number,omitempty"`
	Representative      *string    `json:"representative,omitempty"`
	BusinessDescription *string    `json:"business_description,omitempty"`
	Industry            *string    `json:"industry,omitempty"`
	Website             *string    `json:"website,omitempty"`
	Phone               *string    `json:"phone,omitempty"`
	Email               *string    `json:"email,omitempty"`
	LastUpdated         *time.Time `json:"last_updated,omitempty"`
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
