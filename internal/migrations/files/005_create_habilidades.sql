CREATE TABLE habilidade_racas (
    id SERIAL PRIMARY KEY,
    raca_id INTEGER REFERENCES racas(id),
    nome VARCHAR(100) NOT NULL,
    descricao TEXT,
    opcional BOOLEAN DEFAULT FALSE,
    nivel_minimo INTEGER DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE TABLE habilidade_classes (
    id SERIAL PRIMARY KEY,
    classe_id INTEGER REFERENCES classes(id),
    nome VARCHAR(100) NOT NULL,
    descricao TEXT,
    nivel INTEGER DEFAULT 1,
    opcional BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);
