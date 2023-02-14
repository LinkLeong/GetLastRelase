
## 暴露端口
8091

## 暴露地址

/etc/config

## docker运行

` docker run -d --name get-last-release -p 8089:8089 -v /etc/config:/etc/config a624669980/get-last-release`

## 多平台构建

- 确定buildx可用   `docker buildx ls`
- 创建builder   `docker buildx create --name mybuilder`
- 查看builder   `docker buildx inspect --bootstrap`
- 构建并推送镜像   `docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t a624669980/get-last-release:{Version} --push .`
## docker构建

` docker build -t a624669980/get-last-release:{Version} . `



