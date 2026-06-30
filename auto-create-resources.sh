#!/bin/bash
set -euo pipefail

DOWNLOADS_DIR="/mnt/c/Users/rutherfoordam/Downloads/"
MAX_AGE_SECONDS=600 # 10 minutes

# Onboarding "Download JSON" files are named <provider_id>_<billing_account_uuid>.json.
# Grab the 2 most recently downloaded ones (one per mock CSP) and make sure
# they actually belong together before creating any resources.
mapfile -t CANDIDATE_FILES < <(ls -t "${DOWNLOADS_DIR}"mock-cloud_*.json "${DOWNLOADS_DIR}"imaginary-cloud_*.json 2>/dev/null)

if [[ ${#CANDIDATE_FILES[@]} -lt 2 ]]; then
  echo "Could not find onboarding JSON files for both CSPs in ${DOWNLOADS_DIR}" >&2
  exit 1
fi

FILE_A="${CANDIDATE_FILES[0]}"
FILE_B="${CANDIDATE_FILES[1]}"

UUID_A=$(basename "$FILE_A" | sed -E 's/^(mock-cloud|imaginary-cloud)_(.+)\.json$/\2/')
UUID_B=$(basename "$FILE_B" | sed -E 's/^(mock-cloud|imaginary-cloud)_(.+)\.json$/\2/')

if [[ -z "$UUID_A" || -z "$UUID_B" || "$UUID_A" != "$UUID_B" ]]; then
  echo "Could not find a matching pair of CSP onboarding files with the same billing account ID" >&2
  exit 1
fi

NOW=$(date +%s)
for FILE in "$FILE_A" "$FILE_B"; do
  FILE_AGE=$(( NOW - $(stat -c %Y "$FILE") ))
  if (( FILE_AGE > MAX_AGE_SECONDS )); then
    echo "Could not find onboarding files newer than 10 minutes old (found ${FILE})" >&2
    exit 1
  fi
done

BILLING_ACCOUNT_UUID="$UUID_A"
echo "Billing account UUID: ${BILLING_ACCOUNT_UUID}"

CREATED_RESOURCES=() # entries: "provider_id|resource_type|resource_id"

create_resource() {
  local csp_host="$1" auth_token="$2" billing_account_id="$3" resource_type="$4" provider_id="$5"
  local response resource_id
  response=$(curl -s -X 'POST' \
    "${csp_host}/resources" \
    -H 'accept: application/json' \
    -H "Authorization: Bearer ${auth_token}" \
    -H 'Content-Type: application/json' \
    -d "{
    \"billing_account_id\": \"${billing_account_id}\",
    \"resource_type\": \"${resource_type}\"
  }")
  echo "$response"
  resource_id=$(python3 -c "import json,sys; print(json.load(sys.stdin).get('id','unknown'))" <<< "$response" 2>/dev/null || echo "unknown")
  CREATED_RESOURCES+=("${provider_id}|${resource_type}|${resource_id}")
}

for ACCOUNT_INFO_FILE in "$FILE_A" "$FILE_B"; do
  echo "--- Processing ${ACCOUNT_INFO_FILE} ---"

  # The CSP customer_id doubles as the Authorization bearer token for /resources.
  PROVIDER_ID=$(python3 -c "import json,sys; print(json.load(open(sys.argv[1]))['provider_id'])" "$ACCOUNT_INFO_FILE")
  AUTHORIZATION_TOKEN=$(python3 -c "import json,sys; print(json.load(open(sys.argv[1]))['customer_id'])" "$ACCOUNT_INFO_FILE")
  BILLING_ACCOUNT_ID=$(python3 -c "import json,sys; print(json.load(open(sys.argv[1]))['billing_account_id'])" "$ACCOUNT_INFO_FILE")
  CSP_HOST=$(python3 -c "import json,sys; print(json.load(open(sys.argv[1]))['host'])" "$ACCOUNT_INFO_FILE")

  echo "Provider: ${PROVIDER_ID}"
  echo "CSP host: ${CSP_HOST}"
  echo "Authorization token (customer ID): ${AUTHORIZATION_TOKEN}"
  echo "Billing account ID: ${BILLING_ACCOUNT_ID}"

  if [[ "$PROVIDER_ID" == "mock-cloud" ]]; then
    create_resource "$CSP_HOST" "$AUTHORIZATION_TOKEN" "$BILLING_ACCOUNT_ID" "compute.vm.medium" "$PROVIDER_ID"
    create_resource "$CSP_HOST" "$AUTHORIZATION_TOKEN" "$BILLING_ACCOUNT_ID" "storage.object.standard" "$PROVIDER_ID"
  elif [[ "$PROVIDER_ID" == "imaginary-cloud" ]]; then
    create_resource "$CSP_HOST" "$AUTHORIZATION_TOKEN" "$BILLING_ACCOUNT_ID" "compute.vm.small" "$PROVIDER_ID"
    create_resource "$CSP_HOST" "$AUTHORIZATION_TOKEN" "$BILLING_ACCOUNT_ID" "database.relational.small" "$PROVIDER_ID"
  else
    echo "Unknown provider_id '${PROVIDER_ID}', skipping resource creation" >&2
  fi
done

echo ""
echo "=== Summary ==="
echo "Billing account: ${BILLING_ACCOUNT_UUID}"
printf "%-20s %-35s %-36s\n" "Provider" "Resource Type" "Resource ID"
printf "%-20s %-35s %-36s\n" "--------------------" "-----------------------------------" "------------------------------------"
for entry in "${CREATED_RESOURCES[@]}"; do
  IFS='|' read -r provider resource_type rid <<< "$entry"
  printf "%-20s %-35s %-36s\n" "$provider" "$resource_type" "$rid"
done
echo "${#CREATED_RESOURCES[@]} resource(s) created"
