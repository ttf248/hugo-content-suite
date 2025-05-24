# Usage Guide

[中文](usage.md) | English

## Starting the Program

```bash
go run main.go [content-directory-path]
```

After startup, an interactive menu interface will be displayed.

## Main Menu Features

### 📊 Data Viewing Module

#### 1. Tag Statistics & Analysis
- **View All Tags**: Display all tags and their usage frequency
- **View Specific Tag Details**: Enter tag name to view all articles using that tag
- **View by Frequency Groups**: 
  - High-frequency tags (≥5 articles)
  - Medium-frequency tags (2-4 articles)  
  - Low-frequency tags (1 article)

#### 2. Category Statistics
Display all categories with article count and percentage information

#### 3. Articles Without Tags
List all articles that have no tags for subsequent tag addition

### 🏷️ Tag Page Management

#### 4. Preview Tag Pages
- Scan all tags and generate translation preview
- Display directory structure to be created
- Show statistics for new/updated pages

#### 5. Generate Tag Pages
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

### 📝 Article Slug Management

#### 6. Preview Article Slug
- Scan all articles and analyze slug status
- Display count of articles missing slugs
- Preview generated slug examples (first 5)

#### 7. Generate Article Slug
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

### 💾 Cache Management

#### 8. View Cache Status
Display detailed translation cache information:
- Cache file location
- Version information
- Last update time
- Number of translation entries

#### 9. Clear Translation Cache
Clear all cached translations; next translation will request AI service again

## AI Translation Features

### Translation Principles
1. **Chinese Detection**: Automatically identify Chinese tags
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
1. **First Use**: Run data viewing features first to understand blog status
2. **Tag Management**: Regularly check and clean up low-frequency tags
3. **Backup**: Backup important files before batch operations
4. **Testing**: Test on small scale before batch operations

### Performance Optimization
1. **Cache Utilization**: Make full use of translation cache to reduce API calls
2. **Batch Processing**: Process large numbers of articles in batches
3. **Network Stability**: Ensure stable LM Studio connection

## Troubleshooting

For common issues, please refer to [Troubleshooting Guide](troubleshooting_en.md)
