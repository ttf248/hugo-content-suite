# Installation Guide

English | [中文](installation.md)

> **Version v3.0.0** - Refactored architecture, enterprise logging, high-performance caching

## System Requirements

### Required Environment
- **Go**: Version 1.22.0 or higher (recommended toolchain 1.23.4)
- **Operating System**: Windows, macOS, Linux
- **Hugo Blog**: Markdown files with Front Matter support
- **Memory**: Recommended 4GB+ (for large blog batch processing)
- **Disk Space**: At least 200MB (including hierarchical cache and rotating log files)

### Optional Components
- **LM Studio**: For AI translation functionality (Highly recommended)
  - Recommended models: gemma-3-12b-it, llama-3.1, qwen-2.5, etc.
- **Git**: For version control
- **Visual Studio Code**: Recommended for viewing structured logs and configuration files

## Quick Installation

### 1. Clone the Project
```bash
git clone https://github.com/your-org/hugo-content-suite.git
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

### 4. First Run
```bash
go run main.go [your-content-directory-path]
```

On first run, the program will automatically create:
- Default configuration file `config.json`
- Log directory `logs/`
- Hierarchical cache files (`*_translations_cache.json`)

## v3.0.0 New Features

### 🏗️ Refactored Architecture
- **Processor Pattern**: Modular business logic with unified interface design
- **Hierarchical Cache**: Separated tag/slug/category management for improved precision
- **Unified HTTP Client**: Eliminated code duplication, enhanced performance

### 📝 Enterprise Logging
- **Structured Logging**: JSON format for easy analysis and monitoring
- **Automatic Rotation**: Log file compression and archiving
- **Multi-level Output**: DEBUG/INFO/WARN/ERROR level control
- **Performance Monitoring**: Integrated operation statistics and performance metrics

### ⚡ Performance Optimization
- **Batch Processing**: Intelligent batching reduces API call frequency
- **Cache Preloading**: Early status checking reduces wait times
- **Memory Optimization**: Reduced memory usage by ~30%
- **Concurrency Control**: Configurable concurrent request limits

## Configuration File

### Auto-generated Configuration File
The program generates `config.json` in the project root directory on first run:

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

### Configuration Options Details

#### LM Studio Configuration (lm_studio)
- `url`: LM Studio API endpoint
- `model`: AI model name to use
- `timeout_seconds`: Request timeout duration
- `max_retries`: Maximum retry attempts
- `retry_delay_ms`: Retry delay in milliseconds

#### Cache Configuration (cache)
- `auto_save_count`: Auto-save interval
- `delay_ms`: Delay between requests
- `expire_days`: Cache expiration in days
- `enable_compression`: Enable cache compression

#### Performance Configuration (performance)
- `max_concurrent_requests`: Maximum concurrent requests
- `batch_size`: Batch processing size
- `memory_limit_mb`: Memory limit

## LM Studio Configuration

### Install LM Studio
1. Visit [LM Studio official website](https://lmstudio.ai/) to download
2. Install and launch LM Studio
3. Download recommended language models:
   - **Recommended**: Gemma-3-12B-IT (balanced performance and quality)
   - **Alternative**: LLaMA2-7B (faster speed)
   - **High Quality**: GPT-4 (if you have API access)

### Configure Connection
Modify LM Studio configuration in `config.json`:

```json
{
  "lm_studio": {
    "url": "http://192.168.1.100:2234/v1/chat/completions",  // Change to your LM Studio address
    "model": "your-model-name",                               // Change to your model name
    "timeout_seconds": 45,                                    // Adjust based on network conditions
    "max_retries": 5                                          // Increase retries for unstable networks
  }
}
```

### Test Connection
```bash
go run main.go
# The program will automatically test LM Studio connection on startup
# Or select menu items to test translation
```

## Directory Structure Configuration

### Recommended Project Structure
```
your-hugo-blog/
├── content/
│   ├── post/              # Articles directory
│   │   ├── article1.md
│   │   └── article2.md
│   └── tags/              # Tag pages directory (auto-created by tool)
│       ├── ai/
│       └── tech/
├── hugo-content-suite/    # This tool directory
│   ├── config.json        # Configuration file
│   ├── cache/             # Cache directory (auto-created)
│   │   ├── tag_cache.json
│   │   └── article_cache.json
│   ├── logs/              # Log directory (auto-created)
│   │   └── app.log
│   └── ...
└── ...
```

### Using Different Content Directories

#### Default Directory
```bash
go run main.go  # Program will prompt for content directory path
```

#### Direct Directory Specification
```bash
go run main.go /path/to/your/content/post
```

#### Windows Path Example
```bash
go run main.go "C:\Users\Username\myblog\content\post"
```

#### Relative Path Example
```bash
go run main.go ../content/post
```

## Advanced Configuration

### Performance Optimization Configuration
Configuration recommendations for different use cases:

#### Large Blog (1000+ articles)
```json
{
  "performance": {
    "max_concurrent_requests": 3,
    "batch_size": 50,
    "memory_limit_mb": 1024
  },
  "cache": {
    "auto_save_count": 20,
    "enable_compression": true
  }
}
```

#### Fast Processing Mode
```json
{
  "performance": {
    "max_concurrent_requests": 10,
    "batch_size": 100,
    "memory_limit_mb": 2048
  },
  "lm_studio": {
    "timeout_seconds": 15,
    "max_retries": 1
  }
}
```

#### Stability-First Mode
```json
{
  "performance": {
    "max_concurrent_requests": 1,
    "batch_size": 10,
    "memory_limit_mb": 256
  },
  "lm_studio": {
    "timeout_seconds": 60,
    "max_retries": 10,
    "retry_delay_ms": 2000
  }
}
```

### Logging Configuration
```json
{
  "logging": {
    "level": "DEBUG",        // Use DEBUG for development, INFO for production
    "file": "./logs/app.log",
    "max_size_mb": 200,      // Increase log file size for large blogs
    "max_backups": 30,       // Keep more backup files
    "console_output": false  // Disable console output in production
  }
}
```

## Installation Verification

### Check File Structure
Ensure your Hugo blog has the correct file structure:

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
Ensure Markdown files contain complete Front Matter:

```yaml
---
title: "Article Title"
date: 2024-01-01T12:00:00+08:00
tags: ["AI", "Technology", "Programming"]
categories: ["Development"]
slug: ""                    # Optional, tool will auto-generate
author: "Author Name"
description: "Article description"
---

