# Cloud 管理小助手

一个使用 Golang + Wails + Vue 开发的云资源管理工具。

后端使用 Golang，负责云厂商 SDK 调用、账号鉴权、资源查询、操作封装和安全校验。

前端使用 Vue 页面，负责账号输入、资源展示、环境区分、操作确认和结果展示。

---

## 当前版本

当前版本：`v0.0.10`

`v0.0.10` 分离 GUI 和 CLI，支持多平台交叉编译。

---

## 项目定位

Cloud 管理小助手是一个面向运维人员的云资源管理工具。

**双模式支持**：
- **GUI 桌面应用**：基于 Wails，适合有图形环境的系统
- **CLI 命令行工具**：纯 Go 编译，支持 Linux/macOS/Windows，无 GUI 依赖

**支持平台**：
| 平台 | CLI | GUI |
|------|-----|-----|
| Linux (amd64/arm64) | ✅ | ✅ (需 GTK/WebKit) |
| macOS (Intel/Apple Silicon) | ✅ | ✅ |
| Windows (amd64) | ✅ | ✅ (需 WebView2) |

---

## 技术栈

- **后端**: Golang
- **GUI 框架**: Wails v2
- **前端**: Vue 3 + Vite
- **云厂商**: 阿里云（ECS、CMS、SLS、OSS）

---

## 项目结构

```
cloud-manage/
├── main.go                         # GUI 入口（Wails）
├── app.go                          # Wails 绑定层，暴露方法给前端
├── cmd/
│   └── cli/
│       └── main.go                 # CLI 入口（纯 Go，无 GUI 依赖）
├── Dockerfile                      # Docker 构建环境（Ubuntu 22.04）
├── wails.json                      # Wails 配置
├── go.mod
├── go.sum
├── service/
│   ├── ecs.go                      # ECS 业务编排层
│   ├── ecs_test.go                 # ECS 单元测试
│   ├── cms.go                      # CMS 业务编排层
│   ├── sls.go                      # SLS 业务编排层
│   └── oss.go                      # OSS 业务编排层
├── provider/
│   └── aliyun/
│       ├── ecs.go                  # 阿里云 ECS SDK 封装
│       ├── cms.go                  # 阿里云 CMS SDK 封装
│       ├── sls.go                  # 阿里云 SLS SDK 封装
│       └── oss.go                  # 阿里云 OSS SDK 封装
├── security/
│   ├── sanitize.go                 # 敏感信息脱敏
│   └── sanitize_test.go            # 脱敏单元测试
├── frontend/
│   ├── index.html
│   ├── package.json
│   ├── vite.config.js
│   ├── src/
│   │   ├── main.js
│   │   ├── App.vue                 # 主窗口入口
│   │   ├── assets/
│   │   │   └── main.css            # 全局样式
│   │   ├── components/
│   │   │   ├── EcsResultTable.vue  # ECS 实例表格
│   │   │   ├── InstanceDetailModal.vue  # ECS 详情弹窗
│   │   │   ├── ConfirmDialog.vue   # 确认弹窗
│   │   │   ├── OperationLog.vue    # 操作日志
│   │   │   ├── CmsMonitor.vue      # 云监控组件
│   │   │   ├── SlsQuery.vue        # SLS 日志查询组件
│   │   │   ├── OssBrowser.vue      # OSS 浏览器组件
│   │   │   ├── PaginationBar.vue   # 分页组件
│   │   │   └── VirtualScroller.vue # 虚拟滚动组件
│   │   └── composables/
│   │       ├── useECS.js           # ECS 查询/操作逻辑
│   │       ├── useCMS.js           # CMS 监控逻辑
│   │       ├── useSLS.js           # SLS 日志查询逻辑
│   │       ├── useOSS.js           # OSS 浏览逻辑
│   │       └── usePagination.js    # 分页逻辑
│   └── wailsjs/
│       └── go/
│           └── main/
│               └── App.js          # Wails 前端绑定（自动生成）
├── scripts/
│   ├── dev.sh                      # 开发启动脚本
│   ├── build.sh                    # GUI 构建脚本
│   ├── build-cli.sh                # CLI 交叉编译脚本（多平台）
│   ├── build-gui-docker.sh         # GUI Docker 构建脚本
│   ├── build-all.sh                # 一键构建 CLI + GUI
│   └── test.sh                     # 测试脚本
├── README.md
├── ROADMAP.md
└── CHANGELOG.md
```

---

## 架构设计

```
┌─────────────────────────────────────────────────┐
│                  frontend (Vue)                  │
│      ECS / CMS / SLS / OSS 四大模块              │
└──────────────────────┬──────────────────────────┘
                       │ Wails 调用
┌──────────────────────▼──────────────────────────┐
│                  app.go (绑定层)                  │
│          暴露给前端调用的方法                       │
└──────────────────────┬──────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────┐
│               service 层 (业务编排)               │
│         ecs / cms / sls / oss 四个服务            │
└──────────────────────┬──────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────┐
│           provider/aliyun 层 (SDK 封装)           │
│       ECS / CMS / SLS / OSS 四个 Provider        │
└──────────────────────┬──────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────┐
│              security 层 (安全)                   │
│          敏感信息脱敏 + 错误信息清洗                 │
└─────────────────────────────────────────────────┘
```

### 工程约束

- **模块化原则**：新增功能以组件 / composable / service / provider 方式扩展
- **单一职责**：每个文件只负责一个功能模块
- **分层架构**：严格遵循 frontend → app.go → service → provider → security 的调用链

---

## 快速开始

### 方式一：直接使用 CLI（推荐）

下载对应平台的 CLI 二进制文件即可使用，无需安装依赖：

