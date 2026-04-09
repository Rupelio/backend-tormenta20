-- Migration: Adicionar PV/PM primeiro nivel nas classes e campos extras nos personagens
-- Em T20, o 1o nivel tem PV/PM diferente dos niveis subsequentes

ALTER TABLE classes ADD COLUMN IF NOT EXISTS pv_primeiro_nivel INTEGER;
ALTER TABLE classes ADD COLUMN IF NOT EXISTS pm_primeiro_nivel INTEGER;

-- Valores do T20 (PV primeiro nivel = dado maximo da classe)
-- PV por nivel = metade do dado + 1 (valor fixo padrao)
-- PM precisa ser validado com o livro

-- Barbaro: d12 PV
UPDATE classes SET pv_primeiro_nivel = 24, pv_por_nivel = 6 WHERE nome = 'Bárbaro';
-- Guerreiro: d10 PV
UPDATE classes SET pv_primeiro_nivel = 20, pv_por_nivel = 5 WHERE nome = 'Guerreiro';
-- Cavaleiro: d10 PV
UPDATE classes SET pv_primeiro_nivel = 20, pv_por_nivel = 5 WHERE nome = 'Cavaleiro';
-- Paladino: d10 PV
UPDATE classes SET pv_primeiro_nivel = 20, pv_por_nivel = 5 WHERE nome = 'Paladino';
-- Lutador: d10 PV
UPDATE classes SET pv_primeiro_nivel = 20, pv_por_nivel = 5 WHERE nome = 'Lutador';
-- Bucaneiro: d8 PV
UPDATE classes SET pv_primeiro_nivel = 16, pv_por_nivel = 4 WHERE nome = 'Bucaneiro';
-- Cacador: d8 PV
UPDATE classes SET pv_primeiro_nivel = 16, pv_por_nivel = 4 WHERE nome = 'Caçador';
-- Clerigo: d8 PV
UPDATE classes SET pv_primeiro_nivel = 16, pv_por_nivel = 4 WHERE nome = 'Clérigo';
-- Druida: d8 PV
UPDATE classes SET pv_primeiro_nivel = 16, pv_por_nivel = 4 WHERE nome = 'Druida';
-- Nobre: d8 PV
UPDATE classes SET pv_primeiro_nivel = 16, pv_por_nivel = 4 WHERE nome = 'Nobre';
-- Bardo: d6 PV
UPDATE classes SET pv_primeiro_nivel = 12, pv_por_nivel = 3 WHERE nome = 'Bardo';
-- Inventor: d8 PV
UPDATE classes SET pv_primeiro_nivel = 12, pv_por_nivel = 3 WHERE nome = 'Inventor';
-- Ladino: d8 PV
UPDATE classes SET pv_primeiro_nivel = 12, pv_por_nivel = 3 WHERE nome = 'Ladino';
-- Arcanista: d6 PV
UPDATE classes SET pv_primeiro_nivel = 8, pv_por_nivel = 2 WHERE nome = 'Arcanista';

-- PM primeiro nivel (valores precisam ser validados com o livro T20)
UPDATE classes SET pm_primeiro_nivel = pm_por_nivel WHERE pm_primeiro_nivel IS NULL;

-- Campos extras no personagem
ALTER TABLE personagens ADD COLUMN IF NOT EXISTS dinheiro DECIMAL(10,2) DEFAULT 0;
ALTER TABLE personagens ADD COLUMN IF NOT EXISTS anotacoes TEXT DEFAULT '';
ALTER TABLE personagens ADD COLUMN IF NOT EXISTS historico TEXT DEFAULT '';

-- Tabela de itens do personagem
CREATE TABLE IF NOT EXISTS personagem_itens (
    id SERIAL PRIMARY KEY,
    personagem_id INTEGER NOT NULL REFERENCES personagens(id) ON DELETE CASCADE,
    nome VARCHAR(200) NOT NULL,
    tipo VARCHAR(50) DEFAULT 'item',
    quantidade INTEGER DEFAULT 1,
    peso DECIMAL(8,2) DEFAULT 0,
    valor DECIMAL(10,2) DEFAULT 0,
    descricao TEXT DEFAULT '',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_personagem_itens_personagem_id ON personagem_itens(personagem_id);
