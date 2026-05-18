# Kubernetes Operator 指南

db-connect-demo Operator 通过 Kubernetes Custom Resource Definition (CRD) 管理多个数据库连接实例。

## 架构概述

```
┌─────────────────────────────────────────────────────┐
│              Kubernetes Cluster                      │
├─────────────────────────────────────────────────────┤
│                                                      │
│  ┌──────────────────────────┐  ┌────────────────┐  │
│  │  Operator Pod            │  │  API Server    │  │
│  │  (db-connect-operator)   │  │  Pods (x2)     │  │
│  │                          │  │                │  │
│  │ • Watches CRDs           │  │ • HTTP /ping   │  │
│  │ • Manages connections    │  │ • HTTP /query  │  │
│  │ • Lifecycle management   │  │ • Queries lib  │  │
│  └──────────────────────────┘  └────────────────┘  │
│           ▲                             ▲           │
│           │                             │           │
│           └─────────────┬───────────────┘           │
│                         │                           │
│           ┌─────────────▼──────────────┐           │
│           │   lib package              │           │
│           │   (Connection drivers)     │           │
│           │                            │           │
│           │  PostgreSQL, MySQL,        │           │
│           │  SQLite, Kafka, Solace     │           │
│           └────────────────────────────┘           │
│                                                     │
│  ┌─────────────────────────────────────────────┐  │
│  │ CRD Resources (Custom Resource Definitions)  │  │
│  │                                              │  │
│  │ PostgreSQLConnection, MySQLConnection,       │  │
│  │ SQLiteConnection, KafkaConnection,           │  │
│  │ SolaceConnection                             │  │
│  └─────────────────────────────────────────────┘  │
│                                                      │
└─────────────────────────────────────────────────────┘
```

## 核心概念

### 1. CRD（Custom Resource Definition）
每种数据库驱动对应一种 CRD：

| CRD | 描述 | 示例 |
|-----|------|------|
| PostgreSQLConnection | PostgreSQL 连接 | 主生产库 |
| MySQLConnection | MySQL 连接 | 用户库 |
| SQLiteConnection | SQLite 文件数据库 | 本地缓存 |
| KafkaConnection | Kafka 消息队列 | 事件流 |
| SolaceConnection | Solace MQTT 消息代理 | 消息总线 |

### 2. Operator Controller
每个 Reconciler 监听对应 CRD 的创建/更新/删除事件：

- **PostgreSQLConnectionReconciler**: 监听 PostgreSQLConnection
- **MySQLConnectionReconciler**: 监听 MySQLConnection
- **SQLiteConnectionReconciler**: 监听 SQLiteConnection
- **KafkaConnectionReconciler**: 监听 KafkaConnection
- **SolaceConnectionReconciler**: 监听 SolaceConnection

当 CRD 变化时：
1. Operator 读取 CRD 的 `spec` 字段（连接参数）
2. 调用 lib 包中的驱动代码建立连接
3. 更新 CRD 的 `status` 字段（连接状态）

### 3. API Server
HTTP 微服务，提供 `/ping` 和 `/query` 端点：
- 从已注册的后端连接执行操作
- 返回健康状态和查询结果

## 部署流程

### 步骤 1：前置检查

```bash
# 检查 kubectl 是否可用
kubectl version

# 检查集群访问权限
kubectl get nodes
```

### 步骤 2：安装 CRD

```bash
# 安装所有 5 种 CRD
kubectl apply -f config_crd_crds.yaml

# 验证 CRD 安装
kubectl get crd | grep db.connect.local
```

### 步骤 3：安装 RBAC

```bash
# 创建 ServiceAccount、ClusterRole、ClusterRoleBinding
kubectl apply -f config_rbac_rbac.yaml

# 验证
kubectl get sa -n db-connect-demo
kubectl get clusterrole | grep db-connect
```

### 步骤 4：构建并推送镜像

```bash
# 构建 Operator 和 API Server 镜像
make docker-build

# 推送到仓库（如果不使用本地镜像）
make docker-push
```

### 步骤 5：部署 Operator 和 API Server

```bash
# 部署 Deployment 和 Service
kubectl apply -f config_manager_deployment.yaml

# 验证 Pod 运行状态
kubectl get pods -n db-connect-demo
kubectl logs -n db-connect-demo -f deployment/db-connect-operator
```

### 步骤 6：创建示例连接

```bash
# 应用示例 CRD
kubectl apply -f config_samples_connections.yaml

# 查看所有连接
kubectl get postgresqlconnection -n db-connect-demo
kubectl get mysqlconnection -n db-connect-demo
kubectl get kafkaconnection -n db-connect-demo

# 查看特定连接状态
kubectl describe postgresqlconnection postgres-primary -n db-connect-demo
```

## 常见操作

### 查看连接状态

```bash
# 列出所有 PostgreSQL 连接
kubectl get postgresqlconnection -n db-connect-demo

# 详细信息（包括 Status 和错误）
kubectl describe postgresqlconnection postgres-primary -n db-connect-demo

# 监控连接状态
kubectl get postgresqlconnection -n db-connect-demo --watch
```

