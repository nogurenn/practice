FROM golang:1.23-alpine3.21 AS test

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY mocks/ ./mocks/

CMD ["go", "test", "-v", "./internal/..."]