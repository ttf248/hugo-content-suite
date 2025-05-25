package generator

import (
	"fmt"
	"hugo-content-suite/config"
	"strings"
	"time"
)

// FieldTranslator 字段翻译器
type FieldTranslator struct {
	translationUtils *TranslationUtils
	contentParser    *ContentParser
}

// NewFieldTranslator 创建字段翻译器
func NewFieldTranslator() *FieldTranslator {
	return &FieldTranslator{
		translationUtils: NewTranslationUtils(),
		contentParser:    NewContentParser(),
	}
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

		// 翻译各种字段
		if translatedLine := a.translateFieldLine(line, targetLang); translatedLine != "" {
			translatedLines = append(translatedLines, translatedLine)
		} else {
			translatedLines = append(translatedLines, line)
		}
	}

	return strings.Join(translatedLines, "\n"), nil
}

// translateFieldLine 翻译字段行
func (a *ArticleTranslator) translateFieldLine(line, targetLang string) string {
	trimmedLine := strings.TrimSpace(line)

	// 处理标题字段
	if strings.HasPrefix(trimmedLine, "title:") {
		return a.translateSingleField(line, "title:", targetLang)
	}

	// 处理描述字段
	if strings.HasPrefix(trimmedLine, "description:") {
		return a.translateSingleField(line, "description:", targetLang)
	}

	// 处理副标题
	if strings.HasPrefix(trimmedLine, "subtitle:") {
		return a.translateSingleField(line, "subtitle:", targetLang)
	}

	// 处理摘要
	if strings.HasPrefix(trimmedLine, "summary:") {
		return a.translateSingleField(line, "summary:", targetLang)
	}

	// 处理slug字段
	if strings.HasPrefix(trimmedLine, "slug:") {
		return a.translateSlugField(line, targetLang)
	}

	// 处理数组字段
	if strings.HasPrefix(trimmedLine, "tags:") {
		return a.translateArrayField(line, "tags:", targetLang)
	}

	if strings.HasPrefix(trimmedLine, "categories:") {
		return a.translateArrayField(line, "categories:", targetLang)
	}

	if strings.HasPrefix(trimmedLine, "authors:") {
		return a.translateArrayField(line, "authors:", targetLang)
	}

	return ""
}

// translateSingleField 翻译单个字段
func (a *ArticleTranslator) translateSingleField(line, prefix, targetLang string) string {
	value := a.contentParser.ExtractFieldValue(line, prefix)
	if value != "" && a.translationUtils.ContainsChinese(value) {
		fmt.Printf("  %s: %s -> ", strings.TrimSuffix(prefix, ":"), value)

		// 使用缓存翻译
		translated, err := a.translationUtils.TranslateToLanguage(value, targetLang)
		if err != nil {
			fmt.Printf("翻译失败\n")
			return ""
		} else {
			translated = a.translationUtils.RemoveQuotes(translated)
			fmt.Printf("%s\n", translated)
			return fmt.Sprintf("%s \"%s\"", prefix, translated)
		}
	}
	return ""
}

// translateSlugField 翻译slug字段
func (a *ArticleTranslator) translateSlugField(line, targetLang string) string {
	slug := a.contentParser.ExtractFieldValue(line, "slug:")
	if slug != "" && a.translationUtils.ContainsChinese(slug) {
		fmt.Printf("  slug: %s -> ", slug)

		// 使用缓存翻译
		translated, err := a.translationUtils.TranslateToLanguage(slug, targetLang)
		if err != nil {
			fmt.Printf("翻译失败\n")
			return ""
		} else {
			translated = a.translationUtils.RemoveQuotes(translated)
			translated = a.translationUtils.FormatSlugField(translated)
			fmt.Printf("%s\n", translated)
			return fmt.Sprintf("slug: \"%s\"", translated)
		}
	}
	return ""
}

// translateArrayField 翻译数组字段
func (a *ArticleTranslator) translateArrayField(line, prefix, targetLang string) string {
	items := a.contentParser.ExtractArrayField(line, prefix)
	if len(items) > 0 {
		translatedItems := a.translateArrayItems(items, strings.TrimSuffix(prefix, ":"), targetLang)
		return fmt.Sprintf("%s %s", prefix, a.contentParser.FormatArrayField(translatedItems))
	}
	return ""
}

