# db-connect-demo — Go 微服务 / CLI 可插拔后端负载测试工具

## 简介

轻量级负载测试与演示工具，支持多种后端驱动（通过 lib 包注册）：Postgres、MySQL、SQLite、Kafka、Solace（基于 MQTT）。工具既可作为命令行程序运行压力测试，也可作为 HTTP 服务动态管理多个后端实例并暴露交互式 Web UI。

## 主要功能

- 多驱动支持：postgres、mysql、sqlite、kafka、solace（MQTT）
- 动态注册后端并统一管理（通过 JSON 文件或启动参数传入后端字典）
- HTTP API（Gin）：
  - GET /ping — 返回已注册后端健康信息（"ok" 或错误信息）
  - POST /query — 输入 {"backend":"name","query":"..."}，返回查询结果（消息驱动返回不支持查询）
- 可选前端（React + Vite + Bootstrap），构建产物位于 frontend/dist，服务在 /ui

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

后端（Go）依赖与构建：

```bash
go mod tidy
go build -o db-bench.exe
```

前端（可选，使用 Node.js + npm/yarn）：

```bash
cd frontend
npm ci           # 或 yarn
npm run build     # 生成 frontend/dist
```

Makefile 的 `run-service` 目标会在启动前尝试构建 frontend（若本机可用 npm）。

> 提示：请确保安装了 Node.js（推荐 16+）和 npm/yarn 以构建前端。

## 开发模式（并行运行前后端）

- 后端：
  ```bash
  go run main.go -backends-file=backends.json -port=8080
  ```
- 前端（开发服务器）：
  ```bash
  cd frontend
  npm install
  npm run dev     # Vite dev server，默认 5173
  ```

在开发模式下可通过 Vite 代理或手动配置 CORS 将前端请求代理到后端。

## 运行（作为服务）

示例 backends.json：

```json
{
  "pg1": {"driver":"postgres","dsn":"postgres://user:pass@localhost:5432/db?sslmode=disable"},
  "mysql1": {"driver":"mysql","dsn":"user:pass@tcp(localhost:3306)/dbname"},
  "kafka1": {"driver":"kafka","dsn":"localhost:9092"},
  "solace1": {"driver":"solace","dsn":"tcp://broker:1883?clientid=bench1"}
}
```

启动已构建服务并访问 UI：

```bash
./db-bench.exe -backends-file=backends.json -port=8080
# 浏览器访问 http://localhost:8080/ui
```

或者使用 Makefile：

```bash
make run-service BACKENDS=backends.json PORT=8080
```

## API 示例

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

## 常见问题与提示

- 若某些后端（例如 Kafka、Solace）在启动时无法连接，服务会记录为 warning 并继续运行，/ping 会展示失败原因；对这些后端的 /query 将返回注册失败的错误信息。
- SQLite 在 Windows 上需要启用 CGO（安装 MinGW/MSYS2）；
- 若需 Solace 的专有 API（非 MQTT），请提供官方 SDK 信息以便集成。

## 贡献

欢迎通过 PR 添加更多驱动、示例和改进文档。

