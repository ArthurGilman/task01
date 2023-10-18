FROM golang:1.20.3 AS builder

WORKDIR /src

COPY ./go.mod .

RUN go mod download

COPY ./ ./

RUN GOOS=linux go build -o /app cmd/main/main.go

FROM ubuntu:23.04

WORKDIR /app

COPY --from=builder app app

COPY --from=builder /src/struct.txt struct.txt

CMD ["./app"]