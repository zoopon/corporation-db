package usecase

import (
	"testing"
	"time"

	"corporation-db/internal/domain"

	"github.com/google/uuid"
)

func TestFetchBasesUseCase_IsSimilarBase(t *testing.T) {
	// Initialize usecase (repositories can be nil for this test)
	uc := &FetchBasesUseCase{}

	// Common test data
	corporationID := uuid.New()
	corporateNumber := "2013201003041"
	baseTime := time.Now()

	tests := []struct {
		name     string
		base1    *domain.Base
		base2    *domain.Base
		expected bool
	}{
		{
			name: "Same postal code with different formats (hyphen vs no hyphen)",
			base1: &domain.Base{
				ID:              uuid.New(),
				CorporationID:   corporationID,
				CorporateNumber: corporateNumber,
				BaseName:        stringPtr("本社"),
				PostalCode:      stringPtr("153-8636"),
				Location:        stringPtr("東京都目黒区中目黒2-9-13"),
				IsHeadOffice:    true,
				CreatedAt:       baseTime,
				UpdatedAt:       baseTime,
			},
			base2: &domain.Base{
				ID:              uuid.New(),
				CorporationID:   corporationID,
				CorporateNumber: corporateNumber,
				BaseName:        stringPtr("本社"),
				PostalCode:      stringPtr("1538636"), // Match normalized format
				Location:        stringPtr("東京都目黒区中目黒２丁目９番１３号"),
				IsHeadOffice:    true,
				CreatedAt:       baseTime,
				UpdatedAt:       baseTime,
			},
			expected: true,
		},
		{
			name: "Same location with different Japanese notation",
			base1: &domain.Base{
				ID:              uuid.New(),
				CorporationID:   corporationID,
				CorporateNumber: corporateNumber,
				BaseName:        stringPtr("本社"),
				PostalCode:      stringPtr("153-8636"),
				Location:        stringPtr("東京都目黒区中目黒２丁目９番１３号"),
				IsHeadOffice:    true,
				CreatedAt:       baseTime,
				UpdatedAt:       baseTime,
			},
			base2: &domain.Base{
				ID:              uuid.New(),
				CorporationID:   corporationID,
				CorporateNumber: corporateNumber,
				BaseName:        stringPtr("本社"),
				PostalCode:      stringPtr("153-8636"),
				Location:        stringPtr("東京都目黒区中目黒2-9-13"),
				IsHeadOffice:    true,
				CreatedAt:       baseTime,
				UpdatedAt:       baseTime,
			},
			expected: true,
		},
		{
			name: "Base name variations (本店 vs 本社)",
			base1: &domain.Base{
				ID:              uuid.New(),
				CorporationID:   corporationID,
				CorporateNumber: corporateNumber,
				BaseName:        stringPtr("本店"),
				PostalCode:      stringPtr("153-8636"),
				Location:        stringPtr("東京都目黒区中目黒2-9-13"),
				IsHeadOffice:    true,
				CreatedAt:       baseTime,
				UpdatedAt:       baseTime,
			},
			base2: &domain.Base{
				ID:              uuid.New(),
				CorporationID:   corporationID,
				CorporateNumber: corporateNumber,
				BaseName:        stringPtr("本社"),
				PostalCode:      stringPtr("153-8636"),
				Location:        stringPtr("東京都目黒区中目黒2-9-13"),
				IsHeadOffice:    true,
				CreatedAt:       baseTime,
				UpdatedAt:       baseTime,
			},
			expected: true,
		},
		{
			name: "Location with postal code prefix vs without",
			base1: &domain.Base{
				ID:              uuid.New(),
				CorporationID:   corporationID,
				CorporateNumber: corporateNumber,
				BaseName:        stringPtr("本社"),
				PostalCode:      stringPtr("153-8636"),
				Location:        stringPtr("1530061 東京都目黒区中目黒２丁目９番１３号"),
				IsHeadOffice:    true,
				CreatedAt:       baseTime,
				UpdatedAt:       baseTime,
			},
			base2: &domain.Base{
				ID:              uuid.New(),
				CorporationID:   corporationID,
				CorporateNumber: corporateNumber,
				BaseName:        stringPtr("本社"),
				PostalCode:      stringPtr("153-8636"),
				Location:        stringPtr("東京都目黒区中目黒2-9-13"),
				IsHeadOffice:    true,
				CreatedAt:       baseTime,
				UpdatedAt:       baseTime,
			},
			expected: true,
		},
		{
			name: "Both marked as head office for same corporation",
			base1: &domain.Base{
				ID:              uuid.New(),
				CorporationID:   corporationID,
				CorporateNumber: corporateNumber,
				BaseName:        stringPtr("本社A"),
				PostalCode:      stringPtr("100-0001"),
				Location:        stringPtr("東京都千代田区千代田1-1"),
				IsHeadOffice:    true,
				CreatedAt:       baseTime,
				UpdatedAt:       baseTime,
			},
			base2: &domain.Base{
				ID:              uuid.New(),
				CorporationID:   corporationID,
				CorporateNumber: corporateNumber,
				BaseName:        stringPtr("本社B"),
				PostalCode:      stringPtr("200-0002"),
				Location:        stringPtr("東京都立川市曙町2-2"),
				IsHeadOffice:    true,
				CreatedAt:       baseTime,
				UpdatedAt:       baseTime,
			},
			expected: false,
		},
		{
			name: "Both marked as head office for same corporation with same location - should be similar",
			base1: &domain.Base{
				ID:              uuid.New(),
				CorporationID:   corporationID,
				CorporateNumber: corporateNumber,
				BaseName:        stringPtr("本社"),
				PostalCode:      stringPtr("153-8636"),
				Location:        stringPtr("東京都目黒区中目黒2-9-13"),
				IsHeadOffice:    true,
				CreatedAt:       baseTime,
				UpdatedAt:       baseTime,
			},
			base2: &domain.Base{
				ID:              uuid.New(),
				CorporationID:   corporationID,
				CorporateNumber: corporateNumber,
				BaseName:        stringPtr("Head Office"),
				PostalCode:      stringPtr("153-8636"),
				Location:        stringPtr("東京都目黒区中目黒2-9-13"),
				IsHeadOffice:    true,
				CreatedAt:       baseTime,
				UpdatedAt:       baseTime,
			},
			expected: true,
		},
		{
			name: "Base name営業所 vs 本社 with same location - should be similar due to location",
			base1: &domain.Base{
				ID:              uuid.New(),
				CorporationID:   corporationID,
				CorporateNumber: corporateNumber,
				BaseName:        stringPtr("営業所"),
				PostalCode:      stringPtr("153-8636"),
				Location:        stringPtr("東京都目黒区中目黒2-9-13"),
				IsHeadOffice:    false,
				CreatedAt:       baseTime,
				UpdatedAt:       baseTime,
			},
			base2: &domain.Base{
				ID:              uuid.New(),
				CorporationID:   corporationID,
				CorporateNumber: corporateNumber,
				BaseName:        stringPtr("本社"),
				PostalCode:      stringPtr("153-8636"),
				Location:        stringPtr("東京都目黒区中目黒2-9-13"),
				IsHeadOffice:    true,
				CreatedAt:       baseTime,
				UpdatedAt:       baseTime,
			},
			expected: true,
		},
		{
			name: "Base name営業所 vs 本社 with different locations - should not be similar (different names and locations)",
			base1: &domain.Base{
				ID:              uuid.New(),
				CorporationID:   corporationID,
				CorporateNumber: corporateNumber,
				BaseName:        stringPtr("営業所"),
				PostalCode:      stringPtr("153-8636"),
				Location:        stringPtr("東京都目黒区中目黒2-9-13"),
				IsHeadOffice:    false,
				CreatedAt:       baseTime,
				UpdatedAt:       baseTime,
			},
			base2: &domain.Base{
				ID:              uuid.New(),
				CorporationID:   corporationID,
				CorporateNumber: corporateNumber,
				BaseName:        stringPtr("本社"),
				PostalCode:      stringPtr("100-0001"),
				Location:        stringPtr("東京都千代田区千代田1-1"),
				IsHeadOffice:    true,
				CreatedAt:       baseTime,
				UpdatedAt:       baseTime,
			},
			expected: false, // Different names and locations should not be similar
		},
		{
			name: "English base name variations",
			base1: &domain.Base{
				ID:              uuid.New(),
				CorporationID:   corporationID,
				CorporateNumber: corporateNumber,
				BaseName:        stringPtr("Head Office"),
				PostalCode:      stringPtr("153-8636"),
				Location:        stringPtr("東京都目黒区中目黒2-9-13"),
				IsHeadOffice:    true,
				CreatedAt:       baseTime,
				UpdatedAt:       baseTime,
			},
			base2: &domain.Base{
				ID:              uuid.New(),
				CorporationID:   corporationID,
				CorporateNumber: corporateNumber,
				BaseName:        stringPtr("本社"),
				PostalCode:      stringPtr("153-8636"),
				Location:        stringPtr("東京都目黒区中目黒2-9-13"),
				IsHeadOffice:    true,
				CreatedAt:       baseTime,
				UpdatedAt:       baseTime,
			},
			expected: true,
		},
		{
			name: "Different locations - should not be similar",
			base1: &domain.Base{
				ID:              uuid.New(),
				CorporationID:   corporationID,
				CorporateNumber: corporateNumber,
				BaseName:        stringPtr("本社"),
				PostalCode:      stringPtr("153-8636"),
				Location:        stringPtr("東京都目黒区中目黒2-9-13"),
				IsHeadOffice:    true,
				CreatedAt:       baseTime,
				UpdatedAt:       baseTime,
			},
			base2: &domain.Base{
				ID:              uuid.New(),
				CorporationID:   corporationID,
				CorporateNumber: corporateNumber,
				BaseName:        stringPtr("大阪支店"),
				PostalCode:      stringPtr("530-0001"),
				Location:        stringPtr("大阪府大阪市北区梅田1-1-1"),
				IsHeadOffice:    false,
				CreatedAt:       baseTime,
				UpdatedAt:       baseTime,
			},
			expected: false,
		},
		{
			name: "Different postal codes and locations - should not be similar",
			base1: &domain.Base{
				ID:              uuid.New(),
				CorporationID:   corporationID,
				CorporateNumber: corporateNumber,
				BaseName:        stringPtr("本社"),
				PostalCode:      stringPtr("153-8636"),
				Location:        stringPtr("東京都目黒区中目黒2-9-13"),
				IsHeadOffice:    true,
				CreatedAt:       baseTime,
				UpdatedAt:       baseTime,
			},
			base2: &domain.Base{
				ID:              uuid.New(),
				CorporationID:   corporationID,
				CorporateNumber: corporateNumber,
				BaseName:        stringPtr("大阪支店"), // Different base name
				PostalCode:      stringPtr("100-0001"),
				Location:        stringPtr("東京都千代田区千代田1-1"),
				IsHeadOffice:    false, // Different head office status
				CreatedAt:       baseTime,
				UpdatedAt:       baseTime,
			},
			expected: false,
		},
		{
			name: "Nil postal codes and locations - should not be similar",
			base1: &domain.Base{
				ID:              uuid.New(),
				CorporationID:   corporationID,
				CorporateNumber: corporateNumber,
				BaseName:        stringPtr("本社A"),
				PostalCode:      nil,
				Location:        nil,
				IsHeadOffice:    false,
				CreatedAt:       baseTime,
				UpdatedAt:       baseTime,
			},
			base2: &domain.Base{
				ID:              uuid.New(),
				CorporationID:   corporationID,
				CorporateNumber: corporateNumber,
				BaseName:        stringPtr("本社B"),
				PostalCode:      nil,
				Location:        nil,
				IsHeadOffice:    false,
				CreatedAt:       baseTime,
				UpdatedAt:       baseTime,
			},
			expected: false,
		},
		{
			name: "Empty strings vs nil - should not be similar",
			base1: &domain.Base{
				ID:              uuid.New(),
				CorporationID:   corporationID,
				CorporateNumber: corporateNumber,
				BaseName:        stringPtr("本社A"), // Different base names
				PostalCode:      stringPtr(""),
				Location:        stringPtr(""),
				IsHeadOffice:    false,
				CreatedAt:       baseTime,
				UpdatedAt:       baseTime,
			},
			base2: &domain.Base{
				ID:              uuid.New(),
				CorporationID:   corporationID,
				CorporateNumber: corporateNumber,
				BaseName:        stringPtr("本社B"), // Different base names
				PostalCode:      nil,
				Location:        nil,
				IsHeadOffice:    false,
				CreatedAt:       baseTime,
				UpdatedAt:       baseTime,
			},
			expected: false,
		},
		{
			name: "Exact same base names with different companies - should not be similar (different companies)",
			base1: &domain.Base{
				ID:              uuid.New(),
				CorporationID:   corporationID,
				CorporateNumber: corporateNumber,
				BaseName:        stringPtr("営業所"),
				PostalCode:      stringPtr("100-0001"),
				Location:        stringPtr("東京都千代田区千代田1-1"),
				IsHeadOffice:    false,
				CreatedAt:       baseTime,
				UpdatedAt:       baseTime,
			},
			base2: &domain.Base{
				ID:              uuid.New(),
				CorporationID:   uuid.New(), // Different corporation
				CorporateNumber: "9999999999999",
				BaseName:        stringPtr("営業所"),
				PostalCode:      stringPtr("100-0001"),
				Location:        stringPtr("東京都千代田区千代田1-1"),
				IsHeadOffice:    false,
				CreatedAt:       baseTime,
				UpdatedAt:       baseTime,
			},
			expected: false, // Different companies should not be considered similar
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := uc.IsSimilarBase(tt.base1, tt.base2)
			if result != tt.expected {
				t.Errorf("IsSimilarBase() = %v, expected %v", result, tt.expected)
				t.Logf("Base1: Name=%v, PostalCode=%v, Location=%v, IsHeadOffice=%v",
					getStringValue(tt.base1.BaseName),
					getStringValue(tt.base1.PostalCode),
					getStringValue(tt.base1.Location),
					tt.base1.IsHeadOffice)
				t.Logf("Base2: Name=%v, PostalCode=%v, Location=%v, IsHeadOffice=%v",
					getStringValue(tt.base2.BaseName),
					getStringValue(tt.base2.PostalCode),
					getStringValue(tt.base2.Location),
					tt.base2.IsHeadOffice)
			}
		})
	}
}