// translateArrayItems 翻译数组项目
func (a *ArticleTranslator) translateArrayItems(items []string, fieldType, targetLang string) []string {
	var translated []string

	fmt.Printf("  %s: ", fieldType)

	for _, item := range items {
		if a.translationUtils.ContainsChinese(item) {
			fmt.Printf("%s -> ", item)

			// 使用缓存翻译
			translatedItem, err := a.translationUtils.TranslateToLanguage(item, targetLang)
			if err != nil {
				fmt.Printf("失败 ")
				translated = append(translated, item)
			} else {
				translatedItem = a.translationUtils.RemoveQuotes(translatedItem)
				fmt.Printf("%s ", translatedItem)
				translated = append(translated, translatedItem)
			}
		} else {
			translated = append(translated, item)
		}
	}

	fmt.Printf("\n")
	return translated
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

	totalChars := len(body)
	fmt.Printf("\n翻译正文到 %s (总计 %d 字符)...\n", targetLangName, totalChars)

	return a.translateContentByLinesToLanguage(body, targetLang)
}

// translateContentByLinesToLanguage 按行翻译内容到指定语言
func (a *ArticleTranslator) translateContentByLinesToLanguage(content, targetLang string) (string, error) {
	cfg := config.GetGlobalConfig()
	lines := strings.Split(content, "\n")
	var result []string

	// 翻译统计信息
	totalChars := len(content)
	translatedChars := 0
	startTime := time.Now()

	inCodeBlock := false
	translationCount := 0

	for _, line := range lines {
		// 检测代码块
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			inCodeBlock = !inCodeBlock
			result = append(result, line)
			translatedChars += len(line) + 1 // +1 for newline
			continue
		}

		// 代码块内容直接保留
		if inCodeBlock {
			result = append(result, line)
			translatedChars += len(line) + 1 // +1 for newline
			continue
		}

		// 空行直接保留
		if strings.TrimSpace(line) == "" {
			result = append(result, line)
			translatedChars += len(line) + 1 // +1 for newline
			continue
		}

		// 检查是否包含中文
		if !a.translationUtils.ContainsChinese(line) {
			result = append(result, line)
			translatedChars += len(line) + 1 // +1 for newline
			continue
		}

		// 需要翻译的行
		translationCount++
		lineStartTime := time.Now()

		fmt.Printf("  [%d/%d] 翻译 %d 字符... ", translationCount, len(lines), len(line))

		translatedLine, err := a.translateSingleLineToLanguage(line, translationCount, targetLang)
		if err != nil {
			fmt.Printf("翻译失败\n")
			result = append(result, line)
			translatedChars += len(line) + 1 // +1 for newline
		} else {
			// 计算翻译统计信息
			lineEndTime := time.Now()
			lineDuration := lineEndTime.Sub(lineStartTime)
			translatedChars += len(line) + 1 // +1 for newline

			// 计算总体进度和效率
			totalDuration := lineEndTime.Sub(startTime)
			charsPerSecond := float64(translatedChars) / totalDuration.Seconds()
			charsPer100 := 100.0 / charsPerSecond

			// 计算剩余字符和预估时间
			remainingChars := totalChars - translatedChars
			estimatedRemainingSeconds := float64(remainingChars) / charsPerSecond
			estimatedRemainingTime := time.Duration(estimatedRemainingSeconds) * time.Second

			// 计算进度百分比
			progress := float64(translatedChars) * 100.0 / float64(totalChars)

			fmt.Printf("完成 (%.1fs) | 进度: %.1f%% | 效率: %.1f字符/100s | 预估剩余: %v\n",
				lineDuration.Seconds(),
				progress,
				charsPer100,
				estimatedRemainingTime.Round(time.Second))

			result = append(result, translatedLine)
		}

		// 添加延迟避免API频率限制
		if cfg.Translation.DelayBetweenMs > 0 {
			time.Sleep(time.Duration(cfg.Translation.DelayBetweenMs) * time.Millisecond)
		}
	}

	// 输出最终统计信息
	totalDuration := time.Since(startTime)
	fmt.Printf("\n翻译完成！总计用时: %v | 平均效率: %.1f字符/秒\n",
		totalDuration.Round(time.Second),
		float64(totalChars)/totalDuration.Seconds())

	return strings.Join(result, "\n"), nil
}

// translateSingleLineToLanguage 翻译单行内容到指定语言
func (a *ArticleTranslator) translateSingleLineToLanguage(line string, lineNum int, targetLang string) (string, error) {
	// 使用缓存翻译
	return a.translationUtils.TranslateToLanguage(line, targetLang)
}
