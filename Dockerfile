FROM golang:1.23.5 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o urlshortener .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/urlshortener .

RUN chmod +x /root/urlshortener

EXPOSE 8080

CMD ["./urlshortener"]
