package translator

import (
	"hugo-content-suite/config"
	"net/http"
	"time"
)

type LLMTranslator struct {
	baseURL string
	model   string
	timeout time.Duration
	cache   *TranslationCache
	client  *http.Client
}

func NewLLMTranslator() *LLMTranslator {
	cfg := config.GetGlobalConfig()
	translator := &LLMTranslator{
		baseURL: cfg.LMStudio.URL,
		model:   cfg.LMStudio.Model,
		timeout: time.Duration(cfg.LMStudio.Timeout) * time.Second,
		cache:   NewTranslationCache(),
		client:  &http.Client{Timeout: time.Duration(cfg.LMStudio.Timeout) * time.Second},
	}

	translator.cache.Load()
	return translator
}

// 缓存管理方法
func (t *LLMTranslator) GetCacheInfo() string      { return t.cache.GetInfo() }
func (t *LLMTranslator) ClearCache() error         { return t.cache.ClearAll() }
func (t *LLMTranslator) ClearTagCache() error      { return t.cache.Clear(kTagCache) }
func (t *LLMTranslator) ClearArticleCache() error  { return t.cache.Clear(kSlugCache) }
func (t *LLMTranslator) ClearCategoryCache() error { return t.cache.Clear(kCategoryCache) } // 新增

func (t *LLMTranslator) GetCacheStats() int {
	tagTotal := t.cache.GetStats(kTagCache)
	articleTotal := t.cache.GetStats(kSlugCache)
	categoryTotal := t.cache.GetStats(kCategoryCache)
	return tagTotal + articleTotal + categoryTotal
}
