# Changelog

## v0.0.3 (2026-05-09)

### 功能

- ECS 实例详情页：点击 InstanceId 弹窗展示 CPU、Memory、ImageId、VPC、安全组、过期时间等
- 启动 / 停止 / 重启 ECS 实例：每个操作均需二次确认弹窗
- 生产环境强制二次确认：红色警告文案 + 特殊按钮文案"确认执行（生产环境）"
- 操作日志：前端内存记录，最多 50 条，支持清空
- 实例列表新增操作列：Running 状态显示"停止/重启"，Stopped 状态显示"启动"

### 后端

- provider/aliyun 新增 DescribeInstanceDetail、StartInstance、StopInstance、RebootInstance
- service 新增 GetInstanceDetail、StartInstance、StopInstance、RebootInstance
- app.go 新增 GetECSDetail、StartECS、StopECS、RebootECS
- 放开写操作约束（启动/停止/重启，需二次确认）

### 文档

- 更新 ROADMAP.md、CHANGELOG.md
- 版本号更新为 v0.0.3
- go test ./... 全部通过

## v0.0.2 (2026-05-09)

### 功能

- 支持多 Region 选择：18 个常用阿里云 Region 复选框（国内 + 海外）
- 支持多 Region 并发查询：每个 Region 独立 goroutine，错误不互相影响
- 结果按 Region 分组展示：每个 Region 独立表格，带实例计数和错误提示
- 新增 `QueryECSMultiRegion` 后端方法
- 新增 `ListInstancesMultiRegion` service 方法
- 新增 `ECSInstanceAdapter` 适配类型，app.go 不直接依赖 provider 层
- 版本号更新为 v0.0.2

### 工程

- 新增多 Region 测试用例
- 更新 Wails 前端绑定（JS + TS 类型定义）
- go test ./... 全部通过

## v0.0.1 (2026-05-09)

### 项目初始化

- 使用 Golang + Wails v2 + Vue 3 初始化项目
- 实现 GUI 主窗口
- 实现后端分层架构：app.go → service → provider/aliyun → security

### 功能

- 支持从 GUI 输入阿里云 AccessKey ID / Secret
- 支持从 GUI 输入 Region（默认 cn-hangzhou）
- 支持选择环境类型（dev / pre / prod）
- 点击按钮调用 Go 后端查询 ECS 实例列表
- 表格展示 ECS 实例：InstanceId, InstanceName, Status, RegionId, ZoneId, PublicIp, PrivateIp, CreationTime
- 支持清空输入 / 清空结果
- 错误提示区域

### 安全

- 不硬编码 AK/SK
- 不日志打印 AK/SK
- 不保存 AK/SK 到本地文件
- 不使用 localStorage / sessionStorage 保存 AK/SK
- 错误信息自动脱敏（security 层）
- 仅支持只读查询

### 工程

- 增加 scripts/dev.sh、scripts/build.sh、scripts/test.sh
- 增加 README.md、ROADMAP.md、CHANGELOG.md
- go test ./... 通过
