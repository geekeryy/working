FROM golang:1.17.0 as builder
MAINTAINER jiangyang.me@gmail.com
WORKDIR /app
ENV TZ="Asia/Shanghai"
RUN go env -w GOPROXY=https://goproxy.cn,direct && go install github.com/go-delve/delve/cmd/dlv@latest


FROM golang:1.17.0
COPY --from=0 /go/bin/dlv /usr/bin