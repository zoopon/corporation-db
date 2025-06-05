package infrastructure

import (
	"archive/zip"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"corporation-db/internal/domain"
	"corporation-db/internal/utils"
)

// prefectureCodeMap maps prefecture names to JIS X 0401 codes
var prefectureCodeMap = map[string]string{
	"北海道": "01",
	"青森県": "02", "青森": "02",
	"岩手県": "03", "岩手": "03",
	"宮城県": "04", "宮城": "04",
	"秋田県": "05", "秋田": "05",
	"山形県": "06", "山形": "06",
	"福島県": "07", "福島": "07",
	"茨城県": "08", "茨城": "08",
	"栃木県": "09", "栃木": "09",
	"群馬県": "10", "群馬": "10",
	"埼玉県": "11", "埼玉": "11",
	"千葉県": "12", "千葉": "12",
	"東京都": "13", "東京": "13",
	"神奈川県": "14", "神奈川": "14",
	"新潟県": "15", "新潟": "15",
	"富山県": "16", "富山": "16",
	"石川県": "17", "石川": "17",
	"福井県": "18", "福井": "18",
	"山梨県": "19", "山梨": "19",
	"長野県": "20", "長野": "20",
	"岐阜県": "21", "岐阜": "21",
	"静岡県": "22", "静岡": "22",
	"愛知県": "23", "愛知": "23",
	"三重県": "24", "三重": "24",
	"滋賀県": "25", "滋賀": "25",
	"京都府": "26", "京都": "26",
	"大阪府": "27", "大阪": "27",
	"兵庫県": "28", "兵庫": "28",
	"奈良県": "29", "奈良": "29",
	"和歌山県": "30", "和歌山": "30",
	"鳥取県": "31", "鳥取": "31",
	"島根県": "32", "島根": "32",
	"岡山県": "33", "岡山": "33",
	"広島県": "34", "広島": "34",
	"山口県": "35", "山口": "35",
	"徳島県": "36", "徳島": "36",
	"香川県": "37", "香川": "37",
	"愛媛県": "38", "愛媛": "38",
	"高知県": "39", "高知": "39",
	"福岡県": "40", "福岡": "40",
	"佐賀県": "41", "佐賀": "41",
	"長崎県": "42", "長崎": "42",
	"熊本県": "43", "熊本": "43",
	"大分県": "44", "大分": "44",
	"宮崎県": "45", "宮崎": "45",
	"鹿児島県": "46", "鹿児島": "46",
	"沖縄県": "47", "沖縄": "47",
}

// cityCodeMap maps major city names to prefecture codes
var cityCodeMap = map[string]string{
	"札幌市":   "01", // 北海道
	"青森市":   "02", // 青森県
	"盛岡市":   "03", // 岩手県
	"仙台市":   "04", // 宮城県
	"秋田市":   "05", // 秋田県
	"山形市":   "06", // 山形県
	"福島市":   "07", // 福島県
	"水戸市":   "08", // 茨城県
	"宇都宮市":  "09", // 栃木県
	"前橋市":   "10", // 群馬県
	"さいたま市": "11", // 埼玉県
	"千葉市":   "12", // 千葉県
	"新宿区":   "13", // 東京都
	"渋谷区":   "13", // 東京都
	"港区":    "13", // 東京都
	"中央区":   "13", // 東京都 (note: 中央区 exists in multiple prefectures, Tokyo is most common)
	"横浜市":   "14", // 神奈川県
	"川崎市":   "14", // 神奈川県
	"新潟市":   "15", // 新潟県
	"富山市":   "16", // 富山県
	"金沢市":   "17", // 石川県
	"福井市":   "18", // 福井県
	"甲府市":   "19", // 山梨県
	"長野市":   "20", // 長野県
	"岐阜市":   "21", // 岐阜県
	"静岡市":   "22", // 静岡県
	"浜松市":   "22", // 静岡県
	"名古屋市":  "23", // 愛知県
	"津市":    "24", // 三重県
	"大津市":   "25", // 滋賀県
	"京都市":   "26", // 京都府
	"大阪市":   "27", // 大阪府
	"堺市":    "27", // 大阪府
	"神戸市":   "28", // 兵庫県
	"奈良市":   "29", // 奈良県
	"和歌山市":  "30", // 和歌山県
	"鳥取市":   "31", // 鳥取県
	"松江市":   "32", // 島根県
	"岡山市":   "33", // 岡山県
	"広島市":   "34", // 広島県
	"山口市":   "35", // 山口県
	"徳島市":   "36", // 徳島県
	"高松市":   "37", // 香川県
	"松山市":   "38", // 愛媛県
	"高知市":   "39", // 高知県
	"福岡市":   "40", // 福岡県
	"北九州市":  "40", // 福岡県
	"佐賀市":   "41", // 佐賀県
	"長崎市":   "42", // 長崎県
	"熊本市":   "43", // 熊本県
	"大分市":   "44", // 大分県
	"宮崎市":   "45", // 宮崎県
	"鹿児島市":  "46", // 鹿児島県
	"那覇市":   "47", // 沖縄県
}

