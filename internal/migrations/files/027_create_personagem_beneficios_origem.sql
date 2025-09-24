-- Tabela para armazenar as perícias escolhidas como benefício de origem
CREATE TABLE IF NOT EXISTS personagem_beneficio_pericias (
    personagem_id INT NOT NULL,
    pericia_id INT NOT NULL,
    PRIMARY KEY (personagem_id, pericia_id),
    FOREIGN KEY (personagem_id) REFERENCES personagens(id) ON DELETE CASCADE,
    FOREIGN KEY (pericia_id) REFERENCES pericias(id) ON DELETE CASCADE
);

-- Tabela para armazenar os poderes escolhidos como benefício de origem
CREATE TABLE IF NOT EXISTS personagem_beneficio_poderes (
    personagem_id INT NOT NULL,
    poder_id INT NOT NULL,
    PRIMARY KEY (personagem_id, poder_id),
    FOREIGN KEY (personagem_id) REFERENCES personagens(id) ON DELETE CASCADE,
    FOREIGN KEY (poder_id) REFERENCES poderes(id) ON DELETE CASCADE
);
