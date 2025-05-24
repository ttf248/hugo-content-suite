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
	"tag-scanner/scanner"
	"tag-scanner/translator"
	"tag-scanner/utils"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
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
	utils.Info("开始翻译文章: %s", preview.OriginalFile)
	utils.Info("目标文件: %s", preview.EnglishFile)

	// 读取原文件
	content, err := os.ReadFile(preview.OriginalFile)
	if err != nil {
		utils.Error("读取原文件失败: %s, 错误: %v", preview.OriginalFile, err)
		return fmt.Errorf("读取原文件失败: %v", err)
	}

	utils.Info("原文件读取成功，内容长度: %d 字符", len(content))

	// 解析文章结构
	frontMatter, bodyContent := a.parseArticleContent(string(content))
	utils.Info("文章结构解析完成 - 前置数据长度: %d, 正文长度: %d", len(frontMatter), len(bodyContent))

	// 翻译前置数据
	translatedFrontMatter, err := a.translateFrontMatter(frontMatter)
	if err != nil {
		utils.Error("翻译前置数据失败: %v", err)
		return fmt.Errorf("翻译前置数据失败: %v", err)
	}

	// 分段翻译正文
	translatedBody, err := a.translateArticleBody(bodyContent)
	if err != nil {
		utils.Error("翻译正文失败: %v", err)
		return fmt.Errorf("翻译正文失败: %v", err)
	}

	// 合成最终内容
	finalContent := a.combineTranslatedContent(translatedFrontMatter, translatedBody)
	utils.Info("翻译内容合成完成，最终长度: %d 字符", len(finalContent))

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(preview.EnglishFile), 0755); err != nil {
		utils.Error("创建目录失败: %s, 错误: %v", filepath.Dir(preview.EnglishFile), err)
		return fmt.Errorf("创建目录失败: %v", err)
	}

	// 写入英文文件
	if err := os.WriteFile(preview.EnglishFile, []byte(finalContent), 0644); err != nil {
		utils.Error("写入英文文件失败: %s, 错误: %v", preview.EnglishFile, err)
		return fmt.Errorf("写入英文文件失败: %v", err)
	}

	utils.Info("文章翻译完成: %s", preview.EnglishFile)
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
		utils.Info("无前置数据需要翻译")
		return "", nil
	}

	fmt.Printf("翻译前置数据...\n")
	utils.Info("开始翻译前置数据，原始长度: %d", len(frontMatter))

	lines := strings.Split(frontMatter, "\n")
	var translatedLines []string

	for lineNum, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		utils.Debug("处理前置数据第%d行: %s", lineNum+1, line)

		if trimmedLine == "---" {
			translatedLines = append(translatedLines, line)
			utils.Debug("保留分隔符: %s", line)
			continue
		}

		// 翻译标题
		if strings.HasPrefix(trimmedLine, "title:") {
			title := a.extractFieldValue(trimmedLine, "title:")
			utils.Info("发现标题字段: %s", title)
			if title != "" && a.containsChinese(title) {
				fmt.Printf("  title: %s -> ", title)
				translatedTitle, err := a.translateFieldContent(title)
				if err != nil {
					fmt.Printf("翻译失败\n")
					utils.Warn("标题翻译失败: %s, 错误: %v", title, err)
					translatedLines = append(translatedLines, line)
				} else {
					fmt.Printf("%s\n", translatedTitle)
					translatedLines = append(translatedLines, fmt.Sprintf("title: \"%s\"", translatedTitle))
				}
			} else {
				utils.Info("标题无需翻译: %s", title)
				translatedLines = append(translatedLines, line)
			}
			continue
		}

		// 翻译描述字段
		if strings.HasPrefix(trimmedLine, "description:") {
			description := a.extractFieldValue(trimmedLine, "description:")
			utils.Info("发现描述字段: %s", description)
			if description != "" && a.containsChinese(description) {
				fmt.Printf("  description: %s -> ", description)
				translatedDescription, err := a.translateFieldContent(description)
				if err != nil {
					fmt.Printf("翻译失败\n")
					utils.Warn("描述翻译失败: %s, 错误: %v", description, err)
					translatedLines = append(translatedLines, line)
				} else {
					fmt.Printf("%s\n", translatedDescription)
					translatedLines = append(translatedLines, fmt.Sprintf("description: \"%s\"", translatedDescription))
				}
			} else {
				utils.Info("描述无需翻译: %s", description)
				translatedLines = append(translatedLines, line)
			}
			continue
		}

		// 翻译副标题
		if strings.HasPrefix(trimmedLine, "subtitle:") {
			subtitle := a.extractFieldValue(trimmedLine, "subtitle:")
			utils.Info("发现副标题字段: %s", subtitle)
			if subtitle != "" && a.containsChinese(subtitle) {
				fmt.Printf("  subtitle: %s -> ", subtitle)
				translatedSubtitle, err := a.translateFieldContent(subtitle)
				if err != nil {
					fmt.Printf("翻译失败\n")
					utils.Warn("副标题翻译失败: %s, 错误: %v", subtitle, err)
					translatedLines = append(translatedLines, line)
				} else {
					fmt.Printf("%s\n", translatedSubtitle)
					translatedLines = append(translatedLines, fmt.Sprintf("subtitle: \"%s\"", translatedSubtitle))
				}
			} else {
				utils.Info("副标题无需翻译: %s", subtitle)
				translatedLines = append(translatedLines, line)
			}
			continue
		}

		// 翻译摘要
		if strings.HasPrefix(trimmedLine, "summary:") {
			summary := a.extractFieldValue(trimmedLine, "summary:")
			utils.Info("发现摘要字段: %s", summary)
			if summary != "" && a.containsChinese(summary) {
				fmt.Printf("  summary: %s -> ", summary)
				translatedSummary, err := a.translateFieldContent(summary)
				if err != nil {
					fmt.Printf("翻译失败\n")
					utils.Warn("摘要翻译失败: %s, 错误: %v", summary, err)
					translatedLines = append(translatedLines, line)
				} else {
					fmt.Printf("%s\n", translatedSummary)
					translatedLines = append(translatedLines, fmt.Sprintf("summary: \"%s\"", translatedSummary))
				}
			} else {
				utils.Info("摘要无需翻译: %s", summary)
				translatedLines = append(translatedLines, line)
			}
			continue
		}

		// 翻译分类数组
		if strings.HasPrefix(trimmedLine, "categories:") {
			categories := a.extractArrayField(trimmedLine, "categories:")
			utils.Info("发现分类字段: %v", categories)
			if len(categories) > 0 {
				translatedCategories := a.translateArrayField(categories, "categories")
				translatedLines = append(translatedLines, fmt.Sprintf("categories: %s", a.formatArrayField(translatedCategories)))
			} else {
				utils.Info("分类数组为空")
				translatedLines = append(translatedLines, line)
			}
			continue
		}

		// 翻译作者数组
		if strings.HasPrefix(trimmedLine, "authors:") {
			authors := a.extractArrayField(trimmedLine, "authors:")
			utils.Info("发现作者字段: %v", authors)
			if len(authors) > 0 {
				translatedAuthors := a.translateArrayField(authors, "authors")
				translatedLines = append(translatedLines, fmt.Sprintf("authors: %s", a.formatArrayField(translatedAuthors)))
			} else {
				utils.Info("作者数组为空")
				translatedLines = append(translatedLines, line)
			}
			continue
		}

		// 其他字段保持不变
		utils.Debug("保留其他字段: %s", line)
		translatedLines = append(translatedLines, line)
	}

	result := strings.Join(translatedLines, "\n")
	utils.Info("前置数据翻译完成，结果长度: %d", len(result))
	return result, nil
}

