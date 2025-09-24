-- Fix Suraggel habilidades - inserir habilidades com nome correto das raças
INSERT INTO habilidade_racas (raca_id, nome, descricao, opcional, nivel_minimo) VALUES
((SELECT id FROM racas WHERE nome = 'Suraggel (aggelus)'), 'Herança Divina', 'Você é uma criatura do tipo espírito e recebe visão no escuro.', FALSE, 1),
((SELECT id FROM racas WHERE nome = 'Suraggel (aggelus)'), 'Luz Sagrada', 'Você recebe +2 em Diplomacia e Intuição. Além disso, pode lançar Luz (como uma magia divina: atributo-chave Carisma). Caso aprenda novamente essa magia, seu custo diminui em -1 PM.', FALSE, 1),
((SELECT id FROM racas WHERE nome = 'Suraggel (sulfure)'), 'Herança Divina', 'Você é uma criatura do tipo espírito e recebe visão no escuro.', FALSE, 1),
((SELECT id FROM racas WHERE nome = 'Suraggel (sulfure)'), 'Sombras Profanas', 'Você recebe +2 em Enganação e Furtividade. Além disso, pode lançar Escuridão (como uma magia divina: atributo-chave Inteligência). Caso aprenda novamente essa magia, seu custo diminui em -1 PM.', FALSE, 1);
