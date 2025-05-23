# Hugo 标签自动化管理工具

> 🚀 一款专为Hugo博客设计的智能标签管理工具，支持AI翻译、本地缓存和友好的交互界面

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

## 🚀 快速开始

### 环境要求
- Go 1.21+
- LM Studio (可选，用于AI翻译)

### 安装运行
```bash
# 克隆项目
git clone https://github.com/ttf248/hugo-slug-auto.git
cd hugo-slug-auto

# 安装依赖
go mod tidy

# 运行工具
go run main.go [content目录路径]
```

### 基本使用
1. **标签分析**: 查看博客标签使用统计
2. **生成标签页面**: 为每个标签创建专门的页面
3. **文章Slug管理**: 为文章标题生成SEO友好的URL
4. **缓存管理**: 查看和管理翻译缓存

## 📁 项目架构

```
hugo-slug-auto/
├── main.go              # 主程序入口
├── models/              # 数据模型
├── scanner/             # 文章扫描解析
├── stats/               # 统计分析
├── translator/          # AI翻译模块
├── generator/           # 内容生成器
├── display/             # 界面显示
└── docs/               # 详细文档
```

## 🎮 主要功能

### 标签管理
- 📊 标签统计分析
- 🏷️ 自动生成标签页面
- 🔄 批量翻译处理

### 文章管理
- 📝 Slug自动生成
- 🔍 文章内容分析
- 📋 批量处理支持

### 智能特性
- 🤖 AI驱动的翻译
- 💾 智能缓存机制
- 🎯 精准内容识别

## 📚 文档链接

- [安装配置指南](docs/installation.md)
- [功能使用说明](docs/usage.md)
- [API接口文档](docs/api.md)
- [故障排除](docs/troubleshooting.md)

## 🤝 贡献指南

欢迎提交Issue和Pull Request！详细说明请查看 [贡献指南](docs/contributing.md)

## 📄 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件

---

⭐ 如果这个项目对您有帮助，请给个Star支持一下！
