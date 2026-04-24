import boto3
import os
import json
from collections import defaultdict
client = boto3.client('s3')




resp = client.list_objects(Bucket='backfill-store-390746273208')

keys = []
timelines = []

for key in resp['Contents']:
    keys.append(key["Key"])
    timelines.append("-".join( str(key["LastModified"]).split(" ")).replace(":" , "-"))


        
for index , key  in enumerate(keys):
    timestamp = timelines[index] 
    print(key)
    #split it up add timestamp and  .csv in
    newkey = key.split(".")[0] + "_" + timestamp + ".csv"
    with open (newkey , "wb" ) as f :
        client.download_fileobj('backfill-store-390746273208', key, f)
    