// GBizClient handles communication with gBizINFO API
type GBizClient struct {
	httpClient    *http.Client
	baseURL       string
	textConverter *utils.TextConverter
}

// NewGBizClient creates a new gBizINFO client
func NewGBizClient() *GBizClient {
	return &GBizClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Minute, // Large files need longer timeout
		},
		baseURL:       "https://info.gbiz.go.jp",
		textConverter: utils.NewTextConverter(),
	}
}

// DownloadBasicInfoCSV downloads the basic corporation information CSV ZIP file
func (c *GBizClient) DownloadBasicInfoCSV(ctx context.Context) (string, error) {
	// gBizINFO download endpoint (using POST as confirmed by browser dev tools)
	url := c.baseURL + "/hojin/Download"

	// Prepare form data for POST request using actual gBizINFO parameters
	formData := "downfile=7&downtype=zip&downenc=UTF-8"

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(formData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set appropriate headers for form submission
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")
	req.Header.Set("Accept", "application/zip,*/*")
	req.Header.Set("Referer", "https://info.gbiz.go.jp/hojin/DownloadTop")

	log.Printf("Requesting gBizINFO data from: %s", url)
	log.Printf("Request body: %s", formData)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to download CSV: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("Response status: %d %s", resp.StatusCode, resp.Status)
	log.Printf("Response headers: %v", resp.Header)

	// Check if response is actually a ZIP file
	contentType := resp.Header.Get("Content-Type")
	log.Printf("Content-Type: %s", contentType)

	if resp.StatusCode != http.StatusOK {
		// Read error response body for debugging
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(body))
	}

	// If content type suggests HTML instead of ZIP, read the response for debugging
	if strings.Contains(contentType, "text/html") {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read HTML response: %w", err)
		}
		log.Printf("Received HTML response (first 500 chars): %s", string(body[:min(500, len(body))]))
		return "", fmt.Errorf("received HTML response instead of ZIP file - this may indicate incorrect parameters or authentication required")
	}

	// Create temporary file
	tmpFile, err := os.CreateTemp("", "gbiz_basic_info_*.zip")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close()

	// Copy response body to file
	written, err := io.Copy(tmpFile, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write response to file: %w", err)
	}

	log.Printf("Downloaded %d bytes to %s", written, tmpFile.Name())
	return tmpFile.Name(), nil
}

