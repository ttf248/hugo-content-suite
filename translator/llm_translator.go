package translator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hugo-content-suite/config"
	"hugo-content-suite/utils"
	"io"
	"net/http"
	"strings"
	"time"
)

// 提示词模板常量
const (
	slugPromptTemplate = `请将以下中文标签翻译为适合作为URL的英文slug。要求：
1. 使用小写字母
2. 单词之间用连字符(-)连接
3. 不包含特殊字符
4. 简洁准确
5. 只返回翻译结果，不要任何解释

中文标签: %s

英文slug:`

	categorySlugPromptTemplate = `请将以下中文分类名称翻译为适合作为URL的英文slug。要求：
1. 使用小写字母
2. 单词之间用连字符(-)连接
3. 不包含特殊字符
4. 简洁准确
5. 只返回翻译结果，不要任何解释

中文分类: %s

英文slug:`
)

// LM Studio API 相关类型定义
type LMStudioRequest struct {
	Model            string    `json:"model"`
	Messages         []Message `json:"messages"`
	Stream           bool      `json:"stream"`
	Temperature      float64   `json:"temperature,omitempty"`
	TopP             float64   `json:"top_p,omitempty"`
	MaxTokens        int       `json:"max_tokens,omitempty"`
	Stop             []string  `json:"stop,omitempty"`
	PresencePenalty  float64   `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64   `json:"frequency_penalty,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type LMStudioResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index   int     `json:"index"`
	Message Message `json:"message"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// CacheType、TagCache和ArticleCache已在cache.go中定义

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

// makeRequest 统一的HTTP请求方法
func (t *LLMTranslator) makeRequest(prompt string, timeout time.Duration) (string, error) {
	request := LMStudioRequest{
		Model: t.model,
		Messages: []Message{
			{Role: "user", Content: prompt},
		},
		Stream: false,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %v", err)
	}

	t.client.Timeout = timeout
	resp, err := t.client.Post(t.baseURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("LM Studio返回错误状态: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	var response LMStudioResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("没有获取到翻译结果")
	}

	return strings.TrimSpace(response.Choices[0].Message.Content), nil
}

// TestConnection 测试与LM Studio的连接
func (t *LLMTranslator) TestConnection() error {
	_, err := t.makeRequest("这是一个测试请求，无需处理，直接应答就行", 30*time.Second)
	return err
}

// translateWithCache 通用的带缓存翻译方法
func (t *LLMTranslator) translateWithCache(text string, cacheType CacheType, promptTemplate string) (string, error) {
	// 检查缓存 - 使用正确的缓存键格式
	cacheKey := fmt.Sprintf("en:%s", text)
	if cached, exists := t.cache.Get(cacheKey, cacheType); exists {
		return cached, nil
	}

	// 如果已经是英文，直接处理
	if utils.IsEnglishOnly(text) {
		slug := utils.NormalizeSlug(text)
		t.cache.Set(cacheKey, slug, cacheType)
		return slug, nil
	}

	// 构建提示词
	prompt := fmt.Sprintf(promptTemplate, text)

	// 发送请求
	result, err := t.makeRequest(prompt, t.timeout)
	if err != nil {
		return "", err
	}

	normalizedResult := utils.NormalizeSlug(result)
	t.cache.Set(cacheKey, normalizedResult, cacheType)

	return normalizedResult, nil
}

func (t *LLMTranslator) TranslateToSlug(text string) (string, error) {
	return t.translateWithCache(text, TagCache, slugPromptTemplate)
}

// 新增：翻译分类为slug
func (t *LLMTranslator) TranslateToCategorySlug(category string) (string, error) {
	return t.translateWithCache(category, CategoryCache, categorySlugPromptTemplate)
}

// 简化的缓存相关方法
func (t *LLMTranslator) GetMissingTags(tags []string) []string {
	return t.cache.GetMissingTexts(tags, "en", TagCache)
}

func (t *LLMTranslator) GetMissingArticles(articles []string) []string {
	return t.cache.GetMissingTexts(articles, "en", SlugCache)
}

// 新增：获取缺失的分类
func (t *LLMTranslator) GetMissingCategories(categories []string) []string {
	return t.cache.GetMissingTexts(categories, "en", CategoryCache)
}

func (t *LLMTranslator) PrepareBulkTranslation(allTexts []string) ([]string, int) {
	tags, articles := t.categorizeTexts(allTexts)
	missingTags := t.GetMissingTags(tags)
	missingArticles := t.GetMissingArticles(articles)

	allMissing := append(missingTags, missingArticles...)
	cachedCount := len(allTexts) - len(allMissing)

	return allMissing, cachedCount
}

// categorizeTexts 分类文本为标签和文章
func (t *LLMTranslator) categorizeTexts(allTexts []string) ([]string, []string) {
	var tags, articles []string

	for _, text := range allTexts {
		if len(text) <= 20 && !strings.Contains(text, "：") && !strings.Contains(text, ":") {
			tags = append(tags, text)
		} else {
			articles = append(articles, text)
		}
	}

	return tags, articles
}

// 缓存管理方法
func (t *LLMTranslator) GetCacheInfo() string      { return t.cache.GetInfo() }
func (t *LLMTranslator) ClearCache() error         { return t.cache.ClearAll() }
func (t *LLMTranslator) ClearTagCache() error      { return t.cache.Clear(TagCache) }
func (t *LLMTranslator) ClearArticleCache() error  { return t.cache.Clear(SlugCache) }
func (t *LLMTranslator) ClearCategoryCache() error { return t.cache.Clear(CategoryCache) } // 新增

func (t *LLMTranslator) GetCacheStats() int {
	tagTotal := t.cache.GetStats(TagCache)
	articleTotal := t.cache.GetStats(SlugCache)
	categoryTotal := t.cache.GetStats(CategoryCache)
	return tagTotal + articleTotal + categoryTotal
}
