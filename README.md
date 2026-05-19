# Cloud Manage

阿里云资源管理工具，支持 GUI 和 CLI 两种模式，覆盖 ECS、CMS、SLS、OSS 服务。

## 构建

### 前置依赖

- Go 1.21+
- Node.js 18+
- Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

### 构建命令

```bash
# 安装前端依赖
cd frontend && npm install && cd ..

# 开发模式（GUI 热更新）
wails dev

# 构建生产版本（GUI + CLI）
wails build

# 仅构建 CLI
go build -o cloud-manage ./cmd/cli/
```

## 使用

### GUI 模式

```bash
./build/bin/cloud-manage
```

### CLI 模式

```bash
# 配置凭证
export CLOUD_ACCESS_KEY_ID=your-key-id
export CLOUD_ACCESS_KEY_SECRET=your-key-secret

# ECS 实例管理
./cloud-manage ecs list
./cloud-manage ecs detail i-xxx
./cloud-manage ecs start i-xxx
./cloud-manage ecs stop i-xxx [--force]

# 云监控（分步查询）
./cloud-manage cms products          # 列出支持的云产品
./cloud-manage cms metrics ecs       # 列出 ECS 实例
./cloud-manage cms query i-xxx       # 查询实例监控数据

# 日志服务
./cloud-manage sls logstores <project>
./cloud-manage sls logs <project> <logstore> --query "level: ERROR"

# 对象存储
./cloud-manage oss buckets
./cloud-manage oss objects <bucket> --prefix logs/ --max 100
```

## 通用参数

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `-id` | AccessKey ID | 环境变量 `CLOUD_ACCESS_KEY_ID` |
| `-secret` | AccessKey Secret | 环境变量 `CLOUD_ACCESS_KEY_SECRET` |
| `-region` | 地域 | `cn-hangzhou` |
| `-json` | JSON 输出 | `false` |

## 支持的云产品（CMS）

| 产品 ID | 产品名称 |
|---------|----------|
| ecs | 云服务器 ECS |
| rds | 云数据库 RDS |
| slb | 负载均衡 SLB |
| redis | 云数据库 Redis |

## License

MIT
