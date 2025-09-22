-- 015_create_pericias_relations.sql
-- Tabelas de relacionamento para perícias

-- Tabela para perícias que cada raça oferece automaticamente
CREATE TABLE IF NOT EXISTS raca_pericias (
    raca_id INTEGER REFERENCES racas(id) ON DELETE CASCADE,
    pericia_id INTEGER REFERENCES pericias(id) ON DELETE CASCADE,
    PRIMARY KEY (raca_id, pericia_id)
);

-- Tabela para perícias disponíveis para cada classe
CREATE TABLE IF NOT EXISTS classe_pericias (
    classe_id INTEGER REFERENCES classes(id) ON DELETE CASCADE,
    pericia_id INTEGER REFERENCES pericias(id) ON DELETE CASCADE,
    PRIMARY KEY (classe_id, pericia_id)
);

-- Tabela para perícias que cada origem oferece
CREATE TABLE IF NOT EXISTS origem_pericias (
    origem_id INTEGER REFERENCES origens(id) ON DELETE CASCADE,
    pericia_id INTEGER REFERENCES pericias(id) ON DELETE CASCADE,
    PRIMARY KEY (origem_id, pericia_id)
);

-- Tabela para perícias escolhidas pelo personagem
CREATE TABLE IF NOT EXISTS personagem_pericias (
    personagem_id INTEGER REFERENCES personagens(id) ON DELETE CASCADE,
    pericia_id INTEGER REFERENCES pericias(id) ON DELETE CASCADE,
    fonte VARCHAR(50) NOT NULL, -- 'raca', 'classe', 'origem'
    PRIMARY KEY (personagem_id, pericia_id, fonte)
);

-- Adicionar campo de quantidade de perícias que cada classe oferece
ALTER TABLE classes ADD COLUMN IF NOT EXISTS pericias_quantidade INTEGER DEFAULT 2;
