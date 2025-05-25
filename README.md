# Hugo 博客管理工具

[English](README_EN.md) | 中文

> 🚀 一款专为Hugo博客设计的智能管理工具，支持AI多语言翻译、高性能缓存和一体化工作流

## ✨ 核心特色

### 🤖 AI智能翻译
- 基于LM Studio的本地AI模型翻译
- 支持多语言翻译（英语、日语、韩语等）
- 中文标签自动转换为SEO友好的英文slug
- 统一的HTTP请求处理，提高稳定性和性能
- 智能缓存机制，分类管理标签和文章翻译

### 💾 高性能缓存系统
- 本地JSON分层缓存，避免重复翻译
- 分离式缓存管理（标签缓存、文章缓存）
- 自动检测缓存状态，智能批量处理
- 支持缓存统计和清理功能
- 缓存命中率监控和性能优化

### 🎯 优化的用户体验
- 精简的命令行菜单系统
- 彩色输出和进度条显示
- 一键处理全流程自动化
- 智能错误处理和重试机制
- 详细的操作统计和性能报告

### 📝 企业级日志系统
- 结构化日志记录，支持多级别输出
- 自动日志文件轮转和压缩归档
- 详细的源码位置信息和调用栈
- 支持控制台和文件双重输出
- 集成logrus高性能日志库
- 操作审计和性能监控

### 🚀 一键处理功能
- 全流程自动化处理工作流
- 智能状态分析和预处理检查
- 批量缓存生成和优化
- 多语言文章翻译支持
- 完整的博客管理解决方案

## 🏗️ 重构优化亮点

### 代码架构优化
- **统一HTTP客户端**: 消除重复代码，提高请求处理效率
- **通用翻译方法**: 使用模板化提示词，支持不同翻译类型
- **分层缓存设计**: 标签和文章分离管理，提高缓存命中率
- **函数式设计**: 可组合的翻译处理函数，易于扩展和维护

### 性能提升
- **批量处理优化**: 智能分批处理，减少网络开销
- **缓存预加载**: 提前检查缓存状态，减少重复查询
- **进度追踪**: 实时显示处理进度，提升用户体验
- **内存优化**: 减少重复对象创建，降低内存占用

## 📁 项目架构

```
hugo-content-suite/
├── main.go              # 主程序入口
├── config/              # 配置管理
│   ├── config.go        # 配置结构和加载
│   └── validation.go    # 配置验证
├── models/              # 数据模型
│   ├── article.go       # 文章模型
│   └── metadata.go      # 元数据结构
├── scanner/             # 文章扫描解析
│   ├── scanner.go       # 文件扫描器
│   └── parser.go        # Markdown解析器
├── stats/               # 统计分析
│   ├── collector.go     # 数据收集器
│   └── reporter.go      # 统计报告
├── translator/          # AI翻译模块 (重构优化)
│   ├── llm_translator.go    # LLM翻译器 (统一HTTP处理)
│   ├── cache.go             # 分层缓存管理
│   └── fallback.go          # 备用翻译策略
├── generator/           # 内容生成器
│   ├── tag_generator.go     # 标签页面生成
│   └── slug_generator.go    # Slug生成器
├── display/             # 界面显示
│   ├── table.go         # 表格显示
│   └── progress.go      # 进度显示
├── menu/                # 交互菜单系统
│   ├── main_menu.go     # 主菜单
│   └── handlers.go      # 菜单处理器
├── operations/          # 业务操作模块
│   ├── batch_process.go # 批量处理
│   └── workflow.go      # 工作流程
├── utils/               # 工具函数和日志系统
│   ├── logger.go        # 企业级日志系统
│   ├── progress.go      # 进度条工具
│   └── helpers.go       # 辅助函数
├── config.json          # 主配置文件
├── cache/               # 缓存文件目录
│   ├── tag_cache.json   # 标签翻译缓存
│   └── article_cache.json # 文章翻译缓存
├── logs/                # 日志文件目录
└── docs/               # 详细文档
    ├── architecture.md  # 架构设计文档
    ├── performance.md   # 性能优化指南
    └── caching.md       # 缓存策略说明
```

## 🎮 主要功能

### 🚀 快速处理
- 📦 一键处理全部 (智能工作流自动化)
- 🔄 批量缓存预热和优化

### 📝 内容管理
- 🏷️ 生成标签页面 (支持自定义模板)
- 📝 生成文章Slug (SEO优化)
- 🌐 翻译文章为多语言版本 (段落级翻译)

### 💾 缓存管理
- 📊 查看分层缓存状态 (标签/文章分离)
- 🚀 生成全量翻译缓存 (智能批量处理)
- 🗑️ 清空指定类型缓存 (精细化管理)

### 智能特性
- 🤖 AI驱动的上下文感知翻译
- 💾 多层级智能缓存机制
- 🎯 精准内容识别和处理
- 📋 全链路日志追踪和监控
- ⚡ 高性能批量处理引擎

## ⚙️ 配置说明

### 配置文件 (config.json)
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

## 📊 性能监控

### 缓存统计
- 缓存命中率监控
- 分类缓存使用情况
- 内存和磁盘占用分析

### 翻译性能
- 平均翻译时间统计
- 批量处理效率分析
- 错误率和重试统计

### 系统资源
- CPU和内存使用监控
- 网络请求性能追踪
- 磁盘I/O操作统计

## 📝 日志系统

### 企业级功能
- 多级别日志记录 (DEBUG/INFO/WARN/ERROR)
- 自动日志轮转和压缩归档
- 结构化日志格式，便于分析
- 源码位置和调用栈追踪
- 操作审计和性能指标记录

### 监控查看
```bash
# 查看实时日志
tail -f logs/app.log

# 查看翻译性能统计
grep "translation" logs/app.log | grep "PERF"

# 查看缓存操作记录
grep "cache" logs/app.log

# 分析错误趋势
grep "ERROR" logs/app.log | cut -d' ' -f1-2 | sort | uniq -c
```

## 📚 文档链接

### 中文文档
- [架构设计文档](docs/architecture.md)
- [性能优化指南](docs/performance.md)
- [缓存策略说明](docs/caching.md)
- [配置文件说明](docs/configuration.md)
- [日志系统指南](docs/logging.md)
- [故障排除](docs/troubleshooting.md)

### English Documentation
- [Architecture Guide](docs/architecture_en.md)
- [Performance Guide](docs/performance_en.md)
- [Caching Strategy](docs/caching_en.md)
- [Configuration Guide](docs/configuration_en.md)
- [Logging Guide](docs/logging_en.md)
- [Troubleshooting](docs/troubleshooting_en.md)

## 🤝 贡献指南

欢迎提交Issue和Pull Request！详细说明请查看 [贡献指南](docs/contributing.md)

## 📄 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件

---

⭐ 如果这个项目对您有帮助，请给个Star支持一下！
