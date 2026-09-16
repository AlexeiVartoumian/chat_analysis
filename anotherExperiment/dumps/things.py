#!/usr/bin/env python3

import argparse
import os
import sys
import requests


def get_run(host, token, run_id, include=None):
    url = f"https://{host}/api/v2/runs/{run_id}"
    params = {"include": include} if include else {}
    headers = {
        "Authorization": f"Bearer {token}",
        "Content-Type": "application/vnd.api+json",
    }
    resp = requests.get(url, headers=headers, params=params)
    resp.raise_for_status()
    return resp.json()


def get_plan_json(host, token, plan_id):
    """Fetch the plan's json-output (same as `terraform show -json`)."""
    url = f"https://{host}/api/v2/plans/{plan_id}"
    headers = {
        "Authorization": f"Bearer {token}",
        "Content-Type": "application/vnd.api+json",
    }
    resp = requests.get(url, headers=headers)
    resp.raise_for_status()
    plan_data = resp.json()

    json_output_link = plan_data["data"]["links"].get("json-output")
    if not json_output_link:
        raise RuntimeError("No json-output link found on plan resource")

    # json-output link is relative; join with host
    full_url = f"https://{host}{json_output_link}"
    resp2 = requests.get(full_url, headers=headers)
    resp2.raise_for_status()
    return resp2.json()


def download_configuration_version(host, token, cv_id, out_path):
    url = f"https://{host}/api/v2/configuration-versions/{cv_id}/download"
    headers = {
        "Authorization": f"Bearer {token}",
        "Content-Type": "application/vnd.api+json",
    }
    resp = requests.get(url, headers=headers, stream=True)
    resp.raise_for_status()
    with open(out_path, "wb") as f:
        for chunk in resp.iter_content(chunk_size=8192):
            f.write(chunk)
    return out_path


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--host", required=True, help="TFE hostname, e.g. nameoftfe.com (no https://)")
    parser.add_argument("--run-id", required=True, help="Run ID, e.g. run-XXXXXXXX")
    parser.add_argument("--out-dir", default=".", help="Directory to write plan.json and config.tar.gz")
    args = parser.parse_args()

    token = os.environ.get("TFE_TOKEN")
    if not token:
        print("Error: set TFE_TOKEN environment variable with your API token", file=sys.stderr)
        sys.exit(1)

    os.makedirs(args.out_dir, exist_ok=True)

    print(f"Fetching run {args.run_id} with configuration_version included...")
    run_data = get_run(args.host, token, args.run_id, include="configuration_version")

    # Find the configuration-version in the included array
    cv_id = None
    for item in run_data.get("included", []):
        if item.get("type") == "configuration-versions":
            cv_id = item["id"]
            break

    if not cv_id:
        print("Warning: no configuration-version found in included data.", file=sys.stderr)
    else:
        print(f"Found configuration-version: {cv_id}")
        tar_path = os.path.join(args.out_dir, "config.tar.gz")
        download_configuration_version(args.host, token, cv_id, tar_path)
        print(f"Downloaded rendered config to {tar_path}")
        print(f"Unpack with: tar -xzf {tar_path} -C {args.out_dir}/rendered-source")

    # Also grab the plan id (relationship on the run) and its plan.json, for convenience
    plan_id = run_data["data"]["relationships"].get("plan", {}).get("data", {}).get("id")
    if plan_id:
        print(f"Found plan: {plan_id}, fetching json-output...")
        plan_json = get_plan_json(args.host, token, plan_id)
        import json
        plan_path = os.path.join(args.out_dir, "plan.json")
        with open(plan_path, "w") as f:
            json.dump(plan_json, f)
        print(f"Saved plan JSON to {plan_path}")
    else:
        print("Warning: no plan relationship found on run.", file=sys.stderr)


if __name__ == "__main__":
    main()