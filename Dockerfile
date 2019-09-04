FROM golang:1.12-alpine AS builder

ENV APP_HOME=/go/src/NightWolf007/rclip
WORKDIR $APP_HOME
RUN apk add --no-cache git && \
    go get -u github.com/golang/dep/cmd/dep
COPY . $APP_HOME/
RUN dep ensure && \
    go build -o rclip ./main.go

FROM alpine:latest
WORKDIR /root
RUN echo "listen: '0.0.0.0:9889'" > ./rclipd.yaml
COPY --from=builder /go/src/NightWolf007/rclip/rclip .
EXPOSE 9889
CMD ["./rclip", "server", "-c", "./rclipd.yaml"]
