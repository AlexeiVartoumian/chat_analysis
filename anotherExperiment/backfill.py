import sys 
import json
import boto3
import time
import uuid
from boto3.dynamodb.conditions import Key, Attr
from botocore.exceptions import ClientError


"""
here we do this diretly invokde lamdba reading stdin the json . then 
"""
#search type is suspended or live
search_type = None
auto = None
if len(sys.argv) > 1:
    search_type = sys.argv[1]
    auto = sys.argv[2]

roles = json.loads(sys.stdin.read())

workflow_id = str(uuid.uuid4())

s3 = boto3.client("s3", region_name='eu-west-2')

lamdba_client = boto3.client("lambda" , region_name="eu-west-2")

source_store = "source-store-390746273208"
output_store = "backfill-store-390746273208"



s3.put_object(
        Bucket=output_store,
        Key=f"{search_type}-{workflow_id}-roles.json",
        Body=json.dumps(roles, indent=2),
        ContentType="application/json",
    )

def main():
    

    payload = json.dumps({"s3_source_bucket": source_store ,"output_store": output_store, 
                          "search_type" : search_type , "workflow_id" : workflow_id ,"auto" :auto })
    response = lamdba_client.invoke(
    FunctionName='backfill_orchestrator',
    Payload=payload,
    Qualifier='1',
    )

    print(response)

if __name__ == "__main__":
    main()