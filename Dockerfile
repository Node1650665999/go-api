##############################step one ###########################
FROM golang:1.16-alpine as builder

ENV GOPROXY https://goproxy.cn

# 更新安装源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 安装 git
RUN apk --no-cache add git

# 设置工作目录
WORKDIR /app/go-api

# 将当前项目所在目录代码拷贝到镜像中
COPY . .

# 确保没有将.env打包进去
RUN if [ -e .env ] ; then rm .env; fi

# 下载依赖
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o go-api

################################step two ###############################
FROM alpine:latest

# 更新安装源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 安装相关软件
RUN apk update && apk add --no-cache bash supervisor ca-certificates

RUN mkdir /app
WORKDIR /app

COPY --from=builder /app/go-api/go-api .

ADD supervisord.conf /etc/supervisord.conf

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]