# Start from a GoLang base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files to the container
COPY go.mod go.sum ./

# Download and cache Go modules
RUN go mod download

# Copy the rest of the application source code to the container
COPY . .

# Build the GoLang service binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o elotus-service

# Expose the port on which the GoLang service listens
EXPOSE 8080

# Set the entry point for the container
ENTRYPOINT ["/app/elotus-service"]
