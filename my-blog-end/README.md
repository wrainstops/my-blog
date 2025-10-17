### RUN

```sh
go mod tidy
go run ./

# 访问 http://127.0.0.1:8081/swagger/index.html 查看swagger
```

### 开发部署中问题记录
#### 开发问题
  1. 如果发现包拉不下来, GoLand --- settings中检查GOROOT / GOPATH / GO Modules
                      VsCode --- go env -w GO111MODULE=on
                                 go env -w GOPROXY=https://goproxy.cn,direct
  2. swag使用外部包的变量，init时命令：```swag init --parseDependency --parseInternal```
  3. 入参结构体改用ReqXXX后，使用binding: "required"进行必填校验，此时使用bind接收参数会触发panic，改用ShouldBindJSON，并进行错误return。
#### docker
  1. 配置文件application.yml等需和Dockerfile在同层下，否则COPY不到，**暂未解决**。
  2. docker访问宿主机mysql，将配置文件的host改为host.docker.internal即可。
  3. 需要在Dockerfile中暴露服务端口，该项目是EXPOSE 8081。
