-- Migration: Tabela de itens iniciais por origem + proficiencias por classe

-- Itens iniciais de cada origem (do livro T20)
CREATE TABLE IF NOT EXISTS origem_itens (
    id SERIAL PRIMARY KEY,
    origem_id INTEGER NOT NULL REFERENCES origens(id) ON DELETE CASCADE,
    nome VARCHAR(200) NOT NULL,
    tipo VARCHAR(50) DEFAULT 'item',
    quantidade INTEGER DEFAULT 1,
    descricao TEXT DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_origem_itens_origem_id ON origem_itens(origem_id);

-- Seed: itens iniciais por origem (extraidos do Livro Basico T20)
INSERT INTO origem_itens (origem_id, nome, tipo, quantidade, descricao) VALUES
-- Acolito
((SELECT id FROM origens WHERE nome='Acólito'), 'Símbolo sagrado', 'equipamento', 1, ''),
((SELECT id FROM origens WHERE nome='Acólito'), 'Traje de sacerdote', 'vestuario', 1, ''),
-- Amigo dos Animais
((SELECT id FROM origens WHERE nome='Amigo dos Animais'), 'Cão de caça, cavalo, pônei ou trobo (escolha um)', 'animal', 1, 'Escolha um animal'),
-- Aristocrata
((SELECT id FROM origens WHERE nome='Aristocrata'), 'Joia de família', 'item', 1, 'Valor de T$ 300'),
((SELECT id FROM origens WHERE nome='Aristocrata'), 'Traje da corte', 'vestuario', 1, ''),
-- Artesao
((SELECT id FROM origens WHERE nome='Artesão'), 'Instrumentos de ofício', 'ferramenta', 1, 'Qualquer tipo'),
((SELECT id FROM origens WHERE nome='Artesão'), 'Item fabricável', 'item', 1, 'Um item que você possa fabricar de até T$ 50'),
-- Artista
((SELECT id FROM origens WHERE nome='Artista'), 'Estojo de disfarces ou instrumento musical', 'ferramenta', 1, 'Escolha um'),
-- Assistente de Laboratorio
((SELECT id FROM origens WHERE nome='Assistente de Laboratório'), 'Instrumentos de ofício (alquimista)', 'ferramenta', 1, ''),
-- Batedor
((SELECT id FROM origens WHERE nome='Batedor'), 'Barraca', 'equipamento', 1, ''),
((SELECT id FROM origens WHERE nome='Batedor'), 'Equipamento de viagem', 'equipamento', 1, ''),
((SELECT id FROM origens WHERE nome='Batedor'), 'Arma simples ou marcial de ataque à distância', 'arma', 1, 'Escolha uma'),
-- Capanga
((SELECT id FROM origens WHERE nome='Capanga'), 'Tatuagem ou adereço de gangue', 'item', 1, '+1 em Intimidação'),
((SELECT id FROM origens WHERE nome='Capanga'), 'Arma simples corpo a corpo', 'arma', 1, 'Escolha uma'),
-- Charlatao
((SELECT id FROM origens WHERE nome='Charlatão'), 'Estojo de disfarces', 'ferramenta', 1, ''),
((SELECT id FROM origens WHERE nome='Charlatão'), 'Joia falsificada', 'item', 1, 'Valor aparente T$ 100, sem valor real'),
-- Circense
((SELECT id FROM origens WHERE nome='Circense'), 'Três bolas coloridas', 'item', 1, '+1 em Atuação'),
-- Criminoso
((SELECT id FROM origens WHERE nome='Criminoso'), 'Estojo de disfarces ou gazua', 'ferramenta', 1, 'Escolha um'),
-- Curandeiro
((SELECT id FROM origens WHERE nome='Curandeiro'), 'Bálsamo restaurador', 'alquimico', 2, ''),
((SELECT id FROM origens WHERE nome='Curandeiro'), 'Maleta de medicamentos', 'ferramenta', 1, ''),
-- Eremita
((SELECT id FROM origens WHERE nome='Eremita'), 'Barraca', 'equipamento', 1, ''),
((SELECT id FROM origens WHERE nome='Eremita'), 'Equipamento de viagem', 'equipamento', 1, ''),
-- Escravo
((SELECT id FROM origens WHERE nome='Escravo'), 'Algemas', 'equipamento', 1, ''),
((SELECT id FROM origens WHERE nome='Escravo'), 'Ferramenta pesada', 'arma', 1, 'Mesmas estatísticas de uma maça'),
-- Estudioso
((SELECT id FROM origens WHERE nome='Estudioso'), 'Coleção de livros', 'ferramenta', 1, '+1 em Conhecimento, Guerra, Misticismo ou Nobreza'),
-- Fazendeiro
((SELECT id FROM origens WHERE nome='Fazendeiro'), 'Carroça', 'veiculo', 1, ''),
((SELECT id FROM origens WHERE nome='Fazendeiro'), 'Ferramenta agrícola', 'arma', 1, 'Mesmas estatísticas de uma lança'),
((SELECT id FROM origens WHERE nome='Fazendeiro'), 'Ração de viagem', 'item', 10, ''),
((SELECT id FROM origens WHERE nome='Fazendeiro'), 'Animal não combativo', 'animal', 1, 'Galinha, porco ou ovelha'),
-- Forasteiro
((SELECT id FROM origens WHERE nome='Forasteiro'), 'Equipamento de viagem', 'equipamento', 1, ''),
((SELECT id FROM origens WHERE nome='Forasteiro'), 'Instrumento musical exótico', 'ferramenta', 1, '+1 em uma perícia de Carisma'),
((SELECT id FROM origens WHERE nome='Forasteiro'), 'Traje estrangeiro', 'vestuario', 1, ''),
-- Gladiador
((SELECT id FROM origens WHERE nome='Gladiador'), 'Arma marcial ou exótica', 'arma', 1, 'Escolha uma'),
((SELECT id FROM origens WHERE nome='Gladiador'), 'Item de admirador', 'item', 1, 'Sem valor monetário'),
-- Guarda
((SELECT id FROM origens WHERE nome='Guarda'), 'Apito', 'item', 1, ''),
((SELECT id FROM origens WHERE nome='Guarda'), 'Insígnia da milícia', 'item', 1, ''),
((SELECT id FROM origens WHERE nome='Guarda'), 'Arma marcial', 'arma', 1, 'Escolha uma'),
-- Herdeiro
((SELECT id FROM origens WHERE nome='Herdeiro'), 'Símbolo de herança', 'item', 1, 'Anel de sinete ou manto cerimonial'),
-- Heroi Campones
((SELECT id FROM origens WHERE nome='Herói Camponês'), 'Instrumentos de ofício ou arma simples', 'item', 1, 'Escolha um'),
((SELECT id FROM origens WHERE nome='Herói Camponês'), 'Traje de plebeu', 'vestuario', 1, ''),
-- Marujo
((SELECT id FROM origens WHERE nome='Marujo'), 'Corda', 'equipamento', 1, '10 metros'),
-- Mateiro
((SELECT id FROM origens WHERE nome='Mateiro'), 'Arco curto', 'arma', 1, ''),
((SELECT id FROM origens WHERE nome='Mateiro'), 'Barraca', 'equipamento', 1, ''),
((SELECT id FROM origens WHERE nome='Mateiro'), 'Equipamento de viagem', 'equipamento', 1, ''),
((SELECT id FROM origens WHERE nome='Mateiro'), 'Flechas', 'municao', 20, ''),
-- Membro de Guilda
((SELECT id FROM origens WHERE nome='Membro de Guilda'), 'Gazua ou instrumentos de ofício', 'ferramenta', 1, 'Escolha um'),
-- Mercador
((SELECT id FROM origens WHERE nome='Mercador'), 'Carroça', 'veiculo', 1, ''),
((SELECT id FROM origens WHERE nome='Mercador'), 'Trobo', 'animal', 1, ''),
((SELECT id FROM origens WHERE nome='Mercador'), 'Mercadorias', 'item', 1, 'Valor de T$ 100'),
-- Minerador
((SELECT id FROM origens WHERE nome='Minerador'), 'Gemas preciosas', 'item', 1, 'Valor de T$ 100'),
((SELECT id FROM origens WHERE nome='Minerador'), 'Picareta', 'arma', 1, ''),
-- Nomade
((SELECT id FROM origens WHERE nome='Nômade'), 'Bordão', 'arma', 1, ''),
((SELECT id FROM origens WHERE nome='Nômade'), 'Equipamento de viagem', 'equipamento', 1, ''),
-- Pivete
((SELECT id FROM origens WHERE nome='Pivete'), 'Gazua', 'ferramenta', 1, ''),
((SELECT id FROM origens WHERE nome='Pivete'), 'Traje de plebeu', 'vestuario', 1, ''),
((SELECT id FROM origens WHERE nome='Pivete'), 'Animal urbano', 'animal', 1, 'Cão, gato, rato ou pombo'),
-- Refugiado
((SELECT id FROM origens WHERE nome='Refugiado'), 'Item estrangeiro', 'item', 1, 'Até T$ 100'),
-- Seguidor
((SELECT id FROM origens WHERE nome='Seguidor'), 'Item do mestre', 'item', 1, 'Até T$ 100'),
-- Selvagem
((SELECT id FROM origens WHERE nome='Selvagem'), 'Arma simples', 'arma', 1, 'Escolha uma'),
((SELECT id FROM origens WHERE nome='Selvagem'), 'Animal de estimação', 'animal', 1, 'Pássaro ou esquilo'),
-- Soldado
((SELECT id FROM origens WHERE nome='Soldado'), 'Arma marcial', 'arma', 1, 'Escolha uma'),
((SELECT id FROM origens WHERE nome='Soldado'), 'Uniforme militar', 'vestuario', 1, ''),
((SELECT id FROM origens WHERE nome='Soldado'), 'Insígnia do exército', 'item', 1, ''),
-- Taverneiro
((SELECT id FROM origens WHERE nome='Taverneiro'), 'Clava', 'arma', 1, 'Rolo de macarrão ou martelo de carne'),
((SELECT id FROM origens WHERE nome='Taverneiro'), 'Panela', 'item', 1, ''),
((SELECT id FROM origens WHERE nome='Taverneiro'), 'Avental', 'item', 1, ''),
((SELECT id FROM origens WHERE nome='Taverneiro'), 'Caneca', 'item', 1, ''),
((SELECT id FROM origens WHERE nome='Taverneiro'), 'Pano sujo', 'item', 1, ''),
-- Trabalhador
((SELECT id FROM origens WHERE nome='Trabalhador'), 'Ferramenta pesada', 'arma', 1, 'Mesmas estatísticas de maça ou lança');

-- Proficiencias por classe
ALTER TABLE classes ADD COLUMN IF NOT EXISTS prof_armas_simples BOOLEAN DEFAULT true;
ALTER TABLE classes ADD COLUMN IF NOT EXISTS prof_armas_marciais BOOLEAN DEFAULT false;
ALTER TABLE classes ADD COLUMN IF NOT EXISTS prof_armaduras_leves BOOLEAN DEFAULT false;
ALTER TABLE classes ADD COLUMN IF NOT EXISTS prof_armaduras_pesadas BOOLEAN DEFAULT false;
ALTER TABLE classes ADD COLUMN IF NOT EXISTS prof_escudos BOOLEAN DEFAULT false;

-- Seed proficiencias (Livro Basico T20)
-- Guerreiro: todas
UPDATE classes SET prof_armas_simples=true, prof_armas_marciais=true, prof_armaduras_leves=true, prof_armaduras_pesadas=true, prof_escudos=true WHERE nome='Guerreiro';
-- Cavaleiro: todas
UPDATE classes SET prof_armas_simples=true, prof_armas_marciais=true, prof_armaduras_leves=true, prof_armaduras_pesadas=true, prof_escudos=true WHERE nome='Cavaleiro';
-- Paladino: todas
UPDATE classes SET prof_armas_simples=true, prof_armas_marciais=true, prof_armaduras_leves=true, prof_armaduras_pesadas=true, prof_escudos=true WHERE nome='Paladino';
-- Barbaro: simples+marciais, leves, escudos
UPDATE classes SET prof_armas_simples=true, prof_armas_marciais=true, prof_armaduras_leves=true, prof_armaduras_pesadas=false, prof_escudos=true WHERE nome='Bárbaro';
-- Lutador: simples+marciais, leves
UPDATE classes SET prof_armas_simples=true, prof_armas_marciais=true, prof_armaduras_leves=true, prof_armaduras_pesadas=false, prof_escudos=false WHERE nome='Lutador';
-- Bucaneiro: simples+marciais, leves
UPDATE classes SET prof_armas_simples=true, prof_armas_marciais=true, prof_armaduras_leves=true, prof_armaduras_pesadas=false, prof_escudos=false WHERE nome='Bucaneiro';
-- Cacador: simples+marciais, leves
UPDATE classes SET prof_armas_simples=true, prof_armas_marciais=true, prof_armaduras_leves=true, prof_armaduras_pesadas=false, prof_escudos=true WHERE nome='Caçador';
-- Clerigo: simples+marciais, leves+pesadas, escudos
UPDATE classes SET prof_armas_simples=true, prof_armas_marciais=true, prof_armaduras_leves=true, prof_armaduras_pesadas=true, prof_escudos=true WHERE nome='Clérigo';
-- Druida: simples, leves, escudos
UPDATE classes SET prof_armas_simples=true, prof_armas_marciais=false, prof_armaduras_leves=true, prof_armaduras_pesadas=false, prof_escudos=true WHERE nome='Druida';
-- Nobre: simples+marciais, leves, escudos
UPDATE classes SET prof_armas_simples=true, prof_armas_marciais=true, prof_armaduras_leves=true, prof_armaduras_pesadas=false, prof_escudos=true WHERE nome='Nobre';
-- Bardo: simples+marciais, leves
UPDATE classes SET prof_armas_simples=true, prof_armas_marciais=true, prof_armaduras_leves=true, prof_armaduras_pesadas=false, prof_escudos=false WHERE nome='Bardo';
-- Inventor: simples, leves
UPDATE classes SET prof_armas_simples=true, prof_armas_marciais=false, prof_armaduras_leves=true, prof_armaduras_pesadas=false, prof_escudos=false WHERE nome='Inventor';
-- Ladino: simples, leves
UPDATE classes SET prof_armas_simples=true, prof_armas_marciais=false, prof_armaduras_leves=true, prof_armaduras_pesadas=false, prof_escudos=false WHERE nome='Ladino';
-- Arcanista: simples, nenhuma armadura
UPDATE classes SET prof_armas_simples=true, prof_armas_marciais=false, prof_armaduras_leves=false, prof_armaduras_pesadas=false, prof_escudos=false WHERE nome='Arcanista';
