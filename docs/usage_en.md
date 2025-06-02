# Usage Guide

English | [中文](usage.md)

## Starting the Program

```bash
go run main.go [content-directory-path]
```

After startup, a streamlined interactive menu interface with intelligent workflow and high-performance batch processing capabilities will be displayed.

## Main Menu Features

### 🚀 Quick Processing

#### 1. One-Click Process All (Intelligent Workflow)
Automatically execute the complete processing workflow using refactored architecture:
1. **Intelligent Status Analysis**: Pre-check cache status and file structure
2. **Batch Cache Warming**: Use hierarchical cache system to generate bulk translation cache
3. **Tag Page Generation**: Template-based processing with custom template support
4. **Article Slug Optimization**: Use unified translation methods for SEO-friendly URLs
5. **Multilingual Translation**: Paragraph-level intelligent translation with context awareness

The system performs intelligent analysis and displays detailed processing preview:
- 📊 Hierarchical cache status (tag cache/article cache)
- 🏷️ Number and types of tags to process
- 📝 Number of articles needing slug addition/update
- 🌐 Articles requiring translation and estimated time
- ⚡ Performance optimization suggestions and batch processing strategies

After confirmation, automatic execution with real-time progress tracking and performance monitoring.

### 📝 Content Management (Refactored & Optimized)

#### 2. Generate Tag Pages (High-Performance Batch Processing)
Using refactored tag generator with support for:
- **Intelligent Batch Translation**: Unified HTTP client reduces network overhead
- **Template-based Generation**: Support for custom tag page templates
- **Cache Optimization**: Auto-detect existing translation cache to avoid duplicate requests
- **Progress Monitoring**: Real-time processing progress and performance metrics

Processing modes:
- **New Only**: Create pages only for tags without existing pages
- **Update Only**: Update only existing tag pages
- **Process All**: Create new + update all tag pages

Generated file structure:
```
content/tags/
├── tag-name-1/
│   └── _index.md
├── tag-name-2/
│   └── _index.md
└── ...
```

Each `_index.md` contains optimized metadata:
```yaml
---
title: Tag Name
slug: "english-slug"
description: Contains N articles about Tag Name
layout: "tag"
---
```

#### 3. Generate Article Slugs (Unified Translation Processing)
Using refactored article slug generator with features:
- **Unified Translation Methods**: Use dedicated `TranslateToArticleSlug` interface
- **Batch Optimization**: Intelligent batch processing reduces API calls
- **Hierarchical Cache Management**: Independent article cache improves hit rates
- **Intelligent Retry Mechanism**: Auto-fallback to backup translation on failure

Add or update optimized slug fields for articles:
```yaml
---
title: "My Article Title"
slug: "my-article-title"
date: 2024-01-01
tags: ["AI", "Technology"]
categories: ["Development"]
---
```

#### 4. Translate Articles to Multiple Languages (Context-Aware Translation)
Using refactored translation engine with support for:
- **Paragraph-level Translation**: Maintain article structure integrity
- **Context Awareness**: AI understands technical terms and professional vocabulary
- **Format Preservation**: Auto-preserve Markdown formatting and code blocks
- **Batch Processing**: Intelligent batch translation for optimized performance

Translation features:
- Auto-detect content to skip (code blocks, links, etc.)
- Intelligently handle mixed Chinese-English content
- Maintain original directory structure and file relationships
- Generate Hugo-compliant multilingual files

### 💾 Cache Management (Hierarchical Architecture)

#### 5. View Cache Status (Hierarchical Display)
Display detailed hierarchical cache information:
```
📊 Cache Status Overview:
🏷️ Tag Translation Cache:
  - Total entries: 145
  - Hit rate: 92.3%
  - Last updated: 2024-01-15 14:30
  - Cache size: 45.2 KB

📝 Article Translation Cache:
  - Total entries: 89
  - Hit rate: 87.6%
  - Last updated: 2024-01-15 14:25
  - Cache size: 32.1 KB

📂 Category Translation Cache:
  - Total entries: 23
  - Hit rate: 95.7%
  - Last updated: 2024-01-15 14:20
  - Cache size: 8.9 KB

💾 Total Cache Performance:
  - Combined hit rate: 89.2%
  - Total size: 86.2 KB
  - Memory usage: 12.3 MB
```

