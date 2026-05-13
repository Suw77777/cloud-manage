# Changelog

## v0.0.8 (2026-05-12)

### CLI 支持

- 新增完整的命令行工具 `cmd/cli/main.go`
- 支持 ECS、CMS、SLS、OSS 四大服务
- 支持环境变量 `CLOUD_ACCESS_KEY_ID` / `CLOUD_ACCESS_KEY_SECRET` 配置凭证
- 支持 `-json` 参数输出 JSON 格式
- 支持 `-region all` 查询所有 Region

### CLI 命令

```bash
# ECS 实例管理
cloud-cli ecs list                        # 列出实例
cloud-cli ecs detail <id>                 # 查看详情
cloud-cli ecs start/stop/reboot <id>      # 操作实例

# 云监控
cloud-cli cms metrics <id>               # 查询监控指标

# 日志服务
cloud-cli sls logstores <project>        # 列出 Logstore
cloud-cli sls logs <project> <logstore>  # 查询日志

# 对象存储
cloud-cli oss buckets                    # 列出 Bucket
cloud-cli oss objects <bucket>           # 列出对象
```

## v0.0.7 (2026-05-11)

### 内存优化

- 新增前端分页组件 `PaginationBar.vue`
- 新增虚拟滚动组件 `VirtualScroller.vue`
- 新增分页 composable `usePagination.js`
- SLS 日志查询支持分页显示（默认每页50条）
- SLS 日志支持流式查询（边查边显示，减少等待时间）
- 后端新增 `QuerySLSLogsStream` 分块查询方法

### 优化效果

| 场景 | 优化前 | 优化后 |
|------|--------|--------|
| 1000条日志 | 全部渲染DOM | 仅渲染当前页50条 |
| 大量数据查询 | 等待全部加载 | 流式显示，边查边看 |
| 内存占用 | 持续增长 | 分页控制，稳定可控 |

## v0.0.6 (2026-05-11)

### 功能

- 新增 OSS 对象存储浏览器模块
- Tab 导航新增 "OSS 存储" 标签页
- 支持列出所有 Bucket
- 支持浏览 Bucket 内的对象和目录
- 支持目录导航：面包屑、返回上级、返回根目录
- 对象表格展示：图标、名称、大小、修改时间、存储类型
- 文件夹识别和导航

### 后端

- 新增 `provider/aliyun/oss.go` - OSS Provider 封装阿里云 OSS SDK
- 新增 `service/oss.go` - OSS Service 层
- 新增 `ListOSSBuckets` 和 `ListOSSObjects` 后端方法
- 添加 `github.com/alibabacloud-go/oss-20190517` 依赖

### 前端

- 新增 `composables/useOSS.js` - OSS 浏览器状态管理
- 新增 `components/OssBrowser.vue` - OSS 浏览器组件
- App.vue 添加 OSS Tab

## v0.0.5 (2026-05-11)

### 功能

- 新增 SLS 日志查询模块
- Tab 导航新增 "SLS 日志" 标签页
- 支持输入 SLS Project 名称并获取 Logstore 列表
- 支持日志查询：时间范围选择（15分钟到7天）、最大行数设置
- 支持自定义查询语句（如 `level: ERROR OR status: 500`）
- 日志结果表格展示：时间戳 + 日志内容字段
- 结果截断提示

### 后端

- 新增 `provider/aliyun/sls.go` - SLS Provider 封装阿里云 SLS SDK v2
- 新增 `service/sls.go` - SLS Service 层
- 新增 `ListSLSLogStores` 和 `QuerySLSLogs` 后端方法
- 添加 `github.com/alibabacloud-go/sls-20201230/v2` 依赖

### 前端

- 新增 `composables/useSLS.js` - SLS 日志查询状态管理
- 新增 `components/SlsQuery.vue` - SLS 日志查询组件
- App.vue 添加 SLS Tab

## v0.0.4 (2026-05-11)

### 功能

- 新增云监控 (CloudMonitor) 模块
- Tab 导航：ECS 实例管理 / 云监控
- 云监控页面支持从 ECS 查询结果中选择实例
- 支持批量选择实例进行监控查询
- 监控指标展示：CPU 使用率、内存使用率、磁盘读写 BPS、网络流量
- 进度条可视化展示 CPU/内存使用率（颜色阈值：正常/警告/危险）
- 支持多 Region 实例并发查询监控数据

### 后端

- 新增 `provider/aliyun/cms.go` - CMS Provider 封装阿里云 CloudMonitor SDK
- 新增 `service/cms.go` - CMS Service 层
- 新增 `GetECSMetrics` 和 `GetECSMetricsMultiRegion` 后端方法
- 添加 `github.com/alibabacloud-go/cms-20190101/v8` 依赖

### 前端

- 新增 `composables/useCMS.js` - 云监控状态管理
- 新增 `components/CmsMonitor.vue` - 云监控组件
- App.vue 添加 Tab 切换功能
- main.css 添加 CSS 变量和 Tab 导航样式

## v0.0.3-refactor (2026-05-09)

### 重构

- App.vue 从 1140 行降到 170 行
- 样式抽到 `assets/main.css`
- 逻辑抽到 `composables/useECS.js`
- 组件拆分：`EcsResultTable`、`InstanceDetailModal`、`ConfirmDialog`、`OperationLog`
- 新增工程约束：入口文件不超过 300 行

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
