# db-connect-demo — Go 微服务 / CLI 可插拔后端负载测试工具

## 简介

轻量级负载测试与演示工具，支持多种后端驱动（可通过 lib 包注册）：Postgres、MySQL、SQLite、Kafka、Solace（基于 MQTT）。工具既可作为命令行程序运行压力测试，也可作为 HTTP 服务动态管理多个后端实例。

## 主要功能

- 驱动：postgres、mysql、sqlite、kafka、solace（MQTT）
- 动态注册后端并启动：通过 JSON 文件或命令行传入后端字典
- HTTP API（Gin）：
  - GET /ping — 返回已注册后端 health（"ok" 或错误信息）
  - POST /query — 输入 {"backend":"name","query":"..."}，返回查询结果（消息驱动返回不支持查询）

## 驱动示例 DSN

- Postgres:

  postgres://user:pass@localhost:5432/dbname?sslmode=disable

- MySQL:

  user:pass@tcp(localhost:3306)/dbname

- SQLite (in-memory):

  file::memory:?cache=shared

- Kafka (producer, brokers list):

  localhost:9092,broker2:9092

- Solace (MQTT):

  tcp://broker-host:1883?clientid=bench1

## 构建

```bash
go mod tidy
go build -o db-bench.exe
```

> Windows: 在 PowerShell 中可运行 `.
\db-bench.exe`；在 Unix-like 系统运行 `./db-bench.exe`（或交叉编译为 linux 可执行文件）。

## 运行（作为服务）

示例 backends.json：

```json
{
  "pg1": {"driver":"postgres","dsn":"postgres://user:pass@localhost:5432/db?sslmode=disable"},
  "k1": {"driver":"kafka","dsn":"localhost:9092"},
  "s1": {"driver":"solace","dsn":"tcp://broker:1883?clientid=bench1"}
}
```

启动服务：

```bash
./db-bench.exe -backends-file=backends.json -port=8080
```

打开 UI：

访问 http://localhost:8080/ui 来打开内置的最小 React 前端（或访问 /ui 以在浏览器中打开）。


## HTTP API 示例

- GET /ping

```bash
curl http://localhost:8080/ping
```

- POST /query

```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{"backend":"pg1","query":"SELECT 1"}' \
  http://localhost:8080/query
```

## Makefile

仓库包含 Makefile，常用目标：

- `make tidy` — 同步模块
- `make build` — 本地构建
- `make test-build` — 交叉构建示例
- `make run-service BACKENDS=backends.json PORT=8080` — 启动服务
- `make run-postgres DSN="..."` 等运行示例目标（请编辑或覆盖 DSN）

## 备注

- SQLite 在 Windows 上需要启用 CGO（安装 MinGW/MSYS2）；
- Solace 驱动使用 MQTT（paho）实现发布；如需使用 Solace 官方 SDK，请提供 SDK 详情以实现替代驱动。

