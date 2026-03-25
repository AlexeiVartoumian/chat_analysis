
# from collections import defaultdict

# import csv


# values = set()

# duplicates =[]
# totals = defaultdict(int)

# ##quick script to see company:jobs ratio
# with open(  "companiesid.txt" , "r" ) as file:

#     for line in file:

#         if int(line) in values:
#             duplicates.append(int(line))
#         else:
#             values.add(int(line))
        
#         totals[int(line)] +=1

# # print(len(values))

# # print(len(totals))
# # print(totals)

# inserted = set()
# with open(  "inserted.txt" , "r" ) as file:

#     for line in file:

#         inserted.add(int(line)) 


# diff = values - inserted

# print(diff)

# # with open("C:/Users/wwwal/Downloads/processedJobs.csv" , "r") as file :

# #     csv_data = csv.DictReader(file)

# #     for index , line in enumerate(csv_data):
# #             if line["company_id"] != "N/A" and int(line["company_id"]) not in inserted :
# #                  print(index , int(line["company_id"] ) )


# info = []
# with open("C:/Users/wwwal/Downloads/job_metadata.csv" , "r") as file :

#     csv_data = csv.DictReader(file)

#     for index , line in enumerate(csv_data):
           
#         data = defaultdict(str)
        
#         data["job_id"] = line["job_id"]
#         data["job_state"] = line["applicants_count"]
#         data["company_apply_url"] = line["job_state"]
#         data["applicants_count"] = line["company_apply_url"]

#         info.append(data)

# print(len(info))
# with open ("C:/Users/wwwal/Downloads/job_metadata2.csv" , "w" ) as file :

#     fields = ['job_id', 'applicants_count', 'company_apply_url', "job_state"]

#     new_csv  = csv.DictWriter(file , fields )
#     new_csv.writeheader()
#     for data in info:
#         new_csv.writerow(data)