FROM golang:alpine AS builder
RUN mkdir /godns
ADD . /godns/
WORKDIR /godns
RUN go build -o godns cmd/godns/godns.go

FROM alpine
RUN apk add --no-cache ca-certificates tzdata
RUN mkdir /usr/local/godns
COPY --from=builder /godns/godns /usr/local/godns
RUN chmod +x /usr/local/godns/godns
RUN rm -rf /var/cache/apk/*
WORKDIR /usr/local/godns
ENTRYPOINT ["./godns", "-c", "/usr/local/godns/config.json"]
