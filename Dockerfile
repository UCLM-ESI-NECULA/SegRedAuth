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
COPY ./certs/auth.crt ./certs/auth.key ./

# Build the app
RUN go build -o /app/main ./cmd/app

# Copy the entrypoint script
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh


ENTRYPOINT ["/entrypoint.sh"]
