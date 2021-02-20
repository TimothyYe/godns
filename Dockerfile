FROM golang:alpine AS builder
RUN mkdir /godns
ADD . /godns/
WORKDIR /godns
RUN CGO_ENABLED=0 go build -o godns cmd/godns/godns.go

FROM gcr.io/distroless/base
COPY --from=builder /godns/godns /godns
ENTRYPOINT ["/godns"]