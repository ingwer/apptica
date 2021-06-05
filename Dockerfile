FROM golang:1.16-alpine3.13 as builder
WORKDIR /go/src/apptica
COPY . .
RUN go build -v -o /app

EXPOSE 8080
EXPOSE 9000

FROM alpine:3.13
COPY --from=builder /app /app
ENTRYPOINT ["/app"]
