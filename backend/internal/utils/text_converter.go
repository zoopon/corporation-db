package utils

import (
	"strings"
	"unicode"
)

// TextConverter provides text normalization functions for search
type TextConverter struct{}

// NewTextConverter creates a new TextConverter instance
func NewTextConverter() *TextConverter {
	return &TextConverter{}
}

// NormalizeForSearch normalizes text for search purposes
// - Converts full-width Roman characters to half-width
// - Converts hiragana to katakana
// - Converts to lowercase for Roman characters
func (tc *TextConverter) NormalizeForSearch(text string) string {
	if text == "" {
		return ""
	}

	var result strings.Builder
	for _, r := range text {
		normalized := tc.normalizeRune(r)
		result.WriteRune(normalized)
	}

	return result.String()
}

// normalizeRune normalizes a single rune
func (tc *TextConverter) normalizeRune(r rune) rune {
	// Convert full-width Roman characters to half-width
	if r >= '０' && r <= '９' {
		return r - '０' + '0' // Full-width digits to half-width
	}
	if r >= 'Ａ' && r <= 'Ｚ' {
		return unicode.ToLower(r - 'Ａ' + 'A') // Full-width uppercase to half-width lowercase
	}
	if r >= 'ａ' && r <= 'ｚ' {
		return r - 'ａ' + 'a' // Full-width lowercase to half-width
	}

	// Convert hiragana to katakana
	if r >= 'あ' && r <= 'ん' {
		return r - 'あ' + 'ア'
	}

	// Handle special hiragana characters
	switch r {
	case 'ゃ':
		return 'ャ'
	case 'ゅ':
		return 'ュ'
	case 'ょ':
		return 'ョ'
	case 'っ':
		return 'ッ'
	case 'ー':
		return 'ー' // Keep as is
	}

	// Convert ASCII uppercase to lowercase
	if r >= 'A' && r <= 'Z' {
		return unicode.ToLower(r)
	}

	// Return as is for other characters
	return r
}

// IsHiragana checks if the character is hiragana
func (tc *TextConverter) IsHiragana(r rune) bool {
	return (r >= 'あ' && r <= 'ん') || r == 'ゃ' || r == 'ゅ' || r == 'ょ' || r == 'っ'
}

// IsKatakana checks if the character is katakana
func (tc *TextConverter) IsKatakana(r rune) bool {
	return (r >= 'ア' && r <= 'ン') || r == 'ャ' || r == 'ュ' || r == 'ョ' || r == 'ッ' || r == 'ー'
}

// IsFullWidthRoman checks if the character is full-width Roman
func (tc *TextConverter) IsFullWidthRoman(r rune) bool {
	return (r >= '０' && r <= '９') || (r >= 'Ａ' && r <= 'Ｚ') || (r >= 'ａ' && r <= 'ｚ')
}
