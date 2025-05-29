package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hugo-content-suite/config"
	"hugo-content-suite/translator"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// TranslationUtils 翻译工具
type TranslationUtils struct {
	translator *translator.LLMTranslator
	cache      *translator.TranslationCache
}

// NewTranslationUtils 创建翻译工具实例
func NewTranslationUtils() *TranslationUtils {
	cache := translator.NewTranslationCache()
	cache.Load() // 加载缓存

	return &TranslationUtils{
		translator: translator.NewLLMTranslator(),
		cache:      cache,
	}
}

// TestConnection 测试连接
func (t *TranslationUtils) TestConnection() error {
	return t.translator.TestConnection()
}

// ContainsEnglish 检查文本是否包含英文
func (t *TranslationUtils) ContainsEnglish(text string) bool {
	for _, r := range text {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			return true
		}
	}
	return false
}

// ContainsChinese 检查文本是否包含中文
func (t *TranslationUtils) ContainsChinese(text string) bool {
	for _, r := range text {
		if r >= 0x4e00 && r <= 0x9fff {
			return true
		}
	}
	return false
}

// IsOnlyEnglish 检查文本是否只包含英文（和标点符号、数字等）
func (t *TranslationUtils) IsOnlyEnglish(text string) bool {
	// 移除空白字符后检查
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return false
	}

	// 如果包含中文，则不是纯英文
	if t.ContainsChinese(trimmed) {
		return false
	}

	// 如果包含英文字母，且不包含中文，则认为是英文内容
	return t.ContainsEnglish(trimmed)
}

// SplitMixedText 分离中英文混合文本
func (t *TranslationUtils) SplitMixedText(text string) ([]TextSegment, bool) {
	if !t.ContainsChinese(text) {
		// 没有中文，无需翻译
		return []TextSegment{{Content: text, NeedsTranslation: false}}, false
	}

	if !t.ContainsEnglish(text) {
		// 没有英文，全部翻译
		return []TextSegment{{Content: text, NeedsTranslation: true}}, true
	}

	// 中英文混合，需要分割
	segments := t.segmentMixedText(text)
	hasTranslatableContent := false

	for _, segment := range segments {
		if segment.NeedsTranslation {
			hasTranslatableContent = true
			break
		}
	}

	return segments, hasTranslatableContent
}

// TextSegment 文本片段
type TextSegment struct {
	Content          string
	NeedsTranslation bool
}

