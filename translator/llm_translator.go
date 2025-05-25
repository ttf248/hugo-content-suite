package translator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hugo-content-suite/config"
	"hugo-content-suite/utils"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// LM Studio API 相关类型定义
type LMStudioRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
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

	client := &http.Client{Timeout: timeout}
	resp, err := client.Post(t.baseURL, "application/json", bytes.NewBuffer(jsonData))
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
	_, err := t.makeRequest("test", 5*time.Second)
	return err
}

// translateWithCache 通用的带缓存翻译方法
func (t *LLMTranslator) translateWithCache(text string, cacheType CacheType, promptTemplate string) (string, error) {
	// 检查缓存
	if cached, exists := t.cache.Get(text, cacheType); exists {
		utils.RecordCacheHit()
		return cached, nil
	}
	utils.RecordCacheMiss()

	// 如果已经是英文，直接处理
	if isEnglishOnly(text) {
		slug := normalizeSlug(text)
		t.cache.Set(text, slug, cacheType)
		return slug, nil
	}

	// 构建提示词
	prompt := fmt.Sprintf(promptTemplate, text)

	// 发送请求
	result, err := t.makeRequest(prompt, t.timeout)
	if err != nil {
		return "", err
	}

	normalizedResult := normalizeSlug(result)
	t.cache.Set(text, normalizedResult, cacheType)

	return normalizedResult, nil
}

func (t *LLMTranslator) TranslateToSlug(text string) (string, error) {
	promptTemplate := `请将以下中文标签翻译为适合作为URL的英文slug。要求：
1. 使用小写字母
2. 单词之间用连字符(-)连接
3. 不包含特殊字符
4. 简洁准确
5. 只返回翻译结果，不要任何解释

中文标签: %s

英文slug:`

	return t.translateWithCache(text, TagCache, promptTemplate)
}

func (t *LLMTranslator) TranslateToArticleSlug(title string) (string, error) {
	promptTemplate := `请将以下中文文章标题翻译为简洁的英文slug，要求：
1. 使用小写字母
2. 单词间用连字符(-)连接
3. 去除特殊字符
4. 保持语义准确
5. 适合作为URL路径

标题：%s

请只返回翻译后的slug，不要其他内容。`

	return t.translateWithCache(title, ArticleCache, promptTemplate)
}

func (t *LLMTranslator) TranslateParagraph(paragraph string) (string, error) {
	if strings.TrimSpace(paragraph) == "" {
		return paragraph, nil
	}

	if t.shouldSkipTranslation(paragraph) {
		return paragraph, nil
	}

	promptTemplate := `请将以下中文段落翻译成自然流畅的英文，保持原文的格式和结构：

%s

要求：
1. 翻译要自然流畅，符合英文表达习惯
2. 保持原文的段落结构和格式
3. 如果包含技术术语，请使用准确的英文术语
4. 如果包含Markdown格式，请保留格式标记
5. 直接返回翻译结果，不要添加额外说明
6. 如果原文已经是英文，请保持不变`

	prompt := fmt.Sprintf(promptTemplate, paragraph)
	return t.makeRequest(prompt, t.timeout)
}

// shouldSkipTranslation 判断是否应该跳过翻译
func (t *LLMTranslator) shouldSkipTranslation(text string) bool {
	trimmed := strings.TrimSpace(text)

	// 空内容跳过
	if trimmed == "" {
		return true
	}

	// 检查是否为代码块
	if strings.HasPrefix(trimmed, "```") || strings.HasSuffix(trimmed, "```") {
		return true
	}

	// 检查是否为缩进代码块
	if strings.HasPrefix(trimmed, "    ") {
		return true
	}

	// 检查是否为引用块（但如果包含中文仍需翻译）
	if strings.HasPrefix(trimmed, ">") {
		// 检查引用内容是否包含中文
		if !t.containsChinese(trimmed) {
			return true
		}
	}

	// 检查是否为纯链接行（不包含中文描述）
	if strings.Contains(trimmed, "](") && strings.Contains(trimmed, "[") {
		// 如果链接中包含中文描述，仍需翻译
		if !t.containsChinese(trimmed) {
			return true
		}
	}

	// 检查是否为图片（但如果alt文本包含中文仍需翻译）
	if strings.HasPrefix(trimmed, "![") {
		if !t.containsChinese(trimmed) {
			return true
		}
	}

	// 检查是否为HTML标签
	if strings.HasPrefix(trimmed, "<") && strings.HasSuffix(trimmed, ">") {
		if !t.containsChinese(trimmed) {
			return true
		}
	}

	// 如果没有中文字符，跳过翻译
	if !t.containsChinese(trimmed) {
		return true
	}

	// 只要包含中文就翻译，不再检查中文字符比例
	return false
}

// containsChinese 检查文本是否包含中文
func (t *LLMTranslator) containsChinese(text string) bool {
	for _, r := range text {
		if r >= 0x4e00 && r <= 0x9fff {
			return true
		}
	}
	return false
}

// BatchTranslateTags 批量翻译标签
func (t *LLMTranslator) BatchTranslateTags(tags []string) (map[string]string, error) {
	return t.batchTranslate(tags, TagCache, "标签", t.TranslateToSlug)
}

// BatchTranslateArticles 批量翻译文章标题
func (t *LLMTranslator) BatchTranslateArticles(titles []string) (map[string]string, error) {
	return t.batchTranslate(titles, ArticleCache, "文章标题", t.TranslateToArticleSlug)
}

// BatchTranslate 兼容旧接口，自动判断类型
func (t *LLMTranslator) BatchTranslate(texts []string) (map[string]string, error) {
	// 默认当作标签处理，保持向后兼容
	return t.BatchTranslateTags(texts)
}

