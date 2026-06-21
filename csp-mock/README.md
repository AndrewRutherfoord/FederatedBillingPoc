# csp-mock

A mock Cloud Service Provider (CSP) API, built as part of a federated billing system prototype. It simulates the resource catalogue, customer registration, and billing account management that a real CSP would expose.

Billing data is emitted as [FinOps FOCUS 1.3](https://focus.finops.org/focus-specification/v1-3/) line items.

## Running

```bash
go run .  -config ./config.yaml
# Or from the root of the repo
go run ./csp-mock/main.go -config ./csp-mock/config.yaml

# To specify the database path (optional, defaults to csp-mock.sqlite in the working directory):
DB_PATH=./my-db.sqlite go run ./csp-mock/main.go -config ./csp-mock/config.yaml
```

The server starts on `:8081`. It expects `config.yaml` in the working directory and creates `csp-mock.sqlite` on first run.

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


## Project structure

```
csp-mock/
  config.yaml        # Provider, catalogue, and metering config
  main.go            # Startup: load config, open DB, start server
  internal/
    config/          # Config structs and YAML loader
    db/              # GORM models only
    repository/      # Data access interfaces and implementations
    handlers/        # Gin route handlers (methods on Server)
    middleware/      # Gin middleware (e.g. auth)
    util/            # Shared utilities (e.g. FOCUS line item builder)
```


## Swagger Documentation

The API endpoints are annotated with [Swaggo](https://github.com/swaggo/swag) comments. To generate the Swagger documentation, run:

```bash
$(go env GOPATH)/bin/swag init --parseInternal --parseDependency -g main.go
```

The generated docs will be available at `http://localhost:8081/swagger/index.html` when the server is running. This is a simple way to interact with the API and verify that the endpoints are working as expected.