// DownloadFinanceData downloads the finance information CSV ZIP file
func (c *GBizClient) DownloadFinanceData(ctx context.Context) (io.Reader, error) {
	// gBizINFO finance download endpoint
	url := c.baseURL + "/hojin/Download"

	// Prepare form data for finance data download (using different downfile parameter)
	// Note: downfile=13 is for finance data (Zaimu)
	formData := "downfile=13&downtype=zip&downenc=UTF-8"

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(formData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set appropriate headers for form submission
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")
	req.Header.Set("Accept", "application/zip,*/*")
	req.Header.Set("Referer", "https://info.gbiz.go.jp/hojin/DownloadTop")

	log.Printf("Requesting gBizINFO finance data from: %s", url)
	log.Printf("Request body: %s", formData)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download finance data: %w", err)
	}

	log.Printf("Response status: %d %s", resp.StatusCode, resp.Status)
	log.Printf("Response headers: %v", resp.Header)

	// Check if response is actually a ZIP file
	contentType := resp.Header.Get("Content-Type")
	log.Printf("Content-Type: %s", contentType)

	if resp.StatusCode != http.StatusOK {
		// Read error response body for debugging
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(body))
	}

	// If content type suggests HTML instead of ZIP, read the response for debugging
	if strings.Contains(contentType, "text/html") {
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read HTML response: %w", err)
		}
		log.Printf("Received HTML response (first 500 chars): %s", string(body[:min(500, len(body))]))
		return nil, fmt.Errorf("received HTML response instead of ZIP file - this may indicate incorrect parameters or authentication required")
	}

	return resp.Body, nil
}

// LoadTestCSVFile loads a local test CSV file for testing purposes
func (c *GBizClient) LoadTestCSVFile(csvPath string) ([]*domain.CreateCorporationRequest, error) {
	file, err := os.Open(csvPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open test CSV file: %w", err)
	}
	defer file.Close()

	return c.parseCSV(file)
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ExtractAndParseCSV extracts the ZIP file and parses the CSV content (legacy method)
// Deprecated: Use ExtractAndProcessCSVStream for memory-efficient processing
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

// ExtractAndProcessCSVStream extracts the ZIP file and processes CSV content in streaming batches
func (c *GBizClient) ExtractAndProcessCSVStream(zipPath string, batchProcessor func([]*domain.CreateCorporationRequest) error) (int, error) {
	// Open ZIP file
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return 0, fmt.Errorf("failed to open ZIP file: %w", err)
	}
	defer reader.Close()

	totalProcessed := 0

	// Find and process CSV files in the ZIP
	for _, file := range reader.File {
		if strings.HasSuffix(strings.ToLower(file.Name), ".csv") {
			fmt.Printf("Processing CSV file: %s\n", file.Name)

			fileReader, err := file.Open()
			if err != nil {
				return totalProcessed, fmt.Errorf("failed to open CSV file %s: %w", file.Name, err)
			}
			defer fileReader.Close()

			processed, err := c.parseCSVStream(fileReader, batchProcessor)
			if err != nil {
				return totalProcessed, fmt.Errorf("failed to process CSV file %s: %w", file.Name, err)
			}

			totalProcessed += processed
		}
	}

	return totalProcessed, nil
}

// ExtractAndProcessFinanceCSVStream extracts the ZIP file and processes finance CSV content in streaming batches
func (c *GBizClient) ExtractAndProcessFinanceCSVStream(zipPath string, batchProcessor func([]*domain.Finance) error) (int, error) {
	// Open ZIP file
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return 0, fmt.Errorf("failed to open ZIP file: %w", err)
	}
	defer reader.Close()

	totalProcessed := 0

	// Find and process CSV files in the ZIP
	for _, file := range reader.File {
		if strings.HasSuffix(strings.ToLower(file.Name), ".csv") {
			fmt.Printf("Processing finance CSV file: %s\n", file.Name)

			fileReader, err := file.Open()
			if err != nil {
				return totalProcessed, fmt.Errorf("failed to open CSV file %s: %w", file.Name, err)
			}
			defer fileReader.Close()

			processed, err := c.parseFinanceCSVStream(fileReader, batchProcessor)
			if err != nil {
				return totalProcessed, fmt.Errorf("failed to process finance CSV file %s: %w", file.Name, err)
			}

			totalProcessed += processed
		}
	}

	return totalProcessed, nil
}

