# API Documentation

English | [中文](api.md)

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

### Config Module

#### LoadConfig
Load configuration file.

```go
func LoadConfig() (*Config, error)
```

**Returns:**
- `*Config`: Configuration object
- `error`: Error information

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

#### BatchTranslateTags
Batch translate tags.

```go
func (t *LLMTranslator) BatchTranslateTags(tags []string) (map[string]string, error)
```

#### BatchTranslateArticles
Batch translate article titles.

```go
func (t *LLMTranslator) BatchTranslateArticles(titles []string) (map[string]string, error)
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

#### GenerateTagPagesWithMode
Generate tag pages based on mode.

```go
func (g *TagPageGenerator) GenerateTagPagesWithMode(tagStats []models.TagStats, mode string) error
```

**Parameters:**
- `tagStats`: List of tag statistics
- `mode`: Processing mode ("create", "update", "all")

#### ArticleSlugGenerator
Article slug generator.

```go
type ArticleSlugGenerator struct {
    contentDir string
    translator *translator.LLMTranslator
}
```

#### GenerateArticleSlugsWithMode
Generate article slugs based on mode.

```go
func (g *ArticleSlugGenerator) GenerateArticleSlugsWithMode(mode string) error
```

#### ArticleTranslator
Article translator.

```go
type ArticleTranslator struct {
    contentDir string
    translator *translator.LLMTranslator
}
```

#### TranslateArticles
Translate articles.

```go
func (t *ArticleTranslator) TranslateArticles(mode string) error
```

#### PreviewArticleTranslations
Preview article translations.

```go
func (t *ArticleTranslator) PreviewArticleTranslations() ([]ArticleTranslationPreview, error)
```

### Operations Module

#### Processor
Business processor.

```go
type Processor struct {
    contentDir string
}
```

#### NewProcessor
Create processor instance.

```go
func NewProcessor(contentDir string) *Processor
```

#### QuickProcessAll
One-click process all functionality.

```go
func (p *Processor) QuickProcessAll(tagStats []models.TagStats, reader *bufio.Reader)
```

#### PreviewTagPages
Preview tag pages.

```go
func (p *Processor) PreviewTagPages(tagStats []models.TagStats)
```

#### GenerateTagPages
Generate tag pages.

```go
func (p *Processor) GenerateTagPages(tagStats []models.TagStats, reader *bufio.Reader)
```

#### PreviewArticleSlugs
Preview article slugs.

```go
func (p *Processor) PreviewArticleSlugs()
```

#### GenerateArticleSlugs
Generate article slugs.

```go
func (p *Processor) GenerateArticleSlugs(reader *bufio.Reader)
```

#### PreviewArticleTranslations
Preview article translations.

```go
func (p *Processor) PreviewArticleTranslations()
```

#### TranslateArticles
Translate articles.

```go
func (p *Processor) TranslateArticles(reader *bufio.Reader)
```

#### ShowCacheStatus
Show cache status.

```go
func (p *Processor) ShowCacheStatus()
```

#### ShowBulkTranslationPreview
Show bulk translation preview.

```go
func (p *Processor) ShowBulkTranslationPreview(tagStats []models.TagStats)
```

#### GenerateBulkTranslationCache
Generate bulk translation cache.

```go
func (p *Processor) GenerateBulkTranslationCache(tagStats []models.TagStats, reader *bufio.Reader)
```

#### ClearTranslationCache
Clear translation cache.

```go
func (p *Processor) ClearTranslationCache(reader *bufio.Reader)
```

### Menu Module

#### InteractiveMenu
Interactive menu.

```go
type InteractiveMenu struct {
    reader    *bufio.Reader
    processor *operations.Processor
}
```

#### NewInteractiveMenu
Create interactive menu.

```go
func NewInteractiveMenu(reader *bufio.Reader, contentDir string) *InteractiveMenu
```

#### Show
Display streamlined main menu (7 core features).

```go
func (m *InteractiveMenu) Show(tagStats []models.TagStats, categoryStats []models.CategoryStats, noTagArticles []models.Article)
```

### Utils Module

#### Logging Functions

```go
// Logging
func Info(format string, args ...interface{})
func Debug(format string, args ...interface{})
func Warn(format string, args ...interface{})
func Error(format string, args ...interface{})

// Initialize logger
func InitLogger(filename string, level LogLevel) error
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

#### GroupTagsByFrequency
Group tags by frequency.

```go
func GroupTagsByFrequency(tagStats []models.TagStats) ([]models.TagStats, []models.TagStats, []models.TagStats)
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

### PerformanceStats
Performance statistics model.

```go
type PerformanceStats struct {
    TranslationCount int           // Translation count
    TranslationTime  time.Duration // Total translation time
    CacheHits        int           // Cache hit count
    CacheMisses      int           // Cache miss count
    FileOperations   int           // File operation count
    Errors           int           // Error count
}
```

### Config
Configuration model.

```go
type Config struct {
    LMStudio struct {
        URL        string        `yaml:"url"`
        Model      string        `yaml:"model"`
        Timeout    time.Duration `yaml:"timeout"`
        MaxRetries int           `yaml:"max_retries"`
    } `yaml:"lm_studio"`
    
    Cache struct {
        Directory string `yaml:"directory"`
        FileName  string `yaml:"file_name"`
        AutoSave  bool   `yaml:"auto_save"`
    } `yaml:"cache"`
    
    Logging struct {
        Level         string `yaml:"level"`
        FilePath      string `yaml:"file_path"`
        MaxSize       string `yaml:"max_size"`
        MaxBackups    int    `yaml:"max_backups"`
        MaxAge        int    `yaml:"max_age"`
        ConsoleOutput bool   `yaml:"console_output"`
    } `yaml:"logging"`
    
    Paths struct {
        DefaultContentDir string `yaml:"default_content_dir"`
    } `yaml:"paths"`
}
```

## Preview Models

### TagPagePreview
Tag page preview.

```go
type TagPagePreview struct {
    TagName     string
    Slug        string
    Status      string // "create", "update", "skip"
    Description string
}
```

### ArticleSlugPreview
Article slug preview.

```go
type ArticleSlugPreview struct {
    FilePath string
    Title    string
    Slug     string
    Status   string // "missing", "exists"
}
```

### ArticleTranslationPreview
Article translation preview.

```go
type ArticleTranslationPreview struct {
    SourceFile string
    TargetFile string
    Title      string
    Status     string // "missing", "exists"
}
```

### BulkTranslationPreview
Bulk translation preview.

```go
type BulkTranslationPreview struct {
    TagsToTranslate     []TranslationItem
    ArticlesToTranslate []TranslationItem
    MissingTranslations []TranslationItem
}

type TranslationItem struct {
    Original    string
    Translation string
    Type        string // "tag", "article"
}
```

## Common Constants

```go
const (
    LMStudioURL     = "http://172.19.192.1:2234/v1/chat/completions"
    ModelName       = "gemma-3-12b-it"
    CacheFileName   = "tag_translations_cache.json"
)

// Log levels
const (
    DEBUG LogLevel = iota
    INFO
    WARN
    ERROR
)
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

### Adding New Menu Features

```go
func (m *InteractiveMenu) customFunction() {
    // Implement custom functionality
}
```

Add new options and corresponding processing logic in the `Show` method.
