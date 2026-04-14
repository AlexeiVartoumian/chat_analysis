#!/bin/bash



WORKSPACE_DIRECTORY=$(pwd)
PROVIDER="provider"

TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S')

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