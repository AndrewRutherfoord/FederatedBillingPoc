


# Federated Billing PoC

## Generating Certificates

```bash
openssl req -x509 -newkey rsa:4096 -keyout ./keys/csp-1.key -out ./keys/csp-1.crt -days 365 -nodes   -subj "/CN=csp-adapter"   -addext "subjectAltName=DNS:localhost,IP:127.0.0.1"
openssl req -x509 -newkey rsa:4096 -keyout ./keys/bp-1.key -out ./keys/bp-1.crt -days 365 -nodes   -subj "/CN=billing-provider-1"  -addext "subjectAltName=DNS:localhost,IP:127.0.0.1"
```