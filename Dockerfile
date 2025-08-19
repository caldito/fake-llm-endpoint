# Use Go 1.25.0 for building
FROM golang:1.25.0-alpine AS builder

WORKDIR /app

# Cache dependencies first
COPY go.mod ./
#COPY go.sum ./ #no needed for now as there are no dependencies
RUN go mod download

# Copy source
COPY . .

# Build statically
RUN go build -o fake-llm-endpoint -v cmd/fake-llm-endpoint/main.go

# Minimal runtime image
FROM alpine:3.19
RUN addgroup -S app && adduser -S app -G app

WORKDIR /app
COPY --from=builder /app/fake-llm-endpoint .

USER app

EXPOSE 8080

CMD ["./fake-llm-endpoint"]

