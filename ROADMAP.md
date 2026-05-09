# Roadmap

## v0.0.1 - 项目启动

**目标**: 完成 Golang + Wails + Vue GUI 最小可运行版本，支持阿里云 ECS 实例列表查询。

- [x] 使用 Golang + Wails 初始化项目
- [x] GUI 作为主入口
- [x] 后端使用 Golang
- [x] 前端使用 Vue
- [x] 实现主窗口
- [x] AccessKey ID 输入框
- [x] AccessKey Secret 密码输入框
- [x] Region 输入框
- [x] 环境选择：dev / pre / prod
- [x] 查询 ECS 按钮
- [x] 清空输入按钮
- [x] 清空结果按钮
- [x] ECS 实例结果表格
- [x] 错误提示区域
- [x] Go 后端调用阿里云 ECS SDK 查询实例列表
- [x] 表格展示 InstanceId, InstanceName, Status, RegionId, ZoneId, PublicIp, PrivateIp, CreationTime
- [x] security 层敏感信息脱敏
- [x] scripts/dev.sh, build.sh, test.sh
- [x] README.md, ROADMAP.md, CHANGELOG.md
- [x] go test ./... 通过

---

## v0.0.2 - 多 Region 支持

**目标**: 支持同时查询多个 Region 的 ECS 实例。

- [x] 支持多 Region 下拉选择（18 个常用 Region 复选框）
- [x] 支持多 Region 并发查询（goroutine per region）
- [x] 结果按 Region 分组展示

---

## v0.0.3 - ECS 详情与操作

**目标**: 支持查看 ECS 实例详情和基本操作。

- [x] ECS 实例详情页（弹窗展示 CPU/Memory/Image/VPC/SG 等）
- [x] 启动 / 停止 / 重启 ECS 实例（二次确认）
- [x] 生产环境操作强制二次确认（红色警告 + 特殊按钮文案）
- [x] 操作结果日志记录（前端内存，最多 50 条）

---

## v0.0.4 - SLS 日志查询

**目标**: 支持阿里云 SLS 日志查询。

- [ ] SLS Logstore 列表
- [ ] 日志查询输入框
- [ ] 日志结果展示
- [ ] 日志下载

---

## v0.0.5 - OSS / NAS 支持

**目标**: 支持阿里云 OSS 和 NAS 资源查询。

- [ ] OSS Bucket 列表
- [ ] OSS Object 浏览
- [ ] NAS 文件系统列表

---

## v0.1.0 - 多云支持

**目标**: 支持腾讯云、华为云。

- [ ] 腾讯云 ECS 查询
- [ ] 华为云 ECS 查询
- [ ] 统一资源抽象层
