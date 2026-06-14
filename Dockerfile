# ---- Builder Stage ----
# Use a specific version for reproducibility
FROM golang:1.23.4-alpine AS builder

# 1. Install necessary build tools
RUN apk add --no-cache gcc musl-dev git

WORKDIR /app

# 2. Set Go proxy (optional but good practice)
RUN go env -w GOPROXY=https://goproxy.cn,https://goproxy.io,off && go env -w GONOSUMDB="*" && go env -w GOFLAGS=-mod=mod

# 3. Copy only dependency files first to leverage Docker cache
COPY go.mod go.sum ./

# 4. Download dependencies
RUN go mod download

# 5. Install the Wire code generator
RUN go install github.com/google/wire/cmd/wire@latest

# 6. Copy all your source code
COPY . .

# 7. IMPORTANT: Run wire to generate the necessary Go files
RUN wire gen ./...

# 8. Build the application, placing the binary in the root for easy access
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /main ./cmd/app


# ---- Final Stage ----
# Use a minimal image for a small and secure final container
FROM alpine:latest

# Copy certificates from the builder stage
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy static assets your app needs
COPY --from=builder /app/internal/infrastructure/communication/email/templates/ /templates/

COPY --from=builder /app/internal/infrastructure/jwt/ /app/internal/infrastructure/jwt/

WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /main .

# Expose the port your application listens on
EXPOSE 8080

# The command to run your application
CMD ["./main"]