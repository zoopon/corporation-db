package usecase

import (
	"context"
	"fmt"
	"log"
	"regexp"
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
			if uc.IsSimilarBase(base, existing) {
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

// IsSimilarBase checks if two bases are similar (to avoid duplicates)
func (uc *FetchBasesUseCase) IsSimilarBase(base1, base2 *domain.Base) bool {
	// First check if they belong to the same corporation - only bases from the same corporation can be duplicates
	if base1.CorporateNumber != base2.CorporateNumber {
		return false
	}
	
	// Normalize postal codes for comparison
	postalCode1 := uc.normalizePostalCode(base1.PostalCode)
	postalCode2 := uc.normalizePostalCode(base2.PostalCode)
	
	// Check if normalized postal codes are the same and not empty
	if postalCode1 != "" && postalCode2 != "" && postalCode1 == postalCode2 {
		return true
	}
	
	// Check if locations are similar (normalize and compare)
	location1 := uc.normalizeLocation(base1.Location)
	location2 := uc.normalizeLocation(base2.Location)
	
	if location1 != "" && location2 != "" {
		// Check for exact match after normalization
		if location1 == location2 {
			return true
		}
		
		// Check if one location contains the other (for cases like "東京都目黒区中目黒2-9-13" vs "東京都目黒区中目黒２丁目９番１３号")
		if strings.Contains(location1, location2) || strings.Contains(location2, location1) {
			return true
		}
		
		// Check if addresses are similar by comparing normalized address parts
		if uc.areAddressesSimilar(location1, location2) {
			return true
		}
	}
	
	// For head office detection - if both are marked as head office for same corporation,
	// we should still check if they're in the same location to avoid false positives
	// (corporations can have multiple head offices due to relocations, mergers, etc.)
	if base1.IsHeadOffice && base2.IsHeadOffice {
		// Only consider them similar if they also have similar locations or postal codes
		if (postalCode1 != "" && postalCode2 != "" && postalCode1 == postalCode2) ||
		   (location1 != "" && location2 != "" && (location1 == location2 || 
		    strings.Contains(location1, location2) || strings.Contains(location2, location1) ||
		    uc.areAddressesSimilar(location1, location2))) {
			return true
		}
	}
	
	// Check if base names are similar (normalize and compare)
	baseName1 := uc.normalizeBaseName(base1.BaseName)
	baseName2 := uc.normalizeBaseName(base2.BaseName)
	
	if baseName1 != "" && baseName2 != "" && baseName1 == baseName2 {
		// If names are the same, they might be the same base
		return true
	}
	
	return false
}

// normalizePostalCode normalizes postal code for comparison
func (uc *FetchBasesUseCase) normalizePostalCode(postalCode *string) string {
	if postalCode == nil || *postalCode == "" {
		return ""
	}
	
	// Remove hyphens and spaces
	normalized := strings.ReplaceAll(*postalCode, "-", "")
	normalized = strings.ReplaceAll(normalized, " ", "")
	
	return normalized
}

// normalizeLocation normalizes location string for comparison
func (uc *FetchBasesUseCase) normalizeLocation(location *string) string {
	if location == nil || *location == "" {
		return ""
	}
	
	normalized := *location
	
	// Remove postal code from the beginning if present
	normalized = strings.TrimSpace(normalized)
	// Remove 7-digit postal code prefix (e.g., "1530061 ")
	postalCodePattern := regexp.MustCompile(`^[0-9]{7}\s+`)
	normalized = postalCodePattern.ReplaceAllString(normalized, "")
	
	// Normalize Japanese address representations
	normalized = strings.ReplaceAll(normalized, "１", "1")
	normalized = strings.ReplaceAll(normalized, "２", "2")
	normalized = strings.ReplaceAll(normalized, "３", "3")
	normalized = strings.ReplaceAll(normalized, "４", "4")
	normalized = strings.ReplaceAll(normalized, "５", "5")
	normalized = strings.ReplaceAll(normalized, "６", "6")
	normalized = strings.ReplaceAll(normalized, "７", "7")
	normalized = strings.ReplaceAll(normalized, "８", "8")
	normalized = strings.ReplaceAll(normalized, "９", "9")
	normalized = strings.ReplaceAll(normalized, "０", "0")
	
	// Normalize address format
	normalized = strings.ReplaceAll(normalized, "丁目", "-")
	normalized = strings.ReplaceAll(normalized, "番", "-")
	normalized = strings.ReplaceAll(normalized, "号", "")
	
	// Remove extra spaces
	normalized = strings.Join(strings.Fields(normalized), "")
	
	return normalized
}

// areAddressesSimilar checks if two addresses are similar
func (uc *FetchBasesUseCase) areAddressesSimilar(addr1, addr2 string) bool {
	// Extract main address parts (prefecture, city, district)
	parts1 := uc.extractAddressParts(addr1)
	parts2 := uc.extractAddressParts(addr2)
	
	// If we have at least 3 matching parts (prefecture, city, district), consider them similar
	matches := 0
	minParts := len(parts1)
	if len(parts2) < minParts {
		minParts = len(parts2)
	}
	
	for i := 0; i < minParts && i < 3; i++ {
		if parts1[i] == parts2[i] {
			matches++
		}
	}
	
	return matches >= 3
}

// extractAddressParts extracts main parts of Japanese address
func (uc *FetchBasesUseCase) extractAddressParts(address string) []string {
	// This is a simplified version - in production, you might want to use a proper Japanese address parser
	var parts []string
	
	// Look for prefecture
	prefectures := []string{"北海道", "青森県", "岩手県", "宮城県", "秋田県", "山形県", "福島県", 
		"茨城県", "栃木県", "群馬県", "埼玉県", "千葉県", "東京都", "神奈川県", "新潟県", 
		"富山県", "石川県", "福井県", "山梨県", "長野県", "岐阜県", "静岡県", "愛知県", 
		"三重県", "滋賀県", "京都府", "大阪府", "兵庫県", "奈良県", "和歌山県", "鳥取県", 
		"島根県", "岡山県", "広島県", "山口県", "徳島県", "香川県", "愛媛県", "高知県", 
		"福岡県", "佐賀県", "長崎県", "熊本県", "大分県", "宮崎県", "鹿児島県", "沖縄県"}
	
	for _, pref := range prefectures {
		if strings.Contains(address, pref) {
			parts = append(parts, pref)
			break
		}
	}
	
	// Look for city/ward (市、区、町、村)
	// This is a simplified extraction - you might want to use regex for better parsing
	
	return parts
}

// normalizeBaseName normalizes base name for comparison
func (uc *FetchBasesUseCase) normalizeBaseName(baseName *string) string {
	if baseName == nil || *baseName == "" {
		return ""
	}
	
	normalized := strings.TrimSpace(*baseName)
	
	// Normalize common base name variations
	normalized = strings.ReplaceAll(normalized, "本店", "本社")
	normalized = strings.ReplaceAll(normalized, "本部", "本社")
	normalized = strings.ReplaceAll(normalized, "head office", "本社")
	normalized = strings.ReplaceAll(normalized, "Head Office", "本社")
	
	return normalized
}
