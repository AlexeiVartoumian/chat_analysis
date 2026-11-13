import io
import json
import os
import zipfile
from collections import defaultdict
from datetime import datetime

import boto3

client = boto3.client('s3')
BUCKET = 'output-store-orac-store-390746273208'

# Check available files, pair them off against the local "seen" file, write the
# ones not seen yet. Loose .csv objects (the old Lambda output) and .zip objects
# (the GCP finaliser output) both end up as <name>_<timeline>.csv files here.
# A CSV inside a zip gets the zip's LastModified as its timeline.

paginator = client.get_paginator('list_objects_v2')
pages = paginator.paginate(Bucket=BUCKET)

# TODO keys.txt should come from db

timelines = []
new_keys = []

seen_keys = set()
new_keys_to_upload = []
sources = []          # same order as new_keys_to_upload: ("s3", key) or ("zip", zip_key, member)
zips = {}             # zip_key -> open ZipFile, for the new zips only
new_zip_keys = []

date_format = '%Y-%m-%d'

with open("keys_orac.txt", "r", encoding="utf-8") as f:
    for line in f:
        seen_keys.add(line.strip("\n"))

output_store_standard_timelines = set()


def add(check, last_modified, source):
    """Queue one CSV, unless it was seen on an earlier run."""
    if check in seen_keys:
        return
    datestring = datetime.strftime(last_modified, date_format)
    year_month_date = datetime.strptime(datestring, date_format)
    if year_month_date not in timelines:
        timelines.append(year_month_date)
        output_store_standard_timelines.add(year_month_date)
    new_keys_to_upload.append(check)
    sources.append(source)


for page in pages:
    for obj in page.get('Contents', []):
        key = obj['Key']
        timeline = "-".join(str(obj["LastModified"]).split(" ")).replace(":", "-")

        if key.endswith(".csv"):
            check = key.removesuffix(".csv") + "_" + timeline + ".csv"
            add(check, obj["LastModified"], ("s3", key))

        elif key.endswith(".zip") and key not in seen_keys:
            # GCP export: open it once, queue every CSV inside under its own name
            body = client.get_object(Bucket=BUCKET, Key=key)["Body"].read()
            zf = zipfile.ZipFile(io.BytesIO(body))
            zips[key] = zf
            new_zip_keys.append(key)
            from_lifecycle = os.path.basename(key).startswith("lifecycle")
            for member in zf.namelist():
                name = os.path.basename(member)          # processedJobsOrac-<uuid>.csv
                if from_lifecycle and name.startswith("processedJobsOrac-"):
                    # lifecycle zip: all live jobs for a company, mark it for the diff
                    name = "processedJobsOracLifecycleGcp-" + name.split("-", 1)[1]
                if name.endswith(".csv"):
                    check = name.removesuffix(".csv") + "_" + timeline + ".csv"
                    add(check, obj["LastModified"], ("zip", key, member))


timelines = sorted(timelines)

for i in timelines:
    print(datetime.strftime(i, date_format))

with open("newkeys_orac.txt", "w", encoding="utf-8") as f:
    for key in new_keys_to_upload:
        f.write(key + "\n")


def forgive(key, mydict):
    if key not in mydict:
        # [company, processedJobs, job_descriptions, deadlinks, lifecycle (GCP)]
        mydict[key] = [0] * 5
    return mydict


with open("keys_orac.json", "w", encoding="utf-8") as f:
    output = defaultdict(list)

    # need this since this guarantees the order
    for timeline in timelines:
        standard_timeline = datetime.strftime(timeline, date_format)
        # it could be a standard job does not happen on a given day but a backfill one does
        if timeline in output_store_standard_timelines:
            output[standard_timeline] = [{}]
            print("standard", standard_timeline, output[standard_timeline])

    for key in new_keys_to_upload:
        # processedJobsOrac-<uuid>_<timeline>.csv -> <uuid>
        unique = key.split("-", 1)[1].split(".")[0].split("_")[0]
        document = os.path.basename(key)

        # TODO POTENTIAL FAILURE ON FILE NAME PARSING
        # 2026-04-24-21-02-43+00-00 -> [:10] -> 2026-04-24
        timeline_key = document.split("_")[-1].removesuffix(".csv")
        timeline_key = datetime.strptime(timeline_key[:10], date_format)
        timeline_key = datetime.strftime(timeline_key, date_format)

        new_keys.append(document)
        if document.startswith("company"):
            forgive(unique, output[timeline_key][0])
            output[timeline_key][0][unique][0] = document

        if document.startswith("processedJobsOracLifecycleGcp"):
            # lifecycle (GCP): every live job for one company, for the DB diff
            forgive(unique, output[timeline_key][0])
            output[timeline_key][0][unique][4] = document
        elif document.startswith("processedJobs"):
            forgive(unique, output[timeline_key][0])
            output[timeline_key][0][unique][1] = document

        if document.startswith("job_descriptionsOrac"):
            forgive(unique, output[timeline_key][0])
            output[timeline_key][0][unique][2] = document

        if document.startswith("deadlinks"):
            forgive(unique, output[timeline_key][0])
            output[timeline_key][0][unique][3] = document

    json.dump(output, f)

# write every new file locally, from S3 or straight out of its zip
for document, source in zip(new_keys, sources):
    with open(document, "wb") as f:
        if source[0] == "s3":
            client.download_fileobj(BUCKET, source[1], f)
        else:
            f.write(zips[source[1]].read(source[2]))

# only mark things as seen once they are actually on disk
with open("keys_orac.txt", "a", encoding="utf-8") as f:
    for key in new_keys_to_upload + new_zip_keys:
        f.write(key + "\n")