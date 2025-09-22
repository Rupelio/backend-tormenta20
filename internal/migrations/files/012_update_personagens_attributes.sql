-- 008_update_personagens_attributes.sql
-- Alterar nomes das colunas de atributos para as abreviações
ALTER TABLE personagens
RENAME COLUMN forca TO "for";

ALTER TABLE personagens
RENAME COLUMN destreza TO des;

ALTER TABLE personagens
RENAME COLUMN constituicao TO con;

ALTER TABLE personagens
RENAME COLUMN inteligencia TO "int";

ALTER TABLE personagens
RENAME COLUMN sabedoria TO sab;

ALTER TABLE personagens
RENAME COLUMN carisma TO car;

-- Adicionar coluna divindade_id se não existir
ALTER TABLE personagens
ADD COLUMN IF NOT EXISTS divindade_id INTEGER REFERENCES divindades(id);

-- Atualizar valores padrão para o novo sistema (0-4)
ALTER TABLE personagens
ALTER COLUMN "for" SET DEFAULT 0;

ALTER TABLE personagens
ALTER COLUMN des SET DEFAULT 0;

ALTER TABLE personagens
ALTER COLUMN con SET DEFAULT 0;

ALTER TABLE personagens
ALTER COLUMN "int" SET DEFAULT 0;

ALTER TABLE personagens
ALTER COLUMN sab SET DEFAULT 0;

ALTER TABLE personagens
ALTER COLUMN car SET DEFAULT 0;
