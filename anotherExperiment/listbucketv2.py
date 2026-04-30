import boto3
import os
import json
from collections import defaultdict
from datetime import datetime
client = boto3.client('s3')

##check available files . pair off against local "seen" file . 
#write download obj

paginator = client.get_paginator('list_objects_v2')


resp = client.list_objects(Bucket='output-store-390746273208')

pages = paginator.paginate(Bucket='output-store-390746273208')

pages2 = paginator.paginate(Bucket='backfill-store-390746273208')

#TODO keys.txt should come from db

keys = []
timelines = []
new_keys = []

seen_keys = set()
new_keys_to_upload = []

date_format = '%Y-%m-%d'

with open("keys.txt" , "r" , encoding="utf-8") as f:

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
for page in pages2 :
    for obj in page.get('Contents' , []):
        key = obj['Key']
       
        if key.endswith(".csv"):
            timeline = "-".join( str(obj["LastModified"]).split(" ")).replace(":" , "-")
            check = key.removesuffix(".csv") +"_"+timeline+ ".csv"               
            if check not in seen_keys:

                datestring= datetime.strftime(obj["LastModified"] ,date_format)
                year_month_date = datetime.strptime(datestring , date_format)
                base = os.path.basename(check)
                if year_month_date not in timelines:
                    timelines.append(year_month_date)

                #a large scan should only happen once a day . otherwise this code block will break
                if base.startswith("live"):
                    timelinecheck[year_month_date].add("l")
                    timelineindexcount.append(year_month_date)

                if base.startswith("suspended"):
                    timelinecheck[year_month_date].add("s")
                    timelineindexcount.append(year_month_date)

                if base.startswith("expired"):
                    timelinecheck[year_month_date].add("e")
                    timelineindexcount.append(year_month_date)
                
                new_keys_to_upload.append(check)
                keys.append(key) ##will need to actually download




timelines = sorted(timelines)
backfill_timelineindexes= set()
for index, time in enumerate(timelines):
    if time in timelinecheck:
        backfill_timelineindexes.add(index)

for i in timelines:
    print(datetime.strftime(i ,date_format))
    
#timelines.append("-".join( str(key["LastModified"]).split(" ")).replace(":" , "-"))

with open("keys.txt", "a" , encoding="utf-8") as f:
        
        for index , key  in enumerate(new_keys_to_upload):
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

with open("keys.json" , "w" , encoding="utf-8" ) as f:
    output = defaultdict(list)
   
    ##need this since this guarantee the order
    for index , timeline in enumerate(timelines):
       
        standard_timeline =  datetime.strftime(timeline ,date_format)

        #it could be a standard job does not happen on a given day but a backfill one does
        if timeline in output_store_standard_timelines:
            output[standard_timeline] = [{}]

            print("standard", standard_timeline,output[standard_timeline])
            
        if index in backfill_timelineindexes:

            if "l" in timelinecheck[timeline]:
                
                back_fill_timeline =  datetime.strftime(timeline ,date_format)
                output[f"{back_fill_timeline}_liveRoles"] = [{}]
           
            #order of ops expiry workflow comes before suspend workflow atm can guarantee they appear once
            if "e" in timelinecheck[timeline]:
                 back_fill_timeline =  datetime.strftime(timeline ,date_format)
                 output[f"{back_fill_timeline}_expiredRoles"] = [{}]

            if "s" in timelinecheck[timeline]:
                back_fill_timeline =  datetime.strftime(timeline ,date_format)
                output[f"{back_fill_timeline}_suspendedRoles"] = [{}]
                print(f"{back_fill_timeline}_suspendedRoles" , "yup" , output[f"{back_fill_timeline}_suspendedRoles"])
    
    for index , key  in enumerate(new_keys_to_upload):
        #we extract matching date from filename and use that
        unique = key.split("-", 1)[1].split(".")[0].split("_")[0] #wtf mate regretting my life choices
        # print("unique parsing \n")
        # print(key)
        # print(unique)
        # print("---------------------\n")
        document = os.path.basename(key)
        
        #string parse fill date string to y:m:d to be used be parsed back to string
        timeline_key = document.split("_")[-1].removesuffix(".csv")
        # print(document)
        # print(timeline_key)
        #2026-04-24 20:06:52+00:00
        #2026-04-24-21-02-43+00-00 [:10]
        timeline_key = datetime.strptime(timeline_key[:10] , date_format)
        timeline_key = datetime.strftime(timeline_key ,date_format)

        new_keys.append(document) 
        if document.startswith("processed"):
            
            #output[timeline_key][0][records[unique][0]] = document
            forgive(unique ,output[timeline_key][0])
            output[timeline_key][0][unique][0] = document
        if document.startswith("company_data"):
            #output[timeline_key][0][records[unique][1]] = document
            forgive(unique ,output[timeline_key][0])
            output[timeline_key][0][unique][1] = document
        if document.startswith("job_metadata"):
            #output[timeline_key][0][records[unique][2]] = document
            forgive(unique ,output[timeline_key][0])
            output[timeline_key][0][unique][2] = document
        if document.startswith("job_description"):
            #output[timeline_key][0][records[unique][3]] = document
            forgive(unique ,output[timeline_key][0])
            output[timeline_key][0][unique][3] = document

        if document.startswith("live"):
            
        
            dblforgive(unique ,output[f"{timeline_key}_liveRoles"][0])
            output[f"{timeline_key}_liveRoles"][0][unique][0] = document
        
        if document.startswith("expired"):
            
        
            dblforgive(unique ,output[f"{timeline_key}_expiredRoles"][0])
            output[f"{timeline_key}_expiredRoles"][0][unique][0] = document
        
        if document.startswith("suspended"):
           
            dblforgive(unique ,output[f"{timeline_key}_suspendedRoles"][0])
            output[f"{timeline_key}_suspendedRoles"][0][unique][0] = document
    json.dump(output, f)
  

for index ,key in enumerate(keys):
    #print(key)
    sanitizekey = new_keys[index]
    # print(sanitizekey)
    # print("here we go \n")
    if key.startswith("suspended") or key.startswith("expired")  or key.startswith("live") :
        
        with open (sanitizekey , "wb" ) as f :
            client.download_fileobj('backfill-store-390746273208', key, f)
        continue
    with open (sanitizekey , "wb" ) as f :
        client.download_fileobj('output-store-390746273208', key, f)




# with open ("processedJobs.csv" , "wb" ) as f :
#     client.download_fileobj('alexeitranscribefile', 'input/processedJobs.csv', f)

# with open("company_data.csv" , "wb" ) as f :
#     client.download_fileobj('alexeitranscribefile', '/output/company_data.csv', f)

# with open("job_metadata.csv" , "wb" ) as f :
#     client.download_fileobj('alexeitranscribefile', '/output/job_metadata.csv', f)

# with open("job_description.csv" , "wb" ) as f :
#     client.download_fileobj('alexeitranscribefile', '/output/job_description.csv', f)