# Build stage
FROM golang:latest AS builder

WORKDIR /app

RUN apt-get update && apt-get install -y git && rm -rf /var/lib/apt/lists/*

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o time-server .

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/time-server .
COPY time-state.json .

EXPOSE 9999

CMD ["./time-server"]
