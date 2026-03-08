from chat_analyzer import ChatAnalyzer
from connector import *



mychat = ChatAnalyzer("chat-1.json")
mychat.parse_conversation()
mychat.extract_document_type()
analyzed = mychat.analyze_all()
with get_connection() as conn:

    words = list(mychat.analyze()[0].vocabulary_.keys())
    # insert_concept(conn , words)
    document = mychat.get_parsed_doc()
   
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

    #chat_id, word, role, frequency, tf_idf_score
    rows = build_chat_concept_rows(document["uuid"] ,analyzed)
    entire_conversation = insert_chat_concepts(conn , rows )

    # human_prompts = insert_chat_concept()
    # assistant_prompts = insert_chat_concept()
