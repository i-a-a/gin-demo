### 说明
- 基于gin、gorm的封装。 
- 提供一个简单的框架结构，方便使用。尽量写一些注释，方便初学者参考理解


### 目录结构

```
    ├── config              (配置)
    |── internal            (业务代码)
    |   ├── common          (数据库操作)
    |   |   ├── dto         (请求、返回结构体)
    |   |   ├── enum        (枚举值)
    |   |   ├── helper      (自定义工具函数，具有业务属性)
    |   |   ├── sdk         (SDK)
    |   ├── controller      (控制层、路由器)
    |   ├── middleware      (中间件)                    
    |   ├── model           (模型层)                        
    |   └── service         (逻辑层)                    
    |── pkg                 (工具包，不包含业务逻辑)
    |   ├── db              (mysql、redis)
    |   ├── response        (返回封装)
    |   ├── token           (JWT)
    |   └── util            (工具函数: 加密、文件操作、HTTP请求、随机数、机器人通知)    
    |── script              (脚本)        
    |── static              (静态资源)


```

### 使用

```sh 
git clone https://github.com/i-a-a/gin-demo.git
cd gin-demo
go mod init app
go mod tidy
```

- 根据 `config/config.default.yaml` 配置 `config/config.yaml` 

- 首次进行数据库表迁移  `go run main.go -m true`  

### 联系我
QQ：772532526