hi hello hello 
I have a db on uni servers and i want to migrate it . 
thing is its in mysql and want to migrate it to postgres local.

breakdown : 
1. create postgres schema
2. insert fresh data into it
3. modify the tdift and vector stuff for the description table
4. lets see


5. sanitize text inputs
6. TODO : employee count is getting zeroes values when avaliable number exists . inverstigate.

7. normalise data!

8. add user table for api keys to validate against 

9. add all the routes 
10. add backfill embedding route

11. .update reader lambda env keys . trigger automation for s3 bucket workflow.
create api endpoint pto update this . think about api being in lambda . think about iam permissions . for now test on ec2 .

12. add the semantic search field as a column somewhere ? with timestamp ? also there is thow to connectr the s3 to the ec2