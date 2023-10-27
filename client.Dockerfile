FROM golang:1.19.4 AS builder

WORKDIR /build

COPY . .

RUN go mod download

RUN GOOS=linux GOARCH=amd64 go build -o main ./cmd/client

FROM scratch

ENV SERVER_HOST=localhost
ENV SERVER_PORT=3131

COPY --from=builder /build/main /

ENTRYPOINT ["/main"]