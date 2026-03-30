import boto3
import sys

first_argument = None
if len(sys.argv) >1 :
    first_argument  = sys.argv[1]

client = boto3.client('route53')

resp = client.list_hosted_zones_by_name(DNSName="worldcaffeine.com")

zone_id = resp['HostedZones'][0]['Id']

if first_argument != None:
    response = client.change_resource_record_sets(
            HostedZoneId=zone_id,
            ChangeBatch={
                "Comment": "Automatic DNS update",
                "Changes": [
                    {
                        "Action": "UPSERT",
                        "ResourceRecordSet": {
                            "Name": "jobdice"+'.'+'worldcaffeine.com',
                            
                            "Type": "A",
                        
                            "TTL": 60,
                            "ResourceRecords": [
                                {
                                    "Value": first_argument,
                                },
                            ],
                        }
                    },
                ]
            }
        )
else:
    raise Exception(first_argument , "uh oh")
    

