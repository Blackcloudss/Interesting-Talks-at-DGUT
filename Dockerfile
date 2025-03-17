# 构建过程
FROM golang:1.22.3 AS builder

WORKDIR /app

# 先单独复制依赖文件
COPY go.mod go.sum ./

# 下载依赖（利用Docker缓存层）
RUN go mod download


COPY . .

RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w CGO_ENABLED=0 \
    # 打印环境变量
    && go env \
    && go mod tidy \
    && go build -ldflags="-s -w" -o interesting-talks ./cmd/main.go

# 执行过程
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/interesting-talks  ./
COPY --from=builder /app/config.yaml  ./

#声明容器运行时监听 8080 端口
EXPOSE 8080
#启动应用程序
ENTRYPOINT ./interesting-talks