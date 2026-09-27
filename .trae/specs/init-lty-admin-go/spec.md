# lty-admin-go 后台管理系统 Spec

## Why
用户需要一个基于 Go 语言的高质量后台管理系统，要求代码结构规范，所有的文件、数据库表名等均需包含 `lty` 前缀，并且后台界面需要美观。后端（包括后台管理前端界面）存放在 `e:\lty-wzxt\hd` 目录，未来面向用户的前台将存放在 `e:\lty-wzxt\qd` 目录，数据库使用 MySQL。

## What Changes
- 初始化 Go 后端项目（基于 Gin 等流行框架，采用分层架构），存放于 `e:\lty-wzxt\hd\lty_backend`。
- 初始化美观的后台管理前端项目（基于 Vite + Vue3 + Element Plus，项目名为 `lty_admin_ui`，存放于 `e:\lty-wzxt\hd`）。
- 配置 MySQL 数据库连接，设计带有 `lty_` 前缀的初始数据表（如管理员表 `lty_admins`）。
- 建立规范的目录结构（如 lty_api, lty_router, lty_service, lty_model, lty_config 等），所有核心文件和包名包含 `lty` 前缀。

## Impact
- Affected specs: 无（全新项目）
- Affected code: `e:\lty-wzxt\hd` 目录下的所有新建代码。

## ADDED Requirements
### Requirement: Go 后端基础架构
系统 SHALL 提供一个标准的 Go Web 服务，连接 MySQL 数据库，且数据库表名均带有 `lty_` 前缀。所有代码文件和包命名规范包含 `lty` 标识。

#### Scenario: 成功启动后端服务
- **WHEN** 开发者运行 `go run lty_main.go`
- **THEN** 服务在指定端口启动，成功连接 MySQL，并暴露基础健康检查 API。

### Requirement: 美观的后台管理前端
系统 SHALL 提供一个美观、现代的后台管理前端界面。

#### Scenario: 访问后台界面
- **WHEN** 开发者启动前端项目并访问浏览器
- **THEN** 渲染出带有精美 UI 的后台登录及控制台页面。
