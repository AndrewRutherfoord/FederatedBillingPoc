# csp-mock

A mock Cloud Service Provider (CSP) API, built as part of a federated billing system prototype. It simulates the resource catalogue, customer registration, and billing account management that a real CSP would expose.

Billing data is emitted as [FinOps FOCUS 1.3](https://focus.finops.org/focus-specification/v1-3/) line items.

## Running

```bash
go run . 
# or from the workspace root:
go run ./csp-mock/
```

The server starts on `:8080`. It expects `config.yaml` in the working directory and creates `csp-mock.sqlite` on first run.

## Configuration

The provider identity, resource catalogue, billing providers, and metering settings are all defined in `config.yaml`.

```yaml
provider_id: "mock-csp-1"
provider_name: "MockCloud EU-West"
currency: "EUR"
region_id: "eu-west-1"
region_name: "EU West (Amsterdam)"

resource_types:
  - id: "compute.vm.small"
    display_name: "Virtual Machine — Small"
    service_category: "Compute"          # FOCUS 1.3 ServiceCategory
    service_subcategory: "Virtual Machines"
    sku_id: "MC-VM-S"
    billing_unit: "Hours"
    pricing:
      model: "per_hour"
      unit_price: 0.042
    ...

billing_providers:
  - id: "billing-provider-1"
    name: "Mock Billing Provider"
    api_endpoint: "http://..."

metering:
  tick_interval_seconds: 60
  simulated_hours_per_tick: 1
```

`service_category` and `service_subcategory` accept human-readable title-case strings (e.g. `"Virtual Machines"`); they are normalised to FOCUS snake_case values at load time.

## API

### Public endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/health` | Liveness check. Returns provider ID. |
| `GET` | `/resource-types` | List all resource types from the catalogue. |
| `GET` | `/resource-types/:id` | Get a single resource type by ID. |
| `POST` | `/account/register` | Register a new customer. |

#### Register a customer

```
POST /account/register
Content-Type: application/json

{
  "name": "Acme Corp",
  "email": "billing@acme.example"
}
```

Returns the created customer record including its generated `id`. Use this `id` when creating billing accounts.

### Authenticated endpoints

Authenticated routes require a billing account ID in the `Authorization` header:

```
Authorization: Bearer <account_id>
```

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/account` | Returns the billing account associated with the token. |

## Project structure

```
csp-mock/
  config.yaml                    # Provider, catalogue, and metering config
  main.go                        # Startup: load config, open DB, start server
  internal/
    config/       config.go      # Config structs and YAML loader
    db/                          # GORM models only
      db.go                      # Opens SQLite, runs AutoMigrate
      customer.go
      billing_account.go
      focus_record.go            # Embeds shared-models FocusLineItem
    repository/                  # Data access interfaces and implementations
      repos.go                   # Repos aggregate passed to handlers
      customers.go               # DB-backed
      billing_accounts.go        # DB-backed
      billing_providers.go       # Config-backed
      resource_types.go          # Config-backed
      focus.go                   # DB-backed (FOCUS line items)
    handlers/                    # Gin route handlers (methods on Server)
      server.go                  # Server struct holding *Repos
      routes.go                  # Route registration
      health.go
      catalog.go
      customer.go
      account.go
    middleware/
      auth.go                    # Bearer token → BillingAccount lookup
```

## Data model

- **Customer** — a registered end-user (name, email, UUID).
- **BillingAccount** — links a customer to a billing provider using the provider-assigned `account_id`. Unique per `(billing_provider_id, account_id)`.
- **FocusRecord** — a persisted FOCUS 1.3 line item. Wraps the `FocusLineItem` type from the `shared-models` module.