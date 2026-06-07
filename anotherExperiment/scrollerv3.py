import sys 
import json
import boto3
import time
import uuid
from boto3.dynamodb.conditions import Key, Attr
from botocore.exceptions import ClientError



search_term = "cloud engineer"
search_term_id = -1
company_c = None
s3 = boto3.client("s3", region_name='eu-west-2')


workflow_id = str(uuid.uuid4())

if len(sys.argv) >1:

    search_term = sys.argv[1]
    search_term_id = sys.argv[2]

else:
    company_c = sys.stdin.read()
    s3.put_object(Bucket = "alexeitranscribefile" , Key=f"fresh_companyInd-{workflow_id}.json" , Body=company_c)


def main():
    sqs = boto3.client('sqs', region_name='eu-west-2')

    if company_c is not None:
        response = sqs.send_message(
        QueueUrl='https://sqs.eu-west-2.amazonaws.com/390746273208/workflow-cordinator-deed',
        MessageBody=json.dumps({
            "search_term_id": None,
            "search_term": None,
            "workflow_type": "company",
            "workflow_id": workflow_id
        })
        )
        print(search_term_id)
        print(search_term)
        print("company is workflowtype")
        print(type(search_term_id))
        print(type(search_term))

    else:
        response = sqs.send_message(
        QueueUrl='https://sqs.eu-west-2.amazonaws.com/390746273208/workflow-cordinator-deed',
        MessageBody=json.dumps({
            "search_term_id": search_term_id,
            "search_term": search_term,
            "workflow_type": "search_term",
            "workflow_id": workflow_id
        })
        )
        print(search_term_id)
        print(search_term)
        print("company is workflowtype")
        print(type(search_term_id))
        print(type(search_term))
   
    print("response:", response)
    

if __name__ == "__main__":
    main()