import boto3
import os
from boto3.dynamodb.conditions import Key, Attr
from botocore.exceptions import ClientError
import time
import uuid
dynamodb = boto3.resource('dynamodb')
table = dynamodb.Table('filepool')

response = table.query(
        IndexName='status-index',
        KeyConditionExpression=Key('status').eq('FREE'),
        Limit=1
    )
print(response)
file = response['Items'][0]
now = int(time.time())
ttl = now + (22 * 60)
workflow_id = str(uuid.uuid4())
try:
    table.update_item(
        Key={'file_id': file['file_id']},
        UpdateExpression='''
            SET #s = :locked,
                workflow_id = :wf_id,
                locked_at = :now,
                #t = :ttl
        ''',
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
    
except ClientError as e:
    if e.response['Error']['Code'] == 'ConditionalCheckFailedException':
        print(e)
    raise

response = table.query(
       
        KeyConditionExpression=Key('file_id').eq('BOO'),
        Limit=1
    )
print(response, "\n")


response = table.query(
       
        KeyConditionExpression=Key('file_id').eq('yar'),
        Limit=1
    )
print(response , "\n")