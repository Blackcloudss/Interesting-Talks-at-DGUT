# 构建过程
FROM golang:1.22.3 AS builder

WORKDIR /app

# 先单独复制依赖文件
COPY go.mod go.sum ./

# 设置容器内的 Go 环境变量
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0

# 下载依赖（新增缓存验证步骤）

# 打印环境变量
# -x 显示详细下载过程
# 验证 Go 版本
RUN go mod download -x \
    && go version \
    && go env

# 复制全部代码
COPY . .

# 编译前验证依赖
RUN go mod verify \
    && go build -v -ldflags="-s -w" -o interesting-talks ./cmd/main.go

# 执行过程
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/interesting-talks  ./
COPY --from=builder /app/config.yaml  ./

#声明容器运行时监听 8080 端口
EXPOSE 8080
#启动应用程序
ENTRYPOINT ["./interesting-talks"]