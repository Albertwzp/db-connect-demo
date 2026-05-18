# db-connect-demo - 统一的 Kubernetes Operator + API Server

## 项目简介

**db-connect-demo** 是一个轻量级的多驱动数据库连接管理工具，支持通过 Kubernetes CRD 动态管理多种数据库连接：

- ✅ PostgreSQL
- ✅ MySQL
- ✅ SQLite
- ✅ Kafka (消息队列)
- ✅ Solace (MQTT 消息)

**关键特性**：
- 🎯 单一 Go 二进制 (Operator + API Server 合一)
- 🔄 完全由 K8s CRD 驱动，无需静态配置文件
- 🚀 热更新：修改 CRD 自动重连，无需重启
- 📊 提供 HTTP API (`/ping`, `/query`)
- 🎨 可选 React Web UI (`/ui`)

---

## 两种运行模式

### 模式 A: 独立运行 (不需要 Kubernetes)

```bash
# 1. 构建
make build ui-build

# 2. 运行
make run

# 3. 访问
# 浏览器: http://localhost:8080/ui
# API: curl http://localhost:8080/ping
```

**特点**:
- ✅ 无依赖
- ✅ 快速启动
- ✅ 适合开发测试

### 模式 B: Kubernetes 部署 (推荐生产)

```bash
# 1. 安装 CRD
make k8s-install

# 2. 安装 RBAC 和部署
make k8s-rbac

# 3. 构建和部署服务
make k8s-deploy

# 4. 创建示例连接
make k8s-samples

# 5. 查看状态
kubectl get postgresqlconnection -n db-connect-demo
```

**特点**:
- ✅ 完整 K8s 集成
- ✅ 自动故障转移
- ✅ 支持多副本
- ✅ 企业级可靠性

---

## 快速示例

### 创建 PostgreSQL 连接

```bash
kubectl apply -f - <<EOF
apiVersion: db.connect.local/v1
kind: PostgreSQLConnection
metadata:
  name: my-postgres
  namespace: db-connect-demo
spec:
  host: "postgres.default.svc.cluster.local"
  port: 5432
  username: "postgres"
  password: "secret"
  database: "mydb"
EOF
```

### 查看连接状态

```bash
kubectl describe postgresqlconnection my-postgres -n db-connect-demo
```

### 执行查询

```bash
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"backend":"my-postgres","query":"SELECT 1"}'
```

---

## 文件结构

```
.
├── main.go                       # 统一的 Operator + API Server 入口
├── api/
│   └── v1/
│       └── *_types.go            # 5 种 CRD 类型定义
├── controllers/
│   ├── postgres_controller.go     # PostgreSQL Reconciler
│   ├── mysql_controller.go        # MySQL Reconciler
│   ├── sqlite_controller.go       # SQLite Reconciler
│   ├── kafka_controller.go        # Kafka Reconciler
│   └── solace_controller.go       # Solace Reconciler
├── lib/
│   ├── lib.go                    # Driver 接口
│   ├── postgres.go               # PostgreSQL 驱动
│   ├── mysql.go                  # MySQL 驱动
│   ├── sqlite.go                 # SQLite 驱动
│   ├── kafka.go                  # Kafka 驱动
│   ├── solace.go                 # Solace 驱动
│   └── mgr.go                    # 后端管理器
├── frontend/                     # React UI (可选)
├── config/                       # K8s 配置
│   ├── crd_crds.yaml            # CRD 定义
│   ├── rbac_rbac.yaml           # RBAC 配置
│   ├── manager_deployment.yaml   # Deployment 清单
│   └── samples_connections.yaml  # 示例
├── Dockerfile                    # 容器镜像
└── Makefile                      # 构建配置
```

---

## CRD 类型参考

### PostgreSQLConnection
```yaml
apiVersion: db.connect.local/v1
kind: PostgreSQLConnection
metadata:
  name: postgres-main
  namespace: db-connect-demo
spec:
  host: "localhost"
  port: 5432
  username: "postgres"
  password: "secret"
  database: "mydb"
```

