FROM golang:1.22.0-alpine

WORKDIR /app

COPY . .

RUN go build -o app ./cmd/quotation/main.go

EXPOSE 8080
EXPOSE 8082

CMD ["./app"]
