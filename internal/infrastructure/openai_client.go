package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"corporation-db/internal/domain"
)

const (
	openAIAPIURL       = "https://api.openai.com/v1/responses"
	defaultModel       = "gpt-4o"
	defaultMaxTokens   = 8192
	defaultTemperature = 0.0
)

// OpenAIClient implements OpenAI API operations
type OpenAIClient struct {
	httpClient *http.Client
	apiKey     string
	model      string
}

// NewOpenAIClient creates a new OpenAI API client
func NewOpenAIClient() (*OpenAIClient, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable is required")
	}

	return &OpenAIClient{
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		apiKey: apiKey,
		model:  defaultModel,
	}, nil
}

// DiscoverURLs discovers URLs containing office information for a corporation
func (c *OpenAIClient) DiscoverURLs(ctx context.Context, corporationName string) (*domain.URLDiscoveryResult, error) {
	systemPrompt := "You are a web assistant who identifies pages on a company's official website that list its offices or facilities, including HQ, plants, sales offices, branches, warehouses, service centers, etc."

	userPrompt := fmt.Sprintf(`法人番号を持つ会社名「%s」から、公式Webサイト内の拠点情報ページのURLと、探索過程で収集したすべての公式URLをJSON形式で返してください。拠点情報には、本社、生産拠点、工場、営業所、支社、支店、倉庫、サービスセンター（保守拠点なども含む）などを含みます。`, corporationName)

	response, err := c.advancedChatCompletion(ctx, systemPrompt, userPrompt, "urls_including_locations", getURLDiscoverySchema())
	if err != nil {
		return nil, fmt.Errorf("failed to get URL discovery response: %w", err)
	}

	// Parse the response
	var rawResult struct {
		URLs          []string `json:"urls"`
		CandidateURLs []string `json:"candidate_urls"`
	}

	if err := json.Unmarshal([]byte(response), &rawResult); err != nil {
		// If JSON parsing fails, try to extract URLs manually
		urls := c.extractURLsFromText(response)
		return &domain.URLDiscoveryResult{URLs: urls}, nil
	}

	result := &domain.URLDiscoveryResult{
		URLs:          rawResult.URLs,
		CandidateURLs: rawResult.CandidateURLs,
	}

	return result, nil
}

// ExtractBasesFromURL extracts base information from given URLs
func (c *OpenAIClient) ExtractBasesFromURL(ctx context.Context, urls []string, corporationName string) (*domain.BaseExtractionResult, error) {
	systemPrompt := "あなたは与えられた会社情報とWebページから拠点情報（本社、生産拠点、製造拠点、営業所、事業所、支社、工場、サービスセンターなど）を取得する作業者です"

	// Create URLs list string
	urlsList := strings.Join(urls, "\n")

	userPrompt := fmt.Sprintf(`次のURLsから『%s』の拠点情報（本社、生産拠点、製造拠点、営業所、事業所、支社、工場、サービスセンターなど）を取得してjson形式で返してください。
国内の全ての拠点を漏れなく取得してください。ただし、関連会社等、当該会社以外の情報は含めないでください。

URLs:
%s`, corporationName, urlsList)

	response, err := c.advancedChatCompletion(ctx, systemPrompt, userPrompt, "location_information", getLocationExtractionSchema())
	if err != nil {
		return nil, fmt.Errorf("failed to get base extraction response: %w", err)
	}

	// Parse the response
	var rawResult struct {
		Locations []struct {
			URL         string `json:"url"`
			Name        string `json:"name"`
			Zipcode     string `json:"zipcode"`
			Address     string `json:"address"`
			PhoneNumber string `json:"phone_number"`
		} `json:"locations"`
	}

	if err := json.Unmarshal([]byte(response), &rawResult); err != nil {
		return &domain.BaseExtractionResult{Bases: []domain.ExtractedBase{}}, nil
	}

	// Convert to domain format
	var bases []domain.ExtractedBase
	for _, location := range rawResult.Locations {
		// Clean name by removing numbers at start/end
		cleanName := strings.TrimSpace(location.Name)
		if len(cleanName) > 2 {
			// Remove leading/trailing 2-digit numbers
			if len(cleanName) >= 2 && cleanName[0] >= '0' && cleanName[0] <= '9' && cleanName[1] >= '0' && cleanName[1] <= '9' {
				cleanName = cleanName[2:]
			}
			if len(cleanName) >= 2 && cleanName[len(cleanName)-2] >= '0' && cleanName[len(cleanName)-2] <= '9' && cleanName[len(cleanName)-1] >= '0' && cleanName[len(cleanName)-1] <= '9' {
				cleanName = cleanName[:len(cleanName)-2]
			}
			cleanName = strings.TrimSpace(cleanName)
		}

		// Determine office type based on name
		officeType := "branch_office"
		if strings.Contains(cleanName, "本社") || strings.Contains(cleanName, "本店") || strings.Contains(cleanName, "Head") {
			officeType = "head_office"
		}

		base := domain.ExtractedBase{
			Name:        cleanName,
			Type:        officeType,
			PostalCode:  location.Zipcode,
			Prefecture:  "", // Will be extracted from address
			City:        "", // Will be extracted from address
			Address:     location.Address,
			PhoneNumber: location.PhoneNumber,
			FaxNumber:   "", // Not provided in this format
		}

		// Try to parse prefecture and city from address
		if location.Address != "" {
			parts := strings.Split(location.Address, " ")
			if len(parts) >= 2 {
				base.Prefecture = parts[0]
				base.City = parts[1]
			}
		}

		bases = append(bases, base)
	}

	return &domain.BaseExtractionResult{Bases: bases}, nil
}

