# Stage 1: Build the application
FROM golang:1.25.6 AS build

WORKDIR /app

# Only copy go.mod and go.sum to leverage Docker caching
COPY go.mod go.sum ./

# Download dependencies early to avoid re-downloading if code changes
RUN go mod download

# Copy the remaining source code
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o notification-service

# Stage 2: Create the minimal runtime image
FROM alpine:latest

# Add CA certificates to enable HTTPS requests
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /app

# Copy the build binary from the build stage
COPY --from=build /app/notification-app /app/notification-service

# Copy the .env file (optional, depending on your app's configuration)
COPY .env /app/.env

# Expose the desired port
EXPOSE 8080

# Command to run the app
CMD ["./notification-app"]