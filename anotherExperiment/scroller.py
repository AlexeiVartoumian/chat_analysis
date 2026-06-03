import sys 
import json
import boto3
import time
import uuid
from boto3.dynamodb.conditions import Key, Attr
from botocore.exceptions import ClientError



search_term = "cloud engineer"
company_c = None
s3 = boto3.client("s3", region_name='eu-west-2')
S3_BUCKET = "somebuckethaha"

workflow_id = str(uuid.uuid4())

if len(sys.argv) >1:

    search_term = sys.argv[1]
    search_term = search_term.lstrip('{"')
    search_term = search_term.rstrip('}"')
else:
    company_c = sys.stdin.read()
    s3.put_object(Bucket = "alexeitranscribefile" , Key=f"fresh_companyInd-{workflow_id}.json" , Body=company_c)


def main():
    
    if company_c is not None:
        ecs = boto3.client('ecs' , region_name = "eu-west-2")
        response = ecs.run_task(
        cluster='scroller-cluster',
        taskDefinition='scroller-task',
        launchType='FARGATE',
        networkConfiguration={ 'awsvpcConfiguration': {
            'subnets': ['subnet-093b63a5b7e5ae000' ,'subnet-0ce31ad86252fdc48' ],
            'securityGroups': ['sg-0b6645c0140f96693'],
            'assignPublicIp': 'ENABLED'
        } },
        overrides={
            'containerOverrides': [{
                'name': 'scroller-container',
                'command' : ['compchecker.py'],
                'environment': [
                    {'name': 'S3_BUCKET', 'value': S3_BUCKET},
                    {'name': 'workflow_id' , 'value' : workflow_id },  
                ]
            }]
        }
    )

    else:
        ecs = boto3.client('ecs' , region_name = "eu-west-2")
        response = ecs.run_task(
        cluster='scroller-cluster',
        taskDefinition='scroller-task',
        launchType='FARGATE',
        networkConfiguration={ 'awsvpcConfiguration': {
            'subnets': ['subnet-093b63a5b7e5ae000' ,'subnet-0ce31ad86252fdc48' ],
            'securityGroups': ['sg-0b6645c0140f96693'],
            'assignPublicIp': 'ENABLED'
        } },
        overrides={
            'containerOverrides': [{
                'name': 'scroller-container',
                'command' : ['screenshot.py'],
                'environment': [
                    {'name': 'S3_BUCKET', 'value': S3_BUCKET},
                    {'name': 'search_term' , 'value' : search_term },
                    {'name': 'workflow_id' , 'value' : workflow_id },  
                ]
            }]
        }
    )
    print("response:", response)
    

if __name__ == "__main__":
    main()