package generator

import (
	"fmt"
	"hugo-content-suite/config"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
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
	if strings.TrimSpace(frontMatter) == "" {
		return frontMatter, nil
	}

	cfg := config.GetGlobalConfig()
	targetLangName := cfg.Language.LanguageNames[targetLang]
	if targetLangName == "" {
		targetLangName = targetLang
	}

	fmt.Printf("翻译前置数据到 %s...\n", targetLangName)

	// 解析 YAML
	var frontMatterData map[string]interface{}
	if err := yaml.Unmarshal([]byte(frontMatter), &frontMatterData); err != nil {
		return "", fmt.Errorf("解析前置数据失败: %v", err)
	}

	// 翻译各个字段
	translatedData, err := a.translateFrontMatterFields(frontMatterData, targetLang)
	if err != nil {
		return "", fmt.Errorf("翻译前置数据字段失败: %v", err)
	}

	// 将翻译后的数据转换回 YAML
	translatedYAML, err := yaml.Marshal(translatedData)
	if err != nil {
		return "", fmt.Errorf("生成翻译后的YAML失败: %v", err)
	}

	return "---\r\n" + string(translatedYAML) + "---\r\n", nil
}

// translateFrontMatterFields 翻译前置数据的所有字段
func (a *ArticleTranslator) translateFrontMatterFields(data map[string]interface{}, targetLang string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// 定义需要翻译的字段
	translatableFields := map[string]bool{
		"title":       true,
		"description": true,
		"subtitle":    true,
		"summary":     true,
	}

	// 定义需要翻译的数组字段
	translatableArrayFields := map[string]bool{
		"tags":       true,
		"categories": true,
		"authors":    true,
	}

	for key, value := range data {
		switch {
		case translatableFields[key]:
			// 翻译单个字符串字段
			if strValue, ok := value.(string); ok {
				translatedValue, err := a.translateStringField(key, strValue, targetLang)
				if err != nil {
					fmt.Printf("  警告: 翻译字段 %s 失败: %v\n", key, err)
					result[key] = value // 保持原值
				} else {
					result[key] = translatedValue
				}
			} else {
				result[key] = value
			}

		case translatableArrayFields[key]:
			// 翻译数组字段
			if arrayValue, ok := value.([]interface{}); ok {
				translatedArray, err := a.translateArrayField(key, arrayValue, targetLang)
				if err != nil {
					fmt.Printf("  警告: 翻译数组字段 %s 失败: %v\n", key, err)
					result[key] = value // 保持原值
				} else {
					result[key] = translatedArray
				}
			} else {
				result[key] = value
			}

		case key == "slug":
			// 特殊处理 slug 字段
			if strValue, ok := value.(string); ok {
				translatedSlug, err := a.translateSlugField(strValue, targetLang)
				if err != nil {
					fmt.Printf("  警告: 翻译slug失败: %v\n", err)
					result[key] = value // 保持原值
				} else {
					result[key] = translatedSlug
				}
			} else {
				result[key] = value
			}

		default:
			// 其他字段保持不变
			result[key] = value
		}
	}

	return result, nil
}

// translateStringField 翻译字符串字段
func (a *ArticleTranslator) translateStringField(fieldName, value, targetLang string) (string, error) {
	if value == "" || !a.translationUtils.ContainsChinese(value) {
		return value, nil
	}

	fmt.Printf("  %s: %s -> ", fieldName, value)

	// 使用缓存翻译
	translated, err := a.translationUtils.TranslateToLanguage(value, targetLang)
	if err != nil {
		fmt.Printf("翻译失败\n")
		return value, err
	}

	// 清理翻译结果
	translated = a.translationUtils.RemoveQuotes(translated)
	translated = a.translationUtils.CleanTranslationResult(translated)
	translated = strings.ReplaceAll(translated, "\"", "")
	translated = strings.ReplaceAll(translated, "'", "")

	fmt.Printf("%s\n", translated)
	return translated, nil
}

// translateArrayField 翻译数组字段
func (a *ArticleTranslator) translateArrayField(fieldName string, items []interface{}, targetLang string) ([]interface{}, error) {
	if len(items) == 0 {
		return items, nil
	}

	fmt.Printf("  %s: ", fieldName)

	var translatedItems []interface{}
	for _, item := range items {
		if strItem, ok := item.(string); ok {
			if a.translationUtils.ContainsChinese(strItem) {
				fmt.Printf("%s -> ", strItem)

				// 使用缓存翻译
				translated, err := a.translationUtils.TranslateToLanguage(strItem, targetLang)
				if err != nil {
					fmt.Printf("失败 ")
					translatedItems = append(translatedItems, item)
					continue
				}

				// 清理翻译结果
				translated = a.translationUtils.RemoveQuotes(translated)
				translated = a.translationUtils.CleanTranslationResult(translated)
				translated = strings.ReplaceAll(translated, "\"", "")
				translated = strings.ReplaceAll(translated, "'", "")

				fmt.Printf("%s ", translated)
				translatedItems = append(translatedItems, translated)
			} else {
				fmt.Printf("%s -> %s\t", strItem, strItem)
				translatedItems = append(translatedItems, item)
			}
		} else {
			translatedItems = append(translatedItems, item)
		}
	}

	fmt.Printf("\n")
	return translatedItems, nil
}

