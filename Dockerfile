# Stage 1: Build the Next.js frontend
FROM node:18-alpine AS web-builder
WORKDIR /web
# Copy the Next.js project files into the image
COPY ./web/package.json ./web/package-lock.json ./
# Install dependencies
RUN npm install
# Copy the rest of the Next.js project files
COPY ./web .
# Build the Next.js project
RUN npm run build

# Stage 2: Build the Go backend
FROM golang:alpine AS builder
WORKDIR /godns
ADD . .
# Copy the Next.js build from the previous stage
COPY --from=web-builder /web/out ./web/out
RUN go generate ./...
RUN CGO_ENABLED=0 go build -o godns cmd/godns/godns.go

FROM gcr.io/distroless/base
COPY --from=builder /godns/godns /godns
ENTRYPOINT ["/godns"]