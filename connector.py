import psycopg2
from psycopg2.extras import execute_values
from contextlib import contextmanager
from dotenv import load_dotenv
import os



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