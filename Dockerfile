# Stage 1: Build the Next.js frontend
FROM node:20-alpine AS web-builder

# Install build dependencies
RUN apk add --no-cache libc6-compat

WORKDIR /app

# Copy package files for better caching
COPY web/package.json web/package-lock.json ./web/

# Change to web directory for npm operations  
WORKDIR /app/web

# Install dependencies with npm ci for faster, reliable installs  
RUN npm ci && npm cache clean --force

# Copy source code
COPY web/ .

# Build the Next.js project
RUN npm run build

# Stage 2: Build the Go backend
FROM golang:1.23-alpine AS go-builder

# Install build dependencies
RUN apk add --no-cache ca-certificates git tzdata

# Build arguments - VERSION is passed from CI, TARGETOS/TARGETARCH are automatic from buildx
ARG VERSION
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

# Copy go mod files for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY pkg/ ./pkg/

# Copy the Next.js build from the previous stage
COPY --from=web-builder /app/web/out ./web/out

# Generate embedded assets
RUN go generate ./...

# Build the Go binary with optimizations
RUN echo "Building for TARGETOS=${TARGETOS} TARGETARCH=${TARGETARCH} VERSION=${VERSION}" && \
    CGO_ENABLED=0 \
    GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH} \
    go build \
    -ldflags="-w -s -X main.Version=${VERSION}" \
    -a -installsuffix cgo \
    -o godns \
    ./cmd/godns

# Final stage: Create a minimal final image
FROM --platform=$TARGETOS/$TARGETARCH gcr.io/distroless/static-debian12:nonroot

# Copy ca-certificates and timezone data
COPY --from=go-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=go-builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the binary
COPY --from=go-builder /app/godns /usr/local/bin/godns

# Use non-root user for security
USER nonroot:nonroot

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD ["/usr/local/bin/godns", "-h"] || exit 1

EXPOSE 9000

ENTRYPOINT ["/usr/local/bin/godns"]