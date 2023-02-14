# 使用最新版本的 Golang 镜像
FROM golang:latest

# 设置工作目录
WORKDIR /app

# 将当前目录的所有文件复制到工作目录
COPY . .

# 下载依赖
RUN go mod download

# 构建应用
RUN go build -o main .

# 暴露应用端口
EXPOSE 8080
VOLUME /etc/config

# 运行应用
CMD ["./main"]
