import json
import boto3
from boto3.dynamodb.conditions import Key, Attr
from botocore.exceptions import ClientError
import time
import uuid
import os
import sys
import time
import shlex
import base64

"""
 ssm based workflow . reading from stdin base64 encode input to export
 as file to be read by instance.
 first pass need to create isntance passing following params
 
 firstrun: bool , numinstanc: int , instanceid: str
 eg
 on firstrun -> true , 3 , ""
"""



if len(sys.argv) > 1:
    numberof = sys.argv[1]
    first_run = sys.argv[2]
    instance_id = sys.argv[3]

leads = json.loads(sys.stdin.read())
encoded = base64.b64encode(json.dumps(leads).encode()).decode()
dynamodb = boto3.resource('dynamodb' , region_name='eu-west-2')

filepool_table_name = "filepoolstore_deed"

filepool_table = dynamodb.Table(filepool_table_name)

ec2_client = boto3.client('ec2',  region_name='eu-west-2')







def wait_for_instance_running(ec2_client, instance_id, timeout=300):
    waiter = ec2_client.get_waiter('instance_running')
    waiter.wait(InstanceIds=[instance_id], WaiterConfig={'Delay': 5, 'MaxAttempts': timeout // 5})

def wait_for_ssm_online(ssm_client, instance_id, timeout=300, interval=5):
    elapsed = 0
    while elapsed < timeout:
        resp = ssm_client.describe_instance_information(
            Filters=[{'Key': 'InstanceIds', 'Values': [instance_id]}]
        )
        infos = resp.get('InstanceInformationList', [])
        if infos and infos[0]['PingStatus'] == 'Online':
            return True
        time.sleep(interval)
        elapsed += interval
    raise TimeoutError(f"Instance {instance_id} did not become SSM-online in time")
def wait_for_command(ssm_client, command_id, instance_id, timeout=120, interval=3):
    elapsed = 0
    while elapsed < timeout:
        try:
            inv = ssm_client.get_command_invocation(CommandId=command_id, InstanceId=instance_id)
            if inv['Status'] not in ('Pending', 'InProgress'):
                return inv
        except ssm_client.exceptions.InvocationDoesNotExist:
            pass
        time.sleep(interval)
        elapsed += interval
    raise TimeoutError("Command did not finish in time")




def acquire_lock(workflow_id):
    response = filepool_table.query(
        IndexName='status-index',
        KeyConditionExpression=Key('status').eq('FREE'),
        Limit=1
    )
    if not response['Items']:
        return None
    
    file = response['Items'][0]
    now = int(time.time())
    try:
        filepool_table.update_item(
            Key={'profile_id': file['profile_id']},
            UpdateExpression='SET #s = :locked, workflow_id = :wf_id, locked_at = :now',
            ConditionExpression=Attr('status').eq('FREE'),
            ExpressionAttributeNames={
                '#s': 'status'
            },
            ExpressionAttributeValues={
                ':locked': 'LOCKED',
                ':wf_id': workflow_id,
                ':now': now,
            }
        )
        return file
    except ClientError as e:
        if e.response['Error']['Code'] == 'ConditionalCheckFailedException':
            return None
        raise

scroller_count = 1
numberof = int(numberof)
for count in range(numberof):

    if first_run == "true":
        workflow_id = str(uuid.uuid4())
        file = acquire_lock(workflow_id)
        scroller_worker = f"greenlead-{scroller_count}"
        scroller_count+=1
        response = ec2_client.run_instances(LaunchTemplate={'LaunchTemplateId': 'lt-0989cbde3348f9e83', 'Version': '$Latest'} ,MinCount=1,MaxCount=1, TagSpecifications=[{'ResourceType': 'instance','Tags': [{'Key': 'Name', 'Value': f'{scroller_worker}'}]}])

        instance_id = response['Instances'][0]['InstanceId']

    

        ssm_client = boto3.client('ssm', region_name="eu-west-2")
        wait_for_instance_running(ec2_client, instance_id)
        wait_for_ssm_online(ssm_client, instance_id)
        
        

        ssm_response = ssm_client.send_command(
            InstanceIds=[instance_id],
            DocumentName="AWS-RunShellScript",
            Parameters={
                "commands": [
                    f"sudo -u ubuntu bash -c 'export DISPLAY=:1 && export Instance_id={instance_id} && export S3_BUCKET=alexeitranscribefile && export PROFILE_S3_KEY={file['profile_id']} && cd /opt/myapp && /venv/bin/python3 extractprofile.py'"
                ]
            }
        )



    workflow_id = str(uuid.uuid4())

    ssm_client = boto3.client('ssm', region_name="eu-west-2")

    if instance_id:

        wait_for_ssm_online(ssm_client, instance_id)
        print(instance_id)
      
        inner_cmd = (
            f"export DISPLAY=:1 && "
            f"export workflow_id={shlex.quote(workflow_id)} && "
            f"echo {shlex.quote(encoded)} | base64 -d > /tmp/leads.json && "
            f"cd /opt/myapp && /venv/bin/python3 greenhouse.py < /tmp/leads.json > /tmp/last_run.log 2>&1; echo EXIT_CODE:$?"
        )

        full_cmd = f"sudo -u ubuntu bash -c {shlex.quote(inner_cmd)}"
        ssm_response = ssm_client.send_command(
            InstanceIds=[instance_id],
            DocumentName="AWS-RunShellScript",
            Parameters={"commands": [full_cmd]}
        )

        command_id = ssm_response['Command']['CommandId']

        # give it a moment to start
        time.sleep(3)

        invocation = ssm_client.get_command_invocation(
            CommandId=command_id,
            InstanceId=instance_id
        )

        print( 'commmandid', invocation['CommandId'] )
        print("Status:", invocation['Status'])
        print("--- STDOUT ---")
        print(invocation['StandardOutputContent'])
        print("--- STDERR ---")
        print(invocation['StandardErrorContent'])

            # result = wait_for_command(ssm_client, command_id, instance_id)
            # print(result['Status'])
            # print(result['StandardOutputContent'])
            # print(result['StandardErrorContent'])


            
