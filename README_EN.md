# Hugo Blog Manager

English | [中文](README.md)

> 🚀 An intelligent blog management tool designed for Hugo blogs, featuring AI multilingual translation, high-performance caching, and integrated workflow automation
> 
> **Version v3.0.0** - Refactored architecture, optimized performance, enhanced user experience

## ✨ Key Features

### 🤖 AI-Powered Multilingual Translation
- Local AI model translation based on LM Studio
- Support for multiple languages (English, Japanese, Korean, etc.)
- Automatic conversion of Chinese tags to SEO-friendly English slugs
- Unified HTTP request handling for improved stability and performance
- Intelligent caching mechanism with categorized tag and article management

### 💾 High-Performance Caching System
- Local JSON hierarchical caching to avoid duplicate translations
- Separated cache management (tag cache, article cache)
- Automatic cache status detection with intelligent batch processing
- Cache statistics and cleanup functionality
- Cache hit rate monitoring and performance optimization

### 🎯 Optimized User Experience
- Streamlined command-line menu system
- Colorful output with progress bar display
- One-click processing with full workflow automation
- Intelligent error handling and retry mechanisms
- Detailed operation statistics and performance reporting

### 📝 Enterprise-Grade Logging System
- Structured logging with multi-level output support
- Automatic log file rotation and compression archiving
- Detailed source code location and call stack information
- Support for both console and file output
- Integrated with high-performance logrus library
- Operation auditing and performance monitoring

### 🚀 One-Click Processing
- Full workflow automation processing
- Intelligent status analysis and pre-processing checks
- Batch cache generation and optimization
- Multilingual article translation support
- Complete blog management solution

## 🏗️ v3.0.0 Refactoring Highlights

### Code Architecture Optimization

- **Unified HTTP Client**: Refactored translation module, eliminated code duplication, improved request processing efficiency
- **Generic Translation Methods**: Template-based prompts supporting different translation types (tags, articles, categories, etc.)
- **Hierarchical Cache Design**: Separated tag, slug, and category management for improved cache hit rates and precision
- **Functional Design**: Composable translation processing functions, easy to extend and maintain
- **Processor Architecture**: Modular business logic with unified interface design

### Performance Improvements

- **Batch Processing Optimization**: Intelligent batch processing to reduce network overhead and API call frequency
- **Cache Preloading**: Early cache status checking to reduce redundant queries and wait times
- **Progress Tracking**: Real-time processing progress display for better user experience and operation transparency
- **Memory Optimization**: Reduced duplicate object creation, lower memory footprint and GC pressure
- **Concurrency Control**: Reasonable concurrent request limits to avoid API limitations and resource contention

### Engineering Improvements

- **Enterprise Logging**: Integrated logrus and lumberjack with structured logging and automatic rotation
- **Performance Monitoring**: Detailed statistics and performance metrics tracking
- **Error Handling**: Intelligent retry mechanisms and graceful error recovery
- **Configuration Management**: Enhanced configuration validation and default value handling
- **Modular Design**: Clear separation of responsibilities and component decoupling

## 📁 Project Architecture

```
hugo-content-suite/
├── main.go              # Program entry point - Interactive menu system
├── config/              # Configuration management
│   └── config.go        # Configuration structure and loading logic
├── models/              # Data models
│   └── article.go       # Article, tag, and category statistics models
├── scanner/             # Article scanning and parsing
│   └── parser.go        # Markdown file parser
├── stats/               # Statistical analysis
│   └── calculator.go    # Statistics data calculator
├── translator/          # AI translation module (v3.0 refactored)
│   ├── llm_translator.go    # LLM translator (unified HTTP client)
│   ├── cache.go             # Hierarchical cache management system
│   └── translation_utils.go # Translation utility functions
├── generator/           # Content generators (refactored)
│   ├── page_generator.go        # Tag and category page generator
│   ├── article_slug_generator.go # Article slug generator
│   ├── article_translator.go    # Article translation generator
│   ├── field_translator.go      # Field translation processor
│   └── content_parser.go        # Content parser
├── display/             # Interface display
│   └── tables.go        # Table and progress display
├── operations/          # Business operation modules (processor architecture)
│   ├── processor.go             # Unified processor interface
│   ├── article_operations.go    # Article operation processor
│   ├── article_slug_operations.go # Article slug operations
│   ├── article_del_operations.go  # Article deletion operations
│   └── page_operations.go       # Page generation operations
├── utils/               # Utilities and system services
│   ├── logger.go        # Enterprise logging system (logrus+lumberjack)
│   ├── progress.go      # Progress bar and status display
│   ├── performance.go   # Performance monitoring and statistics
│   └── help.go          # Help and support functions
├── config.json          # Main configuration file
├── *_translations_cache.json # Separated cache files
│   ├── tag_translations_cache.json      # Tag translation cache
│   ├── slug_translations_cache.json     # Slug translation cache
│   └── category_translations_cache.json # Category translation cache
├── markdown/            # Multilingual content examples
└── docs/               # Detailed documentation
    ├── installation.md     # Chinese installation guide
    ├── installation_en.md  # English installation guide
    ├── usage.md           # Chinese usage instructions
    └── usage_en.md        # English usage instructions
```

