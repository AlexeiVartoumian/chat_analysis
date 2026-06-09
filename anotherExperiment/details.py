import sys 
import json
import boto3
import time
import uuid
from boto3.dynamodb.conditions import Key, Attr
from botocore.exceptions import ClientError




s3 = boto3.client("s3", region_name='eu-west-2')
workflow_id = str(uuid.uuid4())
auto = None


if len(sys.argv) >1:
    auto = sys.argv[1]
    companies = json.loads(sys.stdin.read())
    s3.put_object(Bucket = "backfill-store-390746273208" , Key=f"companyDetail-{workflow_id}-roles.json" , Body=companies)


def main():
    
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
            'command' : ['companylink.py'],
            'environment': [
                {'name': 's3_source_bucket', 'value': "source-store-390746273208"},
                {'name': 'output_store', 'value': "backfill-store-390746273208"},
                {'name': 'file_pool_table', 'value': "filepoolstore"},
                {'name': 'file_id' , 'value' : "cookies-grouped-7.json" },
                {'name': 'workflow_id' , 'value' : workflow_id },  
                {'name': 'auto' , 'value' : "false" },  
            ]
        }]
    }
    )
    print("response:", response)
    
if __name__ == "__main__":
    main()