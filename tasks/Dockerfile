FROM golang:1.23.2-alpine AS builder

COPY . /github.com/Gustcat/tasks/source/
WORKDIR /github.com/Gustcat/tasks/source/

RUN go mod download
RUN go build -o ./bin/crud_server cmd/http_server/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/Gustcat/tasks/source/bin/crud_server .
COPY --from=builder /github.com/Gustcat/tasks/source/.env .

CMD ["./crud_server"]