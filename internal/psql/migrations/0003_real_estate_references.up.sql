-- Типы зданий
CREATE TABLE build_types (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

-- Улицы (с привязкой к городам)
CREATE TABLE streets (
    id SERIAL PRIMARY KEY,
    city_id INTEGER REFERENCES cities(id) ON DELETE CASCADE NOT NULL,
    name TEXT NOT NULL,
    UNIQUE (city_id, name)
);