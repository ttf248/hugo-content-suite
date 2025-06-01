package generator

import (
	"fmt"
	"hugo-content-suite/config"
	"hugo-content-suite/models"
	"hugo-content-suite/scanner"
	"hugo-content-suite/translator"
	"hugo-content-suite/utils"
	"os"
	"strings"
	"time"
)

// ArticleTranslator 文章翻译器
type ArticleTranslator struct {
	contentDir       string
	translationUtils *translator.TranslationUtils
	contentParser    *ContentParser
}

// TranslationStatus 翻译状态信息
type TranslationStatus struct {
	MissingArticles  int // 有缺失翻译的文章数
	ExistingArticles int // 所有翻译都存在的文章数
	TotalArticles    int // 文章总数
}

// NewArticleTranslator 创建新的文章翻译器
func NewArticleTranslator(contentDir string) *ArticleTranslator {
	return &ArticleTranslator{
		contentDir:       contentDir,
		translationUtils: translator.NewTranslationUtils(),
		contentParser:    NewContentParser(),
	}
}

// GetTranslationStatus 获取翻译状态统计
func (a *ArticleTranslator) GetTranslationStatus() (*TranslationStatus, error) {
	articles, err := scanner.ScanArticles(a.contentDir)
	if err != nil {
		return nil, fmt.Errorf("扫描文章失败: %v", err)
	}

	cfg := config.GetGlobalConfig()
	targetLanguages := cfg.Language.TargetLanguages

	missingCount := 0
	existingCount := 0
	totalArticles := 0

	for _, article := range articles {
		if article.Title == "" {
			continue
		}
		totalArticles++

		hasMissing := false
		hasExisting := false

		// 检查每种目标语言的翻译状态
		for _, targetLang := range targetLanguages {
			targetFile := utils.BuildTargetFilePath(article.FilePath, targetLang)
			if targetFile == "" {
				continue
			}

			if utils.FileExists(targetFile) {
				hasExisting = true
			} else {
				hasMissing = true
			}
		}

		// 如果有任何语言缺失翻译，则算作需要翻译的文章
		if hasMissing {
			missingCount++
		} else if hasExisting {
			// 只有当所有语言都存在时，才算作已翻译的文章
			existingCount++
		}
	}

	return &TranslationStatus{
		MissingArticles:  missingCount,
		ExistingArticles: existingCount,
		TotalArticles:    totalArticles,
	}, nil
}