### 创建新连接

```bash
# 编辑 YAML 文件
cat <<EOF | kubectl apply -f -
apiVersion: db.connect.local/v1
kind: PostgreSQLConnection
metadata:
  name: postgres-staging
  namespace: db-connect-demo
spec:
  host: "postgres-staging.example.com"
  port: 5432
  username: "admin"
  password: "secret123"
  database: "staging_db"
EOF
```

### 修改连接参数

```bash
# 编辑现有连接
kubectl edit postgresqlconnection postgres-primary -n db-connect-demo

# 或使用 patch 命令
kubectl patch postgresqlconnection postgres-primary -n db-connect-demo \
  --type merge -p '{"spec":{"port":5433}}'
```

### 删除连接

```bash
# 删除单个连接（自动触发清理逻辑）
kubectl delete postgresqlconnection postgres-primary -n db-connect-demo

# 删除所有 PostgreSQL 连接
kubectl delete postgresqlconnection --all -n db-connect-demo
```

## API Server 使用

### 转发端口

```bash
kubectl port-forward svc/db-connect-api 8080:8080 -n db-connect-demo
```

### 健康检查

```bash
curl http://localhost:8080/ping

# 返回示例
{
  "postgres-primary": "ok",
  "mysql-primary": "ok",
  "kafka-cluster": "connection failed: ...",
  "sqlite-local": "ok"
}
```

### 执行查询

```bash
# PostgreSQL 查询
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{
    "backend": "postgres-primary",
    "query": "SELECT COUNT(*) as total FROM users"
  }'

# 返回示例
{
  "rows": [
    {
      "total": "1000"
    }
  ]
}
```

## 故障排查

### Operator Pod 无法启动

```bash
# 查看 Operator 日志
kubectl logs -n db-connect-demo deployment/db-connect-operator --tail=50

# 查看事件
kubectl describe pod -n db-connect-demo -l app=db-connect-operator

# 检查 RBAC 权限
kubectl auth can-i list postgresqlconnection --as=system:serviceaccount:db-connect-demo:db-connect-operator
```

### 连接失败

```bash
# 查看连接 Status
kubectl describe postgresqlconnection <name> -n db-connect-demo

# 查看 Status.Error 字段，通常会显示：
# - "connection refused" - 网络连接问题
# - "invalid credentials" - 认证错误
# - "database not found" - 数据库不存在
```

### API Server 无法连接后端

```bash
# 检查 API Server 是否能看到连接
kubectl exec -n db-connect-demo deployment/db-connect-api -- \
  curl http://localhost:8080/ping

# 检查网络策略是否阻止通信
kubectl get networkpolicy -n db-connect-demo
```

## 扩展指南

### 添加新驱动

1. 在 `lib` 包中实现 `Driver` 接口
2. 创建新的 CRD 类型（如 `RedisConnection`）
3. 在 `api/v1/groupversion_info.go` 中注册 CRD
4. 创建对应的 Reconciler（如 `RedisConnectionReconciler`）
5. 在 `operator_main.go` 中注册 Reconciler
6. 更新 RBAC 配置添加新 CRD 权限

### 修改现有 CRD 字段

例如为 PostgreSQL 添加 SSL 支持：

1. 编辑 `api_v1_types.go`，在 `PostgreSQLConnectionSpec` 中添加字段：
   ```go
   SSLMode string `json:"sslMode,omitempty"`
   ```

2. 更新 `config_crd_crds.yaml` 中的 CRD schema

3. 更新 Reconciler 使用新字段构建 DSN

4. 应用更新：
   ```bash
   kubectl apply -f config_crd_crds.yaml
   ```

## 开发指南

### 本地运行 Operator

```bash
# 需要访问 Kubernetes 集群
export KUBECONFIG=~/.kube/config

make operator-run
```

### 测试新的 Reconciler

1. 创建测试用的 CRD 实例
2. 查看 Operator 日志
3. 验证 Status 字段更新

```bash
kubectl logs -f deployment/db-connect-operator -n db-connect-demo
```

### 修改后重新部署

```bash
# 重建镜像
make docker-build

# 重新部署（会重启 Pod）
kubectl rollout restart deployment/db-connect-operator -n db-connect-demo
kubectl rollout restart deployment/db-connect-api -n db-connect-demo
```

## 生产部署建议

1. **高可用性**：设置 Operator 副本数 > 1，使用 leader election

2. **监控**：
   - 监控 Operator metrics（:8081）
   - 监控 API Server 健康（/ping）
   - 监控 CRD resource 状态

3. **日志**：
   - 集成 ELK/Loki 收集 Operator 日志
   - 监控 "ConnectionFailed" 条件

4. **备份**：
   - 备份 CRD 资源定义
   - 使用 ConfigMap/Secret 存储敏感配置

5. **更新策略**：
   - 测试环境先行
   - 使用蓝绿部署
   - 记录 CRD 变更历史

## 参考文档

- [Kubernetes Operator Pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)
- [controller-runtime](https://pkg.go.dev/sigs.k8s.io/controller-runtime)
- [CRD 最佳实践](https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/)