// parseCSV parses CSV content and converts to Corporation entities (legacy method)
// Deprecated: Use parseCSVStream for memory-efficient processing
func (c *GBizClient) parseCSV(reader io.Reader) ([]*domain.CreateCorporationRequest, error) {
	csvReader := csv.NewReader(reader)

	// Configure CSV reader for gBizINFO format
	csvReader.LazyQuotes = true    // Allow lazy quotes for malformed CSV
	csvReader.FieldsPerRecord = -1 // Allow variable number of fields per record

	// Read header
	headers, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	fmt.Printf("CSV Headers: %v\n", headers)
	fmt.Printf("Number of columns: %d\n", len(headers))

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

// parseCSVStream parses CSV content in streaming batches to minimize memory usage
func (c *GBizClient) parseCSVStream(reader io.Reader, batchProcessor func([]*domain.CreateCorporationRequest) error) (int, error) {
	csvReader := csv.NewReader(reader)

	// Configure CSV reader for gBizINFO format
	csvReader.LazyQuotes = true    // Allow lazy quotes for malformed CSV
	csvReader.FieldsPerRecord = -1 // Allow variable number of fields per record

	// Read header
	headers, err := csvReader.Read()
	if err != nil {
		return 0, fmt.Errorf("failed to read CSV header: %w", err)
	}

	fmt.Printf("CSV Headers: %v\n", headers)
	fmt.Printf("Number of columns: %d\n", len(headers))

	const batchSize = 100 // Further reduced batch size for optimal memory efficiency
	batch := make([]*domain.CreateCorporationRequest, 0, batchSize)
	lineNum := 1
	totalProcessed := 0

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			// Process final batch if any records remain
			if len(batch) > 0 {
				if err := batchProcessor(batch); err != nil {
					return totalProcessed, fmt.Errorf("failed to process final batch at line %d: %w", lineNum, err)
				}
				totalProcessed += len(batch)
				fmt.Printf("Processed final batch of %d records (total: %d)\n", len(batch), totalProcessed)
			}
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
			batch = append(batch, corp)

			// Process batch when it reaches the target size
			if len(batch) >= batchSize {
				if err := batchProcessor(batch); err != nil {
					return totalProcessed, fmt.Errorf("failed to process batch at line %d: %w", lineNum, err)
				}
				totalProcessed += len(batch)
				fmt.Printf("Processed batch of %d records (total: %d)\n", len(batch), totalProcessed)

				// Clear batch for next iteration (reuse slice to minimize allocations)
				batch = batch[:0]
			}
		}
		lineNum++

		// Progress indicator for large files - reduced frequency to minimize I/O overhead
		if lineNum%100000 == 0 {
			fmt.Printf("Processed %d lines...\n", lineNum)
		}
	}

	fmt.Printf("Successfully processed %d corporations from CSV\n", totalProcessed)
	return totalProcessed, nil
}

