-- Populando a tabela de habilidades de divindades (Poderes Concedidos)
INSERT INTO habilidade_divindades (divindade_id, nome, descricao, nivel, opcional) VALUES
-- Aharadak
((SELECT id FROM divindades WHERE nome = 'Aharadak'), 'Afinidade com a Tormenta', 'Você recebe +10 em testes de resistência contra efeitos da Tormenta, de suas criaturas e de devotos de Aharadak. Além disso, seu primeiro poder da Tormenta não conta para perda de Carisma.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Aharadak'), 'Êxtase da Loucura', 'Toda vez que uma ou mais criaturas falham em um teste de Vontade contra uma de suas habilidades mágicas, você recebe 1 PM temporário cumulativo. Você pode ganhar um máximo de PM temporários por cena desta forma igual a sua Sabedoria.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Aharadak'), 'Percepção Temporal', 'Você pode gastar 3 PM para somar sua Sabedoria (limitado por seu nível e não cumulativo com efeitos que somam este atributo) a seus ataques, Defesa e testes de Reflexos até o fim da cena.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Aharadak'), 'Rejeição Divina', 'Você recebe resistência a magia divina +5.', 1, TRUE),

-- Allihanna
((SELECT id FROM divindades WHERE nome = 'Allihanna'), 'Compreender os Ermos', 'Você recebe +2 em Sobrevivência e pode usar Sabedoria para Adestramento (em vez de Carisma).', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Allihanna'), 'Dedo Verde', 'Você aprende e pode lançar Controlar Plantas. Caso aprenda novamente essa magia, seu custo diminui em –1 PM.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Allihanna'), 'Descanso Natural', 'Para você, dormir ao relento conta como condição de descanso confortável.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Allihanna'), 'Voz da Natureza', 'Você pode falar com animais (como o efeito da magia Voz Divina) e aprende e pode lançar Acalmar Animal, mas só contra animais. Caso aprenda novamente essa magia, seu custo diminui em –1 PM.', 1, TRUE),

-- Arsenal
((SELECT id FROM divindades WHERE nome = 'Arsenal'), 'Conjurar Arma', 'Você pode gastar 1 PM para invocar uma arma corpo a corpo ou de arremesso com a qual seja proficiente. A arma surge em sua mão, fornece +1 em testes de ataque e rolagens de dano, é considerada mágica e dura pela cena.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Arsenal'), 'Coragem Total', 'Você é imune a efeitos de medo, mágicos ou não.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Arsenal'), 'Fé Guerreira', 'Você pode usar Sabedoria para Guerra (em vez de Inteligência). Além disso, em combate, quando vai fazer um teste de perícia, você pode gastar 2 PM para substituí-lo por um teste de Guerra.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Arsenal'), 'Sangue de Ferro', 'Você pode pagar 3 PM para receber +2 em rolagens de dano e redução de dano 5 até o fim da cena.', 1, TRUE),

-- Azgher
((SELECT id FROM divindades WHERE nome = 'Azgher'), 'Espada Solar', 'Você pode gastar 1 PM para fazer uma arma corpo a corpo de corte que esteja empunhando causar +1d6 de dano por fogo até o fim da cena.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Azgher'), 'Fulgor Solar', 'Você recebe redução de frio e trevas 5. Além disso, quando é alvo de um ataque você pode gastar 1 PM para emitir um clarão solar que deixa o atacante ofuscado por uma rodada.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Azgher'), 'Habitante do Deserto', 'Você recebe redução de fogo 10 e pode pagar 1 PM para criar água pura e potável suficiente para um odre.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Azgher'), 'Inimigo de Tenebra', 'Seus ataques e habilidades causam +1d6 pontos de dano contra mortos-vivos. Quando você usa um efeito que gera luz, o alcance da iluminação dobra.', 1, TRUE),

-- Hyninn
((SELECT id FROM divindades WHERE nome = 'Hyninn'), 'Apostar com o Trapaceiro', 'Quando faz um teste de perícia, você pode gastar 1 PM para apostar com Hyninn. Você e o mestre rolam 1d20, mas o mestre mantém o resultado dele em segredo. Você então escolhe entre usar seu próprio resultado ou o resultado oculto do mestre.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Hyninn'), 'Farsa do Fingidor', 'Você aprende e pode lançar Criar Ilusão. Caso aprenda novamente essa magia, seu custo diminui em –1 PM.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Hyninn'), 'Forma de Macaco', 'Você pode gastar uma ação completa e 2 PM para se transformar em um macaco (tamanho Minúsculo, deslocamento de escalar 9m). A transformação dura indefinidamente, mas termina caso você faça um ataque, lance uma magia ou sofra dano.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Hyninn'), 'Golpista Divino', 'Você recebe +2 em Enganação, Jogatina e Ladinagem.', 1, TRUE),

