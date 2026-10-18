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
SOURCE_BUCKET = "output-store-orac-store-390746273208"


values = sys.stdin.read()
data = json.loads(values)

items = data if isinstance(data, list) else [data]
first = items[0]

workflow_id = str(uuid.uuid4())
s3 = boto3.client("s3", region_name="eu-west-2")
sqs = boto3.client("sqs", region_name="eu-west-2")

workflow_id = str(uuid.uuid4())
message = None
def to_s3(workflow):
    """Write Go's payload as-is to S3; return the file name."""
    file_name = f"{workflow}-{workflow_id}.json"
    s3.put_object(Bucket=SOURCE_BUCKET, Key=file_name, Body=json.dumps(data).encode("utf-8"))
    return file_name

if first.get("company_id") and first.get("company_apply_url"):
    # lifecycle: company + its known job_ids, in a file
    message = {
        "workflow": "lifecycle",
        "workflow_id": workflow_id,
        "company_id": first["company_id"],
        "file_name": to_s3("lifecycle"),
    }

elif first.get("job_id") and first.get("api_url"):
    # detail: up to 350 job_id + api_url rows, in a file
    message = {
        "workflow": "detail",
        "workflow_id": workflow_id,
        "file_name": to_s3("detail"),
    }

elif first.get("job_id") and first.get("url"):
    # normal: one new link, small enough for the message itself
    message = {
        "workflow": "normal",
        "workflow_id": workflow_id,
        "company_apply_url": first["url"],
        "link_id": first["job_id"],
    }


else:
    raise SystemExit(f"don't know which workflow this is")



sqs = boto3.client("sqs", region_name="eu-west-2")
response = sqs.send_message(QueueUrl=QUEUE_URL, MessageBody=json.dumps(message))
print(f"Sent message {response['MessageId']}")