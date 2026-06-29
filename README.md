# mycc

基于 Gin 框架的 Go Web 服务基础框架，使用 Go 1.25.5 开发。

## 特性

- 🚀 基于 [Gin](https://github.com/gin-gonic/gin) 的高性能 HTTP 框架
- 📝 基于 [Zap](https://github.com/uber-go/zap) 的结构化日志，支持文件轮转
- ⚙️ 基于 [Viper](https://github.com/spf13/viper) 的多环境配置管理
- 🌍 自动环境识别（development/testing/staging/production）
- 🔌 基于 [Wire](https://github.com/google/wire) 的编译时依赖注入
- 🛡️ 内置中间件：请求日志、Panic 恢复、CORS、RequestID
- 🧪 TDD 开发流程，测试覆盖核心功能

## 快速开始

### 前置要求

- Go 1.25.5+
- Wire CLI（`go install github.com/google/wire/cmd/wire@latest`）

### 安装

```bash
git clone https://github.com/sunyaofeng/mycc.git
cd mycc

# 下载依赖
go mod tidy

# 生成 Wire 依赖注入代码
wire ./cmd/mycc

# 验证编译
go build ./cmd/mycc
```

### 运行

```bash
# 默认开发环境
go run ./cmd/mycc

# 指定环境
MYCC_ENV=production go run ./cmd/mycc
```

### 构建

```bash
go build -o bin/mycc ./cmd/mycc
./bin/mycc
```

## API 接口

| 方法 | 路径 | 说明 | 响应示例 |
|------|------|------|---------|
| GET | `/health` | 健康检查 | `{"status":"ok","timestamp":"2026-06-29 15:30:00"}` |
| GET | `/api/v1/ping` | 连通测试 | `{"message":"pong"}` |
| GET | `/api/v1/date` | 获取当前日期 | 见下方 |

### 日期接口响应示例

```json
{
  "date": "2026-06-29",
  "time": "2026-06-29 15:30:00",
  "timestamp": 1782719400,
  "timezone": "Local"
}
```

## 项目结构

```
mycc/
├── cmd/mycc/               # 应用程序入口
│   ├── main.go             # 主函数
│   └── wire.go             # Wire 注入器定义
├── internal/               # 内部包（不对外暴露）
│   ├── config/             # 配置管理（Viper）
│   ├── env/                # 环境识别
│   ├── handler/            # 请求处理器
│   │   ├── date.go         # 日期接口
│   │   └── health.go       # 健康检查接口
│   ├── logger/             # 日志系统（Zap + Lumberjack）
│   ├── middleware/          # Gin 中间件
│   ├── router/             # 路由定义
│   └── wire/               # Wire Provider 集合
├── configs/                # 配置文件
│   ├── config.yaml         # 基础配置
│   ├── config.development.yaml  # 开发环境
│   └── config.production.yaml   # 生产环境
├── go.mod
└── go.sum
```

## 配置

### 多环境配置

配置文件按 `config.{env}.yaml` 命名，根据 `MYCC_ENV` 环境变量自动加载：

1. 先加载 `config.yaml`（基础配置）
2. 再合并 `config.{env}.yaml`（环境特定配置，覆盖基础配置）

### 配置项

```yaml
server:
  port: 8080              # 服务端口
  mode: debug             # Gin 模式：debug / release
  read_timeout: 30        # 读超时（秒）
  write_timeout: 30       # 写超时（秒）

log:
  level: info             # 日志级别：debug / info / warn / error
  format: console         # 日志格式：console / json
  output: stdout          # 输出目标：stdout / file / both
  file_path: ./logs/app.log
  max_size: 100           # 单文件最大 MB
  max_backups: 10         # 保留文件数
  max_age: 30             # 保留天数

database:
  driver: mysql
  host: localhost
  port: 3306
  user: root
  password: ""
  dbname: mycc
  max_open_conns: 100
  max_idle_conns: 10
  conn_max_lifetime: 3600
```

### 环境变量

环境变量优先级高于配置文件，前缀为 `MYCC_`，点号替换为下划线：

| 环境变量 | 对应配置项 | 说明 |
|---------|-----------|------|
| `MYCC_ENV` | - | 运行环境：development / testing / staging / production |
| `MYCC_SERVER_PORT` | server.port | 服务端口 |
| `MYCC_LOG_LEVEL` | log.level | 日志级别 |
| `MYCC_DATABASE_HOST` | database.host | 数据库地址 |

### 环境识别优先级

`MYCC_ENV` > `GO_ENV` > `ENV`，默认为 `development`。

## 日志

基于 Zap 的结构化日志，支持以下配置：

- **格式**：Console（开发友好） / JSON（生产环境）
- **输出**：stdout / file / both（同时输出到控制台和文件）
- **轮转**：文件输出使用 Lumberjack 自动轮转，支持按大小切割、按数量/天数保留

日志中间件自动记录每个请求的状态码、延迟、客户端IP、请求方法、路径等信息。

## 依赖注入

使用 [Wire](https://github.com/google/wire) 实现编译时依赖注入：

```
cmd/mycc/wire.go       → 定义注入器（Injector）
cmd/mycc/wire_gen.go   → Wire 自动生成（已加入 .gitignore）
internal/wire/wire.go  → Provider 集合
```

### 新增 Handler 流程

1. 在 `internal/handler/` 下创建 Handler
2. 在 `cmd/mycc/wire.go` 的 `wire.Build` 中添加 Provider
3. 运行 `wire ./cmd/mycc` 重新生成注入代码

## 测试

```bash
# 运行所有测试
go test ./...

# 运行指定包测试
go test ./internal/handler/ -v

# 运行单个测试
go test -run TestDateHandler_GetCurrentDate ./internal/handler/

# 查看覆盖率
go test -cover ./...
```

## VS Code 调试

`.vscode/` 目录已被 `.gitignore` 忽略，需自行创建 `.vscode/launch.json`：

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch mycc (development)",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/cmd/mycc",
      "cwd": "${workspaceFolder}",
      "env": { "MYCC_ENV": "development" },
      "console": "integratedTerminal"
    },
    {
      "name": "Launch mycc (debug mode)",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/cmd/mycc",
      "cwd": "${workspaceFolder}",
      "env": { "MYCC_ENV": "development" },
      "console": "integratedTerminal"
    },
    {
      "name": "Test All",
      "type": "go",
      "request": "launch",
      "mode": "test",
      "program": "${workspaceFolder}",
      "args": ["./..."],
      "console": "integratedTerminal"
    }
  ]
}
```

推荐扩展：`golang.go`、`eamodio.gitlens`、`mhutchie.git-graph`、`redhat.vscode-yaml`

## 技术栈

| 组件 | 技术 | 版本 |
|------|------|------|
| 语言 | Go | 1.25.5 |
| Web 框架 | Gin | v1.10.0 |
| 配置管理 | Viper | v1.20.1 |
| 结构化日志 | Zap | v1.27.0 |
| 日志轮转 | Lumberjack | v2.2.1 |
| 依赖注入 | Wire | v0.7.0 |
| 测试断言 | Testify | v1.11.1 |

## License

MIT
