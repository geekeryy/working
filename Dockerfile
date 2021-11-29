FROM golang:1.17.0
MAINTAINER jiangyang.me@gmail.com
WORKDIR /app
ENV TZ="Asia/Shanghai"
COPY main .
EXPOSE 8080
CMD ["./main"]
