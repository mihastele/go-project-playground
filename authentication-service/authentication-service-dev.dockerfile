# base
FROM golang:1.18-alpine AS builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN go get github.com/go-delve/delve/cmd/dlv
RUN go install github.com/go-delve/delve/cmd/dlv

RUN CGO_ENABLED=0 go build -gcflags "all=-N -l" -o authApp ./cmd/api

RUN chmod +rwx /app/authApp

RUN chmod +x /go/bin/dlv

# Copy binary to debian
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/authApp /app
COPY --from=builder /go/bin/dlv /app
#set entrypoints
CMD ["/app/dlv","--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/app/authApp"]