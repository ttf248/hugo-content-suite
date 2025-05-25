# API 接口文档

[English](api_en.md) | 中文

## 核心模块接口

### Scanner 模块

#### ScanArticles
扫描指定目录下的所有Markdown文章。

```go
func ScanArticles(dir string) ([]models.Article, error)
```

**参数:**
- `dir`: 要扫描的目录路径

**返回:**
- `[]models.Article`: 文章列表
- `error`: 错误信息

**示例:**
```go
articles, err := scanner.ScanArticles("./content/post")
if err != nil {
    log.Fatal(err)
}
```

### Config 模块

#### LoadConfig
加载配置文件。

```go
func LoadConfig() (*Config, error)
```

**返回:**
- `*Config`: 配置对象
- `error`: 错误信息

### Translator 模块

#### LLMTranslator
基于LM Studio的AI翻译器。

```go
type LLMTranslator struct {
    client  *http.Client
    baseURL string
    model   string
    cache   *TranslationCache
}
```

#### NewLLMTranslator
创建新的翻译器实例。

```go
func NewLLMTranslator() *LLMTranslator
```

#### TranslateToSlug
将中文标签翻译为英文slug。

```go
func (t *LLMTranslator) TranslateToSlug(tag string) (string, error)
```

#### BatchTranslateTags
批量翻译标签。

```go
func (t *LLMTranslator) BatchTranslateTags(tags []string) (map[string]string, error)
```

#### BatchTranslateArticles
批量翻译文章标题。

```go
func (t *LLMTranslator) BatchTranslateArticles(titles []string) (map[string]string, error)
```

### Generator 模块

#### TagPageGenerator
标签页面生成器。

```go
type TagPageGenerator struct {
    contentDir string
    translator *translator.LLMTranslator
    slugCache  map[string]string
}
```

#### GenerateTagPagesWithMode
根据模式生成标签页面。

```go
func (g *TagPageGenerator) GenerateTagPagesWithMode(tagStats []models.TagStats, mode string) error
```

**参数:**
- `tagStats`: 标签统计列表
- `mode`: 处理模式 ("create", "update", "all")

#### ArticleSlugGenerator
文章slug生成器。

```go
type ArticleSlugGenerator struct {
    contentDir string
    translator *translator.LLMTranslator
}
```

#### GenerateArticleSlugsWithMode
根据模式生成文章Slug。

```go
func (g *ArticleSlugGenerator) GenerateArticleSlugsWithMode(mode string) error
```

#### ArticleTranslator
文章翻译器。

```go
type ArticleTranslator struct {
    contentDir string
    translator *translator.LLMTranslator
}
```

#### TranslateArticles
翻译文章。

```go
func (t *ArticleTranslator) TranslateArticles(mode string) error
```

#### PreviewArticleTranslations
预览文章翻译。

```go
func (t *ArticleTranslator) PreviewArticleTranslations() ([]ArticleTranslationPreview, error)
```

### Operations 模块

#### Processor
业务处理器。

```go
type Processor struct {
    contentDir string
}
```

#### NewProcessor
创建处理器实例。

```go
func NewProcessor(contentDir string) *Processor
```

#### QuickProcessAll
一键处理全部功能。

```go
func (p *Processor) QuickProcessAll(tagStats []models.TagStats, reader *bufio.Reader)
```

#### PreviewTagPages
预览标签页面。

```go
func (p *Processor) PreviewTagPages(tagStats []models.TagStats)
```

#### GenerateTagPages
生成标签页面。

```go
func (p *Processor) GenerateTagPages(tagStats []models.TagStats, reader *bufio.Reader)
```

#### PreviewArticleSlugs
预览文章Slug。

```go
func (p *Processor) PreviewArticleSlugs()
```

#### GenerateArticleSlugs
生成文章Slug。

```go
func (p *Processor) GenerateArticleSlugs(reader *bufio.Reader)
```

#### PreviewArticleTranslations
预览文章翻译。

```go
func (p *Processor) PreviewArticleTranslations()
```

#### TranslateArticles
翻译文章。

```go
func (p *Processor) TranslateArticles(reader *bufio.Reader)
```

#### ShowCacheStatus
显示缓存状态。

```go
func (p *Processor) ShowCacheStatus()
```

#### ShowBulkTranslationPreview
显示批量翻译预览。

```go
func (p *Processor) ShowBulkTranslationPreview(tagStats []models.TagStats)
```

#### GenerateBulkTranslationCache
生成批量翻译缓存。

```go
func (p *Processor) GenerateBulkTranslationCache(tagStats []models.TagStats, reader *bufio.Reader)
```

