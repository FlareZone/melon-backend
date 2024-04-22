# 使用 golang 官方镜像作为基础镜像
FROM golang:latest AS builder

# 设置工作目录
WORKDIR /app

# 将当前目录下的所有文件复制到工作目录
COPY . .

# 构建应用程序
RUN go build -o melon-service main.go

# 运行阶段
FROM golang:latest

# 设置工作目录
WORKDIR /app

# 从构建阶段中复制编译后的可执行文件到当前镜像
COPY --from=builder /app/melon-service .

# 运行应用程序
CMD ["./melon-service", "api"]
