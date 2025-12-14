# Stage 1: Build the Next.js frontend
FROM node:20-alpine AS web-builder
RUN apk add --no-cache libc6-compat
WORKDIR /app
COPY web/package.json web/package-lock.json ./web/
WORKDIR /app/web
RUN npm ci && npm cache clean --force
COPY web/ .
RUN npm run build

# Stage 2: Build the Go backend
FROM golang:1.24.1-alpine AS go-builder
RUN apk add --no-cache ca-certificates tzdata
ARG VERSION
ARG TARGETOS
ARG TARGETARCH
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY pkg/ ./pkg/
COPY --from=web-builder /app/web/out ./web/out
RUN go generate ./...
RUN CGO_ENABLED=0 \
    GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH} \
    go build \
    -ldflags="-w -s -X main.Version=${VERSION}" \
    -a -installsuffix cgo \
    -o godns \
    ./cmd/godns

# Final stage: Minimal runtime image
FROM --platform=$TARGETOS/$TARGETARCH gcr.io/distroless/static-debian12:nonroot
COPY --from=go-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=go-builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=go-builder /app/godns /godns
USER nonroot:nonroot
WORKDIR /
ENTRYPOINT ["/godns"]
