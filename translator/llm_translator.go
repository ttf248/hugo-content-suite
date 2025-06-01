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

	articleSlugPromptTemplate = `请将以下中文文章标题翻译为简洁的英文slug，要求：
1. 使用小写字母
2. 单词间用连字符(-)连接
3. 去除特殊字符
4. 保持语义准确
5. 适合作为URL路径

标题：%s

请只返回翻译后的slug，不要其他内容。`
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

func (t *LLMTranslator) TranslateToArticleSlug(title string) (string, error) {
	return t.translateWithCache(title, SlugCache, articleSlugPromptTemplate)
}

// BatchTranslateTags 批量翻译标签
func (t *LLMTranslator) BatchTranslateTags(tags []string) (map[string]string, error) {
	return t.batchTranslate(tags, TagCache, "标签", t.TranslateToSlug)
}

// BatchTranslateSlugs 批量翻译文章标题
func (t *LLMTranslator) BatchTranslateSlugs(titles []string) (map[string]string, error) {
	return t.batchTranslate(titles, SlugCache, "Slug", t.TranslateToArticleSlug)
}

// batchTranslate 通用批量翻译方法
func (t *LLMTranslator) batchTranslate(texts []string, cacheType CacheType, typeName string, translateFunc func(string) (string, error)) (map[string]string, error) {
	cfg := config.GetGlobalConfig()
	startTime := time.Now()
	result := make(map[string]string)

	// 处理缓存
	_, missingTexts := t.processCacheAndGetMissing(texts, cacheType, typeName, result)
	if len(missingTexts) == 0 {
		return result, nil
	}

	// 批量翻译
	progressBar := utils.NewProgressBar(len(missingTexts))
	newTranslationsAdded := t.translateMissingTexts(missingTexts, result, translateFunc, progressBar, cfg)

	// 保存缓存
	if newTranslationsAdded > 0 {
		t.saveCacheAndLog(startTime)
	}

	return result, nil
}

// processCacheAndGetMissing 处理缓存并获取缺失的文本
func (t *LLMTranslator) processCacheAndGetMissing(texts []string, cacheType CacheType, typeName string, result map[string]string) (int, []string) {
	fmt.Printf("🔍 检查%s缓存...\n", typeName)

	// 批量检查缓存
	cachedCount := t.loadFromCache(texts, cacheType, result)
	if cachedCount > 0 {
		fmt.Printf("📋 从缓存获取 %d 个%s翻译\n", cachedCount, typeName)
	}

	// 获取需要翻译的文本
	missingTexts := t.cache.GetMissingTexts(texts, "en", cacheType)
	if len(missingTexts) == 0 {
		fmt.Printf("✅ 所有%s都已有缓存，无需重新翻译\n", typeName)
		return cachedCount, nil
	}

	fmt.Printf("🔄 需要翻译 %d 个新%s\n", len(missingTexts), typeName)
	return cachedCount, missingTexts
}

// loadFromCache 从缓存加载已有翻译
func (t *LLMTranslator) loadFromCache(texts []string, cacheType CacheType, result map[string]string) int {
	cachedCount := 0
	for _, text := range texts {
		cacheKey := fmt.Sprintf("en:%s", text)
		if translation, exists := t.cache.Get(cacheKey, cacheType); exists {
			result[text] = translation
			cachedCount++
		}
	}
	return cachedCount
}

// translateMissingTexts 翻译缺失的文本
func (t *LLMTranslator) translateMissingTexts(missingTexts []string, result map[string]string, translateFunc func(string) (string, error), progressBar *utils.ProgressBar, cfg *config.Config) int {
	newTranslationsAdded := 0

	for i, text := range missingTexts {
		slug, err := translateFunc(text)
		if err != nil {
			// 翻译失败，跳过此项
			utils.Error("翻译失败: %s - %v", text, err)
			progressBar.Update(i + 1)
			continue
		}

		// 保存翻译结果
		newTranslationsAdded += t.saveTranslationResult(text, slug, result)
		progressBar.Update(i + 1)

		// 处理自动保存和延迟
		t.handleAutoSaveAndDelay(newTranslationsAdded, i, len(missingTexts), cfg)
	}

	return newTranslationsAdded
}

// saveTranslationResult 保存翻译结果到缓存和结果集
func (t *LLMTranslator) saveTranslationResult(text, translation string, result map[string]string) int {
	cacheKey := fmt.Sprintf("en:%s", text)
	cacheType := t.determineCacheType(text)
	t.cache.Set(cacheKey, translation, cacheType)
	result[text] = translation
	return 1
}

// determineCacheType 根据文本特征确定缓存类型
func (t *LLMTranslator) determineCacheType(text string) CacheType {
	if len(text) <= 20 && !strings.Contains(text, "：") && !strings.Contains(text, ":") {
		return TagCache
	}
	return SlugCache
}

// handleAutoSaveAndDelay 处理自动保存和延迟
func (t *LLMTranslator) handleAutoSaveAndDelay(translationsCount, currentIndex, totalCount int, cfg *config.Config) {
	// 自动保存
	if translationsCount%cfg.Cache.AutoSaveCount == 0 {
		if err := t.cache.Save(); err != nil {
			utils.Error("中间保存缓存失败: %v", err)
		}
	}

	// 添加延迟
	if currentIndex < totalCount-1 {
		time.Sleep(time.Duration(cfg.Cache.DelayMs) * time.Millisecond)
	}
}

// saveCacheAndLog 保存缓存并记录日志
func (t *LLMTranslator) saveCacheAndLog(startTime time.Time) {
	if err := t.cache.Save(); err != nil {
		utils.Error("保存缓存失败: %v", err)
	} else {
		utils.Info("批量翻译完成，耗时: %v", time.Since(startTime))
	}
}

// 简化的缓存相关方法
func (t *LLMTranslator) GetMissingTags(tags []string) []string {
	return t.cache.GetMissingTexts(tags, "en", TagCache)
}

func (t *LLMTranslator) GetMissingArticles(articles []string) []string {
	return t.cache.GetMissingTexts(articles, "en", SlugCache)
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
func (t *LLMTranslator) GetCacheInfo() string     { return t.cache.GetInfo() }
func (t *LLMTranslator) ClearCache() error        { return t.cache.ClearAll() }
func (t *LLMTranslator) ClearTagCache() error     { return t.cache.Clear(TagCache) }
func (t *LLMTranslator) ClearArticleCache() error { return t.cache.Clear(SlugCache) }

func (t *LLMTranslator) GetCacheStats() int {
	tagTotal := t.cache.GetStats(TagCache)
	articleTotal := t.cache.GetStats(SlugCache)
	return tagTotal + articleTotal
}
