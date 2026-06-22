#!/bin/bash

AUTHORIZATION_TOKEN="16f833e7-a11f-46d0-a422-c790add6d623"
BILLING_ACCOUNT_ID="d79d8057-3d6e-4342-b157-8d30077c47d0"

echo "Authorization token: ${AUTHORIZATION_TOKEN}"
echo "Billing account ID: ${BILLING_ACCOUNT_ID}"


# Create a small VM 
curl -X 'POST' \
  'http://localhost:8080/resources' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer ${AUTHORIZATION_TOKEN}' \
  -H 'Content-Type: application/json' \
  -d '{
  "billing_account_id": "${BILLING_ACCOUNT_ID}",
  "resource_type": "compute.vm.small"
}'

# Create a medium VM
curl -X 'POST' \
  'http://localhost:8080/resources' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer ${AUTHORIZATION_TOKEN}' \
  -H 'Content-Type: application/json' \
  -d '{
  "billing_account_id": "${BILLING_ACCOUNT_ID}",
  "resource_type": "compute.vm.medium"
}'

curl -X 'POST' \
  'http://localhost:8080/resources' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer ${AUTHORIZATION_TOKEN}' \
  -H 'Content-Type: application/json' \
  -d '{
  "billing_account_id": "${BILLING_ACCOUNT_ID}",
  "resource_type": "storage.object.standard"
}'