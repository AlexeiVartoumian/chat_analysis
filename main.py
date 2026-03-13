from chat_analyzer import ChatAnalyzer
from connector import *
import os
from pathlib import Path


cur_path = os.getcwd()

folder_directory = os.path.join(cur_path , "chats")

path = Path.glob(Path(folder_directory) , "*.json")
for file in path :
    print(file)

    mychat = ChatAnalyzer(file)
    mychat.parse_conversation()

    #chats are stored by ordered number
    if 0 not in mychat.parse_conversation(): 
        continue

    document = mychat.get_parsed_doc()

    texts = mychat.extract_document_type()
    #chats could contain empty records for some reason ?
    if len(texts["human"][0]) == 0 or len(texts["assistant"][0]) == 0:
        print(texts["human"])
        continue

    analyzed = mychat.analyze_all()
    code_blocks = mychat.get_code_blocks()
    message_embeddings = mychat.get_embeddable_text()
    with get_connection() as conn:

        words = list(mychat.analyze()[0].vocabulary_.keys())
        # insert_concept(conn , words)
        
    
        messages = [
        { "message_id" :  val["uuid"] ,
            "chat_id" : document["uuid"],
            "role" : val["sender"] ,
            "content": val["text"],
            "created_at": val["created_at"]
            } 
        for key , val in document.items()
        if isinstance(key , int)
        ]
        print(document["uuid"] , document["name"] , document["created_at"])
        insert_chat(conn , document["uuid"] , document["name"] , document["created_at"])
        insert_message(conn , messages)
        insert_code_blocks(conn , code_blocks)
        insert_message_embeddings(conn , message_embeddings)
        #chat_id, word, role, frequency, tf_idf_score
        rows = build_chat_concept_rows(document["uuid"] ,analyzed)
        entire_conversation = insert_chat_concepts(conn , rows )

        # human_prompts = insert_chat_concept()
        # assistant_prompts = insert_chat_concept()
