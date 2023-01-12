FROM golang:1.18.3-alpine as build

ENV CGO_ENABLED=0 \
    GOPATH=/go \
    GOOS=linux \
    GOARCH=amd64

WORKDIR $GOPATH/tools.liu.app

COPY . $GOPATH/tools.liu.app

RUN go build -v -a -o tools .

# 新增一个容器，用来运行应用
FROM alpine as run

RUN apk --no-cache add ca-certificates

WORKDIR /tools.liu.app

RUN mkdir -p /var/log/tools.liu.app/

COPY ./templates ./templates
COPY ./pages ./pages
COPY ./static ./static
COPY ./lang ./lang
COPY ./*.yaml .
COPY ./sitemap.xml .
COPY ./robots.txt .

# 从 build 镜像中把/tools 拷贝到当前目录
COPY --from=build /go/tools.liu.app/tools .

EXPOSE 8088

ENTRYPOINT ["./tools"]

CMD ["-config", "./config_production.yaml"]
