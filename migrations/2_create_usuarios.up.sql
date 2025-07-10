
CREATE TABLE usuarios (
    id serial PRIMARY KEY,
    login VARCHAR(100) NOT NULL,
    nome VARCHAR(100) NOT NULL,
    senha VARCHAR(50)
);