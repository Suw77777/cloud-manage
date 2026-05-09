# Changelog

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
