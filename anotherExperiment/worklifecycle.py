import sys
import json
import boto3

LAMBDA_ARN = "arn:aws:lambda:eu-west-2:390746273208:function:seekworkposts"

values = sys.stdin.read()

# stdin is JSON from Go's SendWorkweek(): [{"job_id": "...", "url": "..."}, ...]
payload = json.loads(values)

print(payload)
lambda_client = boto3.client('lambda' , region_name="eu-west-2")

response = lambda_client.invoke(
    FunctionName=LAMBDA_ARN,
    InvocationType='Event',  # async - fire and forget
    Payload=json.dumps(payload).encode('utf-8')
)

print(f"Status code: {response['StatusCode']}")

if 'FunctionError' in response:
    print(f"Function error: {response['FunctionError']}")
    print(response['Payload'].read().decode('utf-8'))