FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o distributed-calculator cmd/distributed-calculator/main.go

EXPOSE 8080

CMD ["./distributed-calculator"]