// batchTranslate 通用批量翻译方法
func (t *LLMTranslator) batchTranslate(texts []string, cacheType CacheType, typeName string, translateFunc func(string) (string, error)) (map[string]string, error) {
	cfg := config.GetGlobalConfig()
	startTime := time.Now()
	result := make(map[string]string)

	fmt.Printf("🔍 检查%s缓存...\n", typeName)

	// 批量检查缓存
	cachedCount := t.loadFromCache(texts, cacheType, result)
	if cachedCount > 0 {
		fmt.Printf("📋 从缓存获取 %d 个%s翻译\n", cachedCount, typeName)
	}

	// 获取需要翻译的文本
	missingTexts := t.cache.GetMissingTexts(texts, cacheType)
	if len(missingTexts) == 0 {
		fmt.Printf("✅ 所有%s都已有缓存，无需重新翻译\n", typeName)
		return result, nil
	}

	fmt.Printf("🔄 需要翻译 %d 个新%s\n", len(missingTexts), typeName)
	progressBar := utils.NewProgressBar(len(missingTexts))

	// 批量翻译
	newTranslationsAdded := t.translateMissingTexts(missingTexts, result, translateFunc, progressBar, cfg)

	// 保存缓存
	if newTranslationsAdded > 0 {
		t.saveCacheAndLog(startTime)
	}

	return result, nil
}

// loadFromCache 从缓存加载已有翻译
func (t *LLMTranslator) loadFromCache(texts []string, cacheType CacheType, result map[string]string) int {
	cachedCount := 0
	for _, text := range texts {
		if translation, exists := t.cache.Get(text, cacheType); exists {
			result[text] = translation
			cachedCount++
			utils.RecordCacheHit()
		} else {
			utils.RecordCacheMiss()
		}
	}
	return cachedCount
}

// translateMissingTexts 翻译缺失的文本
func (t *LLMTranslator) translateMissingTexts(missingTexts []string, result map[string]string, translateFunc func(string) (string, error), progressBar *utils.ProgressBar, cfg *config.Config) int {
	newTranslationsAdded := 0

	for i, text := range missingTexts {
		translationStart := time.Now()

		slug, err := translateFunc(text)
		if err != nil {
			utils.RecordError()
			slug = fallbackSlug(text)
		}

		utils.RecordTranslation(time.Since(translationStart))
		result[text] = slug
		newTranslationsAdded++

		progressBar.Update(i + 1)

		// 中间保存
		if newTranslationsAdded%cfg.Cache.AutoSaveCount == 0 {
			if err := t.cache.Save(); err != nil {
				utils.Error("中间保存缓存失败: %v", err)
			} else {
				utils.RecordFileOperation()
			}
		}

		// 添加延迟
		if i < len(missingTexts)-1 {
			time.Sleep(time.Duration(cfg.Cache.DelayMs) * time.Millisecond)
		}
	}

	return newTranslationsAdded
}

// saveCacheAndLog 保存缓存并记录日志
func (t *LLMTranslator) saveCacheAndLog(startTime time.Time) {
	if err := t.cache.Save(); err != nil {
		utils.Error("保存缓存失败: %v", err)
	} else {
		utils.RecordFileOperation()
		utils.Info("批量翻译完成，耗时: %v", time.Since(startTime))
	}
}

// 简化的缓存相关方法
func (t *LLMTranslator) GetMissingTags(tags []string) []string {
	return t.cache.GetMissingTexts(tags, TagCache)
}

func (t *LLMTranslator) GetMissingArticles(articles []string) []string {
	return t.cache.GetMissingTexts(articles, ArticleCache)
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
func (t *LLMTranslator) ClearArticleCache() error { return t.cache.Clear(ArticleCache) }

func (t *LLMTranslator) GetCacheStats() (int, int) {
	tagTotal, tagExpired := t.cache.GetStats(TagCache)
	articleTotal, articleExpired := t.cache.GetStats(ArticleCache)
	return tagTotal + articleTotal, tagExpired + articleExpired
}

// isEnglishOnly 检查字符串是否只包含英文字符
func isEnglishOnly(s string) bool {
	for _, r := range s {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '-' || r == '_' || r == ' ') {
			return false
		}
	}
	return true
}

// normalizeSlug 标准化slug格式
func normalizeSlug(s string) string {
	// 转为小写
	s = strings.ToLower(s)

	// 移除引号和其他特殊字符
	s = strings.Trim(s, "\"'`")

	// 替换空格为连字符
	s = strings.ReplaceAll(s, " ", "-")

	// 移除非法字符，只保留字母、数字和连字符
	reg := regexp.MustCompile(`[^a-z0-9\-]`)
	s = reg.ReplaceAllString(s, "")

	// 移除多个连续的连字符
	reg = regexp.MustCompile(`-+`)
	s = reg.ReplaceAllString(s, "-")

	// 移除开头和结尾的连字符
	s = strings.Trim(s, "-")

	return s
}

// fallbackSlug 当翻译失败时的备用方案
func fallbackSlug(tag string) string {
	// 预定义的映射表作为备用
	fallbackTranslations := map[string]string{
		"人工智能":       "artificial-intelligence",
		"机器学习":       "machine-learning",
		"深度学习":       "deep-learning",
		"前端开发":       "frontend-development",
		"后端开发":       "backend-development",
		"JavaScript": "javascript",
		"Python":     "python",
		"Go":         "golang",
		"技术":         "technology",
		"教程":         "tutorial",
		"编程":         "programming",
		"开发":         "development",
	}

	if slug, exists := fallbackTranslations[tag]; exists {
		return slug
	}

	// 最后的备用方案：简单处理
	return normalizeSlug(tag)
}

// FallbackSlug 导出的备用slug生成函数
func FallbackSlug(tag string) string {
	return fallbackSlug(tag)
}
