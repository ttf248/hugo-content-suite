# 性能监控指南

[English](performance_en.md) | 中文

## 性能监控概述

本工具内置了完整的性能监控系统，可以实时跟踪程序运行状态、资源使用情况和操作效率，帮助用户优化处理性能。

## 监控指标

### 系统资源指标

| 指标类型 | 指标名称 | 单位 | 说明 |
|----------|----------|------|------|
| 内存 | 已用内存 | MB | 程序当前使用的内存量 |
| 内存 | 内存使用率 | % | 占系统总内存的百分比 |
| CPU | CPU使用率 | % | 程序CPU使用率 |
| 磁盘 | 磁盘读取速度 | MB/s | 文件读取速度 |
| 磁盘 | 磁盘写入速度 | MB/s | 文件写入速度 |

### 业务性能指标

| 指标类型 | 指标名称 | 单位 | 说明 |
|----------|----------|------|------|
| 处理速度 | 文章扫描速度 | 篇/秒 | 每秒扫描的文章数量 |
| 处理速度 | 翻译速度 | 标签/秒 | 每秒翻译的标签数量 |
| 网络 | API响应时间 | ms | LM Studio API平均响应时间 |
| 网络 | API成功率 | % | API调用成功率 |
| 缓存 | 缓存命中率 | % | 缓存命中的百分比 |
| 缓存 | 缓存大小 | KB | 缓存文件大小 |

## 性能监控配置

### 基本配置

```yaml
# config.yaml
performance:
  enable_monitoring: true        # 启用性能监控
  metrics_interval: 10s         # 指标收集间隔
  memory_threshold: 500MB       # 内存使用警告阈值
  cpu_threshold: 80             # CPU使用警告阈值(%)
  enable_profiling: false       # 启用Go性能分析
  profile_port: 6060           # 性能分析端口
  metrics_retention: 24h       # 指标保留时间
  export_metrics: true         # 导出Prometheus指标
```

### 高级配置

```yaml
performance:
  # 详细监控配置
  detailed_monitoring:
    goroutine_count: true      # 监控协程数量
    gc_stats: true            # 监控垃圾回收统计
    memory_breakdown: true    # 详细内存分析
    
  # 告警配置
  alerts:
    memory_critical: 1GB      # 内存严重告警阈值
    response_time_critical: 5s # 响应时间严重告警
    error_rate_critical: 10   # 错误率严重告警(%)
    
  # 性能分析配置
  profiling:
    cpu_profile_duration: 30s # CPU分析持续时间
    memory_profile_interval: 5m # 内存分析间隔
    block_profile: true       # 阻塞分析
    mutex_profile: true       # 互斥锁分析
```

## 实时性能监控

### 控制台显示

程序运行时会在界面底部显示实时性能信息：

```
=====================================
📊 实时性能监控
=====================================
💾 内存使用: 245.2MB / 500MB (49.0%)
⚡ CPU使用率: 15.3%
🏃 处理速度: 12.5 篇/秒
🌐 API响应: 1.2s (平均)
💾 缓存命中率: 87.5%
📊 协程数: 8
⏱️  运行时间: 2m 15s
=====================================
```

### 性能仪表板

启用Web仪表板查看详细性能指标：

```bash
go run main.go --dashboard --port=8080
```

访问 `http://localhost:8080/dashboard` 查看：
- 实时性能图表
- 历史趋势分析
- 资源使用详情
- 操作统计报告

## 性能分析工具

### Go pprof 集成

启用内置的Go性能分析工具：

```yaml
performance:
  enable_profiling: true
  profile_port: 6060
```

#### CPU性能分析

```bash
# 收集30秒CPU性能数据
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# 分析CPU热点
(pprof) top10
(pprof) list function_name
(pprof) web
```

#### 内存性能分析

```bash
# 内存使用分析
go tool pprof http://localhost:6060/debug/pprof/heap

# 内存分配分析
go tool pprof http://localhost:6060/debug/pprof/allocs
```

#### 协程分析

```bash
# 协程状态分析
go tool pprof http://localhost:6060/debug/pprof/goroutine

# 阻塞分析
go tool pprof http://localhost:6060/debug/pprof/block
```

### 性能报告生成

#### 自动报告

程序运行结束后自动生成性能报告：

```
=====================================
📊 性能统计报告
=====================================
⏱️  总运行时间: 5m 32s
📝 处理文章总数: 1,247 篇
🏷️  翻译标签总数: 156 个
📊 平均处理速度: 3.7 篇/秒
🌐 API调用次数: 89 次
💾 缓存命中率: 82.3%
📈 最大内存使用: 287.5MB
⚡ 平均CPU使用: 23.1%
🔄 垃圾回收次数: 12 次
=====================================
```

