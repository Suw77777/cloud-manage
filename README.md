# Cloud Manage

阿里云资源管理工具，支持 GUI、TUI 和 CLI 三种模式，覆盖 ECS、CMS、SLS、OSS、VPC、SLB 服务。

## 构建

### 前置依赖

- Go 1.21+
- Node.js 18+
- Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- Linux GUI 依赖（仅 GUI 模式需要，TUI/CLI 无需安装）:

```bash
# 一键安装（自动检测发行版和 webkit 版本）
./scripts/install-deps.sh

# 或手动安装:
# Ubuntu 22.04 / Linux Mint 21:
sudo apt install libgtk-3-dev libwebkit2gtk-4.0-dev
# Ubuntu 24.04 / Linux Mint 22+:
sudo apt install libgtk-3-dev libwebkit2gtk-4.1-dev
# Arch:
sudo pacman -S gtk3 webkit2gtk
# Fedora:
sudo dnf install gtk3-devel webkit2gtk4.1-devel
```

> **不需要 GUI？** CLI 和 TUI 模式零系统依赖，直接 `go build -o cloud-manage .` 即可。

### 构建命令

```bash
# 安装前端依赖
cd frontend && npm install && cd ..

# 开发模式（GUI 热更新）
wails dev

# 构建生产版本（GUI + CLI）
wails build

# 仅构建命令行版本（无 GUI 依赖）
go build -o cloud-manage .
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

### TUI 模式（终端界面）

```bash
# 强制 TUI 模式
./cloud-manage --tui

# 功能
# - 登录界面输入凭证
# - Tab 切换：ECS、CMS、SLS、OSS、VPC、SLB
# - 方向键导航，Enter 查看详情
```

### CLI 模式（推荐用于服务器/无 GUI 环境）

```bash
# 方式一：通过环境变量配置凭证（推荐）
export CLOUD_ACCESS_KEY_ID=your-key-id
export CLOUD_ACCESS_KEY_SECRET=your-key-secret
./cloud-manage ecs list

# 方式二：通过 flag 传入凭证（flag 必须放在子命令前面）
./cloud-manage -id your-key-id -secret your-key-secret ecs list
./cloud-manage -id your-key-id -secret your-key-secret -region cn-beijing vpc list

# 方式三：通过配置文件（需先初始化）
./cloud-manage config init              # 初始化配置文件
./cloud-manage config add prod --save   # 添加账号（加密保存凭证）
./cloud-manage -profile prod ecs list   # 使用指定账号

# ECS 实例管理
./cloud-manage ecs list
./cloud-manage ecs detail i-xxx
./cloud-manage ecs start i-xxx
./cloud-manage ecs stop i-xxx [--force]

# 云监控
./cloud-manage cms products             # 列出支持的云产品
./cloud-manage cms metrics i-xxx        # 查询实例监控数据

# 日志服务
./cloud-manage sls logstores <project>
./cloud-manage sls logs <project> <logstore> --query "level: ERROR"
./cloud-manage sls export <project> <logstore> --format csv --output logs.csv

# 对象存储
./cloud-manage oss buckets
./cloud-manage oss objects <bucket> --prefix logs/ --max 100

# VPC 网络
./cloud-manage vpc list
./cloud-manage vpc detail <vpc-id>
./cloud-manage vpc vswitches <vpc-id>

# 负载均衡
./cloud-manage slb list
./cloud-manage slb detail <slb-id>
./cloud-manage slb listeners <slb-id>

# JSON 格式输出
./cloud-manage -json ecs list
```

### AppImage 模式

```bash
# 直接运行，无需安装依赖
./cloud-manage-v0.1.0-x86_64.AppImage
```

## 通用参数

> **重要**: 所有 flag 参数必须放在子命令**之前**，否则会被忽略。

```bash
# 正确 - flag 在子命令前面
./cloud-manage -id xxx -secret xxx -region cn-shanghai ecs list
./cloud-manage -json -region cn-beijing vpc list

# 错误 - flag 在子命令后面，会被忽略
./cloud-manage ecs list -id xxx -secret xxx
```

### 全局参数

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `-id` | AccessKey ID | 环境变量 `CLOUD_ACCESS_KEY_ID` |
| `-secret` | AccessKey Secret | 环境变量 `CLOUD_ACCESS_KEY_SECRET` |
| `-region` | 地域 | `cn-hangzhou` |
| `-json` | JSON 格式输出 | `false` |
| `-profile` | 使用指定配置账号 | 配置文件中的 `current_profile` |

### 模式参数

| 参数 | 说明 |
|------|------|
| `--gui` | 强制 GUI 模式（需要图形环境） |
| `--tui` | 强制 TUI 模式（终端界面） |
| `--cli` | 强制 CLI 模式 |

> 不带模式参数时自动检测：有图形环境启动 GUI，检测到子命令进入 CLI，否则进入 TUI。

### 参数优先级

命令行参数 > 环境变量 > `--profile` 指定的配置 > `current_profile` 配置 > 默认值

## 支持的服务

| 服务 | 功能 |
|------|------|
| ECS | 云服务器实例管理（列表、详情、启动、停止、重启） |
| CMS | 云监控指标查询（CPU、内存、磁盘、网络） |
| SLS | 日志服务查询（Logstore 列表、日志查询） |
| OSS | 对象存储管理（Bucket 列表、对象列表） |
| VPC | VPC 网络管理（VPC 列表、详情、VSwitch 列表） |
| SLB | 负载均衡管理（SLB 列表、详情、监听器列表） |

## Docker 构建

```bash
# 构建镜像（使用国内镜像）
docker build -t cloud-manage .

# 运行 CLI
docker run --rm cloud-manage ./cloud-manage ecs list
```

## License

MIT
