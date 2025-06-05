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
	baseUsecase        *usecase.BaseUsecase
	fetchBasesUsecase  *usecase.FetchBasesUseCase
}

func NewCorporationHandler(corporationUsecase *usecase.CorporationUsecase, baseUsecase *usecase.BaseUsecase, fetchBasesUsecase *usecase.FetchBasesUseCase) *CorporationHandler {
	return &CorporationHandler{
		corporationUsecase: corporationUsecase,
		baseUsecase:        baseUsecase,
		fetchBasesUsecase:  fetchBasesUsecase,
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
	corp, err := h.corporationUsecase.GetByCorporateNumber(r.Context(), corporateNumber)
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

	// Get bases for this corporation
	bases, err := h.baseUsecase.GetBasesByCorporateNumber(r.Context(), corporateNumber)
	if err != nil {
		// Log error but don't fail the request - bases are optional
		// Just return corporation without bases
		bases = []*domain.Base{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(convertCorporationToAPIWithBases(corp, bases))
}

// FetchCorporationBases implements the POST /corporations/{corporate_number}/fetch-bases endpoint
func (h *CorporationHandler) FetchCorporationBases(w http.ResponseWriter, r *http.Request, corporateNumber string) {
	// Validate corporate number format
	if len(corporateNumber) != 13 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid corporate number format",
		})
		return
	}

	// Execute fetch bases use case
	result, err := h.fetchBasesUsecase.Execute(r.Context(), corporateNumber)
	if err != nil {
		if err.Error() == "failed to get corporation: sql: no rows in result set" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Corporation not found",
			})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to fetch bases information",
		})
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
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

// Helper function to convert base domain to API model
func convertBaseToAPI(base *domain.Base) api.Base {
	apiBase := api.Base{
		Id:              openapi_types.UUID(base.ID),
		CorporationId:   openapi_types.UUID(base.CorporationID),
		CorporateNumber: base.CorporateNumber,
		CountryCode:     base.CountryCode,
		Location:        base.Location,
		DataObtainedAt:  base.DataObtainedAt,
		DataSourceUrl:   base.DataSourceURL,
		IsHeadOffice:    base.IsHeadOffice,
		CreatedAt:       base.CreatedAt,
		UpdatedAt:       base.UpdatedAt,
	}

	if base.BaseName != nil {
		apiBase.BaseName = base.BaseName
	}
	if base.PostalCode != nil {
		apiBase.PostalCode = base.PostalCode
	}
	if base.PhoneNumber != nil {
		apiBase.PhoneNumber = base.PhoneNumber
	}
	if base.FaxNumber != nil {
		apiBase.FaxNumber = base.FaxNumber
	}

	return apiBase
}

// Helper function to convert corporation with bases to API model
func convertCorporationToAPIWithBases(corp *domain.Corporation, bases []*domain.Base) api.Corporation {
	apiCorp := convertCorporationToAPI(corp)

	if len(bases) > 0 {
		apiBases := make([]api.Base, len(bases))
		for i, base := range bases {
			apiBases[i] = convertBaseToAPI(base)
		}
		apiCorp.Bases = &apiBases
	}

	return apiCorp
}
