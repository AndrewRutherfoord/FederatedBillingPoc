# Build stage
FROM golang:latest AS builder

WORKDIR /app

RUN apt-get update && apt-get install -y git && rm -rf /var/lib/apt/lists/*

COPY csp-mock/go.mod csp-mock/go.sum ./
COPY shared /shared
RUN go mod download

COPY csp-mock/ .

# Target 1: Customer API
FROM builder AS customer-api
RUN CGO_ENABLED=1 go build -o customer-api ./cmd/customer-api

# Target 2: Billing Provider Interface
FROM builder AS billing-provider-interface
RUN CGO_ENABLED=1 go build -o billing-provider-interface ./cmd/billing-provider-interface

# Runtime stage for customer-api
FROM alpine:latest AS runtime-customer-api
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=customer-api /app/customer-api .
EXPOSE 8080
CMD ["./customer-api"]

# Runtime stage for billing-provider-interface
FROM alpine:latest AS runtime-billing-provider-interface
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=billing-provider-interface /app/billing-provider-interface .
EXPOSE 8443
CMD ["./billing-provider-interface"]