// translateSlugField 翻译slug字段
func (a *ArticleTranslator) translateSlugField(slug, targetLang string) (string, error) {
	if slug == "" || !a.translationUtils.ContainsChinese(slug) {
		return slug, nil
	}

	fmt.Printf("  slug: %s -> ", slug)

	// 使用缓存翻译
	translated, err := a.translationUtils.TranslateToLanguage(slug, targetLang)
	if err != nil {
		fmt.Printf("翻译失败\n")
		return slug, err
	}

	// 清理翻译结果
	translated = a.translationUtils.RemoveQuotes(translated)
	translated = a.translationUtils.CleanTranslationResult(translated)
	translated = strings.ReplaceAll(translated, "\"", "")
	translated = strings.ReplaceAll(translated, "'", "")
	translated = a.translationUtils.FormatSlugField(translated)

	fmt.Printf("%s\n", translated)
	return translated, nil
}

// translateArticleBodyToLanguage 翻译正文到指定语言（段落级，支持拆分合并）
func (a *ArticleTranslator) translateArticleBodyToLanguage(body, targetLang string) (string, error) {
	if strings.TrimSpace(body) == "" {
		return body, nil
	}

	cfg := config.GetGlobalConfig()
	targetLangName := cfg.Language.LanguageNames[targetLang]
	if targetLangName == "" {
		targetLangName = targetLang
	}

	fmt.Printf("\n翻译正文到 %s...\n", targetLangName)

	// 解析为段落并获取映射关系
	splitResult, err := a.contentParser.ParseContentIntoParagraphsWithMapping(body)
	if err != nil {
		return "", fmt.Errorf("解析段落失败: %v", err)
	}

	paragraphs := splitResult.Paragraphs
	mappings := splitResult.Mappings
	totalParagraphs := len(paragraphs)
	translatableParagraphs := a.contentParser.CountTranslatableParagraphs(paragraphs)

	fmt.Printf("📖 总段落数: %d | 需翻译: %d | 跳过: %d\n",
		totalParagraphs, translatableParagraphs, totalParagraphs-translatableParagraphs)

	// 翻译段落
	translatedParagraphs, err := a.translateParagraphsToLanguageWithMapping(paragraphs, targetLang)
	if err != nil {
		return "", err
	}

	// 如果启用了合并功能，则合并拆分后的段落
	if cfg.Paragraph.MergeAfterTranslation {
		fmt.Printf("🔄 合并拆分的段落...\n")
		mergedParagraphs, err := a.contentParser.MergeTranslatedParagraphs(translatedParagraphs, mappings)
		if err != nil {
			fmt.Printf("⚠️ 段落合并失败，使用原始翻译结果: %v\n", err)
			return strings.Join(translatedParagraphs, "\n\n"), nil
		}

		fmt.Printf("✅ 段落合并完成: %d个翻译段落 → %d个合并段落\n",
			len(translatedParagraphs), len(mergedParagraphs))
		return strings.Join(mergedParagraphs, "\n\n"), nil
	}

	return strings.Join(translatedParagraphs, "\n\n"), nil
}

