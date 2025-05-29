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

// ContainsChinese 检查文本是否包含中文
func (t *TranslationUtils) ContainsChinese(text string) bool {
	for _, r := range text {
		if r >= 0x4e00 && r <= 0x9fff {
			return true
		}
	}
	return false
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

// ProtectMarkdownElements 保护关键markdown元素（简化版）
func (t *TranslationUtils) ProtectMarkdownElements(text string) (string, map[string]string) {
	protectedElements := make(map[string]string)
	counter := 0

	// 1. 保护代码块（优先级最高）
	codeBlockRegex := regexp.MustCompile("(?s)```[^`]*```")
	text = codeBlockRegex.ReplaceAllStringFunc(text, func(match string) string {
		placeholder := fmt.Sprintf("__CODE_BLOCK_%d__", counter)
		protectedElements[placeholder] = match
		counter++
		return placeholder
	})

	// 2. 保护内联代码
	inlineCodeRegex := regexp.MustCompile("`[^`\n]+`")
	text = inlineCodeRegex.ReplaceAllStringFunc(text, func(match string) string {
		placeholder := fmt.Sprintf("__INLINE_CODE_%d__", counter)
		protectedElements[placeholder] = match
		counter++
		return placeholder
	})

	// 3. 保护完整链接
	linkRegex := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	text = linkRegex.ReplaceAllStringFunc(text, func(match string) string {
		placeholder := fmt.Sprintf("__LINK_%d__", counter)
		protectedElements[placeholder] = match
		counter++
		return placeholder
	})

	// 4. 保护图片
	imageRegex := regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)
	text = imageRegex.ReplaceAllStringFunc(text, func(match string) string {
		placeholder := fmt.Sprintf("__IMAGE_%d__", counter)
		protectedElements[placeholder] = match
		counter++
		return placeholder
	})

	// 5. 保护URL
	urlRegex := regexp.MustCompile(`https?://[^\s<>"{}|\\^` + "`" + `\[\]]+`)
	text = urlRegex.ReplaceAllStringFunc(text, func(match string) string {
		placeholder := fmt.Sprintf("__URL_%d__", counter)
		protectedElements[placeholder] = match
		counter++
		return placeholder
	})

	return text, protectedElements
}

// RestoreMarkdownElements 恢复保护的markdown元素
func (t *TranslationUtils) RestoreMarkdownElements(text string, protectedElements map[string]string) string {
	for placeholder, original := range protectedElements {
		text = strings.ReplaceAll(text, placeholder, original)
	}
	return text
}

// TranslateParagraphToLanguage 翻译段落到指定语言
func (t *TranslationUtils) TranslateParagraphToLanguage(paragraph, targetLang string) (string, error) {
	// 检查是否为标题行
	if t.isHeaderLine(paragraph) {
		return t.translateHeaderLine(paragraph, targetLang)
	}

	// 保护关键元素
	protectedContent, protectedElements := t.ProtectMarkdownElements(paragraph)

	// 翻译处理后的内容
	translatedContent, err := t.TranslateToLanguage(protectedContent, targetLang)
	if err != nil {
		return "", err
	}

	// 清理翻译结果
	translatedContent = t.CleanTranslationResult(translatedContent)

	// 恢复保护的元素
	finalContent := t.RestoreMarkdownElements(translatedContent, protectedElements)

	return finalContent, nil
}

// isHeaderLine 检查是否为标题行
func (t *TranslationUtils) isHeaderLine(line string) bool {
	trimmed := strings.TrimSpace(line)

	// 检查是否以#开头
	if !strings.HasPrefix(trimmed, "#") {
		return false
	}

	// 计算连续的#号数量
	hashCount := 0
	for _, r := range trimmed {
		if r == '#' {
			hashCount++
		} else {
			break
		}
	}

	// 必须是1-6个#号，且后面要么是空格要么是结尾
	if hashCount >= 1 && hashCount <= 6 {
		if len(trimmed) == hashCount {
			// 只有#号
			return true
		}
		if len(trimmed) > hashCount && trimmed[hashCount] == ' ' {
			// #号后面跟空格
			return true
		}
	}

	return false
}

// translateHeaderLine 翻译标题行
func (t *TranslationUtils) translateHeaderLine(line, targetLang string) (string, error) {
	trimmed := strings.TrimSpace(line)

	// 提取标题前缀和内容
	prefix, content := t.extractHeaderPrefix(trimmed)

	// 如果没有内容需要翻译，直接返回原行
	if content == "" || !t.ContainsChinese(content) {
		return line, nil
	}

	// 翻译标题内容
	translatedContent, err := t.TranslateToLanguage(content, targetLang)
	if err != nil {
		return "", err
	}

	// 清理翻译结果
	translatedContent = t.CleanTranslationResult(translatedContent)
	translatedContent = t.RemoveQuotes(translatedContent)

	// 重构标题行，保持原有的缩进
	originalIndent := ""
	if len(line) > len(trimmed) {
		originalIndent = line[:len(line)-len(trimmed)]
	}

	return originalIndent + prefix + translatedContent, nil
}

// extractHeaderPrefix 提取标题前缀
func (t *TranslationUtils) extractHeaderPrefix(line string) (string, string) {
	if !strings.HasPrefix(line, "#") {
		return "", line
	}

	// 计算连续的#号数量
	hashCount := 0
	for _, r := range line {
		if r == '#' {
			hashCount++
		} else {
			break
		}
	}

	// 构建前缀
	prefix := strings.Repeat("#", hashCount)

	// 提取内容
	content := ""
	if len(line) > hashCount {
		if line[hashCount] == ' ' {
			// 有空格，提取空格后的内容
			content = strings.TrimSpace(line[hashCount+1:])
			prefix += " "
		} else {
			// 没有空格，提取#号后的内容
			content = strings.TrimSpace(line[hashCount:])
			prefix += " " // 补充空格
		}
	}

	return prefix, content
}
