# Hugo Blog Manager

English | [中文](README.md)

> 🚀 An intelligent blog management tool designed for Hugo blogs, featuring AI multilingual translation, local caching, and streamlined workflow

## ✨ Key Features

### 🤖 AI-Powered Multilingual Translation
- Local AI model translation based on LM Studio
- Support for multiple languages (English, Japanese, Korean, etc.)
- Automatic conversion of Chinese tags to SEO-friendly English slugs
- Batch translation support for improved efficiency

### 💾 Smart Caching System
- Local JSON caching to avoid duplicate translations
- Automatic cache status detection to save API calls
- Cache management and cleanup functionality

### 🎯 User-Friendly Interface
- Streamlined command-line menu system (7 core features)
- Colorful output for clear readability
- One-click processing for efficiency

### 📝 Advanced Logging System
- Structured logging with multiple levels (DEBUG, INFO, WARN, ERROR)
- Automatic log file rotation and archiving
- Detailed source code location information for troubleshooting
- Support for both console and file output
- Integrated with high-performance logrus library

### 🚀 One-Click Processing
- Full workflow automation
- Intelligent status analysis
- Batch cache generation and multilingual content processing
- Complete blog management solution

## 🚀 Quick Start

### Requirements
- Go 1.21+
- LM Studio (optional, for AI translation)

### Installation & Running
```bash
# Clone the project
git clone https://github.com/ttf248/hugo-content-suite.git
cd hugo-content-suite

# Install dependencies
go mod tidy

# Run the tool
go run main.go [content-directory-path]
```

### Basic Usage
1. **One-Click Process All**: Complete blog processing workflow
2. **Generate Tag Pages**: Create dedicated pages for each tag
3. **Generate Article Slugs**: Add SEO-friendly URLs for articles
4. **Translate Articles**: Translate Chinese articles to multiple languages
5. **Cache Management**: Manage translation cache efficiently

## 🎮 Main Features

### 🚀 Quick Processing
- 📦 One-Click Process All (automatic multilingual workflow)

### 📝 Content Management
- 🏷️ Generate tag pages
- 📝 Generate article slugs
- 🌐 Translate articles to multiple languages

### 💾 Cache Management
- 📊 View cache status
- 🚀 Generate bulk translation cache
- 🗑️ Clear translation cache

### Smart Features
- 🤖 AI-driven multilingual translation
- 💾 Intelligent caching mechanism
- 🎯 Precise content recognition
- 📋 Structured logging with source tracking

## ⚙️ Configuration

### Configuration File (config.json)
```json
{
  "lm_studio": {
    "url": "http://localhost:2234/v1/chat/completions",
    "model": "gemma-3-12b-it",
    "timeout_seconds": 30
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
    "file": "./logs/app.log"
  }
}
```

## 📝 Logging Features

### Structured Logging
- Detailed source code location (filename:line:function)
- Operation tracking and performance statistics
- Specialized logging for translation and cache operations
- Automatic log rotation and compression

### Log Viewing
```bash
# View real-time logs
tail -f logs/app.log

# View error logs
grep "ERROR" logs/app.log

# View translation operations
grep "translation" logs/app.log
```

## 📚 Documentation Links

### 中文文档
- [安装配置指南](docs/installation.md)
- [功能使用说明](docs/usage.md)
- [配置文件说明](docs/configuration.md)
- [故障排除](docs/troubleshooting.md)

### English Documentation
- [Installation Guide](docs/installation_en.md)
- [Usage Guide](docs/usage_en.md)
- [Configuration Guide](docs/configuration_en.md)
- [Troubleshooting](docs/troubleshooting_en.md)

## 🤝 Contributing

Issues and Pull Requests are welcome!

## 📄 License

MIT License - See [LICENSE](LICENSE) file for details

---

⭐ If this project helps you, please give it a Star!
