FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server main.go

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /app/server .
COPY web ./web

ENV TODO_PORT=7540
ENV TODO_DBFILE=scheduler.db
ENV TODO_PASSWORD=12345

EXPOSE 7540

CMD ["./server"]