package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"tag-scanner/scanner"
	"tag-scanner/translator"
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

	var previews []ArticleTranslationPreview

	for _, article := range articles {
		if article.Title == "" {
			continue
		}

		// 构建英文文件路径
		dir := filepath.Dir(article.FilePath)
		baseName := filepath.Base(article.FilePath)

		var englishFile string
		if strings.HasSuffix(baseName, ".md") {
			englishFile = filepath.Join(dir, "index.en.md")
		} else {
			continue
		}

		// 检查英文文件是否存在
		status := "missing"
		if _, err := os.Stat(englishFile); err == nil {
			status = "exists"
		}

		// 分析文章内容
		wordCount, paragraphCount := a.analyzeArticleContent(article.FilePath)
		estimatedTime := a.estimateTranslationTime(paragraphCount)

		preview := ArticleTranslationPreview{
			OriginalFile:   article.FilePath,
			EnglishFile:    englishFile,
			Title:          article.Title,
			WordCount:      wordCount,
			ParagraphCount: paragraphCount,
			Status:         status,
			EstimatedTime:  estimatedTime,
		}

		previews = append(previews, preview)
	}

	return previews, nil
}

// TranslateArticles 翻译文章
func (a *ArticleTranslator) TranslateArticles(mode string) error {
	previews, err := a.PreviewArticleTranslations()
	if err != nil {
		return fmt.Errorf("获取翻译预览失败: %v", err)
	}

	// 根据模式过滤文章
	var targetPreviews []ArticleTranslationPreview
	for _, preview := range previews {
		switch mode {
		case "missing":
			if preview.Status == "missing" {
				targetPreviews = append(targetPreviews, preview)
			}
		case "all":
			targetPreviews = append(targetPreviews, preview)
		case "update":
			if preview.Status == "exists" {
				targetPreviews = append(targetPreviews, preview)
			}
		}
	}

	if len(targetPreviews) == 0 {
		fmt.Println("根据选择的模式，没有需要翻译的文章")
		return nil
	}

	// 测试连接
	fmt.Println("正在测试与LM Studio的连接...")
	if err := a.translator.TestConnection(); err != nil {
		return fmt.Errorf("无法连接到LM Studio: %v", err)
	}
	fmt.Println("LM Studio连接成功！")

	successCount := 0
	errorCount := 0

	for i, preview := range targetPreviews {
		fmt.Printf("\n处理文章 (%d/%d): %s\n", i+1, len(targetPreviews), preview.Title)
		fmt.Printf("预计需要时间: %s\n", preview.EstimatedTime)

		if err := a.translateSingleArticle(preview); err != nil {
			fmt.Printf("❌ 翻译失败: %v\n", err)
			errorCount++
		} else {
			fmt.Printf("✅ 翻译完成: %s\n", preview.EnglishFile)
			successCount++
		}
	}

	fmt.Printf("\n文章翻译完成！\n")
	fmt.Printf("- 成功翻译: %d 篇\n", successCount)
	fmt.Printf("- 翻译失败: %d 篇\n", errorCount)
	fmt.Printf("- 总计处理: %d 篇\n", len(targetPreviews))

	return nil
}

// translateSingleArticle 翻译单篇文章
func (a *ArticleTranslator) translateSingleArticle(preview ArticleTranslationPreview) error {
	// 读取原文件
	content, err := os.ReadFile(preview.OriginalFile)
	if err != nil {
		return fmt.Errorf("读取原文件失败: %v", err)
	}

	// 解析文章结构
	frontMatter, bodyContent := a.parseArticleContent(string(content))

	// 翻译前置数据
	translatedFrontMatter, err := a.translateFrontMatter(frontMatter)
	if err != nil {
		return fmt.Errorf("翻译前置数据失败: %v", err)
	}

	// 分段翻译正文
	translatedBody, err := a.translateArticleBody(bodyContent)
	if err != nil {
		return fmt.Errorf("翻译正文失败: %v", err)
	}

	// 合成最终内容
	finalContent := a.combineTranslatedContent(translatedFrontMatter, translatedBody)

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(preview.EnglishFile), 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	// 写入英文文件
	if err := os.WriteFile(preview.EnglishFile, []byte(finalContent), 0644); err != nil {
		return fmt.Errorf("写入英文文件失败: %v", err)
	}

	return nil
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

// translateFrontMatter 翻译前置数据
func (a *ArticleTranslator) translateFrontMatter(frontMatter string) (string, error) {
	if frontMatter == "" {
		return "", nil
	}

	lines := strings.Split(frontMatter, "\n")
	var translatedLines []string

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		if trimmedLine == "---" {
			translatedLines = append(translatedLines, line)
			continue
		}

		// 翻译标题
		if strings.HasPrefix(trimmedLine, "title:") {
			title := strings.TrimSpace(strings.TrimPrefix(trimmedLine, "title:"))
			title = strings.Trim(title, "\"'")

			if title != "" {
				translatedTitle, err := a.translator.TranslateToArticleSlug(title)
				if err != nil {
					fmt.Printf("⚠️ 标题翻译失败，保持原文: %v\n", err)
					translatedLines = append(translatedLines, line)
				} else {
					translatedLines = append(translatedLines, fmt.Sprintf("title: \"%s\"", translatedTitle))
				}
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

// translateArticleBody 分段翻译正文
func (a *ArticleTranslator) translateArticleBody(body string) (string, error) {
	if strings.TrimSpace(body) == "" {
		return body, nil
	}

	// 按段落分割内容
	paragraphs := a.splitIntoParagraphs(body)
	var translatedParagraphs []string

	for i, paragraph := range paragraphs {
		if strings.TrimSpace(paragraph) == "" {
			translatedParagraphs = append(translatedParagraphs, paragraph)
			continue
		}

		fmt.Printf("  翻译段落 (%d/%d)...\n", i+1, len(paragraphs))

		// 翻译当前段落
		translatedParagraph, err := a.translator.TranslateParagraph(paragraph)
		if err != nil {
			fmt.Printf("⚠️ 段落翻译失败，保持原文: %v\n", err)
			translatedParagraphs = append(translatedParagraphs, paragraph)
		} else {
			translatedParagraphs = append(translatedParagraphs, translatedParagraph)
		}

		// 添加延迟避免API频率限制
		time.Sleep(time.Millisecond * 500)
	}

	return strings.Join(translatedParagraphs, "\n\n"), nil
}

// splitIntoParagraphs 将文本分割成段落
func (a *ArticleTranslator) splitIntoParagraphs(text string) []string {
	// 按双换行符分割段落
	paragraphs := strings.Split(text, "\n\n")

	var result []string
	for _, paragraph := range paragraphs {
		trimmed := strings.TrimSpace(paragraph)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
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