// parseCSVRecord parses a single CSV record into Corporation (updated for gBizINFO API compliance)
func (c *GBizClient) parseCSVRecord(headers []string, record []string) (*domain.CreateCorporationRequest, error) {
	// Create a map for easier field access
	fieldMap := make(map[string]string)
	for i, header := range headers {
		if i < len(record) {
			// Clean header (remove BOM and quotes)
			cleanHeader := strings.Trim(strings.TrimPrefix(header, "\ufeff"), "\"")
			fieldMap[cleanHeader] = strings.TrimSpace(record[i])
		}
	}

	// Extract corporate number from first column or specific header
	corporateNumber := ""
	if len(record) > 0 {
		corporateNumber = strings.TrimSpace(record[0])
	}
	// Also try from named field
	if val, exists := fieldMap["法人番号"]; exists && val != "" {
		corporateNumber = val
	}

	// Validate corporate number (should be 13 digits)
	if len(corporateNumber) != 13 {
		return nil, fmt.Errorf("invalid corporate number: %s", corporateNumber)
	}

	// Extract name from second column or specific header
	name := ""
	if len(record) > 1 {
		name = strings.TrimSpace(record[1])
	}
	// Also try from named field
	if val, exists := fieldMap["法人名"]; exists && val != "" {
		name = val
	}

	if name == "" {
		return nil, fmt.Errorf("empty corporation name")
	}

	corp := &domain.CreateCorporationRequest{
		CorporateNumber: corporateNumber,
		Name:            name,
		Status:          "01", // Default status code (matches gBizINFO API)
	}

	// Generate search_name using text converter
	searchName := c.textConverter.NormalizeForSearch(name)
	corp.SearchName = &searchName

	// Map CSV fields to gBizINFO API structure

	// Basic Information
	if val, exists := fieldMap["法人名ふりがな"]; exists && val != "" {
		corp.Kana = &val
	}
	// Also try alternative field name used in some CSV files
	if val, exists := fieldMap["法人名カナ"]; exists && val != "" {
		corp.Kana = &val
	}

	if val, exists := fieldMap["法人名英語"]; exists && val != "" {
		corp.NameEn = &val
	}

	if val, exists := fieldMap["郵便番号"]; exists && val != "" {
		corp.PostalCode = &val
	}

	if val, exists := fieldMap["本社所在地"]; exists && val != "" {
		corp.Location = &val
	}
	// Also try alternative field name used in some CSV files
	if val, exists := fieldMap["本店所在地"]; exists && val != "" {
		corp.Location = &val
	}

	// Extract prefecture code from location
	if corp.Location != nil {
		prefCode := extractPrefectureCode(*corp.Location)
		if prefCode != "" {
			corp.PrefectureCode = &prefCode
		}
	}

	if val, exists := fieldMap["ステータス"]; exists && val != "" {
		corp.Status = val
	}

	// Registration Information
	if val, exists := fieldMap["登記記録の閉鎖等年月日"]; exists && val != "" {
		if date, err := parseDate(val); err == nil && date != nil {
			corp.CloseDate = date
		}
	}

	if val, exists := fieldMap["登記記録の閉鎖等の事由"]; exists && val != "" {
		corp.CloseCause = &val
	}

	// Representative Information
	if val, exists := fieldMap["法人代表者名"]; exists && val != "" {
		corp.RepresentativeName = &val
	}
	// Also try alternative field name used in some CSV files
	if val, exists := fieldMap["代表者名"]; exists && val != "" {
		corp.RepresentativeName = &val
	}

	if val, exists := fieldMap["法人代表者役職"]; exists && val != "" {
		corp.RepresentativePosition = &val
	}

	// Company Details
	if val, exists := fieldMap["設立年月日"]; exists && val != "" {
		if date, err := parseDate(val); err == nil && date != nil {
			corp.DateOfEstablishment = date
		}
	}

	if val, exists := fieldMap["創業年"]; exists && val != "" {
		if year, err := strconv.ParseInt(val, 10, 32); err == nil {
			foundingYear := int32(year)
			corp.FoundingYear = &foundingYear
		}
	}

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

	if val, exists := fieldMap["企業規模詳細(男性)"]; exists && val != "" {
		if male, err := strconv.ParseInt(val, 10, 32); err == nil {
			sizeMale := int32(male)
			corp.CompanySizeMale = &sizeMale
		}
	}

	if val, exists := fieldMap["企業規模詳細(女性)"]; exists && val != "" {
		if female, err := strconv.ParseInt(val, 10, 32); err == nil {
			sizeFemale := int32(female)
			corp.CompanySizeFemale = &sizeFemale
		}
	}

	// Business Information
	if val, exists := fieldMap["営業品目リスト"]; exists && val != "" {
		corp.BusinessItems = &val
	}

	if val, exists := fieldMap["事業概要"]; exists && val != "" {
		corp.BusinessSummary = &val
	}

	if val, exists := fieldMap["企業ホームページ"]; exists && val != "" {
		corp.CompanyUrl = &val
	}

	if val, exists := fieldMap["資格等級"]; exists && val != "" {
		corp.QualificationGrade = &val
	}

	// Financial Information
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

	// gBizINFO Metadata
	if val, exists := fieldMap["最終更新日"]; exists && val != "" {
		if date, err := parseDate(val); err == nil && date != nil {
			corp.UpdateDate = date
		}
	}

	return corp, nil
}

