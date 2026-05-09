# Cloud 管理小助手

一个使用 Golang + Wails + Vue 开发的云资源管理桌面工具。

后端使用 Golang，负责云厂商 SDK 调用、账号鉴权、资源查询、操作封装和安全校验。

前端使用 Vue 页面，负责账号输入、资源展示、环境区分、操作确认和结果展示。

---

## 当前版本

当前版本：`v0.0.1`

`v0.0.1` 是项目启动版本，只完成最小可运行骨架，不追求功能完整。

---

## 项目定位

Cloud 管理小助手是一个面向运维人员的桌面 GUI 工具。

GUI 是主要入口。

CLI 可以保留，但只作为调试入口、自动化入口和后续脚本化入口。

---

## 技术栈

- **后端**: Golang
- **GUI 框架**: Wails v2
- **前端**: Vue 3 + Vite
- **云厂商**: 阿里云 ECS SDK

---

## 项目结构

```
cloud-manage/
├── main.go                     # Wails 主入口 (GUI)
├── app.go                      # Wails 绑定层，暴露方法给前端
├── wails.json                  # Wails 配置
├── go.mod
├── go.sum
├── cmd/
│   └── cli/
│       └── main.go             # CLI 调试入口（非主入口）
├── service/
│   └── ecs.go                  # 业务编排层
├── provider/
│   └── aliyun/
│       └── ecs.go              # 阿里云 ECS SDK 封装
├── security/
│   └── sanitize.go             # 敏感信息脱敏
├── frontend/
│   ├── index.html
│   ├── package.json
│   ├── vite.config.js
│   ├── src/
│   │   ├── main.js
│   │   └── App.vue             # 主窗口组件
│   └── wailsjs/
│       └── go/
│           └── main/
│               └── App.js      # Wails 前端绑定（自动生成）
├── scripts/
│   ├── dev.sh                  # 开发启动脚本
│   ├── build.sh                # 构建脚本
│   └── test.sh                 # 测试脚本
├── README.md
├── ROADMAP.md
└── CHANGELOG.md
```

---

## 架构设计

```
┌─────────────────────────────────────────────────┐
│                  frontend (Vue)                  │
│            页面展示 + 用户交互                     │
└──────────────────────┬──────────────────────────┘
                       │ Wails 调用
┌──────────────────────▼──────────────────────────┐
│                  app.go (绑定层)                  │
│          暴露给前端调用的方法                       │
└──────────────────────┬──────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────┐
│               service 层 (业务编排)               │
│              参数校验 + 结果组装                    │
└──────────────────────┬──────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────┐
│           provider/aliyun 层 (SDK 封装)           │
│           阿里云 ECS SDK 调用                      │
└──────────────────────┬──────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────┐
│              security 层 (安全)                   │
│          敏感信息脱敏 + 错误信息清洗                 │
└─────────────────────────────────────────────────┘
```

---

## 快速开始

### 环境要求

- Go 1.21+
- Node.js 18+
- npm 9+
- Wails v2 CLI (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)
- Linux: `libgtk-3-dev`, `libwebkit2gtk-4.0-dev`

### 安装依赖

```bash
# 安装 Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Linux 系统依赖
sudo apt install libgtk-3-dev libwebkit2gtk-4.0-dev

# 前端依赖
cd frontend && npm install
```

### 开发模式

```bash
./scripts/dev.sh
```

### 构建

```bash
./scripts/build.sh
```

构建产物位于 `build/bin/` 目录。

### 测试

```bash
./scripts/test.sh
```

---

## 功能说明

### v0.0.1 支持的功能

- 从 GUI 输入阿里云 AK/SK
- 从 GUI 输入 Region
- 选择环境类型（dev / pre / prod）
- 点击按钮查询 ECS 实例列表
- 表格展示 ECS 实例信息（InstanceId、InstanceName、Status、RegionId、ZoneId、PublicIp、PrivateIp、CreationTime）
- 清空输入 / 清空结果
- 错误提示区域

### 安全设计

- 不硬编码 AK/SK
- 不日志打印 AK/SK
- 不保存 AK/SK 到本地文件
- 不使用 localStorage / sessionStorage 保存 AK/SK
- 错误信息自动脱敏
- 仅支持只读查询，不支持创建/删除/重启/停止 ECS

---

## Roadmap

详见 [ROADMAP.md](./ROADMAP.md)

## Changelog

详见 [CHANGELOG.md](./CHANGELOG.md)
