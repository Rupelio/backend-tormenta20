-- 025_add_atributos_livres.sql
-- Adiciona coluna para armazenar atributos livres escolhidos

ALTER TABLE personagens ADD COLUMN IF NOT EXISTS atributos_livres JSONB DEFAULT '[]'::jsonb;

-- Adicionar comentário para documentar o campo
COMMENT ON COLUMN personagens.atributos_livres IS 'Array JSON com os atributos livres escolhidos pelo personagem (ex: ["FOR", "CON"])';