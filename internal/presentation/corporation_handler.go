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
		Id:              corp.ID,
		CorporateNumber: corp.CorporateNumber,
		Name:            corp.Name,
		Status:          corp.Status,
		CreatedAt:       corp.CreatedAt,
		UpdatedAt:       corp.UpdatedAt,
	}

	if corp.NameKana != nil {
		apiCorp.NameKana = corp.NameKana
	}
	if corp.EnglishName != nil {
		apiCorp.EnglishName = corp.EnglishName
	}
	if corp.PostalCode != nil {
		apiCorp.PostalCode = corp.PostalCode
	}
	if corp.Address != nil {
		apiCorp.Address = corp.Address
	}
	if corp.PrefectureCode != nil {
		apiCorp.PrefectureCode = corp.PrefectureCode
	}
	if corp.CityCode != nil {
		apiCorp.CityCode = corp.CityCode
	}
	if corp.Phone != nil {
		apiCorp.Phone = corp.Phone
	}
	if corp.Email != nil {
		apiCorp.Email = corp.Email
	}
	if corp.Website != nil {
		apiCorp.Website = corp.Website
	}
	if corp.FoundingDate != nil {
		// Convert time.Time to openapi_types.Date
		date := openapi_types.Date{Time: *corp.FoundingDate}
		apiCorp.FoundingDate = &date
	}
	if corp.CapitalStock != nil {
		apiCorp.CapitalStock = corp.CapitalStock
	}
	if corp.EmployeeNumber != nil {
		apiCorp.EmployeeNumber = corp.EmployeeNumber
	}
	if corp.BusinessDescription != nil {
		apiCorp.BusinessDescription = corp.BusinessDescription
	}
	if corp.CorporateType != nil {
		apiCorp.CorporateType = corp.CorporateType
	}
	if corp.Representative != nil {
		apiCorp.Representative = corp.Representative
	}
	if corp.Industry != nil {
		apiCorp.Industry = corp.Industry
	}
	if corp.LastUpdated != nil {
		apiCorp.LastUpdated = corp.LastUpdated
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
