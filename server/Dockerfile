FROM golang:1.25-bookworm AS builder

WORKDIR /usr/app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/main .

FROM alpine:latest

COPY --from=builder /out/main /out/main

ENTRYPOINT ["/out/main"]
