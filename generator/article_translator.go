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
			title := a.extractFieldValue(trimmedLine, "title:")
			if title != "" && a.containsChinese(title) {
				translatedTitle, err := a.translateFieldContent(title)
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

		// 翻译描述字段
		if strings.HasPrefix(trimmedLine, "description:") {
			description := a.extractFieldValue(trimmedLine, "description:")
			if description != "" && a.containsChinese(description) {
				translatedDescription, err := a.translateFieldContent(description)
				if err != nil {
					fmt.Printf("⚠️ 描述翻译失败，保持原文: %v\n", err)
					translatedLines = append(translatedLines, line)
				} else {
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
				translatedSubtitle, err := a.translateFieldContent(subtitle)
				if err != nil {
					fmt.Printf("⚠️ 副标题翻译失败，保持原文: %v\n", err)
					translatedLines = append(translatedLines, line)
				} else {
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
				translatedSummary, err := a.translateFieldContent(summary)
				if err != nil {
					fmt.Printf("⚠️ 摘要翻译失败，保持原文: %v\n", err)
					translatedLines = append(translatedLines, line)
				} else {
					translatedLines = append(translatedLines, fmt.Sprintf("summary: \"%s\"", translatedSummary))
				}
			} else {
				translatedLines = append(translatedLines, line)
			}
			continue
		}

		// 翻译分类数组
		if strings.HasPrefix(trimmedLine, "categories:") {
			categories := a.extractArrayField(trimmedLine, "categories:")
			if len(categories) > 0 {
				translatedCategories := a.translateArrayField(categories, "分类")
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
				// 作者名通常不翻译，但如果是中文描述可以翻译
				translatedAuthors := a.translateArrayField(authors, "作者")
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

// extractFieldValue 提取字段值
func (a *ArticleTranslator) extractFieldValue(line, prefix string) string {
	value := strings.TrimSpace(strings.TrimPrefix(line, prefix))
	value = strings.Trim(value, "\"'")
	return value
}

// translateFieldContent 翻译字段内容，使用简化的提示词
func (a *ArticleTranslator) translateFieldContent(content string) (string, error) {
	prompt := fmt.Sprintf(`Translate the following Chinese text to English. Keep it concise and natural:

%s`, content)

	request := translator.LMStudioRequest{
		Model: "gemma-3-12b-it",
		Messages: []translator.Message{
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

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Post("http://172.19.192.1:2234/v1/chat/completions", "application/json", bytes.NewBuffer(jsonData))
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

	// 清理可能的额外内容
	result = a.cleanTranslationResult(result)

	return result, nil
}

// cleanTranslationResult 清理翻译结果，移除多余的提示词或格式
func (a *ArticleTranslator) cleanTranslationResult(result string) string {
	// 移除常见的提示词残留
	patterns := []string{
		"Translation:",
		"Translated:",
		"English:",
		"Result:",
		"Output:",
	}

	for _, pattern := range patterns {
		if strings.HasPrefix(result, pattern) {
			result = strings.TrimSpace(strings.TrimPrefix(result, pattern))
		}
	}

	// 移除引号
	result = strings.Trim(result, "\"'")

	// 移除多余的换行符
	result = strings.ReplaceAll(result, "\n", " ")
	result = strings.TrimSpace(result)

	return result
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

// translateArrayField 翻译数组字段
func (a *ArticleTranslator) translateArrayField(items []string, fieldType string) []string {
	var translated []string

	for _, item := range items {
		if a.containsChinese(item) {
			translatedItem, err := a.translateFieldContent(item)
			if err != nil {
				fmt.Printf("⚠️ %s翻译失败，保持原文: %s\n", fieldType, item)
				translated = append(translated, item)
			} else {
				translated = append(translated, translatedItem)
			}
		} else {
			translated = append(translated, item)
		}
	}

	return translated
}

// formatArrayField 格式化数组字段
func (a *ArticleTranslator) formatArrayField(items []string) string {
	if len(items) == 0 {
		return "[]"
	}

	var quotedItems []string
	for _, item := range items {
		quotedItems = append(quotedItems, fmt.Sprintf("\"%s\"", item))
	}

	return fmt.Sprintf("[%s]", strings.Join(quotedItems, ", "))
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

		fmt.Printf("  翻译段落 (%d/%d): %s...\n", i+1, len(paragraphs), a.truncateText(paragraph, 50))

		// 检查段落是否包含中文
		if !a.containsChinese(paragraph) {
			fmt.Printf("    跳过：无中文内容\n")
			translatedParagraphs = append(translatedParagraphs, paragraph)
			continue
		}

		// 对于包含中文的内容，强制翻译
		translatedParagraph, err := a.translator.TranslateParagraph(paragraph)
		if err != nil {
			fmt.Printf("⚠️ 段落翻译失败，保持原文: %v\n", err)
			translatedParagraphs = append(translatedParagraphs, paragraph)
		} else {
			// 验证翻译结果是否还包含中文
			if a.containsChinese(translatedParagraph) {
				fmt.Printf("    ⚠️ 翻译结果仍包含中文，尝试重新翻译...\n")
				// 尝试重新翻译
				retryTranslated, retryErr := a.retryTranslation(paragraph)
				if retryErr == nil && !a.containsChinese(retryTranslated) {
					translatedParagraphs = append(translatedParagraphs, retryTranslated)
					fmt.Printf("    ✓ 重新翻译成功\n")
				} else {
					fmt.Printf("    ⚠️ 重新翻译失败，使用首次结果\n")
					translatedParagraphs = append(translatedParagraphs, translatedParagraph)
				}
			} else {
				fmt.Printf("    ✓ 翻译完成\n")
				translatedParagraphs = append(translatedParagraphs, translatedParagraph)
			}
		}

		// 添加延迟避免API频率限制
		time.Sleep(time.Millisecond * 800)
	}

	return strings.Join(translatedParagraphs, "\n\n"), nil
}

// retryTranslation 重新翻译段落，使用更强的提示词
func (a *ArticleTranslator) retryTranslation(paragraph string) (string, error) {
	prompt := fmt.Sprintf(`请将以下中文内容完全翻译成英文，绝对不要保留任何中文字符：

%s

严格要求：
1. 必须将所有中文字符翻译为英文
2. 保持Markdown格式标记
3. 翻译要自然流畅
4. 技术术语使用准确的英文表达
5. 绝对不能在结果中保留任何中文字符
6. 直接返回翻译结果，不要添加任何解释`, paragraph)

	// 直接使用 translator 的 TranslateParagraph 方法，但需要临时修改提示词
	// 创建一个临时的翻译器实例来处理重试翻译
	request := translator.LMStudioRequest{
		Model: "gemma-3-12b-it", // 直接使用模型名称
		Messages: []translator.Message{
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

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Post("http://172.19.192.1:2234/v1/chat/completions", "application/json", bytes.NewBuffer(jsonData))
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
	return result, nil
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
