import boto3
import os
import json
from collections import defaultdict
client = boto3.client('s3')

##check available files . pair off against local "seen" file . 
#write download obj


resp = client.list_objects(Bucket='output-store-390746273208')

keys = []
timelines = []
for key in resp["Contents"]:
    keys.append(key["Key"])
    
    timelines.append("-".join( str(key["LastModified"]).split(" ")).replace(":" , "-"))

with open("keys.txt", "w" , encoding="utf-8") as f:
        
        for index , key  in enumerate(keys):
            timestamp = timelines[index] 
            
            #split it up add timestamp and  .csv in
            key = key.split(".")[0] + "_" + timestamp + ".csv"
            f.writelines(key)
            f.writelines("\n")

with open("keys.json" , "w" , encoding="utf-8" ) as f:
    records = defaultdict(lambda: [0] * 4)
    for index , key  in enumerate(keys):
        unique = key.split("-", 1)[1].split(".")[0] 
        document = os.path.basename(key)
        timestamp = timelines[index]
        key = key.split(".")[0] + "_" + timestamp + ".csv" 
        if document.startswith("processed"):
            records[unique][0] = os.path.basename(key)
        if document.startswith("company_data"):
            records[unique][1] = os.path.basename(key)
        if document.startswith("job_metadata"):
            records[unique][2] = os.path.basename(key)
        if document.startswith("job_description"):
            records[unique][3] = os.path.basename(key)
  
    json.dump(records , f)

count = 0
for key in keys:
    sanitizekey = os.path.basename(key)


    with open (sanitizekey , "wb" ) as f :
        client.download_fileobj('output-store-390746273208', key, f)

    with open(sanitizekey , "wb" ) as f :
        client.download_fileobj('output-store-390746273208', key, f)

    with open(sanitizekey , "wb" ) as f :
        client.download_fileobj('output-store-390746273208', key, f)

    with open(sanitizekey, "wb" ) as f :
        client.download_fileobj('output-store-390746273208', key, f)


# with open ("processedJobs.csv" , "wb" ) as f :
#     client.download_fileobj('alexeitranscribefile', 'input/processedJobs.csv', f)

# with open("company_data.csv" , "wb" ) as f :
#     client.download_fileobj('alexeitranscribefile', '/output/company_data.csv', f)

# with open("job_metadata.csv" , "wb" ) as f :
#     client.download_fileobj('alexeitranscribefile', '/output/job_metadata.csv', f)

# with open("job_description.csv" , "wb" ) as f :
#     client.download_fileobj('alexeitranscribefile', '/output/job_description.csv', f)