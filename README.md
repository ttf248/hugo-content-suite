# Hugo 博客管理工具

[English](README_EN.md) | 中文

> 🚀 一款专为Hugo博客设计的智能管理工具，支持AI翻译、本地缓存和友好的交互界面

## ✨ 核心特色

### 🤖 AI智能翻译
- 基于LM Studio的本地AI模型翻译
- 中文标签自动转换为SEO友好的英文slug
- 支持批量翻译，提高处理效率

### 💾 智能缓存系统
- 本地JSON缓存，避免重复翻译
- 自动检测缓存状态，节省API调用
- 支持缓存管理和清理功能

### 🎯 友好交互界面
- 直观的命令行菜单系统
- 彩色输出，清晰易读
- 预览功能，安全可靠

### 📊 数据分析统计
- 标签使用频率分析
- 分类统计与可视化
- 无标签文章检测

### 📝 完善日志系统
- 多级别日志记录（DEBUG、INFO、WARN、ERROR）
- 自动日志文件轮转和归档
- 详细的操作记录和错误追踪
- 支持控制台和文件双重输出

### ⚡ 性能监控
- 实时处理速度统计
- 内存使用情况监控
- API调用次数和响应时间统计
- 缓存命中率分析

### ⚙️ 灵活配置管理
- YAML格式配置文件
- 支持环境变量覆盖
- 热重载配置更新
- 配置验证和默认值处理

### 🚀 一键处理功能
- 全流程自动化处理
- 智能状态分析和预览
- 批量缓存生成
- 文章翻译支持

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
1. **一键处理全部**: 自动执行缓存→标签页面→文章Slug→文章翻译
2. **标签分析**: 查看博客标签使用统计
3. **生成标签页面**: 为每个标签创建专门的页面
4. **文章Slug管理**: 为文章标题生成SEO友好的URL
5. **文章翻译**: 将中文文章翻译为英文
6. **缓存管理**: 查看和管理翻译缓存
7. **性能监控**: 查看处理性能和系统资源使用情况
8. **日志分析**: 查看详细的操作日志和错误记录

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
├── utils/               # 工具函数和性能监控
├── config.yaml          # 主配置文件
├── logs/                # 日志文件目录
└── docs/               # 详细文档
```

## 🎮 主要功能

### 🚀 快速处理
- 📦 一键处理全部 (缓存→标签页面→文章Slug→文章翻译)

### 📊 数据查看
- 🔍 标签统计与分析
- 📊 分类统计
- 📝 无标签文章检测

### 🏷️ 标签页面管理
- 👀 预览标签页面
- 🏷️ 生成标签页面

### 📝 文章管理
- 👀 预览文章Slug
- 📝 生成文章Slug
- 👀 预览文章翻译
- 🌐 翻译文章为英文

### 💾 缓存管理
- 📊 查看缓存状态
- 👀 预览全量翻译缓存
- 🚀 生成全量翻译缓存
- 🗑️ 清空翻译缓存 (支持分类)

### 🔧 系统工具
- 📈 查看性能统计
- 🔄 重置性能统计

### 智能特性
- 🤖 AI驱动的翻译
- 💾 智能缓存机制
- 🎯 精准内容识别

### 系统监控
- 📈 实时性能统计
- 📋 详细日志记录
- ⚙️ 灵活配置管理
- 🔍 操作审计追踪

## ⚙️ 配置说明

### 配置文件 (config.yaml)
```yaml
# LM Studio 配置
lm_studio:
  url: "http://localhost:2234/v1/chat/completions"
  model: "gemma-3-12b-it"
  timeout: 30s
  max_retries: 3

# 缓存配置
cache:
  directory: "./cache"
  file_name: "tag_translations_cache.json"
  auto_save: true
  max_entries: 10000

# 日志配置
logging:
  level: "INFO"
  file_path: "./logs/app.log"
  max_size: 100MB
  max_backups: 5
  max_age: 30
  console_output: true

# 性能监控
performance:
  enable_monitoring: true
  metrics_interval: 10s
  memory_threshold: 500MB

# 路径配置
paths:
  default_content_dir: "../../content/post"
```

### 环境变量覆盖
```bash
export LM_STUDIO_URL="http://192.168.1.100:2234/v1/chat/completions"
export LOG_LEVEL="DEBUG"
export CACHE_DIR="./custom_cache"
```

## 📝 日志功能

### 日志级别
- **DEBUG**: 详细的调试信息
- **INFO**: 一般信息记录
- **WARN**: 警告信息
- **ERROR**: 错误信息

### 日志文件管理
- 自动按大小轮转日志文件
- 保留指定数量的历史日志
- 按时间自动清理过期日志

### 日志查看
```bash
# 查看实时日志
tail -f logs/app.log

# 查看错误日志
grep "ERROR" logs/app.log

# 查看特定时间的日志
grep "2024-01-01" logs/app.log
```

## 📈 性能监控

### 实时统计
- 处理速度 (文章/秒)
- 内存使用量
- CPU使用率
- 网络请求延迟

### 性能报告
- 翻译次数统计
- 缓存命中率分析
- 平均翻译时间
- 文件操作次数
- 错误次数统计

## 📚 文档链接

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

## 🤝 贡献指南

欢迎提交Issue和Pull Request！详细说明请查看 [贡献指南](docs/contributing.md)

## 📄 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件

---

⭐ 如果这个项目对您有帮助，请给个Star支持一下！
