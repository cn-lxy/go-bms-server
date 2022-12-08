# 笔记


## 1. 构建docker image

```shell
$ docker build -t <tag name> .
```

## 2. 后台运行 docker 镜像
> -p: 端口映射， -d：后台运行

```shell
sudo docker run -p 8080:8080 -d go-bms-server
```

## 3. 重命名docker image

```shell
$ sudo docker tag <image id> <new tag name>
```

## 3. 删除docker image

```shell
$ sudo docker rmi <image tag name>
```
## 4. 我的数据库配置
```go
// windows 配置
const (
	DataBaseUser     = "root"
	DataBasePassword = "123456"
	DataBaseHost     = "127.0.0.1"
	DataBasePort     = "3306"
	DataBaseName     = "web_bms"
)

// wsl docker 配置
const (
	DataBaseUser     = "lxy"
	DataBasePassword = "123456"
	DataBaseHost     = "172.17.0.1" // 宿主相对于docker的地址: 172.17.0.1
	DataBasePort     = "3306"
	DataBaseName     = "web_bms"
)

// 腾讯云配置
const (
	DataBaseUser     = "lxy" // GOOS=linux GOARCH=arm64 go build xxx.go
	DataBasePassword = "LXY1019XYXYZ"
	DataBaseHost     = "42.192.149.39"
	DataBasePort     = "3306"
	DataBaseName     = "bms"
)
```