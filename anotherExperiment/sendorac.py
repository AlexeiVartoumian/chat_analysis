# import sys
# import json
# import boto3

# LAMBDA_ARN = "arn:aws:lambda:eu-west-2:390746273208:function:orac"

# values = sys.stdin.read()

# # stdin is JSON from Go's SendWorkweek(): [{"job_id": "...", "url": "..."}, ...]
# payload = json.loads(values)

# print(payload)
# lambda_client = boto3.client('lambda' , region_name="eu-west-2")

# response = lambda_client.invoke(
#     FunctionName=LAMBDA_ARN,
#     InvocationType='Event',  # async - fire and forget
#     Payload=json.dumps(payload).encode('utf-8')
# )

# print(f"Status code: {response['StatusCode']}")

# if 'FunctionError' in response:
#     print(f"Function error: {response['FunctionError']}")
#     print(response['Payload'].read().decode('utf-8'))




import sys
import json
import uuid
import boto3

QUEUE_URL = "https://sqs.eu-west-2.amazonaws.com/390746273208/workflow-cordinator-orac"

values = sys.stdin.read()

# stdin is JSON from Go: {"job_id": "...", "url": "..."} (a one-item list also works)
link = json.loads(values)
if isinstance(link, list):
    link = link[0]

workflow_id = str(uuid.uuid4())
message = None

if link.get("company_id") and link.get("company_apply_url"):
    # lifecycle: rescan a company we already know
    message = {
        "workflow": "lifecycle",
        "workflow_id": workflow_id,
        "company_id": link["company_id"],
        "company_apply_url": link["company_apply_url"],
    }

elif link.get("job_id") and link.get("api_url"):
    # detail: fetch the full description for one job
    message = {
        "workflow": "detail",
        "workflow_id": workflow_id,
        "job_id": link["job_id"],
        "api_url": link["api_url"],
    }

elif link.get("job_id") and link.get("url"):
    # normal: a new LinkedIn link to chase
    message = {
        "workflow": "normal",
        "workflow_id": workflow_id,
        "company_apply_url": link["url"],
        "link_id": link["job_id"],
    }

else:
    raise SystemExit(f"don't know which workflow this is: {link}")



sqs = boto3.client("sqs", region_name="eu-west-2")
response = sqs.send_message(QueueUrl=QUEUE_URL, MessageBody=json.dumps(message))
print(f"Sent message {response['MessageId']}")