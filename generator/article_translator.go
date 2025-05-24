package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"tag-scanner/config"
	"tag-scanner/models"
	"tag-scanner/scanner"
	"tag-scanner/translator"
	"tag-scanner/utils"
	"time"
)

// ArticleTranslator 文章翻译器
type ArticleTranslator struct {
	contentDir string
	translator *translator.LLMTranslator
}

// ArticleTranslationPreview 文章翻译预览信息
type ArticleTranslationPreview struct {
	OriginalFile   string
	EnglishFile    string
	Title          string
	WordCount      int
	ParagraphCount int
	Status         string // "missing", "exists"
	EstimatedTime  string
}

// MarkdownNode 表示需要翻译的Markdown节点
type MarkdownNode struct {
	Type     string // text, heading, listitem, etc.
	Content  string // 原始中文内容
	Position int    // 在文档中的位置
	Level    int    // 标题级别（仅用于标题）
}

// NewArticleTranslator 创建新的文章翻译器
func NewArticleTranslator(contentDir string) *ArticleTranslator {
	return &ArticleTranslator{
		contentDir: contentDir,
		translator: translator.NewLLMTranslator(),
	}
}

// PreviewArticleTranslations 预览需要翻译的文章
func (a *ArticleTranslator) PreviewArticleTranslations() ([]ArticleTranslationPreview, error) {
	articles, err := scanner.ScanArticles(a.contentDir)
	if err != nil {
		return nil, fmt.Errorf("扫描文章失败: %v", err)
	}

	cfg := config.GetGlobalConfig()
	targetLanguages := cfg.Language.TargetLanguages

	var previews []ArticleTranslationPreview

	for _, article := range articles {
		if article.Title == "" {
			continue
		}

		// 为每种目标语言生成预览
		for _, targetLang := range targetLanguages {
			// 根据目标语言构建文件路径
			dir := filepath.Dir(article.FilePath)
			baseName := filepath.Base(article.FilePath)

			var targetFile string
			if strings.HasSuffix(baseName, ".md") {
				switch targetLang {
				case "ja":
					targetFile = filepath.Join(dir, "index.ja.md")
				case "ko":
					targetFile = filepath.Join(dir, "index.ko.md")
				default: // "en" 或其他
					targetFile = filepath.Join(dir, "index.en.md")
				}
			} else {
				continue
			}

			// 检查目标文件是否存在
			status := "missing"
			if _, err := os.Stat(targetFile); err == nil {
				status = "exists"
			}

			// 分析文章内容
			wordCount, paragraphCount := a.analyzeArticleContent(article.FilePath)
			estimatedTime := a.estimateTranslationTime(paragraphCount)

			preview := ArticleTranslationPreview{
				OriginalFile:   article.FilePath,
				EnglishFile:    targetFile,
				Title:          fmt.Sprintf("%s (%s)", article.Title, cfg.Language.LanguageNames[targetLang]),
				WordCount:      wordCount,
				ParagraphCount: paragraphCount,
				Status:         status,
				EstimatedTime:  estimatedTime,
			}

			previews = append(previews, preview)
		}
	}

	return previews, nil
}

