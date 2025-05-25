# Hugo 博客管理工具

[English](README_EN.md) | 中文

> 🚀 一款专为Hugo博客设计的智能管理工具，支持AI多语言翻译、本地缓存和友好的交互界面

## ✨ 核心特色

### 🤖 AI智能翻译
- 基于LM Studio的本地AI模型翻译
- 支持多语言翻译（英语、日语、韩语等）
- 中文标签自动转换为SEO友好的英文slug
- 支持批量翻译，提高处理效率

### 💾 智能缓存系统
- 本地JSON缓存，避免重复翻译
- 自动检测缓存状态，节省API调用
- 支持缓存管理和清理功能

### 🎯 友好交互界面
- 精简的命令行菜单系统（7个核心功能）
- 彩色输出，清晰易读
- 一键处理功能，提高效率

### 📝 完善日志系统
- 结构化日志记录，支持多级别（DEBUG、INFO、WARN、ERROR）
- 自动日志文件轮转和归档
- 详细的源码位置信息，便于问题排查
- 支持控制台和文件双重输出
- 集成logrus高性能日志库

### 🚀 一键处理功能
- 全流程自动化处理
- 智能状态分析和预览
- 批量缓存生成
- 多语言文章翻译支持

## 🚀 快速开始

### 环境要求
- Go 1.21+
- LM Studio (可选，用于AI翻译)

### 安装运行
```bash
# 克隆项目
git clone https://github.com/ttf248/hugo-content-suite.git
cd hugo-content-suite

# 安装依赖
go mod tidy

# 运行工具
go run main.go [content目录路径]
```

### 基本使用
1. **一键处理全部**: 自动执行缓存→标签页面→文章Slug→多语言翻译
2. **生成标签页面**: 为每个标签创建专门的页面
3. **文章Slug管理**: 为文章标题生成SEO友好的URL
4. **多语言翻译**: 将中文文章翻译为多种语言
5. **缓存管理**: 查看和管理翻译缓存

## 📁 项目架构

```
hugo-content-suite/
├── main.go              # 主程序入口
├── config/              # 配置管理
├── models/              # 数据模型
├── scanner/             # 文章扫描解析
├── stats/               # 统计分析
├── translator/          # AI翻译模块
├── generator/           # 内容生成器
├── display/             # 界面显示
├── menu/                # 交互菜单系统
├── operations/          # 业务操作模块
├── utils/               # 工具函数和日志系统
├── config.json          # 主配置文件
├── logs/                # 日志文件目录
└── docs/               # 详细文档
```

## 🎮 主要功能

### 🚀 快速处理
- 📦 一键处理全部 (缓存→标签页面→文章Slug→多语言翻译)

### 📝 内容管理
- 🏷️ 生成标签页面
- 📝 生成文章Slug
- 🌐 翻译文章为多语言版本

### 💾 缓存管理
- 📊 查看缓存状态
- 🚀 生成全量翻译缓存
- 🗑️ 清空翻译缓存

### 智能特性
- 🤖 AI驱动的多语言翻译
- 💾 智能缓存机制
- 🎯 精准内容识别
- 📋 结构化日志记录

## ⚙️ 配置说明

### 配置文件 (config.json)
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

## 📝 日志功能

### 结构化日志
- 详细的源码位置信息（文件名:行号:函数名）
- 操作追踪和性能统计
- 翻译和缓存操作专门记录
- 自动日志轮转和压缩

### 日志查看
```bash
# 查看实时日志
tail -f logs/app.log

# 查看错误日志
grep "ERROR" logs/app.log

# 查看翻译操作
grep "translation" logs/app.log
```

## 📚 文档链接

### 中文文档
- [安装配置指南](docs/installation.md)
- [功能使用说明](docs/usage.md)


### English Documentation
- [Installation Guide](docs/installation_en.md)
- [Usage Guide](docs/usage_en.md)


## 🤝 贡献指南

欢迎提交Issue和Pull Request！详细说明请查看 [贡献指南](docs/contributing.md)

## 📄 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件

---

⭐ 如果这个项目对您有帮助，请给个Star支持一下！
