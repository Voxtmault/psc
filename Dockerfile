FROM golang:1.24.3 as builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

# Build the application.
# -ldflags="-w -s" reduces the binary size.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /spektra ./cmd/spektra

# Stage 2: Create the final, lean image
FROM alpine:3.20

# Install ffmpeg and ca-certificates.
# ca-certificates is required for making secure HTTPS requests from within the container.
# Using --no-cache avoids storing the package index, keeping the image smaller.
RUN apk add --no-cache ffmpeg ca-certificates

WORKDIR /app

COPY --from=builder /spektra /usr/local/bin/spektra
COPY --from=builder /build/assets ./assets

# For better security and clarity, create a dedicated non-root user.
RUN addgroup -S spektra && adduser -S spektra -G spektra
USER spektra

# Expose the port the application runs on for documentation and tooling
# EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/spektra"]