Article content...
```

### Verify Functionality
Run the following commands to verify various functions:

```bash
# 1. Verify basic functionality
go run main.go /path/to/content

# 2. Check configuration file
cat config.json

# 3. View generated directory structure
ls -la cache/
ls -la logs/

# 4. Test translation functionality (if LM Studio is configured)
# Select "Generate Bulk Translation Cache" in the program menu
```

## Troubleshooting

### Common Issues

#### 1. Go Version Issues
```bash
go version  # Check current version
# If version is too old, upgrade to 1.21+
```

#### 2. Dependency Issues
```bash
go clean -modcache
go mod download
go mod tidy
```

#### 3. Permission Issues
Ensure necessary permissions:
```bash
# Linux/macOS
chmod 755 hugo-content-suite/
chmod 666 config.json

# Windows (run as administrator)
icacls hugo-content-suite /grant Everyone:F
```

#### 4. LM Studio Connection Issues
- Check if LM Studio is running
- Verify the port is correct (default 2234)
- Test network connection:
```bash
curl -X POST http://localhost:2234/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{"model":"test","messages":[{"role":"user","content":"test"}]}'
```

#### 5. Cache-Related Issues
```bash
# Clear cache
rm -rf cache/
mkdir cache

# Check disk space
df -h
```

### Corrupted Configuration File
If configuration file is corrupted:
```bash
# Delete configuration file, program will recreate default configuration
rm config.json
go run main.go
```

### Log Viewing
View detailed log information:
```bash
# View real-time logs
tail -f logs/app.log

# View error logs
grep "ERROR" logs/app.log

# View performance information
grep "PERF" logs/app.log
```

## Next Steps

### Recommended Workflow
1. **After Installation**: View [Usage Guide](usage_en.md)
2. **Configuration Optimization**: Refer to [Configuration Guide](configuration_en.md)
3. **Performance Tuning**: Check [Performance Optimization Guide](performance_en.md)
4. **Troubleshooting**: Refer to [Troubleshooting Guide](troubleshooting_en.md)

### Advanced Usage
- [Architecture Guide](architecture_en.md) - Understand system architecture
- [Caching Strategy](caching_en.md) - Optimize cache usage
- [Logging Guide](logging_en.md) - Monitoring and debugging

---

After installation, we recommend using the "One-Click Process All" feature to experience the complete workflow!
