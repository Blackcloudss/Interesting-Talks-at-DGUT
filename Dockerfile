# 构建过程
FROM golang:1.22.3 AS builder

WORKDIR /app

# 先单独复制依赖文件
COPY go.mod go.sum ./

# 设置容器内的 Go 环境变量
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0

# 下载依赖
RUN go mod download -x

# 完整代码复制（在整理依赖前确保所有代码就位）
COPY . .

# 整理验证依赖
RUN go mod tidy \
    && go mod verify \
    && go version \
    && go env


# 编译
RUN go build -v -ldflags="-s -w" -o interesting-forum ./cmd/main.go

# 执行过程
FROM alpine:latest

# 设置时区
RUN apk add --no-cache tzdata && \
    ln -snf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app

COPY --from=builder /app/interesting-forum ./
COPY --chmod=644 config.yaml ./

#声明容器运行时监听 8080 端口
EXPOSE 8080
#启动应用程序
ENTRYPOINT ["./interesting-forum"]