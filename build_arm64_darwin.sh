# 替换生产环境所使用的配置文件
# cp config/$1.yaml ./env.yaml

# 构建出应用程序包（Linux）
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o faw_vw_auto_darwin
