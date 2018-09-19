FROM golang:1.11-alpine3.8 as builder
WORKDIR /go/src/git.trj.tw/golang/mtfosbot
COPY . .
RUN GOOS=linux go build -o mtfosbot .

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /data
COPY --from=builder /go/src/git.trj.tw/golang/mtfosbot/mtfosbot /usr/bin
COPY config.default.yml config.yml
EXPOSE 10230
CMD ["/usr/bin/mtfosbot", "-f", "/data/config.yml", "-dbtool"]