import boto3
import os
import json
from collections import defaultdict
client = boto3.client('s3')

##check available files . pair off against local "seen" file . 
#write download obj


resp = client.list_objects(Bucket='alexeitranscribefile')

keys = []
for key in resp["Contents"]:
    keys.append(key["Key"])
#print(resp)

#print(keys)

with open("keys.txt", "w" , encoding="utf-8") as f:
        
        for key in keys:
            f.writelines(key)
            f.writelines("\n")

with open("keys.json" , "w" , encoding="utf-8" ) as f:
    records = defaultdict(list)

    for key in keys:
        unique = key.split("-", 1)[1].split(".")[0] 
        print(unique)
        
        records[unique].append(os.path.basename(key))
       
    json.dump(records , f)


# for key in keys:
#     sanitizekey = os.path.basename(key)
#     with open (sanitizekey , "wb" ) as f :
#         client.download_fileobj('alexeitranscribefile', key, f)

#     with open(sanitizekey , "wb" ) as f :
#         client.download_fileobj('alexeitranscribefile', key, f)

#     with open(sanitizekey , "wb" ) as f :
#         client.download_fileobj('alexeitranscribefile', key, f)

#     with open(sanitizekey, "wb" ) as f :
#         client.download_fileobj('alexeitranscribefile', key, f)


# with open ("processedJobs.csv" , "wb" ) as f :
#     client.download_fileobj('alexeitranscribefile', 'input/processedJobs.csv', f)

# with open("company_data.csv" , "wb" ) as f :
#     client.download_fileobj('alexeitranscribefile', '/output/company_data.csv', f)

# with open("job_metadata.csv" , "wb" ) as f :
#     client.download_fileobj('alexeitranscribefile', '/output/job_metadata.csv', f)

# with open("job_description.csv" , "wb" ) as f :
#     client.download_fileobj('alexeitranscribefile', '/output/job_description.csv', f)


