FROM golang:1.22

WORKDIR /app

# Copy and download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN go build -o ingest .

EXPOSE 8080

CMD ["./ingest"]