## 🎮 Main Features

### 🚀 Quick Processing

- 📦 One-Click Process All (intelligent workflow automation)
- 🔄 Batch cache warming and optimization

### 📝 Content Management

- 🏷️ Generate tag and category pages (custom template support)
- 📝 Generate article slugs (SEO optimization)
- 🌐 Translate articles to multiple languages (paragraph-level translation)
- 🔄 Article field translation (titles, descriptions, tags, etc.)

### 💾 Cache Management

- 📊 View hierarchical cache status (tag/article/category separation)
- 🚀 Generate bulk translation cache (intelligent batch processing)
- 🗑️ Clear specific cache types (fine-grained management)
- 📈 Cache performance monitoring and statistics

### 🔧 Processor Architecture

- 🎯 Modular processor design (unified interface)
- 📋 Article operation processors (create, update, delete)
- 🏷️ Page generation processors (tag pages, category pages)
- 🔗 Slug operation processors (generation and management)

### Smart Features

- 🤖 AI-driven context-aware translation
- 💾 Multi-tier intelligent caching mechanism (tag/slug/category)
- 🎯 Precise content identification and processing
- 📋 Full-chain log tracking and monitoring
- ⚡ High-performance batch processing engine
- 🔄 Unified HTTP client optimization

## ⚙️ Configuration

### Configuration File (config.json)
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

## 📊 Performance Monitoring

### Cache Statistics
- Cache hit rate monitoring
- Categorized cache usage analysis
- Memory and disk usage analysis

### Translation Performance
- Average translation time statistics
- Batch processing efficiency analysis
- Error rate and retry statistics

### System Resources
- CPU and memory usage monitoring
- Network request performance tracking
- Disk I/O operation statistics

## 📝 Logging Features

### Enterprise-Grade Functionality
- Multi-level logging (DEBUG/INFO/WARN/ERROR)
- Automatic log rotation and compression archiving
- Structured log format for easy analysis
- Source code location and call stack tracking
- Operation auditing and performance metrics recording

### Monitoring and Viewing
```bash
# View real-time logs
tail -f logs/app.log

# View translation performance statistics
grep "translation" logs/app.log | grep "PERF"

# View cache operation records
grep "cache" logs/app.log

# Analyze error trends
grep "ERROR" logs/app.log | cut -d' ' -f1-2 | sort | uniq -c
```

## 📚 Documentation Links

### 中文文档
- [架构设计文档](docs/architecture.md)
- [性能优化指南](docs/performance.md)
- [缓存策略说明](docs/caching.md)
- [配置文件说明](docs/configuration.md)
- [故障排除](docs/troubleshooting.md)

### English Documentation
- [Architecture Guide](docs/architecture_en.md)
- [Performance Guide](docs/performance_en.md)
- [Caching Strategy](docs/caching_en.md)
- [Configuration Guide](docs/configuration_en.md)
- [Logging Guide](docs/logging_en.md)
- [Troubleshooting](docs/troubleshooting_en.md)

## 🤝 Contributing

Issues and Pull Requests are welcome!

## 📄 License

MIT License - See [LICENSE](LICENSE) file for details

---

⭐ If this project helps you, please give it a Star!
