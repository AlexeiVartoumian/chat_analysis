import sys 
import json
import boto3

roles = json.loads(sys.stdin.read())

s3 = boto3.client("s3")
s3.put_object(
    Bucket="backfill-store-390746273208",
    Key=f"live-roles.json",
    Body=json.dumps(roles, indent=2),
    ContentType="application/json",
)
