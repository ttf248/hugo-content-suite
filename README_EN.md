# Hugo Blog Manager

English | [中文](README.md)

> 🚀 An intelligent blog management tool designed for Hugo blogs, featuring AI translation, local caching, and streamlined workflow

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
- Streamlined command-line menu system
- Colorful output for clear readability
- One-click processing for efficiency

### 🚀 One-Click Processing
- Full workflow automation
- Intelligent status analysis
- Batch cache generation and content processing
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
4. **Translate Articles**: Translate Chinese articles to English
5. **Cache Management**: Manage translation cache efficiently

## 🎮 Main Features

### 🚀 Quick Processing
- 📦 One-Click Process All (automatic workflow)

### 📝 Content Management
- 🏷️ Generate tag pages
- 📝 Generate article slugs
- 🌐 Translate articles to English

### 💾 Cache Management
- 📊 View cache status
- 🚀 Generate bulk translation cache
- 🗑️ Clear translation cache

### Smart Features
- 🤖 AI-driven translation
- 💾 Intelligent caching mechanism
- 🎯 Precise content recognition

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

# Path Configuration
paths:
  default_content_dir: "../../content/post"
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
