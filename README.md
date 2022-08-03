### 说明
- 基于gin进行的二次封装项目
- 没有希望提供一个简单的框架结构，方便使用。尽量写一些注释，方便初学者参考理解

### 环境
- go > 1.15
- mysql
- redis (非必需)

### 目录结构

```
    ├── config              (配置解析)
    |── internal
    |   ├── common          (数据库操作)
    |   |   ├── constant    (常量)
    |   |   ├── dto         (请求、返回结构体)
    |   ├── controller      (控制层、路由器)
    |   ├── middleware      (中间件)                    
    |   ├── model           (模型层)                        
    |   └── service         (逻辑层)                    
    └── pkg  
        ├── curl            (HTTP请求)
        ├── db              (mysql、redis)
        ├── response        (返回封装)
        ├── cron            (定时任务)
        ├── token           (JWT)
        └── util            (工具包)          
```

### 使用

- `git clone https://github.com/i-a-a/gin-demo.git`

- 根据 `config/config.default.yaml` 配置 `config/config.yaml` 

- `go mod tidy` 

- 首次进行数据库表迁移  `go run main.go -m true`  


### 联系我
QQ：772532526