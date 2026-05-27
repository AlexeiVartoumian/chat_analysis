import sys 
import json
import boto3
import time
import uuid
from boto3.dynamodb.conditions import Key, Attr
from botocore.exceptions import ClientError

#search type is suspended or live
search_type = None
auto = None
if len(sys.argv) > 1:
    search_type = sys.argv[1]
    auto = sys.argv[2]

roles = json.loads(sys.stdin.read())

workflow_id = str(uuid.uuid4())

s3 = boto3.client("s3", region_name='eu-west-2')

S3_BUCKET = "somebuckethaha"







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
            'name': 'reader-container',
            'environment': [
                {'name': 'S3_BUCKET', 'value': S3_BUCKET}, 
            ]
        }]
    }
)
    print("response:", response)
    

if __name__ == "__main__":
    main()