-- Kallyadranoch
((SELECT id FROM divindades WHERE nome = 'Kallyadranoch'), 'Aura de Medo', 'Você pode gastar 2 PM para gerar uma aura de medo de 9m de raio e duração até o fim da cena. Inimigos que entrem na aura devem fazer um teste de Vontade (CD Car) ou ficam abalados.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Kallyadranoch'), 'Escamas Dracônicas', 'Você recebe +2 na Defesa e em Fortitude.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Kallyadranoch'), 'Presas Primordiais', 'Você pode gastar 1 PM para ganhar uma arma natural de mordida (dano 1d6) pela cena. Se já possui mordida, o dano aumenta em dois passos.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Kallyadranoch'), 'Servos do Dragão', 'Você pode gastar uma ação completa e 2 PM para invocar 2d4+1 kobolds capangas em alcance curto que desaparecem no fim da cena.', 1, TRUE),

-- Khalmyr
((SELECT id FROM divindades WHERE nome = 'Khalmyr'), 'Coragem Total', 'Você é imune a efeitos de medo, mágicos ou não.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Khalmyr'), 'Dom da Verdade', 'Você pode pagar 2 PM para receber +5 em testes de Intuição, e em testes de Percepção contra Enganação e Furtividade, até o fim da cena.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Khalmyr'), 'Espada Justiceira', 'Você pode gastar 1 PM para encantar sua espada (ou outra arma corpo a corpo de corte). Ela tem seu dano aumentado em um passo até o fim da cena.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Khalmyr'), 'Reparar Injustiça', 'Uma vez por rodada, quando um oponente em alcance curto acerta um ataque em você ou em um de seus aliados, você pode gastar 2 PM para fazer este oponente repetir o ataque, escolhendo o pior resultado.', 1, TRUE),

-- Lena
((SELECT id FROM divindades WHERE nome = 'Lena'), 'Ataque Piedoso', 'Você pode usar armas corpo a corpo para causar dano não letal sem sofrer a penalidade de –5 no teste de ataque.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Lena'), 'Aura Restauradora', 'Efeitos de cura usados por você e seus aliados em um raio de 9m recuperam +1 PV por dado.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Lena'), 'Cura Gentil', 'Você soma seu Carisma aos PV restaurados por seus efeitos mágicos de cura.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Lena'), 'Curandeira Perfeita', 'Você sempre pode escolher 10 em testes de Cura e não sofre penalidade por usar essa perícia sem uma maleta de medicamentos. Com o item, recebe +2 no teste de Cura (ou +5, se aprimorado).', 1, TRUE),

-- Lin-Wu
((SELECT id FROM divindades WHERE nome = 'Lin-Wu'), 'Coragem Total', 'Você é imune a efeitos de medo, mágicos ou não.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Lin-Wu'), 'Kiai Divino', 'Uma vez por rodada, quando faz um ataque corpo a corpo, você pode pagar 3 PM. Se acertar o ataque, causa dano máximo, sem necessidade de rolar dados.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Lin-Wu'), 'Mente Vazia', 'Você recebe +2 em Iniciativa, Percepção e Vontade.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Lin-Wu'), 'Tradição de Lin-Wu', 'Você considera a katana uma arma simples e, se for proficiente em armas marciais, recebe +1 na margem de ameaça com ela.', 1, TRUE),

-- Marah
((SELECT id FROM divindades WHERE nome = 'Marah'), 'Aura de Paz', 'Você pode gastar 2 PM para gerar uma aura de paz com 9m de raio por uma cena. Qualquer inimigo na aura que tente uma ação hostil contra você deve fazer um teste de Vontade (CD Car) ou perderá a ação.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Marah'), 'Dom da Esperança', 'Você soma sua Sabedoria em seus PV em vez de Constituição, e se torna imune às condições alquebrado, esmorecido e frustrado.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Marah'), 'Palavras de Bondade', 'Você aprende e pode lançar Enfeitiçar. Caso aprenda novamente essa magia, seu custo diminui em –1 PM.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Marah'), 'Talento Artístico', 'Você recebe +2 em Acrobacia, Atuação e Diplomacia.', 1, TRUE),

