# 替换生产环境所使用的配置文件
# cp config/$1.yaml ./env.yaml

# 构建出应用程序包（Linux）
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o faw_vw_auto