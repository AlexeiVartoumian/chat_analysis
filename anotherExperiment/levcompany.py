import sys
import json
import boto3

LAMBDA_ARN = "arn:aws:lambda:eu-west-2:390746273208:function:orchestratorlev"

# stdin: {"company_url": "...", "job_ids": ["...", ...]}
raw = sys.stdin.read().strip()
if not raw or raw == "null":
    print("Nothing to send")
    sys.exit(0)

payload = json.loads(raw)
print(payload)

lambda_client = boto3.client("lambda", region_name="eu-west-2")
response = lambda_client.invoke(
    FunctionName=LAMBDA_ARN,
    InvocationType="Event",
    Payload=json.dumps(payload).encode("utf-8"),
)
print(f"Status code: {response['StatusCode']}")  # 202 = queued