#### 6. Generate Bulk Translation Cache (Intelligent Batch Processing)
Using refactored batch processing engine:
- **Interactive Confirmation**: Detailed preview with user confirmation before processing
- **Categorized Batch Processing**: Separate processing for tags, slugs, and categories
- **Intelligent Deduplication**: Auto-detect existing cache, translate only missing items
- **Performance Optimization**: Connection pooling and concurrency control
- **Progress Tracking**: Real-time progress display and statistics

Processing workflow:
```
🔍 Analyzing content to translate...
📊 Found 45 tags, 89 article titles, 12 categories
💾 Checking existing cache...
🔄 Need to translate: 12 tags, 23 articles, 3 categories

🏷️ Batch translate tags (12/12) ████████████ 100%
📝 Batch translate articles (23/23) ████████████ 100%
📂 Batch translate categories (3/3) ████████████ 100%

✅ Batch translation completed!
  - New tag translations: 12
  - New article translations: 23
  - New category translations: 3
  - Total time: 45.6s
  - Average per item: 1.2s
```

#### 7. Clear Translation Cache (Fine-grained Management)
Support categorized cache clearing:
- **Safety Confirmation**: User confirmation required before clearing
- **Tag Cache Clear**: Clear only tag translation cache
- **Article Cache Clear**: Clear only article translation cache
- **Category Cache Clear**: Clear only category translation cache
- **Full Cache Clear**: Clear all translation caches
- **Expired Cache Cleanup**: Auto-clean expired cache entries

### 📊 Performance Monitoring (Enterprise-grade)

#### New: View Performance Statistics
Display detailed performance monitoring information:
```
📊 Performance Statistics:
🔄 Translation Operations:
  - Total translations: 156
  - Tag translations: 89
  - Article translations: 67
  - Average response time: 1.2s

💾 Cache Performance:
  - Total queries: 1,234
  - Cache hits: 1,078 (87.4%)
  - Cache misses: 156 (12.6%)
  - Average query time: 0.02ms

🌐 Network Performance:
  - HTTP requests: 156
  - Success rate: 98.7% (154/156)
  - Average request time: 1.18s
  - Retry count: 3

📁 File Operations:
  - File reads: 245
  - File writes: 89
  - Cache saves: 12
  - Average I/O time: 0.05s

❌ Error Statistics:
  - Network errors: 2
  - Translation errors: 0
  - File errors: 1
  - Total error rate: 1.9%
```

#### New: Reset Performance Statistics
Clear all performance statistics and restart monitoring.

## AI Translation Features (Refactored & Optimized)

### Unified Translation Architecture
Using refactored translator architecture:
1. **Unified HTTP Client**: Eliminates code duplication, improves request efficiency
2. **Template-based Prompts**: Specialized prompts for different content types
3. **Intelligent Cache Strategy**: Hierarchical cache management improves hit rates
4. **Optimized Error Handling**: Smart retry and fallback strategies

### Translation Method Types
- **TranslateToSlug**: Tag translation, generates URL-safe slugs
- **TranslateToArticleSlug**: Article title translation, maintains semantic accuracy
- **TranslateParagraph**: Paragraph translation, preserves format and context

### Intelligent Content Detection
Refactored content detection mechanism:
```go
// Auto-skip content that doesn't need translation
- Code blocks (content surrounded by ```)
- Indented code blocks (lines starting with 4 spaces)
- Pure English content
- HTML tags
- Image links (but translate Chinese in alt text)
- Quote blocks (but translate Chinese content within)
```

### Batch Processing Optimization
- **Smart Batching**: Auto-adjust batch size based on network conditions and API limits
- **Cache Preloading**: Pre-check cache status to reduce duplicate queries
- **Progress Monitoring**: Real-time progress display and estimated completion time
- **Error Recovery**: Single translation failures don't affect overall process

## Cache Mechanism (Hierarchical Architecture)

### Hierarchical Cache Design
```json
// Tag Cache (tag_translations_cache.json)
{
  "version": "3.0",
  "cache_type": "tag",
  "last_updated": "2024-01-15T14:30:25Z",
  "entries": {
    "人工智能": {
      "translation": "artificial-intelligence",
      "created_at": "2024-01-15T10:00:00Z",
      "updated_at": "2024-01-15T10:00:00Z",
      "hit_count": 15
    }
  },
  "stats": {
    "total_entries": 145,
    "hit_rate": 0.923,
    "avg_hit_count": 8.2
  }
}

