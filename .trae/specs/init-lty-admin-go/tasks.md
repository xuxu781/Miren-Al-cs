# Tasks
- [x] Task 1: 创建后端项目目录并初始化 Go 模块
  - [x] SubTask 1.1: 在 `e:\lty-wzxt\hd` 创建后端项目目录 `lty_backend`
  - [x] SubTask 1.2: 使用 `go mod init lty_backend` 初始化模块
  - [x] SubTask 1.3: 建立规范的项目结构（如 `lty_config`, `lty_models`, `lty_routers`, `lty_services`, `lty_controllers`），并编写 `lty_main.go`
- [x] Task 2: 配置 MySQL 数据库和基础模型
  - [x] SubTask 2.1: 引入 GORM 和 MySQL 驱动
  - [x] SubTask 2.2: 编写数据库连接配置，并设置表名前缀为 `lty_`
  - [x] SubTask 2.3: 创建基础的管理员数据模型（对应表 `lty_admins`）
- [x] Task 3: 编写核心 API 和路由
  - [x] SubTask 3.1: 搭建 Gin Web 框架
  - [x] SubTask 3.2: 编写健康检查及测试 API 接口
  - [x] SubTask 3.3: 配置路由组，允许跨域请求（CORS）
- [x] Task 4: 初始化美观的后台前端项目
  - [x] SubTask 4.1: 在 `e:\lty-wzxt\hd` 中初始化前端项目 `lty_admin_ui`（使用 Vite + Vue3 + Element Plus）
  - [x] SubTask 4.2: 配置前端项目的基本结构、路由、状态管理（如 Pinia）
  - [x] SubTask 4.3: 设计并实现基础的登录页和后台布局（Layout）

# Task Dependencies
- [Task 2] depends on [Task 1]
- [Task 3] depends on [Task 2]
- [Task 4] depends on [Task 1]
