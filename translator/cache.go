package translator

import (
	"encoding/json"
	"fmt"
	"hugo-content-suite/config"
	"hugo-content-suite/utils"
	"os"
	"path/filepath"
	"time"
)

type CacheType string

const (
	kTagCache      CacheType = "tag"
	kSlugCache     CacheType = "article"
	kCategoryCache CacheType = "category" // 新增
)

type CacheEntry struct {
	Translation string    `json:"translation"`
	Timestamp   time.Time `json:"timestamp"`
	Type        CacheType `json:"type"`
}

type TranslationCache struct {
	tagCacheFile      string
	slugCacheFile     string
	categoryCacheFile string // 新增
	tagCache          map[string]CacheEntry
	slugCache         map[string]CacheEntry
	categoryCache     map[string]CacheEntry // 新增
}

func NewTranslationCache() *TranslationCache {
	return NewTranslationCacheWithConfig(config.GetGlobalConfig())
}

func NewTranslationCacheWithConfig(cfg *config.Config) *TranslationCache {
	return &TranslationCache{
		tagCacheFile:      cfg.Cache.TagFileName,
		slugCacheFile:     cfg.Cache.ArticleFileName,
		categoryCacheFile: cfg.Cache.CategoryFileName, // 新增
		tagCache:          make(map[string]CacheEntry),
		slugCache:         make(map[string]CacheEntry),
		categoryCache:     make(map[string]CacheEntry), // 新增
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

	// 加载分类缓存
	if err := c.loadCacheFile(c.categoryCacheFile, &c.categoryCache); err != nil {
		utils.WarnWithFields("加载分类缓存失败", map[string]interface{}{
			"file":  c.categoryCacheFile,
			"error": err.Error(),
		})
		c.categoryCache = make(map[string]CacheEntry)
	}

	utils.InfoWithFields("缓存加载完成", map[string]interface{}{
		"tag_count":      len(c.tagCache),
		"article_count":  len(c.slugCache),
		"category_count": len(c.categoryCache),
	})

	fmt.Printf("📄 已加载缓存文件 - 标签: %d 个, Slug: %d 个, 分类: %d 个\n",
		len(c.tagCache), len(c.slugCache), len(c.categoryCache))
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

	// 保存分类缓存
	if err := c.saveCacheFile(c.categoryCacheFile, c.categoryCache); err != nil {
		return fmt.Errorf("保存分类缓存失败: %v", err)
	}

	fmt.Printf("💾 已保存缓存文件 - 标签: %d 个, 文章: %d 个, 分类: %d 个\n",
		len(c.tagCache), len(c.slugCache), len(c.categoryCache))
	return nil
}

func (c *TranslationCache) saveCacheFile(filename string, cache map[string]CacheEntry) error {
	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(filename), 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(filename), ".cache-*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	// Rename keeps a malformed/partial write from replacing an existing cache.
	if err := os.Rename(tmpName, filename); err != nil {
		if removeErr := os.Remove(filename); removeErr != nil && !os.IsNotExist(removeErr) {
			return err
		}
		return os.Rename(tmpName, filename)
	}
	return nil
}

func (c *TranslationCache) Get(text string, cacheType CacheType) (string, bool) {
	var cache map[string]CacheEntry
	switch cacheType {
	case kTagCache:
		cache = c.tagCache
	case kSlugCache:
		cache = c.slugCache
	case kCategoryCache:
		cache = c.categoryCache
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
	case kTagCache:
		c.tagCache[text] = entry
	case kSlugCache:
		c.slugCache[text] = entry
	case kCategoryCache:
		c.categoryCache[text] = entry
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
	case kTagCache:
		cache = c.tagCache
	case kSlugCache:
		cache = c.slugCache
	case kCategoryCache:
		cache = c.categoryCache
	default:
		return 0
	}

	return len(cache)
}

func (c *TranslationCache) Clear(cacheType CacheType) error {
	switch cacheType {
	case kTagCache:
		c.tagCache = make(map[string]CacheEntry)
		return c.saveCacheFile(c.tagCacheFile, c.tagCache)
	case kSlugCache:
		c.slugCache = make(map[string]CacheEntry)
		return c.saveCacheFile(c.slugCacheFile, c.slugCache)
	case kCategoryCache:
		c.categoryCache = make(map[string]CacheEntry)
		return c.saveCacheFile(c.categoryCacheFile, c.categoryCache)
	default:
		return fmt.Errorf("未知的缓存类型: %v", cacheType)
	}
}

func (c *TranslationCache) ClearAll() error {
	c.tagCache = make(map[string]CacheEntry)
	c.slugCache = make(map[string]CacheEntry)
	c.categoryCache = make(map[string]CacheEntry)

	if err := c.saveCacheFile(c.tagCacheFile, c.tagCache); err != nil {
		return err
	}
	if err := c.saveCacheFile(c.slugCacheFile, c.slugCache); err != nil {
		return err
	}
	if err := c.saveCacheFile(c.categoryCacheFile, c.categoryCache); err != nil {
		return err
	}
	return nil
}

func (c *TranslationCache) GetInfo() string {
	tagTotal := c.GetStats(kTagCache)
	slugTotal := c.GetStats(kSlugCache)
	categoryTotal := c.GetStats(kCategoryCache)

	return fmt.Sprintf(`📊 缓存状态信息:
🏷️  标签缓存:
   📁 文件: %s
   📄 总条目: %d 个

📝 Slug缓存:
   📁 文件: %s
   📄 总条目: %d 个

📂 分类缓存:
   📁 文件: %s
   📄 总条目: %d 个`,
		c.tagCacheFile, tagTotal,
		c.slugCacheFile, slugTotal,
		c.categoryCacheFile, categoryTotal)
}
