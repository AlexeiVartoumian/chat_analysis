import sys 
import json
import boto3
import time
import uuid
from boto3.dynamodb.conditions import Key, Attr
from botocore.exceptions import ClientError

#search type is suspended or live
# search_type = None
# auto = None
# if len(sys.argv) > 1:
#     search_type = sys.argv[1]
#     auto = sys.argv[2]

# roles = json.loads(sys.stdin.read())

# workflow_id = str(uuid.uuid4())

search_term = "cloud engineer"
if len(sys.argv) >1:

    search_term = sys.argv[1]
    search_term = search_term.lstrip('{"')
    search_term = search_term.rstrip('}"')

s3 = boto3.client("s3", region_name='eu-west-2')

S3_BUCKET = "somebuckethaha"

workflow_id = str(uuid.uuid4())

def main():
    
    #file["fil
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