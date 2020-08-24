FROM golang:latest as builder

LABEL maintainer="Galileo Daras <galileo@getcoffee.io>"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY main.go server.go tournament.go ./
COPY dock/ dock/

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
RUN apk add --update --no-cache ca-certificates docker

WORKDIR /root/

COPY --from=builder /app/main .
EXPOSE 8080

CMD ["./main"]