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

"""
 order of ops . 
 if first then create instance and capture the instance id then send command for profile.
 getting profile will be a dynamodb query . 
 need the id for future runs.
 thinking to send it as a env variable to query in the script afterwards 
 else : send search term search term id to be consumed .
"""



if len(sys.argv) > 1:
    first_run = sys.argv[1]
    instance_id = sys.argv[2]


search_terms = json.loads(sys.stdin.read())
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
file = None
workflow_id = str(uuid.uuid4())
for term in search_terms:

    search_term_id = term["search_term_id"]
    search_term = term["search_term"]

    # if first_run == "yes":
    
    file = acquire_lock(workflow_id)
    scroller_worker = f"scroller-{scroller_count}"
    scroller_count+=1
    response = ec2_client.run_instances(LaunchTemplate={'LaunchTemplateId': 'lt-0989cbde3348f9e83', 'Version': '$Latest'} ,InstanceInitiatedShutdownBehavior='terminate',MinCount=1,MaxCount=1, TagSpecifications=[{'ResourceType': 'instance','Tags': [{'Key': 'Name', 'Value': f'{scroller_worker}'},{'Key': 'Role', 'Value': 'time-to-go'}]}])

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



    

    ssm_client = boto3.client('ssm', region_name="eu-west-2")

    if instance_id:

        wait_for_ssm_online(ssm_client, instance_id)
        print(instance_id)
      
        search_term_id = shlex.quote(str(search_term_id))
        inner_cmd = (
            f"export DISPLAY=:1 && "
            f"export workflow_id={shlex.quote(workflow_id)} && "
            f"export profile_id={shlex.quote(file['profile_id'])} && "
            f"export search_term={shlex.quote(search_term)} && "
            f"export search_term_id={search_term_id} &&"
            f"cd /opt/myapp && " 
            f"timeout 1800 /venv/bin/python3 screenshot.py > /tmp/last_run.log 2>&1; echo EXIT_CODE:$?"
            f"echo EXIT_CODE:$? > /tmp/exit_code.txt"
        )

        full_cmd = f"sudo -u ubuntu bash -c {shlex.quote(inner_cmd)}"
        commands = [
            full_cmd,
            "cat /tmp/exit_code.txt",   
            "sleep 15",                 # CloudWatch flush logs
            "shutdown -h now",          
        ]#want this to be self terminating

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

