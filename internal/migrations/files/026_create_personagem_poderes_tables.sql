-- Migration para criar tabelas de relacionamento de poderes
-- Up migration

-- Tabela para relacionar personagens com poderes divinos
CREATE TABLE IF NOT EXISTS personagem_poderes_divinos (
    personagem_id INTEGER NOT NULL,
    poder_id INTEGER NOT NULL,
    nivel INTEGER DEFAULT 1,
    PRIMARY KEY (personagem_id, poder_id),
    FOREIGN KEY (personagem_id) REFERENCES personagens(id) ON DELETE CASCADE,
    FOREIGN KEY (poder_id) REFERENCES poderes(id) ON DELETE CASCADE
);

-- Tabela para relacionar personagens com poderes de classe
CREATE TABLE IF NOT EXISTS personagem_poderes_classe (
    personagem_id INTEGER NOT NULL,
    poder_id INTEGER NOT NULL,
    nivel INTEGER DEFAULT 1,
    PRIMARY KEY (personagem_id, poder_id),
    FOREIGN KEY (personagem_id) REFERENCES personagens(id) ON DELETE CASCADE,
    FOREIGN KEY (poder_id) REFERENCES poderes(id) ON DELETE CASCADE
);

-- Índices para melhor performance
CREATE INDEX IF NOT EXISTS idx_personagem_poderes_divinos_personagem ON personagem_poderes_divinos(personagem_id);
CREATE INDEX IF NOT EXISTS idx_personagem_poderes_divinos_poder ON personagem_poderes_divinos(poder_id);
CREATE INDEX IF NOT EXISTS idx_personagem_poderes_classe_personagem ON personagem_poderes_classe(personagem_id);
CREATE INDEX IF NOT EXISTS idx_personagem_poderes_classe_poder ON personagem_poderes_classe(poder_id);
