package utils

import (
	"testing"
)

func TestTextConverter_NormalizeForSearch(t *testing.T) {
	converter := NewTextConverter()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Full-width digits to half-width",
			input:    "１２３４５",
			expected: "12345",
		},
		{
			name:     "Full-width uppercase to half-width lowercase",
			input:    "ＡＢＣＤＥ",
			expected: "abcde",
		},
		{
			name:     "Full-width lowercase to half-width",
			input:    "ａｂｃｄｅ",
			expected: "abcde",
		},
		{
			name:     "Hiragana to katakana",
			input:    "あいうえお",
			expected: "アイウエオ",
		},
		{
			name:     "Hiragana with small characters",
			input:    "きゃきゅきょ",
			expected: "キャキュキョ",
		},
		{
			name:     "Hiragana with sokuon",
			input:    "がっこう",
			expected: "ガッコウ",
		},
		{
			name:     "ASCII uppercase to lowercase",
			input:    "ABCDE",
			expected: "abcde",
		},
		{
			name:     "Mixed Japanese and Roman",
			input:    "株式会社ＡＢＣＤ",
			expected: "株式会社abcd",
		},
		{
			name:     "Company name with hiragana and full-width",
			input:    "かぶしきがいしゃＴＯＹＯＴＡ",
			expected: "カブシキガイシャtoyota",
		},
		{
			name:     "Mixed case with numbers",
			input:    "ソニー１２３ＡＢＣＤ",
			expected: "ソニー123abcd",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Katakana should remain unchanged",
			input:    "カタカナ",
			expected: "カタカナ",
		},
		{
			name:     "Long vowel mark",
			input:    "コンピューター",
			expected: "コンピューター",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.NormalizeForSearch(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizeForSearch(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestTextConverter_IsHiragana(t *testing.T) {
	converter := NewTextConverter()

	tests := []struct {
		char     rune
		expected bool
	}{
		{'あ', true},
		{'ん', true},
		{'ゃ', true},
		{'っ', true},
		{'ア', false},
		{'A', false},
		{'1', false},
	}

	for _, tt := range tests {
		t.Run(string(tt.char), func(t *testing.T) {
			result := converter.IsHiragana(tt.char)
			if result != tt.expected {
				t.Errorf("IsHiragana(%q) = %v, expected %v", tt.char, result, tt.expected)
			}
		})
	}
}

func TestTextConverter_IsKatakana(t *testing.T) {
	converter := NewTextConverter()

	tests := []struct {
		char     rune
		expected bool
	}{
		{'ア', true},
		{'ン', true},
		{'ャ', true},
		{'ッ', true},
		{'ー', true},
		{'あ', false},
		{'A', false},
		{'1', false},
	}

	for _, tt := range tests {
		t.Run(string(tt.char), func(t *testing.T) {
			result := converter.IsKatakana(tt.char)
			if result != tt.expected {
				t.Errorf("IsKatakana(%q) = %v, expected %v", tt.char, result, tt.expected)
			}
		})
	}
}

func TestTextConverter_IsFullWidthRoman(t *testing.T) {
	converter := NewTextConverter()

	tests := []struct {
		char     rune
		expected bool
	}{
		{'０', true},
		{'９', true},
		{'Ａ', true},
		{'Ｚ', true},
		{'ａ', true},
		{'ｚ', true},
		{'0', false},
		{'A', false},
		{'あ', false},
	}

	for _, tt := range tests {
		t.Run(string(tt.char), func(t *testing.T) {
			result := converter.IsFullWidthRoman(tt.char)
			if result != tt.expected {
				t.Errorf("IsFullWidthRoman(%q) = %v, expected %v", tt.char, result, tt.expected)
			}
		})
	}
}
