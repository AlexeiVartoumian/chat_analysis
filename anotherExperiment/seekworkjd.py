import sys
import json
import boto3
import csv
import uuid

values = sys.stdin.read()
# stdin is JSON from Go's SendWorkweek(): [{"job_id": "...", "url": "..."}, ...]
payload = json.loads(values)

#print(payload)
if not payload:
    sys.exit(0)

S3_BUCKET = "output-store-work-store-390746273208"
SQS_QUEUE_URL = "https://sqs.eu-west-2.amazonaws.com/390746273208/workflow-cordinator-work"
def send_csv(leads, workflow_id, urls , path):
    s3 = boto3.client('s3')
    sqs = boto3.client('sqs')
    with open("/tmp/processedJobs.csv" , "w" , newline='' , encoding="utf-8") as csv_file:
        field_names = ["jobCategory" , "jobCategoryId" , "title" ,"externalPath" , "locationsText" , "job_id" ,"bulletFields", "companyName", "apiendpoint" , "job_url", "public_url","wd_instance"]
                
        writer = csv.DictWriter(csv_file ,fieldnames=field_names)
    
                #print(all_jobs)
        writer.writeheader()
        for job in leads:
            writer.writerow({k : job.get(k , '') for k in field_names})

    with open("/tmp/processedJobs.csv" , "rb" ) as csv_file:
                        #s3 = boto3.client("s3")
        contents = csv_file.read()
        s3.put_object(Bucket =S3_BUCKET , Key=f"processedJobs-{path}-{workflow_id}.csv" ,Body=contents)


    payload =  {
                'workflow_id': workflow_id,
                'file_name' : f"processedJobs-{path}-{workflow_id}.csv",
                'url_origin': urls["url"],
                'url_origin_id' :urls["job_id"]
            }
    sqs.send_message(   
                    QueueUrl=SQS_QUEUE_URL,
                    MessageBody=json.dumps(payload)
                )



workflow_id = str(uuid.uuid4())
send_csv(payload, workflow_id, {"url" :"postjdcheck", "job_id" : "oh"} , "jd")