import requests


def  lambda_handler(event , context):
    
    ip = requests.get("https://api.ipify.org?format=json").json()["ip"]
    print(f"outbound IP {ip}")
    return {"statusCode": 200, "body": ip}

