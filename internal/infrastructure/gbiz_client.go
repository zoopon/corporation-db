package infrastructure

import (
	"archive/zip"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"corporation-db/internal/domain"
)

// GBizClient handles communication with gBizINFO API
type GBizClient struct {
	httpClient *http.Client
	baseURL    string
}

// NewGBizClient creates a new gBizINFO client
func NewGBizClient() *GBizClient {
	return &GBizClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Minute, // Large files need longer timeout
		},
		baseURL: "https://info.gbiz.go.jp",
	}
}

// DownloadBasicInfoCSV downloads the basic corporation information CSV ZIP file
func (c *GBizClient) DownloadBasicInfoCSV(ctx context.Context) (string, error) {
	// gBizINFO basic info CSV URL (UTF-8 version)
	url := c.baseURL + "/hojin/DownloadToCSV?type=01&format=csv&encoding=utf8"

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to download CSV: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Create temporary file
	tmpFile, err := os.CreateTemp("", "gbiz_basic_info_*.zip")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close()

	// Copy response body to file
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write response to file: %w", err)
	}

	return tmpFile.Name(), nil
}

// ExtractAndParseCSV extracts the ZIP file and parses the CSV content
func (c *GBizClient) ExtractAndParseCSV(zipPath string) ([]*domain.CreateCorporationRequest, error) {
	// Open ZIP file
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open ZIP file: %w", err)
	}
	defer reader.Close()

	var corporations []*domain.CreateCorporationRequest

	// Find and process CSV files in the ZIP
	for _, file := range reader.File {
		if strings.HasSuffix(strings.ToLower(file.Name), ".csv") {
			fmt.Printf("Processing CSV file: %s\n", file.Name)

			fileReader, err := file.Open()
			if err != nil {
				return nil, fmt.Errorf("failed to open CSV file %s: %w", file.Name, err)
			}
			defer fileReader.Close()

			csvCorporations, err := c.parseCSV(fileReader)
			if err != nil {
				return nil, fmt.Errorf("failed to parse CSV file %s: %w", file.Name, err)
			}

			corporations = append(corporations, csvCorporations...)
		}
	}

	return corporations, nil
}

// parseCSV parses CSV content and converts to Corporation entities
func (c *GBizClient) parseCSV(reader io.Reader) ([]*domain.CreateCorporationRequest, error) {
	csvReader := csv.NewReader(reader)

	// Read header
	headers, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	fmt.Printf("CSV Headers: %v\n", headers)

	var corporations []*domain.CreateCorporationRequest
	lineNum := 1

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Warning: failed to read line %d: %v\n", lineNum, err)
			lineNum++
			continue
		}

		// Skip if not enough columns
		if len(record) < 3 {
			lineNum++
			continue
		}

		corp, err := c.parseCSVRecord(headers, record)
		if err != nil {
			fmt.Printf("Warning: failed to parse line %d: %v\n", lineNum, err)
			lineNum++
			continue
		}

		if corp != nil {
			corporations = append(corporations, corp)
		}
		lineNum++

		// Progress indicator for large files
		if lineNum%10000 == 0 {
			fmt.Printf("Processed %d lines...\n", lineNum)
		}
	}

	fmt.Printf("Successfully parsed %d corporations from CSV\n", len(corporations))
	return corporations, nil
}

// parseCSVRecord parses a single CSV record into Corporation
func (c *GBizClient) parseCSVRecord(headers []string, record []string) (*domain.CreateCorporationRequest, error) {
	// Create a map for easier field access
	fieldMap := make(map[string]string)
	for i, header := range headers {
		if i < len(record) {
			fieldMap[header] = strings.TrimSpace(record[i])
		}
	}

	// Extract corporate number (assumed to be in first column or specific header)
	corporateNumber := ""
	if len(record) > 0 {
		corporateNumber = strings.TrimSpace(record[0])
	}

	// Validate corporate number (should be 13 digits)
	if len(corporateNumber) != 13 {
		return nil, fmt.Errorf("invalid corporate number: %s", corporateNumber)
	}

	// Extract name (assumed to be in second column or specific header)
	name := ""
	if len(record) > 1 {
		name = strings.TrimSpace(record[1])
	}

	if name == "" {
		return nil, fmt.Errorf("empty corporation name")
	}

	corp := &domain.CreateCorporationRequest{
		CorporateNumber: corporateNumber,
		Name:            name,
		Status:          "活動中", // Default status
	}

	// Parse optional fields based on common gBizINFO CSV structure
	// Note: Actual field mapping should be adjusted based on real CSV format

	if val, exists := fieldMap["法人名カナ"]; exists && val != "" {
		corp.NameKana = &val
	}

	if val, exists := fieldMap["英語法人名"]; exists && val != "" {
		corp.EnglishName = &val
	}

	if val, exists := fieldMap["郵便番号"]; exists && val != "" {
		corp.PostalCode = &val
	}

	if val, exists := fieldMap["所在地"]; exists && val != "" {
		corp.Address = &val
	}

	if val, exists := fieldMap["都道府県コード"]; exists && val != "" {
		corp.PrefectureCode = &val
	}

	if val, exists := fieldMap["市区町村コード"]; exists && val != "" {
		corp.CityCode = &val
	}

	if val, exists := fieldMap["法人種別"]; exists && val != "" {
		corp.CorporateType = &val
	}

	if val, exists := fieldMap["代表者名"]; exists && val != "" {
		corp.Representative = &val
	}

	if val, exists := fieldMap["事業内容"]; exists && val != "" {
		corp.BusinessDescription = &val
	}

	if val, exists := fieldMap["業種"]; exists && val != "" {
		corp.Industry = &val
	}

	// Parse numeric fields
	if val, exists := fieldMap["資本金"]; exists && val != "" {
		if capital, err := strconv.ParseInt(val, 10, 64); err == nil {
			corp.CapitalStock = &capital
		}
	}

	if val, exists := fieldMap["従業員数"]; exists && val != "" {
		if employees, err := strconv.ParseInt(val, 10, 32); err == nil {
			employeeNum := int32(employees)
			corp.EmployeeNumber = &employeeNum
		}
	}

	// Parse date fields
	if val, exists := fieldMap["設立年月日"]; exists && val != "" {
		if date, err := time.Parse("2006-01-02", val); err == nil {
			corp.FoundingDate = &date
		} else if date, err := time.Parse("2006/01/02", val); err == nil {
			corp.FoundingDate = &date
		}
	}

	return corp, nil
}

// Cleanup removes temporary files
func (c *GBizClient) Cleanup(filePath string) error {
	if filePath != "" {
		return os.Remove(filePath)
	}
	return nil
}
