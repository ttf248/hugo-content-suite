# 安装配置指南

[English](installation_en.md) | 中文

## 系统要求

### 必需环境
- **Go**: 版本 1.21 或更高
- **操作系统**: Windows, macOS, Linux
- **Hugo博客**: 支持Front Matter的Markdown文件
- **内存**: 建议 4GB 以上 (支持大型博客批量处理)
- **磁盘空间**: 至少 100MB (包含缓存和日志文件)

### 可选组件
- **LM Studio**: 用于AI翻译功能 (强烈推荐)
- **Git**: 用于版本控制
- **Visual Studio Code**: 推荐用于查看日志和配置文件

## 快速安装

### 1. 克隆项目
```bash
git clone https://github.com/ttf248/hugo-content-suite.git
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

首次运行时，程序会自动创建默认配置文件 `config.json`。

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
