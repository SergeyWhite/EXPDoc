-- Тип отчета
CREATE TYPE report_type AS ENUM ('auto', 'real_estate', 'other');

-- Таблица городов
CREATE TABLE cities (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

-- Таблица пользователей
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    full_name TEXT NOT NULL,
    position TEXT NOT NULL,
    phone TEXT NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT false,
    login TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);

-- Таблица контрагентов
CREATE TABLE counteragents (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    phone TEXT,
    email TEXT,
    inn TEXT,
    ogrn TEXT,
    bik TEXT,
    passport_series TEXT,
    passport_number TEXT
);