
CREATE TABLE documentos (
    id serial PRIMARY KEY,
    conteudo TEXT,
    contexto VARCHAR(300),
    link VARCHAR(200),
    embedding vector(1024),
    criado_em TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP,
    atualizado_em TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP
); 