// segmentMixedText 分割混合文本为片段
func (t *TranslationUtils) segmentMixedText(text string) []TextSegment {
	var segments []TextSegment
	var currentSegment strings.Builder
	var currentType bool // true表示中文，false表示英文
	var hasContent bool

	runes := []rune(text)

	for _, r := range runes {
		isChinese := r >= 0x4e00 && r <= 0x9fff
		isEnglish := (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')

		if isChinese || isEnglish {
			// 如果是第一个字符，或者类型改变了
			if !hasContent || (isChinese != currentType) {
				// 保存之前的片段
				if hasContent && currentSegment.Len() > 0 {
					segments = append(segments, TextSegment{
						Content:          currentSegment.String(),
						NeedsTranslation: currentType,
					})
					currentSegment.Reset()
				}

				currentType = isChinese
				hasContent = true
			}

			currentSegment.WriteRune(r)
		} else {
			// 标点符号、数字、空格等，附加到当前片段
			currentSegment.WriteRune(r)
		}
	}

	// 保存最后一个片段
	if hasContent && currentSegment.Len() > 0 {
		segments = append(segments, TextSegment{
			Content:          currentSegment.String(),
			NeedsTranslation: currentType,
		})
	}

	return segments
}

// CleanTranslationResult 清理翻译结果
func (t *TranslationUtils) CleanTranslationResult(result string) string {
	cfg := config.GetGlobalConfig()

	// 移除首尾空白
	result = strings.TrimSpace(result)

	// 移除常见的多余前缀
	unwantedPrefixes := []string{
		"Translation:", "English:", "Japanese:", "Korean:",
		"The translation is:", "Here is the translation:",
		"Translated:", "Answer:", "Result:", "Output:",
		"翻译:", "英文:", "日文:", "韩文:",
	}

	for _, prefix := range unwantedPrefixes {
		if strings.HasPrefix(result, prefix) {
			result = strings.TrimSpace(strings.TrimPrefix(result, prefix))
		}
	}

	// 使用配置中的清理模式
	for _, pattern := range cfg.Translation.CleanupPatterns {
		if strings.HasPrefix(result, pattern) {
			result = strings.TrimSpace(strings.TrimPrefix(result, pattern))
		}
	}

	// 移除多层引号
	for strings.HasPrefix(result, "\"") && strings.HasSuffix(result, "\"") && len(result) > 2 {
		inner := result[1 : len(result)-1]
		if !strings.Contains(inner, "\"") || strings.Count(inner, "\"")%2 == 0 {
			result = strings.TrimSpace(inner)
		} else {
			break
		}
	}

	// 移除句号结尾
	if strings.HasSuffix(result, ".") && !strings.Contains(result, ". ") {
		result = strings.TrimSpace(strings.TrimSuffix(result, "."))
	}

	// 移除多余的换行符和空格
	result = strings.ReplaceAll(result, "\n", " ")
	result = strings.ReplaceAll(result, "\r", " ")
	result = regexp.MustCompile(`\s+`).ReplaceAllString(result, " ")

	return strings.TrimSpace(result)
}

// RemoveQuotes 移除译文中的所有引号
func (t *TranslationUtils) RemoveQuotes(text string) string {
	quotes := []string{"\"", "'", "'", "'", "„", "‚", "‹", "›", "«", "»"}
	for _, quote := range quotes {
		text = strings.ReplaceAll(text, quote, "")
	}
	return strings.TrimSpace(text)
}

// FormatSlugField 格式化slug字段
func (t *TranslationUtils) FormatSlugField(slug string) string {
	slug = strings.ToLower(slug)
	slug = strings.ReplaceAll(slug, " ", "-")

	reg := regexp.MustCompile(`[^a-z0-9\-]`)
	slug = reg.ReplaceAllString(slug, "")

	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	return strings.Trim(slug, "-")
}

// TranslateToLanguage 翻译文本到指定语言（带缓存）
func (t *TranslationUtils) TranslateToLanguage(content, targetLang string) (string, error) {
	// 检查是否只包含英文，如果是则直接返回
	if t.IsOnlyEnglish(content) {
		return content, nil
	}

	// 分离中英文内容
	segments, hasTranslatableContent := t.SplitMixedText(content)

	// 如果没有需要翻译的内容，直接返回原文
	if !hasTranslatableContent {
		return content, nil
	}

	// 如果只有一个片段且需要翻译，使用原有逻辑
	if len(segments) == 1 && segments[0].NeedsTranslation {
		return t.translateSingleText(content, targetLang)
	}

	// 处理混合文本
	var result strings.Builder
	for _, segment := range segments {
		if segment.NeedsTranslation {
			translated, err := t.translateSingleText(segment.Content, targetLang)
			if err != nil {
				return "", err
			}
			result.WriteString(translated)
		} else {
			result.WriteString(segment.Content)
		}
	}

	return result.String(), nil
}

// translateSingleText 翻译单个文本片段
func (t *TranslationUtils) translateSingleText(content, targetLang string) (string, error) {
	// 先检查缓存
	cacheKey := fmt.Sprintf("%s:%s", targetLang, content)
	if cached, found := t.cache.Get(cacheKey, translator.TagCache); found {
		return cached, nil
	}

	// 缓存未命中，调用翻译服务
	result, err := t.translateWithAPI(content, targetLang)
	if err != nil {
		return "", err
	}

	// 保存到缓存
	t.cache.Set(cacheKey, result, translator.TagCache)

	return result, nil
}

// BatchTranslateWithCache 批量翻译（优先使用缓存）
func (t *TranslationUtils) BatchTranslateWithCache(texts []string, targetLang string, cacheType translator.CacheType) (map[string]string, error) {
	result := make(map[string]string)
	var missingTexts []string

	// 检查缓存
	for _, text := range texts {
		cacheKey := fmt.Sprintf("%s:%s", targetLang, text)
		if cached, found := t.cache.Get(cacheKey, cacheType); found {
			result[text] = cached
		} else {
			missingTexts = append(missingTexts, text)
		}
	}

	fmt.Printf("📄 缓存命中: %d 个, 需要翻译: %d 个\n",
		len(texts)-len(missingTexts), len(missingTexts))

	// 翻译缺失的文本
	if len(missingTexts) > 0 {
		for i, text := range missingTexts {
			fmt.Printf("  [%d/%d] 翻译: %s -> ", i+1, len(missingTexts), text)

			translated, err := t.translateWithAPI(text, targetLang)
			if err != nil {
				fmt.Printf("失败 (%v)\n", err)
				// 使用备用方案
				translated = t.FallbackSlug(text)
			} else {
				translated = t.CleanTranslationResult(translated)
			}

			result[text] = translated

			// 保存到缓存
			cacheKey := fmt.Sprintf("%s:%s", targetLang, text)
			t.cache.Set(cacheKey, translated, cacheType)

			fmt.Printf("%s\n", translated)

			// 添加延迟避免API频率限制
			if i < len(missingTexts)-1 {
				time.Sleep(100 * time.Millisecond)
			}
		}

		// 保存缓存到文件
		if err := t.cache.Save(); err != nil {
			fmt.Printf("⚠️ 保存缓存失败: %v\n", err)
		}
	}

	return result, nil
}

// translateWithAPI 直接调用API翻译
func (t *TranslationUtils) translateWithAPI(content, targetLang string) (string, error) {
	cfg := config.GetGlobalConfig()
	targetLangName := cfg.Language.LanguageNames[targetLang]
	if targetLangName == "" {
		targetLangName = targetLang
	}

	var prompt string
	switch targetLang {
	case "ja":
		prompt = fmt.Sprintf(`Please translate this Chinese text to Japanese. Return ONLY the Japanese translation:

%s`, content)
	case "ko":
		prompt = fmt.Sprintf(`Please translate this Chinese text to Korean. Return ONLY the Korean translation:

%s`, content)
	default:
		prompt = fmt.Sprintf(`Please translate this Chinese text to English. Return ONLY the English translation:

%s`, content)
	}

	request := translator.LMStudioRequest{
		Model: cfg.LMStudio.Model,
		Messages: []translator.Message{
			{
				Role:    "system",
				Content: fmt.Sprintf("You are a professional translator. You translate Chinese to %s accurately and concisely.", targetLangName),
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Stream: false,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %v", err)
	}

	client := &http.Client{Timeout: time.Duration(cfg.LMStudio.Timeout) * time.Second}
	resp, err := client.Post(cfg.LMStudio.URL, "application/json", bytes.NewBuffer(jsonData))
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

	var response translator.LMStudioResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("没有获取到翻译结果")
	}

	result := strings.TrimSpace(response.Choices[0].Message.Content)
	return t.CleanTranslationResult(result), nil
}

// FallbackSlug 备用slug生成方案
func (t *TranslationUtils) FallbackSlug(tag string) string {
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
		"数据库":        "database",
		"网络":         "network",
		"安全":         "security",
		"算法":         "algorithm",
		"框架":         "framework",
		"工具":         "tools",
		"设计":         "design",
		"产品":         "product",
	}

	if slug, exists := fallbackTranslations[tag]; exists {
		return slug
	}

	// 简单处理
	slug := strings.ToLower(tag)
	slug = strings.ReplaceAll(slug, " ", "-")
	reg := regexp.MustCompile(`[^\w\x{4e00}-\x{9fff}\-]`)
	slug = reg.ReplaceAllString(slug, "")
	return strings.Trim(slug, "-")
}

// SaveCache 保存缓存
func (t *TranslationUtils) SaveCache() error {
	return t.cache.Save()
}

// GetCacheStats 获取缓存统计
func (t *TranslationUtils) GetCacheStats() string {
	return t.cache.GetInfo()
}

// IsMarkdownStructuralElement 检查是否为markdown结构元素
func (t *TranslationUtils) IsMarkdownStructuralElement(line string) bool {
	trimmed := strings.TrimSpace(line)

	// 代码块
	if strings.HasPrefix(trimmed, "```") {
		return true
	}

	// 水平分割线
	if matched, _ := regexp.MatchString(`^(-{3,}|\*{3,}|_{3,})$`, trimmed); matched {
		return true
	}

	// HTML标签
	if matched, _ := regexp.MatchString(`^<[^>]+>.*</[^>]+>$`, trimmed); matched {
		return true
	}

	// 链接定义
	if matched, _ := regexp.MatchString(`^\[.+\]:\s+https?://`, trimmed); matched {
		return true
	}

	return false
}

// ProtectMarkdownSyntax 保护markdown语法
func (t *TranslationUtils) ProtectMarkdownSyntax(text string) (string, map[string]string) {
	protectedElements := make(map[string]string)
	counter := 0

	// 保护内联代码（优先级最高）
	inlineCodeRegex := regexp.MustCompile("`[^`\n]+`")
	text = inlineCodeRegex.ReplaceAllStringFunc(text, func(match string) string {
		placeholder := fmt.Sprintf("__INLINE_CODE_%d__", counter)
		protectedElements[placeholder] = match
		counter++
		return placeholder
	})

	// 保护链接
	linkRegex := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	text = linkRegex.ReplaceAllStringFunc(text, func(match string) string {
		placeholder := fmt.Sprintf("__LINK_%d__", counter)
		protectedElements[placeholder] = match
		counter++
		return placeholder
	})

	// 保护图片
	imageRegex := regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)
	text = imageRegex.ReplaceAllStringFunc(text, func(match string) string {
		placeholder := fmt.Sprintf("__IMAGE_%d__", counter)
		protectedElements[placeholder] = match
		counter++
		return placeholder
	})

	// 保护粗体（两个星号或下划线）
	boldRegex := regexp.MustCompile(`(\*\*|__)[^*_\n]+(\*\*|__)`)
	text = boldRegex.ReplaceAllStringFunc(text, func(match string) string {
		placeholder := fmt.Sprintf("__BOLD_%d__", counter)
		protectedElements[placeholder] = match
		counter++
		return placeholder
	})

	// 保护斜体（单个星号或下划线，但要避免与粗体冲突）
	italicRegex := regexp.MustCompile(`(?:^|[^*_])(\*|_)([^*_\n]+)(\*|_)(?:[^*_]|$)`)
	text = italicRegex.ReplaceAllStringFunc(text, func(match string) string {
		// 检查是否已经被保护
		for _, protected := range protectedElements {
			if strings.Contains(protected, match) {
				return match
			}
		}
		placeholder := fmt.Sprintf("__ITALIC_%d__", counter)
		protectedElements[placeholder] = match
		counter++
		return placeholder
	})

	// 保护删除线
	strikeRegex := regexp.MustCompile(`~~[^~\n]+~~`)
	text = strikeRegex.ReplaceAllStringFunc(text, func(match string) string {
		placeholder := fmt.Sprintf("__STRIKE_%d__", counter)
		protectedElements[placeholder] = match
		counter++
		return placeholder
	})

	// 保护HTML标签
	htmlRegex := regexp.MustCompile(`<[^>]+>`)
	text = htmlRegex.ReplaceAllStringFunc(text, func(match string) string {
		placeholder := fmt.Sprintf("__HTML_%d__", counter)
		protectedElements[placeholder] = match
		counter++
		return placeholder
	})

	return text, protectedElements
}

// RestoreMarkdownSyntax 恢复markdown语法
func (t *TranslationUtils) RestoreMarkdownSyntax(text string, protectedElements map[string]string) string {
	for placeholder, original := range protectedElements {
		text = strings.ReplaceAll(text, placeholder, original)
	}
	return text
}
