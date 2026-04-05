#!/bin/bash

HUB_ACCOUNT_ID=$1
SPOKE_ACCOUNT_ID=$2
ROLE_NAME="AWSControlTowerExecution"

# Step 1: Assume hub role
export $(aws sts assume-role \
  --role-arn "arn:aws:iam::${HUB_ACCOUNT_ID}:role/${ROLE_NAME}" \
  --role-session-name "hub-session" | \
jq -r '.Credentials | "AWS_ACCESS_KEY_ID=\(.AccessKeyId) AWS_SECRET_ACCESS_KEY=\(.SecretAccessKey) AWS_SESSION_TOKEN=\(.SessionToken)"')

# Step 2: Assume spoke role using those creds
export $(aws sts assume-role \
  --role-arn "arn:aws:iam::${SPOKE_ACCOUNT_ID}:role/${ROLE_NAME}" \
  --role-session-name "spoke-session" | \
jq -r '.Credentials | "AWS_ACCESS_KEY_ID=\(.AccessKeyId) AWS_SECRET_ACCESS_KEY=\(.SecretAccessKey) AWS_SESSION_TOKEN=\(.SessionToken)"')


#unset AWS_ACCESS_KEY_ID AWS_SECRET_ACCESS_KEY AWS_SESSION_TOKEN

