# Roadmap

## v0.2.0 - 凭证管理 + SLS 增强 + 体验优化

**目标**: 多账号凭证管理、SLS 日志导出、性能和体验提升

### 设计决策记录

#### 决策 1: 凭证加密方案
- **选择**: 主密码 + AES-256-GCM 加密
- **理由**: 简单，用户容易理解，不依赖外部工具
- **工作流程**: 首次使用设置主密码，后续每次启动输入主密码解锁

#### 决策 2: 主密码缓存策略
- **选择**: 仅当前会话（程序退出就忘记）
- **理由**: 安全性优先，避免密码泄露风险

#### 决策 3: 配置文件格式
- **选择**: YAML
- **理由**: 可读性好，Go 支持成熟

#### 决策 4: 配置文件位置
- **选择**: 跨平台标准位置（`os.UserConfigDir()`）
  - Linux: `~/.config/cloud-manage/config.yaml`
  - macOS: `~/Library/Application Support/cloud-manage/config.yaml`
  - Windows: `%APPDATA%\cloud-manage\config.yaml`
- **理由**: 遵循各平台标准

#### 决策 5: Profile 配置结构
- **选择**: 每个 profile 独立配置 region/endpoint
- **理由**: 简单明确，没有隐式行为

#### 决策 6: 参数优先级
- **选择**: 命令行参数 > 环境变量 > 配置文件 > 默认值
- **理由**: 符合 Unix 惯例，适合 CI/CD 场景

#### 决策 7: 环境变量与配置文件不一致处理
- **选择**: 打印警告，继续执行
- **范围**: 所有参数（AK/SK、region、endpoint）不一致都提醒
- **理由**: 不中断流程，但让用户知道当前使用的值

#### 决策 8: SLS 导出范围
- **选择**: 可以选择范围或已选的日志
- **限制**: 大于 5000 条日志直接拒绝
- **理由**: 避免大数据量问题

#### 决策 9: SLS 导出文件命名
- **选择**: 默认自动生成（包含项目、logstore、时间戳），`--output` 可覆盖
- **理由**: 灵活，用户不指定时有有意义的文件名

---

### 1. 凭证管理（多账号切换）

**存储位置**: 跨平台标准位置（见决策 4）

**配置文件格式**:
```yaml
current_profile: prod
profiles:
  prod:
    access_key_id: "LTAI4xxx"
    access_key_secret: "encrypted:xxxx"  # AES-256-GCM 加密
    region: "cn-hangzhou"
  dev:
    access_key_id: "LTAI4xxx"
    access_key_secret: "encrypted:xxxx"
    region: "cn-shanghai"
    endpoint: "ecs.cn-shanghai.aliyuncs.com"  # 可选自定义 endpoint
```

**加密方案** (决策 1):
- 算法: AES-256-GCM
- 密钥派生: 用户主密码 → PBKDF2 (100000 iterations) → 256-bit key
- 首次使用时设置主密码，后续每次启动输入主密码解锁
- 缓存策略: 仅当前会话（决策 2）

**CLI 命令**:
```bash
# 配置管理
cloud-manage config init              # 初始化配置文件
cloud-manage config add <profile>     # 添加账号
cloud-manage config remove <profile>  # 删除账号
cloud-manage config list              # 列出所有账号
cloud-manage config switch <profile>  # 切换默认账号
cloud-manage config show              # 显示当前账号

# 使用指定账号
cloud-manage --profile dev ecs list
```

**参数优先级** (决策 6):
1. 命令行参数 (`--region cn-shanghai`)
2. 环境变量 (`CLOUD_REGION=cn-shanghai`)
3. 配置文件 (`region: "cn-hangzhou"`)
4. 默认值 (`cn-hangzhou`)

**不一致提醒** (决策 7):
- 当环境变量与配置文件不一致时，打印警告，继续执行
- 所有参数（AK/SK、region、endpoint）都提醒

**GUI/TUI**:
- 登录界面增加"选择账号"下拉框
- 支持在界面中切换账号
- 首次使用引导设置主密码

**优先级**: 高 - 这是基础设施，其他功能依赖它

---

### 2. SLS 日志增强

**导出功能** (决策 8、9):
```bash
# CLI 导出
cloud-manage sls export <project> <logstore> --format csv --output logs.csv
cloud-manage sls export <project> <logstore> --format json --output logs.json

# 支持查询过滤
cloud-manage sls export <project> <logstore> --query "level: ERROR" --from "2024-01-01" --to "2024-01-02"

# 不指定文件名时自动生成
cloud-manage sls export <project> <logstore> --format csv
# 生成: sls_project_logstore_20240101_120000.csv
```

