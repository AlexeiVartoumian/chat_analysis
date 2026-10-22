#!/usr/bin/env bash
# Assume an IAM role via the AWS CLI, then run commands as that role.
# Requirements: aws CLI v2 (works in Git Bash, no jq needed)
# Usage: ./assume-role.sh arn:aws:iam::123456789012:role/MyRole [session-name]

set -euo pipefail

ROLE_ARN="${1:?Usage: $0 <role-arn> [session-name]}"
SESSION_NAME="${2:-cli-session-$(date +%s)}"
DURATION=3600   # seconds (max depends on the role's MaxSessionDuration)

echo "Current identity:"
aws sts get-caller-identity --output table

echo "Assuming role: $ROLE_ARN"
CREDS=$(aws sts assume-role \
  --role-arn "$ROLE_ARN" \
  --role-session-name "$SESSION_NAME" \
  --duration-seconds "$DURATION" \
  --query 'Credentials.[AccessKeyId,SecretAccessKey,SessionToken,Expiration]' \
  --output text | tr -d '\r')
# If the role requires MFA, add:
#   --serial-number arn:aws:iam::<acct>:mfa/<user> --token-code <6-digit-code>
# If it requires an external ID, add:
#   --external-id <id>

read -r AWS_ACCESS_KEY_ID AWS_SECRET_ACCESS_KEY AWS_SESSION_TOKEN EXPIRY <<< "$CREDS"
export AWS_ACCESS_KEY_ID AWS_SECRET_ACCESS_KEY AWS_SESSION_TOKEN

echo "Role assumed. Credentials expire at: $EXPIRY"
echo "New identity:"
aws sts get-caller-identity --output table

# ---------------------------------------------------------------
# Your actions as the assumed role go here
# ---------------------------------------------------------------
echo "Listing S3 buckets as assumed role..."
aws s3 ls

# e.g.
# aws ec2 describe-instances --region eu-west-2 --output table
# aws s3 cp ./file.txt s3://my-bucket/

# ---------------------------------------------------------------
# Drop back to your original credentials (optional)
# ---------------------------------------------------------------
unset AWS_ACCESS_KEY_ID AWS_SECRET_ACCESS_KEY AWS_SESSION_TOKEN
echo "Reverted to original identity."