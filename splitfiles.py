import os
import json


input_file = "conversations.json"




path = os.getcwd()

outputpath = os.path.join(path , "chats")
myfileSize = os.path.join(path , input_file)

print(os.path.getsize(myfileSize))
  
with open(input_file,  "rb" ) as f :
        # num_lines = sum(1 for _  in input)
        # print(num_lines)
        data = json.load(f)
        print(len(data[0]))
        print(len(data))

        
        for i in range(len(data)):
            outputfile = f"chat-{i}.json"

            outputfilepath = os.path.join(outputpath , outputfile)
            jsonto_append = []
            mydict = {}
            for key , val in data[i].items() :

                # print(key , type(val))
                # print(key , val)
                mydict[key] = val
            
            jsonto_append.append(mydict)

            with open(outputfilepath , "w") as f:
                 
                 json.dump(jsonto_append , f)
                
              
            print("new line \n\n\n")
              



# output = "chats"

# if not os.path.exists(output):
#     os.makedirs(output)

# def write_chunk(chunk_lines , output_file_path):
#     with open(output_file_path, "wb" ) as output_file:

#         output_file.writelines(chunk_lines)