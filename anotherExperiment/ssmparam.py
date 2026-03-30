import boto3
import os

client = boto3.client('ssm' , region_name="eu-west-2" )

resp = client.get_parameters_by_path(Path="/api/store/")

keys= {}
for val in resp["Parameters"]:

    Name = val["Name"].split("/")[-1]
    Value = val["Value"]
    keys[Name]= Value

curPath = os.getcwd()

curPath1 = os.path.join(curPath , ".env")
with open( curPath1 ,"w" , encoding="utf-8" ) as f :

    for key , val in keys.items() :
        f.write(f'{key}="{val}"\n')

curPath2 =os.path.join(curPath , ".env2")
with open( curPath2 ,"w" , encoding="utf-8" ) as f :
        f.write(f'POSTGRES_DB={keys["db_database"]}\n')
        f.write(f'POSTGRES_PASSWORD={keys["db_password"]}\n')