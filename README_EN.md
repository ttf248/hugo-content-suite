# Hugo Blog Manager

English | [中文](README.md)

> 🚀 An intelligent blog management tool designed for Hugo blogs, featuring AI multilingual translation, high-performance caching, and integrated workflow automation

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

## 🏗️ Refactoring Highlights

### Code Architecture Optimization
- **Unified HTTP Client**: Eliminated code duplication, improved request processing efficiency
- **Generic Translation Methods**: Template-based prompts supporting different translation types
- **Hierarchical Cache Design**: Separated tag and article management for improved cache hit rates
- **Functional Design**: Composable translation processing functions, easy to extend and maintain

### Performance Improvements
- **Batch Processing Optimization**: Intelligent batch processing to reduce network overhead
- **Cache Preloading**: Early cache status checking to reduce redundant queries
- **Progress Tracking**: Real-time processing progress display for better user experience
- **Memory Optimization**: Reduced duplicate object creation, lower memory footprint

## 📁 Project Architecture

```
hugo-content-suite/
├── main.go              # Program entry point
├── config/              # Configuration management
│   ├── config.go        # Configuration structure and loading
│   └── validation.go    # Configuration validation
├── models/              # Data models
│   ├── article.go       # Article model
│   └── metadata.go      # Metadata structures
├── scanner/             # Article scanning and parsing
│   ├── scanner.go       # File scanner
│   └── parser.go        # Markdown parser
├── stats/               # Statistical analysis
│   ├── collector.go     # Data collector
│   └── reporter.go      # Statistics reporter
├── translator/          # AI translation module (refactored)
│   ├── llm_translator.go    # LLM translator (unified HTTP handling)
│   ├── cache.go             # Hierarchical cache management
│   └── fallback.go          # Fallback translation strategy
├── generator/           # Content generators
│   ├── tag_generator.go     # Tag page generator
│   └── slug_generator.go    # Slug generator
├── display/             # Interface display
│   ├── table.go         # Table display
│   └── progress.go      # Progress display
├── menu/                # Interactive menu system
│   ├── main_menu.go     # Main menu
│   └── handlers.go      # Menu handlers
├── operations/          # Business operation modules
│   ├── batch_process.go # Batch processing
│   └── workflow.go      # Workflow management
├── utils/               # Utilities and logging system
│   ├── logger.go        # Enterprise logging system
│   ├── progress.go      # Progress bar utilities
│   └── helpers.go       # Helper functions
├── config.json          # Main configuration file
├── cache/               # Cache file directory
│   ├── tag_cache.json   # Tag translation cache
│   └── article_cache.json # Article translation cache
├── logs/                # Log file directory
└── docs/               # Detailed documentation
    ├── architecture.md  # Architecture design documentation
    ├── performance.md   # Performance optimization guide
    └── caching.md       # Caching strategy documentation
```

## 🎮 Main Features

### 🚀 Quick Processing
- 📦 One-Click Process All (intelligent workflow automation)
- 🔄 Batch cache warming and optimization

### 📝 Content Management
- 🏷️ Generate tag pages (custom template support)
- 📝 Generate article slugs (SEO optimization)
- 🌐 Translate articles to multiple languages (paragraph-level translation)

### 💾 Cache Management
- 📊 View hierarchical cache status (tag/article separation)
- 🚀 Generate bulk translation cache (intelligent batch processing)
- 🗑️ Clear specific cache types (fine-grained management)

### Smart Features
- 🤖 AI-driven context-aware translation
- 💾 Multi-tier intelligent caching mechanism
- 🎯 Precise content identification and processing
- 📋 Full-chain log tracking and monitoring
- ⚡ High-performance batch processing engine

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
