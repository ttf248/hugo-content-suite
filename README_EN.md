# Hugo Tag Auto Management Tool

English | [中文](README.md)

> 🚀 An intelligent tag management tool designed for Hugo blogs, featuring AI translation, local caching, and user-friendly interactive interface

## ✨ Key Features

### 🤖 AI-Powered Translation
- Local AI model translation based on LM Studio
- Automatic conversion of Chinese tags to SEO-friendly English slugs
- Batch translation support for improved efficiency

### 💾 Smart Caching System
- Local JSON caching to avoid duplicate translations
- Automatic cache status detection to save API calls
- Cache management and cleanup functionality

### 🎯 User-Friendly Interface
- Intuitive command-line menu system
- Colorful output for clear readability
- Preview functionality for safe operations

### 📊 Data Analysis & Statistics
- Tag usage frequency analysis
- Category statistics and visualization
- Detection of articles without tags

### 📝 Comprehensive Logging System
- Multi-level logging (DEBUG, INFO, WARN, ERROR)
- Automatic log file rotation and archiving
- Detailed operation records and error tracking
- Support for both console and file output

### ⚡ Performance Monitoring
- Real-time processing speed statistics
- Memory usage monitoring
- API call count and response time statistics
- Cache hit rate analysis

### ⚙️ Flexible Configuration Management
- YAML format configuration files
- Environment variable override support
- Hot reload configuration updates
- Configuration validation and default value handling

## 🚀 Quick Start

### Requirements
- Go 1.21+
- LM Studio (optional, for AI translation)

### Installation & Running
```bash
# Clone the project
git clone https://github.com/ttf248/hugo-slug-auto.git
cd hugo-slug-auto

# Install dependencies
go mod tidy

# Run the tool
go run main.go [content-directory-path]
```

### Basic Usage
1. **Tag Analysis**: View blog tag usage statistics
2. **Generate Tag Pages**: Create dedicated pages for each tag
3. **Article Slug Management**: Generate SEO-friendly URLs for article titles
4. **Cache Management**: View and manage translation cache
5. **Performance Monitoring**: View processing performance and system resource usage
6. **Log Analysis**: View detailed operation logs and error records

## 📁 Project Architecture

```
hugo-slug-auto/
├── main.go              # Main program entry
├── config/              # Configuration management
├── models/              # Data models
├── scanner/             # Article scanning and parsing
├── stats/               # Statistical analysis
├── translator/          # AI translation module
├── generator/           # Content generators
├── display/             # User interface
├── logger/              # Logging system
├── performance/         # Performance monitoring
├── config.yaml          # Main configuration file
├── logs/                # Log files directory
└── docs/               # Detailed documentation
```

## 🎮 Main Features

### Tag Management
- 📊 Tag statistics analysis
- 🏷️ Automatic tag page generation
- 🔄 Batch translation processing

### Article Management
- 📝 Automatic slug generation
- 🔍 Article content analysis
- 📋 Batch processing support

### Smart Features
- 🤖 AI-driven translation
- 💾 Intelligent caching mechanism
- 🎯 Precise content recognition

### System Monitoring
- 📈 Real-time performance statistics
- 📋 Detailed logging
- ⚙️ Flexible configuration management
- 🔍 Operation audit tracking

## ⚙️ Configuration

### Configuration File (config.yaml)
```yaml
# LM Studio Configuration
lm_studio:
  url: "http://localhost:2234/v1/chat/completions"
  model: "gemma-3-12b-it"
  timeout: 30s
  max_retries: 3

# Cache Configuration
cache:
  directory: "./cache"
  file_name: "tag_translations_cache.json"
  auto_save: true
  max_entries: 10000

# Logging Configuration
logging:
  level: "INFO"
  file_path: "./logs/app.log"
  max_size: 100MB
  max_backups: 5
  max_age: 30
  console_output: true

# Performance Monitoring
performance:
  enable_monitoring: true
  metrics_interval: 10s
  memory_threshold: 500MB
```

### Environment Variable Override
```bash
export LM_STUDIO_URL="http://192.168.1.100:2234/v1/chat/completions"
export LOG_LEVEL="DEBUG"
export CACHE_DIR="./custom_cache"
```

## 📝 Logging Features

### Log Levels
- **DEBUG**: Detailed debugging information
- **INFO**: General information logging
- **WARN**: Warning messages
- **ERROR**: Error messages

### Log File Management
- Automatic log file rotation by size
- Retain specified number of historical logs
- Automatic cleanup of expired logs by time

### Log Viewing
```bash
# View real-time logs
tail -f logs/app.log

# View error logs
grep "ERROR" logs/app.log

# View logs for specific time
grep "2024-01-01" logs/app.log
```

## 📈 Performance Monitoring

### Real-time Statistics
- Processing speed (articles/second)
- Memory usage
- CPU usage
- Network request latency

### Performance Reports
- Total processing time
- Average response time
- Cache hit rate
- Error rate statistics

## 📚 Documentation Links

### 中文文档
- [安装配置指南](docs/installation.md)
- [功能使用说明](docs/usage.md)
- [配置文件说明](docs/configuration.md)
- [日志系统指南](docs/logging.md)
- [性能监控指南](docs/performance.md)
- [API接口文档](docs/api.md)
- [故障排除](docs/troubleshooting.md)
- [贡献指南](docs/contributing.md)

### English Documentation
- [Installation Guide](docs/installation_en.md)
- [Usage Guide](docs/usage_en.md)
- [Configuration Guide](docs/configuration_en.md)
- [Logging Guide](docs/logging_en.md)
- [Performance Guide](docs/performance_en.md)
- [API Documentation](docs/api_en.md)
- [Troubleshooting](docs/troubleshooting_en.md)
- [Contributing Guide](docs/contributing_en.md)

## 🤝 Contributing

Issues and Pull Requests are welcome! Please see the [Contributing Guide](docs/contributing_en.md) for details.

## 📄 License

MIT License - See [LICENSE](LICENSE) file for details

---

⭐ If this project helps you, please give it a Star!