// TranslateArticles 翻译文章到多种语言
func (a *ArticleTranslator) TranslateArticles(mode string) error {
	cfg := config.GetGlobalConfig()
	targetLanguages := cfg.Language.TargetLanguages

	utils.LogOperation("开始多语言翻译", map[string]interface{}{
		"mode":             mode,
		"target_languages": targetLanguages,
		"content_dir":      a.contentDir,
	})

	// 获取所有文章，使用新的扫描函数读取完整内容
	articles, err := scanner.ScanArticlesForTranslation(a.contentDir)
	if err != nil {
		utils.ErrorWithFields("扫描文章失败", map[string]interface{}{
			"content_dir": a.contentDir,
			"error":       err.Error(),
		})
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
	if err := a.translationUtils.TestConnection(); err != nil {
		utils.ErrorWithFields("LM Studio连接失败", map[string]interface{}{
			"error": err.Error(),
		})
		return fmt.Errorf("无法连接到LM Studio: %v", err)
	}
	fmt.Printf("LM Studio连接成功！\n")

	return a.processArticlesByLanguage(targetArticles, targetLanguages, mode)
}

// processArticlesByLanguage 按语言处理文章
func (a *ArticleTranslator) processArticlesByLanguage(targetArticles []models.Article, targetLanguages []string, mode string) error {
	totalSuccessCount := 0
	totalErrorCount := 0

	// 1. 统计所有需要翻译的正文总字符数 - 直接使用缓存的字符数
	totalCharsAllArticles := 0
	for _, article := range targetArticles {
		for _, targetLang := range targetLanguages {
			targetFile := utils.BuildTargetFilePath(article.FilePath, targetLang)
			if targetFile == "" {
				continue
			}
			shouldTranslate := a.shouldTranslateArticle(targetFile, mode)
			if shouldTranslate {
				totalCharsAllArticles += article.CharCount
			}
		}
	}

	globalTranslatedChars := 0
	startTime := time.Now()

	// 按文章顺序翻译，每篇文章完成所有语言后再处理下一篇
	for i, article := range targetArticles {
		fmt.Printf("\n📄 处理文章 (%d/%d): %s\n", i+1, len(targetArticles), article.Title)

		articleSuccessCount := 0
		articleErrorCount := 0

		// 统计当前文章剩余语言数
		remainingLangsOfCurrentArticle := 0
		for _, targetLang := range targetLanguages {
			targetFile := utils.BuildTargetFilePath(article.FilePath, targetLang)
			if targetFile == "" {
				continue
			}
			shouldTranslate := a.shouldTranslateArticle(targetFile, mode)
			if shouldTranslate {
				remainingLangsOfCurrentArticle++
			}
		}

		cfg := config.GetGlobalConfig()
		for langIndex, targetLang := range targetLanguages {
			targetLangName := cfg.Language.LanguageNames[targetLang]
			if targetLangName == "" {
				targetLangName = targetLang
			}

			targetFile := utils.BuildTargetFilePath(article.FilePath, targetLang)
			if targetFile == "" {
				continue
			}
			shouldTranslate := a.shouldTranslateArticle(targetFile, mode)
			if !shouldTranslate {
				fmt.Printf("  ⏭️  跳过 %s (已存在)\n", targetLangName)
				continue
			}

			// 统计全局剩余文章数
			remainingArticles := 0
			for j := i + 1; j < len(targetArticles); j++ {
				for _, tl := range targetLanguages {
					tf := utils.BuildTargetFilePath(targetArticles[j].FilePath, tl)
					if tf == "" {
						continue
					}
					if a.shouldTranslateArticle(tf, mode) {
						remainingArticles++
						break
					}
				}
			}

			fmt.Printf("  🌐 翻译为 %s (%d/%d)\n", targetLangName, langIndex+1, len(targetLanguages))
			fmt.Printf("     目标文件: %s\n", targetFile)

			if err := a.translateSingleArticleToLanguage(
				article, targetFile, targetLang,
				totalCharsAllArticles, &globalTranslatedChars, startTime,
				remainingArticles, remainingLangsOfCurrentArticle-1,
			); err != nil {
				fmt.Printf("     ❌ 翻译失败: %v\n", err)
				articleErrorCount++
				totalErrorCount++
			} else {
				fmt.Printf("     ✅ 翻译完成\n")
				articleSuccessCount++
				totalSuccessCount++
			}
			remainingLangsOfCurrentArticle--
		}

		fmt.Printf("  📊 当前文章翻译结果: 成功 %d, 失败 %d\n", articleSuccessCount, articleErrorCount)
	}

	fmt.Printf("\n🎉 多语言翻译全部完成！\n")
	fmt.Printf("- 目标语言: %v\n", targetLanguages)
	fmt.Printf("- 总成功翻译: %d 篇\n", totalSuccessCount)
	fmt.Printf("- 总翻译失败: %d 篇\n", totalErrorCount)

	return nil
}

// translateSingleArticleToLanguage 翻译单篇文章到指定语言
func (a *ArticleTranslator) translateSingleArticleToLanguage(
	article models.Article, targetFile, targetLang string,
	totalCharsAllArticles int, globalTranslatedChars *int, globalStartTime time.Time,
	remainingArticles int, remainingLangsOfCurrentArticle int,
) error {
	utils.Info("开始翻译文章到 %s: %s", targetLang, article.FilePath)

	// 直接使用缓存的前置信息和正文内容
	frontMatter := article.FrontMatter
	bodyParagraphs := article.BodyContent

	// 翻译前置数据和正文
	translatedFrontMatter, err := a.translateFrontMatterToLanguage(frontMatter, targetLang)
	if err != nil {
		fmt.Printf("⚠️ 翻译前置数据失败: %v\n", err)
		return fmt.Errorf("翻译前置数据失败: %v", err)
	}

	translatedBody, err := a.translateArticleBodyParagraphsWithProgress(
		bodyParagraphs, targetLang, totalCharsAllArticles, globalTranslatedChars, globalStartTime,
		remainingArticles, remainingLangsOfCurrentArticle,
	)
	if err != nil {
		return fmt.Errorf("翻译正文失败: %v", err)
	}

	// 合成并写入最终内容
	finalContent := a.contentParser.CombineTranslatedContent(translatedFrontMatter, translatedBody)
	if err := utils.WriteFileContent(targetFile, finalContent); err != nil {
		return fmt.Errorf("写入目标文件失败: %v", err)
	}

	utils.Info("文章翻译完成 (%s): %s", targetLang, targetFile)
	return nil
}

// translateArticleBodyParagraphsWithProgress 翻译段落数组
func (a *ArticleTranslator) translateArticleBodyParagraphsWithProgress(
	paragraphs []string, targetLang string,
	totalCharsAllArticles int, globalTranslatedChars *int, globalStartTime time.Time,
	remainingArticles int, remainingLangsOfCurrentArticle int,
) (string, error) {
	if len(paragraphs) == 0 {
		return "", nil
	}

	cfg := config.GetGlobalConfig()
	targetLangName := cfg.Language.LanguageNames[targetLang]
	if targetLangName == "" {
		targetLangName = targetLang
	}

	fmt.Printf("\n翻译正文到 %s...\n", targetLangName)

	// 应用段落拆分并获取映射关系
	splitResult := a.contentParser.applySplittingWithMapping(paragraphs)

	splitParagraphs := splitResult.Paragraphs
	mappings := splitResult.Mappings
	totalParagraphs := len(splitParagraphs)
	translatableParagraphs := len(splitParagraphs)

	// 统计总字符数
	totalChars := 0
	for _, p := range paragraphs {
		totalChars += len([]rune(p))
	}

	fmt.Printf("📖 总段落数: %d | 需翻译: %d | 跳过: %d\n",
		totalParagraphs, translatableParagraphs, totalParagraphs-translatableParagraphs)
	fmt.Printf("🔢 总字符数: %d\n", totalChars)

	// 翻译段落，传递全局进度参数
	translatedParagraphs, err := a.translateParagraphsToLanguageWithMappingAndGlobalProgress(
		splitParagraphs, targetLang, totalChars, totalCharsAllArticles, globalTranslatedChars, globalStartTime,
		remainingArticles, remainingLangsOfCurrentArticle,
	)
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

// 新增：带全局进度的段落翻译
func (a *ArticleTranslator) translateParagraphsToLanguageWithMappingAndGlobalProgress(
	paragraphs []string, targetLang string, totalChars int, totalCharsAllArticles int,
	globalTranslatedChars *int, globalStartTime time.Time,
	remainingArticles int, remainingLangsOfCurrentArticle int,
) ([]string, error) {
	cfg := config.GetGlobalConfig()
	var translatedParagraphs []string

	// 统计信息
	totalParagraphs := len(paragraphs)
	translatableParagraphs := len(paragraphs)
	translatedCount := 0
	successCount := 0
	errorCount := 0
	startTime := time.Now()

	// 新增：累计已翻译字符数
	translatedChars := 0

	for _, paragraph := range paragraphs {
		trimmed := strings.TrimSpace(paragraph)
		paraLen := len([]rune(trimmed))

		translatedCount++
		translatedChars += paraLen
		if globalTranslatedChars != nil {
			*globalTranslatedChars += paraLen
		}

		// 仅每N个段落输出一次进度，减少刷屏
		const progressStep = 1
		showProgress := translatedCount == 1 || translatedCount == translatableParagraphs || translatedCount%progressStep == 0
		if showProgress {
			// 进度信息
			progressPercent := float64(translatedCount) * 100.0 / float64(translatableParagraphs)

			// 文章级进度（按字符数）
			charProgressPercent := 0.0
			if totalChars > 0 {
				charProgressPercent = float64(translatedChars) * 100.0 / float64(totalChars)
			}
			avgTimePerChar := 0.0
			elapsed := time.Since(startTime)
			if translatedChars > 0 {
				avgTimePerChar = elapsed.Seconds() / float64(translatedChars)
			}
			remainingChars := totalChars - translatedChars
			estimatedCharRemaining := time.Duration(float64(remainingChars) * avgTimePerChar * float64(time.Second))

			// 全局进度
			globalProgressLine := ""
			if globalTranslatedChars != nil && totalCharsAllArticles > 0 {
				globalPercent := float64(*globalTranslatedChars) * 100.0 / float64(totalCharsAllArticles)
				globalElapsed := time.Since(globalStartTime)
				globalAvgTimePerChar := globalElapsed.Seconds() / float64(*globalTranslatedChars)
				globalRemainingChars := totalCharsAllArticles - *globalTranslatedChars
				globalEstimatedRemaining := time.Duration(float64(globalRemainingChars) * globalAvgTimePerChar * float64(time.Second))
				globalProgressLine = fmt.Sprintf(
					"\n🌏 总进度: %d/%d 字符 (%.1f%%) | 剩余文章: %d | 当前文章剩余语言: %d | 总用时: %v | 预计剩余: %v\n",
					*globalTranslatedChars, totalCharsAllArticles, globalPercent,
					remainingArticles, remainingLangsOfCurrentArticle,
					globalElapsed.Round(time.Second), globalEstimatedRemaining.Round(time.Second))
			}

			// 先打印总进度，再打印全局进度
			if globalProgressLine != "" {
				fmt.Print(globalProgressLine)
			}
			fmt.Printf("\n📊 文章进度: %d/%d 字符 (%.1f%%) | 段落 %d/%d %.1f%% | 预计剩余: %v\n",
				translatedChars, totalChars, charProgressPercent,
				translatedCount, translatableParagraphs, progressPercent,
				estimatedCharRemaining.Round(time.Second))
		}

		preview := trimmed
		if len(preview) > 80 {
			preview = preview[:80] + "..."
		}
		fmt.Printf("📖 内容: %s\n", preview)

		// 翻译段落
		paragraphStartTime := time.Now()
		translatedParagraph, err := a.translationUtils.TranslateToLanguage(paragraph, targetLang)
		paragraphDuration := time.Since(paragraphStartTime)

		if err != nil {
			fmt.Printf("❌ 翻译失败 (%.1fs): %v\n", paragraphDuration.Seconds(), err)
			fmt.Printf("📝 保留原文\n")
			translatedParagraphs = append(translatedParagraphs, paragraph)
			errorCount++
		} else {
			translatedPreview := strings.TrimSpace(translatedParagraph)
			if len(translatedPreview) > 80 {
				translatedPreview = translatedPreview[:80] + "..."
			}
			fmt.Printf("📝 译文: %s\n", translatedPreview)
			translatedParagraphs = append(translatedParagraphs, translatedParagraph)
			successCount++
		}

		// 添加延迟避免API频率限制
		if cfg.Translation.DelayBetweenMs > 0 && translatedCount < translatableParagraphs {
			time.Sleep(time.Duration(cfg.Translation.DelayBetweenMs) * time.Millisecond)
		}

		// 每10个段落输出阶段报告
		if translatedCount%10 == 0 {
			elapsed := time.Since(startTime)
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

// shouldTranslateArticle 判断是否应该翻译文章
func (a *ArticleTranslator) shouldTranslateArticle(targetFile, mode string) bool {
	if _, err := os.Stat(targetFile); os.IsNotExist(err) {
		return mode == "missing" || mode == "all"
	} else {
		return mode == "update" || mode == "all"
	}
}