// parseFinanceCSVStream parses finance CSV content in streaming batches to minimize memory usage
func (c *GBizClient) parseFinanceCSVStream(reader io.Reader, batchProcessor func([]*domain.Finance) error) (int, error) {
	csvReader := csv.NewReader(reader)

	// Configure CSV reader for gBizINFO format
	csvReader.LazyQuotes = true    // Allow lazy quotes for malformed CSV
	csvReader.FieldsPerRecord = -1 // Allow variable number of fields per record

	// Read header
	headers, err := csvReader.Read()
	if err != nil {
		return 0, fmt.Errorf("failed to read CSV header: %w", err)
	}

	fmt.Printf("Finance CSV Headers: %v\n", headers)
	fmt.Printf("Number of columns: %d\n", len(headers))

	const batchSize = 100 // Further reduced batch size for optimal memory efficiency
	batch := make([]*domain.Finance, 0, batchSize)
	lineNum := 1
	totalProcessed := 0

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			// Process final batch if any records remain
			if len(batch) > 0 {
				if err := batchProcessor(batch); err != nil {
					return totalProcessed, fmt.Errorf("failed to process final batch at line %d: %w", lineNum, err)
				}
				totalProcessed += len(batch)
				fmt.Printf("Processed final batch of %d finance records (total: %d)\n", len(batch), totalProcessed)
			}
			break
		}
		if err != nil {
			fmt.Printf("Warning: failed to read line %d: %v\n", lineNum, err)
			lineNum++
			continue
		}

		// Skip if not enough columns
		if len(record) < 5 {
			lineNum++
			continue
		}

		finance, err := c.parseFinanceCSVRecord(headers, record)
		if err != nil {
			fmt.Printf("Warning: failed to parse finance line %d: %v\n", lineNum, err)
			lineNum++
			continue
		}

		if finance != nil {
			batch = append(batch, finance)

			// Process batch when it reaches the target size
			if len(batch) >= batchSize {
				if err := batchProcessor(batch); err != nil {
					return totalProcessed, fmt.Errorf("failed to process batch at line %d: %w", lineNum, err)
				}
				totalProcessed += len(batch)
				fmt.Printf("Processed batch of %d finance records (total: %d)\n", len(batch), totalProcessed)

				// Clear batch for next iteration (reuse slice to minimize allocations)
				batch = batch[:0]
			}
		}
		lineNum++

		// Progress indicator for large files - reduced frequency to minimize I/O overhead
		if lineNum%100000 == 0 {
			fmt.Printf("Processed %d lines...\n", lineNum)
		}
	}

	fmt.Printf("Successfully processed %d finance records from CSV\n", totalProcessed)
	return totalProcessed, nil
}