**导出限制**:
- 导出范围: 可以选择范围或已选的日志
- 数量限制: 大于 5000 条日志直接拒绝
- 理由: 避免大数据量问题

**GUI 导出**:
- 查询结果区域增加"导出"按钮
- 支持选择导出格式 (CSV/JSON)
- 支持选择导出范围（全部/已选/查询范围）
- 超过 5000 条时提示用户缩小范围

**大结果集优化**:
- 虚拟滚动：只渲染可见区域的 DOM
- 分页查询：支持分页显示

**优先级**: 中 - 高频需求，但不阻塞其他功能

---

### 3. 性能优化

**优化清单**:

| 优化点 | 说明 | 状态 |
|--------|------|------|
| A. 模式相关初始化 | 只初始化当前模式需要的组件 | 待做 |
| B. SDK 客户端复用 | 多次查询复用同一个客户端 | 待做 |
| C. 前端产物压缩 | Vite 构建时启用 gzip/brotli 压缩 | 待做 |
| D. Go 编译优化 | `-ldflags "-s -w"` 减小二进制体积，`-trimpath` 移除路径信息 | 待做 |
| E. GC 调优 | 设置 `GOMEMLIMIT` 控制内存上限 | 待做 |
| F. 懒加载 embed.FS | GUI 模式才解压前端资源 | 待做 |
| G. 并发查询优化 | 多 Region 查询时控制并发数，避免打满连接 | 待做 |

**优先级**: 中 - 预防性优化，提升体验

---

### 4. 体验优化

**GUI 界面**:
- 统一色彩方案 (主色、辅色、强调色)
- 统一按钮、输入框、表格样式
- 增加加载动画、过渡效果
- 响应式布局优化

**TUI 体验**:
- 完善键盘快捷键 (j/k 导航, / 搜索, q 退出)
- 增加帮助面板 (? 键显示)
- 状态栏显示当前操作
- 颜色主题支持

**错误处理**:
- 友好的错误消息（非技术语言）
- 错误码 + 详细说明链接
- 操作建议（如"请检查 AK/SK 是否正确"）

**优先级**: 低 - 锦上添花，最后做

---

### 实施计划

不分阶段，全部一起完成。

**凭证管理**:
- 配置文件读写
- 加密存储 (AES-256-GCM + 主密码)
- CLI config 子命令
- GUI/TUI 账号选择器
- 参数优先级处理
- 不一致提醒

**SLS 增强**:
- CLI export 子命令
- GUI 导出按钮（直接下载，默认 CSV）
- 大结果集优化（虚拟滚动）
- 时间格式混合支持（ISO 8601 默认，支持相对时间和 Unix 时间戳）
- 5000 条限制

**性能优化**:
- 模式相关初始化
- SDK 客户端复用（按需创建 + 缓存）
- 前端产物压缩
- Go 编译优化 (`-ldflags "-s -w"` + `-trimpath`)
- GC 调优 (256MB 默认，环境变量/配置文件可修改)
- 懒加载 embed.FS
- 并发查询优化（默认 3，可配置）

**体验优化**:
- GUI 微调统一风格（v0.2.0 微调，未来版本 Material Design）
- TUI 快捷键（vim + 方向键都支持）
- 错误处理友好化（CLI 技术细节，GUI/TUI 友好提示）
- 暗色主题（跟随系统，可手动切换，默认 light）

---

## v0.1.0 - VPC/SLB 管理

**目标**: 添加 VPC 和 SLB 管理功能，完善工程质量。

- [x] VPC 网络管理（VPC 列表、详情、VSwitch）
- [x] SLB 负载均衡管理（SLB 列表、详情、监听器）
- [x] 统一入口（删除独立 CLI/TUI）
- [x] 统一版本号管理
- [x] 完善测试覆盖（Provider、Service、TUI）
- [x] 前端 VPC/SLB 组件
- [x] 清理构建产物

---

## v0.0.9 - 统一入口 + TUI

**目标**: 合并 GUI、CLI、TUI 为单一二进制，自动检测环境选择模式。

### 统一入口
- [x] 合并 GUI 和 CLI 到单一入口 `main.go`
- [x] 自动检测图形环境（DISPLAY/WAYLAND_DISPLAY）
- [x] 支持 `--gui`、`--cli`、`--tui` 标志强制指定模式
- [x] 检测到子命令自动进入 CLI 模式
- [x] 启动时显示模式提示信息
- [x] 删除独立的 `cmd/cli/` 和 `cmd/tui/` 目录

### TUI 终端界面
- [x] 使用 Bubble Tea 实现 TUI 框架
- [x] 登录界面输入凭证
- [x] Tab 切换服务（ECS、CMS、SLS、OSS、VPC、SLB）
- [x] 方向键导航，Enter 查看详情

---

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