func TestFetchBasesUseCase_normalizePostalCode(t *testing.T) {
	uc := &FetchBasesUseCase{}

	tests := []struct {
		name     string
		input    *string
		expected string
	}{
		{
			name:     "Postal code with hyphen",
			input:    stringPtr("153-8636"),
			expected: "1538636",
		},
		{
			name:     "Postal code without hyphen",
			input:    stringPtr("1530061"),
			expected: "1530061",
		},
		{
			name:     "Postal code with spaces",
			input:    stringPtr("153 8636"),
			expected: "1538636",
		},
		{
			name:     "Postal code with hyphen and spaces",
			input:    stringPtr(" 153-8636 "),
			expected: "1538636",
		},
		{
			name:     "Empty string",
			input:    stringPtr(""),
			expected: "",
		},
		{
			name:     "Nil input",
			input:    nil,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := uc.normalizePostalCode(tt.input)
			if result != tt.expected {
				t.Errorf("normalizePostalCode() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestFetchBasesUseCase_normalizeLocation(t *testing.T) {
	uc := &FetchBasesUseCase{}

	tests := []struct {
		name     string
		input    *string
		expected string
	}{
		{
			name:     "Location with Japanese numbers",
			input:    stringPtr("東京都目黒区中目黒２丁目９番１３号"),
			expected: "東京都目黒区中目黒2-9-13",
		},
		{
			name:     "Location with Arabic numbers",
			input:    stringPtr("東京都目黒区中目黒2-9-13"),
			expected: "東京都目黒区中目黒2-9-13",
		},
		{
			name:     "Location with postal code prefix",
			input:    stringPtr("1530061 東京都目黒区中目黒２丁目９番１３号"),
			expected: "東京都目黒区中目黒2-9-13",
		},
		{
			name:     "Location with spaces",
			input:    stringPtr("東京都 目黒区 中目黒 2-9-13"),
			expected: "東京都目黒区中目黒2-9-13",
		},
		{
			name:     "Empty string",
			input:    stringPtr(""),
			expected: "",
		},
		{
			name:     "Nil input",
			input:    nil,
			expected: "",
		},
		{
			name:     "Location with mixed Japanese and Arabic numbers",
			input:    stringPtr("東京都目黒区中目黒２丁目9番13号"),
			expected: "東京都目黒区中目黒2-9-13",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := uc.normalizeLocation(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeLocation() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestFetchBasesUseCase_normalizeBaseName(t *testing.T) {
	uc := &FetchBasesUseCase{}

	tests := []struct {
		name     string
		input    *string
		expected string
	}{
		{
			name:     "本店 to 本社",
			input:    stringPtr("本店"),
			expected: "本社",
		},
		{
			name:     "本部 to 本社",
			input:    stringPtr("本部"),
			expected: "本社",
		},
		{
			name:     "head office to 本社",
			input:    stringPtr("head office"),
			expected: "本社",
		},
		{
			name:     "Head Office to 本社",
			input:    stringPtr("Head Office"),
			expected: "本社",
		},
		{
			name:     "Already 本社",
			input:    stringPtr("本社"),
			expected: "本社",
		},
		{
			name:     "Branch office (no change)",
			input:    stringPtr("営業所"),
			expected: "営業所",
		},
		{
			name:     "Empty string",
			input:    stringPtr(""),
			expected: "",
		},
		{
			name:     "Nil input",
			input:    nil,
			expected: "",
		},
		{
			name:     "String with spaces",
			input:    stringPtr("  本店  "),
			expected: "本社",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := uc.normalizeBaseName(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeBaseName() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestFetchBasesUseCase_extractAddressParts(t *testing.T) {
	uc := &FetchBasesUseCase{}

	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Tokyo address",
			input:    "東京都目黒区中目黒2-9-13",
			expected: []string{"東京都"},
		},
		{
			name:     "Osaka address",
			input:    "大阪府大阪市北区梅田1-1-1",
			expected: []string{"大阪府"},
		},
		{
			name:     "Hokkaido address",
			input:    "北海道札幌市中央区北1条西1-1-1",
			expected: []string{"北海道"},
		},
		{
			name:     "Address without prefecture",
			input:    "梅田1-1-1",
			expected: []string{},
		},
		{
			name:     "Empty string",
			input:    "",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := uc.extractAddressParts(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("extractAddressParts() length = %v, expected %v", len(result), len(tt.expected))
				return
			}
			for i, part := range result {
				if part != tt.expected[i] {
					t.Errorf("extractAddressParts()[%d] = %v, expected %v", i, part, tt.expected[i])
				}
			}
		})
	}
}

// Helper function to get string value from pointer
func getStringValue(s *string) string {
	if s == nil {
		return "<nil>"
	}
	return *s
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
