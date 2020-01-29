FROM golang:latest

                       ENV GO111MODULE=on
                       ENV GOPROXY https://goproxy.cn,direct

                       WORKDIR "F:\code\002Golang\lwx-crontab"
                       COPY . .
                       RUN cd master/main
                       RUN go build

                       EXPOSE 8000
                       ENTRYPOINT ["./master"]