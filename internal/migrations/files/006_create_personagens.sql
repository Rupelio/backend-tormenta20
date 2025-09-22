CREATE TABLE personagens (
    id SERIAL PRIMARY KEY,
    nome VARCHAR(100) NOT NULL,
    nivel INTEGER DEFAULT 1,
    forca INTEGER DEFAULT 0,
    destreza INTEGER DEFAULT 0,
    constituicao INTEGER DEFAULT 0,
    inteligencia INTEGER DEFAULT 0,
    sabedoria INTEGER DEFAULT 0,
    carisma INTEGER DEFAULT 0,
    raca_id INTEGER REFERENCES racas(id),
    classe_id INTEGER REFERENCES classes(id),
    origem_id INTEGER REFERENCES origens(id),
    divindade_id INTEGER REFERENCES divindades(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
