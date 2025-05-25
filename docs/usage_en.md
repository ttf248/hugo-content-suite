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
