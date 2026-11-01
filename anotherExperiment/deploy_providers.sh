#!/bin/bash

set -e


WORKSPACE_DIRECTORY=$(pwd)
PROVIDER="provider"

TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S')


case "${1:-}" in
  on)  echo "polling_enabled = true"  > polling.auto.tfvars ;;
  off) echo "polling_enabled = false" > polling.auto.tfvars ;;
  "")  ;;   # no argument: keep the current setting
  *)   echo "usage: $0 [on|off]"; exit 1 ;;
esac

# if ls *.jinja > /dev/null 2>&1; then
#     for f in *.jinja; do
#         jinja2 $f vars.yaml  >> "./$(basename $f .jinja).tf"; 
#     done
# fi

for f in *.jinja; do
    # inject timestamp into vars.yaml dynamically
    sed "s/timestamp:.*/timestamp: \"$TIMESTAMP\"/" vars.yaml > vars_tmp.yaml
    jinja2 $f vars_tmp.yaml > "./$(basename $f .jinja).tf"
done
terraform init
terraform plan
terraform apply