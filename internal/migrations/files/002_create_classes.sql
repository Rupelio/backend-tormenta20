CREATE TABLE classes (
    id SERIAL PRIMARY KEY,
    nome VARCHAR(100) NOT NULL,
    pv_por_nivel INTEGER NOT NULL,
    pm_por_nivel INTEGER NOT NULL,
    atributo_principal VARCHAR(20),
    pericias_quantidade INTEGER DEFAULT 2,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);