-- Megalokk
((SELECT id FROM divindades WHERE nome = 'Megalokk'), 'Olhar Amedrontador', 'Você aprende e pode lançar Amedrontar. Caso aprenda novamente essa magia, seu custo diminui em –1 PM.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Megalokk'), 'Presas Primordiais', 'Você pode gastar 1 PM para ganhar uma arma natural de mordida (dano 1d6) pela cena. Se já possui mordida, o dano aumenta em dois passos.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Megalokk'), 'Urro Divino', 'Quando faz um ataque ou lança uma magia, você pode pagar 1 PM para somar sua Constituição (mínimo +1) à rolagem de dano.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Megalokk'), 'Voz dos Monstros', 'Você conhece os idiomas de todos os monstros inteligentes e pode se comunicar com monstros não inteligentes.', 1, TRUE),

-- Nimb
((SELECT id FROM divindades WHERE nome = 'Nimb'), 'Êxtase da Loucura', 'Toda vez que uma ou mais criaturas falham em um teste de Vontade contra uma de suas habilidades mágicas, você recebe 1 PM temporário cumulativo. Máximo por cena igual à sua Sabedoria.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Nimb'), 'Poder Oculto', 'Você pode gastar uma ação de movimento e 2 PM para invocar a força, a rapidez ou o vigor dos loucos. Role 1d6 para receber +2 em Força (1-2), Destreza (3-4) ou Constituição (5-6) até o fim da cena.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Nimb'), 'Sorte dos Loucos', 'Quando faz um teste, você pode pagar 1 PM para rolá-lo novamente. Se ainda assim falhar, perde 1d6 PM para cada vez que utilizou este poder neste teste.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Nimb'), 'Transmissão da Loucura', 'Você pode lançar Sussurros Insanos (CD Car). Caso aprenda novamente essa magia, seu custo diminui em –1 PM.', 1, TRUE),

-- Oceano
((SELECT id FROM divindades WHERE nome = 'Oceano'), 'Anfíbio', 'Você pode respirar embaixo d''água e adquire deslocamento de natação igual a seu deslocamento terrestre.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Oceano'), 'Arsenal das Profundezas', 'Você recebe +2 nas rolagens de dano com azagaias, lanças e tridentes e seu multiplicador de crítico com essas armas aumenta em +1.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Oceano'), 'Mestre dos Mares', 'Você pode falar com animais aquáticos e aprende a magia Acalmar Animal (apenas contra criaturas aquáticas). Caso aprenda novamente essa magia, seu custo diminui em –1 PM.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Oceano'), 'Sopro do Mar', 'Você pode gastar uma ação padrão e 1 PM para soprar um cone de 6m, causando 2d6 de dano de frio (Reflexos CD Sab reduz à metade). Você pode aprender Sopro das Uivantes como magia divina, com custo reduzido em -1 PM.', 1, TRUE),

-- Sszzaas
((SELECT id FROM divindades WHERE nome = 'Sszzaas'), 'Astúcia da Serpente', 'Você recebe +2 em Enganação, Furtividade e Intuição.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Sszzaas'), 'Familiar Ofídico', 'Você recebe um familiar cobra que não conta em seu limite de parceiros.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Sszzaas'), 'Presas Venenosas', 'Você pode gastar uma ação de movimento e 1 PM para envenenar uma arma corpo a corpo. Em caso de acerto, a arma causa perda de 1d12 pontos de vida.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Sszzaas'), 'Sangue Ofídico', 'Você recebe resistência a veneno +5 e a CD para resistir aos seus venenos aumenta em +2.', 1, TRUE),

-- Tanna-Toh
((SELECT id FROM divindades WHERE nome = 'Tanna-Toh'), 'Conhecimento Enciclopédico', 'Você se torna treinado em duas perícias baseadas em Inteligência a sua escolha.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Tanna-Toh'), 'Mente Analítica', 'Você recebe +2 em Intuição, Investigação e Vontade.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Tanna-Toh'), 'Pesquisa Abençoada', 'Se passar uma hora pesquisando, pode rolar novamente um teste de perícia baseada em Inteligência ou Sabedoria que tenha feito desde a última cena.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Tanna-Toh'), 'Voz da Civilização', 'Você está sempre sob efeito de Compreensão.', 1, TRUE),

