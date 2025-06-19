
CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE IF NOT EXISTS documentos (
    id serial PRIMARY KEY,
    conteudo TEXT,
    contexto VARCHAR(300),
    link VARCHAR(200),
    embedding vector(1024),
    criado_em TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP,
    atualizado_em TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP
); 

CREATE TABLE usuarios (
    id serial PRIMARY KEY,
    login VARCHAR(100) NOT NULL,
    nome VARCHAR(100) NOT NULL,
    senha VARCHAR(50)
);

CREATE TABLE mensagens (
    id serial PRIMARY KEY NOT NULL,
    id_chat uuid NOT NULL,
    conteudo TEXT NOT NULL,
    criado_em TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    assistente BOOLEAN NOT NULL
);

-- CreateTable
CREATE TABLE chats (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    id_usuario INTEGER NOT NULL,
    titulo  TEXT NOT NULL,
    criado_em TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    atualizado_em TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP
);

-- CreateIndex
CREATE UNIQUE INDEX "Mensagem_idChat_key" ON mensagens(id_chat);

-- AddForeignKey
ALTER TABLE mensagens ADD CONSTRAINT "Mensagem_idChat_fkey" FOREIGN KEY (id_chat) REFERENCES chats("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE chats ADD CONSTRAINT "Chat_idUsuario_fkey" FOREIGN KEY (id_usuario) REFERENCES usuarios("id") ON DELETE RESTRICT ON UPDATE CASCADE;


-- Função genérica para atualizar a coluna atualizadoEm
CREATE OR REPLACE FUNCTION fn_update_timestamp_generic()
RETURNS TRIGGER AS $$
BEGIN
    -- Verifica se a coluna atualizadoEm existe na tabela
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = TG_TABLE_SCHEMA
        AND table_name = TG_TABLE_NAME
        AND column_name = 'atualizado_em'
    ) THEN
        NEW.atualizado_em = CURRENT_TIMESTAMP;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER trg_atualizar_documentos
BEFORE UPDATE ON documentos
FOR EACH ROW
EXECUTE FUNCTION fn_update_timestamp_generic();

