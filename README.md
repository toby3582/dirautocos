# dirwatchsynccos
监听指定文件夹自动同步腾讯云存储

### 依赖

- github.com/fsnotify/fsnotify
- github.com/tencentyun/cos-go-sdk-v5
- github.com/urfave/cli/v2

### 使用方式

```
git clone https://github.com/toby3582/dirwatchsynccos.git
#修改config.yaml 完善腾讯云cos信息
cd dirwatchsynccos
go mod tidy
#首次运行初始化 同步全部到腾讯云
go run main.go cos-init

# 后台守护进程运行 
go run main.go cos-watch
```