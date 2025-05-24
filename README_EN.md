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

## 📁 Project Architecture

```
hugo-slug-auto/
├── main.go              # Main program entry
├── models/              # Data models
├── scanner/             # Article scanning and parsing
├── stats/               # Statistical analysis
├── translator/          # AI translation module
├── generator/           # Content generators
├── display/             # User interface
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

## 📚 Documentation Links

### 中文文档
- [安装配置指南](docs/installation.md)
- [功能使用说明](docs/usage.md)
- [API接口文档](docs/api.md)
- [故障排除](docs/troubleshooting.md)
- [贡献指南](docs/contributing.md)

### English Documentation
- [Installation Guide](docs/installation_en.md)
- [Usage Guide](docs/usage_en.md)
- [API Documentation](docs/api_en.md)
- [Troubleshooting](docs/troubleshooting_en.md)
- [Contributing Guide](docs/contributing_en.md)

## 🤝 Contributing

Issues and Pull Requests are welcome! Please see the [Contributing Guide](docs/contributing_en.md) for details.

## 📄 License

MIT License - See [LICENSE](LICENSE) file for details

---

⭐ If this project helps you, please give it a Star!
