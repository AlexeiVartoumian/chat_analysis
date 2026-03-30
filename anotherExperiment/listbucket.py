import boto3
import os

client = boto3.client('s3')



with open ("processedJobs.csv" , "wb" ) as f :
    client.download_fileobj('alexeitranscribefile', 'input/processedJobs.csv', f)

with open("company_data.csv" , "wb" ) as f :
    client.download_fileobj('alexeitranscribefile', '/output/company_data.csv', f)

with open("company_data.csv" , "wb" ) as f :
    client.download_fileobj('alexeitranscribefile', '/output/job_metadata.csv', f)

with open("company_data.csv" , "wb" ) as f :
    client.download_fileobj('alexeitranscribefile', '/output/job_description.csv', f)


