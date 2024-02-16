FROM golang:1.22

WORKDIR /app

ENV MAX_GOROUTINES=15
ENV POSTGRES_HOST=localhost
ENV POSTGRES_PORT=5432
ENV POSTGRES_USER=admin
ENV POSTGRES_PASSWORD=password
ENV POSTGRES_DB=postgres

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o distributed-calculator cmd/distributed-calculator/main.go

EXPOSE 8080

CMD ["./distributed-calculator"]