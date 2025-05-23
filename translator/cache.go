package translator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const CacheFileName = "tag_translations_cache.json"

type TranslationCache struct {
	Version      string                `json:"version"`
	LastUpdated  time.Time             `json:"last_updated"`
	Translations map[string]CacheEntry `json:"translations"`
	filePath     string
}

type CacheEntry struct {
	Translation string    `json:"translation"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewTranslationCache 创建新的翻译缓存
func NewTranslationCache(cacheDir string) *TranslationCache {
	if cacheDir == "" {
		cacheDir = "."
	}

	filePath := filepath.Join(cacheDir, CacheFileName)

	cache := &TranslationCache{
		Version:      "1.0",
		LastUpdated:  time.Now(),
		Translations: make(map[string]CacheEntry),
		filePath:     filePath,
	}

	// 尝试加载现有缓存
	cache.Load()

	return cache
}

// Load 从文件加载缓存
func (c *TranslationCache) Load() error {
	if _, err := os.Stat(c.filePath); os.IsNotExist(err) {
		fmt.Println("📄 缓存文件不存在，将创建新的缓存")
		return nil
	}

	data, err := os.ReadFile(c.filePath)
	if err != nil {
		return fmt.Errorf("读取缓存文件失败: %v", err)
	}

	if err := json.Unmarshal(data, c); err != nil {
		return fmt.Errorf("解析缓存文件失败: %v", err)
	}

	fmt.Printf("📄 已加载缓存文件，包含 %d 个翻译记录\n", len(c.Translations))
	return nil
}

// Save 保存缓存到文件
func (c *TranslationCache) Save() error {
	c.LastUpdated = time.Now()

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化缓存失败: %v", err)
	}

	if err := os.WriteFile(c.filePath, data, 0644); err != nil {
		return fmt.Errorf("写入缓存文件失败: %v", err)
	}

	fmt.Printf("💾 已保存缓存文件，包含 %d 个翻译记录\n", len(c.Translations))
	return nil
}

// Get 获取缓存的翻译
func (c *TranslationCache) Get(tag string) (string, bool) {
	entry, exists := c.Translations[tag]
	if !exists {
		return "", false
	}
	return entry.Translation, true
}

// Set 设置缓存的翻译
func (c *TranslationCache) Set(tag, translation string) {
	now := time.Now()

	if entry, exists := c.Translations[tag]; exists {
		// 更新现有条目
		entry.Translation = translation
		entry.UpdatedAt = now
		c.Translations[tag] = entry
	} else {
		// 创建新条目
		c.Translations[tag] = CacheEntry{
			Translation: translation,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
	}
}

// GetMissingTags 获取需要翻译的标签（缓存中不存在的）
func (c *TranslationCache) GetMissingTags(tags []string) []string {
	var missing []string

	for _, tag := range tags {
		if _, exists := c.Translations[tag]; !exists {
			missing = append(missing, tag)
		}
	}

	return missing
}

// GetStats 获取缓存统计信息
func (c *TranslationCache) GetStats() (int, int) {
	return len(c.Translations), 0 // 总数，过期数（暂时未实现过期机制）
}

// Clear 清空缓存
func (c *TranslationCache) Clear() {
	c.Translations = make(map[string]CacheEntry)
	c.LastUpdated = time.Now()
}

// GetCacheInfo 获取缓存文件信息
func (c *TranslationCache) GetCacheInfo() string {
	info := fmt.Sprintf("缓存文件: %s\n", c.filePath)
	info += fmt.Sprintf("版本: %s\n", c.Version)
	info += fmt.Sprintf("最后更新: %s\n", c.LastUpdated.Format("2006-01-02 15:04:05"))
	info += fmt.Sprintf("翻译条目: %d 个", len(c.Translations))
	return info
}
