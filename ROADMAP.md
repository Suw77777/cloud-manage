# Roadmap

## v0.0.8 - CLI 命令行工具

**目标**: 新增完整 CLI 命令行工具。

- [x] CLI 支持 ECS、CMS、SLS、OSS 四大服务
- [x] 支持环境变量配置凭证
- [x] 支持 `-json` 参数输出 JSON 格式
- [x] 支持 `-region all` 查询所有 Region

---

## v0.0.7 - 内存优化

**目标**: 优化大数据量场景下的内存占用。

- [x] 新增分页组件 PaginationBar.vue
- [x] 新增虚拟滚动组件 VirtualScroller.vue
- [x] SLS 日志查询支持分页显示（默认每页 50 条）
- [x] SLS 日志支持流式查询（边查边显示）

---

## v0.0.6 - OSS 对象存储

**目标**: 支持阿里云 OSS 对象存储浏览。

- [x] OSS Bucket 列表
- [x] OSS Object 浏览
- [x] 目录导航（面包屑）
- [x] 文件大小和时间显示

---

## v0.0.5 - SLS 日志查询

**目标**: 支持阿里云 SLS 日志查询。

- [x] SLS Logstore 列表
- [x] 日志查询输入框
- [x] 日志结果展示
- [x] 时间范围选择
- [x] 查询语句支持

---

## v0.0.4 - 云监控

**目标**: 支持阿里云 CloudMonitor 监控数据查询。

- [x] Tab 导航：ECS 实例管理 / 云监控
- [x] 云监控页面：从 ECS 结果选择实例
- [x] 监控指标查询：CPU、内存、磁盘、网络
- [x] 监控数据可视化展示
- [x] 多实例批量查询

---

## v0.0.3 - ECS 详情与操作

**目标**: 支持查看 ECS 实例详情和基本操作。

- [x] ECS 实例详情页（弹窗展示 CPU/Memory/Image/VPC/SG 等）
- [x] 启动 / 停止 / 重启 ECS 实例（二次确认）
- [x] 生产环境操作强制二次确认（红色警告 + 特殊按钮文案）
- [x] 操作结果日志记录（前端内存，最多 50 条）

---

## v0.0.2 - 多 Region 支持

**目标**: 支持同时查询多个 Region 的 ECS 实例。

- [x] 支持多 Region 下拉选择（18 个常用 Region 复选框）
- [x] 支持多 Region 并发查询（goroutine per region）
- [x] 结果按 Region 分组展示

---

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

## 未来规划（暂不实施）

以下功能在项目边界之外，暂不考虑：

- 多云支持（腾讯云、华为云）- 复杂度高，ROI 低
- 完整的资源管理（创建/删除 VPC、RDS 等）- 风险高
- 用户认证系统 - 超出工具定位
- 数据持久化（数据库）- 增加复杂度
- 多用户协作 - 超出桌面工具定位
