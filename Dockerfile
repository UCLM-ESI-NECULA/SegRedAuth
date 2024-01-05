FROM golang:alpine

WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./
COPY go.sum ./

# Download all dependencies.
RUN go mod download

COPY . .

# Build the app
RUN go build -o main ./cmd/app

# Expose port
EXPOSE 8080

CMD ["/app/main"]