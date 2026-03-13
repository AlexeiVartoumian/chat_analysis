from chat_analyzer import ChatAnalyzer
from connector import *
import os
from pathlib import Path
from openai import OpenAI
import logging
load_dotenv()

OPEN_API_KEY = os.environ.get("OPEN_API_KEY")
embedding_client = OpenAI(api_key=OPEN_API_KEY)


logging.basicConfig(
    level = logging.INFO,
    format="%(asctime)s  %(levelname)-8s  %(message)s",
)
logger = logging.getLogger(__name__)

cur_path = os.getcwd()

folder_directory = os.path.join(cur_path , "chats")

path = Path.glob(Path(folder_directory) , "*.json")

for file in path :
    print(file)
    print(file, type(file))
    mychat = ChatAnalyzer(file , embedding_client=embedding_client)

    conversation, messages = mychat.parse_conversation()
    texts = mychat.extract_texts(messages)

    if not texts.full.strip() or not texts.assistant.strip() or not texts.human.strip():
        logger.warning("Skipping %s — no text content", file.name)
        continue
 

    analyzed = mychat.analyze_all(texts)
    code_blocks = mychat.extract_code_blocks(messages)
    message_embeddings = mychat.extract_embeddings(messages)
    with get_connection() as conn:


        insert_chat(conn, conversation)
        insert_message(conn, conversation, messages)
        insert_code_blocks(conn, code_blocks)
        insert_message_embeddings(conn, message_embeddings)
 
        rows = build_chat_concept_rows(conversation.uuid, analyzed)
        insert_chat_concepts(conn, rows)