// parseFinanceCSVRecord parses a single CSV record into Finance
func (c *GBizClient) parseFinanceCSVRecord(headers []string, record []string) (*domain.Finance, error) {
	// Create a map for easier field access
	fieldMap := make(map[string]string)
	for i, header := range headers {
		if i < len(record) {
			// Clean header (remove BOM and quotes)
			cleanHeader := strings.Trim(strings.TrimPrefix(header, "\ufeff"), "\"")
			fieldMap[cleanHeader] = strings.TrimSpace(record[i])
		}
	}

	// Extract corporate number from first column or specific header
	corporateNumber := ""
	if len(record) > 0 {
		corporateNumber = strings.TrimSpace(record[0])
	}
	// Also try from named field
	if val, exists := fieldMap["法人番号"]; exists && val != "" {
		corporateNumber = val
	}

	// Validate corporate number (should be 13 digits)
	if len(corporateNumber) != 13 {
		return nil, fmt.Errorf("invalid corporate number: %s", corporateNumber)
	}

	// Use domain constructor to create Finance with UUIDv7
	finance := domain.NewFinance(corporateNumber)

	// Map CSV fields to Finance structure based on gBizINFO finance CSV specification
	if val, exists := fieldMap["法人名（法人番号）"]; exists && val != "" {
		finance.CorporateNameFromNumber = val
	}
	if val, exists := fieldMap["本社所在地（法人番号）"]; exists && val != "" {
		finance.HeadOfficeLocationFromNumber = val
	}
	if val, exists := fieldMap["法人名"]; exists && val != "" {
		finance.CorporateName = val
	}
	if val, exists := fieldMap["本社所在地"]; exists && val != "" {
		finance.HeadOfficeLocation = val
	}
	if val, exists := fieldMap["会計基準"]; exists && val != "" {
		finance.AccountingStandards = val
	}
	if val, exists := fieldMap["事業年度"]; exists && val != "" {
		finance.BusinessYear = val
	}
	if val, exists := fieldMap["回次"]; exists && val != "" {
		finance.PeriodNumber = val
	}
	if val, exists := fieldMap["売上高"]; exists && val != "" {
		finance.SalesRevenue = val
	}
	if val, exists := fieldMap["売上高（単位)"]; exists && val != "" {
		finance.SalesRevenueUnit = val
	}
	if val, exists := fieldMap["営業収益"]; exists && val != "" {
		finance.OperatingRevenue1 = val
	}
	if val, exists := fieldMap["営業収益（単位）"]; exists && val != "" {
		finance.OperatingRevenue1Unit = val
	}
	if val, exists := fieldMap["営業収入"]; exists && val != "" {
		finance.OperatingRevenue2 = val
	}
	if val, exists := fieldMap["営業収入（単位）"]; exists && val != "" {
		finance.OperatingRevenue2Unit = val
	}
	if val, exists := fieldMap["営業総収入"]; exists && val != "" {
		finance.GrossOperatingRevenue = val
	}
	if val, exists := fieldMap["営業総収入（単位）"]; exists && val != "" {
		finance.GrossOperatingRevenueUnit = val
	}
	if val, exists := fieldMap["経常収益"]; exists && val != "" {
		finance.OrdinaryRevenue = val
	}
	if val, exists := fieldMap["経常収益（単位）"]; exists && val != "" {
		finance.OrdinaryRevenueUnit = val
	}
	if val, exists := fieldMap["正味収入保険料"]; exists && val != "" {
		finance.NetPremiumsWritten = val
	}
	if val, exists := fieldMap["正味収入保険料（単位）"]; exists && val != "" {
		finance.NetPremiumsWrittenUnit = val
	}
	if val, exists := fieldMap["経常利益又は経常損失（△）"]; exists && val != "" {
		finance.OrdinaryIncome = val
	}
	if val, exists := fieldMap["経常利益又は経常損失（△）(単位)"]; exists && val != "" {
		finance.OrdinaryIncomeUnit = val
	}
	if val, exists := fieldMap["当期純利益又は当期純損失（△）"]; exists && val != "" {
		finance.NetIncome = val
	}
	if val, exists := fieldMap["当期純利益又は当期純損失（△）(単位)"]; exists && val != "" {
		finance.NetIncomeUnit = val
	}
	if val, exists := fieldMap["資本金"]; exists && val != "" {
		finance.CapitalStock = val
	}
	if val, exists := fieldMap["資本金(単位)"]; exists && val != "" {
		finance.CapitalStockUnit = val
	}
	if val, exists := fieldMap["純資産額"]; exists && val != "" {
		finance.NetAssets = val
	}
	if val, exists := fieldMap["純資産額(単位)"]; exists && val != "" {
		finance.NetAssetsUnit = val
	}
	if val, exists := fieldMap["総資産額"]; exists && val != "" {
		finance.TotalAssets = val
	}
	if val, exists := fieldMap["総資産額(単位)"]; exists && val != "" {
		finance.TotalAssetsUnit = val
	}
	if val, exists := fieldMap["従業員数"]; exists && val != "" {
		finance.NumberOfEmployees = val
	}
	if val, exists := fieldMap["従業員数(単位)"]; exists && val != "" {
		finance.NumberOfEmployeesUnit = val
	}
	if val, exists := fieldMap["大株主1"]; exists && val != "" {
		finance.MajorShareholder1 = val
	}
	if val, exists := fieldMap["発行済株式総数に対する所有株式数の割合1"]; exists && val != "" {
		finance.ShareholdingRatio1 = val
	}
	if val, exists := fieldMap["大株主2"]; exists && val != "" {
		finance.MajorShareholder2 = val
	}
	if val, exists := fieldMap["発行済株式総数に対する所有株式数の割合2"]; exists && val != "" {
		finance.ShareholdingRatio2 = val
	}
	if val, exists := fieldMap["大株主3"]; exists && val != "" {
		finance.MajorShareholder3 = val
	}
	if val, exists := fieldMap["発行済株式総数に対する所有株式数の割合3"]; exists && val != "" {
		finance.ShareholdingRatio3 = val
	}
	if val, exists := fieldMap["大株主4"]; exists && val != "" {
		finance.MajorShareholder4 = val
	}
	if val, exists := fieldMap["発行済株式総数に対する所有株式数の割合4"]; exists && val != "" {
		finance.ShareholdingRatio4 = val
	}
	if val, exists := fieldMap["大株主5"]; exists && val != "" {
		finance.MajorShareholder5 = val
	}
	if val, exists := fieldMap["発行済株式総数に対する所有株式数の割合5"]; exists && val != "" {
		finance.ShareholdingRatio5 = val
	}
	return finance, nil
}

