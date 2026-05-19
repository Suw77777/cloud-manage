# Cloud Manage

阿里云资源管理工具，支持 GUI 和 CLI 两种模式，覆盖 ECS、CMS、SLS、OSS 服务。

## 构建

### 前置依赖

- Go 1.21+
- Node.js 18+
- Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- Linux GUI 依赖: `libgtk-3-dev libwebkit2gtk-4.0-dev`

### 构建命令

```bash
# 安装前端依赖
cd frontend && npm install && cd ..

# 开发模式（GUI 热更新）
wails dev

# 构建生产版本（GUI + CLI）
wails build

# 仅构建 CLI（无 GUI 依赖）
go build -o cloud-manage ./cmd/cli/
```

### 打包分发

```bash
# 安装 GUI 运行时依赖
./scripts/install-deps.sh

# 构建 AppImage（独立可执行文件，包含依赖）
./scripts/build-appimage.sh
```

## 使用

### GUI 模式

```bash
# 需要先安装依赖
./scripts/install-deps.sh

# 运行
./build/bin/cloud-manage
```

### CLI 模式（推荐用于服务器/无 GUI 环境）

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

### AppImage 模式

```bash
# 直接运行，无需安装依赖
./cloud-manage-0.0.12-x86_64.AppImage
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

## Docker 构建

```bash
# 构建镜像（使用国内镜像）
docker build -t cloud-manage .

# 运行 CLI
docker run --rm cloud-manage ./cloud-manage ecs list
```

## License

MIT
