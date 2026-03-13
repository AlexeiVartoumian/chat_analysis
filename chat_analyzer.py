import json
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.feature_extraction.text import CountVectorizer
import re
import os
from dotenv import load_dotenv
from openai import OpenAI
from dataclasses import dataclass , field
from pathlib import Path
from typing import Optional

load_dotenv()

OPEN_API_KEY = os.environ.get("OPEN_API_KEY")
client = OpenAI(api_key=OPEN_API_KEY)

"""
problems with this code ..
1. too much going on 
1.1 init class has too many variables 
1.2 api key being used inside class prevents mock tests
1.3 too many methonds on class
1.4 functions calling other functions
1.5 json parsing does not handle errors
1.6 no use of Constants for entities such as human ,assitant
1.7 no use of dataclasses
"""

SENDER_HUMAN = "human"
SENDER_ASSISTANT = "assistant"
SENDER_FULL = "full"

SENDERS = (SENDER_HUMAN , SENDER_ASSISTANT , SENDER_FULL)
@dataclass
class Message:
    sender:str
    text:str 
    created_at: str
    uuid : str


@dataclass
class ConversationTexts:
    """ extracted text by json type """
    human : str = ""
    assistant : str = ""
    full : str = ""

@dataclass
class Codeblock:
    """extracted Code block is present from given Assistant response"""

    block_index : int
    language : Optional[str]
    content : str
    created_at : str
    message_uuid: str


@dataclass
class AnalysisResult:
    pass

@dataclass
class EmbeddedMessage:
    pass

@dataclass
class AnalysisResult:
    pass


def parse_conversation_file(path: str | Path  ) -> dict:
    """
    This function accepts a path to JSon file and loads it to create a file object.
    It then parses the data into the Message Dataclass
    """
    
    parsed_document : list[Message] = []

    path = Path(path)

    if not path.exists():
        raise FileNotFoundError(f"error finding the file")
    
    with path.open(encoding="utf-8") as file_handler:
        data = json.load(file_handler)
    
    if not isinstance(data , list) or not data:
        raise ValueError("Expected non empty array")
    
    conversation = data[0]

    messages = conversation.get("chat_messages")

    for  message in messages:
        m = Message()
        try:
            m["sender"] = message["sender"]
            m["text"] = message["text"]
            m["created_at"] = message["created_at"]
            m["uuid"] = message["uuid"]
        except KeyError as e: 
            raise ValueError(f"missing a required field {e}" ) from e
        
        parsed_document.append(m)
    
    return parsed_document

def extract_document_type( messages : list[Message] ) -> ConversationTexts:
        """parse the message list to extract text by sender"""

        parts = {
            SENDER_FULL : [],
            SENDER_HUMAN :[],
            SENDER_ASSISTANT: []
        }

        for message in messages:
            parts[SENDER_FULL].append(message.text) # full convo append whatever happens

            if message.sender in (SENDER_HUMAN , SENDER_ASSISTANT):
                parts[message.sender].append(message.text)

        return ConversationTexts(
            human  = " ".join(parts[SENDER_HUMAN]).strip(),
            assistant= " ".join(parts[SENDER_ASSISTANT]).strip(),
            full = " ".join(parts[SENDER_FULL]).strip()
        )    

def extract_code_blocks(messages: list[Message]) -> list[Codeblock] :
    """ Iterate through Messages . then grep the code from assistant response"""
    pattern = re.compile('```(\w+)?\n([\s\S]*?)```') # compile once use several times

    blocks : list[Codeblock] = []

    for message in messages:

        if message.sender != SENDER_ASSISTANT:
            continue
        for index , (language , content) in enumerate(pattern.findall(message.text)):
            blocks.append(
                Codeblock(
                    block_index= index,
                    language= language or None ,
                    content = content.strip(), 
                    created_at = message.created_at,
                    message_uuid=message.uuid
                )
            )
    
    return blocks

def strip_code_from_text(text: str) -> str:
    """ remove all code from assistant response to extract natural language response from assitant . 
    to be embedded andused for semantic search"""
    cleaned = re.sub(r'```[\s\S]*?```', '', text)
    cleaned = re.sub(r'`[^`]*`', '', cleaned)
    cleaned = ' '.join(cleaned.split())
    return cleaned.strip()

    
