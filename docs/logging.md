# 日志系统指南

[English](logging_en.md) | 中文

## 日志系统概述

本工具采用分级日志系统，支持多种输出格式和自动轮转功能，帮助用户跟踪程序运行状态和排查问题。

## 日志级别

### 级别说明

| 级别 | 数值 | 用途 | 示例 |
|------|------|------|------|
| DEBUG | 0 | 详细调试信息 | 函数调用、变量值、执行流程 |
| INFO | 1 | 一般信息记录 | 操作开始/完成、配置加载、统计信息 |
| WARN | 2 | 警告信息 | 配置不当、性能问题、兼容性问题 |
| ERROR | 3 | 错误信息 | 操作失败、网络错误、文件读写错误 |

### 级别设置

```yaml
# config.yaml
logging:
  level: "INFO"  # 只输出INFO及以上级别的日志
```

环境变量覆盖：
```bash
export HUGO_SLUG_LOGGING_LEVEL="DEBUG"
```

## 日志输出配置

### 文件输出

```yaml
logging:
  file_path: "./logs/app.log"
  max_size: 100MB      # 单文件最大大小
  max_backups: 5       # 保留的历史文件数
  max_age: 30          # 文件保留天数
```

### 控制台输出

```yaml
logging:
  console_output: true  # 同时输出到控制台
  json_format: false    # 控制台使用纯文本格式
```

### 双重输出

程序支持同时输出到文件和控制台，便于开发调试和生产监控。

## 日志格式

### 文本格式（默认）

```
2024-01-15 14:30:25 [INFO] [scanner] 开始扫描文章目录: ./content/post
2024-01-15 14:30:25 [DEBUG] [scanner] 找到文件: article1.md
2024-01-15 14:30:25 [INFO] [translator] 开始批量翻译 5 个标签
2024-01-15 14:30:26 [WARN] [cache] 缓存文件不存在，将创建新文件
2024-01-15 14:30:27 [ERROR] [translator] LM Studio连接失败: connection refused
```

### JSON格式

```yaml
logging:
  json_format: true
```

JSON输出示例：
```json
{
  "timestamp": "2024-01-15T14:30:25.123Z",
  "level": "INFO",
  "module": "scanner",
  "message": "开始扫描文章目录",
  "data": {
    "directory": "./content/post",
    "file_count": 0
  },
  "correlation_id": "req_123456"
}
```

## 日志模块

### 核心模块标识

| 模块 | 标识 | 功能 |
|------|------|------|
| 扫描器 | scanner | 文章扫描和解析 |
| 翻译器 | translator | AI翻译功能 |
| 生成器 | generator | 内容生成 |
| 缓存 | cache | 缓存管理 |
| 性能 | performance | 性能监控 |
| 配置 | config | 配置管理 |
| 主程序 | main | 主程序流程 |

## 日志文件管理

### 自动轮转

当日志文件达到指定大小时自动轮转：

```
logs/
├── app.log           # 当前日志文件
├── app.log.1         # 最近的备份
├── app.log.2         # 次新的备份
├── app.log.3
├── app.log.4
└── app.log.5         # 最旧的备份
```

### 自动清理

- 超过 `max_backups` 数量的文件会被删除
- 超过 `max_age` 天数的文件会被删除
- 清理操作在每次轮转时执行

### 压缩存储

历史日志文件会自动压缩（需要配置）：

```yaml
logging:
  compress_backups: true  # 压缩历史日志
```

压缩后的文件格式：
```
logs/
├── app.log
├── app.log.1.gz
├── app.log.2.gz
└── ...
```

## 日志查看和分析

### 实时查看

```bash
# 查看实时日志
tail -f logs/app.log

# 查看最近100行
tail -n 100 logs/app.log

# 实时查看并高亮错误
tail -f logs/app.log | grep --color=always -E "(ERROR|WARN|$)"
```

### 按级别过滤

```bash
# 查看错误日志
grep "ERROR" logs/app.log

# 查看警告和错误
grep -E "(WARN|ERROR)" logs/app.log

# 查看特定模块日志
grep "[translator]" logs/app.log
```

### 按时间过滤

