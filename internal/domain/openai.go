package domain

import "context"

// OpenAIRequest represents a request to OpenAI API
type OpenAIRequest struct {
	Model       string          `json:"model"`
	Messages    []OpenAIMessage `json:"messages"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Temperature float64         `json:"temperature,omitempty"`
	Stream      bool            `json:"stream,omitempty"`
}

// OpenAIMessage represents a message in OpenAI API request
type OpenAIMessage struct {
	Role    string `json:"role"` // "system", "user", "assistant"
	Content string `json:"content"`
}

// OpenAIResponse represents a response from OpenAI API
type OpenAIResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []OpenAIChoice `json:"choices"`
	Usage   OpenAIUsage    `json:"usage"`
}

// OpenAIChoice represents a choice in OpenAI API response
type OpenAIChoice struct {
	Index        int           `json:"index"`
	Message      OpenAIMessage `json:"message"`
	FinishReason string        `json:"finish_reason"`
}

// OpenAIUsage represents token usage in OpenAI API response
type OpenAIUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// OpenAIErrorResponse represents an error response from OpenAI API
type OpenAIErrorResponse struct {
	Error OpenAIError `json:"error"`
}

// OpenAIError represents an error from OpenAI API
type OpenAIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code,omitempty"`
}

// URLDiscoveryResult represents the result of URL discovery
type URLDiscoveryResult struct {
	URLs          []string `json:"urls"`
	CandidateURLs []string `json:"candidate_urls,omitempty"`
}

// BaseExtractionResult represents the result of base information extraction
type BaseExtractionResult struct {
	Bases []ExtractedBase `json:"bases"`
}

// ExtractedBase represents base information extracted from a URL
type ExtractedBase struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"` // "head_office", "branch_office"
	PostalCode  string  `json:"postal_code"`
	Prefecture  string  `json:"prefecture"`
	City        string  `json:"city"`
	Address     string  `json:"address"`
	Latitude    float64 `json:"latitude,omitempty"`
	Longitude   float64 `json:"longitude,omitempty"`
	PhoneNumber string  `json:"phone_number,omitempty"`
	FaxNumber   string  `json:"fax_number,omitempty"`
}

// OpenAIService defines the interface for OpenAI API operations
type OpenAIService interface {
	// DiscoverURLs discovers URLs containing office information for a corporation
	DiscoverURLs(ctx context.Context, corporationName string) (*URLDiscoveryResult, error)

	// ExtractBasesFromURL extracts base information from given URLs
	ExtractBasesFromURL(ctx context.Context, urls []string, corporationName string) (*BaseExtractionResult, error)
}