// Article Cache (slug_translations_cache.json)
{
  "version": "3.0", 
  "cache_type": "article",
  "last_updated": "2024-01-15T14:25:12Z",
  "entries": {
    "如何使用AI提升开发效率": {
      "translation": "how-to-improve-development-efficiency-with-ai",
      "created_at": "2024-01-15T09:30:00Z",
      "updated_at": "2024-01-15T09:30:00Z",
      "hit_count": 3
    }
  },
  "stats": {
    "total_entries": 89,
    "hit_rate": 0.876,
    "avg_hit_count": 4.1
  }
}

// Category Cache (category_translations_cache.json)
{
  "version": "3.0",
  "cache_type": "category", 
  "last_updated": "2024-01-15T14:20:15Z",
  "entries": {
    "技术教程": {
      "translation": "tech-tutorials",
      "created_at": "2024-01-15T09:15:00Z",
      "updated_at": "2024-01-15T09:15:00Z",
      "hit_count": 8
    }
  },
  "stats": {
    "total_entries": 23,
    "hit_rate": 0.957,
    "avg_hit_count": 6.4
  }
}
```

### Cache Performance Optimization
- **Compressed Storage**: Auto-compress cache files to reduce disk usage
- **Expiration Management**: Auto-clean expired cache entries
- **Hit Rate Monitoring**: Real-time cache usage efficiency statistics

## v3.0.0 Processor Architecture

### Unified Processor Interface
All business operations use standardized processor interface design:

```go
type Processor interface {
    Process() error
    GetType() string
    GetDescription() string
}
```

### Processor Types

#### 1. Article Processor (ArticleProcessor)
- **Function**: Handle article slug generation and translation
- **Features**: Support batch processing, intelligent caching
- **Configuration**: Configurable concurrency and batch size

#### 2. Page Processor (PageProcessor) 
- **Function**: Generate tag and category pages
- **Features**: Template-based generation, custom template support
- **Optimization**: Incremental updates, process only changed pages

#### 3. Delete Processor (DeleteProcessor)
- **Function**: Safely delete generated pages and cache
- **Features**: Preview deletion content, execute after confirmation
- **Protection**: Prevent accidental deletion of important files

### Error Handling Mechanism
- **Unified Error Handling**: All processors use same error handling strategy
- **Error Recovery**: Support recovery from errors, continue subsequent operations
- **Detailed Logging**: Record detailed error information and call stacks

## Logging System (Enterprise-grade)

### Structured Log Format
```json
{
  "timestamp": "2024-01-15T14:30:25.123Z",
  "level": "INFO",
  "module": "translator",
  "operation": "translate_tag",
  "message": "Tag translation completed",
  "data": {
    "tag": "人工智能",
    "translation": "artificial-intelligence",
    "duration_ms": 1250,
    "cache_hit": false
  },
  "performance": {
    "memory_mb": 45.2,
    "cpu_percent": 12.5,
    "goroutines": 8
  }
}
```

### Log Level Description
- **DEBUG**: Detailed debugging information, including function calls and variable states
- **INFO**: General operation information, including processing progress and results
- **WARN**: Warning information, non-fatal errors and performance issues
- **ERROR**: Error information, issues requiring user attention

### Log Rotation Configuration
```json
{
  "logging": {
    "level": "INFO",
    "file": "./logs/app.log",
    "max_size_mb": 100,     // Max 100MB per file
    "max_backups": 10,      // Keep 10 backup files
    "max_age_days": 30,     // Keep logs within 30 days
    "compress": true,       // Compress old log files
    "console_output": true  // Also output to console
  }
}
```

## Advanced Features

### Custom Template Support
v3.0.0 supports custom page templates:

#### Tag Page Template (templates/tag_page.md)
```markdown
---
title: "{{.Name}} - Tag Page"
slug: "{{.Slug}}"
description: "{{.Count}} technical articles related to {{.Name}}"
type: "tag"
layout: "tag"
---

