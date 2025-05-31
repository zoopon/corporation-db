package infrastructure

import (
	"strings"
	"testing"
)

func TestExtractPrefectureCode(t *testing.T) {
	tests := []struct {
		name     string
		location string
		expected string
	}{
		// 都道府県名のテスト
		{
			name:     "北海道",
			location: "北海道札幌市中央区大通1-1",
			expected: "01",
		},
		{
			name:     "青森県",
			location: "青森県青森市長島1-1-1",
			expected: "02",
		},
		{
			name:     "東京都",
			location: "東京都新宿区西新宿2-8-1",
			expected: "13",
		},
		{
			name:     "大阪府",
			location: "大阪府大阪市中央区北浜1-1",
			expected: "27",
		},
		{
			name:     "福岡県",
			location: "福岡県福岡市博多区博多駅前1-1",
			expected: "40",
		},
		{
			name:     "沖縄県",
			location: "沖縄県那覇市久茂地1-1",
			expected: "47",
		},
		// 主要都市名のテスト
		{
			name:     "札幌市（北海道）",
			location: "札幌市中央区大通1-1",
			expected: "01",
		},
		{
			name:     "仙台市（宮城県）",
			location: "仙台市青葉区本町3-8-1",
			expected: "04",
		},
		{
			name:     "横浜市（神奈川県）",
			location: "横浜市中区本町6-50-10",
			expected: "14",
		},
		{
			name:     "名古屋市（愛知県）",
			location: "名古屋市中区栄1-1",
			expected: "23",
		},
		{
			name:     "京都市（京都府）",
			location: "京都市下京区烏丸通七条下ル東塩小路町735-1",
			expected: "26",
		},
		{
			name:     "神戸市（兵庫県）",
			location: "神戸市中央区加納町6-5-1",
			expected: "28",
		},
		{
			name:     "広島市（広島県）",
			location: "広島市中区基町10-52",
			expected: "34",
		},
		// 都道府県名と「県」「府」「都」「道」なしのテスト
		{
			name:     "愛知（愛知県）",
			location: "愛知名古屋市中区栄1-1",
			expected: "23",
		},
		{
			name:     "兵庫（兵庫県）",
			location: "兵庫神戸市中央区加納町6-5-1",
			expected: "28",
		},
		// マッチしないケース
		{
			name:     "不明な地名",
			location: "不明県不明市不明町1-1",
			expected: "",
		},
		{
			name:     "空文字列",
			location: "",
			expected: "",
		},
		{
			name:     "アルファベットのみ",
			location: "Tokyo Shibuya 1-1-1",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractPrefectureCode(tt.location)
			if result != tt.expected {
				t.Errorf("extractPrefectureCode(%q) = %q; want %q", tt.location, result, tt.expected)
			}
		})
	}
}

func TestPrefectureMapping(t *testing.T) {
	// 全47都道府県のテスト
	expectedMappings := map[string]string{
		"北海道": "01", "青森": "02", "岩手": "03", "宮城": "04", "秋田": "05",
		"山形": "06", "福島": "07", "茨城": "08", "栃木": "09", "群馬": "10",
		"埼玉": "11", "千葉": "12", "東京": "13", "神奈川": "14", "新潟": "15",
		"富山": "16", "石川": "17", "福井": "18", "山梨": "19", "長野": "20",
		"岐阜": "21", "静岡": "22", "愛知": "23", "三重": "24", "滋賀": "25",
		"京都": "26", "大阪": "27", "兵庫": "28", "奈良": "29", "和歌山": "30",
		"鳥取": "31", "島根": "32", "岡山": "33", "広島": "34", "山口": "35",
		"徳島": "36", "香川": "37", "愛媛": "38", "高知": "39", "福岡": "40",
		"佐賀": "41", "長崎": "42", "熊本": "43", "大分": "44", "宮崎": "45",
		"鹿児島": "46", "沖縄": "47",
	}

	for prefName, expectedCode := range expectedMappings {
		t.Run(prefName, func(t *testing.T) {
			// 県名のみ
			result := extractPrefectureCode(prefName + "市test町1-1")
			if result != expectedCode {
				t.Errorf("extractPrefectureCode(%q) = %q; want %q", prefName+"市test町1-1", result, expectedCode)
			}

			// 県名＋県/府/都/道
			var suffix string
			switch prefName {
			case "北海道":
				suffix = ""
			case "東京":
				suffix = "都"
			case "大阪", "京都":
				suffix = "府"
			default:
				suffix = "県"
			}

			fullName := prefName + suffix
			result = extractPrefectureCode(fullName + "test市test町1-1")
			if result != expectedCode {
				t.Errorf("extractPrefectureCode(%q) = %q; want %q", fullName+"test市test町1-1", result, expectedCode)
			}
		})
	}
}

