import boto3
import os
import json
from collections import defaultdict
import csv
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
    if key.startswith("live-roles"):
        with open (key , "wb" ) as f :
            client.download_fileobj('backfill-store-390746273208', key, f)

        
        with open(key , "r" , encoding="utf-8") as jsonfile:

            data = json.load(jsonfile)

            with open(newkey , "w" , newline='' , encoding="utf-8") as csv_file:    
                
                fieldnames = ["job_id"]
                writer = csv.DictWriter(csv_file , fieldnames=fieldnames)
                writer.writeheader()    
                for item in data:
                    writer.writerow({"job_id" : item})


    else:
        with open (newkey , "wb" ) as f :
            client.download_fileobj('backfill-store-390746273208', key, f)
    