# {{.Name}}

{{.Description}}

> Total **{{.Count}}** articles about {{.Name}}

## Latest Articles

{{range .Articles}}
- [{{.Title}}]({{.URL}}) - {{.Date}}
{{end}}

## Related Tags

{{range .RelatedTags}}
- [{{.Name}}]({{.URL}})
{{end}}
```

#### Category Page Template (templates/category_page.md)
```markdown
---
title: "{{.Name}} - Category Page"
slug: "{{.Slug}}"
description: "{{.Count}} articles in {{.Name}} category"
type: "category"
layout: "category"
---

# {{.Name}} Category

{{.Description}}

## Article List ({{.Count}} articles)

{{range .Articles}}
### [{{.Title}}]({{.URL}})
{{.Summary}}
*Published: {{.Date}}* | *Tags: {{.Tags}}*

---
{{end}}
```

### Batch Operation Strategy

#### Intelligent Batch Processing
```go
// Auto-adjust based on content volume and network conditions
Small operations (< 10 items):  concurrency=2, batch=5
Medium operations (10-50 items): concurrency=3, batch=10
Large operations (> 50 items):  concurrency=5, batch=20
```

#### Performance Monitoring Metrics
- **Throughput**: Items processed per second
- **Latency**: Average response time for single operations
- **Error Rate**: Percentage of failed operations
- **Resource Usage**: CPU and memory consumption

### Multi-language Support Configuration

#### Supported Target Languages
```json
{
  "language": {
    "target_languages": ["en", "ja", "ko", "fr", "de", "es"],
    "language_names": {
      "en": "English",
      "ja": "日本語", 
      "ko": "한국어",
      "fr": "Français",
      "de": "Deutsch",
      "es": "Español"
    },
    "fallback_language": "en"
  }
}
```

#### Translation Quality Control
- **Terminology Consistency**: Auto-detect and maintain consistent professional term translations
- **Context Awareness**: Adjust translation strategy based on article type
- **Format Preservation**: Perfect preservation of Markdown formatting and special characters
- **Quality Assessment**: Auto-evaluate translation quality and provide improvement suggestions

## Performance Optimization Recommendations

### System Configuration Optimization

#### Hardware Configuration Recommendations
```
Small blog (< 100 articles):
- CPU: 2+ cores
- Memory: 4GB+
- Disk: SSD 100MB available space

Medium blog (100-500 articles):
- CPU: 4+ cores  
- Memory: 8GB+
- Disk: SSD 500MB available space

Large blog (> 500 articles):
- CPU: 8+ cores
- Memory: 16GB+
- Disk: NVMe SSD 1GB available space
```

#### Network Optimization Configuration
```json
{
  "performance": {
    "max_concurrent_requests": 5,
    "request_timeout_seconds": 30,
    "retry_delay_ms": 1000,
    "max_retries": 3,
    "connection_pool_size": 10,
    "keep_alive_seconds": 300
  }
}
```

### Usage Recommendations

#### First-time Use Workflow
1. **Environment Check**: Ensure LM Studio is running properly
2. **Configuration Optimization**: Adjust performance parameters based on blog size
3. **Cache Warming**: First execute "Generate Bulk Translation Cache"
4. **Incremental Processing**: Use "One-Click Process All" for new content
5. **Regular Maintenance**: Regularly clean expired cache and logs

#### Daily Maintenance Workflow
```bash
# Daily automation script example
# 1. Clean expired cache
go run main.go --cleanup-cache

# 2. Check system status  
go run main.go --health-check

# 3. Process new content
go run main.go /path/to/content --auto-process

