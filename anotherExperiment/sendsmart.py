import json
import sys
import uuid

import boto3

REGION = "eu-west-2"
LAMBDA_ARN = "arn:aws:lambda:eu-west-2:390746273208:function:seeksmartposts"
SOURCE_BUCKET = "output-store-smart-store-390746273208" 
QUEUE_URL = "https://sqs.eu-west-2.amazonaws.com/390746273208/workflow-cordinator-smart"  

# stdin is JSON from Go, either a single object or a list of rows
data = json.loads(sys.stdin.read())
items = data if isinstance(data, list) else [data]
if not items:
    raise SystemExit("empty payload, nothing to do")
first = items[0]

workflow_id = str(uuid.uuid4())


def invoke_lambda(payload):
    """Hand the payload straight to the Lambda, async."""
    client = boto3.client("lambda", region_name=REGION)
    resp = client.invoke(
        FunctionName=LAMBDA_ARN,
        InvocationType="Event",  # async, fire and forget
        Payload=json.dumps(payload).encode("utf-8"),
    )
    # async invokes return 202 once Lambda has accepted the event
    if resp["StatusCode"] != 202:
        raise SystemExit(f"Lambda invoke failed: {resp['StatusCode']}")
    print(f"Invoked Lambda async ({len(items)} items)")


def to_s3(workflow):
    """Write Go's payload as-is to S3; return the file name."""
    s3 = boto3.client("s3", region_name=REGION)
    file_name = f"{workflow}-{workflow_id}.json"
    s3.put_object(
        Bucket=SOURCE_BUCKET,
        Key=file_name,
        Body=json.dumps(data).encode("utf-8"),
    )
    return file_name


def to_sqs(message):
    sqs = boto3.client("sqs", region_name=REGION)
    resp = sqs.send_message(QueueUrl=QUEUE_URL, MessageBody=json.dumps(message))
    print(f"Sent message {resp['MessageId']}")


if first.get("company_id") and first.get("company_apply_url"):
    # lifecycle: company + its known job_ids, in a file
    to_sqs({
        "workflow": "lifecycle",
        "workflow_id": workflow_id,
        "company_id": first["company_id"],
        "file_name": to_s3("lifecycle"),
    })

elif first.get("job_id") and first.get("url"):
    # new links: skip the queue, hand straight to the Lambda
    invoke_lambda(items)

elif first.get("job_id") and first.get("api_endpoint"):
    # detail: batch of job_id + api_endpoint rows, via S3 + SQS
    to_sqs({
        "workflow": "detail",
        "workflow_id": workflow_id,
        "file_name": to_s3("detail"),
    })

else:
    raise SystemExit(f"don't know which workflow this is: keys={sorted(first)}")