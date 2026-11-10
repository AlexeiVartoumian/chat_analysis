import json
import os

import boto3

client = boto3.client("ssm", region_name="eu-west-2")


keys = {}
for page in client.get_paginator("get_parameters_by_path").paginate(
        Path="/api/store/", WithDecryption=True):
    for p in page["Parameters"]:
        keys[p["Name"].split("/")[-1]] = p["Value"]

curPath = os.getcwd()


wif_path = os.path.join(curPath, "gcp-wif.json")
with open(wif_path, "w", encoding="utf-8") as f:
    f.write(keys["gcp-wif"])
json.load(open(wif_path))          


with open(os.path.join(curPath, ".env"), "w", encoding="utf-8") as f:
    for key, val in keys.items():
        if key == "gcp-wif":
            continue
        f.write(f'{key}="{val}"\n')
    f.write(f'GOOGLE_APPLICATION_CREDENTIALS="{wif_path}"\n')

# 3. .env2
with open(os.path.join(curPath, ".env2"), "w", encoding="utf-8") as f:
    f.write(f'POSTGRES_DB={keys["db_database"]}\n')
    f.write(f'POSTGRES_PASSWORD={keys["db_password"]}\n')
    f.write(f"GOOGLE_APPLICATION_CREDENTIALS={wif_path}\n")