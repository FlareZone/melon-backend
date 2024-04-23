# 使用 golang 官方镜像作为基础镜像
FROM golang:latest
MAINTAINER ltll <lantianlaoli@gmail.com>

# 将工作目录设置为 /app
WORKDIR /app

# 将当前目录下的所有文件复制到工作目录中
COPY . .

# 运行应用程序
CMD ["go", "run", "main.go", "api"]