// advancedChatCompletion makes a request to OpenAI API with tools and structured output
func (c *OpenAIClient) advancedChatCompletion(ctx context.Context, systemPrompt, userPrompt, schemaName string, schema map[string]interface{}) (string, error) {
	requestBody := map[string]interface{}{
		"model": c.model,
		"tools": []map[string]interface{}{
			{
				"type":                "web_search_preview",
				"search_context_size": "high",
			},
		},
		"input": []map[string]interface{}{
			{
				"role":    "system",
				"content": systemPrompt,
			},
			{
				"role":    "user",
				"content": userPrompt,
			},
		},
		"temperature":       defaultTemperature,
		"max_output_tokens": defaultMaxTokens,
		"text": map[string]interface{}{
			"format": map[string]interface{}{
				"type":   "json_schema",
				"name":   schemaName,
				"strict": true,
				"schema": schema,
			},
		},
	}

	reqBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", openAIAPIURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OpenAI API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse advanced API response format
	var response struct {
		Output []struct {
			Type    string `json:"type"`
			Content []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"content"`
		} `json:"output"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Extract the text content
	for _, output := range response.Output {
		if output.Type == "message" {
			for _, content := range output.Content {
				if content.Type == "output_text" {
					// Clean up control characters but keep \n and \t
					cleaned := strings.Map(func(r rune) rune {
						if r >= 0 && r <= 31 && r != '\n' && r != '\t' {
							return -1
						}
						return r
					}, content.Text)
					return cleaned, nil
				}
			}
		}
	}

	return "", fmt.Errorf("no valid content found in response")
}

// getURLDiscoverySchema returns the JSON schema for URL discovery
func getURLDiscoverySchema() map[string]interface{} {
	return map[string]interface{}{
		"type":     "object",
		"required": []string{"urls", "candidate_urls"},
		"properties": map[string]interface{}{
			"urls": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type":        "string",
					"description": "拠点情報が記載されているページと判断したページのURL",
				},
				"description": "拠点情報が記載されているページと判断したページのURLの一覧",
			},
			"candidate_urls": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type":        "string",
					"description": "探索で取得したURL",
				},
				"description": "探索で取得した全てのURLの一覧",
			},
		},
		"additionalProperties": false,
	}
}

// getLocationExtractionSchema returns the JSON schema for location extraction
func getLocationExtractionSchema() map[string]interface{} {
	return map[string]interface{}{
		"type":     "object",
		"required": []string{"locations"},
		"properties": map[string]interface{}{
			"locations": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type":     "object",
					"required": []string{"url", "name", "zipcode", "address", "phone_number"},
					"properties": map[string]interface{}{
						"url": map[string]interface{}{
							"type":        "string",
							"description": "拠点情報を取得したWebページのパスを含むURL",
						},
						"name": map[string]interface{}{
							"type":        "string",
							"description": "拠点の名前",
						},
						"zipcode": map[string]interface{}{
							"type":        "string",
							"description": "拠点の郵便番号",
						},
						"address": map[string]interface{}{
							"type":        "string",
							"description": "郵便番号を含まない拠点の住所",
						},
						"phone_number": map[string]interface{}{
							"type":        "string",
							"description": "拠点の電話番号",
						},
					},
					"additionalProperties": false,
				},
				"description": "An array of base information objects.",
			},
		},
		"additionalProperties": false,
	}
}

// extractURLsFromText extracts URLs from text response when JSON parsing fails
func (c *OpenAIClient) extractURLsFromText(text string) []string {
	var urls []string
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "http://") || strings.HasPrefix(line, "https://") {
			urls = append(urls, line)
		}
		// Also check for quoted URLs
		if strings.Contains(line, "http") {
			parts := strings.Split(line, "\"")
			for _, part := range parts {
				if strings.HasPrefix(part, "http://") || strings.HasPrefix(part, "https://") {
					urls = append(urls, part)
				}
			}
		}
	}

	return urls
}
