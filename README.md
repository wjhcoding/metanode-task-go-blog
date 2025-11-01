## 初始化 Go 模块

在项目根目录下执行以下命令，初始化 Go 模块并自动整理依赖：

```bash
go mod init github.com/wjhcoding/metanode-task-go-blog
go mod tidy
```


```
metanode-task-go-blog/
├── api/                     # 控制器层（HTTP接口层）
    ├── v1/
    │   ├── user_api.go      # 用户注册、登录
    │   ├── post_api.go      # 文章 CRUD
    │   └── comment_api.go   # 评论创建与查询
├── cmd/
│   └── main.go              # 程序入口（初始化配置、日志、路由）
├── config/
│   └── toml_config.go       # 配置文件解析逻辑
├── internal/
│   ├── dao/                 # 数据访问层（数据库操作）
│   │   └── pool/mysql_tool.go
│   ├── model/               # 数据模型（对应数据库表）
│   ├── router/              # 路由定义
│   │   └── router.go
│   └── service/             # 业务逻辑层
├── pkg/
│   ├── common/
│   │   └── response/        # 通用响应结构封装
│   │       └── response_msg.go
│   └── global/
│       └── log/             # 日志封装
│           └── logger.go
├── go.mod
├── go.sum
└── README.md
```