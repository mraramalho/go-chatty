FROM golang:1.24-alpine

WORKDIR /app

COPY . .

RUN  go build -o go-chatty ./cmd/socket/main.go

CMD ["./go-chatty"]