```bash
# Linux amd64
wget https://github.com/Suw77777/cloud-manage/releases/latest/download/cloud-cli-linux-amd64
chmod +x cloud-cli-linux-amd64
mv cloud-cli-linux-amd64 cloud-cli

# macOS (Apple Silicon)
wget https://github.com/Suw77777/cloud-manage/releases/latest/download/cloud-cli-darwin-arm64
chmod +x cloud-cli-darwin-arm64
mv cloud-cli-darwin-arm64 cloud-cli

# Windows (PowerShell)
Invoke-WebRequest -Uri "https://github.com/Suw77777/cloud-manage/releases/latest/download/cloud-cli-windows-amd64.exe" -OutFile "cloud-cli.exe"
```

使用示例：
```bash
# 设置凭证
export CLOUD_ACCESS_KEY_ID=your-key-id
export CLOUD_ACCESS_KEY_SECRET=your-key-secret

# 列出 ECS 实例
./cloud-cli ecs list

# 查看实例详情
./cloud-cli ecs detail i-xxx

# 查询监控指标
./cloud-cli cms metrics i-xxx
```

### 方式二：从源码构建

#### 构建 CLI（所有平台）

```bash
# 交叉编译 Linux/macOS/Windows
./scripts/build-cli.sh

# 输出在 release/ 目录
ls release/
```

#### 构建 GUI

**Ubuntu 22.04 / Linux Mint 21 及以下：**
```bash
# 安装依赖
sudo apt install libgtk-3-dev libwebkit2gtk-4.0-dev

# 安装 Wails
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# 构建
./scripts/build.sh
```

**Ubuntu 24.04 / Linux Mint 22 及以上：**
```bash
# 使用 Docker 构建（解决 libwebkit2gtk-4.1 兼容问题）
./scripts/build-gui-docker.sh
```

#### 一键构建全部

```bash
./scripts/build-all.sh
```

### 开发模式

```bash
./scripts/dev.sh
```

### 测试

```bash
./scripts/test.sh
```

---

## 功能说明

### ECS 实例管理

- 多 Region 选择（18 个常用 Region 复选框，支持国内 + 海外）
- 多 Region 并发查询，结果按 Region 分组展示
- ECS 实例详情页（CPU、Memory、ImageId、VPC、安全组、过期时间等）
- 启动 / 停止 / 重启 ECS 实例（二次确认弹窗）
- 生产环境操作强制二次确认（红色警告 + 特殊按钮文案）
- 操作日志记录（前端内存，最多 50 条）

### 云监控 (CMS)

- 从 ECS 结果选择实例查看监控数据
- 监控指标：CPU 使用率、内存使用率、磁盘读写 BPS、网络流量
- 进度条可视化展示（颜色阈值：正常/警告/危险）
- 支持多实例批量查询

### 日志服务 (SLS)

- 输入 SLS Project 名称获取 Logstore 列表
- 日志查询：时间范围选择（15 分钟到 7 天）、最大行数设置
- 自定义查询语句支持
- 日志结果表格展示
- 分页显示（默认每页 50 条）
- 流式查询（边查边显示）

### 对象存储 (OSS)

- 列出所有 Bucket
- 浏览 Bucket 内的对象和目录
- 目录导航：面包屑、返回上级、返回根目录
- 对象表格展示：图标、名称、大小、修改时间、存储类型

### CLI 命令行

```bash
# ECS 实例管理
cloud-cli ecs list                            # 列出实例
cloud-cli ecs detail <id>                     # 查看详情
cloud-cli ecs start/stop/reboot <id>          # 操作实例

# 云监控
cloud-cli cms metrics <id>                   # 查询监控指标

# 日志服务
cloud-cli sls logstores <project>            # 列出 Logstore
cloud-cli sls logs <project> <logstore>      # 查询日志

# 对象存储
cloud-cli oss buckets                        # 列出 Bucket
cloud-cli oss objects <bucket>               # 列出对象

# 所有 Region 查询
cloud-cli -region all ecs list
```

**认证方式：**
```bash
# 方式一：环境变量（推荐）
export CLOUD_ACCESS_KEY_ID=your-key-id
export CLOUD_ACCESS_KEY_SECRET=your-key-secret

# 方式二：命令行参数
cloud-cli -id your-key-id -secret your-key-secret ecs list
```

**输出格式：**
```bash
# JSON 格式输出
cloud-cli -json ecs list
```

---

## 安全设计

- 不硬编码 AK/SK
- 不日志打印 AK/SK
- 不保存 AK/SK 到本地文件
- 不使用 localStorage / sessionStorage 保存 AK/SK
- 错误信息自动脱敏
- ECS 写操作（启动/停止/重启）需二次确认
- 生产环境操作强制二次确认

---

## Roadmap

详见 [ROADMAP.md](./ROADMAP.md)

## Changelog

详见 [CHANGELOG.md](./CHANGELOG.md)

---

## Contributors

感谢以下贡献者对本项目的支持：

| 贡献者 | 角色 | 贡献内容 |
|--------|------|----------|
| [Suw77777](https://github.com/Suw77777) | 项目发起人 & 主要开发者 | 项目架构设计、需求定义、代码审查、测试验证 |
| [MiMo-V2.5-Pro](https://platform.xiaomimimo.com) | AI 编程助手 (小米) | 代码实现、功能开发、Bug 修复、文档编写 |

### AI 贡献说明

本项目 v0.0.1 ~ v0.0.9 版本的核心代码由小米 MiMo-V2.5-Pro AI 模型辅助生成，包括：

- 后端架构：Go 服务层、Provider 层、安全层实现
- 前端组件：Vue 3 组件、Composables、样式开发
- 功能模块：ECS、CMS、SLS、OSS 四大云服务模块
- 工程化：构建脚本、测试用例、项目文档
