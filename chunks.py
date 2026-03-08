import json
from sklearn.feature_extraction.text import TfidfVectorizer


file_input = "chat-0.json"

count = 0

docs = []
# doc1 = {}

# doc2 = {}

# doc3 = {}

for i in range(3):

    file_input = f"chat-{i}.json"
    doc = {}
    with open(file_input,  "r" , encoding="utf-8") as f :
        data = json.load(f)
        print(f"on file {file_input}")
        for j in range(len(data)):

            print(len(data) ,"size of data")
            doc["name"] = data[j]["name"]
            snippet = {}
            
            for k in range(len(data[j]["chat_messages"])):
                snippet[data[j]["chat_messages"][k]["sender"]] = data[j]["chat_messages"][k]["text"]
                if k % 2 == 1:
                    snippet["created_at"]= data[j]["chat_messages"][k]["created_at"]
                    doc[k//2] = snippet
                    snippet = {}
              
                    
    docs.append(doc)

    #if i == 1:

        #print(doc)
        #print(data)
    # print( "this islength " , len(doc))
    # print("\n\n")
texts = []
humans = []
assistants = []
for i in range(len(docs)):
    
        for key , val in docs[i].items():
            #print(docs[i][key])
            if key != "name":
                text = ""
                human = ""
                assistant = ""
                for key2,val2 in val.items():
                    #if key2 == "human":
                    conversation+= docs[i][key][key2]
                    
                    if key2 == "human":
                        human +=  docs[i][key][key2]
                    
                    if key2 == "assistant":
                        assistant += docs[i][key][key2]
                texts.append(text)
                humans.append(human) 
                assistants.append(assistant)

        print("end of doc")
        print("\n\n")

#print(text)

tfidf =TfidfVectorizer()

conversation_result = tfidf.fit_transform(texts)

human_result =  tfidf.fit_transform(humans)

assistant_result =  tfidf.fit_transform(assistants)

print("len of texts" , len(texts))

print(texts)

print("\nidf values")
for el1 , el2 in zip(tfidf.get_feature_names_out() , tfidf.idf_):
    print(el1 , ":" , el2)

print("\n word indexes")
print(tfidf.vocabulary_)

print("\ntf-idf value")
print(conversation_result)