// TranslateArticles 翻译文章到多种语言
func (a *ArticleTranslator) TranslateArticles(mode string) error {
	cfg := config.GetGlobalConfig()
	targetLanguages := cfg.Language.TargetLanguages

	// 获取所有文章
	articles, err := scanner.ScanArticles(a.contentDir)
	if err != nil {
		return fmt.Errorf("扫描文章失败: %v", err)
	}

	var targetArticles []models.Article
	for _, article := range articles {
		if article.Title == "" {
			continue
		}
		targetArticles = append(targetArticles, article)
	}

	if len(targetArticles) == 0 {
		fmt.Printf("没有需要翻译的文章\n")
		return nil
	}

	// 测试连接
	fmt.Printf("正在测试与LM Studio的连接...\n")
	if err := a.translator.TestConnection(); err != nil {
		return fmt.Errorf("无法连接到LM Studio: %v", err)
	}
	fmt.Printf("LM Studio连接成功！\n")

	totalSuccessCount := 0
	totalErrorCount := 0

	// 按语言顺序翻译
	for langIndex, targetLang := range targetLanguages {
		targetLangName := cfg.Language.LanguageNames[targetLang]
		if targetLangName == "" {
			targetLangName = targetLang
		}

		fmt.Printf("\n🌐 开始翻译为 %s (%d/%d)\n", targetLangName, langIndex+1, len(targetLanguages))
		utils.Info("开始翻译为 %s (%d/%d)", targetLangName, langIndex+1, len(targetLanguages))

		successCount := 0
		errorCount := 0

		for i, article := range targetArticles {
			// 构建目标文件路径
			dir := filepath.Dir(article.FilePath)
			var targetFile string
			switch targetLang {
			case "ja":
				targetFile = filepath.Join(dir, "index.ja.md")
			case "ko":
				targetFile = filepath.Join(dir, "index.ko.md")
			default: // "en" 或其他
				targetFile = filepath.Join(dir, "index.en.md")
			}

			// 检查是否需要翻译
			shouldTranslate := false
			if _, err := os.Stat(targetFile); os.IsNotExist(err) {
				if mode == "missing" || mode == "all" {
					shouldTranslate = true
				}
			} else {
				if mode == "update" || mode == "all" {
					shouldTranslate = true
				}
			}

			if !shouldTranslate {
				continue
			}

			fmt.Printf("\n处理文章 (%d/%d): %s\n", i+1, len(targetArticles), article.Title)
			fmt.Printf("目标语言: %s\n", targetLangName)
			fmt.Printf("目标文件: %s\n", targetFile)

			preview := ArticleTranslationPreview{
				OriginalFile: article.FilePath,
				EnglishFile:  targetFile,
				Title:        article.Title,
			}

			if err := a.translateSingleArticleToLanguage(preview, targetLang); err != nil {
				fmt.Printf("❌ 翻译失败: %v\n", err)
				errorCount++
				totalErrorCount++
			} else {
				fmt.Printf("✅ 翻译完成: %s\n", targetFile)
				successCount++
				totalSuccessCount++
			}
		}

		fmt.Printf("\n%s 翻译完成:\n", targetLangName)
		fmt.Printf("- 成功翻译: %d 篇\n", successCount)
		fmt.Printf("- 翻译失败: %d 篇\n", errorCount)
	}

	fmt.Printf("\n🎉 多语言翻译全部完成！\n")
	fmt.Printf("- 目标语言: %v\n", targetLanguages)
	fmt.Printf("- 总成功翻译: %d 篇\n", totalSuccessCount)
	fmt.Printf("- 总翻译失败: %d 篇\n", totalErrorCount)

	return nil
}

// translateSingleArticleToLanguage 翻译单篇文章到指定语言
func (a *ArticleTranslator) translateSingleArticleToLanguage(preview ArticleTranslationPreview, targetLang string) error {
	utils.Info("开始翻译文章到 %s: %s", targetLang, preview.OriginalFile)
	utils.Info("目标文件: %s", preview.EnglishFile)

	// 读取原文件
	content, err := os.ReadFile(preview.OriginalFile)
	if err != nil {
		utils.Error("读取原文件失败: %s, 错误: %v", preview.OriginalFile, err)
		return fmt.Errorf("读取原文件失败: %v", err)
	}

	// 解析文章结构
	frontMatter, bodyContent := a.parseArticleContent(string(content))

	// 翻译前置数据
	translatedFrontMatter, err := a.translateFrontMatterToLanguage(frontMatter, targetLang)
	if err != nil {
		utils.Error("翻译前置数据失败: %v", err)
		return fmt.Errorf("翻译前置数据失败: %v", err)
	}

	// 翻译正文
	translatedBody, err := a.translateArticleBodyToLanguage(bodyContent, targetLang)
	if err != nil {
		utils.Error("翻译正文失败: %v", err)
		return fmt.Errorf("翻译正文失败: %v", err)
	}

	// 合成最终内容
	finalContent := a.combineTranslatedContent(translatedFrontMatter, translatedBody)

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(preview.EnglishFile), 0755); err != nil {
		utils.Error("创建目录失败: %s, 错误: %v", filepath.Dir(preview.EnglishFile), err)
		return fmt.Errorf("创建目录失败: %v", err)
	}

	// 写入目标文件
	if err := os.WriteFile(preview.EnglishFile, []byte(finalContent), 0644); err != nil {
		utils.Error("写入目标文件失败: %s, 错误: %v", preview.EnglishFile, err)
		return fmt.Errorf("写入目标文件失败: %v", err)
	}

	utils.Info("文章翻译完成 (%s): %s", targetLang, preview.EnglishFile)
	return nil
}

