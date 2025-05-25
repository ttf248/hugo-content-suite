# Usage Guide

English | [中文](usage.md)

## Starting the Program

```bash
go run main.go [content-directory-path]
```

After startup, a streamlined interactive menu interface will be displayed.

## Main Menu Features

### 🚀 Core Functionality

#### 1. One-Click Process All
Automatically execute the complete processing workflow:
1. Generate bulk translation cache
2. Generate new tag pages
3. Generate missing article slugs
4. Translate new articles to English

The system performs intelligent analysis and executes all necessary steps automatically.

### 📝 Content Management

#### 2. Generate Tag Pages
Create dedicated pages for each tag with automatic translation.

Generated file structure:
```
content/tags/
├── tag-name-1/
│   └── _index.md
├── tag-name-2/
│   └── _index.md
└── ...
```

Each `_index.md` contains:
```yaml
---
title: Tag Name
slug: "english-slug"
description: Contains N articles about Tag Name
---
```

#### 3. Generate Article Slugs
Add SEO-friendly URLs for article titles.

Adds slug field in article Front Matter:
```yaml
---
title: "My Article Title"
slug: "my-article-title"
date: 2024-01-01
---
```

#### 4. Translate Articles to English
Translate Chinese articles to English with intelligent content processing.

### 💾 Cache Management

#### 5. View Cache Status
Display translation cache information and statistics:
- Cache file location
- Version information
- Last update time
- Number of translation entries

#### 6. Generate Bulk Translation Cache
Batch translate all tags and article titles for faster processing.

#### 7. Clear Translation Cache
Clear translation cache with support for selective clearing:
- Option to clear tag cache or article cache
- Or clear all cache

## AI Translation Features

### Translation Principles
1. **Chinese Detection**: Automatically identify Chinese content
2. **AI Translation**: Use LM Studio for intelligent translation
3. **Format Standardization**: URL-compliant slug format
4. **Cache Storage**: Local caching to avoid duplicate requests

### Fallback Translation
Predefined translation mappings when AI translation fails:
```go
"人工智能" -> "artificial-intelligence"
"机器学习" -> "machine-learning"
"前端开发" -> "frontend-development"
// ... more mappings
```

## Best Practices

### Usage Recommendations
1. **First Use**: Use "One-Click Process All" for complete setup
2. **Regular Maintenance**: Use individual functions for specific updates
3. **Cache Management**: Monitor cache status and clear when needed

## Related Documentation

- [Configuration Guide](configuration_en.md)
- [Troubleshooting](troubleshooting_en.md)
