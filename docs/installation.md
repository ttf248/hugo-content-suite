# 安装配置指南

## 系统要求

### 必需环境
- **Go**: 版本 1.21 或更高
- **操作系统**: Windows, macOS, Linux
- **Hugo博客**: 支持Front Matter的Markdown文件

### 可选组件
- **LM Studio**: 用于AI翻译功能
- **Git**: 用于版本控制

## 安装步骤

### 1. 克隆项目
```bash
git clone https://github.com/ttf248/hugo-slug-auto.git
cd hugo-slug-auto
```

### 2. 安装依赖
```bash
go mod tidy
```

### 3. 验证安装
```bash
go run main.go --help
```

## LM Studio 配置

### 安装 LM Studio
1. 访问 [LM Studio官网](https://lmstudio.ai/) 下载
2. 安装并启动LM Studio
3. 下载推荐的语言模型（如 Gemma-3-12B-IT）

### 配置连接
在 `translator/llm_translator.go` 中修改配置：
```go
const (
    LMStudioURL = "http://localhost:2234/v1/chat/completions"  // 修改为你的LM Studio地址
    ModelName   = "your-model-name"                           // 修改为你的模型名称
)
```

### 测试连接
```bash
go run main.go
# 选择菜单项 "6. 预览文章Slug" 来测试AI翻译功能
```

## 配置选项

### 缓存设置
默认缓存文件保存在当前目录下的 `tag_translations_cache.json`

修改缓存位置：
```go
// 在 translator/llm_translator.go 的 NewLLMTranslator 函数中
cache: NewTranslationCache("./your-cache-directory"),
```

### 翻译超时设置
```go
// 在 translator/llm_translator.go 中
client: &http.Client{
    Timeout: 30 * time.Second,  // 修改超时时间
},
```

## 使用不同内容目录

### 默认目录
```bash
go run main.go  # 使用默认路径 ../../content/post
```

### 自定义目录
```bash
go run main.go /path/to/your/content
```

### Windows 路径示例
```bash
go run main.go "C:\Users\Username\myblog\content\post"
```

## 验证配置

### 检查文件结构
确保你的Hugo博客具有以下结构：
```
your-blog/
├── content/
│   ├── post/           # 文章目录
│   │   ├── article1.md
│   │   └── article2.md
│   └── tags/           # 标签页面目录（工具会自动创建）
└── ...
```

### 检查文章格式
确保Markdown文件包含Front Matter：
```yaml
---
title: "文章标题"
date: 2024-01-01
tags: ["标签1", "标签2"]
categories: ["分类"]
---

文章内容...
```

## 常见问题

### Go版本问题
如果遇到Go版本兼容问题：
```bash
go version  # 检查当前版本
go mod edit -go=1.21  # 修改go.mod中的版本要求
```

### 依赖问题
如果依赖下载失败：
```bash
go clean -modcache
go mod download
go mod tidy
```

### 权限问题
确保有读写权限：
- 内容目录的读取权限
- 标签目录的写入权限
- 缓存文件的写入权限

## 下一步

安装完成后，请查看 [使用说明](usage.md) 了解具体功能使用方法。