// translateParagraphsToLanguageWithMapping 翻译段落列表到指定语言（支持映射关系）
func (a *ArticleTranslator) translateParagraphsToLanguageWithMapping(paragraphs []string, targetLang string) ([]string, error) {
	cfg := config.GetGlobalConfig()
	var translatedParagraphs []string

	// 统计信息
	totalParagraphs := len(paragraphs)
	translatableParagraphs := a.contentParser.CountTranslatableParagraphs(paragraphs)
	translatedCount := 0
	successCount := 0
	errorCount := 0
	startTime := time.Now()

	fmt.Printf("\n开始段落级翻译...\n")

	for _, paragraph := range paragraphs {
		trimmed := strings.TrimSpace(paragraph)

		// 检查是否需要翻译
		if !a.contentParser.needsTranslation(paragraph) {
			translatedParagraphs = append(translatedParagraphs, paragraph)
			continue
		}

		translatedCount++

		// 生成进度信息
		progressPercent := float64(translatedCount) * 100.0 / float64(translatableParagraphs)
		progressBar := a.generateProgressBar(translatedCount, translatableParagraphs, 30)

		// 计算效率和预估时间
		elapsed := time.Since(startTime)
		avgTimePerParagraph := float64(elapsed.Nanoseconds()) / float64(translatedCount) / 1e9
		remainingParagraphs := translatableParagraphs - translatedCount
		estimatedRemaining := time.Duration(float64(remainingParagraphs) * avgTimePerParagraph * 1e9)

		fmt.Printf("\n📝 段落 %d/%d %s %.1f%%\n",
			translatedCount, translatableParagraphs, progressBar, progressPercent)
		fmt.Printf("📄 长度: %d 字符 | 预计剩余: %v\n",
			len(trimmed), estimatedRemaining.Round(time.Second))

		// 显示段落预览（前80字符）
		preview := trimmed
		if len(preview) > 200 {
			preview = preview[:200] + "..."
		}
		fmt.Printf("📖 内容: %s\n", preview)

		// 翻译段落
		paragraphStartTime := time.Now()
		translatedParagraph, err := a.translationUtils.TranslateParagraphToLanguage(paragraph, targetLang)
		paragraphDuration := time.Since(paragraphStartTime)

		if err != nil {
			fmt.Printf("❌ 翻译失败 (%.1fs): %v\n", paragraphDuration.Seconds(), err)
			fmt.Printf("📝 保留原文\n")
			translatedParagraphs = append(translatedParagraphs, paragraph)
			errorCount++
		} else {
			fmt.Printf("✅ 翻译完成 (%.1fs)\n", paragraphDuration.Seconds())
			translatedParagraphs = append(translatedParagraphs, translatedParagraph)
			successCount++

			// 显示翻译结果预览
			translatedPreview := strings.TrimSpace(translatedParagraph)
			if len(translatedPreview) > 200 {
				translatedPreview = translatedPreview[:200] + "..."
			}
			fmt.Printf("📝 译文: %s\n", translatedPreview)
		}

		// 添加延迟避免API频率限制
		if cfg.Translation.DelayBetweenMs > 0 && translatedCount < translatableParagraphs {
			time.Sleep(time.Duration(cfg.Translation.DelayBetweenMs) * time.Millisecond)
		}

		// 每10个段落输出阶段报告
		if translatedCount%10 == 0 {
			a.printParagraphStageReport(translatedCount, translatableParagraphs, elapsed, successCount, errorCount)
		}
	}

	// 输出最终统计
	totalDuration := time.Since(startTime)
	successRate := float64(successCount) * 100.0 / float64(translatedCount)
	avgParagraphTime := totalDuration.Seconds() / float64(translatedCount)

	fmt.Printf("\n🎉 段落翻译完成！\n")
	fmt.Printf("   ⏱️  总用时: %v\n", totalDuration.Round(time.Second))
	fmt.Printf("   📊 成功率: %.1f%% (%d/%d)\n", successRate, successCount, translatedCount)
	fmt.Printf("   ⚡ 平均速度: %.1f 秒/段落\n", avgParagraphTime)
	fmt.Printf("   📖 处理: %d 段落 (翻译 %d | 跳过 %d)\n",
		totalParagraphs, translatedCount, totalParagraphs-translatedCount)

	return translatedParagraphs, nil
}

// printParagraphStageReport 打印段落翻译阶段报告
func (a *ArticleTranslator) printParagraphStageReport(current, total int, elapsed time.Duration, success, error int) {
	stage := (current + 9) / 10
	successRate := float64(success) * 100.0 / float64(current)
	avgTime := elapsed.Seconds() / float64(current)

	fmt.Printf("\n🏁 阶段 %d 完成 | 已翻译 %d/%d 段落\n", stage, current, total)
	fmt.Printf("   ⏱️  阶段用时: %v | 平均: %.1f 秒/段落\n",
		elapsed.Round(time.Second), avgTime)
	fmt.Printf("   📊 成功率: %.1f%% (%d 成功, %d 失败)\n", successRate, success, error)
	fmt.Printf("   ▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔\n")
}

// generateProgressBar 生成进度条
func (a *ArticleTranslator) generateProgressBar(current, total, width int) string {
	if total == 0 {
		return strings.Repeat("░", width)
	}

	progress := float64(current) / float64(total)
	completed := int(progress * float64(width))

	bar := strings.Repeat("█", completed)
	remaining := strings.Repeat("░", width-completed)

	return fmt.Sprintf("[%s%s]", bar, remaining)
}

// 移除不再需要的行级翻译相关方法
// translateContentByLinesToLanguage, translateMarkdownAwareLine 等方法可以删除或标记为废弃
