# 快速开始 - 5 分钟部署指南

## 两种部署模式

### 🚀 方案 A：独立 HTTP 服务（适合快速测试）

**无需 Kubernetes，开箱即用！**

```bash
# 1. 构建
make build-cli
make ui-build

# 2. 运行
make run-service

# 3. 访问
# 浏览器打开: http://localhost:8080/ui
# API: http://localhost:8080/ping
```

**特点**：
- ✅ 快速启动
- ✅ 无依赖
- ✅ 适合开发测试

---

### ☸️  方案 B：Kubernetes Operator（推荐生产）

**完整的 K8s 集成，可动态管理连接！**

#### 前置检查

```bash
# 1. 确保有 Kubernetes 集群
kubectl version

# 2. 确保能访问集群
kubectl get nodes
```

#### 部署（一键）

```bash
# 使用自动部署脚本（推荐）
bash deploy.sh

# 或手动部署
make k8s-deploy
make k8s-samples
```

#### 验证

```bash
# 1. 查看 Pod 状态
kubectl get pods -n db-connect-demo

# 2. 转发 API 端口
kubectl port-forward svc/db-connect-api 8080:8080 -n db-connect-demo

# 3. 测试健康检查
curl http://localhost:8080/ping
```

**特点**：
- ✅ 热更新（无需重启）
- ✅ 原生 K8s 工具支持
- ✅ 自动故障转移
- ✅ 高可用部署

---

## 常见操作

### 创建数据库连接

```bash
# PostgreSQL
kubectl apply -f - <<EOF
apiVersion: db.connect.local/v1
kind: PostgreSQLConnection
metadata:
  name: postgres-main
  namespace: db-connect-demo
spec:
  host: "postgres.default.svc.cluster.local"
  port: 5432
  username: "postgres"
  password: "password123"
  database: "mydb"
EOF

# MySQL
kubectl apply -f - <<EOF
apiVersion: db.connect.local/v1
kind: MySQLConnection
metadata:
  name: mysql-main
  namespace: db-connect-demo
spec:
  host: "mysql.default.svc.cluster.local"
  port: 3306
  username: "root"
  password: "password123"
  database: "mydb"
EOF
```

### 查看连接状态

```bash
# 列表视图
kubectl get postgresqlconnection -n db-connect-demo

# 详细信息
kubectl describe postgresqlconnection postgres-main -n db-connect-demo

# 实时监控
kubectl get postgresqlconnection -n db-connect-demo -w
```

### 执行查询

```bash
# 通过 API 执行查询
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{
    "backend": "postgres-main",
    "query": "SELECT VERSION()"
  }'
```

### 修改连接

```bash
# 编辑连接参数
kubectl edit postgresqlconnection postgres-main -n db-connect-demo

# 或者使用 patch
kubectl patch postgresqlconnection postgres-main -n db-connect-demo \
  --type merge -p '{"spec":{"port":5433}}'
```

### 删除连接

```bash
kubectl delete postgresqlconnection postgres-main -n db-connect-demo
```

---

## Makefile 参考

```bash
# 构建
make build                  # 构建所有二进制
make build-cli              # 只构建 CLI
make build-operator         # 只构建 Operator
make build-api              # 只构建 API Server
make ui-build               # 构建前端

# 独立服务模式
make run-service            # 运行 HTTP 服务

# Operator 模式
make k8s-install            # 安装 CRD
make k8s-deploy             # 完整部署
make k8s-samples            # 创建示例
make k8s-uninstall          # 卸载

# Docker
make docker-build           # 构建镜像
make docker-push            # 推送镜像
```

---

## 故障排查

### Operator 无法启动

```bash
kubectl logs -n db-connect-demo deployment/db-connect-operator --tail=50
```

### 连接失败

```bash
kubectl describe postgresqlconnection <name> -n db-connect-demo
# 查看 Status.Error 字段
```

### API Server 无响应

```bash
kubectl logs -n db-connect-demo deployment/db-connect-api
# 或检查 Pod 状态
kubectl get pods -n db-connect-demo -l app=db-connect-api
```

---

## 相关文档

- 📖 [README.md](README.md) - 完整项目说明
- 📋 [OPERATOR_GUIDE.md](OPERATOR_GUIDE.md) - 详细 Operator 指南
- 📊 [TRANSFORMATION_SUMMARY.md](TRANSFORMATION_SUMMARY.md) - 改造总结

---

## 支持的驱动

| 驱动 | DSN 示例 | 用途 |
|------|----------|------|
| PostgreSQL | `host=localhost port=5432 user=... password=... dbname=...` | 关系数据库 |
| MySQL | `user:pass@tcp(localhost:3306)/dbname` | 关系数据库 |
| SQLite | `file::memory:?cache=shared` | 本地数据库 |
| Kafka | `localhost:9092,broker2:9092` | 消息队列 |
| Solace | `tcp://broker-host:1883` | MQTT 消息 |

---

**选择一种模式开始部署吧！** 🎉
