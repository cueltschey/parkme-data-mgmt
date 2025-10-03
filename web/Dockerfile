FROM golang:1.24.6 AS builder

WORKDIR /app

COPY . .

RUN go mod init github.com/cueltschey/parkme-data-mgmt || true

RUN go mod tidy

RUN go get github.com/gin-gonic/gin

RUN CGO_ENABLED=0 GOOS=linux go build

FROM alpine:latest

WORKDIR /

COPY --from=builder /app/ /usr/bin

EXPOSE 3000

CMD ["/usr/bin/surety-db", "--port", "3000", "--static", "/web", "--db", "/main.db"]
