
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