# 4. Generate statistics report
go run main.go --generate-report
```

#### Best Practices
- **Backup Strategy**: Regularly backup configuration files and cache data
- **Monitoring Setup**: Configure log monitoring and performance alerts
- **Version Control**: Include generated content in Git version control
- **Testing Validation**: Validate in test environment before production

## Troubleshooting

### Common Issue Diagnosis

#### LM Studio Connection Issues
```bash
# Check connection status
curl -X POST http://localhost:2234/v1/chat/completions \
     -H "Content-Type: application/json" \
     -d '{"model":"test","messages":[{"role":"user","content":"hello"}]}'

# View error logs
grep "LM Studio" logs/app.log | tail -10
```

#### Cache Issue Diagnosis
```bash
# Check cache file integrity
go run main.go --validate-cache

# Rebuild corrupted cache
rm *_translations_cache.json
go run main.go --rebuild-cache
```

#### Performance Issue Diagnosis
```bash
# View performance logs
grep "performance" logs/app.log | tail -20

# Analyze memory usage
grep "memory" logs/app.log | tail -10

# Check disk space
df -h
```

### Error Code Reference
- **E001**: LM Studio connection failed
- **E002**: Configuration file format error
- **E003**: Cache file corrupted
- **E004**: Network request timeout
- **E005**: Insufficient memory
- **E006**: Insufficient disk space
- **E007**: Insufficient permissions

---

## Upgrade Guide

### Upgrading from v2.x to v3.0.0

#### 1. Backup Existing Data
```bash
# Backup configuration and cache
cp config.json config.json.v2.backup
cp translations_cache.json translations_cache.json.v2.backup
```

#### 2. Migrate Configuration Format
v3.0.0 will auto-detect old configuration and prompt for upgrade:
```bash
go run main.go --migrate-config
```

#### 3. Rebuild Cache Structure
```bash
# Delete old single cache file
rm translations_cache.json

# Run program to auto-create new hierarchical cache
go run main.go /path/to/content
```

### New Feature Migration

#### Enable Enterprise Logging
```json
{
  "logging": {
    "level": "INFO",
    "file": "./logs/app.log",
    "max_size_mb": 100,
    "max_backups": 10,
    "console_output": true
  }
}
```

#### Configure Processor Parameters
```json
{
  "processors": {
    "article_processor": {
      "batch_size": 20,
      "concurrent_workers": 5
    },
    "page_processor": {
      "template_dir": "./templates",
      "output_format": "hugo"
    }
  }
}
```

---

After completing installation and configuration, we recommend using the "One-Click Process All" feature to experience the complete intelligent workflow!

For any issues, please check the log files or GitHub Issues.
Display refactored hierarchical cache information:
```
📊 Cache Status Overview:
🏷️ Tag Cache:
  - File: cache/tag_cache.json
  - Entries: 145 translations
  - Hit Rate: 92.3%
  - Last Updated: 2024-01-15 14:30:25

📝 Article Cache:
  - File: cache/article_cache.json  
  - Entries: 89 translations
  - Hit Rate: 87.6%
  - Last Updated: 2024-01-15 14:25:12

💾 Cache Performance:
  - Total Memory Usage: 2.3MB
  - Disk Usage: 456KB (compressed)
  - Average Query Time: 0.02ms
```

#### 6. Generate Bulk Translation Cache (Intelligent Batch Processing)
Using refactored batch processing engine:
- **Categorized Batch Processing**: Handle tag and article translations separately
- **Intelligent Deduplication**: Auto-detect existing cache, translate only missing items
- **Performance Optimization**: Use connection pooling and concurrency control
- **Progress Tracking**: Real-time translation progress and statistics

Processing flow:
```
🔍 Analyzing content to translate...
📊 Found 45 tags, 89 article titles
💾 Checking existing cache...
🔄 Need to translate: 12 tags, 23 article titles

🏷️ Batch translating tags (12/12) ████████████ 100%
📝 Batch translating articles (23/23) ████████████ 100%

✅ Batch translation completed!
  - New tag translations: 12
  - New article translations: 23
  - Total time: 45.6s
  - Average per item: 1.3s
