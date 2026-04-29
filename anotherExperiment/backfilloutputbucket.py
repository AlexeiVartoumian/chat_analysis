import boto3
import os
import json
from collections import defaultdict
import csv
client = boto3.client('s3')


"""
TODO : YOU MESSED IT! problem is upstream filetype is called live . problem is insertion logic wants to 
scan ALL ROLES first suspended or no . then it wants scan closed state roles. you have two diff files backfilloutputbucket and listbucketv2 doing thier own thing
also improve the pass down filename logic from container.

flow is -> 1.live_roles.csv 2. closed/expired roles.csv
"""

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
    if key.endswith(".json"):
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
    