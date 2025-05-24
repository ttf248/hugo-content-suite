# Usage Guide

English | [中文](usage.md)

## Starting the Program

```bash
go run main.go [content-directory-path]
```

After startup, an interactive menu interface will be displayed.

## Main Menu Features

### 🚀 Quick Processing

#### 1. One-Click Process All
Automatically execute the complete processing workflow:
1. Generate bulk translation cache
2. Generate new tag pages
3. Generate missing article slugs
4. Translate new articles to English

The system will first perform intelligent analysis and display preview of content to be processed:
- 📊 Number of items to translate
- 🏷️ Number of tag pages to create
- 📝 Number of articles needing slugs
- 🌐 Number of articles to translate

After confirmation, all steps are executed automatically with progress display and status feedback for each step.

### 📊 Data Viewing Module

#### 2. Tag Statistics & Analysis
- **View All Tags**: Display all tags and their usage frequency
- **View Specific Tag Details**: Enter tag name to view all articles using that tag
- **View by Frequency Groups**: 
  - High-frequency tags (≥5 articles)
  - Medium-frequency tags (2-4 articles)  
  - Low-frequency tags (1 article)

#### 3. Category Statistics
Display all categories with article count and percentage information

#### 4. Articles Without Tags
List all articles that have no tags for subsequent tag addition

### 🏷️ Tag Page Management

#### 5. Preview Tag Pages
- Scan all tags and generate translation preview
- Display directory structure to be created
- Show statistics for new/updated pages

#### 6. Generate Tag Pages
- **Processing Mode Selection**:
  - Create Only: Only create pages for tags without existing pages
  - Update Only: Only update existing tag pages
  - Process All: Create + update all tag pages

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

### 📝 Article Management Module

#### 7. Preview Article Slug
- Scan all articles and analyze slug status
- Display count of articles missing slugs
- Preview generated slug examples (first 5)

#### 8. Generate Article Slug
- **Processing Mode Selection**:
  - Create Only: Only add slugs for articles missing them
  - Update Only: Only update existing slugs
  - Process All: Create + update all articles

Adds or updates slug field in article Front Matter:
```yaml
---
title: "My Article Title"
slug: "my-article-title"
date: 2024-01-01
---
```

#### 9. Preview Article Translations
- Scan all Chinese articles
- Check if corresponding English versions exist
- Display list of articles needing translation

#### 10. Translate Articles to English
- **Translation Mode Selection**:
  - Only translate missing articles
  - Re-translate existing articles
  - Translate all articles

Translation features:
- Intelligent detection of Chinese articles
- Generate corresponding English directory structure
- Translate titles, content, and Front Matter
- Maintain original format and structure

### 💾 Cache Management Module

#### 11. View Cache Status
Display detailed translation cache information:
- Cache file location
- Version information
- Last update time
- Number of translation entries

#### 12. Preview Bulk Translation Cache
- Analyze all content needing translation
- Distinguish between tag and article translations
- Display list of missing cache items
- Estimate translation workload

#### 13. Generate Bulk Translation Cache
- Batch translate all tags and article titles
- Process tag and article translations separately
- Real-time translation progress display
- Automatically save to cache file

#### 14. Clear Translation Cache (with categorization support)
- Support categorized cache clearing
- Option to clear tag cache or article cache
- Or clear all cache

### 🔧 System Tools Module

#### 15. View Performance Statistics
Display detailed performance monitoring information:
```
📊 Performance Statistics:
🔄 Translation Count: 156
⚡ Cache Hit Rate: 87.5% (140/160)
⏱️ Average Translation Time: 1.2s
📁 File Operations: 89
❌ Error Count: 2
```

#### 16. Reset Performance Statistics
Clear all performance statistics and restart counting.

## AI Translation Features

### Translation Principles
1. **Chinese Detection**: Automatically identify Chinese tags and articles
2. **AI Translation**: Use LM Studio for intelligent translation
3. **Format Standardization**: Automatically convert to URL-compliant slug format
4. **Cache Storage**: Translation results cached locally to avoid duplicate requests

### Translation Rules
- Use lowercase letters
- Connect words with hyphens (-)
- Remove special characters
- Keep concise and accurate

### Fallback Translation
When AI translation fails, automatically use predefined translation mappings:
```go
"人工智能" -> "artificial-intelligence"
"机器学习" -> "machine-learning"
"前端开发" -> "frontend-development"
// ... more mappings
```

## Batch Processing

### Batch Tag Translation
- Automatically collect all tags that need translation
- Use cache to reduce duplicate translations
- Display processing progress and statistics

### Batch Article Processing
- Scan all Markdown files in specified directory
- Parse Front Matter content
- Batch generate or update slug fields

### Article Translation Processing
- Intelligent detection of Chinese articles
- Batch translate article content
- Generate English version files
- Maintain consistent directory structure

## Caching Mechanism

### Cache File Format
```json
{
  "version": "1.0",
  "last_updated": "2024-01-01T12:00:00Z",
  "translations": {
    "tag-name": {
      "translation": "tag-name",
      "created_at": "2024-01-01T12:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z"
    }
  }
}
```

### Cache Strategy
- Write to cache on first translation
- Read directly from cache for subsequent use
- Support manual clearing for re-translation
- Distinguish between tag cache and article cache

## Output Format

### Statistical Tables
Use colorful tables to display statistical data:
- 🔴 Red: High-frequency tags/create operations
- 🟡 Yellow: Medium-frequency tags/update operations  
- 🔵 Blue: Low-frequency tags/skip operations
- 🟢 Green: Success status

### Progress Indicators
- Real-time processing progress display
- Detailed error message prompts
- Operation result statistics summary

## Best Practices

### Usage Recommendations
1. **First Use**: Use "One-Click Process All" feature for quick initialization
2. **Daily Maintenance**: Regularly run data viewing functions to understand blog status
3. **Tag Management**: Regularly check and clean up low-frequency tags
4. **Backup**: Backup important files before batch operations
5. **Testing**: Test on small scale before batch operations

### Performance Optimization
1. **Cache Utilization**: Make full use of translation cache to reduce API calls
2. **Batch Processing**: Process large numbers of articles in batches
3. **Network Stability**: Ensure stable LM Studio connection
4. **Usage Monitoring**: Regularly check performance statistics to optimize usage

## Troubleshooting

For common issues, please refer to [Troubleshooting Guide](troubleshooting_en.md)

## Related Documentation

- [Configuration Guide](configuration_en.md)
- [Logging Guide](logging_en.md)
- [Performance Guide](performance_en.md)
- [API Documentation](api_en.md)
