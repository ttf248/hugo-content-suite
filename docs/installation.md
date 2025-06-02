# 安装配置指南

[English](installation_en.md) | 中文

> **Version v3.0.0** - 重构架构，企业级日志，高性能缓存

## 系统要求

### 必需环境
- **Go**: 版本 1.22.0 或更高 (推荐工具链 1.23.4)
- **操作系统**: Windows, macOS, Linux
- **Hugo博客**: 支持Front Matter的Markdown文件
- **内存**: 建议 4GB 以上 (支持大型博客批量处理)
- **磁盘空间**: 至少 200MB (包含分层缓存和轮转日志文件)

### 可选组件
- **LM Studio**: 用于AI翻译功能 (强烈推荐)
  - 推荐模型: gemma-3-12b-it, llama-3.1, qwen-2.5 等
- **Git**: 用于版本控制
- **Visual Studio Code**: 推荐用于查看结构化日志和配置文件

## 快速安装

### 1. 克隆项目
```bash
git clone https://github.com/your-org/hugo-content-suite.git
cd hugo-content-suite
```

### 2. 安装依赖
```bash
go mod tidy
```

### 3. 验证安装
```bash
go run main.go --help
```

### 4. 首次运行
```bash
go run main.go [你的content目录路径]
```

首次运行时，程序会自动创建：
- 默认配置文件 `config.json`
- 日志目录 `logs/`
- 分层缓存文件 (`*_translations_cache.json`)

## v3.0.0 新特性

### 🏗️ 重构架构
- **处理器模式**: 模块化业务逻辑，统一接口设计
- **分层缓存**: 标签/Slug/分类分离管理，提高精准度
- **统一HTTP客户端**: 消除重复代码，提升性能

### 📝 企业级日志
- **结构化日志**: JSON格式，便于分析和监控
- **自动轮转**: 日志文件自动压缩和归档
- **多级别输出**: DEBUG/INFO/WARN/ERROR级别控制
- **性能监控**: 集成操作统计和性能指标

### ⚡ 性能优化
- **批量处理**: 智能分批减少API调用次数
- **缓存预加载**: 提前检查状态，减少等待时间
- **内存优化**: 降低内存占用约30%
- **并发控制**: 可配置的并发请求限制

## 配置文件说明

### 自动生成的配置文件
程序首次运行时会在项目根目录生成 `config.json`：

```json
{
  "lm_studio": {
    "url": "http://localhost:2234/v1/chat/completions",
    "model": "gemma-3-12b-it",
    "timeout_seconds": 30,
    "max_retries": 3,
    "retry_delay_ms": 1000
  },
  "cache": {
    "auto_save_count": 10,
    "delay_ms": 500,
    "expire_days": 30,
    "enable_compression": true
  },
  "language": {
    "target_languages": ["en", "ja", "ko"],
    "language_names": {
      "en": "English",
      "ja": "Japanese", 
      "ko": "Korean"
    }
  },
  "logging": {
    "level": "INFO",
    "file": "./logs/app.log",
    "max_size_mb": 100,
    "max_backups": 10,
    "console_output": true
  },
  "performance": {
    "max_concurrent_requests": 5,
    "batch_size": 20,
    "memory_limit_mb": 512
  }
}
```

### 配置项详细说明

#### LM Studio 配置 (lm_studio)
- `url`: LM Studio API地址
- `model`: 使用的AI模型名称
- `timeout_seconds`: 请求超时时间
- `max_retries`: 最大重试次数
- `retry_delay_ms`: 重试延迟时间

#### 缓存配置 (cache)
- `auto_save_count`: 自动保存间隔
- `delay_ms`: 请求间延迟
- `expire_days`: 缓存过期天数
- `enable_compression`: 启用缓存压缩

#### 性能配置 (performance)
- `max_concurrent_requests`: 最大并发请求数
- `batch_size`: 批量处理大小
- `memory_limit_mb`: 内存限制

## LM Studio 配置

