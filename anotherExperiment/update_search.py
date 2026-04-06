# import boto3
# import os
# import json
# import sys

# search_term = None

# if len(sys.argv) >1 :

#     search_term = sys.argv[1]


# client = boto3.client('lambda', region_name='eu-west-2')


# if search_term != None:
#     response = client.invoke(
#         FunctionName='reader',
#         InvocationType='Event', 
#         Payload=json.dumps({
#             "search_term": search_term
#         })
#     )
# else:
#     raise Exception(search_term , "whoopsie")

import boto3
import json
import sys
import uuid
# dynamodb = boto3.resource('dynamodb')
# table = dynamodb.Table('workflowstate')

search_term = None

if len(sys.argv) > 1:
    search_term = sys.argv[1]

if search_term is None:
    raise Exception("search_term is required")

sqs = boto3.client('sqs', region_name='eu-west-2')

response = sqs.send_message(
    QueueUrl='https://sqs.eu-west-2.amazonaws.com/390746273208/workflow-requests',
    MessageBody=json.dumps({
        "search_term": search_term,
        "target_bucket": "somebuckethaha"  
    })
)

print(f"Message sent to queue: {response['MessageId']}")




# sns = boto3.client('sns' , region_name='eu-west-2')
# workflow_id = str(uuid.uuid4())
# item = {
#         "workflow_id":  workflow_id,
#         "completed_at": None,
#         "locked_at": 0,
#         "status": "FREE",
#         "ttl": None
#     }

# table.put_item(
#     Item=item
# )


# response = sns.publish(
#     TopicArn='arn:aws:sns:eu-west-2:390746273208:requestsSns',
#     Message=json.dumps({
#         "search_term": search_term,
#         "target_bucket": "somebuckethaha",
#         "workflow_id" : workflow_id  
#     })
# )