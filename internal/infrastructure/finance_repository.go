package infrastructure

import (
"context"
"database/sql"

"corporation-db/internal/domain"
"corporation-db/internal/infrastructure/db"
)

type FinanceRepository struct {
	db      *sql.DB
	queries *db.Queries
}

func NewFinanceRepository(database *sql.DB) domain.FinanceRepository {
	return &FinanceRepository{
		db:      database,
		queries: db.New(database),
	}
}

func (r *FinanceRepository) Create(finance *domain.Finance) error {
	ctx := context.Background()
	return r.queries.CreateFinance(ctx, db.CreateFinanceParams{
CorporateNumber:              finance.CorporateNumber,
CorporateNameFromNumber:      sql.NullString{String: finance.CorporateNameFromNumber, Valid: finance.CorporateNameFromNumber != ""},
HeadOfficeLocationFromNumber: sql.NullString{String: finance.HeadOfficeLocationFromNumber, Valid: finance.HeadOfficeLocationFromNumber != ""},
CorporateName:                sql.NullString{String: finance.CorporateName, Valid: finance.CorporateName != ""},
HeadOfficeLocation:           sql.NullString{String: finance.HeadOfficeLocation, Valid: finance.HeadOfficeLocation != ""},
AccountingStandards:          sql.NullString{String: finance.AccountingStandards, Valid: finance.AccountingStandards != ""},
BusinessYear:                 sql.NullString{String: finance.BusinessYear, Valid: finance.BusinessYear != ""},
PeriodNumber:                 sql.NullString{String: finance.PeriodNumber, Valid: finance.PeriodNumber != ""},
SalesRevenue:                 sql.NullString{String: finance.SalesRevenue, Valid: finance.SalesRevenue != ""},
SalesRevenueUnit:             sql.NullString{String: finance.SalesRevenueUnit, Valid: finance.SalesRevenueUnit != ""},
OperatingRevenue1:            sql.NullString{String: finance.OperatingRevenue1, Valid: finance.OperatingRevenue1 != ""},
OperatingRevenue1Unit:        sql.NullString{String: finance.OperatingRevenue1Unit, Valid: finance.OperatingRevenue1Unit != ""},
OperatingRevenue2:            sql.NullString{String: finance.OperatingRevenue2, Valid: finance.OperatingRevenue2 != ""},
OperatingRevenue2Unit:        sql.NullString{String: finance.OperatingRevenue2Unit, Valid: finance.OperatingRevenue2Unit != ""},
GrossOperatingRevenue:        sql.NullString{String: finance.GrossOperatingRevenue, Valid: finance.GrossOperatingRevenue != ""},
GrossOperatingRevenueUnit:    sql.NullString{String: finance.GrossOperatingRevenueUnit, Valid: finance.GrossOperatingRevenueUnit != ""},
OrdinaryRevenue:              sql.NullString{String: finance.OrdinaryRevenue, Valid: finance.OrdinaryRevenue != ""},
OrdinaryRevenueUnit:          sql.NullString{String: finance.OrdinaryRevenueUnit, Valid: finance.OrdinaryRevenueUnit != ""},
NetPremiumsWritten:           sql.NullString{String: finance.NetPremiumsWritten, Valid: finance.NetPremiumsWritten != ""},
NetPremiumsWrittenUnit:       sql.NullString{String: finance.NetPremiumsWrittenUnit, Valid: finance.NetPremiumsWrittenUnit != ""},
OrdinaryIncome:               sql.NullString{String: finance.OrdinaryIncome, Valid: finance.OrdinaryIncome != ""},
OrdinaryIncomeUnit:           sql.NullString{String: finance.OrdinaryIncomeUnit, Valid: finance.OrdinaryIncomeUnit != ""},
NetIncome:                    sql.NullString{String: finance.NetIncome, Valid: finance.NetIncome != ""},
NetIncomeUnit:                sql.NullString{String: finance.NetIncomeUnit, Valid: finance.NetIncomeUnit != ""},
CapitalStock:                 sql.NullString{String: finance.CapitalStock, Valid: finance.CapitalStock != ""},
CapitalStockUnit:             sql.NullString{String: finance.CapitalStockUnit, Valid: finance.CapitalStockUnit != ""},
NetAssets:                    sql.NullString{String: finance.NetAssets, Valid: finance.NetAssets != ""},
NetAssetsUnit:                sql.NullString{String: finance.NetAssetsUnit, Valid: finance.NetAssetsUnit != ""},
TotalAssets:                  sql.NullString{String: finance.TotalAssets, Valid: finance.TotalAssets != ""},
TotalAssetsUnit:              sql.NullString{String: finance.TotalAssetsUnit, Valid: finance.TotalAssetsUnit != ""},
NumberOfEmployees:            sql.NullString{String: finance.NumberOfEmployees, Valid: finance.NumberOfEmployees != ""},
NumberOfEmployeesUnit:        sql.NullString{String: finance.NumberOfEmployeesUnit, Valid: finance.NumberOfEmployeesUnit != ""},
MajorShareholder1:            sql.NullString{String: finance.MajorShareholder1, Valid: finance.MajorShareholder1 != ""},
ShareholdingRatio1:           sql.NullString{String: finance.ShareholdingRatio1, Valid: finance.ShareholdingRatio1 != ""},
MajorShareholder2:            sql.NullString{String: finance.MajorShareholder2, Valid: finance.MajorShareholder2 != ""},
ShareholdingRatio2:           sql.NullString{String: finance.ShareholdingRatio2, Valid: finance.ShareholdingRatio2 != ""},
MajorShareholder3:            sql.NullString{String: finance.MajorShareholder3, Valid: finance.MajorShareholder3 != ""},
ShareholdingRatio3:           sql.NullString{String: finance.ShareholdingRatio3, Valid: finance.ShareholdingRatio3 != ""},
MajorShareholder4:            sql.NullString{String: finance.MajorShareholder4, Valid: finance.MajorShareholder4 != ""},
ShareholdingRatio4:           sql.NullString{String: finance.ShareholdingRatio4, Valid: finance.ShareholdingRatio4 != ""},
MajorShareholder5:            sql.NullString{String: finance.MajorShareholder5, Valid: finance.MajorShareholder5 != ""},
ShareholdingRatio5:           sql.NullString{String: finance.ShareholdingRatio5, Valid: finance.ShareholdingRatio5 != ""},
})
}

