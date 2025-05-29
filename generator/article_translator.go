package generator

import (
	"fmt"
	"hugo-content-suite/config"
	"hugo-content-suite/models"
	"hugo-content-suite/scanner"
	"hugo-content-suite/utils"
	"os"
)

// ArticleTranslator 文章翻译器
type ArticleTranslator struct {
	contentDir       string
	translationUtils *TranslationUtils
	fileUtils        *FileUtils
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
		translationUtils: NewTranslationUtils(),
		fileUtils:        NewFileUtils(),
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
			targetFile := a.fileUtils.BuildTargetFilePath(article.FilePath, targetLang)
			if targetFile == "" {
				continue
			}

			if a.fileUtils.FileExists(targetFile) {
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

	// 获取所有文章
	articles, err := scanner.ScanArticles(a.contentDir)
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
	fmt.Printf("正在测试与LM Studio的连接...\n")
	if err := a.translationUtils.TestConnection(); err != nil {
		utils.ErrorWithFields("LM Studio连接失败", map[string]interface{}{
			"error": err.Error(),
		})
		return fmt.Errorf("无法连接到LM Studio: %v", err)
	}
	fmt.Printf("LM Studio连接成功！\n")

	return a.processArticlesByLanguage(targetArticles, targetLanguages, mode)
}

// translateSingleArticleToLanguage 翻译单篇文章到指定语言
func (a *ArticleTranslator) translateSingleArticleToLanguage(originalFile, targetFile, targetLang string) error {
	utils.Info("开始翻译文章到 %s: %s", targetLang, originalFile)

	// 读取原文件
	content, err := a.fileUtils.ReadFileContent(originalFile)
	if err != nil {
		utils.Error("读取原文件失败: %s, 错误: %v", originalFile, err)
		return fmt.Errorf("读取原文件失败: %v", err)
	}

	// 解析文章结构
	frontMatter, bodyContent := a.contentParser.ParseArticleContent(content)

	// 翻译前置数据和正文
	translatedFrontMatter, err := a.translateFrontMatterToLanguage(frontMatter, targetLang)
	if err != nil {
		return fmt.Errorf("翻译前置数据失败: %v", err)
	}

	translatedBody, err := a.translateArticleBodyToLanguage(bodyContent, targetLang)
	if err != nil {
		return fmt.Errorf("翻译正文失败: %v", err)
	}

	// 合成并写入最终内容
	finalContent := a.contentParser.CombineTranslatedContent(translatedFrontMatter, translatedBody)
	if err := a.fileUtils.WriteFileContent(targetFile, finalContent); err != nil {
		return fmt.Errorf("写入目标文件失败: %v", err)
	}

	utils.Info("文章翻译完成 (%s): %s", targetLang, targetFile)
	return nil
}

// processArticlesByLanguage 按语言处理文章
func (a *ArticleTranslator) processArticlesByLanguage(targetArticles []models.Article, targetLanguages []string, mode string) error {
	cfg := config.GetGlobalConfig()
	totalSuccessCount := 0
	totalErrorCount := 0

	// 按语言顺序翻译
	for langIndex, targetLang := range targetLanguages {
		targetLangName := cfg.Language.LanguageNames[targetLang]
		if targetLangName == "" {
			targetLangName = targetLang
		}

		fmt.Printf("\n🌐 开始翻译为 %s (%d/%d)\n", targetLangName, langIndex+1, len(targetLanguages))

		successCount := 0
		errorCount := 0

		for i, article := range targetArticles {
			targetFile := a.fileUtils.BuildTargetFilePath(article.FilePath, targetLang)
			if targetFile == "" {
				continue
			}

			// 检查是否需要翻译
			shouldTranslate := a.shouldTranslateArticle(targetFile, mode)
			if !shouldTranslate {
				continue
			}

			fmt.Printf("\n处理文章 (%d/%d): %s\n", i+1, len(targetArticles), article.Title)
			fmt.Printf("目标语言: %s\n", targetLangName)
			fmt.Printf("目标文件: %s\n", targetFile)

			if err := a.translateSingleArticleToLanguage(article.FilePath, targetFile, targetLang); err != nil {
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

// shouldTranslateArticle 判断是否应该翻译文章
func (a *ArticleTranslator) shouldTranslateArticle(targetFile, mode string) bool {
	if _, err := os.Stat(targetFile); os.IsNotExist(err) {
		return mode == "missing" || mode == "all"
	} else {
		return mode == "update" || mode == "all"
	}
}
