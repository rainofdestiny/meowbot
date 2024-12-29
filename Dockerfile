FROM golang:1.20

WORKDIR /app

COPY . .

COPY .env .env

RUN go mod tidy

RUN go build -o bot .

CMD ["./bot"]