```

#### 7. Clear Translation Cache (Fine-grained Management)
Support categorized cache clearing:
- **Clear Tag Cache**: Clear only tag translation cache
- **Clear Article Cache**: Clear only article translation cache
- **Clear All Cache**: Clear all translation caches
- **Clean Expired Cache**: Auto-clean expired cache entries

### 📊 Performance Monitoring (Enterprise-Grade)

#### New: View Performance Statistics
Display detailed performance monitoring information:
```
📊 Performance Statistics Detail:
🔄 Translation Operations:
  - Total Translations: 156
  - Tag Translations: 89 times
  - Article Translations: 67 times
  - Average Response Time: 1.2s

💾 Cache Performance:
  - Total Queries: 1,234
  - Cache Hits: 1,078 (87.4%)
  - Cache Misses: 156 (12.6%)
  - Average Query Time: 0.02ms

🌐 Network Performance:
  - HTTP Requests: 156
  - Success Rate: 98.7% (154/156)
  - Average Request Time: 1.18s
  - Retry Count: 3

📁 File Operations:
  - File Reads: 245
  - File Writes: 89
  - Cache Saves: 12
  - Average I/O Time: 0.05s

❌ Error Statistics:
  - Network Errors: 2
  - Translation Errors: 0
  - File Errors: 1
  - Total Error Rate: 1.9%
