-- 013_create_pericias.sql
CREATE TABLE IF NOT EXISTS pericias (
    id SERIAL PRIMARY KEY,
    nome VARCHAR(255) NOT NULL UNIQUE,
    atributo VARCHAR(3) NOT NULL, -- FOR, DES, CON, INT, SAB, CAR
    descricao TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Índices para melhor performance
CREATE INDEX IF NOT EXISTS idx_pericias_atributo ON pericias(atributo);
CREATE INDEX IF NOT EXISTS idx_pericias_nome ON pericias(nome);
