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
    postal_code VARCHAR(10),                             -- postal_code
    location TEXT,                                       -- location (所在地)
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
CREATE INDEX idx_corporations_status ON corporations(status);
CREATE INDEX idx_corporations_location ON corporations(location);
CREATE INDEX idx_corporations_representative_name ON corporations(representative_name);
CREATE INDEX idx_corporations_date_of_establishment ON corporations(date_of_establishment);
CREATE INDEX idx_corporations_update_date ON corporations(update_date);
CREATE INDEX idx_corporations_created_at ON corporations(created_at);
