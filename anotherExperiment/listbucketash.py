import boto3
import os
import json
from collections import defaultdict
from datetime import datetime
client = boto3.client('s3')

##check available files . pair off against local "seen" file . 
#write download obj

paginator = client.get_paginator('list_objects_v2')


resp = client.list_objects(Bucket='output-store-ash-cache-390746273208')

pages = paginator.paginate(Bucket='output-store-ash-cache-390746273208')


#TODO keys.txt should come from db

keys = []
timelines = []
new_keys = []

seen_keys = set()
new_keys_to_upload = []

date_format = '%Y-%m-%d'

with open("keys_ash.txt" , "r" , encoding="utf-8") as f:

    for line in f:
        #print(line)
        seen_keys.add(line.strip("\n"))

output_store_standard_timelines = set()
for page in pages :
    for obj in page.get('Contents' , []):
        key = obj['Key']
       
        if key.endswith(".csv"):
            timeline = "-".join( str(obj["LastModified"]).split(" ")).replace(":" , "-")
            
            check = key.removesuffix(".csv") +"_"+timeline+ ".csv"
            #print(check)        
            if check not in seen_keys:

                #datetime_key = datetime.strptime(obj["LastModified"] ,date_format)
                datestring= datetime.strftime(obj["LastModified"] ,date_format)
                year_month_date = datetime.strptime(datestring , date_format)
                
             
                if year_month_date not in timelines:
                    timelines.append(year_month_date)
                    output_store_standard_timelines.add(year_month_date)
                new_keys_to_upload.append(check)
                keys.append(key) ##will need to actually download


timelinecheck = defaultdict(set)
timelineindexcount = []


timelines = sorted(timelines)

for i in timelines:
    print(datetime.strftime(i ,date_format))
    
#timelines.append("-".join( str(key["LastModified"]).split(" ")).replace(":" , "-"))

with open("keys_ash.txt", "a" , encoding="utf-8") as f:
        
        for index , key  in enumerate(new_keys_to_upload):
            f.writelines(key)
            f.writelines("\n")

with open("newkeys_ash.txt" , "w" , encoding="utf-8") as f :

        for index , key in enumerate(new_keys_to_upload):
            f.writelines(key)
            f.writelines("\n")


def forgive(key , mydict ):

    if key not in mydict:
        mydict[key] =  [0] *4    
    return mydict

def dblforgive(key , mydict ):

    
    if key not in mydict:
        mydict[key] =  [0]     
    return mydict

with open("keys_ash.json" , "w" , encoding="utf-8" ) as f:
    output = defaultdict(list)
   
    ##need this since this guarantee the order
    for index , timeline in enumerate(timelines):
       
        standard_timeline =  datetime.strftime(timeline ,date_format)

        #it could be a standard job does not happen on a given day but a backfill one does
        if timeline in output_store_standard_timelines:
            output[standard_timeline] = [{}]

            print("standard", standard_timeline,output[standard_timeline])
            
       
    
    for index , key  in enumerate(new_keys_to_upload):
        #we extract matching date from filename and use that
        unique = key.split("-", 1)[1].split(".")[0].split("_")[0] #wtf mate regretting my life choices
        # print("unique parsing \n")
        # print(key)
        # print(unique)
        # print("---------------------\n")
        document = os.path.basename(key)
        
        #TODO POTENTIAL FAILURE ON FILE NAME PARSING

        #string parse fill date string to y:m:d to be used be parsed back to string // 
        timeline_key = document.split("_")[-1].removesuffix(".csv")
        # print(document)
        # print(timeline_key)
        #2026-04-24 20:06:52+00:00
        #2026-04-24-21-02-43+00-00 [:10]
        #TODO POTENTIAL FAILURE ON FILE NAME PARSING
        timeline_key = datetime.strptime(timeline_key[:10] , date_format)
        timeline_key = datetime.strftime(timeline_key ,date_format)

        new_keys.append(document) 
        if document.startswith("processedJobsAsh"):
            
            #output[timeline_key][0][records[unique][0]] = document
            forgive(unique ,output[timeline_key][0])
            output[timeline_key][0][unique][0] = document
        
        
        if document.startswith("JobdescriptionAsh"):
            #output[timeline_key][0][records[unique][0]] = document
            forgive(unique ,output[timeline_key][0])
            output[timeline_key][0][unique][1] = document

        if document.startswith("deadlink"):
            forgive(unique ,output[timeline_key][0])
            output[timeline_key][0][unique][2] = document
        
        if document.startswith("AshJobsByCompany"):
            forgive(unique , output[timeline_key][0])
            output[timeline_key][0][unique][3] = document
      
       
    json.dump(output, f)
  
for index ,key in enumerate(keys):
    #print(key)
    sanitizekey = new_keys[index]
    # print(sanitizekey)
    # print("here we go \n")
    with open (sanitizekey , "wb" ) as f :
        client.download_fileobj('output-store-ash-cache-390746273208', key, f)



