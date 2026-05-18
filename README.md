# Cloud Manage

阿里云资源管理 CLI 工具，支持 ECS、CMS、SLS、OSS 服务。

## 安装

```bash
go build -o cloud-manage ./cmd/cli/
```

## 配置

通过环境变量或命令行参数设置凭证：

```bash
export CLOUD_ACCESS_KEY_ID=your-key-id
export CLOUD_ACCESS_KEY_SECRET=your-key-secret
```

或使用参数：`-id <key-id> -secret <key-secret>`

## 使用

```bash
# ECS 实例管理
cloud-manage ecs list
cloud-manage ecs detail i-xxx
cloud-manage ecs start i-xxx
cloud-manage ecs stop i-xxx [--force]

# 云监控
cloud-manage cms metrics i-xxx

# 日志服务
cloud-manage sls logstores <project>
cloud-manage sls logs <project> <logstore> --query "level: ERROR"

# 对象存储
cloud-manage oss buckets
cloud-manage oss objects <bucket> --prefix logs/ --max 100
```

## 通用参数

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `-id` | AccessKey ID | 环境变量 `CLOUD_ACCESS_KEY_ID` |
| `-secret` | AccessKey Secret | 环境变量 `CLOUD_ACCESS_KEY_SECRET` |
| `-region` | 地域 | `cn-hangzhou` |
| `-json` | JSON 输出 | `false` |

## License

MIT
