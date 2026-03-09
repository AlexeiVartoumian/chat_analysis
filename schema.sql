
-- CREATE EXTENSION IF NOT EXISTS vector;


CREATE TABLE chats(
    chat_id UUID  PRIMARY KEY,
    title   TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TYPE message_role AS ENUM('human' , 'assistant');

CREATE TABLE messages(
    message_id UUID PRIMARY KEY,
    chat_id UUID NOT NULL REFERENCES chats(chat_id) ON DELETE CASCADE,
    role message_role NOT NULL, 
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_messages_chat_id ON messages(chat_id);


CREATE TYPE analysis_role AS ENUM('human' , 'assistant' , 'full'); --refactor

CREATE TABLE chat_concepts(
    chat_id UUID NOT NULL REFERENCES chats(chat_id) ON DELETE CASCADE,
    word TEXT NOT NULL, --REFERENCES concepts(word) ON DELETE CASCADE,
    role analysis_role NOT NULL,
    frequency bigint NOT NULL ,
    tf_idf_score FLOAT NOT NULL ,
    PRIMARY KEY(chat_id , word , role) -- to uniquely identify row
);


CREATE TABLE code_blocks(
    code_block_index int NOT NULL,
    message_id UUID NOT NULL REFERENCES messages(message_id) ON DELETE CASCADE ,
    language TEXT,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY(code_block_index , message_id)
)

CREATE INDEX idx_chat_concepts_word ON chat_concepts(word);
CREATE INDEX idx_chat_concepts_chat_id ON chat_concepts(chat_id);
CREATE INDEX idx_chat_concepts_chat_role ON chat_concepts(chat_id, role);


-- CREATE TABLE embeddings(
--     message_id UUID PRIMARY KEY REFERENCES messages(message_id) ON DELETE CASCADE,
--     embedding VECTOR(1536)
-- );


-- CREATE INDEX idx_embeddings_vector 
--     ON embeddings USING hnsw (vector vector_cosine_ops);


