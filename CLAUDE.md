# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

mycc 是一个基于 Gin 框架的 Go Web 服务基础框架，使用 Go 1.25.5 开发。项目提供了基础的日志、配置管理、环境识别和路由等核心功能。

## 常用命令

### VS Code 调试启动（推荐）

项目已配置 `.vscode/launch.json`，支持通过 VS Code 调试面板启动：

- **Launch mycc (development)** — 开发环境运行（`MYCC_ENV=development`）
- **Launch mycc (production)** — 生产环境运行（`MYCC_ENV=production`）
- **Launch mycc (debug mode)** — 带调试器运行（可断点调试）
- **Test Current Package** — 运行当前文件所在包的测试
- **Test Current File** — 运行当前文件中的测试（需选中测试函数名）
- **Test All** — 运行所有测试

快捷键：`Ctrl+Shift+D` 打开调试面板，选择配置后按 `F5` 启动。

### 命令行

```bash
# 下载依赖
go mod tidy

# 运行服务（开发环境）
go run ./cmd/mycc

# 或者指定环境变量
MYCC_ENV=development go run ./cmd/mycc

# 构建二进制文件
go build -o bin/mycc ./cmd/mycc

# 运行测试
go test ./...

# 运行单个测试
go test -run TestFunctionName ./path/to/package

# 运行测试并显示覆盖率
go test -cover ./...

# 代码格式化
go fmt ./...

# 代码检查
go vet ./...

# 生成 Wire 依赖注入代码
wire ./cmd/mycc
```

### 环境变量

- `MYCC_ENV` - 运行环境，可选值：`development`（默认）、`testing`、`staging`、`production`
- `MYCC_SERVER_PORT` - 服务器端口，覆盖配置文件中的端口
- `MYCC_LOG_LEVEL` - 日志级别：`debug`、`info`、`warn`、`error`

### 首次克隆后的初始化

```bash
# 1. 下载依赖
go mod tidy

# 2. 安装 Wire 工具（如未安装）
go install github.com/google/wire/cmd/wire@latest

# 3. 生成 Wire 依赖注入代码（wire_gen.go 被 .gitignore 忽略，需手动生成）
wire ./cmd/mycc

# 4. 验证编译
go build ./cmd/mycc
```

## 项目架构

### 目录结构

```
mycc/
├── cmd/mycc/           # 应用程序入口
│   └── main.go         # 主函数，初始化配置、日志、路由并启动服务
├── internal/           # 内部包，不对外暴露
│   ├── config/         # 配置管理（Viper）
│   ├── env/            # 环境识别
│   ├── logger/         # 日志（Zap + Lumberjack）
│   ├── middleware/     # Gin 中间件
│   └── router/         # 路由定义
├── configs/            # 配置文件
│   ├── config.yaml              # 基础配置
│   ├── config.development.yaml  # 开发环境配置
│   └── config.production.yaml   # 生产环境配置
├── go.mod              # Go 模块定义
└── README.md           # 项目说明
```

### 核心模块说明

**配置系统** (`internal/config/`)
- 使用 Viper 实现多环境配置加载
- 配置文件按 `config.{env}.yaml` 命名，自动根据 `MYCC_ENV` 环境变量加载
- 支持环境变量覆盖，前缀为 `MYCC_`，点号替换为下划线（如 `MYCC_SERVER_PORT`）
- 配置项包括：服务器端口/模式/超时、日志级别/格式/输出、数据库连接信息

**日志系统** (`internal/logger/`)
- 基于 Zap 实现结构化日志
- 支持多种输出：stdout、file、both（同时输出到控制台和文件）
- 支持 JSON 和 Console 两种格式
- 文件输出使用 Lumberjack 实现自动轮转（按大小/时间切割）
- 全局日志函数：`logger.Info()`、`logger.Debug()`、`logger.Warn()`、`logger.Error()`、`logger.Fatal()`

**环境识别** (`internal/env/`)
- 自动识别运行环境，优先级：`MYCC_ENV` > `GO_ENV` > `ENV`
- 提供便捷判断函数：`IsDevelopment()`、`IsTesting()`、`IsStaging()`、`IsProduction()`

**中间件** (`internal/middleware/`)
- `Logger()` - 请求日志，自动记录状态码、延迟、客户端IP、请求方法、路径
- `Recovery()` -  panic 恢复
- `CORS()` - 跨域支持
- `RequestID()` - 请求追踪ID

**路由** (`internal/router/`)
- 集中管理所有路由定义
- 默认注册全局中间件
- 提供 `/health` 健康检查端点和 `/api/v1/ping` 测试接口

## 开发规范

- 所有内部包放在 `internal/` 目录下
- 新增模块在 `internal/` 下创建子目录
- 路由在 `internal/router/` 中注册，按版本分组（`/api/v1/`）
- 配置项在 `Config` 结构体中添加，并在对应配置文件设置默认值
- 新增 Handler 后需在 `cmd/mycc/wire.go` 中添加 Provider，然后运行 `wire ./cmd/mycc` 重新生成注入代码
- VS Code 用户可直接使用调试面板启动服务，无需手动输入命令
