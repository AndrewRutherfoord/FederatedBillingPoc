# Federated Billing PoC

Proof-of-concept implementation for my MSc thesis on federated cloud billing. It simulates a small ecosystem where multiple Cloud Service Providers (CSPs) meter resource usage, a Billing Provider aggregates and validates that usage across CSPs, and a Customer Billing service turns it into invoices for the end customer.

Usage data is exchanged as [FinOps FOCUS 1.3](https://focus.finops.org/focus-specification/v1-3/) line items. Charge batches passed from a CSP to the Billing Provider are hashed into a Merkle tree so the Billing Provider can verify their integrity.

## Architecture

This is the same C4 container diagram used in the thesis, updated with the extra pieces the PoC actually has: the Billing Provider split into its two real APIs, and a "Billing Control Plane" boundary made up of the customer-billing service and the frontend that sits on top of it.

```plantuml
@startuml
!includeurl https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Container.puml

title Federated Billing PoC - System Diagram

Person(customer, "Customer", "Uses the frontend to onboard CSP accounts, link them to a billing provider, and view usage and invoices.")

System_Boundary(cspBoundary, "Cloud Service Provider (x2 in this PoC)") {
    System(cspCustomerApi, "CSP Customer API", "Resource catalogue, provisioning, and billing account onboarding.")
    System(cspBillingAdapter, "CSP Billing Adapter", "Meters usage, builds FOCUS line items, and batches them with a Merkle root for integrity.")

    Rel(cspCustomerApi, cspBillingAdapter, "Metered resource usage")
}

System_Boundary(bpBoundary, "Billing Provider") {
    System(bpCspApi, "Billing Provider - CSP API", "Pulls charge batches from each CSP over mTLS and verifies their Merkle roots.")
    System(bpCustomerApi, "Billing Provider - Customer API", "Pricing, aggregation, invoicing, and credit, using verified charge batches.")

    Rel(bpCspApi, bpCustomerApi, "Verified charge batches")
}

System_Boundary(bcpBoundary, "Billing Control Plane") {
    System(customerBilling, "Customer Billing", "Collects and correlates usage and billing data from CSPs and the billing provider for unified visibility.")
    System(frontend, "Frontend", "Vue SPA for onboarding, and browsing charge batches, resource charges, and invoices.")

    Rel(frontend, customerBilling, "Reads usage, batches, and invoices")
}

Rel(customer, cspCustomerApi, "Provisions resources, onboards billing account")
Rel(customer, frontend, "Views analytics & insights, controls data sharing scope")

Rel(bpCspApi, cspBillingAdapter, "Retrieves charge batches (FOCUS, hash & Merkle root)", "mTLS")
Rel(customerBilling, cspBillingAdapter, "Retrieves usage & cost data for cross-verification (hash & Merkle root)")
Rel(customerBilling, bpCustomerApi, "Retrieves invoices and charge batches")

SHOW_LEGEND()
@enduml
```

The mock clock (time-server) is left out of the diagram since it is a testing utility rather than a conceptual part of the architecture.

- **csp-mock** - a mock CSP. Meters simulated resource usage, exposes a customer-facing API and a billing-provider-facing API (mTLS). Two instances run side by side to demonstrate multi-CSP support.
- **billing-provider** - sits between CSPs and the customer-facing billing service. Pulls charge batches from each CSP over mTLS, verifies them, and exposes its own customer-facing API.
- **customer-billing** - the customer-facing billing service. Talks to one or more billing providers and presents consolidated invoices.
- **frontend** - Vue 3 SPA used to onboard CSPs/billing providers and browse charge batches, resource charges, and invoices.
- **shared** - common Go module: mTLS client, Merkle tree/canonical hashing, models, and the mock-clock client.
- **time-server** - a mock clock. Every service reads simulated time from here so usage, billing, and settlement periods can be fast-forwarded instead of waiting on real time.

## Getting started

### 1. Generate mTLS certificates

The CSPs and the billing provider authenticate to each other with mTLS. Certificates are not committed, so generate them first:

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

### 2. Start the backend

```bash
make up
# or: docker compose up -d
```

This builds and starts the time server, both mock CSPs (customer API + billing-provider interface), the billing provider (customer API + CSP-facing API), and the customer billing service, each with its own SQLite database under `./dbs`.

Useful ports once everything is up:

| Service | Port |
| --- | --- |
| time-server | 9999 |
| csp-mock (mock-cloud) customer API | 8080 |
| csp-mock (mock-cloud) billing-provider interface | 8443 |
| csp-mock (imaginary-cloud) customer API | 8090 |
| csp-mock (imaginary-cloud) billing-provider interface | 8445 |
| billing-provider customer API | 8081 |
| billing-provider CSP-facing API | 8444 |
| customer-billing API | 8082 |

Tail logs with `make logs`.

### 3. Run the frontend

```bash
cd frontend
pnpm install
pnpm dev
```

### 4. Advance the mock clock

Since billing periods and settlement cycles happen on simulated time, use the time server to fast-forward instead of waiting:

```bash
make step-hour
# or open http://localhost:9999 for a simple web UI with time-advance buttons
```

## Project layout

```
billing-provider/   Billing Provider service (CSP-facing API + customer-facing API)
csp-mock/           Mock CSP service, run twice to simulate two independent providers
customer-billing/   Customer-facing billing/invoicing service
frontend/           Vue 3 SPA
shared/             Shared Go module (mTLS client, Merkle tree, models, scheduler)
time-server/        Mock clock used by every other service
keys/               Generated mTLS certificates (not committed)
dbs/                SQLite databases created at runtime (not committed)
```

See [csp-mock/README.md](csp-mock/README.md) for details on the mock CSP's configuration and Swagger docs.
