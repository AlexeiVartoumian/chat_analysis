
CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE message_embeddings(
    message_id UUID PRIMARY KEY REFERENCES messages(message_id) ON DELETE CASCADE,
    embedding VECTOR(1536),
    embedded_text TEXT NOT NULL ,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);


CREATE INDEX idx_embeddings_vector 
    ON message_embeddings USING hnsw (embedding vector_cosine_ops);