// translateFieldContentToLanguage 翻译字段内容到指定语言
func (a *ArticleTranslator) translateFieldContentToLanguage(content, targetLang string) (string, error) {
	cfg := config.GetGlobalConfig()
	targetLangName := cfg.Language.LanguageNames[targetLang]
	if targetLangName == "" {
		targetLangName = targetLang
	}

	// 根据目标语言调整提示词
	var prompt string
	switch targetLang {
	case "ja":
		prompt = fmt.Sprintf(`Please translate this Chinese text to Japanese. Return ONLY the Japanese translation, no explanations or additional text:

%s`, content)
	case "ko":
		prompt = fmt.Sprintf(`Please translate this Chinese text to Korean. Return ONLY the Korean translation, no explanations or additional text:

%s`, content)
	default: // "en" 或其他
		prompt = fmt.Sprintf(`Please translate this Chinese text to English. Return ONLY the English translation, no explanations or additional text:

%s`, content)
	}

	request := translator.LMStudioRequest{
		Model: cfg.LMStudio.Model,
		Messages: []translator.Message{
			{
				Role:    "system",
				Content: fmt.Sprintf("You are a professional translator. You translate Chinese to %s accurately and concisely. You only return the translation without any additional text, explanations, or formatting.", targetLangName),
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
	result = a.cleanTranslationResult(result)

	utils.Info("字段翻译完成 (%s) - 原文: %s, 译文: %s", targetLangName, content, result)

	return result, nil
}

// translateFrontMatterToLanguage 翻译前置数据到指定语言
func (a *ArticleTranslator) translateFrontMatterToLanguage(frontMatter, targetLang string) (string, error) {
	if frontMatter == "" {
		return "", nil
	}

	cfg := config.GetGlobalConfig()
	targetLangName := cfg.Language.LanguageNames[targetLang]
	if targetLangName == "" {
		targetLangName = targetLang
	}

	fmt.Printf("翻译前置数据到 %s...\n", targetLangName)

	lines := strings.Split(frontMatter, "\n")
	var translatedLines []string

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		if trimmedLine == "---" {
			translatedLines = append(translatedLines, line)
			continue
		}

		// 翻译标题字段
		if strings.HasPrefix(trimmedLine, "title:") {
			title := a.extractFieldValue(trimmedLine, "title:")
			if title != "" && a.containsChinese(title) {
				fmt.Printf("  title: %s -> ", title)
				translatedTitle, err := a.translateFieldContentToLanguage(title, targetLang)
				if err != nil {
					fmt.Printf("翻译失败\n")
					translatedLines = append(translatedLines, line)
				} else {
					translatedTitle = a.removeQuotes(translatedTitle)
					fmt.Printf("%s\n", translatedTitle)
					translatedLines = append(translatedLines, fmt.Sprintf("title: \"%s\"", translatedTitle))
				}
			} else {
				translatedLines = append(translatedLines, line)
			}
			continue
		}

		// 翻译描述字段
		if strings.HasPrefix(trimmedLine, "description:") {
			description := a.extractFieldValue(trimmedLine, "description:")
			if description != "" && a.containsChinese(description) {
				fmt.Printf("  description: %s -> ", description)
				translatedDescription, err := a.translateFieldContentToLanguage(description, targetLang)
				if err != nil {
					fmt.Printf("翻译失败\n")
					translatedLines = append(translatedLines, line)
				} else {
					translatedDescription = a.removeQuotes(translatedDescription)
					fmt.Printf("%s\n", translatedDescription)
					translatedLines = append(translatedLines, fmt.Sprintf("description: \"%s\"", translatedDescription))
				}
			} else {
				translatedLines = append(translatedLines, line)
			}
			continue
		}

		// 翻译副标题
		if strings.HasPrefix(trimmedLine, "subtitle:") {
			subtitle := a.extractFieldValue(trimmedLine, "subtitle:")
			if subtitle != "" && a.containsChinese(subtitle) {
				fmt.Printf("  subtitle: %s -> ", subtitle)
				translatedSubtitle, err := a.translateFieldContentToLanguage(subtitle, targetLang)
				if err != nil {
					fmt.Printf("翻译失败\n")
					translatedLines = append(translatedLines, line)
				} else {
					translatedSubtitle = a.removeQuotes(translatedSubtitle)
					fmt.Printf("%s\n", translatedSubtitle)
					translatedLines = append(translatedLines, fmt.Sprintf("subtitle: \"%s\"", translatedSubtitle))
				}
			} else {
				translatedLines = append(translatedLines, line)
			}
			continue
		}

		// 翻译摘要
		if strings.HasPrefix(trimmedLine, "summary:") {
			summary := a.extractFieldValue(trimmedLine, "summary:")
			if summary != "" && a.containsChinese(summary) {
				fmt.Printf("  summary: %s -> ", summary)
				translatedSummary, err := a.translateFieldContentToLanguage(summary, targetLang)
				if err != nil {
					fmt.Printf("翻译失败\n")
					translatedLines = append(translatedLines, line)
				} else {
					translatedSummary = a.removeQuotes(translatedSummary)
					fmt.Printf("%s\n", translatedSummary)
					translatedLines = append(translatedLines, fmt.Sprintf("summary: \"%s\"", translatedSummary))
				}
			} else {
				translatedLines = append(translatedLines, line)
			}
			continue
		}

		// 翻译slug字段
		if strings.HasPrefix(trimmedLine, "slug:") {
			slug := a.extractFieldValue(trimmedLine, "slug:")
			if slug != "" && a.containsChinese(slug) {
				fmt.Printf("  slug: %s -> ", slug)
				translatedSlug, err := a.translateFieldContentToLanguage(slug, targetLang)
				if err != nil {
					fmt.Printf("翻译失败\n")
					translatedLines = append(translatedLines, line)
				} else {
					translatedSlug = a.removeQuotes(translatedSlug)
					translatedSlug = a.formatSlugField(translatedSlug)
					fmt.Printf("%s\n", translatedSlug)
					translatedLines = append(translatedLines, fmt.Sprintf("slug: \"%s\"", translatedSlug))
				}
			} else {
				translatedLines = append(translatedLines, line)
			}
			continue
		}

		// 翻译标签数组
		if strings.HasPrefix(trimmedLine, "tags:") {
			tags := a.extractArrayField(trimmedLine, "tags:")
			if len(tags) > 0 {
				translatedTags := a.translateArrayFieldToLanguage(tags, "tags", targetLang)
				translatedLines = append(translatedLines, fmt.Sprintf("tags: %s", a.formatArrayField(translatedTags)))
			} else {
				translatedLines = append(translatedLines, line)
			}
			continue
		}

		// 翻译分类数组
		if strings.HasPrefix(trimmedLine, "categories:") {
			categories := a.extractArrayField(trimmedLine, "categories:")
			if len(categories) > 0 {
				translatedCategories := a.translateArrayFieldToLanguage(categories, "categories", targetLang)
				translatedLines = append(translatedLines, fmt.Sprintf("categories: %s", a.formatArrayField(translatedCategories)))
			} else {
				translatedLines = append(translatedLines, line)
			}
			continue
		}

		// 翻译作者数组
		if strings.HasPrefix(trimmedLine, "authors:") {
			authors := a.extractArrayField(trimmedLine, "authors:")
			if len(authors) > 0 {
				translatedAuthors := a.translateArrayFieldToLanguage(authors, "authors", targetLang)
				translatedLines = append(translatedLines, fmt.Sprintf("authors: %s", a.formatArrayField(translatedAuthors)))
			} else {
				translatedLines = append(translatedLines, line)
			}
			continue
		}

		// 其他字段保持不变
		translatedLines = append(translatedLines, line)
	}

	return strings.Join(translatedLines, "\n"), nil
}

// translateArticleBodyToLanguage 翻译正文到指定语言
func (a *ArticleTranslator) translateArticleBodyToLanguage(body, targetLang string) (string, error) {
	if strings.TrimSpace(body) == "" {
		return body, nil
	}

	cfg := config.GetGlobalConfig()
	targetLangName := cfg.Language.LanguageNames[targetLang]
	if targetLangName == "" {
		targetLangName = targetLang
	}

	fmt.Printf("\n翻译正文到 %s (%d 字符)...\n", targetLangName, len(body))

	return a.translateContentByLinesToLanguage(body, targetLang)
}

// translateContentByLinesToLanguage 按行翻译内容到指定语言
func (a *ArticleTranslator) translateContentByLinesToLanguage(content, targetLang string) (string, error) {
	cfg := config.GetGlobalConfig()
	lines := strings.Split(content, "\n")
	var result []string

	inCodeBlock := false
	translationCount := 0

	for _, line := range lines {
		// 检测代码块
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			inCodeBlock = !inCodeBlock
			result = append(result, line)
			continue
		}

		// 代码块内容直接保留
		if inCodeBlock {
			result = append(result, line)
			continue
		}

		// 空行直接保留
		if strings.TrimSpace(line) == "" {
			result = append(result, line)
			continue
		}

		// 检查是否包含中文
		if !a.containsChinese(line) {
			result = append(result, line)
			continue
		}

		// 需要翻译的行
		translationCount++
		fmt.Printf("  [%d] ", translationCount)

		translatedLine, err := a.translateSingleLineToLanguage(line, translationCount, targetLang)
		if err != nil {
			fmt.Printf("翻译失败\n")
			result = append(result, line) // 翻译失败保持原文
		} else {
			fmt.Printf("完成\n")
			result = append(result, translatedLine)
		}

		// 添加延迟避免API频率限制
		if cfg.Translation.DelayBetweenMs > 0 {
			time.Sleep(time.Duration(cfg.Translation.DelayBetweenMs) * time.Millisecond)
		}
	}

	return strings.Join(result, "\n"), nil
}

// translateSingleLineToLanguage 翻译单行内容到指定语言
func (a *ArticleTranslator) translateSingleLineToLanguage(line string, lineNum int, targetLang string) (string, error) {
	trimmedLine := strings.TrimSpace(line)

	// 提取Markdown格式前缀
	var prefix, content, suffix string

	// 处理标题
	if strings.HasPrefix(trimmedLine, "#") {
		match := regexp.MustCompile(`^(#+\s*)`).FindString(trimmedLine)
		if match != "" {
			prefix = match
			content = strings.TrimSpace(strings.TrimPrefix(trimmedLine, match))
		}
	} else if strings.HasPrefix(trimmedLine, "- ") || strings.HasPrefix(trimmedLine, "* ") {
		// 处理无序列表
		if strings.HasPrefix(trimmedLine, "- ") {
			prefix = "- "
			content = strings.TrimSpace(strings.TrimPrefix(trimmedLine, "- "))
		} else {
			prefix = "* "
			content = strings.TrimSpace(strings.TrimPrefix(trimmedLine, "* "))
		}
	} else {
		// 普通段落
		content = trimmedLine
	}

	// 如果没有可翻译的内容，直接返回
	if strings.TrimSpace(content) == "" || !a.containsChinese(content) {
		return line, nil
	}

	// 翻译纯文本内容到指定语言
	translatedContent, err := a.translatePlainTextToLanguage(content, lineNum, targetLang)
	if err != nil {
		return "", err
	}

	// 重新组合
	leadingSpaces := ""
	if len(line) > len(strings.TrimLeft(line, " \t")) {
		leadingSpaces = line[:len(line)-len(strings.TrimLeft(line, " \t"))]
	}

	return leadingSpaces + prefix + translatedContent + suffix, nil
}

// translatePlainTextToLanguage 翻译纯文本内容到指定语言
func (a *ArticleTranslator) translatePlainTextToLanguage(text string, lineNum int, targetLang string) (string, error) {
	cfg := config.GetGlobalConfig()
	targetLangName := cfg.Language.LanguageNames[targetLang]
	if targetLangName == "" {
		targetLangName = targetLang
	}

	// 清理文本
	cleanText := strings.TrimSpace(text)
	cleanText = regexp.MustCompile(`\s+`).ReplaceAllString(cleanText, " ")

	// 根据目标语言调整提示词
	var prompt string
	switch targetLang {
	case "ja":
		prompt = fmt.Sprintf(`Translate this Chinese text to Japanese. Return ONLY the Japanese translation:

%s`, cleanText)
	case "ko":
		prompt = fmt.Sprintf(`Translate this Chinese text to Korean. Return ONLY the Korean translation:

%s`, cleanText)
	default: // "en" 或其他
		prompt = fmt.Sprintf(`Translate this Chinese text to English. Return ONLY the English translation:

%s`, cleanText)
	}

	request := translator.LMStudioRequest{
		Model: cfg.LMStudio.Model,
		Messages: []translator.Message{
			{
				Role:    "system",
				Content: fmt.Sprintf("You are a professional translator. Translate Chinese to %s accurately. Return only the translation without explanations or formatting.", targetLangName),
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
	result = a.cleanTranslationResult(result)

	return result, nil
}

// removeQuotes 移除译文中的所有引号
func (a *ArticleTranslator) removeQuotes(text string) string {
	// 移除双引号
	text = strings.ReplaceAll(text, "\"", "")
	// 移除单引号
	text = strings.ReplaceAll(text, "'", "")
	// 移除中文引号
	text = strings.ReplaceAll(text, "“", "")
	text = strings.ReplaceAll(text, "”", "")
	text = strings.ReplaceAll(text, "‘", "")
	text = strings.ReplaceAll(text, "’", "")
	// 移除其他类型的引号
	text = strings.ReplaceAll(text, "„", "")
	text = strings.ReplaceAll(text, "‚", "")
	text = strings.ReplaceAll(text, "‹", "")
	text = strings.ReplaceAll(text, "›", "")
	text = strings.ReplaceAll(text, "«", "")
	text = strings.ReplaceAll(text, "»", "")

	// 清理空格
	text = strings.TrimSpace(text)

	return text
}

// formatSlugField 格式化slug字段，转换为URL友好格式
func (a *ArticleTranslator) formatSlugField(slug string) string {
	// 转换为小写
	slug = strings.ToLower(slug)

	// 替换空格为连字符
	slug = strings.ReplaceAll(slug, " ", "-")

	// 移除特殊字符，只保留字母、数字和连字符
	reg := regexp.MustCompile(`[^a-z0-9\-]`)
	slug = reg.ReplaceAllString(slug, "")

	// 移除连续的连字符
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	// 移除首尾的连字符
	slug = strings.Trim(slug, "-")

	utils.Debug("格式化slug: %s", slug)
	return slug
}

// extractFieldValue 提取字段值
func (a *ArticleTranslator) extractFieldValue(line, prefix string) string {
	value := strings.TrimSpace(strings.TrimPrefix(line, prefix))
	value = strings.Trim(value, "\"'")
	return value
}

// extractArrayField 提取数组字段
func (a *ArticleTranslator) extractArrayField(line, prefix string) []string {
	value := strings.TrimSpace(strings.TrimPrefix(line, prefix))

	// 移除方括号
	value = strings.Trim(value, "[]")

	if value == "" {
		return []string{}
	}

	// 分割数组元素
	parts := strings.Split(value, ",")
	var result []string

	for _, part := range parts {
		part = strings.TrimSpace(part)
		part = strings.Trim(part, "\"'")
		if part != "" {
			result = append(result, part)
		}
	}

	return result
}

// translateArrayFieldToLanguage 翻译数组字段到指定语言
func (a *ArticleTranslator) translateArrayFieldToLanguage(items []string, fieldType, targetLang string) []string {
	var translated []string

	cfg := config.GetGlobalConfig()
	targetLangName := cfg.Language.LanguageNames[targetLang]
	if targetLangName == "" {
		targetLangName = targetLang
	}

	fmt.Printf("  %s: ", fieldType)
	utils.Info("开始翻译%s数组到%s: %v", fieldType, targetLangName, items)

	for i, item := range items {
		utils.Debug("处理数组项目 [%d/%d]: %s", i+1, len(items), item)

		if a.containsChinese(item) {
			fmt.Printf("%s -> ", item)
			utils.Info("翻译数组项目 [%d/%d]: %s", i+1, len(items), item)

			translatedItem, err := a.translateFieldContentToLanguage(item, targetLang)
			if err != nil {
				fmt.Printf("失败 ")
				utils.Warn("数组项目翻译失败 [%d/%d] - %s: %s, 错误: %v", i+1, len(items), fieldType, item, err)
				translated = append(translated, item)
			} else {
				// 移除译文中的引号
				translatedItem = a.removeQuotes(translatedItem)
				fmt.Printf("%s ", translatedItem)
				utils.Info("数组项目翻译成功 [%d/%d] - %s: %s -> %s", i+1, len(items), fieldType, item, translatedItem)
				translated = append(translated, translatedItem)
			}
		} else {
			utils.Debug("跳过数组项目 [%d/%d] - 无中文: %s", i+1, len(items), item)
			translated = append(translated, item)
		}
	}

	fmt.Printf("\n")
	utils.Info("%s数组翻译完成: %v -> %v", fieldType, items, translated)
	return translated
}

// formatArrayField 格式化数组字段，避免多余引号
func (a *ArticleTranslator) formatArrayField(items []string) string {
	if len(items) == 0 {
		return "[]"
	}

	var quotedItems []string
	for _, item := range items {
		// 清理可能存在的多余引号，并确保不包含引号
		cleanItem := a.removeQuotes(item)
		quotedItems = append(quotedItems, fmt.Sprintf("\"%s\"", cleanItem))
	}

	return fmt.Sprintf("[%s]", strings.Join(quotedItems, ", "))
}

// cleanTranslationResult 清理翻译结果，移除多余的提示词或格式
func (a *ArticleTranslator) cleanTranslationResult(result string) string {
	cfg := config.GetGlobalConfig()

	// 移除首尾空白
	result = strings.TrimSpace(result)

	// 移除常见的多余前缀
	unwantedPrefixes := []string{
		"Translation:",
		"English:",
		"Japanese:",
		"Korean:",
		"The translation is:",
		"Here is the translation:",
		"The English translation is:",
		"The Japanese translation is:",
		"The Korean translation is:",
		"Translated:",
		"Answer:",
		"Result:",
		"Output:",
		"English translation:",
		"Japanese translation:",
		"Korean translation:",
		"翻译:",
		"英文:",
		"日文:",
		"韩文:",
		"Translation: ",
		"English: ",
		"Japanese: ",
		"Korean: ",
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

	// 移除多层引号（更严格的处理）
	for strings.HasPrefix(result, "\"") && strings.HasSuffix(result, "\"") && len(result) > 2 {
		inner := result[1 : len(result)-1]
		if !strings.Contains(inner, "\"") || strings.Count(inner, "\"")%2 == 0 {
			result = inner
			result = strings.TrimSpace(result)
		} else {
			break
		}
	}
	for strings.HasPrefix(result, "'") && strings.HasSuffix(result, "'") && len(result) > 2 {
		inner := result[1 : len(result)-1]
		if !strings.Contains(inner, "'") || strings.Count(inner, "'")%2 == 0 {
			result = inner
			result = strings.TrimSpace(result)
		} else {
			break
		}
	}

	// 移除句号结尾（对于标题、描述等字段不需要句号）
	if strings.HasSuffix(result, ".") && !strings.Contains(result, ". ") {
		result = strings.TrimSuffix(result, ".")
		result = strings.TrimSpace(result)
	}

	// 移除多余的换行符和空格
	result = strings.ReplaceAll(result, "\n", " ")
	result = strings.ReplaceAll(result, "\r", " ")

	// 合并多个连续空格为单个空格
	for strings.Contains(result, "  ") {
		result = strings.ReplaceAll(result, "  ", " ")
	}

	result = strings.TrimSpace(result)

	return result
}

// parseArticleContent 解析文章内容，分离前置数据和正文
func (a *ArticleTranslator) parseArticleContent(content string) (string, string) {
	lines := strings.Split(content, "\n")

	if len(lines) < 3 || strings.TrimSpace(lines[0]) != "---" {
		return "", content // 没有前置数据
	}

	frontMatterEnd := -1
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			frontMatterEnd = i
			break
		}
	}

	if frontMatterEnd == -1 {
		return "", content // 没有找到前置数据结束标记
	}

	frontMatter := strings.Join(lines[0:frontMatterEnd+1], "\n")
	body := strings.Join(lines[frontMatterEnd+1:], "\n")

	return frontMatter, body
}

// combineTranslatedContent 合并翻译后的内容
func (a *ArticleTranslator) combineTranslatedContent(frontMatter, body string) string {
	if frontMatter == "" {
		return body
	}

	return frontMatter + "\n\n" + body
}

// analyzeArticleContent 分析文章内容统计
func (a *ArticleTranslator) analyzeArticleContent(filePath string) (int, int) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return 0, 0
	}

	_, body := a.parseArticleContent(string(content))

	// 统计字数（粗略估算）
	wordCount := len(strings.Fields(body))

	// 统计段落数
	paragraphs := a.splitIntoParagraphs(body)
	paragraphCount := len(paragraphs)

	return wordCount, paragraphCount
}

// estimateTranslationTime 估算翻译时间
func (a *ArticleTranslator) estimateTranslationTime(paragraphCount int) string {
	// 假设每段落需要2秒翻译时间（包括网络延迟）
	seconds := paragraphCount * 2

	if seconds < 60 {
		return fmt.Sprintf("%d秒", seconds)
	} else if seconds < 3600 {
		minutes := seconds / 60
		return fmt.Sprintf("%d分钟", minutes)
	} else {
		hours := seconds / 3600
		minutes := (seconds % 3600) / 60
		return fmt.Sprintf("%d小时%d分钟", hours, minutes)
	}
}

// splitIntoParagraphs 将文本分割成段落，更细致的处理
func (a *ArticleTranslator) splitIntoParagraphs(text string) []string {
	// 先按双换行符分割
	preliminaryParagraphs := strings.Split(text, "\n\n")

	var finalParagraphs []string

	for _, p := range preliminaryParagraphs {
		trimmed := strings.TrimSpace(p)
		if trimmed == "" {
			continue
		}

		// 进一步处理包含代码块的段落
		if strings.Contains(trimmed, "```") {
			// 代码块保持原样，但检查注释是否包含中文
			finalParagraphs = append(finalParagraphs, trimmed)
		} else {
			// 对于普通段落，按行进一步分割，确保每个有意义的部分都能被翻译
			lines := strings.Split(trimmed, "\n")
			var currentParagraph []string

			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" {
					if len(currentParagraph) > 0 {
						finalParagraphs = append(finalParagraphs, strings.Join(currentParagraph, "\n"))
						currentParagraph = nil
					}
				} else {
					// 检查是否为特殊格式行，但如果包含中文也要翻译
					if a.isSpecialFormatLine(line) && a.containsChinese(line) {
						// 特殊格式但包含中文，单独翻译
						if len(currentParagraph) > 0 {
							finalParagraphs = append(finalParagraphs, strings.Join(currentParagraph, "\n"))
							currentParagraph = nil
						}
						finalParagraphs = append(finalParagraphs, line)
					} else if a.isSpecialFormatLine(line) {
						// 特殊格式且无中文，单独保留
						if len(currentParagraph) > 0 {
							finalParagraphs = append(finalParagraphs, strings.Join(currentParagraph, "\n"))
							currentParagraph = nil
						}
						finalParagraphs = append(finalParagraphs, line)
					} else {
						currentParagraph = append(currentParagraph, line)
					}
				}
			}

			if len(currentParagraph) > 0 {
				finalParagraphs = append(finalParagraphs, strings.Join(currentParagraph, "\n"))
			}
		}
	}

	return finalParagraphs
}

// isSpecialFormatLine 判断是否为特殊格式行
func (a *ArticleTranslator) isSpecialFormatLine(line string) bool {
	trimmed := strings.TrimSpace(line)

	// 标题
	if strings.HasPrefix(trimmed, "#") {
		return true
	}

	// 无序列表
	if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") || strings.HasPrefix(trimmed, "+ ") {
		return true
	}

	// 有序列表
	if matched, _ := regexp.MatchString(`^\d+\. `, trimmed); matched {
		return true
	}

	// 引用
	if strings.HasPrefix(trimmed, ">") {
		return true
	}

	// 水平线
	if trimmed == "---" || trimmed == "***" || trimmed == "___" {
		return true
	}

	return false
}

// containsChinese 检查文本是否包含中文
func (a *ArticleTranslator) containsChinese(text string) bool {
	for _, r := range text {
		if r >= 0x4e00 && r <= 0x9fff {
			return true
		}
	}
	return false
}

// truncateText 截断文本用于显示
func (a *ArticleTranslator) truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen] + "..."
}
