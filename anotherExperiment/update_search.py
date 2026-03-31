import boto3
import os
import json
import sys

search_term = None

if len(sys.argv) >1 :

    search_term = sys.argv[1]


client = boto3.client('lambda', region_name='eu-west-2')

#TODO paramertirize
if search_term != None:
    response = client.invoke(
        FunctionName='reader',
        InvocationType='Event', 
        Payload=json.dumps({
            "search_term": search_term,
            "target_bucket" : "somebuckethaha" 
        })
    )
else:
    raise Exception(search_term , "whoopsie")