-- Tenebra
((SELECT id FROM divindades WHERE nome = 'Tenebra'), 'Carícia Sombria', 'Você pode gastar 1 PM e uma ação padrão para tocar uma criatura, causando 2d6 de dano de trevas (Fortitude CD Sab reduz à metade) e recuperando PV iguais à metade do dano causado.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Tenebra'), 'Manto da Penumbra', 'Você aprende e pode lançar Escuridão. Caso aprenda novamente essa magia, seu custo diminui em –1 PM.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Tenebra'), 'Visão nas Trevas', 'Você enxerga perfeitamente no escuro, incluindo em magias de escuridão.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Tenebra'), 'Zumbificar', 'Você pode gastar uma ação completa e 3 PM para reanimar um cadáver Pequeno ou Médio como um parceiro iniciante (combatente, fortão ou guardião).', 1, TRUE),

-- Thwor
((SELECT id FROM divindades WHERE nome = 'Thwor'), 'Almejar o Impossível', 'Quando faz um teste de perícia, um resultado de 19 ou mais no dado sempre é um sucesso.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Thwor'), 'Fúria Divina', 'Você pode gastar 2 PM para receber +2 em testes de ataque e rolagens de dano corpo a corpo pela cena. Se usar com a habilidade Fúria, ela também dura uma cena.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Thwor'), 'Olhar Amedrontador', 'Você aprende e pode lançar Amedrontar. Caso aprenda novamente essa magia, seu custo diminui em –1 PM.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Thwor'), 'Tropas Duyshidakk', 'Você pode gastar uma ação completa e 2 PM para invocar 1d4+1 goblinoides capangas em alcance curto que desaparecem no fim da cena.', 1, TRUE),

-- Thyatis
((SELECT id FROM divindades WHERE nome = 'Thyatis'), 'Ataque Piedoso', 'Você pode usar armas corpo a corpo para causar dano não letal sem sofrer a penalidade de –5 no teste de ataque.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Thyatis'), 'Dom da Imortalidade', 'Você é imortal. Sempre que morre, volta à vida após 3d6 dias. Apenas paladinos podem escolher este poder.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Thyatis'), 'Dom da Profecia', 'Você pode lançar Augúrio. Caso aprenda novamente essa magia, seu custo diminui em –1 PM. Você também pode gastar 2 PM para receber +2 em um teste.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Thyatis'), 'Dom da Ressurreição', 'Você pode gastar uma ação completa e todos os seus PM para tocar um cadáver e ressuscitá-lo. A criatura volta com 1 PV e 0 PM, e perde 1 de Constituição permanentemente. Apenas clérigos podem escolher este poder.', 1, TRUE),

-- Valkaria
((SELECT id FROM divindades WHERE nome = 'Valkaria'), 'Almejar o Impossível', 'Quando faz um teste de perícia, um resultado de 19 ou mais no dado sempre é um sucesso.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Valkaria'), 'Armas da Ambição', 'Você recebe +1 em testes de ataque e na margem de ameaça com armas nas quais é proficiente.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Valkaria'), 'Coragem Total', 'Você é imune a efeitos de medo, mágicos ou não.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Valkaria'), 'Liberdade Divina', 'Você pode gastar 2 PM para receber imunidade a efeitos de movimento por uma rodada.', 1, TRUE),

-- Wynna
((SELECT id FROM divindades WHERE nome = 'Wynna'), 'Bênção do Mana', 'Você recebe +1 PM a cada nível ímpar.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Wynna'), 'Centelha Mágica', 'Escolha uma magia arcana ou divina de 1º círculo. Você aprende e pode lançar essa magia.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Wynna'), 'Escudo Mágico', 'Quando lança uma magia, você recebe um bônus na Defesa igual ao círculo da magia lançada até o início do seu próximo turno.', 1, TRUE),
((SELECT id FROM divindades WHERE nome = 'Wynna'), 'Teurgista Místico', 'Até uma magia de cada círculo que você aprender poderá ser escolhida entre magias divinas (se for conjurador arcano) ou arcanas (se for divino).', 1, TRUE);
