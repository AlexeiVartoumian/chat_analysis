from pgvector.psycopg2 import register_vector
import psycopg2
from psycopg2.extras import execute_values

from contextlib import contextmanager
from dotenv import load_dotenv
import os
import numpy as np


load_dotenv()


CONNECTION_STRING = os.environ.get("CONNECTION_STRING")




@contextmanager
def get_connection():
    conn = psycopg2.connect(CONNECTION_STRING)
    try:
        yield conn
        conn.commit()
    except Exception:
        conn.rollback()
        raise
    finally:
        conn.close()

def insert_chat(conn , chat_id , title , created_at):

    with conn.cursor() as cur:

        cur.execute(""" 
                    INSERT INTO chats (chat_id,title, created_at)
                    VALUES (%s, %s , %s)
                    ON CONFLICT (chat_id) DO NOTHING
                    """, (chat_id, title , created_at))

def insert_message(conn, messages):
    
    with conn.cursor() as cur:
        
        execute_values(cur,""" 
            INSERT INTO messages (message_id, chat_id, role, content, created_at)
            VALUES %s
            ON CONFLICT (message_id) DO NOTHING
        """, [ (m["message_id"], m["chat_id"], m["role"], m["content"], m["created_at"])
               for m in messages ] ) 

# code_block_index int NOT NULL,
#     message_id UUID NOT NULL REFERENCES messages(message_id) ON DELETE CASCADE ,
#     language TEXT,
#     content TEXT NOT NULL,
#     created_at TIMESTAMP NOT NULL DEFAULT NOW(),
#     PRIMARY KEY(code_block_index , message_id)
def insert_code_blocks(conn , code_blocks):

    with conn.cursor() as cur:
        
        execute_values(cur , """
        INSERT INTO code_blocks (code_block_index , message_id, language , content , created_at  )
        VALUES %s
        ON CONFLICT (code_block_index , message_id) DO NOTHING
""" , [ (code_block["block_index"],code_block["message_uuid"] ,code_block["language"] ,code_block["content"] ,code_block["created_at"] ) 
       for code_block in code_blocks  ])

def insert_message_embeddings(conn , message_embeddings):

    with conn.cursor() as cur:
        
        execute_values(cur , """
        INSERT INTO message_embeddings ( message_id, embedding , embedded_text , created_at  )
        VALUES %s
        ON CONFLICT (message_id) DO NOTHING
""" , [ ( message["message_uuid"] ,message["embedding"] ,message["embedded_text"] ,message["created_at"] ) 
       for message in message_embeddings  ])
        
def insert_chat_concepts(conn, rows):
    with conn.cursor() as cur:
        execute_values(cur, """
            INSERT INTO chat_concepts (chat_id, word, role, frequency, tf_idf_score)
            VALUES %s
            ON CONFLICT (chat_id, word, role) DO UPDATE
                SET frequency = EXCLUDED.frequency,
                    tf_idf_score = EXCLUDED.tf_idf_score
        """, rows)

def build_chat_concept_rows(chat_id, analyzed):
    rows = []
    print(analyzed , type(analyzed))
    
    for role, (vectorizer, tfidf_matrix, count_matrix) in analyzed.items():

        # print(len(vectorizer.vocabulary_.items()))
        # print("haha",  tfidf_matrix[0 , 0])
        # print(count_matrix.shape[0])
        for word, idx in vectorizer.vocabulary_.items():
           
            tf_idf_score = float(tfidf_matrix[0, idx])
            frequency = int(count_matrix[0, idx])
            if tf_idf_score > 0:
                rows.append((chat_id, word, role, frequency, tf_idf_score))
        print("iteration complete \n\n\n")
    return rows
# def insert_message(conn, message_id, chat_id, role, content, created_at):
    
#     with conn.cursor() as cur:
        
#         cur.execute("""
#             INSERT INTO messages (message_id, chat_id, role, content, created_at)
#             VALUES (%s, %s, %s, %s, %s)
#             ON CONFLICT (message_id) DO NOTHING
#         """, (message_id, chat_id, role, content, created_at))

# def insert_concept(conn, words):
    
#     with conn.cursor() as cur:
        
#         execute_values(cur ,"""
#             INSERT INTO concepts (word)
#             VALUES %s
#             ON CONFLICT (word) DO NOTHING
#         """, [(word,) for word in words])

# def insert_chat_concept(conn, chat_id, word, role, frequency, tf_idf_score):
    
#     with conn.cursor() as cur:
        
#         cur.execute("""
#             INSERT INTO chat_concepts (chat_id, word, role, frequency, tf_idf_score)
#             VALUES (%s, %s, %s, %s, %s)
#             ON CONFLICT (chat_id, word, role) DO UPDATE
#                 SET frequency = EXCLUDED.frequency,
#                     tf_idf_score = EXCLUDED.tf_idf_score
#         """, (chat_id, word, role, frequency, tf_idf_score))



def get_message(conn ):

    with conn.cursor() as cur:
        
        cur.execute("""
        SELECT * FROM message_embeddings
""")
        
        for record in cur:
            print(record)
        

def get_embeddings(conn , vector):

    with conn.cursor() as cur:
        register_vector(conn)
        vector =np.array(vector)
        cur.execute("""
        SELECT 1 - (embedding <=> %s ) AS cosine_similarity , embedded_text FROM message_embeddings ORDER BY embedding <=> %s LIMIT 5
""", (vector,vector))
        
        for record in cur :
            print(record)

# with get_connection() as conn:
#     get_message(conn)