```

#### New: Reset Performance Statistics
Clear all performance statistics data and restart monitoring.

## AI Translation Features (Refactored & Optimized)

### Unified Translation Architecture
Using refactored translator architecture:
1. **Unified HTTP Client**: Eliminate code duplication, improve request efficiency
2. **Template-based Prompts**: Specialized prompts for different content types
3. **Intelligent Cache Strategy**: Hierarchical cache management improves hit rates
4. **Optimized Error Handling**: Intelligent retry and backup strategies

### Translation Method Types
- **TranslateToSlug**: Tag translation, generate URL-safe slugs
- **TranslateToArticleSlug**: Article title translation, maintain semantic accuracy
- **TranslateParagraph**: Paragraph translation, preserve format and context

### Intelligent Content Detection
Refactored content detection mechanism:
```go
// Auto-skip content that doesn't need translation
- Code blocks (content surrounded by ```)
- Indented code blocks (starting with 4 spaces)
- Pure English content
- HTML tags
- Image links (but translate Chinese in alt text)
- Quote blocks (but translate Chinese content within)
```

### Batch Processing Optimization
- **Intelligent Batching**: Auto-adjust batch size based on network conditions and API limits
- **Cache Preloading**: Pre-check cache status to reduce redundant queries
- **Progress Monitoring**: Real-time translation progress and estimated completion time
- **Error Recovery**: Single translation failures don't affect overall process

## Cache Mechanism (Hierarchical Architecture)

### Hierarchical Cache Design
```json
// Tag Cache (cache/tag_cache.json)
{
  "version": "2.0",
  "cache_type": "tag",
  "last_updated": "2024-01-15T14:30:25Z",
  "entries": {
    "人工智能": {
      "translation": "artificial-intelligence",
      "created_at": "2024-01-15T10:00:00Z",
      "updated_at": "2024-01-15T10:00:00Z",
      "hit_count": 15
    }
  },
  "stats": {
    "total_entries": 145,
    "hit_rate": 0.923,
    "avg_hit_count": 8.2
  }
}

// Article Cache (cache/article_cache.json)
{
  "version": "2.0", 
  "cache_type": "article",
  "last_updated": "2024-01-15T14:25:12Z",
  "entries": {
    "如何使用AI提升开发效率": {
      "translation": "how-to-improve-development-efficiency-with-ai",
      "created_at": "2024-01-15T09:30:00Z",
      "updated_at": "2024-01-15T09:30:00Z",
      "hit_count": 3
    }
  },
  "stats": {
    "total_entries": 89,
    "hit_rate": 0.876,
    "avg_hit_count": 4.1
  }
}
```

### Cache Performance Optimization
- **Compressed Storage**: Auto-compress cache files to reduce disk usage
- **Expiration Management**: Auto-clean expired cache entries
- **Hit Rate Monitoring**: Real-time cache usage efficiency statistics
- **Preloading Strategy**: Preload related cache during batch operations

## Logging and Monitoring (Enterprise-Grade)

### Structured Logging
```
2024-01-15 14:30:25.123 [INFO] [translator] Batch translation started
  type=tag_batch count=12 estimated_time=15s source=llm_translator.go:156

2024-01-15 14:30:26.456 [PERF] [cache] Cache operation completed  
  operation=batch_load hit_rate=87.3% duration=0.02s source=cache.go:89

2024-01-15 14:30:27.789 [INFO] [workflow] One-click processing progress update
  step=2/4 progress=50% current_operation="Generate tag pages" source=workflow.go:67
```

### Performance Monitoring Metrics
- **Translation Performance**: Average translation time, success rate, error rate
- **Cache Performance**: Hit rate, query time, memory usage
- **Network Performance**: Request latency, retry count, timeout rate
- **File Operations**: I/O time, read/write success rate

## Output Format (Enhanced Display)

### Enhanced Progress Bars
```
🔄 Batch translating tags (8/12) ████████░░░░ 67% [ETA: 6s]
  Current: "机器学习" -> "machine-learning"
  Speed: 1.2 items/sec | Errors: 0 | Cache hits: 4/8
```

### Performance Dashboard
```
📊 Real-time Performance Monitoring:
┌─────────────────┬──────────┬──────────┬──────────┐
│     Metric      │ Current  │ Average  │   Best   │
├─────────────────┼──────────┼──────────┼──────────┤
│ Translation/s   │   1.2    │   1.1    │   1.8    │
│ Cache Hit Rate  │  87.3%   │  85.2%   │  92.1%   │
│ Network Latency │  1.18s   │  1.25s   │  0.89s   │
│ Memory Usage    │   45MB   │   42MB   │   38MB   │
└─────────────────┴──────────┴──────────┴──────────┘
```

## Best Practices (Performance Optimization)

### Usage Recommendations
1. **First Use**: Use "One-Click Process All" feature, system auto-optimizes processing order
2. **Cache Management**: Regularly check cache status, maintain high hit rate (>80%)
3. **Batch Operations**: Use batch processing for large numbers of articles to improve efficiency
4. **Performance Monitoring**: Regularly check performance statistics to identify bottlenecks
5. **Network Optimization**: Ensure stable LM Studio connection, avoid frequent retries

### Performance Optimization Tips
1. **Cache Warming**: Generate bulk cache before batch operations
2. **Batch Processing**: Set appropriate batch size (default 20) to avoid API limits
3. **Concurrency Control**: Use reasonable concurrency (default 5) to balance speed and stability
4. **Memory Management**: Regularly clean unnecessary cache to control memory usage
5. **Error Handling**: Enable intelligent retry to improve operation success rate

### Advanced Configuration
```json
{
  "performance": {
    "max_concurrent_requests": 5,
    "batch_size": 20,
    "memory_limit_mb": 512,
    "cache_hit_rate_threshold": 0.8,
    "enable_compression": true,
    "auto_cleanup": true
  },
  "monitoring": {
    "enable_performance_logging": true,
    "metrics_interval_seconds": 60,
    "slow_operation_threshold_ms": 5000
  }
}
```

## Troubleshooting

### Common Performance Issues
1. **Low Cache Hit Rate**: Check cache file integrity, consider regenerating cache
2. **Slow Translation Speed**: Check network connection, adjust concurrency settings
3. **High Memory Usage**: Clean expired cache, adjust batch processing size
4. **Frequent Retries**: Check LM Studio status, optimize network configuration

For detailed troubleshooting, please refer to [Troubleshooting Guide](troubleshooting_en.md)

## Related Documentation

- [Architecture Guide](architecture_en.md)
- [Performance Guide](performance_en.md)
- [Caching Strategy](caching_en.md)
- [Configuration Guide](configuration_en.md)
- [Logging Guide](logging_en.md)
