FROM golang:1.17.0
MAINTAINER jiangyang.me@gmail.com
WORKDIR /app
ENV TZ="Asia/Shanghai"
RUN go env -w GOPROXY=https://goproxy.cn,direct && go install github.com/go-delve/delve/cmd/dlv@latest
COPY main ./
EXPOSE 8080 8081 6060 2345
CMD ["dlv", "--listen=:2345", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "./main"]
