import sys 
import json
import boto3
import time
import uuid
from boto3.dynamodb.conditions import Key, Attr
from botocore.exceptions import ClientError



values = sys.stdin.read()

