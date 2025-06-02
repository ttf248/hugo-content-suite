# Hugo 博客管理工具

[English](README_EN.md) | 中文

> 🚀 一款专为Hugo博客设计的智能管理工具，支持AI多语言翻译、高性能缓存和一体化工作流
> 
> **版本 v3.0.0** - 重构架构，优化性能，提升用户体验

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

## 🏗️ v3.0.0 重构优化亮点

### 代码架构优化

- **统一HTTP客户端**: 重构翻译模块，消除重复代码，提高请求处理效率
- **通用翻译方法**: 使用模板化提示词，支持不同翻译类型（标签、文章、分类等）
- **分层缓存设计**: 标签、Slug、分类分离管理，提高缓存命中率和精准度
- **函数式设计**: 可组合的翻译处理函数，易于扩展和维护
- **处理器架构**: 模块化业务逻辑，统一接口设计

### 性能提升

- **批量处理优化**: 智能分批处理，减少网络开销和API调用次数
- **缓存预加载**: 提前检查缓存状态，减少重复查询和等待时间
- **进度追踪**: 实时显示处理进度，提升用户体验和操作透明度
- **内存优化**: 减少重复对象创建，降低内存占用和GC压力
- **并发控制**: 合理的并发请求限制，避免API限制和资源竞争

### 工程化改进

- **企业级日志**: 集成logrus和lumberjack，支持结构化日志和自动轮转
- **性能监控**: 详细的统计数据和性能指标追踪
- **错误处理**: 智能重试机制和优雅的错误恢复
- **配置管理**: 更完善的配置验证和默认值处理
- **模块化设计**: 清晰的职责分离和组件解耦

## 📁 项目架构

```
hugo-content-suite/
├── main.go              # 主程序入口 - 交互式菜单系统
├── config/              # 配置管理
│   └── config.go        # 配置结构和加载逻辑
├── models/              # 数据模型
│   └── article.go       # 文章、标签、分类统计模型
├── scanner/             # 文章扫描解析
│   └── parser.go        # Markdown文件解析器
├── stats/               # 统计分析
│   └── calculator.go    # 统计数据计算器
├── translator/          # AI翻译模块 (v3.0重构)
│   ├── llm_translator.go    # LLM翻译器 (统一HTTP客户端)
│   ├── cache.go             # 分层缓存管理系统
│   └── translation_utils.go # 翻译工具函数
├── generator/           # 内容生成器 (重构优化)
│   ├── page_generator.go        # 标签和分类页面生成
│   ├── article_slug_generator.go # 文章Slug生成器
│   ├── article_translator.go    # 文章翻译生成器
│   ├── field_translator.go      # 字段翻译处理器
│   └── content_parser.go        # 内容解析器
├── display/             # 界面显示
│   └── tables.go        # 表格和进度显示
├── operations/          # 业务操作模块 (处理器架构)
│   ├── processor.go             # 统一处理器接口
│   ├── article_operations.go    # 文章操作处理
│   ├── article_slug_operations.go # 文章Slug操作
│   ├── article_del_operations.go  # 文章删除操作
│   └── page_operations.go       # 页面生成操作
├── utils/               # 工具函数和系统服务
│   ├── logger.go        # 企业级日志系统 (logrus+lumberjack)
│   ├── progress.go      # 进度条和状态显示
│   ├── performance.go   # 性能监控和统计
│   └── help.go          # 帮助和支持功能
├── config.json          # 主配置文件
├── *_translations_cache.json # 分离式缓存文件
│   ├── tag_translations_cache.json      # 标签翻译缓存
│   ├── slug_translations_cache.json     # Slug翻译缓存
│   └── category_translations_cache.json # 分类翻译缓存
├── markdown/            # 多语言内容示例
└── docs/               # 详细文档
    ├── installation.md     # 中文安装指南
    ├── installation_en.md  # 英文安装指南
    ├── usage.md           # 中文使用说明
    └── usage_en.md        # 英文使用说明
```

## 🎮 主要功能

### 🚀 快速处理

- 📦 一键处理全部 (智能工作流自动化)
- 🔄 批量缓存预热和优化

### 📝 内容管理

- 🏷️ 生成标签和分类页面 (支持自定义模板)
- 📝 生成文章Slug (SEO优化)
- 🌐 翻译文章为多语言版本 (段落级翻译)
- 🔄 文章字段翻译 (标题、描述、标签等)

### 💾 缓存管理

- 📊 查看分层缓存状态 (标签/文章/分类分离)
- 🚀 生成全量翻译缓存 (智能批量处理)
- 🗑️ 清空指定类型缓存 (精细化管理)
- 📈 缓存性能监控和统计

### 🔧 处理器架构

- 🎯 模块化处理器设计 (统一接口)
- 📋 文章操作处理器 (创建、更新、删除)
- 🏷️ 页面生成处理器 (标签页、分类页)
- 🔗 Slug操作处理器 (生成和管理)

### 智能特性

- 🤖 AI驱动的上下文感知翻译
- 💾 多层级智能缓存机制 (标签/Slug/分类)
- 🎯 精准内容识别和处理
- 📋 全链路日志追踪和监控
- ⚡ 高性能批量处理引擎
- 🔄 统一HTTP客户端优化

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

- [安装部署](docs/installation.md)
- [使用说明](docs/usage.md)

### English Documentation

- [Installation Guide](docs/installation_en.md)
- [Usage Instructions](docs/usage_en.md)

## 🤝 贡献指南

欢迎提交Issue和Pull Request！详细说明请查看 [贡献指南](docs/contributing.md)

## 📄 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件

---

⭐ 如果这个项目对您有帮助，请给个Star支持一下！
