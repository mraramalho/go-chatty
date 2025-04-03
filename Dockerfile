FROM golang:1.24-alpine

WORKDIR /app

COPY . .

RUN  add

CMD ["./go-chatty"]
