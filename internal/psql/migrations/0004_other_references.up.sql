-- Типы объектов
CREATE TABLE item_types (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

-- Бренды объектов
CREATE TABLE item_brands (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

-- Модели объектов
CREATE TABLE item_models (
    id SERIAL PRIMARY KEY,
    brand_id INTEGER REFERENCES item_brands(id) ON DELETE CASCADE NOT NULL,
    name TEXT NOT NULL,
    UNIQUE (brand_id, name)
);

-- Цвета объектов
CREATE TABLE item_colors (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

-- Виды повреждений
CREATE TABLE damages (
    id SERIAL PRIMARY KEY,
    description TEXT NOT NULL UNIQUE
);