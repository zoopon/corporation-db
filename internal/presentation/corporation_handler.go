package presentation

import (
	"encoding/json"
	"net/http"

	"corporation-db/internal/api"
	"corporation-db/internal/domain"
	"corporation-db/internal/usecase"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

type CorporationHandler struct {
	corporationUsecase *usecase.CorporationUsecase
}

func NewCorporationHandler(corporationUsecase *usecase.CorporationUsecase) *CorporationHandler {
	return &CorporationHandler{
		corporationUsecase: corporationUsecase,
	}
}

// GetCorporations implements the GET /corporations endpoint
func (h *CorporationHandler) GetCorporations(w http.ResponseWriter, r *http.Request, params api.GetCorporationsParams) {
	// Set defaults
	limit := 100
	offset := 0
	if params.Limit != nil {
		limit = *params.Limit
	}
	if params.Offset != nil {
		offset = *params.Offset
	}

	// Create filter
	filter := domain.CorporationFilter{
		CorporateNumber: params.CorporateNumber,
		Name:            params.Name,
		Prefecture:      params.Prefecture,
		PrefectureCode:  params.PrefectureCode,
		Limit:           limit,
		Offset:          offset,
	}

	// Get corporations from usecase
	corporations, total, err := h.corporationUsecase.GetCorporations(r.Context(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to response format
	response := map[string]interface{}{
		"corporations": convertCorporationsToAPI(corporations),
		"total":        total,
		"limit":        limit,
		"offset":       offset,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetCorporationsCorporateNumber implements the GET /corporations/{corporate_number} endpoint
func (h *CorporationHandler) GetCorporationsCorporateNumber(w http.ResponseWriter, r *http.Request, corporateNumber string) {
	corp, err := h.corporationUsecase.GetCorporationByCorporateNumber(r.Context(), corporateNumber)
	if err != nil {
		switch err {
		case domain.ErrCorporationNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		case domain.ErrInvalidCorporateNumber:
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(convertCorporationToAPI(corp))
}

// Helper functions for conversion between domain and API models

func convertCorporationToAPI(corp *domain.Corporation) api.Corporation {
	apiCorp := api.Corporation{
		Id:              openapi_types.UUID(corp.ID),
		CorporateNumber: corp.CorporateNumber,
		Name:            corp.Name,
		Status:          corp.Status,
		CreatedAt:       corp.CreatedAt,
		UpdatedAt:       corp.UpdatedAt,
	}

	// Basic Information
	if corp.Kana != nil {
		apiCorp.Kana = corp.Kana
	}
	if corp.NameEn != nil {
		apiCorp.NameEn = corp.NameEn
	}
	if corp.PostalCode != nil {
		apiCorp.PostalCode = corp.PostalCode
	}
	if corp.Location != nil {
		apiCorp.Location = corp.Location
	}
	if corp.PrefectureCode != nil {
		apiCorp.PrefectureCode = corp.PrefectureCode
	}

	// Registration Information
	if corp.CloseDate != nil {
		date := openapi_types.Date{Time: *corp.CloseDate}
		apiCorp.CloseDate = &date
	}
	if corp.CloseCause != nil {
		apiCorp.CloseCause = corp.CloseCause
	}

	// Representative Information
	if corp.RepresentativeName != nil {
		apiCorp.RepresentativeName = corp.RepresentativeName
	}
	if corp.RepresentativePosition != nil {
		apiCorp.RepresentativePosition = corp.RepresentativePosition
	}

	// Company Details
	if corp.DateOfEstablishment != nil {
		date := openapi_types.Date{Time: *corp.DateOfEstablishment}
		apiCorp.DateOfEstablishment = &date
	}
	if corp.FoundingYear != nil {
		foundingYear := int(*corp.FoundingYear)
		apiCorp.FoundingYear = &foundingYear
	}
	if corp.CapitalStock != nil {
		apiCorp.CapitalStock = corp.CapitalStock
	}
	if corp.EmployeeNumber != nil {
		employeeNumber := int(*corp.EmployeeNumber)
		apiCorp.EmployeeNumber = &employeeNumber
	}
	if corp.CompanySizeMale != nil {
		companySizeMale := int(*corp.CompanySizeMale)
		apiCorp.CompanySizeMale = &companySizeMale
	}
	if corp.CompanySizeFemale != nil {
		companySizeFemale := int(*corp.CompanySizeFemale)
		apiCorp.CompanySizeFemale = &companySizeFemale
	}

	// Business Information
	if corp.BusinessItems != nil {
		apiCorp.BusinessItems = corp.BusinessItems
	}
	if corp.BusinessSummary != nil {
		apiCorp.BusinessSummary = corp.BusinessSummary
	}
	if corp.CompanyUrl != nil {
		apiCorp.CompanyUrl = corp.CompanyUrl
	}
	if corp.QualificationGrade != nil {
		apiCorp.QualificationGrade = corp.QualificationGrade
	}
	if corp.NumberOfActivity != nil {
		apiCorp.NumberOfActivity = corp.NumberOfActivity
	}

	// gBizINFO Metadata
	if corp.UpdateDate != nil {
		date := openapi_types.Date{Time: *corp.UpdateDate}
		apiCorp.UpdateDate = &date
	}

	return apiCorp
}

func convertCorporationsToAPI(corps []*domain.Corporation) []api.Corporation {
	result := make([]api.Corporation, len(corps))
	for i, corp := range corps {
		result[i] = convertCorporationToAPI(corp)
	}
	return result
}
