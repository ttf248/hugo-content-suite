# API 接口文档

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

**参数:**
- `tag`: 要翻译的中文标签

**返回:**
- `string`: 翻译后的英文slug
- `error`: 错误信息

#### BatchTranslate
批量翻译标签（支持缓存）。

```go
func (t *LLMTranslator) BatchTranslate(tags []string) (map[string]string, error)
```

**参数:**
- `tags`: 要翻译的标签列表

**返回:**
- `map[string]string`: 标签到slug的映射
- `error`: 错误信息

#### TestConnection
测试与LM Studio的连接。

```go
func (t *LLMTranslator) TestConnection() error
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

#### NewTagPageGenerator
创建标签页面生成器。

```go
func NewTagPageGenerator(contentDir string) *TagPageGenerator
```

#### GenerateTagPages
生成所有标签页面。

```go
func (g *TagPageGenerator) GenerateTagPages(tagStats []models.TagStats) error
```

#### PreviewTagPages
预览标签页面生成。

```go
func (g *TagPageGenerator) PreviewTagPages(tagStats []models.TagStats) []TagPagePreview
```

#### ArticleSlugGenerator
文章slug生成器。

```go
type ArticleSlugGenerator struct {
    contentDir string
    translator *translator.LLMTranslator
}
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

## LM Studio API

### 请求格式

#### Chat Completions
```json
{
  "model": "gemma-3-12b-it",
  "messages": [
    {
      "role": "user",
      "content": "翻译请求内容"
    }
  ],
  "stream": false
}
```

#### 响应格式
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
        "content": "翻译结果"
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

## 缓存接口

### TranslationCache
翻译缓存管理。

```go
type TranslationCache struct {
    Version      string                `json:"version"`
    LastUpdated  time.Time             `json:"last_updated"`
    Translations map[string]CacheEntry `json:"translations"`
    filePath     string
}
```

#### 主要方法

```go
// 创建缓存
func NewTranslationCache(cacheDir string) *TranslationCache

// 获取翻译
func (c *TranslationCache) Get(tag string) (string, bool)

// 设置翻译
func (c *TranslationCache) Set(tag, translation string)

// 保存缓存
func (c *TranslationCache) Save() error

// 加载缓存
func (c *TranslationCache) Load() error
```

## 错误处理

### 常见错误类型

```go
// 网络连接错误
fmt.Errorf("发送请求失败: %v", err)

// 文件读写错误
fmt.Errorf("读取文件失败: %v", err)

// JSON解析错误
fmt.Errorf("解析响应失败: %v", err)

// LM Studio错误
fmt.Errorf("LM Studio返回错误状态: %d", resp.StatusCode)
```

### 错误处理策略

1. **网络错误**: 自动重试机制
2. **翻译失败**: 回退到预定义映射
3. **文件错误**: 详细错误提示
4. **缓存错误**: 继续执行但给出警告

## 配置参数

### 常量配置

```go
const (
    LMStudioURL     = "http://172.19.192.1:2234/v1/chat/completions"
    ModelName       = "gemma-3-12b-it"
    CacheFileName   = "tag_translations_cache.json"
)
```

### 可调参数

```go
// HTTP超时设置
Timeout: 30 * time.Second

// 请求间隔
time.Sleep(500 * time.Millisecond)

// 表格显示限制
defaultLimit := 20
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