func TestCityMapping(t *testing.T) {
	// 主要都市名のテスト
	cityMappings := map[string]string{
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
		"横浜市":   "14", // 神奈川県
		"新潟市":   "15", // 新潟県
		"富山市":   "16", // 富山県
		"金沢市":   "17", // 石川県
		"福井市":   "18", // 福井県
		"甲府市":   "19", // 山梨県
		"長野市":   "20", // 長野県
		"岐阜市":   "21", // 岐阜県
		"静岡市":   "22", // 静岡県
		"名古屋市":  "23", // 愛知県
		"津市":    "24", // 三重県
		"大津市":   "25", // 滋賀県
		"京都市":   "26", // 京都府
		"大阪市":   "27", // 大阪府
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
		"佐賀市":   "41", // 佐賀県
		"長崎市":   "42", // 長崎県
		"熊本市":   "43", // 熊本県
		"大分市":   "44", // 大分県
		"宮崎市":   "45", // 宮崎県
		"鹿児島市":  "46", // 鹿児島県
		"那覇市":   "47", // 沖縄県
	}

	for cityName, expectedCode := range cityMappings {
		t.Run(cityName, func(t *testing.T) {
			result := extractPrefectureCode(cityName + "test町1-1")
			if result != expectedCode {
				t.Errorf("extractPrefectureCode(%q) = %q; want %q", cityName+"test町1-1", result, expectedCode)
			}
		})
	}
}

func TestParseCSVWithPrefectureCodeExtraction(t *testing.T) {
	// テスト用のCSVデータ
	csvData := `法人番号,法人名,法人名ふりがな,法人名英語,郵便番号,本社所在地,ステータス,登記記録の閉鎖等年月日,登記記録の閉鎖等の事由,法人代表者名,法人代表者役職,資本金,従業員数,企業規模詳細(男性),企業規模詳細(女性),営業品目リスト,事業概要,企業ホームページ,設立年月日,創業年,最終更新日,資格等級
4567890123456,北海道株式会社,ホッカイドウカブシキガイシャ,Hokkaido Corporation,064-0001,北海道札幌市中央区大通1-1,01,,,鈴木一郎,代表取締役,8000000,75,40,35,情報通信業,IT関連事業,https://hokkaido.co.jp,2020-01-01,2020,2025-05-31,
5678901234567,沖縄株式会社,オキナワカブシキガイシャ,Okinawa Corporation,900-0015,沖縄県那覇市久茂地1-1,01,,,田中美香,代表取締役,6000000,30,15,15,観光業,観光サービス業,https://okinawa.co.jp,2018-05-01,2018,2025-05-31,
6789012345678,福岡有限会社,フクオカユウゲンガイシャ,Fukuoka Company,812-0011,福岡県福岡市博多区博多駅前1-1,01,,,高橋健太,代表取締役,4000000,40,25,15,製造業,製品製造業,https://fukuoka.co.jp,2019-03-15,2019,2025-05-31,`

	client := &GBizClient{}
	reader := strings.NewReader(csvData)

	corporations, err := client.parseCSV(reader)
	if err != nil {
		t.Fatalf("parseCSV() error = %v", err)
	}

	if len(corporations) != 3 {
		t.Fatalf("Expected 3 corporations, got %d", len(corporations))
	}

	// Test cases with expected prefecture codes
	testCases := []struct {
		corporateNumber  string
		expectedName     string
		expectedLocation string
		expectedPrefCode string
	}{
		{
			corporateNumber:  "4567890123456",
			expectedName:     "北海道株式会社",
			expectedLocation: "北海道札幌市中央区大通1-1",
			expectedPrefCode: "01",
		},
		{
			corporateNumber:  "5678901234567",
			expectedName:     "沖縄株式会社",
			expectedLocation: "沖縄県那覇市久茂地1-1",
			expectedPrefCode: "47",
		},
		{
			corporateNumber:  "6789012345678",
			expectedName:     "福岡有限会社",
			expectedLocation: "福岡県福岡市博多区博多駅前1-1",
			expectedPrefCode: "40",
		},
	}

	for i, tc := range testCases {
		corp := corporations[i]

		if corp.CorporateNumber != tc.corporateNumber {
			t.Errorf("Corporation %d: CorporateNumber = %q; want %q", i, corp.CorporateNumber, tc.corporateNumber)
		}

		if corp.Name != tc.expectedName {
			t.Errorf("Corporation %d: Name = %q; want %q", i, corp.Name, tc.expectedName)
		}

		if corp.Location == nil || *corp.Location != tc.expectedLocation {
			var location string
			if corp.Location != nil {
				location = *corp.Location
			}
			t.Errorf("Corporation %d: Location = %q; want %q", i, location, tc.expectedLocation)
		}

		if corp.PrefectureCode == nil || *corp.PrefectureCode != tc.expectedPrefCode {
			var prefCode string
			if corp.PrefectureCode != nil {
				prefCode = *corp.PrefectureCode
			}
			t.Errorf("Corporation %d: PrefectureCode = %q; want %q", i, prefCode, tc.expectedPrefCode)
		}
	}
}