### MySQLConnection
```yaml
apiVersion: db.connect.local/v1
kind: MySQLConnection
metadata:
  name: mysql-main
  namespace: db-connect-demo
spec:
  host: "localhost"
  port: 3306
  username: "root"
  password: "secret"
  database: "mydb"
```

### SQLiteConnection
```yaml
apiVersion: db.connect.local/v1
kind: SQLiteConnection
metadata:
  name: sqlite-local
  namespace: db-connect-demo
spec:
  filePath: "/data/test.db"
```

### KafkaConnection
```yaml
apiVersion: db.connect.local/v1
kind: KafkaConnection
metadata:
  name: kafka-prod
  namespace: db-connect-demo
spec:
  brokers:
    - "kafka-0:9092"
    - "kafka-1:9092"
```

### SolaceConnection
```yaml
apiVersion: db.connect.local/v1
kind: SolaceConnection
metadata:
  name: solace-msg
  namespace: db-connect-demo
spec:
  brokerURL: "tcp://solace-broker:55555"
  username: "admin"
  password: "secret"
```

---

## Makefile 目标

| 目标 | 说明 |
|------|------|
| `make build` | 构建二进制 |
| `make run` | 运行服务 (独立模式) |
| `make ui-build` | 构建前端 |
| `make clean` | 清理构建物 |
| `make cleanup` | 删除开发文件 |
| `make k8s-install` | 安装 CRD |
| `make k8s-rbac` | 安装 RBAC |
| `make k8s-deploy` | 完整部署 |
| `make k8s-samples` | 创建示例 |
| `make k8s-uninstall` | 卸载 |
| `make docker-build` | 构建镜像 |
| `make docker-push` | 推送镜像 |
| `make help` | 显示帮助 |

---

## API 端点

### GET /ping
返回所有连接的健康状态

```bash
curl http://localhost:8080/ping

# 返回示例
{
  "postgres-main": "ok",
  "mysql-main": "ok",
  "kafka-prod": "connection refused"
}
```

### POST /query
执行查询

```bash
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"backend":"postgres-main","query":"SELECT 1"}'

# 返回示例
{
  "rows": [{"?column?": "1"}]
}
```

### GET /ui
Web 界面 (可选，需要构建前端)

```
http://localhost:8080/ui
```

---

## 故障排查

### 连接失败

```bash
# 查看 CRD 状态
kubectl describe postgresqlconnection my-postgres -n db-connect-demo

# 查看错误详情在 Status.Error 字段
```

### Operator 日志

```bash
# 独立模式 (查看控制台输出)
make run

# K8s 模式
kubectl logs -f deployment/db-connect-operator -n db-connect-demo
```

### 健康检查

```bash
# 检查 API 服务是否运行
curl http://localhost:8080/ping
```

---

## 支持的驱动

| 驱动 | DSN 格式 | 备注 |
|------|----------|------|
| PostgreSQL | `host=... port=... user=... password=... dbname=...` | 关系数据库 |
| MySQL | `user:pass@tcp(host:port)/dbname` | 关系数据库 |
| SQLite | `file:path` 或 `:memory:` | 本地文件/内存 |
| Kafka | `broker1:9092,broker2:9092` | 消息队列 |
| Solace | `tcp://host:port` | MQTT 消息代理 |

---

## 清理

运行以下命令删除开发期间创建的临时文件：

```bash
make cleanup
```

这会删除以下文件：
- `operator_main.go` (已合并到 main.go)
- `init_dirs.go` (初始化工具)
- `test_operator.sh` (旧测试脚本)
- `deploy.sh` (单一用途脚本)
- `validate.go` (验证工具)
- `backends.json` (静态配置，已用 CRD 替代)
- `db-bench.exe` (旧二进制)

---

## 下一步

1. **快速体验**: `make run`
2. **K8s 部署**: `make k8s-deploy`
3. **创建连接**: `kubectl apply -f config_samples_connections.yaml`
4. **查看文档**: 
   - [QUICKSTART.md](QUICKSTART.md) - 快速开始
   - [OPERATOR_GUIDE.md](OPERATOR_GUIDE.md) - 详细指南

---

**简洁、统一、生产就绪！** 🚀
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

