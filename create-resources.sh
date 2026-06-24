#!/bin/bash

AUTHORIZATION_TOKEN="7ca48d95-dab3-45d5-92c5-0ab888636f89"
BILLING_ACCOUNT_ID="e2046d78-3734-4118-932d-ca3325b61e9d"

echo "Authorization token: ${AUTHORIZATION_TOKEN}"
echo "Billing account ID: ${BILLING_ACCOUNT_ID}"


# Create a small VM
curl -X 'POST' \
  'http://localhost:8080/resources' \
  -H 'accept: application/json' \
  -H "Authorization: Bearer ${AUTHORIZATION_TOKEN}" \
  -H 'Content-Type: application/json' \
  -d "{
  \"billing_account_id\": \"${BILLING_ACCOUNT_ID}\",
  \"resource_type\": \"compute.vm.small\"
}"

# Create a medium VM
curl -X 'POST' \
  'http://localhost:8080/resources' \
  -H 'accept: application/json' \
  -H "Authorization: Bearer ${AUTHORIZATION_TOKEN}" \
  -H 'Content-Type: application/json' \
  -d "{
  \"billing_account_id\": \"${BILLING_ACCOUNT_ID}\",
  \"resource_type\": \"compute.vm.medium\"
}"

curl -X 'POST' \
  'http://localhost:8080/resources' \
  -H 'accept: application/json' \
  -H "Authorization: Bearer ${AUTHORIZATION_TOKEN}" \
  -H 'Content-Type: application/json' \
  -d "{
  \"billing_account_id\": \"${BILLING_ACCOUNT_ID}\",
  \"resource_type\": \"storage.object.standard\"
}"