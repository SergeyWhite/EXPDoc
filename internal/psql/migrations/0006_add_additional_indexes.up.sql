-- Для поиска по типу отчета
CREATE INDEX idx_reports_type ON reports (type);

-- Для поиска по датам
CREATE INDEX idx_reports_date ON reports (date);
CREATE INDEX idx_reports_event_date ON reports (event_date);

-- Для поиска по срокам
CREATE INDEX idx_reports_deadline ON reports (deadline_days);

-- Для JSONB-полей (пример для VIN и кадастрового номера)
CREATE INDEX idx_auto_vin ON reports ((auto_data->>'vin')) WHERE type = 'auto';
CREATE INDEX idx_real_estate_kadastr ON reports ((real_estate_data->>'cadastral')) WHERE type = 'real_estate';