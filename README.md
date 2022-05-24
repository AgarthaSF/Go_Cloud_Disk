# CloudDisk

> 轻量级云盘系统，基于`go-zero`, `xorm`实现。

使用到的命令

```text
# 创建API服务
goctl api new core
cd core
# 启动服务
go run core.go -f etc/core-api.yaml
# 使用API文件生成代码
goctl api go -api core.api -dir . -style go_zero

# 配置邮件服务
go get github.com/jordan-wright/email

# 配置 go-redis服务
go get github.com/go-redis/redis/v8

# 配置 go-UUID服务
go get github.com/satori/go.uuid

# 腾讯云 COS Go SDK安装
go get -u github.com/tencentyun/cos-go-sdk-v5


```