```bash
# 查看今天的日志
grep "$(date +%Y-%m-%d)" logs/app.log

# 查看特定时间段
grep "2024-01-15 14:" logs/app.log

# 查看最近1小时的错误
grep "ERROR" logs/app.log | grep "$(date -d '1 hour ago' +%Y-%m-%d\ %H):"
```

### 统计分析

```bash
# 统计各级别日志数量
grep -c "ERROR" logs/app.log
grep -c "WARN" logs/app.log
grep -c "INFO" logs/app.log

# 统计每小时的错误数
grep "ERROR" logs/app.log | cut -d' ' -f2 | cut -d':' -f1 | sort | uniq -c

# 找出最频繁的错误
grep "ERROR" logs/app.log | cut -d']' -f3- | sort | uniq -c | sort -nr | head -10
```

## 日志监控和告警

### 日志监控脚本

```bash
#!/bin/bash
# monitor_logs.sh

LOG_FILE="./logs/app.log"
ERROR_THRESHOLD=10

# 统计最近5分钟的错误数
ERROR_COUNT=$(tail -n 1000 "$LOG_FILE" | grep "$(date -d '5 minutes ago' +%Y-%m-%d\ %H:%M)" | grep -c "ERROR")

if [ "$ERROR_COUNT" -gt "$ERROR_THRESHOLD" ]; then
    echo "警告：最近5分钟内发生 $ERROR_COUNT 个错误，超过阈值 $ERROR_THRESHOLD"
    # 发送告警通知
fi
```

### 集成监控系统

#### Prometheus + Grafana

```yaml
# prometheus.yml
- job_name: 'hugo-slug-auto'
  static_configs:
    - targets: ['localhost:8080']
  metrics_path: '/metrics'
```

#### ELK Stack

```yaml
# filebeat.yml
filebeat.inputs:
- type: log
  enabled: true
  paths:
    - ./logs/*.log
  fields:
    service: hugo-slug-auto
```

## 调试技巧

### 开启详细日志

```yaml
logging:
  level: "DEBUG"
  console_output: true
```

### 跟踪特定操作

在代码中添加关联ID：
```go
correlationID := uuid.New().String()
logger.WithField("correlation_id", correlationID).Info("开始处理请求")
```

### 性能调试

开启性能日志：
```yaml
logging:
  performance_logging: true
```

性能日志示例：
```
2024-01-15 14:30:25 [PERF] [translator] 翻译耗时: 1.23s, 标签数: 10
2024-01-15 14:30:26 [PERF] [cache] 缓存命中率: 85.2%
```

## 最佳实践

### 开发环境

```yaml
logging:
  level: "DEBUG"
  console_output: true
  file_path: "./logs/dev.log"
  json_format: false
  max_size: 10MB
  max_backups: 3
```

### 生产环境

```yaml
logging:
  level: "INFO"
  console_output: false
  file_path: "/var/log/hugo-slug-auto/app.log"
  json_format: true
  max_size: 100MB
  max_backups: 10
  max_age: 30
  compress_backups: true
```

### 性能考虑

1. **避免过度日志记录**
   - 生产环境使用INFO级别
   - 避免在循环中记录DEBUG日志

2. **异步日志写入**
   ```yaml
   logging:
     async_write: true
     buffer_size: 1000
   ```

3. **日志轮转策略**
   - 根据磁盘空间设置合理的文件大小
   - 平衡历史保留和存储成本

## 故障排除

### 常见问题

1. **日志文件无法创建**
   ```
   原因：目录权限不足
   解决：检查并修改目录权限
   ```

2. **日志轮转失败**
   ```
   原因：磁盘空间不足
   解决：清理磁盘空间或调整日志策略
   ```

3. **JSON格式错误**
   ```
   原因：日志内容包含特殊字符
   解决：启用字符转义或使用文本格式
   ```

4. **性能影响**
   ```
   原因：DEBUG级别日志过多
   解决：调整日志级别或启用异步写入
   ```

## 相关文档

- [配置文件说明](configuration.md)
- [性能监控指南](performance.md)
- [故障排除](troubleshooting.md)
- [API接口文档](api.md)
