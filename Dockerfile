# 使用 golang 官方镜像作为基础镜像
FROM golang:latest

# 设置工作目录
WORKDIR /app

# 将当前目录下的所有文件复制到工作目录
COPY . .

# 构建并运行应用
CMD ["go", "run", "main.go", "api"]
