import boto3
import os

client = boto3.client('s3')

##check available files . pair off against local "seen" file . 
#write download obj


resp = client.list_objects(Bucket='alexeitranscribefile')

keys = []
for key in resp["Contents"]:
    print(key["Key"])
    #print("ok \n")
    keys.append(key)
#print(resp)

for key in keys:
    with open ("processedJobs.csv" , "wb" ) as f :
        client.download_fileobj('alexeitranscribefile', key, f)

    with open("company_data.csv" , "wb" ) as f :
        client.download_fileobj('alexeitranscribefile', key, f)

    with open("job_metadata.csv" , "wb" ) as f :
        client.download_fileobj('alexeitranscribefile', key, f)

    with open("job_description.csv" , "wb" ) as f :
        client.download_fileobj('alexeitranscribefile', key, f)


# with open ("processedJobs.csv" , "wb" ) as f :
#     client.download_fileobj('alexeitranscribefile', 'input/processedJobs.csv', f)

# with open("company_data.csv" , "wb" ) as f :
#     client.download_fileobj('alexeitranscribefile', '/output/company_data.csv', f)

# with open("job_metadata.csv" , "wb" ) as f :
#     client.download_fileobj('alexeitranscribefile', '/output/job_metadata.csv', f)

# with open("job_description.csv" , "wb" ) as f :
#     client.download_fileobj('alexeitranscribefile', '/output/job_description.csv', f)


