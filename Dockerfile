FROM golang:1.21-alpine

WORKDIR /app

COPY chat.go .

RUN go build -o go-chatty ./cmd/socket/main.go

CMD ["./go-chatty"]
