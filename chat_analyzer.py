import json
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.feature_extraction.text import CountVectorizer
import typing 



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
    def __init__(self , path_to_chat ):

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
        

    """
    restructures the conversation into snippets where snippets consist two items.
    1.Human prompt . 2. Assistant response .   
    """
    def parse_conversation( self ) -> dict:
        
        parsed_document = {}
        with open(self.chat ,"r" , encoding = "utf-8") as f:
            data = json.load(f)
            for j in range(len(data)):
                parsed_document["name"] = data[j]["name"]
                parsed_document["uuid"] = data[j]["uuid"]
                parsed_document["created_at"] = data[j]["created_at"]
                # snippet = {}                
                # for k in range(len(data[j]["chat_messages"])):
                #     snippet[data[j]["chat_messages"][k]["sender"]] = data[j]["chat_messages"][k]["text"]
                #     if k % 2 == 1:
                #         snippet["created_at"]= data[j]["chat_messages"][k]["created_at"]
                #         parsed_document[k//2] = snippet
                #         snippet = {}

                for k , message in enumerate(data[j]["chat_messages"]):
                    parsed_document[k] = {
                        "sender": message["sender"],
                        "text" : message["text"],
                        "created_at": message["created_at"],
                        "uuid" : message["uuid"]
                    }
        self.parsed_document = parsed_document

        
        return self.parsed_document
    
    def get_parsed_doc(self):
        return self.parsed_document
    
    def extract_document_type(self):
        full_text = ""
        human = ""
        assistant = ""
        for metadata , val in self.parsed_document.items():
            

            #chat snippets are stored by order they came in 
            if isinstance(metadata , int): 
                           
                text = val["text"]
                full_text += " " + text 

                if val["sender"] == "human":
                        human +=  " " + text
                    
                if val["sender"] == "assistant":
                    assistant += " " + text
              
        self.texts["human"].append( human.strip())
        self.texts["full"].append(full_text.strip())
        self.texts["assistant"].append(assistant.strip())
        
        return self.texts
    
    def analyze(self, doctype="full"):
        tfidf , count = self.vectorizers[doctype]
      
        tfidf_matrix = tfidf.fit_transform(self.texts[doctype])
        count_matrix = count.fit_transform(self.texts[doctype])
        self.results[doctype] = (tfidf_matrix, count_matrix)
        
        return tfidf, tfidf_matrix, count_matrix
  
    # def analyze(self , doctype="conversation"):
        
        

    #     self.results[doctype] = self.vectorizer.fit_transform(self.texts[doctype])
        
    #     return self.vectorizer , self.results[doctype]
    
    def analyze_all(self):

        return {

            doctype : self.analyze(doctype) for doctype in ["full" , "human" , "assistant"  ] 
                
        } 


# mychat = ChatAnalyzer("chat-1.json")

# mychat.parse_conversation()
# mychat.extract_document_type()

# print(mychat.get_parsed_doc())

# tfidresp , results , count = mychat.analyze()

# print( list(mychat.analyze()[0].vocabulary_.keys() ), mychat.analyze()[1] )



