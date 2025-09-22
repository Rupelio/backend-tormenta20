-- 018_create_classe_pericias_automaticas.sql
-- Separar perícias automáticas das classes das perícias de escolha

-- Tabela para perícias automáticas que a classe dá (todo personagem da classe ganha)
CREATE TABLE IF NOT EXISTS classe_pericias_automaticas (
    classe_id INTEGER REFERENCES classes(id) ON DELETE CASCADE,
    pericia_id INTEGER REFERENCES pericias(id) ON DELETE CASCADE,
    PRIMARY KEY (classe_id, pericia_id)
);

-- Renomear a tabela atual classe_pericias para classe_pericias_disponiveis
-- pois ela representa as perícias que podem ser escolhidas, não as automáticas
DROP TABLE IF EXISTS classe_pericias_disponiveis;
ALTER TABLE classe_pericias RENAME TO classe_pericias_disponiveis;
