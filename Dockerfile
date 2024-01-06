FROM golang:alpine

WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./
COPY go.sum ./

# Download all dependencies.
RUN go mod download

COPY . .

COPY ./certs/mycert.crt ./certs/mycert.key ./

# Build the app
RUN go build -o main ./cmd/app

CMD ["/app/main"]