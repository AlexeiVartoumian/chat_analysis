import json, sys
import boto3, requests

BUCKET = "output-store-orac-store-390746273208"   # or wherever your diff reads from

req = json.loads(sys.stdin.read())
key = req["object"].rsplit("/", 1)[-1]            # lifecycle-20261111T201500Z.zip
with requests.get(req["url"], stream=True, timeout=60) as r:
    r.raise_for_status()
    boto3.client("s3", region_name="eu-west-2").upload_fileobj(r.raw, BUCKET, key)
print(f"copied {req['object']} ({req['files']} files) -> s3://{BUCKET}/{key}")