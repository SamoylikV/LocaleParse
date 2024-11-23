FROM golang:1.23.3

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o locale-parser ./cmd/main.go

EXPOSE 8083

CMD ["./locale-parser"]