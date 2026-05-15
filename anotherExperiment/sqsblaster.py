import sys
import json
import time 
import boto3



search_terms = json.loads(sys.stdin.read())

print("*" * 75)




searchtermid = None
search_term = None

sqs = boto3.client('sqs', region_name='eu-west-2')

for term in search_terms:

        print(term)
        print("\n")

        search_term_id = term["search_term_id"]
        search_term = term["search_term"]
        response = sqs.send_message(
        QueueUrl='https://sqs.eu-west-2.amazonaws.com/390746273208/workflow-cordinator-test',
        MessageBody=json.dumps({
            "search_term_id": search_term_id,
            "search_term": search_term,
            "target_bucket": "somebuckethaha"  #todo remove
        })
        )
        print(search_term_id)
        print(search_term)
        print(type(search_term_id))
        print(type(search_term))

        print(f"Message sent to queue: {response['MessageId']}")
        time.sleep(1)










#print("farambula")