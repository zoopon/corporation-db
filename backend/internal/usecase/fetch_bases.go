package usecase

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"corporation-db/internal/domain"

	"github.com/google/uuid"
)

// FetchBasesUseCase handles fetching bases using OpenAI API
type FetchBasesUseCase struct {
	corporationRepo domain.CorporationRepository
	baseRepo        domain.BaseRepository
	openAIService   domain.OpenAIService
}

// NewFetchBasesUseCase creates a new FetchBasesUseCase
func NewFetchBasesUseCase(
	corporationRepo domain.CorporationRepository,
	baseRepo domain.BaseRepository,
	openAIService domain.OpenAIService,
) *FetchBasesUseCase {
	return &FetchBasesUseCase{
		corporationRepo: corporationRepo,
		baseRepo:        baseRepo,
		openAIService:   openAIService,
	}
}

// FetchBasesResult represents the result of fetching bases
type FetchBasesResult struct {
	Message    string   `json:"message"`
	BasesCount int      `json:"bases_count"`
	URLsFound  []string `json:"urls_found"`
}

// Execute fetches base/branch offices using OpenAI API
func (uc *FetchBasesUseCase) Execute(ctx context.Context, corporateNumber string) (*FetchBasesResult, error) {
	// Get corporation by corporate number
	corporation, err := uc.corporationRepo.GetByCorporateNumber(ctx, corporateNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get corporation: %w", err)
	}

	// Discover URLs containing office information
	log.Printf("Discovering URLs for corporation: %s", corporation.Name)
	urlResult, err := uc.openAIService.DiscoverURLs(ctx, corporation.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to discover URLs: %w", err)
	}

	if len(urlResult.URLs) == 0 {
		return &FetchBasesResult{
			Message:    "No URLs found containing office information",
			BasesCount: 0,
			URLsFound:  []string{},
		}, nil
	}

	log.Printf("Found %d URLs: %v", len(urlResult.URLs), urlResult.URLs)

	// Extract base information from all URLs at once
	log.Printf("Extracting bases from %d URLs", len(urlResult.URLs))
	baseResult, err := uc.openAIService.ExtractBasesFromURL(ctx, urlResult.URLs, corporation.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to extract bases from URLs: %w", err)
	}

	if len(baseResult.Bases) == 0 {
		return &FetchBasesResult{
			Message:    "No base information could be extracted from the found URLs",
			BasesCount: 0,
			URLsFound:  urlResult.URLs,
		}, nil
	}

	log.Printf("Extracted %d bases", len(baseResult.Bases))

	// Convert extracted bases to domain entities and save to database
	savedCount := 0
	for _, extractedBase := range baseResult.Bases {
		// Create location string from prefecture, city, and address (fix whitespace issues)
		var location *string
		if extractedBase.Prefecture != "" || extractedBase.City != "" || extractedBase.Address != "" {
			var parts []string
			if strings.TrimSpace(extractedBase.Prefecture) != "" {
				parts = append(parts, strings.TrimSpace(extractedBase.Prefecture))
			}
			if strings.TrimSpace(extractedBase.City) != "" {
				parts = append(parts, strings.TrimSpace(extractedBase.City))
			}
			if strings.TrimSpace(extractedBase.Address) != "" {
				parts = append(parts, strings.TrimSpace(extractedBase.Address))
			}
			if len(parts) > 0 {
				locationStr := strings.Join(parts, " ")
				location = &locationStr
			}
		}

		// Create postal code pointer
		var postalCode *string
		if extractedBase.PostalCode != "" {
			postalCode = &extractedBase.PostalCode
		}

		// Create phone number pointer
		var phoneNumber *string
		if extractedBase.PhoneNumber != "" {
			phoneNumber = &extractedBase.PhoneNumber
		}

		// Create fax number pointer
		var faxNumber *string
		if extractedBase.FaxNumber != "" {
			faxNumber = &extractedBase.FaxNumber
		}

		// Create base name pointer
		var baseName *string
		if extractedBase.Name != "" {
			baseName = &extractedBase.Name
		}

		// Use the source URL instead of hardcoded "OpenAI API"
		dataSourceURL := extractedBase.SourceURL
		if dataSourceURL == "" {
			dataSourceURL = "OpenAI API" // Fallback if URL is not available
		}

		base := &domain.Base{
			ID:              uuid.New(),
			CorporationID:   corporation.ID,
			CorporateNumber: corporation.CorporateNumber,
			BaseName:        baseName,
			CountryCode:     "JP", // Default to Japan
			PostalCode:      postalCode,
			Location:        location,
			PhoneNumber:     phoneNumber,
			FaxNumber:       faxNumber,
			DataObtainedAt:  time.Now(),
			DataSourceURL:   dataSourceURL,
			IsHeadOffice:    extractedBase.Type == "head_office",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		// Check if a similar base already exists to avoid duplicates
		existingBases, err := uc.baseRepo.GetByCorporationID(ctx, corporation.ID)
		if err != nil {
			log.Printf("Failed to check existing bases: %v", err)
			continue
		}

		isDuplicate := false
		for _, existing := range existingBases {
			if uc.isSimilarBase(base, existing) {
				log.Printf("Skipping duplicate base: %s", *base.BaseName)
				isDuplicate = true
				break
			}
		}

		if !isDuplicate {
			_, err := uc.baseRepo.Create(ctx, base)
			if err != nil {
				log.Printf("Failed to save base %s: %v", *base.BaseName, err)
				continue
			}
			savedCount++
			log.Printf("Saved base: %s", *base.BaseName)
		}
	}

	message := fmt.Sprintf("Successfully fetched and saved %d base/branch offices", savedCount)
	if savedCount != len(baseResult.Bases) {
		message += fmt.Sprintf(" (%d duplicates skipped)", len(baseResult.Bases)-savedCount)
	}

	return &FetchBasesResult{
		Message:    message,
		BasesCount: savedCount,
		URLsFound:  urlResult.URLs,
	}, nil
}

// isSimilarBase checks if two bases are similar (to avoid duplicates)
func (uc *FetchBasesUseCase) isSimilarBase(base1, base2 *domain.Base) bool {
	// Check if locations are similar
	if base1.Location != nil && base2.Location != nil && *base1.Location == *base2.Location {
		return true
	}

	// Check if names are similar
	if base1.BaseName != nil && base2.BaseName != nil && *base1.BaseName == *base2.BaseName {
		return true
	}

	// Check if postal codes are the same
	if base1.PostalCode != nil && base2.PostalCode != nil && *base1.PostalCode == *base2.PostalCode {
		return true
	}

	return false
}
