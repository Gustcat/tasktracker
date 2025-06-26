FROM golang:1.23.2-alpine AS builder

COPY . /github.com/Gustcat/auth/source/
WORKDIR /github.com/Gustcat/auth/source/

RUN go mod download
RUN go build -o ./bin/crud_server cmd/grpc_server/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/Gustcat/auth/source/bin/crud_server .
COPY --from=builder /github.com/Gustcat/auth/source/local.env .

CMD ["./crud_server", "-config-path", "local.env" ]