import sys 
import json
import boto3
import time
import uuid
from boto3.dynamodb.conditions import Key, Attr
from botocore.exceptions import ClientError

roles = json.loads(sys.stdin.read())

s3 = boto3.client("s3")

source_store = "source-store-390746273208"
output_store = "backfill-store-390746273208"
s3.put_object(
    Bucket=output_store,
    Key=f"live-roles.json",
    Body=json.dumps(roles, indent=2),
    ContentType="application/json",
)
dynamodb = boto3.resource('dynamodb')

filepool_table_name = "filepoolstore"
accountpool_table_name = "accountpoolstore"

filepool_table = dynamodb.Table(filepool_table_name)
account_table = dynamodb.Table(accountpool_table_name)



def acquire_lock(workflow_id):
    response = filepool_table.query(
        IndexName='status-index',
        KeyConditionExpression=Key('status').eq('FREE'),
        Limit=1
    )
    
    if not response['Items']:
        return None
    
    file = response['Items'][0]
    now = int(time.time())
    ttl = now + (60 * 60)
    
    try:
        filepool_table.update_item(
            Key={'file_id': file['file_id']},
            UpdateExpression='SET #s = :locked, workflow_id = :wf_id, locked_at = :now, #t = :ttl',
            ConditionExpression=Attr('status').eq('FREE'),
            ExpressionAttributeNames={
                '#s': 'status',
                '#t': 'ttl'
            },
            ExpressionAttributeValues={
                ':locked': 'LOCKED',
                ':wf_id': workflow_id,
                ':now': now,
                ':ttl': ttl
            }
        )
        return file
    except ClientError as e:
        if e.response['Error']['Code'] == 'ConditionalCheckFailedException':
            return None
        raise

def main():

    workflow_id = str(uuid.uuid4())
    file = acquire_lock(workflow_id)

    if not file:
        raise Exception
    
    #file["file_id"]
    

    ecs = boto3.client('ecs' , region_name = "eu-west-2")
    ecs.run_task(
    cluster='reader-cluster',
    taskDefinition='reader-task:3',
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
                {'name': 's3_source_bucket', 'value': source_store},
                {'name': 'output_store',    'value': output_store},
                {'name': 'file_id',    'value': file["file_id"] },
            ]
        }]
    }
)