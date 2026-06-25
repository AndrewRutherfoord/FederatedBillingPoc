#!/bin/bash
set -euo pipefail

ACCOUNT_INFO_FILE="${1:-}"

if [[ -z "$ACCOUNT_INFO_FILE" ]]; then
  echo "Usage: $0 <path-to-account-info.json>" >&2
  echo "  (the JSON file downloaded from the onboarding 'Download JSON' button)" >&2
  exit 1
fi

if [[ ! -f "$ACCOUNT_INFO_FILE" ]]; then
  echo "File not found: ${ACCOUNT_INFO_FILE}" >&2
  exit 1
fi

# The CSP customer_id doubles as the Authorization bearer token for /resources.
AUTHORIZATION_TOKEN=$(python3 -c "import json,sys; print(json.load(open(sys.argv[1]))['customer_id'])" "$ACCOUNT_INFO_FILE")
BILLING_ACCOUNT_ID=$(python3 -c "import json,sys; print(json.load(open(sys.argv[1]))['billing_account_id'])" "$ACCOUNT_INFO_FILE")
CSP_HOST=$(python3 -c "import json,sys; print(json.load(open(sys.argv[1]))['host'])" "$ACCOUNT_INFO_FILE")

if [[ -z "$AUTHORIZATION_TOKEN" || -z "$BILLING_ACCOUNT_ID" || -z "$CSP_HOST" ]]; then
  echo "Could not read customer_id / billing_account_id / host from ${ACCOUNT_INFO_FILE}" >&2
  exit 1
fi

echo "CSP host: ${CSP_HOST}"
echo "Authorization token (customer ID): ${AUTHORIZATION_TOKEN}"
echo "Billing account ID: ${BILLING_ACCOUNT_ID}"


# # Create a small VM
# curl -X 'POST' \
#   "${CSP_HOST}/resources" \
#   -H 'accept: application/json' \
#   -H "Authorization: Bearer ${AUTHORIZATION_TOKEN}" \
#   -H 'Content-Type: application/json' \
#   -d "{
#   \"billing_account_id\": \"${BILLING_ACCOUNT_ID}\",
#   \"resource_type\": \"compute.vm.small\"
# }"

# Create a medium VM
curl -X 'POST' \
  "${CSP_HOST}/resources" \
  -H 'accept: application/json' \
  -H "Authorization: Bearer ${AUTHORIZATION_TOKEN}" \
  -H 'Content-Type: application/json' \
  -d "{
  \"billing_account_id\": \"${BILLING_ACCOUNT_ID}\",
  \"resource_type\": \"compute.vm.medium\"
}"

# curl -X 'POST' \
#   "${CSP_HOST}/resources" \
#   -H 'accept: application/json' \
#   -H "Authorization: Bearer ${AUTHORIZATION_TOKEN}" \
#   -H 'Content-Type: application/json' \
#   -d "{
#   \"billing_account_id\": \"${BILLING_ACCOUNT_ID}\",
#   \"resource_type\": \"storage.object.standard\"
# }"
