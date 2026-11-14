import json
import os
import sys
import uuid

import oci

# Oracle side of sendoracgcp.py: same input, same messages, different transport.
#   jobs     -> intake Queue (Connector Hub hands them to the dispatcher)
#   finalize -> finalize Notifications topic (triggers the finaliser)
QUEUE_ID = "ocid1.queue.oc1.phx.amaaaaaa3uy77uaacq5okejv6timeab4ncel7vlimj26hrtsn4vqzpgjfeaa"
QUEUE_ENDPOINT = "https://cell-1.queue.messaging.us-phoenix-1.oci.oraclecloud.com"
FINALIZE_TOPIC_ID = "ocid1.onstopic.oc1.phx.amaaaaaa3uy77uaanvb4jyey44mvpecjpibgoxsohoopukd2e6px22ux5b5q"
FINALIZE_ENDPOINT = "https://cell1.notification.us-phoenix-1.oci.oraclecloud.com"

OCI_CONFIG = os.environ.get("OCI_CONFIG", "/etc/oci/config")
OCI_PROFILE = os.environ.get("OCI_PROFILE", "fanout-publisher")

# Smaller than GCP's 300 for two Oracle limits:
#   - the dispatcher forwards each job through Notifications, which caps messages at ~64 KB
#     (each job_id + api_url pair is ~250-300 bytes)
#   - the worker stops at a 250 s budget (Oracle kills functions at 300 s), ~1 s per job
DETAIL_BATCH = 150
QUEUE_BATCH = 20            # most messages one PutMessages call accepts

# Everything that should land in one zip must share a run_id; pass RUN_ID to group calls.
RUN_ID = os.environ.get("RUN_ID") or str(uuid.uuid4())

data = json.loads(sys.stdin.read() or "[]")
items = data if isinstance(data, list) else [data]
if not items:
    print("nothing to send")
    sys.exit(0)
first = items[0]
config = oci.config.from_file(OCI_CONFIG, profile_name=OCI_PROFILE)

if first.get("company_apply_url") == "final":
    # Go's "all done" marker: ask the finaliser to zip what the workers wrote.
    ons = oci.ons.NotificationDataPlaneClient(config, service_endpoint=FINALIZE_ENDPOINT)
    resp = ons.publish_message(FINALIZE_TOPIC_ID, oci.ons.models.MessageDetails(
        title=f"finalize {RUN_ID}", body=json.dumps({"run_id": RUN_ID})))
    print(f"Sent finalize {resp.data.message_id}")
    sys.exit(0)

if first.get("company_id") and first.get("company_apply_url"):
    # lifecycle: no job_ids needed any more, AWS already has them for the diff
    messages = [{"workflow": "lifecycle", "company_id": row["company_id"],
                 "company_apply_url": row["company_apply_url"]} for row in items]

elif first.get("job_id") and first.get("api_url"):
    jobs = [{"job_id": row["job_id"], "api_url": row["api_url"]} for row in items]
    messages = [{"workflow": "detail", "jobs": jobs[i:i + DETAIL_BATCH]}
                for i in range(0, len(jobs), DETAIL_BATCH)]

elif first.get("url") and first.get("company"):
    messages = [{"workflow": "normal", "company_apply_url": first["url"],
                 "company": first["company"]}]

elif first.get("job_id") and first.get("url"):
    # renamed to link_id: "job_id" is the dispatch id and must not be overwritten
    messages = [{"workflow": "normal", "company_apply_url": first["url"],
                 "link_id": first["job_id"]}]

else:
    raise SystemExit("don't know which workflow this is")

queue = oci.queue.QueueClient(config, service_endpoint=QUEUE_ENDPOINT)

# dispatch fields last so nothing in the payload can clobber them;
# job_id must be unique per message (the dispatcher dedupes on it)
bodies = [{**m, "job_id": str(uuid.uuid4()), "run_id": RUN_ID} for m in messages]

for i in range(0, len(bodies), QUEUE_BATCH):
    batch = bodies[i:i + QUEUE_BATCH]
    resp = queue.put_messages(QUEUE_ID, oci.queue.models.PutMessagesDetails(
        messages=[oci.queue.models.PutMessagesDetailsEntry(content=json.dumps(b)) for b in batch]))
    for sent in resp.data.messages:
        print(f"Sent message {sent.id}")
print(f"run_id {RUN_ID}")