// extractFieldValue 提取字段值
func (a *ArticleTranslator) extractFieldValue(line, prefix string) string {
	value := strings.TrimSpace(strings.TrimPrefix(line, prefix))
	value = strings.Trim(value, "\"'")
	return value
}

// translateFieldContent 翻译字段内容，使用优化的提示词
func (a *ArticleTranslator) translateFieldContent(content string) (string, error) {
	cfg := config.GetGlobalConfig()

	// 优化的prompt，更加精确和简洁
	prompt := fmt.Sprintf(`Please translate this Chinese text to English. Return ONLY the English translation, no explanations or additional text:

%s`, content)

	request := translator.LMStudioRequest{
		Model: cfg.LMStudio.Model,
		Messages: []translator.Message{
			{
				Role:    "system",
				Content: "You are a professional translator. You translate Chinese to English accurately and concisely. You only return the translation without any additional text, explanations, or formatting.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Stream: false,
	}

	// 记录详细请求信息到日志
	utils.Debug("LLM翻译请求 - Model: %s", request.Model)
	utils.Debug("LLM翻译请求 - 原文: %s", content)
	utils.Debug("LLM翻译请求 - Prompt: %s", prompt)

	jsonData, err := json.Marshal(request)
	if err != nil {
		utils.Error("LLM请求序列化失败: %v", err)
		return "", fmt.Errorf("序列化请求失败: %v", err)
	}

	utils.Debug("LLM请求JSON: %s", string(jsonData))

	startTime := time.Now()
	client := &http.Client{Timeout: time.Duration(cfg.LMStudio.Timeout) * time.Second}
	resp, err := client.Post(cfg.LMStudio.URL, "application/json", bytes.NewBuffer(jsonData))
	requestDuration := time.Since(startTime)

	if err != nil {
		utils.Error("LLM请求网络错误: %v, 耗时: %v", err, requestDuration)
		return "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	utils.Debug("LLM响应状态: %d, 耗时: %v", resp.StatusCode, requestDuration)

	if resp.StatusCode != http.StatusOK {
		utils.Error("LLM响应错误状态码: %d", resp.StatusCode)
		return "", fmt.Errorf("LM Studio返回错误状态: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.Error("LLM响应读取失败: %v", err)
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	utils.Debug("LLM响应原始数据: %s", string(body))

	var response translator.LMStudioResponse
	if err := json.Unmarshal(body, &response); err != nil {
		utils.Error("LLM响应解析失败: %v", err)
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	if len(response.Choices) == 0 {
		utils.Error("LLM响应无翻译结果")
		return "", fmt.Errorf("没有获取到翻译结果")
	}

	result := strings.TrimSpace(response.Choices[0].Message.Content)

	// 增强的结果清理，移除常见的多余内容
	result = a.cleanTranslationResult(result)

	// 记录翻译完成信息到日志
	utils.Info("字段翻译完成 - 原文: %s, 译文: %s, 耗时: %v", content, result, requestDuration)

	return result, nil
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
		"The translation is:",
		"Here is the translation:",
		"The English translation is:",
		"Translated:",
		"Answer:",
		"Result:",
		"Output:",
		"English translation:",
		"翻译:",
		"英文:",
		"Translation: ",
		"English: ",
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
	for strings.HasPrefix(result, "\"") && strings.HasSuffix(result, "\"") {
		result = strings.Trim(result, "\"")
		result = strings.TrimSpace(result)
	}
	for strings.HasPrefix(result, "'") && strings.HasSuffix(result, "'") {
		result = strings.Trim(result, "'")
		result = strings.TrimSpace(result)
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

	fmt.Printf("  %s: ", fieldType)
	utils.Info("开始翻译%s数组: %v", fieldType, items)

	for i, item := range items {
		utils.Debug("处理数组项目 [%d/%d]: %s", i+1, len(items), item)

		if a.containsChinese(item) {
			fmt.Printf("%s -> ", item)
			utils.Info("翻译数组项目 [%d/%d]: %s", i+1, len(items), item)

			translatedItem, err := a.translateFieldContent(item)
			if err != nil {
				fmt.Printf("失败 ")
				utils.Warn("数组项目翻译失败 [%d/%d] - %s: %s, 错误: %v", i+1, len(items), fieldType, item, err)
				translated = append(translated, item)
			} else {
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

// translateArticleBody 分段翻译正文，使用Markdown解析器
func (a *ArticleTranslator) translateArticleBody(body string) (string, error) {
	cfg := config.GetGlobalConfig()

	if strings.TrimSpace(body) == "" {
		utils.Info("正文为空，跳过翻译")
		return body, nil
	}

	fmt.Printf("\n翻译正文 (%d 字符)...\n", len(body))
	utils.Info("开始翻译正文内容，原文长度: %d 字符", len(body))

	// 使用Markdown解析器提取需要翻译的文本
	translationNodes, err := a.extractTranslatableNodes(body)
	if err != nil {
		utils.Error("Markdown解析失败: %v", err)
		return "", fmt.Errorf("Markdown解析失败: %v", err)
	}

	utils.Info("Markdown解析完成，找到 %d 个需要翻译的节点", len(translationNodes))

	// 翻译所有需要翻译的节点
	translatedNodes := make(map[int]string)
	translationCount := 0

	for i, node := range translationNodes {
		if !a.containsChinese(node.Content) {
			utils.Debug("节点 %d 无中文内容，跳过: %s", i+1, a.truncateText(node.Content, 50))
			continue
		}

		translationCount++
		fmt.Printf("  [%d/%d] 翻译 %s...", translationCount, len(translationNodes), node.Type)
		utils.Info("翻译节点 %d/%d - 类型: %s, 内容: %s", translationCount, len(translationNodes), node.Type, node.Content)

		// 翻译纯文本内容
		translatedText, err := a.translatePlainText(node.Content, translationCount, len(translationNodes))
		if err != nil {
			fmt.Printf("失败\n")
			utils.Error("节点翻译失败 %d/%d: %v", translationCount, len(translationNodes), err)
			// 翻译失败时保持原文
			translatedNodes[node.Position] = node.Content
		} else {
			fmt.Printf("完成\n")
			utils.Info("节点翻译成功 %d/%d", translationCount, len(translationNodes))
			utils.Debug("翻译结果: %s -> %s", node.Content, translatedText)
			translatedNodes[node.Position] = translatedText
		}

		// 添加延迟避免API频率限制
		if translationCount < len(translationNodes) && cfg.Translation.DelayBetweenMs > 0 {
			utils.Debug("等待 %dms 避免API频率限制", cfg.Translation.DelayBetweenMs)
			time.Sleep(time.Duration(cfg.Translation.DelayBetweenMs) * time.Millisecond)
		}
	}

	// 重新组装Markdown内容
	result := a.reassembleMarkdown(body, translationNodes, translatedNodes)

	fmt.Printf("正文翻译完成 (%d 字符)\n", len(result))
	utils.Info("正文翻译完成 - 原文长度: %d, 译文长度: %d, 翻译节点: %d", len(body), len(result), translationCount)

	return result, nil
}

// extractTranslatableNodes 提取可翻译的Markdown节点
func (a *ArticleTranslator) extractTranslatableNodes(content string) ([]MarkdownNode, error) {
	var nodes []MarkdownNode

	md := goldmark.New()
	doc := md.Parser().Parse(text.NewReader([]byte(content)))

	source := []byte(content)
	position := 0

	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		switch node := n.(type) {
		case *ast.Text:
			text := string(node.Segment.Value(source))
			if strings.TrimSpace(text) != "" {
				nodes = append(nodes, MarkdownNode{
					Type:     "text",
					Content:  text,
					Position: position,
				})
				position++
			}

		case *ast.Heading:
			// 提取标题文本
			var headingText strings.Builder
			for child := node.FirstChild(); child != nil; child = child.NextSibling() {
				if textNode, ok := child.(*ast.Text); ok {
					headingText.WriteString(string(textNode.Segment.Value(source)))
				}
			}
			if headingText.Len() > 0 {
				nodes = append(nodes, MarkdownNode{
					Type:     "heading",
					Content:  headingText.String(),
					Position: position,
					Level:    node.Level,
				})
				position++
			}

		case *ast.ListItem:
			// 提取列表项文本
			var itemText strings.Builder
			for child := node.FirstChild(); child != nil; child = child.NextSibling() {
				if para, ok := child.(*ast.Paragraph); ok {
					for grandChild := para.FirstChild(); grandChild != nil; grandChild = grandChild.NextSibling() {
						if textNode, ok := grandChild.(*ast.Text); ok {
							itemText.WriteString(string(textNode.Segment.Value(source)))
						}
					}
				}
			}
			if itemText.Len() > 0 {
				nodes = append(nodes, MarkdownNode{
					Type:     "listitem",
					Content:  itemText.String(),
					Position: position,
				})
				position++
			}

		case *ast.Blockquote:
			// 提取引用文本
			var quoteText strings.Builder
			for child := node.FirstChild(); child != nil; child = child.NextSibling() {
				if para, ok := child.(*ast.Paragraph); ok {
					for grandChild := para.FirstChild(); grandChild != nil; grandChild = grandChild.NextSibling() {
						if textNode, ok := grandChild.(*ast.Text); ok {
							quoteText.WriteString(string(textNode.Segment.Value(source)))
						}
					}
				}
			}
			if quoteText.Len() > 0 {
				nodes = append(nodes, MarkdownNode{
					Type:     "blockquote",
					Content:  quoteText.String(),
					Position: position,
				})
				position++
			}
		}

		return ast.WalkContinue, nil
	})

	if err != nil {
		return nil, err
	}

	return nodes, nil
}

// translatePlainText 翻译纯文本内容
func (a *ArticleTranslator) translatePlainText(text string, nodeNum, totalNodes int) (string, error) {
	cfg := config.GetGlobalConfig()

	// 清理文本，移除多余的空白字符
	cleanText := strings.TrimSpace(text)
	cleanText = regexp.MustCompile(`\s+`).ReplaceAllString(cleanText, " ")

	prompt := fmt.Sprintf(`Translate this Chinese text to English. Return ONLY the English translation:

%s`, cleanText)

	request := translator.LMStudioRequest{
		Model: cfg.LMStudio.Model,
		Messages: []translator.Message{
			{
				Role:    "system",
				Content: "You are a professional translator. Translate Chinese to English accurately. Return only the translation without explanations or formatting.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Stream: false,
	}

	// 记录详细的翻译请求信息
	utils.Info("文本翻译请求 %d/%d", nodeNum, totalNodes)
	utils.Debug("文本翻译请求 %d/%d - Model: %s", nodeNum, totalNodes, request.Model)
	utils.Debug("文本翻译请求 %d/%d - 原文: %s", nodeNum, totalNodes, cleanText)
	utils.Debug("文本翻译请求 %d/%d - Prompt: %s", nodeNum, totalNodes, prompt)

	jsonData, err := json.Marshal(request)
	if err != nil {
		utils.Error("文本翻译请求 %d/%d 序列化失败: %v", nodeNum, totalNodes, err)
		return "", fmt.Errorf("序列化请求失败: %v", err)
	}

	utils.Debug("文本翻译请求 %d/%d JSON: %s", nodeNum, totalNodes, string(jsonData))

	startTime := time.Now()
	client := &http.Client{Timeout: time.Duration(cfg.LMStudio.Timeout) * time.Second}
	resp, err := client.Post(cfg.LMStudio.URL, "application/json", bytes.NewBuffer(jsonData))
	requestDuration := time.Since(startTime)

	if err != nil {
		utils.Error("文本翻译请求 %d/%d 网络错误: %v, 耗时: %v", nodeNum, totalNodes, err, requestDuration)
		return "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	utils.Info("文本翻译响应 %d/%d - 状态码: %d, 耗时: %v", nodeNum, totalNodes, resp.StatusCode, requestDuration)

	if resp.StatusCode != http.StatusOK {
		utils.Error("文本翻译响应 %d/%d 错误状态码: %d", nodeNum, totalNodes, resp.StatusCode)
		return "", fmt.Errorf("LM Studio返回错误状态: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.Error("文本翻译响应 %d/%d 读取失败: %v", nodeNum, totalNodes, err)
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	utils.Debug("文本翻译响应 %d/%d 原始数据: %s", nodeNum, totalNodes, string(body))

	var response translator.LMStudioResponse
	if err := json.Unmarshal(body, &response); err != nil {
		utils.Error("文本翻译响应 %d/%d 解析失败: %v", nodeNum, totalNodes, err)
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	if len(response.Choices) == 0 {
		utils.Error("文本翻译响应 %d/%d 无结果", nodeNum, totalNodes)
		return "", fmt.Errorf("没有获取到翻译结果")
	}

	result := strings.TrimSpace(response.Choices[0].Message.Content)

	// 清理翻译结果
	result = a.cleanTranslationResult(result)

	// 记录翻译完成信息
	utils.Info("文本翻译完成 %d/%d - 原文长度: %d, 译文长度: %d, 耗时: %v",
		nodeNum, totalNodes, len(cleanText), len(result), requestDuration)
	utils.Debug("文本翻译结果 %d/%d: %s", nodeNum, totalNodes, result)

	return result, nil
}

// reassembleMarkdown 重新组装Markdown内容
func (a *ArticleTranslator) reassembleMarkdown(originalContent string, nodes []MarkdownNode, translations map[int]string) string {
	lines := strings.Split(originalContent, "\n")
	var result []string

	nodeIndex := 0

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// 跳过空行
		if trimmedLine == "" {
			result = append(result, line)
			continue
		}

		// 检查是否为代码块
		if strings.HasPrefix(trimmedLine, "```") {
			result = append(result, line)
			continue
		}

		// 处理标题
		if strings.HasPrefix(trimmedLine, "#") {
			if nodeIndex < len(nodes) && nodes[nodeIndex].Type == "heading" {
				if translation, exists := translations[nodes[nodeIndex].Position]; exists {
					// 保持原始的标题格式，只替换文本内容
					level := strings.Repeat("#", nodes[nodeIndex].Level)
					result = append(result, fmt.Sprintf("%s %s", level, translation))
				} else {
					result = append(result, line)
				}
				nodeIndex++
			} else {
				result = append(result, line)
			}
			continue
		}

		// 处理列表项
		if strings.HasPrefix(trimmedLine, "- ") || strings.HasPrefix(trimmedLine, "* ") ||
			regexp.MustCompile(`^\d+\. `).MatchString(trimmedLine) {
			if nodeIndex < len(nodes) && nodes[nodeIndex].Type == "listitem" {
				if translation, exists := translations[nodes[nodeIndex].Position]; exists {
					// 保持原始的列表格式，只替换文本内容
					prefix := ""
					if strings.HasPrefix(trimmedLine, "- ") {
						prefix = "- "
					} else if strings.HasPrefix(trimmedLine, "* ") {
						prefix = "* "
					} else if match := regexp.MustCompile(`^(\d+\. )`).FindString(trimmedLine); match != "" {
						prefix = match
					}
					result = append(result, prefix+translation)
				} else {
					result = append(result, line)
				}
				nodeIndex++
			} else {
				result = append(result, line)
			}
			continue
		}

		// 处理引用
		if strings.HasPrefix(trimmedLine, ">") {
			if nodeIndex < len(nodes) && nodes[nodeIndex].Type == "blockquote" {
				if translation, exists := translations[nodes[nodeIndex].Position]; exists {
					result = append(result, "> "+translation)
				} else {
					result = append(result, line)
				}
				nodeIndex++
			} else {
				result = append(result, line)
			}
			continue
		}

		// 处理普通文本段落
		if nodeIndex < len(nodes) && nodes[nodeIndex].Type == "text" {
			if translation, exists := translations[nodes[nodeIndex].Position]; exists {
				result = append(result, translation)
			} else {
				result = append(result, line)
			}
			nodeIndex++
		} else {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
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
