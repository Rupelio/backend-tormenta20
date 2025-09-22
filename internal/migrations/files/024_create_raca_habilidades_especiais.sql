-- Tabela para habilidades especiais de raças
CREATE TABLE IF NOT EXISTS raca_habilidades_especiais (
    id SERIAL PRIMARY KEY,
    raca_id INTEGER REFERENCES racas(id) ON DELETE CASCADE,
    nome VARCHAR(100) NOT NULL,
    descricao TEXT NOT NULL,
    tipo VARCHAR(50) NOT NULL, -- 'atributos_livres', 'deformidade', 'versatilidade'
    opcoes JSONB DEFAULT '{}', -- Opções específicas da habilidade
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Seed para habilidades especiais de raças
INSERT INTO raca_habilidades_especiais (raca_id, nome, descricao, tipo, opcoes) VALUES
-- Humano - Versatilidade
((SELECT id FROM racas WHERE nome = 'Humano'),
 'Versatilidade',
 'Humanos recebem +1 em três atributos diferentes à sua escolha e podem escolher duas perícias treinadas adicionais ou um poder geral.',
 'versatilidade',
 '{"atributos_livres": 3, "pode_escolher": "pericias_ou_poder", "pericias_bonus": 2, "tipos_poder": ["Combate", "Destino", "Magia"]}'),

-- Lefou - Deformidade
((SELECT id FROM racas WHERE nome = 'Lefou'),
 'Deformidade',
 'Todo lefou possui defeitos físicos que, embora desagradáveis, conferem certas vantagens. Você recebe +2 em duas perícias à sua escolha. Cada um desses bônus conta como um poder da Tormenta. Você pode trocar um desses bônus por um poder da Tormenta à sua escolha.',
 'deformidade',
 '{"atributos_livres": 3, "pericias_bonus": 2, "pode_trocar_por_poder": true, "tipos_poder": ["Tormenta"], "max_trocas": 1}'),

-- Osteon - Versatilidade dos Mortos
((SELECT id FROM racas WHERE nome = 'Osteon'),
 'Versatilidade dos Mortos',
 'Osteon recebem +1 em três atributos diferentes à sua escolha e podem escolher duas perícias treinadas adicionais ou um poder geral.',
 'versatilidade',
 '{"atributos_livres": 3, "pode_escolher": "pericias_ou_poder", "pericias_bonus": 2, "tipos_poder": ["Combate", "Destino", "Magia"]}'),

-- Sereia/Tritão - Versatilidade Aquática
((SELECT id FROM racas WHERE nome = 'Sereia/Tritão'),
 'Versatilidade Aquática',
 'Sereias e Tritões recebem +1 em três atributos diferentes à sua escolha e podem escolher duas perícias treinadas adicionais ou um poder geral.',
 'versatilidade',
 '{"atributos_livres": 3, "pode_escolher": "pericias_ou_poder", "pericias_bonus": 2, "tipos_poder": ["Combate", "Destino", "Magia"]}');
