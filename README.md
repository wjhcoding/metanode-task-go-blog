# 📝 Metanode Task Go Blog

一个使用 **Go (Golang)** + **Gin** + **GORM** + **MySQL** 开发的简易博客系统后端，  
支持用户注册登录（JWT认证）、文章管理（CRUD）、评论功能与统一响应结构。

---

## 🚀 功能概览

### ✅ 用户模块
- 用户注册（密码加密存储）
- 用户登录（JWT生成与验证）
- 用户信息获取

### ✅ 文章模块
- 创建文章（需登录）
- 获取文章列表与详情
- 更新文章（仅作者）
- 删除文章（仅作者）

### ✅ 评论模块
- 文章评论（需登录）
- 查看文章下所有评论
- 删除评论（仅作者）

### ✅ 系统支持
- JWT 鉴权中间件
- 全局统一响应结构
- 全局异常恢复（Recovery）
- CORS 跨域支持
- 结构化日志（基于 zap）
- TOML 配置文件加载（基于 viper）


## 📂 项目结构
````
metanode-task-go-blog/
├── README.md
├── api/                      # 控制器层
│   └── v1/
│       ├── user_api.go
│       ├── post_api.go
│       └── comment_api.go
├── cmd/
│   └── main.go                # 程序入口
├── config/
│   └── toml_config.go         # 配置加载
├── internal/                  # 内部逻辑
│   ├── dao/pool/mysql_tool.go # 数据库连接池
│   ├── model/                 # 模型定义（User、Post、Comment）
│   ├── router/router.go       # 路由定义
│   └── service/               # 业务逻辑层（可扩展）
├── pkg/
│   ├── common/response/       # 统一响应结构
│   │   └── response_msg.go
│   └── global/log/            # 日志封装
│       └── logger.go
├── go.mod
├── go.sum
└── config.toml                # 配置文件

````

---

## ⚙️ 环境依赖

- **Go** >= 1.19  
- **MySQL** >= 5.7  
- **Gin**  
- **GORM**  
- **Viper**  
- **JWT-go**  
- **Zap**

安装依赖：
```bash
go mod init github.com/wjhcoding/metanode-task-go-blog
go mod tidy
````

---

## 🧱 数据库设计

### users 表

| 字段         | 类型       | 描述     |
| ---------- | -------- | ------ |
| id         | bigint   | 主键     |
| username   | varchar  | 用户名    |
| password   | varchar  | 加密后的密码 |
| email      | varchar  | 邮箱     |
| created_at | datetime | 注册时间   |

### posts 表

| 字段         | 类型       | 描述   |
| ---------- | -------- | ---- |
| id         | bigint   | 主键   |
| title      | varchar  | 文章标题 |
| content    | text     | 文章内容 |
| user_id    | bigint   | 作者ID |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### comments 表

| 字段         | 类型       | 描述      |
| ---------- | -------- | ------- |
| id         | bigint   | 主键      |
| content    | text     | 评论内容    |
| user_id    | bigint   | 评论者ID   |
| post_id    | bigint   | 评论的文章ID |
| created_at | datetime | 评论时间    |

---

## 🧩 配置文件（config.toml）

````toml
AppName = "metanode-task-go-blog"

[MySQL]
Host = "127.0.0.1"
Port = 3306
User = "root"
Password = "123456"
Name = "blogdb"
TablePrefix = "blog_"

[Log]
Path = "./logs"
Level = "info"

[StaticPath]
FilePath = "./uploads"

````

---

## 🏃‍♂️ 启动项目

### 1️⃣ 运行 MySQL 并导入表结构

````bash
mysql -u root -p
CREATE DATABASE blogdb CHARACTER SET utf8mb4;
````

### 2️⃣ 启动项目

````bash
go run cmd/main.go
````

服务器默认启动在：

````
http://localhost:8888
````

---

## 🔗 API 接口示例

### 🧍 用户注册

`POST /api/v1/user/register`

````json
{
  "username": "wjh",
  "password": "123456",
  "email": "wjh@example.com"
}
````

### 🔑 用户登录

`POST /api/v1/user/login`

````json
{
  "username": "wjh",
  "password": "123456"
}
````

返回：

````json
{
  "code": 200,
  "msg": "success",
  "data": {
    "token": "<JWT_TOKEN>"
  }
}
````

### 📝 创建文章

`POST /api/v1/posts`

````json
{
  "title": "我的第一篇博客",
  "content": "Hello, world!"
}
````

Header：

````
Authorization: Bearer <JWT_TOKEN>
````

---

## 🪵 日志示例

项目运行后会在 `logs` 目录生成日志文件：

````
logs/
 ├── app.log
````

---

## 🧠 后续扩展建议

* ✅ 支持分页与搜索；
* ✅ 管理员角色与后台管理；
* ✅ 支持文件上传与头像；
* ✅ 评论回复树形结构；
* ✅ 前端可视化（Vue3 + Element Plus）。

---

## 👨‍💻 作者信息

* 作者：**wjhcoding**
* 项目地址：[GitHub](https://github.com/wjhcoding/metanode-task-go-blog)
* 邮箱：`wjhcoding@example.com`