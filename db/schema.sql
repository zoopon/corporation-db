-- db/schema.sql
-- Full database schema for sqldef

-- 法人情報テーブル
CREATE TABLE corporations (
    id SERIAL PRIMARY KEY,
    corporate_number CHAR(13) UNIQUE NOT NULL,
    name VARCHAR(500) NOT NULL,
    name_kana VARCHAR(500),
    english_name VARCHAR(500),
    postal_code VARCHAR(10),
    address TEXT,
    prefecture_code VARCHAR(5),
    city_code VARCHAR(10),
    founding_date DATE,
    status VARCHAR(50) NOT NULL DEFAULT '活動中',
    corporate_type VARCHAR(100),
    capital_stock BIGINT,
    employee_number INTEGER,
    representative VARCHAR(255),
    business_description TEXT,
    industry VARCHAR(255),
    website VARCHAR(500),
    phone VARCHAR(20),
    email VARCHAR(255),
    last_updated TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 法人情報のインデックス
CREATE INDEX idx_corporations_corporate_number ON corporations(corporate_number);
CREATE INDEX idx_corporations_name ON corporations(name);
CREATE INDEX idx_corporations_status ON corporations(status);
CREATE INDEX idx_corporations_prefecture_code ON corporations(prefecture_code);
CREATE INDEX idx_corporations_industry ON corporations(industry);
CREATE INDEX idx_corporations_created_at ON corporations(created_at);
CREATE INDEX idx_corporations_last_updated ON corporations(last_updated);
