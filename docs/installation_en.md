# Installation Guide

[中文](installation.md) | English

## System Requirements

### Required Environment
- **Go**: Version 1.21 or higher
- **Operating System**: Windows, macOS, Linux
- **Hugo Blog**: Markdown files with Front Matter support

### Optional Components
- **LM Studio**: For AI translation functionality
- **Git**: For version control

## Installation Steps

### 1. Clone the Project
```bash
git clone https://github.com/ttf248/hugo-content-suite.git
cd hugo-content-suite
```

### 2. Install Dependencies
```bash
go mod tidy
```

### 3. Verify Installation
```bash
go run main.go --help
```

## LM Studio Configuration

### Install LM Studio
1. Visit [LM Studio official website](https://lmstudio.ai/) to download
2. Install and launch LM Studio
3. Download recommended language model (e.g., Gemma-3-12B-IT)

### Configure Connection
Modify configuration in `translator/llm_translator.go`:
```go
const (
    LMStudioURL = "http://localhost:2234/v1/chat/completions"  // Change to your LM Studio address
    ModelName   = "your-model-name"                           // Change to your model name
)
```

### Test Connection
```bash
go run main.go
# Select menu item "6. Preview Article Slug" to test AI translation functionality
```

## Configuration Options

### Cache Settings
Default cache file is saved as `tag_translations_cache.json` in the current directory

To change cache location:
```go
// In translator/llm_translator.go NewLLMTranslator function
cache: NewTranslationCache("./your-cache-directory"),
```

### Translation Timeout Settings
```go
// In translator/llm_translator.go
client: &http.Client{
    Timeout: 30 * time.Second,  // Modify timeout duration
},
```

## Using Different Content Directories

### Default Directory
```bash
go run main.go  # Uses default path ../../content/post
```

### Custom Directory
```bash
go run main.go /path/to/your/content
```

### Windows Path Example
```bash
go run main.go "C:\Users\Username\myblog\content\post"
```

## Configuration Verification

### Check File Structure
Ensure your Hugo blog has the following structure:
```
your-blog/
├── content/
│   ├── post/           # Articles directory
│   │   ├── article1.md
│   │   └── article2.md
│   └── tags/           # Tag pages directory (auto-created by tool)
└── ...
```

### Check Article Format
Ensure Markdown files contain Front Matter:
```yaml
---
title: "Article Title"
date: 2024-01-01
tags: ["tag1", "tag2"]
categories: ["category"]
---

Article content...
```

## Common Issues

### Go Version Issues
If encountering Go version compatibility issues:
```bash
go version  # Check current version
go mod edit -go=1.21  # Modify version requirement in go.mod
```

### Dependency Issues
If dependency download fails:
```bash
go clean -modcache
go mod download
go mod tidy
```

### Permission Issues
Ensure proper permissions:
- Read permission for content directory
- Write permission for tags directory
- Write permission for cache file

## Next Steps

After installation, please refer to [Usage Guide](usage_en.md) to learn about specific feature usage.
