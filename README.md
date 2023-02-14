
## 暴露端口
8091

## 暴露地址

/etc/config

## docker构建

` docker build -t a624669980/get-last-release:{Version} . `



## docker运行

` docker run -d --name get-last-release -p 8091:8091 -v /etc/config:/etc/config a624669980/get-last-release:{Version} `