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

// ArticleTranslationPreview 文章翻译预览信息
type ArticleTranslationPreview struct {
	Article      models.Article
	TargetLang   string
	TargetFile   string
	Status       string // "missing", "update", "skip"
	LanguageName string
}

// 实现 StatusLike 接口
func (a ArticleTranslationPreview) GetStatus() string {
	if a.Status == "missing" {
		return "create"
	}
	return "update"
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

// PrepareArticleTranslations 预处理文章翻译
func (a *ArticleTranslator) PrepareArticleTranslations() ([]ArticleTranslationPreview, int, int, error) {
	var previews []ArticleTranslationPreview

	// 测试LM Studio连接
	fmt.Print("🔗 测试LM Studio连接... ")
	if err := a.translationUtils.TestConnection(); err != nil {
		fmt.Printf("❌ 失败 (%v)\n", err)
		fmt.Println("⚠️ 无法连接AI翻译，终止操作")
		return nil, 0, 0, fmt.Errorf("AI翻译连接失败: %v", err)
	} else {
		fmt.Println("✅ 成功")
	}

	// 获取所有文章，使用翻译扫描函数读取完整内容
	articles, err := scanner.ScanArticlesForTranslation(a.contentDir)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("扫描文章失败: %v", err)
	}

	cfg := config.GetGlobalConfig()
	targetLanguages := cfg.Language.TargetLanguages

	var validArticles []models.Article
	for _, article := range articles {
		if article.Title != "" {
			validArticles = append(validArticles, article)
		}
	}

	if len(validArticles) == 0 {
		return previews, 0, 0, nil
	}

	fmt.Printf("📊 正在分析 %d 篇文章的翻译状态...\n", len(validArticles))

	createCount := 0
	updateCount := 0

	for i, article := range validArticles {
		fmt.Printf("  [%d/%d] 检查: %s", i+1, len(validArticles), article.Title)

		articleHasMissing := false
		articleHasExisting := false

		for _, targetLang := range targetLanguages {
			targetFile := utils.BuildTargetFilePath(article.FilePath, targetLang)
			if targetFile == "" {
				continue
			}

			targetLangName := cfg.Language.LanguageNames[targetLang]
			if targetLangName == "" {
				targetLangName = targetLang
			}

			var status string
			if !utils.FileExists(targetFile) {
				status = "missing"
				articleHasMissing = true
			} else {
				status = "update"
				articleHasExisting = true
			}

			preview := ArticleTranslationPreview{
				Article:      article,
				TargetLang:   targetLang,
				TargetFile:   targetFile,
				Status:       status,
				LanguageName: targetLangName,
			}
			previews = append(previews, preview)
		}

		// 统计文章级别的状态
		if articleHasMissing {
			createCount++
		}
		if articleHasExisting {
			updateCount++
		}

		statusText := ""
		if articleHasMissing && articleHasExisting {
			statusText = " 🔄 部分翻译"
		} else if articleHasMissing {
			statusText = " ✨ 需要翻译"
		} else {
			statusText = " ✅ 已完全翻译"
		}
		fmt.Printf("%s\n", statusText)
	}

	fmt.Printf("\n📈 统计结果:\n")
	fmt.Printf("   ✨ 有缺失翻译的文章: %d 篇\n", createCount)
	fmt.Printf("   🔄 已有翻译的文章: %d 篇\n", updateCount)
	fmt.Printf("   📦 总计: %d 篇文章，%d 个翻译任务\n", len(validArticles), len(previews))

	return previews, createCount, updateCount, nil
}

// TranslateArticlesWithMode 根据模式翻译文章
func (a *ArticleTranslator) TranslateArticlesWithMode(targetPreviews []ArticleTranslationPreview, mode string) error {
	fmt.Println("\n📝 文章翻译器 (模式选择)")
	fmt.Println("===============================")

	if len(targetPreviews) == 0 {
		fmt.Printf("ℹ️  根据选择的模式 '%s'，没有需要处理的翻译任务\n", mode)
		return nil
	}

	fmt.Printf("📊 将处理 %d 个翻译任务 (模式: %s)\n", len(targetPreviews), mode)

	return a.processTargetPreviews(targetPreviews)
}

// processTargetPreviews 处理目标预览
func (a *ArticleTranslator) processTargetPreviews(targetPreviews []ArticleTranslationPreview) error {

	utils.LogOperation("开始多语言翻译", map[string]interface{}{
		"translation_tasks": len(targetPreviews),
		"content_dir":       a.contentDir,
	})

	// 1. 统计所有需要翻译的正文总字符数
	totalCharsAllArticles := 0
	for _, preview := range targetPreviews {
		totalCharsAllArticles += preview.Article.CharCount
	}

	globalTranslatedChars := 0
	startTime := time.Now()
	totalSuccessCount := 0
	totalErrorCount := 0

	// 按文章分组处理翻译任务
	articleGroups := a.groupPreviewsByArticle(targetPreviews)

	for i, group := range articleGroups {
		article := group[0].Article
		fmt.Printf("\n📄 处理文章 (%d/%d): %s\n", i+1, len(articleGroups), article.Title)

		articleSuccessCount := 0
		articleErrorCount := 0

		// 统计当前文章剩余语言数
		remainingLangsOfCurrentArticle := len(group)

		// 统计全局剩余文章数
		remainingArticles := len(articleGroups) - i - 1

		for langIndex, preview := range group {
			fmt.Printf("  🌐 翻译为 %s (%d/%d)\n", preview.LanguageName, langIndex+1, len(group))
			fmt.Printf("     目标文件: %s\n", preview.TargetFile)

			if err := a.translateSingleArticleToLanguage(
				preview.Article, preview.TargetFile, preview.TargetLang,
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
	fmt.Printf("- 总成功翻译: %d 个任务\n", totalSuccessCount)
	fmt.Printf("- 总翻译失败: %d 个任务\n", totalErrorCount)

	return nil
}

// groupPreviewsByArticle 按文章分组翻译预览
func (a *ArticleTranslator) groupPreviewsByArticle(previews []ArticleTranslationPreview) [][]ArticleTranslationPreview {
	articleMap := make(map[string][]ArticleTranslationPreview)
	var articleOrder []string

	for _, preview := range previews {
		filePath := preview.Article.FilePath
		if _, exists := articleMap[filePath]; !exists {
			articleOrder = append(articleOrder, filePath)
			articleMap[filePath] = []ArticleTranslationPreview{}
		}
		articleMap[filePath] = append(articleMap[filePath], preview)
	}

	var groups [][]ArticleTranslationPreview
	for _, filePath := range articleOrder {
		groups = append(groups, articleMap[filePath])
	}

	return groups
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