func (r *FinanceRepository) CreateBatch(finances []*domain.Finance) error {
	if len(finances) == 0 {
		return nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	queries := r.queries.WithTx(tx)
	ctx := context.Background()

	for _, finance := range finances {
		err := queries.CreateFinance(ctx, db.CreateFinanceParams{
CorporateNumber:              finance.CorporateNumber,
CorporateNameFromNumber:      sql.NullString{String: finance.CorporateNameFromNumber, Valid: finance.CorporateNameFromNumber != ""},
HeadOfficeLocationFromNumber: sql.NullString{String: finance.HeadOfficeLocationFromNumber, Valid: finance.HeadOfficeLocationFromNumber != ""},
CorporateName:                sql.NullString{String: finance.CorporateName, Valid: finance.CorporateName != ""},
HeadOfficeLocation:           sql.NullString{String: finance.HeadOfficeLocation, Valid: finance.HeadOfficeLocation != ""},
AccountingStandards:          sql.NullString{String: finance.AccountingStandards, Valid: finance.AccountingStandards != ""},
BusinessYear:                 sql.NullString{String: finance.BusinessYear, Valid: finance.BusinessYear != ""},
PeriodNumber:                 sql.NullString{String: finance.PeriodNumber, Valid: finance.PeriodNumber != ""},
SalesRevenue:                 sql.NullString{String: finance.SalesRevenue, Valid: finance.SalesRevenue != ""},
SalesRevenueUnit:             sql.NullString{String: finance.SalesRevenueUnit, Valid: finance.SalesRevenueUnit != ""},
OperatingRevenue1:            sql.NullString{String: finance.OperatingRevenue1, Valid: finance.OperatingRevenue1 != ""},
OperatingRevenue1Unit:        sql.NullString{String: finance.OperatingRevenue1Unit, Valid: finance.OperatingRevenue1Unit != ""},
OperatingRevenue2:            sql.NullString{String: finance.OperatingRevenue2, Valid: finance.OperatingRevenue2 != ""},
OperatingRevenue2Unit:        sql.NullString{String: finance.OperatingRevenue2Unit, Valid: finance.OperatingRevenue2Unit != ""},
GrossOperatingRevenue:        sql.NullString{String: finance.GrossOperatingRevenue, Valid: finance.GrossOperatingRevenue != ""},
GrossOperatingRevenueUnit:    sql.NullString{String: finance.GrossOperatingRevenueUnit, Valid: finance.GrossOperatingRevenueUnit != ""},
OrdinaryRevenue:              sql.NullString{String: finance.OrdinaryRevenue, Valid: finance.OrdinaryRevenue != ""},
OrdinaryRevenueUnit:          sql.NullString{String: finance.OrdinaryRevenueUnit, Valid: finance.OrdinaryRevenueUnit != ""},
NetPremiumsWritten:           sql.NullString{String: finance.NetPremiumsWritten, Valid: finance.NetPremiumsWritten != ""},
NetPremiumsWrittenUnit:       sql.NullString{String: finance.NetPremiumsWrittenUnit, Valid: finance.NetPremiumsWrittenUnit != ""},
OrdinaryIncome:               sql.NullString{String: finance.OrdinaryIncome, Valid: finance.OrdinaryIncome != ""},
OrdinaryIncomeUnit:           sql.NullString{String: finance.OrdinaryIncomeUnit, Valid: finance.OrdinaryIncomeUnit != ""},
NetIncome:                    sql.NullString{String: finance.NetIncome, Valid: finance.NetIncome != ""},
NetIncomeUnit:                sql.NullString{String: finance.NetIncomeUnit, Valid: finance.NetIncomeUnit != ""},
CapitalStock:                 sql.NullString{String: finance.CapitalStock, Valid: finance.CapitalStock != ""},
CapitalStockUnit:             sql.NullString{String: finance.CapitalStockUnit, Valid: finance.CapitalStockUnit != ""},
NetAssets:                    sql.NullString{String: finance.NetAssets, Valid: finance.NetAssets != ""},
NetAssetsUnit:                sql.NullString{String: finance.NetAssetsUnit, Valid: finance.NetAssetsUnit != ""},
TotalAssets:                  sql.NullString{String: finance.TotalAssets, Valid: finance.TotalAssets != ""},
TotalAssetsUnit:              sql.NullString{String: finance.TotalAssetsUnit, Valid: finance.TotalAssetsUnit != ""},
NumberOfEmployees:            sql.NullString{String: finance.NumberOfEmployees, Valid: finance.NumberOfEmployees != ""},
NumberOfEmployeesUnit:        sql.NullString{String: finance.NumberOfEmployeesUnit, Valid: finance.NumberOfEmployeesUnit != ""},
MajorShareholder1:            sql.NullString{String: finance.MajorShareholder1, Valid: finance.MajorShareholder1 != ""},
ShareholdingRatio1:           sql.NullString{String: finance.ShareholdingRatio1, Valid: finance.ShareholdingRatio1 != ""},
MajorShareholder2:            sql.NullString{String: finance.MajorShareholder2, Valid: finance.MajorShareholder2 != ""},
ShareholdingRatio2:           sql.NullString{String: finance.ShareholdingRatio2, Valid: finance.ShareholdingRatio2 != ""},
MajorShareholder3:            sql.NullString{String: finance.MajorShareholder3, Valid: finance.MajorShareholder3 != ""},
ShareholdingRatio3:           sql.NullString{String: finance.ShareholdingRatio3, Valid: finance.ShareholdingRatio3 != ""},
MajorShareholder4:            sql.NullString{String: finance.MajorShareholder4, Valid: finance.MajorShareholder4 != ""},
ShareholdingRatio4:           sql.NullString{String: finance.ShareholdingRatio4, Valid: finance.ShareholdingRatio4 != ""},
MajorShareholder5:            sql.NullString{String: finance.MajorShareholder5, Valid: finance.MajorShareholder5 != ""},
ShareholdingRatio5:           sql.NullString{String: finance.ShareholdingRatio5, Valid: finance.ShareholdingRatio5 != ""},
})
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *FinanceRepository) GetByCorporateNumber(corporateNumber string) ([]*domain.Finance, error) {
	ctx := context.Background()
	finances, err := r.queries.GetFinancesByCorporateNumber(ctx, corporateNumber)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.Finance, len(finances))
	for i, finance := range finances {
		result[i] = r.convertToFinance(finance)
	}

	return result, nil
}

func (r *FinanceRepository) GetLatest(corporateNumber string) (*domain.Finance, error) {
	ctx := context.Background()
	finance, err := r.queries.GetLatestFinanceByCorporateNumber(ctx, corporateNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return r.convertToFinance(finance), nil
}

func (r *FinanceRepository) Count() (int64, error) {
	ctx := context.Background()
	return r.queries.CountFinances(ctx)
}

func (r *FinanceRepository) DeleteAll() error {
	ctx := context.Background()
	return r.queries.DeleteAllFinances(ctx)
}

func (r *FinanceRepository) DeleteByCorporateNumber(corporateNumber string) error {
	ctx := context.Background()
	return r.queries.DeleteFinancesByCorporateNumber(ctx, corporateNumber)
}

func (r *FinanceRepository) convertToFinance(finance db.Finance) *domain.Finance {
	return &domain.Finance{
		ID:                           int(finance.ID),
		CorporateNumber:              finance.CorporateNumber,
		CorporateNameFromNumber:      finance.CorporateNameFromNumber.String,
		HeadOfficeLocationFromNumber: finance.HeadOfficeLocationFromNumber.String,
		CorporateName:                finance.CorporateName.String,
		HeadOfficeLocation:           finance.HeadOfficeLocation.String,
		AccountingStandards:          finance.AccountingStandards.String,
		BusinessYear:                 finance.BusinessYear.String,
		PeriodNumber:                 finance.PeriodNumber.String,
		SalesRevenue:                 finance.SalesRevenue.String,
		SalesRevenueUnit:             finance.SalesRevenueUnit.String,
		OperatingRevenue1:            finance.OperatingRevenue1.String,
		OperatingRevenue1Unit:        finance.OperatingRevenue1Unit.String,
		OperatingRevenue2:            finance.OperatingRevenue2.String,
		OperatingRevenue2Unit:        finance.OperatingRevenue2Unit.String,
		GrossOperatingRevenue:        finance.GrossOperatingRevenue.String,
		GrossOperatingRevenueUnit:    finance.GrossOperatingRevenueUnit.String,
		OrdinaryRevenue:              finance.OrdinaryRevenue.String,
		OrdinaryRevenueUnit:          finance.OrdinaryRevenueUnit.String,
		NetPremiumsWritten:           finance.NetPremiumsWritten.String,
		NetPremiumsWrittenUnit:       finance.NetPremiumsWrittenUnit.String,
		OrdinaryIncome:               finance.OrdinaryIncome.String,
		OrdinaryIncomeUnit:           finance.OrdinaryIncomeUnit.String,
		NetIncome:                    finance.NetIncome.String,
		NetIncomeUnit:                finance.NetIncomeUnit.String,
		CapitalStock:                 finance.CapitalStock.String,
		CapitalStockUnit:             finance.CapitalStockUnit.String,
		NetAssets:                    finance.NetAssets.String,
		NetAssetsUnit:                finance.NetAssetsUnit.String,
		TotalAssets:                  finance.TotalAssets.String,
		TotalAssetsUnit:              finance.TotalAssetsUnit.String,
		NumberOfEmployees:            finance.NumberOfEmployees.String,
		NumberOfEmployeesUnit:        finance.NumberOfEmployeesUnit.String,
		MajorShareholder1:            finance.MajorShareholder1.String,
		ShareholdingRatio1:           finance.ShareholdingRatio1.String,
		MajorShareholder2:            finance.MajorShareholder2.String,
		ShareholdingRatio2:           finance.ShareholdingRatio2.String,
		MajorShareholder3:            finance.MajorShareholder3.String,
		ShareholdingRatio3:           finance.ShareholdingRatio3.String,
		MajorShareholder4:            finance.MajorShareholder4.String,
		ShareholdingRatio4:           finance.ShareholdingRatio4.String,
		MajorShareholder5:            finance.MajorShareholder5.String,
		ShareholdingRatio5:           finance.ShareholdingRatio5.String,
		CreatedAt:                    finance.CreatedAt.Time,
		UpdatedAt:                    finance.UpdatedAt.Time,
	}
}
