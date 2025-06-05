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
	openAIAPIURL       = "https://api.openai.com/v1/chat/completions"
	defaultModel       = "gpt-3.5-turbo"
	defaultMaxTokens   = 2048
	defaultTemperature = 0.3
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
	systemPrompt := `You are a helpful assistant that finds official company website URLs containing office/location information. 
You should respond with a JSON object containing an array of URLs that are likely to contain office or branch location information for the given company.
Only include official company website URLs, not third-party sites.
Format your response as JSON: {"urls": ["url1", "url2", ...]}
If no URLs are found, return {"urls": []}`

	userPrompt := fmt.Sprintf(`Find official website URLs for the company "%s" that contain information about their offices, branches, or locations. 
Look for pages like:
- Office locations page
- About us / Company info page with office details
- Contact page with multiple office addresses
- Branch/store locator pages

Company name: %s`, corporationName, corporationName)

	response, err := c.chatCompletion(ctx, systemPrompt, userPrompt)
	if err != nil {
		return nil, fmt.Errorf("failed to get URL discovery response: %w", err)
	}

	// Parse the response
	var result domain.URLDiscoveryResult
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		// If JSON parsing fails, try to extract URLs manually
		urls := c.extractURLsFromText(response)
		result = domain.URLDiscoveryResult{URLs: urls}
	}

	return &result, nil
}

// ExtractBasesFromURL extracts base information from a given URL
func (c *OpenAIClient) ExtractBasesFromURL(ctx context.Context, url string) (*domain.BaseExtractionResult, error) {
	// First, fetch the content of the URL
	content, err := c.fetchURLContent(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL content: %w", err)
	}

	systemPrompt := `You are a helpful assistant that extracts office/branch location information from website content.
Extract all office and branch locations from the provided text and return them as a JSON object.
For each location, determine if it's a head office or branch office based on context.
Format your response as JSON: {"bases": [{"name": "Office Name", "type": "head_office" or "branch_office", "postal_code": "123-4567", "prefecture": "Prefecture", "city": "City", "address": "Full Address", "phone_number": "Phone", "fax_number": "Fax"}]}
If no locations are found, return {"bases": []}`

	userPrompt := fmt.Sprintf(`Extract all office and branch location information from this website content.
Look for:
- Office names
- Addresses (postal codes, prefectures, cities, street addresses)
- Contact information (phone, fax)
- Determine if each location is a head office or branch office

Website URL: %s
Content: %s`, url, content)

	response, err := c.chatCompletion(ctx, systemPrompt, userPrompt)
	if err != nil {
		return nil, fmt.Errorf("failed to get base extraction response: %w", err)
	}

	// Parse the response
	var result domain.BaseExtractionResult
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return &domain.BaseExtractionResult{Bases: []domain.ExtractedBase{}}, nil
	}

	return &result, nil
}

// chatCompletion makes a chat completion request to OpenAI API
func (c *OpenAIClient) chatCompletion(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	req := domain.OpenAIRequest{
		Model: c.model,
		Messages: []domain.OpenAIMessage{
			{
				Role:    "system",
				Content: systemPrompt,
			},
			{
				Role:    "user",
				Content: userPrompt,
			},
		},
		MaxTokens:   defaultMaxTokens,
		Temperature: defaultTemperature,
	}

	reqBody, err := json.Marshal(req)
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
		var errorResp domain.OpenAIErrorResponse
		if err := json.Unmarshal(body, &errorResp); err == nil {
			return "", fmt.Errorf("OpenAI API error: %s", errorResp.Error.Message)
		}
		return "", fmt.Errorf("OpenAI API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var openAIResp domain.OpenAIResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(openAIResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return openAIResp.Choices[0].Message.Content, nil
}

// fetchURLContent fetches the content of a URL (simplified implementation)
func (c *OpenAIClient) fetchURLContent(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Corporation DB Bot/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch URL, status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Limit content size to avoid token limits
	content := string(body)
	if len(content) > 8000 {
		content = content[:8000] + "..."
	}

	return content, nil
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