#### ClearTranslationCache
清空翻译缓存。

```go
func (p *Processor) ClearTranslationCache(reader *bufio.Reader)
```

### Menu 模块

#### InteractiveMenu
交互菜单。

```go
type InteractiveMenu struct {
    reader    *bufio.Reader
    processor *operations.Processor
}
```

#### NewInteractiveMenu
创建交互菜单。

```go
func NewInteractiveMenu(reader *bufio.Reader, contentDir string) *InteractiveMenu
```

#### Show
显示精简主菜单 (9个核心功能)。

```go
func (m *InteractiveMenu) Show(tagStats []models.TagStats, categoryStats []models.CategoryStats, noTagArticles []models.Article)
```

### Utils 模块

#### 性能监控函数

```go
// 记录翻译操作
func RecordTranslation(duration time.Duration)

// 记录缓存命中
func RecordCacheHit()

// 记录缓存失效
func RecordCacheMiss()

// 记录文件操作
func RecordFileOperation()

// 记录错误
func RecordError()

// 获取全局统计
func GetGlobalStats() PerformanceStats

// 重置全局统计
func ResetGlobalStats()
```

#### 日志函数

```go
// 日志记录
func Info(format string, args ...interface{})
func Debug(format string, args ...interface{})
func Warn(format string, args ...interface{})
func Error(format string, args ...interface{})

// 初始化日志
func InitLogger(filename string, level LogLevel) error
```

### Stats 模块

#### CalculateTagStats
计算标签统计信息。

```go
func CalculateTagStats(articles []models.Article) []models.TagStats
```

#### CalculateCategoryStats
计算分类统计信息。

```go
func CalculateCategoryStats(articles []models.Article) []models.CategoryStats
```

#### FindNoTagArticles
查找无标签的文章。

```go
func FindNoTagArticles(articles []models.Article) []models.Article
```

#### GroupTagsByFrequency
按频率分组标签。

```go
func GroupTagsByFrequency(tagStats []models.TagStats) ([]models.TagStats, []models.TagStats, []models.TagStats)
```

## 数据模型

### Article
文章模型。

```go
type Article struct {
    FilePath string   // 文件路径
    Title    string   // 文章标题
    Tags     []string // 标签列表
    Category string   // 分类
    Date     string   // 发布日期
}
```

### TagStats
标签统计模型。

```go
type TagStats struct {
    Name  string   // 标签名称
    Count int      // 使用次数
    Files []string // 使用该标签的文件列表
}
```

### CategoryStats
分类统计模型。

```go
type CategoryStats struct {
    Name  string // 分类名称
    Count int    // 文章数量
}
```

### PerformanceStats
性能统计模型。

```go
type PerformanceStats struct {
    TranslationCount int           // 翻译次数
    TranslationTime  time.Duration // 总翻译时间
    CacheHits        int           // 缓存命中次数
    CacheMisses      int           // 缓存失效次数
    FileOperations   int           // 文件操作次数
    Errors           int           // 错误次数
}
```

### Config
配置模型。

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

## 预览模型

### TagPagePreview
标签页面预览。

```go
type TagPagePreview struct {
    TagName     string
    Slug        string
    Status      string // "create", "update", "skip"
    Description string
}
```

### ArticleSlugPreview
文章Slug预览。

```go
type ArticleSlugPreview struct {
    FilePath string
    Title    string
    Slug     string
    Status   string // "missing", "exists"
}
```

### ArticleTranslationPreview
文章翻译预览。

```go
type ArticleTranslationPreview struct {
    SourceFile string
    TargetFile string
    Title      string
    Status     string // "missing", "exists"
}
```

### BulkTranslationPreview
批量翻译预览。

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

## 常用常量

```go
const (
    LMStudioURL     = "http://172.19.192.1:2234/v1/chat/completions"
    ModelName       = "gemma-3-12b-it"
    CacheFileName   = "tag_translations_cache.json"
)

// 日志级别
const (
    DEBUG LogLevel = iota
    INFO
    WARN
    ERROR
)
```

## 扩展开发

### 添加新的翻译器

```go
type CustomTranslator struct {
    // 自定义字段
}

func (t *CustomTranslator) TranslateToSlug(tag string) (string, error) {
    // 实现翻译逻辑
}
```

### 添加新的生成器

```go
type CustomGenerator struct {
    // 自定义字段
}

func (g *CustomGenerator) Generate(data interface{}) error {
    // 实现生成逻辑
}
```

### 添加新的菜单功能

```go
func (m *InteractiveMenu) customFunction() {
    // 实现自定义功能
}
```

在 `Show` 方法中添加新的选项和对应的处理逻辑。
