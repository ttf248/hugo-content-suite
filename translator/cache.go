package translator

import (
	"encoding/json"
	"fmt"
	"os"
	"tag-scanner/config"
	"time"
)

type CacheType string

const (
	TagCache     CacheType = "tag"
	ArticleCache CacheType = "article"
)

type CacheEntry struct {
	Translation string    `json:"translation"`
	Timestamp   time.Time `json:"timestamp"`
	Type        CacheType `json:"type"`
}

type TranslationCache struct {
	tagCacheFile     string
	articleCacheFile string
	tagCache         map[string]CacheEntry
	articleCache     map[string]CacheEntry
	expireDuration   time.Duration
}

func NewTranslationCache() *TranslationCache {
	cfg := config.GetGlobalConfig()
	return &TranslationCache{
		tagCacheFile:     cfg.Cache.TagFileName,
		articleCacheFile: cfg.Cache.ArticleFileName,
		tagCache:         make(map[string]CacheEntry),
		articleCache:     make(map[string]CacheEntry),
		expireDuration:   time.Duration(cfg.Cache.ExpireDays) * 24 * time.Hour,
	}
}

func (c *TranslationCache) Load() error {
	// 加载标签缓存
	if err := c.loadCacheFile(c.tagCacheFile, &c.tagCache); err != nil {
		fmt.Printf("⚠️ 加载标签缓存失败: %v\n", err)
		c.tagCache = make(map[string]CacheEntry)
	}

	// 加载文章缓存
	if err := c.loadCacheFile(c.articleCacheFile, &c.articleCache); err != nil {
		fmt.Printf("⚠️ 加载文章缓存失败: %v\n", err)
		c.articleCache = make(map[string]CacheEntry)
	}

	fmt.Printf("📄 已加载缓存文件 - 标签: %d 个, 文章: %d 个\n",
		len(c.tagCache), len(c.articleCache))
	return nil
}

func (c *TranslationCache) loadCacheFile(filename string, cache *map[string]CacheEntry) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil // 文件不存在，不是错误
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return nil // 空文件
	}

	return json.Unmarshal(data, cache)
}

func (c *TranslationCache) Save() error {
	// 保存标签缓存
	if err := c.saveCacheFile(c.tagCacheFile, c.tagCache); err != nil {
		return fmt.Errorf("保存标签缓存失败: %v", err)
	}

	// 保存文章缓存
	if err := c.saveCacheFile(c.articleCacheFile, c.articleCache); err != nil {
		return fmt.Errorf("保存文章缓存失败: %v", err)
	}

	fmt.Printf("💾 已保存缓存文件 - 标签: %d 个, 文章: %d 个\n",
		len(c.tagCache), len(c.articleCache))
	return nil
}

func (c *TranslationCache) saveCacheFile(filename string, cache map[string]CacheEntry) error {
	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func (c *TranslationCache) Get(text string, cacheType CacheType) (string, bool) {
	var cache map[string]CacheEntry
	switch cacheType {
	case TagCache:
		cache = c.tagCache
	case ArticleCache:
		cache = c.articleCache
	default:
		return "", false
	}

	entry, exists := cache[text]
	if !exists {
		return "", false
	}

	// 检查是否过期
	if time.Since(entry.Timestamp) > c.expireDuration {
		delete(cache, text)
		return "", false
	}

	return entry.Translation, true
}

func (c *TranslationCache) Set(text, translation string, cacheType CacheType) {
	entry := CacheEntry{
		Translation: translation,
		Timestamp:   time.Now(),
		Type:        cacheType,
	}

	switch cacheType {
	case TagCache:
		c.tagCache[text] = entry
	case ArticleCache:
		c.articleCache[text] = entry
	}
}

func (c *TranslationCache) GetMissingTexts(texts []string, cacheType CacheType) []string {
	var missing []string
	for _, text := range texts {
		if _, exists := c.Get(text, cacheType); !exists {
			missing = append(missing, text)
		}
	}
	return missing
}

func (c *TranslationCache) GetStats(cacheType CacheType) (total int, expired int) {
	var cache map[string]CacheEntry
	switch cacheType {
	case TagCache:
		cache = c.tagCache
	case ArticleCache:
		cache = c.articleCache
	default:
		return 0, 0
	}

	total = len(cache)
	for _, entry := range cache {
		if time.Since(entry.Timestamp) > c.expireDuration {
			expired++
		}
	}
	return
}

func (c *TranslationCache) Clear(cacheType CacheType) error {
	switch cacheType {
	case TagCache:
		c.tagCache = make(map[string]CacheEntry)
		return c.saveCacheFile(c.tagCacheFile, c.tagCache)
	case ArticleCache:
		c.articleCache = make(map[string]CacheEntry)
		return c.saveCacheFile(c.articleCacheFile, c.articleCache)
	default:
		return fmt.Errorf("未知的缓存类型: %v", cacheType)
	}
}

func (c *TranslationCache) ClearAll() error {
	c.tagCache = make(map[string]CacheEntry)
	c.articleCache = make(map[string]CacheEntry)

	if err := c.saveCacheFile(c.tagCacheFile, c.tagCache); err != nil {
		return err
	}
	if err := c.saveCacheFile(c.articleCacheFile, c.articleCache); err != nil {
		return err
	}
	return nil
}

func (c *TranslationCache) GetInfo() string {
	tagTotal, tagExpired := c.GetStats(TagCache)
	articleTotal, articleExpired := c.GetStats(ArticleCache)

	return fmt.Sprintf(`📊 缓存状态信息:
🏷️  标签缓存:
   📁 文件: %s
   📄 总条目: %d 个
   ⏰ 过期条目: %d 个
   ✅ 有效条目: %d 个

📝 文章缓存:
   📁 文件: %s
   📄 总条目: %d 个
   ⏰ 过期条目: %d 个
   ✅ 有效条目: %d 个`,
		c.tagCacheFile, tagTotal, tagExpired, tagTotal-tagExpired,
		c.articleCacheFile, articleTotal, articleExpired, articleTotal-articleExpired)
}
