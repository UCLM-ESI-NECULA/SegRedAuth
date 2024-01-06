FROM midebian

# Install Go and other build dependencies
RUN apt-get update && \
    apt-get install -y \
    golang-go \
    git

# Set up the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./
COPY go.sum ./

# Download all dependencies
RUN go mod download

COPY . .
COPY ./certs/mycert.crt ./certs/mycert.key ./

# Build the app
RUN go build -o /app/main ./cmd/app

# Expose the port of the service and the SSH port
EXPOSE 8080 22

# Copy the entrypoint script
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
