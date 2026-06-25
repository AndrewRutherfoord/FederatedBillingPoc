


# Federated Billing PoC

## Generating Certificates

```bash
openssl req -x509 -newkey rsa:4096 -days 365 -nodes \
  -keyout ./keys/csp-1.key -out ./keys/csp-1.crt \
  -subj "/CN=csp-adapter" \
  -addext "subjectAltName=DNS:localhost,DNS:csp-mock-adapter,IP:127.0.0.1"

openssl req -x509 -newkey rsa:4096 -days 365 -nodes \
  -keyout ./keys/bp-1.key -out ./keys/bp-1.crt \
  -subj "/CN=billing-provider-1" \
  -addext "subjectAltName=DNS:localhost,DNS:billing-provider-csp-api,IP:127.0.0.1"

# Second mock CSP, used to demonstrate the system handling multiple CSPs.
openssl req -x509 -newkey rsa:4096 -days 365 -nodes \
  -keyout ./keys/csp-2.key -out ./keys/csp-2.crt \
  -subj "/CN=csp-2-adapter" \
  -addext "subjectAltName=DNS:localhost,DNS:csp-mock-2-adapter,IP:127.0.0.1"
```

The Docker service hostnames (`csp-mock-adapter`, `billing-provider-csp-api`) must be included as SANs so that mTLS connections between containers pass certificate verification.