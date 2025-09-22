CREATE TABLE racas (
    id SERIAL PRIMARY KEY,
    nome VARCHAR(100) NOT NULL,
    atributo_bonus_1 VARCHAR(20),
    valor_bonus_1 INTEGER DEFAULT 0,
    atributo_bonus_2 VARCHAR(20),
    valor_bonus_2 INTEGER DEFAULT 0,
    atributo_bonus_3 VARCHAR(20),
    valor_bonus_3 INTEGER DEFAULT 0,
    atributo_penalidade VARCHAR(20),
    valor_penalidade INTEGER DEFAULT 0,
    tamanho VARCHAR(20),
    deslocamento INTEGER DEFAULT 9,
    descricao TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);
