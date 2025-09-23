-- Migration: Add user identification to personagens table
-- Adiciona campos para identificar o usuário/sessão que criou o personagem

ALTER TABLE personagens
ADD COLUMN user_session_id VARCHAR(36),
ADD COLUMN user_ip INET,
ADD COLUMN created_by_type VARCHAR(20) DEFAULT 'session';

-- Adiciona comentários para documentar os campos
COMMENT ON COLUMN personagens.user_session_id IS 'UUID da sessão do usuário (cookie)';
COMMENT ON COLUMN personagens.user_ip IS 'IP do usuário que criou o personagem';
COMMENT ON COLUMN personagens.created_by_type IS 'Tipo de identificação: session, ip, hybrid';

-- Índice para melhorar performance nas consultas
CREATE INDEX idx_personagens_user_session_id ON personagens(user_session_id);
CREATE INDEX idx_personagens_user_ip ON personagens(user_ip);
CREATE INDEX idx_personagens_user_identification ON personagens(user_session_id, user_ip);
