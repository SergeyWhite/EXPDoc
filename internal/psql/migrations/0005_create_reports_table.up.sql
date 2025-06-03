CREATE TABLE reports (
    id SERIAL PRIMARY KEY,
    type report_type NOT NULL,
    date DATE NOT NULL,
    city_id INTEGER REFERENCES cities(id) ON DELETE SET NULL,
    doc_type TEXT NOT NULL,
    service TEXT NOT NULL,
    client_id INTEGER REFERENCES counteragents(id) ON DELETE RESTRICT NOT NULL,
    client_basis TEXT,
    affected_id INTEGER REFERENCES counteragents(id) ON DELETE SET NULL,
    event_date DATE,
    deadline_days INTEGER,
    expert_id INTEGER REFERENCES users(id) ON DELETE RESTRICT NOT NULL,
    
    -- JSONB данные
    auto_data JSONB,
    real_estate_data JSONB,
    other_data JSONB,
    
    -- Ограничения для JSONB полей
    CONSTRAINT auto_data_check 
        CHECK (type <> 'auto' OR auto_data IS NOT NULL),
    
    CONSTRAINT real_estate_data_check 
        CHECK (type <> 'real_estate' OR real_estate_data IS NOT NULL),
    
    CONSTRAINT other_data_check 
        CHECK (type <> 'other' OR other_data IS NOT NULL)
);

-- Индексы для JSONB полей
CREATE INDEX idx_reports_auto_data ON reports USING gin (auto_data);
CREATE INDEX idx_reports_real_estate_data ON reports USING gin (real_estate_data);
CREATE INDEX idx_reports_other_data ON reports USING gin (other_data);