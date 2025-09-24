
CREATE TABLE documento (
    id VARCHAR(32) PRIMARY KEY,
    content TEXT,
    context VARCHAR(300),
    link VARCHAR(200),
    embedding vector(1024),
    created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP
); 