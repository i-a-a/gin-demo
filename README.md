### 说明
- 基于gin、gorm的封装。 
- 提供一个简单的框架结构，方便使用。尽量写一些注释，方便初学者参考理解


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
        ├── db              (mysql、redis)
        ├── logger          (日志)
        ├── response        (返回封装)
        ├── token           (JWT)
      
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