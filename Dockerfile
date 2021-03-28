FROM golang:1.16 AS builder

WORKDIR /build
COPY . .

RUN go build -o /go/bin/api cmd/api/api.go
COPY config.yml /go/bin

EXPOSE 8080

CMD ["/go/bin/api"]
