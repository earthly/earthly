#!/bin/sh
set -e

# This script should be run via eval $(assume-developer-role.sh)

test -n "$MFA_ARN" || (echo "echo MFA_ARN not set && exit 1" && exit 1)
test -n "$MFA_KEY" || (echo "echo MFA_KEY not set && exit 1" && exit 1)

offset="$1"
if [ -n "$offset" ]; then
    tokencode="$(oathtool -b --totp "$MFA_KEY" -w "$offset" | tail -n 1)"
else
    tokencode="$(oathtool -b --totp "$MFA_KEY")"
fi

role_path="$(mktemp)"

aws sts assume-role \
  --role-arn arn:aws:iam::404851345508:role/developer \
  --role-session-name $(date +%s) \
  --duration-seconds 3600 \
  --serial-number $MFA_ARN \
  --token-code "$tokencode" > "$role_path"

AWS_ACCESS_KEY_ID=$(jq -r '.Credentials.AccessKeyId' "$role_path")
AWS_SECRET_ACCESS_KEY=$(jq -r '.Credentials.SecretAccessKey' "$role_path")
AWS_SESSION_TOKEN=$(jq -r '.Credentials.SessionToken' "$role_path")

rm "$role_path"

echo "export AWS_ACCESS_KEY_ID=\"$AWS_ACCESS_KEY_ID\";"
echo "export AWS_SECRET_ACCESS_KEY=\"$AWS_SECRET_ACCESS_KEY\";"
echo "export AWS_SESSION_TOKEN=\"$AWS_SESSION_TOKEN\";"
