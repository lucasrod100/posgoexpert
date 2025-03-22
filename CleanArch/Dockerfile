FROM golang:1.22 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy

EXPOSE 8000 50051 8080

CMD ["go", "run", "main.go", "wire_gen.go"]