class ChatAnalyzer():


    """
    what should the class do? A v. basic nlp 
    accepts and analyzes a single conversation in json format .
    it parses the fields and will return three "documents" . 
    1. Conversation
    2. human
    3. Assistant. 

    it then parses those categories of documents and appliest the term-frequency 
    to extract the unique keywords with the highest scores in that chat .

    Could in future reference other chats.   
    """
    def __init__(self , path_to_chat  ):

        self.chat = path_to_chat
        self.parsed_document = {
        }
        #need separate vectorizer for each doctype since diff wordlength 
        self.vectorizers = {
        "full":  (TfidfVectorizer(stop_words='english'), CountVectorizer(stop_words='english')),
        "human":        (TfidfVectorizer(stop_words='english'), CountVectorizer(stop_words='english')),
        "assistant":     (TfidfVectorizer(stop_words='english'), CountVectorizer(stop_words='english'))
    }
        self.texts = {"human": [],
        "assistant": [],
        "full": []}
        self.results = {}

        self.assistant_code_blocks =[]
        self.embeddable_text= []

    """
    restructures the conversation into snippets where snippets consist two items.
    1.Human prompt . 2. Assistant response .   
    """
    def parse_conversation( self ) -> dict:
        
        parsed_document = {}
        with open(self.chat ,"r" , encoding = "utf-8") as f:
            data = json.load(f) # prob 1 what if dile not found ? 
            for j in range(len(data)): #prob 2 what is expected document structure / what is documentations? 
                parsed_document["name"] = data[j]["name"]
                parsed_document["uuid"] = data[j]["uuid"]
                parsed_document["created_at"] = data[j]["created_at"]
           
                for k , message in enumerate(data[j]["chat_messages"]):
                    parsed_document[k] = {
                        "sender": message["sender"],
                        "text" : message["text"],
                        "created_at": message["created_at"],
                        "uuid" : message["uuid"]
                    }
        self.parsed_document = parsed_document #parsed document dictionary could be a ordered list of class message 

        
        
        return self.parsed_document
    
    def get_parsed_doc(self):
        return self.parsed_document
    
    def extract_document_type(self):
        full_text = ""
        human = ""
        assistant = ""
        #print(self.parsed_document)
        for metadata , val in self.parsed_document.items():
            
            
            #chat snippets are stored by order they came in 
            if isinstance(metadata , int): 
                           
                text = val["text"]
                full_text += " " + text 

                if val["sender"] == "human":
                        human +=  " " + text
                    
                if val["sender"] == "assistant":
                    assistant += " " + text

                    #sneaking it in
                    self.extract_code_blocks(text ,  val["created_at"] , val["uuid"] )
                   
                    embedded_text = self.extract_embeddable_text(text)
                    
                    vector_embedding_api_call = client.embeddings.create(
                    input = embedded_text, 
                    model = "text-embedding-3-small"
                    )
                    self.embeddable_text.append ({
                        "message_uuid" : val["uuid"] ,
                        "embedding" : vector_embedding_api_call.data[0].embedding,
                        "embedded_text" : embedded_text ,
                        "created_at" : val["created_at"]
                    })

                    
         
        self.texts["human"].append( human.strip())
        self.texts["full"].append(full_text.strip())
        self.texts["assistant"].append(assistant.strip())
        

        #print(self.get_embeddable_text(assistant))

        

        return self.texts
    
    def analyze(self, doctype="full"):
        tfidf , count = self.vectorizers[doctype]
      
        tfidf_matrix = tfidf.fit_transform(self.texts[doctype])
        count_matrix = count.fit_transform(self.texts[doctype])
        self.results[doctype] = (tfidf_matrix, count_matrix)
        
        return tfidf, tfidf_matrix, count_matrix
  

    
    def analyze_all(self):

        return {

            doctype : self.analyze(doctype) for doctype in ["full" , "human" , "assistant"  ] 
                
        } 

    def extract_code_blocks(self , text ,created_at, uuid ):
        #only assitant code snippets for now will be called directly in the extract document type a little bit of coupling

        pattern = r'```(\w+)?\n([\s\S]*?)```'

        m = re.findall(pattern , text)


        if m :  
            for index , (language , content) in enumerate(m):
                self.assistant_code_blocks .append( {
                    "block_index" : index ,
                    "language": language or None, 
                    "content" : content.strip(),
                    "created_at": created_at,
                    "message_uuid": uuid
                })
        

    def get_code_blocks(self):
        return self.assistant_code_blocks
    
    def extract_embeddable_text(self, text):
        """
        may have to think about chunking this somehow where there is more than one concept in a chat . 
        """
        cleaned = re.sub(r'```[\s\S]*?```', '', text)
        cleaned = re.sub(r'`[^`]*`', '', cleaned)
        cleaned = ' '.join(cleaned.split())
        return cleaned.strip()

    def get_embeddable_text(self):
        return self.embeddable_text

mychat = ChatAnalyzer("chat-3.json")

mychat.parse_conversation()
mychat.extract_document_type()