#### 详细报告导出

```bash
# 导出详细性能报告
go run main.go --export-performance-report=./reports/performance.json

# 导出Prometheus格式指标
go run main.go --export-metrics=./metrics/metrics.txt
```

## 性能优化建议

### 内存优化

#### 监控内存泄漏

```bash
# 定期检查内存使用趋势
watch -n 5 'curl -s http://localhost:6060/debug/pprof/heap | head -20'
```

#### 优化配置

```yaml
# 减少内存使用的配置
cache:
  max_entries: 5000      # 降低缓存条目数
  expire_after: 24h      # 缩短缓存过期时间

scanner:
  max_concurrent: 5      # 减少并发处理数

translator:
  batch_size: 5          # 减少批处理大小
```

### CPU优化

#### 并发处理优化

```yaml
scanner:
  max_concurrent: 16     # 根据CPU核数调整

performance:
  enable_profiling: true # 启用性能分析找出热点
```

#### I/O优化

```yaml
cache:
  auto_save: false       # 手动控制缓存保存时机
  save_interval: 10m     # 增加保存间隔
```

### 网络优化

#### API调用优化

```yaml
lm_studio:
  timeout: 10s           # 降低超时时间
  max_retries: 1         # 减少重试次数

translator:
  batch_size: 20         # 增加批处理大小
  request_delay: 200ms   # 减少请求间隔
```

## 性能监控告警

### 阈值告警

```yaml
performance:
  alerts:
    memory_warning: 400MB    # 内存警告阈值
    memory_critical: 800MB   # 内存严重阈值
    cpu_warning: 70         # CPU警告阈值
    cpu_critical: 90        # CPU严重阈值
    response_time_warning: 3s # 响应时间警告
    cache_hit_rate_low: 60  # 缓存命中率过低阈值
```

### 告警通知

#### 邮件通知

```yaml
alerts:
  email:
    enabled: true
    smtp_server: "smtp.gmail.com:587"
    username: "your-email@gmail.com"
    password: "your-password"
    recipients:
      - "admin@example.com"
```

#### Webhook通知

```yaml
alerts:
  webhook:
    enabled: true
    url: "https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK"
    timeout: 5s
```

## 性能基准测试

### 基准测试脚本

```bash
#!/bin/bash
# benchmark.sh

echo "开始性能基准测试..."

# 准备测试数据
mkdir -p test_data
for i in {1..1000}; do
    echo "---
title: \"测试文章 $i\"
tags: [\"标签1\", \"标签2\", \"标签3\"]
date: 2024-01-01
---
测试内容" > "test_data/article_$i.md"
done

# 运行基准测试
echo "测试开始时间: $(date)"
time go run main.go test_data --benchmark
echo "测试结束时间: $(date)"

# 清理测试数据
rm -rf test_data
```

### 基准测试结果

```
=====================================
📊 基准测试结果
=====================================
📝 测试文章数量: 1,000 篇
🏷️  测试标签数量: 50 个
⏱️  总处理时间: 2m 15s
📊 处理速度: 7.4 篇/秒
🌐 翻译速度: 3.2 标签/秒
💾 内存峰值: 156.7MB
⚡ CPU峰值: 45.2%
💾 缓存命中率: 94.1%
=====================================
```

## 故障排除

### 性能问题诊断

#### 高内存使用

```bash
# 检查内存分配热点
go tool pprof http://localhost:6060/debug/pprof/heap
(pprof) top10
(pprof) list function_name
```

#### 高CPU使用

```bash
# 检查CPU热点
go tool pprof http://localhost:6060/debug/pprof/profile
(pprof) top10
(pprof) web
```

#### 慢响应诊断

```bash
# 检查阻塞点
go tool pprof http://localhost:6060/debug/pprof/block
```

### 常见性能问题

1. **内存持续增长**
   - 检查是否有内存泄漏
   - 减少缓存大小
   - 启用定期垃圾回收

2. **处理速度慢**
   - 增加并发处理数
   - 优化文件I/O
   - 使用SSD存储

3. **API响应慢**
   - 检查网络连接
   - 调整超时设置
   - 使用本地模型

## 相关文档

- [配置文件说明](configuration.md)
- [日志系统指南](logging.md)
- [故障排除](troubleshooting.md)
- [API接口文档](api.md)
