package translator

import (
	"encoding/json"
	"fmt"
	"hugo-content-suite/config"
	"hugo-content-suite/utils"
	"os"
	"time"
)

type CacheType string

const (
	TagCache  CacheType = "tag"
	SlugCache CacheType = "article"
)

type CacheEntry struct {
	Translation string    `json:"translation"`
	Timestamp   time.Time `json:"timestamp"`
	Type        CacheType `json:"type"`
}

type TranslationCache struct {
	tagCacheFile  string
	slugCacheFile string
	tagCache      map[string]CacheEntry
	slugCache     map[string]CacheEntry
}

func NewTranslationCache() *TranslationCache {
	cfg := config.GetGlobalConfig()
	return &TranslationCache{
		tagCacheFile:  cfg.Cache.TagFileName,
		slugCacheFile: cfg.Cache.ArticleFileName,
		tagCache:      make(map[string]CacheEntry),
		slugCache:     make(map[string]CacheEntry),
	}
}

func (c *TranslationCache) Load() error {
	// 加载标签缓存
	if err := c.loadCacheFile(c.tagCacheFile, &c.tagCache); err != nil {
		utils.WarnWithFields("加载标签缓存失败", map[string]interface{}{
			"file":  c.tagCacheFile,
			"error": err.Error(),
		})
		c.tagCache = make(map[string]CacheEntry)
	}

	// 加载文章缓存
	if err := c.loadCacheFile(c.slugCacheFile, &c.slugCache); err != nil {
		utils.WarnWithFields("加载文章缓存失败", map[string]interface{}{
			"file":  c.slugCacheFile,
			"error": err.Error(),
		})
		c.slugCache = make(map[string]CacheEntry)
	}

	utils.InfoWithFields("缓存加载完成", map[string]interface{}{
		"tag_count":     len(c.tagCache),
		"article_count": len(c.slugCache),
	})

	fmt.Printf("📄 已加载缓存文件 - 标签: %d 个, 文章: %d 个\n",
		len(c.tagCache), len(c.slugCache))
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
	if err := c.saveCacheFile(c.slugCacheFile, c.slugCache); err != nil {
		return fmt.Errorf("保存文章缓存失败: %v", err)
	}

	fmt.Printf("💾 已保存缓存文件 - 标签: %d 个, 文章: %d 个\n",
		len(c.tagCache), len(c.slugCache))
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
	case SlugCache:
		cache = c.slugCache
	default:
		return "", false
	}

	entry, exists := cache[text]
	if !exists {
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
	case SlugCache:
		c.slugCache[text] = entry
	}
}

// GetMissingTexts 获取缓存中缺失的文本
func (c *TranslationCache) GetMissingTexts(texts []string, targetLang string, cacheType CacheType) []string {
	var missing []string
	for _, text := range texts {
		cacheKey := fmt.Sprintf("%s:%s", targetLang, text)
		if _, exists := c.Get(cacheKey, cacheType); !exists {
			missing = append(missing, text)
		}
	}
	return missing
}

// GetCachedTranslations 获取已缓存的翻译
func (c *TranslationCache) GetCachedTranslations(texts []string, targetLang string, cacheType CacheType) map[string]string {
	result := make(map[string]string)
	for _, text := range texts {
		cacheKey := fmt.Sprintf("%s:%s", targetLang, text)
		if translation, exists := c.Get(cacheKey, cacheType); exists {
			result[text] = translation
		}
	}
	return result
}

func (c *TranslationCache) GetStats(cacheType CacheType) (total int) {
	var cache map[string]CacheEntry
	switch cacheType {
	case TagCache:
		cache = c.tagCache
	case SlugCache:
		cache = c.slugCache
	default:
		return 0
	}

	return len(cache)
}

func (c *TranslationCache) Clear(cacheType CacheType) error {
	switch cacheType {
	case TagCache:
		c.tagCache = make(map[string]CacheEntry)
		return c.saveCacheFile(c.tagCacheFile, c.tagCache)
	case SlugCache:
		c.slugCache = make(map[string]CacheEntry)
		return c.saveCacheFile(c.slugCacheFile, c.slugCache)
	default:
		return fmt.Errorf("未知的缓存类型: %v", cacheType)
	}
}

func (c *TranslationCache) ClearAll() error {
	c.tagCache = make(map[string]CacheEntry)
	c.slugCache = make(map[string]CacheEntry)

	if err := c.saveCacheFile(c.tagCacheFile, c.tagCache); err != nil {
		return err
	}
	if err := c.saveCacheFile(c.slugCacheFile, c.slugCache); err != nil {
		return err
	}
	return nil
}

func (c *TranslationCache) GetInfo() string {
	tagTotal := c.GetStats(TagCache)
	slugTotal := c.GetStats(SlugCache)

	return fmt.Sprintf(`📊 缓存状态信息:
🏷️  标签缓存:
   📁 文件: %s
   📄 总条目: %d 个

📝 Slug缓存:
   📁 文件: %s
   📄 总条目: %d 个`,
		c.tagCacheFile, tagTotal,
		c.slugCacheFile, slugTotal)
}