### 安装 LM Studio
1. 访问 [LM Studio官网](https://lmstudio.ai/) 下载
2. 安装并启动LM Studio
3. 下载推荐的语言模型：
   - **推荐**: Gemma-3-12B-IT (平衡性能和质量)
   - **备选**: LLaMA2-7B (较快速度)
   - **高质量**: GPT-4 (如果有API访问权限)

### 配置连接
修改 `config.json` 中的LM Studio配置：

```json
{
  "lm_studio": {
    "url": "http://192.168.1.100:2234/v1/chat/completions",  // 修改为你的LM Studio地址
    "model": "your-model-name",                               // 修改为你的模型名称
    "timeout_seconds": 45,                                    // 可根据网络情况调整
    "max_retries": 5                                          // 网络不稳定时可增加重试次数
  }
}
```

### 验证LM Studio连接
```bash
# 运行连接测试
go run main.go --test-connection

# 或启动程序后在菜单中选择测试
go run main.go [你的content目录]
```

## v3.0.0 日志系统

### 日志级别配置
程序支持多级别日志输出，可通过配置文件调整：

```json
{
  "logging": {
    "level": "INFO",        // DEBUG/INFO/WARN/ERROR
    "file": "./logs/app.log",
    "max_size_mb": 100,     // 单个日志文件最大大小
    "max_backups": 10,      // 保留的备份文件数量
    "console_output": true  // 是否同时输出到控制台
  }
}
```

### 日志文件结构
```
logs/
├── app.log              # 当前日志文件
├── app-2024-01-01.log   # 按日期轮转的备份文件
├── app-2024-01-02.log.gz # 压缩的历史日志
└── performance.log      # 性能监控日志
```

### 日志分析示例
```bash
# 查看最新错误日志
grep "ERROR" logs/app.log | tail -10

# 分析API调用性能
grep "api_call" logs/performance.log | jq '.duration'

# 监控缓存命中率
grep "cache_hit" logs/app.log | wc -l
```

## v3.0.0 缓存系统

### 分层缓存文件
v3.0.0引入了分离的缓存管理：

```
project_root/
├── tag_translations_cache.json      # 标签翻译缓存
├── slug_translations_cache.json     # Slug翻译缓存
├── category_translations_cache.json # 分类翻译缓存
└── config.json                      # 主配置文件
```

### 缓存管理
```bash
# 清理特定类型缓存
rm tag_translations_cache.json

# 清理所有缓存
rm *_translations_cache.json

# 查看缓存统计
go run main.go --cache-stats
```

### 缓存优化建议
- **过期时间**: 根据内容更新频率设置合理的过期天数
- **压缩功能**: 对于大型博客启用缓存压缩
- **预热策略**: 首次运行时建议执行完整缓存预热

## 性能优化配置

### 批量处理配置
根据系统配置调整性能参数：

```json
{
  "performance": {
    "max_concurrent_requests": 5,  // 并发请求数 (1-10)
    "batch_size": 20,             // 批量处理大小 (10-50)
    "memory_limit_mb": 512        // 内存限制 (256-1024)
  }
}
```

### 性能调优建议

#### 小型博客 (< 100篇文章)
```json
{
  "max_concurrent_requests": 3,
  "batch_size": 10,
  "memory_limit_mb": 256
}
```

#### 中型博客 (100-500篇文章)
```json
{
  "max_concurrent_requests": 5,
  "batch_size": 20,
  "memory_limit_mb": 512
}
```

#### 大型博客 (> 500篇文章)
```json
{
  "max_concurrent_requests": 8,
  "batch_size": 30,
  "memory_limit_mb": 1024
}
```

## 故障排除

### 常见问题

#### 1. LM Studio连接失败
```bash
# 检查LM Studio是否运行
curl http://localhost:2234/v1/models

# 检查网络连接
ping localhost

# 查看详细错误日志
tail -f logs/app.log
```

#### 2. 缓存问题
```bash
# 清理并重建缓存
rm *_translations_cache.json
go run main.go [content目录] --rebuild-cache
```

#### 3. 内存不足
```bash
# 减少并发数和批量大小
# 在config.json中调整:
{
  "performance": {
    "max_concurrent_requests": 2,
    "batch_size": 10
  }
}
```

#### 4. 翻译质量问题
- 检查LM Studio模型是否适合翻译任务
- 考虑更换更大的模型 (如Gemma-3-12B)
- 调整翻译提示词模板

### 日志分析
```bash
# 查看启动错误
grep "FATAL\|ERROR" logs/app.log

# 分析处理性能
grep "duration" logs/performance.log | tail -20

# 监控缓存使用
grep "cache" logs/app.log | grep "hit\|miss"
```

## 高级配置

### 自定义翻译模板
创建 `templates/translation_prompt.txt` 自定义翻译提示词：

```text
请将以下{source_language}文本翻译成{target_language}:

原文: {content}

要求:
1. 保持Markdown格式不变
2. 保持专业术语准确性
3. 符合{target_language}表达习惯
4. 不要翻译代码块内容

翻译:
```

### 自定义标签页模板
创建 `templates/tag_page.md` 自定义标签页模板：

```markdown
---
title: "标签: {{.Name}}"
slug: "{{.Slug}}"
description: "包含{{.Count}}篇文章的{{.Name}}标签页面"
---

# {{.Name}}

{{.Description}}

## 相关文章 ({{.Count}}篇)
```

### 环境变量配置
支持通过环境变量覆盖配置：

```bash
export HCS_LOG_LEVEL=DEBUG
export HCS_LM_STUDIO_URL=http://192.168.1.100:2234/v1/chat/completions
export HCS_CACHE_EXPIRE_DAYS=7

go run main.go [content目录]
```

## 部署建议

### 服务器部署
```bash
# 编译二进制文件
go build -o hugo-content-suite main.go

# 创建系统服务 (Linux)
sudo cp hugo-content-suite /usr/local/bin/
sudo chmod +x /usr/local/bin/hugo-content-suite

# 配置定时任务
crontab -e
# 每天凌晨2点自动处理
0 2 * * * /usr/local/bin/hugo-content-suite /path/to/content --auto-process
```

### Docker部署
创建 `Dockerfile`：

```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o hugo-content-suite main.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /root/
COPY --from=builder /app/hugo-content-suite .
COPY --from=builder /app/config.json .
EXPOSE 8080
CMD ["./hugo-content-suite"]
```

---

## 版本升级

### 从v2.x升级到v3.0.0

1. **备份现有配置和缓存**
```bash
cp config.json config.json.backup
cp translations_cache.json translations_cache.json.backup
```

2. **更新代码**
```bash
git pull origin main
go mod tidy
```

3. **迁移配置**
v3.0.0会自动检测旧配置格式并提示升级。

4. **重建缓存**
```bash
# 删除旧的单一缓存文件
rm translations_cache.json

# 运行程序，自动创建新的分层缓存
go run main.go [content目录]
```

### 配置迁移指南

#### v2.x配置格式
```json
{
  "lm_studio_url": "http://localhost:2234/v1/chat/completions",
  "cache_file": "translations_cache.json"
}
```

#### v3.0.0配置格式
```json
{
  "lm_studio": {
    "url": "http://localhost:2234/v1/chat/completions",
    "model": "gemma-3-12b-it",
    "timeout_seconds": 30
  },
  "cache": {
    "auto_save_count": 10,
    "delay_ms": 500
  }
}
```

---

## 技术支持

### 获取帮助
- **命令行帮助**: `go run main.go --help`
- **配置示例**: 查看自动生成的 `config.json`
- **日志分析**: 检查 `logs/app.log` 文件
- **GitHub Issues**: 报告问题和功能请求

### 贡献代码
欢迎提交Pull Request和Issue，帮助改进Hugo Content Suite。

### 许可证
本项目采用MIT许可证，详见LICENSE文件。
  }
}
```

### 测试连接
```bash
go run main.go
# 程序启动时会自动测试LM Studio连接
# 或选择菜单项进行翻译测试
```

## 目录结构配置

### 推荐的项目结构
```
your-hugo-blog/
├── content/
│   ├── post/              # 文章目录
│   │   ├── article1.md
│   │   └── article2.md
│   └── tags/              # 标签页面目录 (工具自动创建)
│       ├── ai/
│       └── tech/
├── hugo-content-suite/    # 本工具目录
│   ├── config.json        # 配置文件
│   ├── cache/             # 缓存目录 (自动创建)
│   │   ├── tag_cache.json
│   │   └── article_cache.json
│   ├── logs/              # 日志目录 (自动创建)
│   │   └── app.log
│   └── ...
└── ...
```

### 使用不同内容目录

#### 默认目录
```bash
go run main.go  # 程序会提示输入content目录路径
```

#### 直接指定目录
```bash
go run main.go /path/to/your/content/post
```

#### Windows 路径示例
```bash
go run main.go "C:\Users\Username\myblog\content\post"
```

#### 相对路径示例
```bash
go run main.go ../content/post
```

## 高级配置

### 性能优化配置
针对不同使用场景的配置建议：

#### 大型博客 (1000+ 文章)
```json
{
  "performance": {
    "max_concurrent_requests": 3,
    "batch_size": 50,
    "memory_limit_mb": 1024
  },
  "cache": {
    "auto_save_count": 20,
    "enable_compression": true
  }
}
```

#### 快速处理模式
```json
{
  "performance": {
    "max_concurrent_requests": 10,
    "batch_size": 100,
    "memory_limit_mb": 2048
  },
  "lm_studio": {
    "timeout_seconds": 15,
    "max_retries": 1
  }
}
```

#### 稳定性优先模式
```json
{
  "performance": {
    "max_concurrent_requests": 1,
    "batch_size": 10,
    "memory_limit_mb": 256
  },
  "lm_studio": {
    "timeout_seconds": 60,
    "max_retries": 10,
    "retry_delay_ms": 2000
  }
}
```

### 日志配置
```json
{
  "logging": {
    "level": "DEBUG",        // 开发环境使用DEBUG，生产环境使用INFO
    "file": "./logs/app.log",
    "max_size_mb": 200,      // 大型博客可增加日志文件大小
    "max_backups": 30,       // 保留更多备份文件
    "console_output": false  // 生产环境可关闭控制台输出
  }
}
```

## 验证安装

### 检查文件结构
确保你的Hugo博客具有正确的文件结构：

```
your-blog/
├── content/
│   ├── post/           # 文章目录
│   │   ├── article1.md
│   │   └── article2.md
│   └── tags/           # 标签页面目录（工具会自动创建）
└── ...
```

### 检查文章格式
确保Markdown文件包含完整的Front Matter：

```yaml
---
title: "文章标题"
date: 2024-01-01T12:00:00+08:00
tags: ["AI", "技术", "编程"]
categories: ["开发"]
slug: ""                    # 可选，工具会自动生成
author: "作者名"
description: "文章描述"
---

