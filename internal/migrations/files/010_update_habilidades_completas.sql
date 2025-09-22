-- Criar tabela para habilidades de origem
CREATE TABLE habilidade_origens (
    id SERIAL PRIMARY KEY,
    origem_id INTEGER REFERENCES origens(id),
    nome VARCHAR(100) NOT NULL,
    descricao TEXT,
    opcional BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Criar tabela para habilidades de divindade
CREATE TABLE habilidade_divindades (
    id SERIAL PRIMARY KEY,
    divindade_id INTEGER REFERENCES divindades(id),
    nome VARCHAR(100) NOT NULL,
    descricao TEXT,
    nivel INTEGER DEFAULT 1,
    opcional BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);
