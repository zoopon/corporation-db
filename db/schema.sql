-- db/schema.sql
-- Full database schema for sqldef based on gBizINFO official API specification

-- 法人情報テーブル (Based on gBizINFO REST API specification)
CREATE TABLE corporations (
    id SERIAL PRIMARY KEY,
    
    -- Basic Information (基本情報)
    corporate_number CHAR(13) UNIQUE NOT NULL,           -- corporate_number
    name VARCHAR(500) NOT NULL,                          -- name
    kana VARCHAR(500),                                   -- kana (法人名ふりがな)
    name_en VARCHAR(500),                                -- name_en (英語法人名)
    search_name VARCHAR(500),                            -- search_name (検索用正規化された名前)
    postal_code VARCHAR(10),                             -- postal_code
    location TEXT,                                       -- location (所在地)
    prefecture_code VARCHAR(2),                          -- prefecture_code (都道府県コード)
    status VARCHAR(100) NOT NULL DEFAULT '01',           -- status (法人状態)
    
    -- Registration Information (登記情報)
    close_date DATE,                                     -- close_date (登記記録の閉鎖等年月日)
    close_cause TEXT,                                    -- close_cause (登記記録の閉鎖等の事由)
    
    -- Representative Information (代表者情報)
    representative_name VARCHAR(255),                    -- representative_name
    representative_position VARCHAR(255),                -- representative_position
    
    -- Company Details (企業詳細)
    date_of_establishment DATE,                          -- date_of_establishment (設立年月日)
    founding_year INTEGER,                               -- founding_year (創業年)
    capital_stock BIGINT,                                -- capital_stock (資本金)
    employee_number INTEGER,                             -- employee_number (従業員数)
    company_size_male INTEGER,                           -- company_size_male
    company_size_female INTEGER,                         -- company_size_female
    
    -- Business Information (事業情報)
    business_items TEXT,                                 -- business_items (JSON array as text)
    business_summary TEXT,                               -- business_summary (事業概要)
    company_url VARCHAR(500),                            -- company_url (企業ホームページ)
    qualification_grade VARCHAR(100),                    -- qualification_grade
    number_of_activity VARCHAR(100),                     -- number_of_activity
    
    -- gBizINFO Metadata
    update_date DATE,                                    -- update_date (gBizINFOでの最終更新日)
    
    -- Database Metadata
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 法人情報のインデックス (Based on gBizINFO API structure)
CREATE INDEX idx_corporations_corporate_number ON corporations(corporate_number);
CREATE INDEX idx_corporations_name ON corporations(name);
CREATE INDEX idx_corporations_kana ON corporations(kana);
CREATE INDEX idx_corporations_search_name ON corporations(search_name);
CREATE INDEX idx_corporations_status ON corporations(status);
CREATE INDEX idx_corporations_location ON corporations(location);
CREATE INDEX idx_corporations_prefecture_code ON corporations(prefecture_code);
CREATE INDEX idx_corporations_representative_name ON corporations(representative_name);
CREATE INDEX idx_corporations_date_of_establishment ON corporations(date_of_establishment);
CREATE INDEX idx_corporations_update_date ON corporations(update_date);
CREATE INDEX idx_corporations_created_at ON corporations(created_at);

-- 財務情報テーブル (Based on gBizINFO Finance CSV and API specification)
CREATE TABLE finances (
    id SERIAL PRIMARY KEY,
    
    -- 基本情報 (Basic Information)
    corporate_number CHAR(13) NOT NULL,                  -- 法人番号
    corporate_name_from_number VARCHAR(500),             -- 法人名（法人番号）
    head_office_location_from_number TEXT,               -- 本社所在地（法人番号）
    corporate_name VARCHAR(500),                         -- 法人名
    head_office_location TEXT,                           -- 本社所在地
    
    -- 会計情報 (Accounting Information)
    accounting_standards VARCHAR(100),                   -- 会計基準
    business_year VARCHAR(100),                          -- 事業年度
    period_number VARCHAR(50),                           -- 回次
    
    -- 売上・収益 (Sales and Revenue)
    sales_revenue VARCHAR(50),                           -- 売上高
    sales_revenue_unit VARCHAR(20),                      -- 売上高（単位）
    operating_revenue1 VARCHAR(50),                      -- 営業収益
    operating_revenue1_unit VARCHAR(20),                 -- 営業収益（単位）
    operating_revenue2 VARCHAR(50),                      -- 営業収入
    operating_revenue2_unit VARCHAR(20),                 -- 営業収入（単位）
    gross_operating_revenue VARCHAR(50),                 -- 営業総収入
    gross_operating_revenue_unit VARCHAR(20),            -- 営業総収入（単位）
    ordinary_revenue VARCHAR(50),                        -- 経常収益
    ordinary_revenue_unit VARCHAR(20),                   -- 経常収益（単位）
    net_premiums_written VARCHAR(50),                    -- 正味収入保険料
    net_premiums_written_unit VARCHAR(20),               -- 正味収入保険料（単位）
    
    -- 利益 (Profit)
    ordinary_income VARCHAR(50),                         -- 経常利益又は経常損失（△）
    ordinary_income_unit VARCHAR(20),                    -- 経常利益又は経常損失（△）（単位）
    net_income VARCHAR(50),                              -- 当期純利益又は当期純損失（△）
    net_income_unit VARCHAR(20),                         -- 当期純利益又は当期純損失（△）（単位）
    
    -- 資本・資産 (Capital and Assets)
    capital_stock VARCHAR(50),                           -- 資本金
    capital_stock_unit VARCHAR(20),                      -- 資本金（単位）
    net_assets VARCHAR(50),                              -- 純資産額
    net_assets_unit VARCHAR(20),                         -- 純資産額（単位）
    total_assets VARCHAR(50),                            -- 総資産額
    total_assets_unit VARCHAR(20),                       -- 総資産額（単位）
    
    -- 従業員 (Employees)
    number_of_employees VARCHAR(50),                     -- 従業員数
    number_of_employees_unit VARCHAR(20),                -- 従業員数（単位）
    
    -- 大株主情報 (Major Shareholders)
    major_shareholder1 VARCHAR(500),                     -- 大株主1
    shareholding_ratio1 VARCHAR(20),                     -- 発行済株式総数に対する所有株式数の割合1
    major_shareholder2 VARCHAR(500),                     -- 大株主2
    shareholding_ratio2 VARCHAR(20),                     -- 発行済株式総数に対する所有株式数の割合2
    major_shareholder3 VARCHAR(500),                     -- 大株主3
    shareholding_ratio3 VARCHAR(20),                     -- 発行済株式総数に対する所有株式数の割合3
    major_shareholder4 VARCHAR(500),                     -- 大株主4
    shareholding_ratio4 VARCHAR(20),                     -- 発行済株式総数に対する所有株式数の割合4
    major_shareholder5 VARCHAR(500),                     -- 大株主5
    shareholding_ratio5 VARCHAR(20),                     -- 発行済株式総数に対する所有株式数の割合5
    
    -- Database Metadata
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Foreign Key
    FOREIGN KEY (corporate_number) REFERENCES corporations(corporate_number) ON DELETE CASCADE
);

-- 財務情報のインデックス
CREATE INDEX idx_finances_corporate_number ON finances(corporate_number);
CREATE INDEX idx_finances_business_year ON finances(business_year);
CREATE INDEX idx_finances_period_number ON finances(period_number);
CREATE INDEX idx_finances_accounting_standards ON finances(accounting_standards);
CREATE INDEX idx_finances_created_at ON finances(created_at);