文章内容...
```

### 验证功能
运行以下命令验证各项功能：

```bash
# 1. 验证基本功能
go run main.go /path/to/content

# 2. 检查配置文件
cat config.json

# 3. 查看生成的目录结构
ls -la cache/
ls -la logs/

# 4. 测试翻译功能（如果配置了LM Studio）
# 在程序菜单中选择 "生成全量翻译缓存"
```

## 故障排除

### 常见问题

#### 1. Go版本问题
```bash
go version  # 检查当前版本
# 如果版本过低，升级到1.21+
```

#### 2. 依赖问题
```bash
go clean -modcache
go mod download
go mod tidy
```

#### 3. 权限问题
确保有必要的权限：
```bash
# Linux/macOS
chmod 755 hugo-content-suite/
chmod 666 config.json

# Windows（以管理员身份运行）
icacls hugo-content-suite /grant Everyone:F
```

#### 4. LM Studio连接问题
- 检查LM Studio是否在运行
- 验证端口是否正确 (默认2234)
- 测试网络连接：
```bash
curl -X POST http://localhost:2234/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{"model":"test","messages":[{"role":"user","content":"test"}]}'
```

#### 5. 缓存相关问题
```bash
# 清理缓存
rm -rf cache/
mkdir cache

# 检查磁盘空间
df -h
```

### 配置文件损坏
如果配置文件损坏：
```bash
# 删除配置文件，程序会重新创建默认配置
rm config.json
go run main.go
```

### 日志查看
查看详细日志信息：
```bash
# 查看实时日志
tail -f logs/app.log

# 查看错误日志
grep "ERROR" logs/app.log

# 查看性能信息
grep "PERF" logs/app.log
```

## 下一步

### 推荐工作流
1. **安装完成后**：查看 [使用说明](usage.md)
2. **配置优化**：参考 [配置指南](configuration.md)
3. **性能调优**：查看 [性能优化指南](performance.md)
4. **问题排查**：参考 [故障排除指南](troubleshooting.md)

### 进阶使用
- [架构设计文档](architecture.md) - 了解系统架构
- [缓存策略说明](caching.md) - 优化缓存使用
- [日志系统指南](logging.md) - 监控和调试

---

安装完成后，建议先使用"一键处理全部"功能体验完整工作流！
