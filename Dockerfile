# 阶段1：构建
FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o faw_vw_auto .

# 阶段2：运行
FROM alpine:3.18
WORKDIR /app

# 安装时区数据并设置为上海时区（Asia/Shanghai）
RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

# 从构建阶段复制文件
COPY --from=builder /app/faw_vw_auto .
COPY --from=builder /app/config/prod.yaml ./env.yaml
COPY --from=builder /app/FAWVW.db ./FAWVW.db
# 设置非root用户
RUN adduser -D appuser && chown -R appuser /app
USER appuser

# 允许通过环境变量覆盖配置路径
# ENV CONFIG_PATH=/app/env.yaml
CMD ["./faw_vw_auto"]