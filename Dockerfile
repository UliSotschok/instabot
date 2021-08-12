FROM golang:latest AS builder

RUN mkdir /app
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -o instabot

FROM alpine:latest AS production

RUN mkdir /app
WORKDIR /app
COPY --from=builder /app .
CMD ["./instabot"]
