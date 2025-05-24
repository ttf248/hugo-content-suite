# API Documentation

[中文](api.md) | English

## Core Module Interfaces

### Scanner Module

#### ScanArticles
Scan all Markdown articles in the specified directory.

```go
func ScanArticles(dir string) ([]models.Article, error)
```

**Parameters:**
- `dir`: Directory path to scan

**Returns:**
- `[]models.Article`: List of articles
- `error`: Error information

**Example:**
```go
articles, err := scanner.ScanArticles("./content/post")
if err != nil {
    log.Fatal(err)
}
```

### Translator Module

#### LLMTranslator
AI translator based on LM Studio.

```go
type LLMTranslator struct {
    client  *http.Client
    baseURL string
    model   string
    cache   *TranslationCache
}
```

#### NewLLMTranslator
Create a new translator instance.

```go
func NewLLMTranslator() *LLMTranslator
```

#### TranslateToSlug
Translate Chinese tags to English slugs.

```go
func (t *LLMTranslator) TranslateToSlug(tag string) (string, error)
```

**Parameters:**
- `tag`: Chinese tag to translate

**Returns:**
- `string`: Translated English slug
- `error`: Error information

#### BatchTranslate
Batch translate tags (with caching support).

```go
func (t *LLMTranslator) BatchTranslate(tags []string) (map[string]string, error)
```

**Parameters:**
- `tags`: List of tags to translate

**Returns:**
- `map[string]string`: Mapping from tags to slugs
- `error`: Error information

#### TestConnection
Test connection to LM Studio.

```go
func (t *LLMTranslator) TestConnection() error
```

### Generator Module

#### TagPageGenerator
Tag page generator.

```go
type TagPageGenerator struct {
    contentDir string
    translator *translator.LLMTranslator
    slugCache  map[string]string
}
```

#### NewTagPageGenerator
Create tag page generator.

```go
func NewTagPageGenerator(contentDir string) *TagPageGenerator
```

#### GenerateTagPages
Generate all tag pages.

```go
func (g *TagPageGenerator) GenerateTagPages(tagStats []models.TagStats) error
```

#### PreviewTagPages
Preview tag page generation.

```go
func (g *TagPageGenerator) PreviewTagPages(tagStats []models.TagStats) []TagPagePreview
```

#### ArticleSlugGenerator
Article slug generator.

```go
type ArticleSlugGenerator struct {
    contentDir string
    translator *translator.LLMTranslator
}
```

### Stats Module

#### CalculateTagStats
Calculate tag statistics.

```go
func CalculateTagStats(articles []models.Article) []models.TagStats
```

#### CalculateCategoryStats
Calculate category statistics.

```go
func CalculateCategoryStats(articles []models.Article) []models.CategoryStats
```

#### FindNoTagArticles
Find articles without tags.

```go
func FindNoTagArticles(articles []models.Article) []models.Article
```

## Data Models

### Article
Article model.

```go
type Article struct {
    FilePath string   // File path
    Title    string   // Article title
    Tags     []string // Tag list
    Category string   // Category
    Date     string   // Publication date
}
```

### TagStats
Tag statistics model.

```go
type TagStats struct {
    Name  string   // Tag name
    Count int      // Usage count
    Files []string // Files using this tag
}
```

### CategoryStats
Category statistics model.

```go
type CategoryStats struct {
    Name  string // Category name
    Count int    // Article count
}
```

## LM Studio API

### Request Format

#### Chat Completions
```json
{
  "model": "gemma-3-12b-it",
  "messages": [
    {
      "role": "user",
      "content": "Translation request content"
    }
  ],
  "stream": false
}
```

#### Response Format
```json
{
  "id": "chatcmpl-xxx",
  "object": "chat.completion",
  "created": 1234567890,
  "model": "gemma-3-12b-it",
  "choices": [
    {
      "index": 0,
      "message": {
        "role": "assistant",
        "content": "Translation result"
      }
    }
  ],
  "usage": {
    "prompt_tokens": 50,
    "completion_tokens": 10,
    "total_tokens": 60
  }
}
```

## Cache Interface

### TranslationCache
Translation cache management.

```go
type TranslationCache struct {
    Version      string                `json:"version"`
    LastUpdated  time.Time             `json:"last_updated"`
    Translations map[string]CacheEntry `json:"translations"`
    filePath     string
}
```

#### Main Methods

```go
// Create cache
func NewTranslationCache(cacheDir string) *TranslationCache

// Get translation
func (c *TranslationCache) Get(tag string) (string, bool)

// Set translation
func (c *TranslationCache) Set(tag, translation string)

// Save cache
func (c *TranslationCache) Save() error

// Load cache
func (c *TranslationCache) Load() error
```

## Error Handling

### Common Error Types

```go
// Network connection error
fmt.Errorf("failed to send request: %v", err)

// File read/write error
fmt.Errorf("failed to read file: %v", err)

// JSON parsing error
fmt.Errorf("failed to parse response: %v", err)

// LM Studio error
fmt.Errorf("LM Studio returned error status: %d", resp.StatusCode)
```

### Error Handling Strategy

1. **Network Errors**: Automatic retry mechanism
2. **Translation Failures**: Fallback to predefined mappings
3. **File Errors**: Detailed error prompts
4. **Cache Errors**: Continue execution with warnings

## Configuration Parameters

### Constant Configuration

```go
const (
    LMStudioURL     = "http://172.19.192.1:2234/v1/chat/completions"
    ModelName       = "gemma-3-12b-it"
    CacheFileName   = "tag_translations_cache.json"
)
```

### Adjustable Parameters

```go
// HTTP timeout setting
Timeout: 30 * time.Second

// Request interval
time.Sleep(500 * time.Millisecond)

// Table display limit
defaultLimit := 20
```

## Extension Development

### Adding New Translators

```go
type CustomTranslator struct {
    // Custom fields
}

func (t *CustomTranslator) TranslateToSlug(tag string) (string, error) {
    // Implement translation logic
}
```

### Adding New Generators

```go
type CustomGenerator struct {
    // Custom fields
}

func (g *CustomGenerator) Generate(data interface{}) error {
    // Implement generation logic
}
```
