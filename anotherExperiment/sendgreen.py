import sys
import json
import uuid
import boto3

QUEUE_URL = "https://sqs.eu-west-2.amazonaws.com/390746273208/workflow-cordinator-green"



values = sys.stdin.read()
data = json.loads(values)

items = data if isinstance(data, list) else [data]
first = items[0]

workflow_id = str(uuid.uuid4())
s3 = boto3.client("s3", region_name="eu-west-2")
sqs = boto3.client("sqs", region_name="eu-west-2")

workflow_id = str(uuid.uuid4())
message = None


if first.get("url") and first.get("company"): 
    message = {
        "workflow": "normal",
        "workflow_id": workflow_id,
        "items": [{
                "company_apply_url": first["url"],
                "company": first["company"],
            }],
    }

else:
    raise SystemExit(f"don't know which workflow this is")



sqs = boto3.client("sqs", region_name="eu-west-2")
response = sqs.send_message(QueueUrl=QUEUE_URL, MessageBody=json.dumps(message))
print(f"Sent message:  - Green {response['MessageId']}")