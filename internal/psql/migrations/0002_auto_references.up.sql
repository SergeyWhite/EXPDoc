-- Марки автомобилей
CREATE TABLE auto_brands (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

-- Модели автомобилей
CREATE TABLE auto_models (
    id SERIAL PRIMARY KEY,
    brand_id INTEGER REFERENCES auto_brands(id) ON DELETE CASCADE NOT NULL,
    name TEXT NOT NULL,
    UNIQUE (brand_id, name)
);

-- Цвета автомобилей
CREATE TABLE auto_colors (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);