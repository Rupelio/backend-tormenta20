-- 007_create_divindades.sql
CREATE TABLE IF NOT EXISTS divindades (
    id SERIAL PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    descricao TEXT,
    dominio VARCHAR(255),
    alinhamento VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Inserir dados das divindades de Tormenta20
INSERT INTO divindades (nome, descricao, dominio, alinhamento) VALUES
('Azgher', 'O Deus-Sol, venerado por povos do deserto, viajantes e mercadores honestos. É um deus generoso, mas severo, que combate as trevas e a mentira.', 'Sol, Luz, Honestidade, Justiça', 'Leal e Bom'),
('Valkaria', 'A Deusa da Ambição, criadora dos humanos e atual líder do Panteão. Inspira a evolução, a rebeldia, o desafio a limites e a proteção da liberdade.', 'Ambição, Humanidade, Liberdade', 'Caótico e Bom'),
('Oceano', 'O Deus dos Mares, uma entidade serena e imutável que governa a vastidão de seus domínios. É reverenciado por marinheiros e povos marinhos.', 'Mares, Água, Viagens', 'Neutro'),
('Thyatis', 'O Deus da Ressurreição e Profecia, que representa o perdão, a tolerância e as segundas chances. Seus clérigos possuem dons de profecia e ressurreição.', 'Ressurreição, Profecia, Perdão', 'Bom'),
('Tenebra', 'A Deusa das Trevas e da Noite, mãe de tudo que se move no escuro, desde anões a mortos-vivos. Protege segredos e a prática da necromancia.', 'Trevas, Noite, Mortos-Vivos, Segredos', 'Neutro e Mal'),
('Arsenal', 'O Deus da Guerra, um antigo clérigo que derrotou seu patrono para ascender ao Panteão. Promove o conflito e a vitória a qualquer custo.', 'Guerra, Conflito, Estratégia', 'Leal e Mal'),
('Lena', 'A Deusa da Vida, conhecida como a Deusa-Criança, é a provedora da fertilidade, do sustento e das mais poderosas curas.', 'Vida, Cura, Fertilidade, Maternidade', 'Bom'),
('Wynna', 'A exuberante Deusa da Magia, louvada por todos que usam poder arcano. Generosa, concede magia a todos que a buscam, seja para o bem ou para o mal.', 'Magia, Conhecimento, Fadas', 'Caótico e Bom'),
('Sszzaas', 'O Deus da Traição, a mais inteligente e perigosa das divindades. Promove a mentira, a trapaça e a desconfiança como ferramentas para alcançar objetivos.', 'Traição, Intriga, Serpentes', 'Caótico e Mal'),
('Nimb', 'O insano Deus do Caos, da sorte e do azar. Imprevisível, reverencia a aleatoriedade e o desafio a regras e ao bom senso.', 'Caos, Sorte, Azar, Loucura', 'Caótico e Neutro'),
('Allihanna', 'A Deusa da Natureza, que representa a pureza das plantas e dos animais selvagens. É protetora da vida selvagem e busca a harmonia entre natureza e civilização.', 'Natureza, Animais, Vida Selvagem', 'Neutro e Bom'),
('Megalokk', 'O Deus dos Monstros, uma divindade de selvageria, violência e descontrole. Prega a soberania do mais forte e a eliminação dos fracos.', 'Monstros, Selvageria, Força', 'Caótico e Mal'),
('Khalmyr', 'O Deus da Justiça, antigo líder do Panteão, louvado por aqueles que lutam pela ordem. Defende a lei, combate o crime e oferece perdão e redenção.', 'Justiça, Ordem, Lei, Proteção', 'Leal e Bom'),
('Hyninn', 'O ardiloso Deus da Trapaça, favorito de ladrões, piratas e outros fora da lei. Prega a astúcia, a esperteza e o desafio à ordem para levar vantagem em tudo.', 'Trapaça, Ladrões, Astúcia', 'Caótico e Neutro');
