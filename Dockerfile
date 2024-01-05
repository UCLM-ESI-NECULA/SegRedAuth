FROM golang:alpine
WORKDIR /app
COPY . .
RUN go build -o main .
EXPOSE 8080
CMD ["/cmd/app/main"]