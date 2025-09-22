-- Migração para adicionar suporte a escolhas de raça
-- (perícias ou poderes especiais)

-- Adicionar campos para controlar escolhas especiais de raça
ALTER TABLE personagens ADD COLUMN IF NOT EXISTS escolhas_raca JSONB DEFAULT '{}';

-- Comentários explicativos:
-- escolhas_raca será um campo JSON que armazenará as escolhas específicas de cada raça
-- Exemplos:
-- Humano: {"pericias_escolhidas": [1,2], "poder_escolhido": null}
-- Lefou: {"pericias_escolhidas": [1], "poder_tormenta_escolhido": 5, "deformidade_trocada": true}
-- Osteon: {"pericias_escolhidas": [3,4], "poder_escolhido": null}