// Cleanup removes temporary files
func (c *GBizClient) Cleanup(filePath string) error {
	if filePath != "" {
		return os.Remove(filePath)
	}
	return nil
}

// DisplayCSVHeaders extracts and displays CSV headers from a ZIP file
func (c *GBizClient) DisplayCSVHeaders(zipPath string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open ZIP file: %w", err)
	}
	defer reader.Close()

	// Find and process the first CSV file in the ZIP
	for _, file := range reader.File {
		if strings.HasSuffix(strings.ToLower(file.Name), ".csv") {
			log.Printf("Processing CSV file: %s", file.Name)

			fileReader, err := file.Open()
			if err != nil {
				return fmt.Errorf("failed to open CSV file %s: %w", file.Name, err)
			}
			defer fileReader.Close()

			csvReader := csv.NewReader(fileReader)
			csvReader.LazyQuotes = true
			csvReader.FieldsPerRecord = -1

			// Read only the header
			headers, err := csvReader.Read()
			if err != nil {
				return fmt.Errorf("failed to read CSV header: %w", err)
			}

			log.Printf("CSV Headers (%d columns):", len(headers))
			for i, header := range headers {
				log.Printf("  [%d] %s", i+1, header)
			}

			// Only process the first CSV file
			break
		}
	}

	return nil
}

// extractPrefectureCode extracts prefecture code from location string
func extractPrefectureCode(location string) string {
	if location == "" {
		return ""
	}

	// Match prefecture names (北海道、〜県、〜都、〜府) from the beginning of the location string
	for prefName, code := range prefectureCodeMap {
		if strings.HasPrefix(location, prefName) {
			return code
		}
	}

	// If no prefecture match found, try matching major cities
	for cityName, code := range cityCodeMap {
		if strings.HasPrefix(location, cityName) {
			return code
		}
	}

	return "" // No prefecture or city match found
}

// parseDate parses various date formats commonly found in CSV files
func parseDate(dateStr string) (*time.Time, error) {
	if dateStr == "" {
		return nil, nil
	}

	// List of date formats to try in order
	formats := []string{
		"2006-01-02T15:04:05Z07:00", // RFC3339 with timezone
		"2006-01-02T15:04:05Z",      // RFC3339 UTC
		"2006-01-02",                // YYYY-MM-DD
		"2006/01/02",                // YYYY/MM/DD
		"2006-01-02T15:04:05",       // ISO format without timezone
	}

	for _, format := range formats {
		if parsed, err := time.Parse(format, dateStr); err == nil {
			return &parsed, nil
		}
	}

	return nil, fmt.Errorf("unable to parse date: %s", dateStr)
}

// ImportFromCSVFile imports corporations from a CSV file directly
func (c *GBizClient) ImportFromCSVFile(ctx context.Context, filePath string) ([]*domain.CreateCorporationRequest, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %s", filePath)
	}

	// Open CSV file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	log.Printf("Starting CSV import from file: %s", filePath)

	// Parse CSV
	corporations, err := c.parseCSV(file)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV: %w", err)
	}

	log.Printf("Successfully imported %d corporations from CSV file", len(corporations))
	return corporations, nil
}
