package generator

import (
	"fmt"
	"regexp"
	"strings"
)

// ContentParser 内容解析器
type ContentParser struct {
	translationUtils *TranslationUtils
}

// NewContentParser 创建内容解析器
func NewContentParser() *ContentParser {
	return &ContentParser{
		translationUtils: NewTranslationUtils(),
	}
}

// ParseArticleContent 解析文章内容，分离前置数据和正文
func (c *ContentParser) ParseArticleContent(content string) (string, string) {
	lines := strings.Split(content, "\n")

	if len(lines) < 3 || strings.TrimSpace(lines[0]) != "---" {
		return "", content
	}

	frontMatterEnd := -1
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			frontMatterEnd = i
			break
		}
	}

	if frontMatterEnd == -1 {
		return "", content
	}

	frontMatter := strings.Join(lines[1:frontMatterEnd], "\n") // 不包含前后的 ---
	body := strings.Join(lines[frontMatterEnd+1:], "\n")

	return frontMatter, body
}

// ExtractFieldValue 提取字段值
func (c *ContentParser) ExtractFieldValue(line, prefix string) string {
	value := strings.TrimSpace(strings.TrimPrefix(line, prefix))
	return strings.Trim(value, "\"'")
}

// ExtractArrayField 提取数组字段
func (c *ContentParser) ExtractArrayField(line, prefix string) []string {
	value := strings.TrimSpace(strings.TrimPrefix(line, prefix))
	value = strings.Trim(value, "[]")

	if value == "" {
		return []string{}
	}

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

// FormatArrayField 格式化数组字段
func (c *ContentParser) FormatArrayField(items []string) string {
	if len(items) == 0 {
		return "[]"
	}

	var quotedItems []string
	for _, item := range items {
		// 彻底清理引号
		cleanItem := c.translationUtils.RemoveQuotes(item)
		// 再次确保移除双引号
		cleanItem = strings.ReplaceAll(cleanItem, "\"", "")
		cleanItem = strings.ReplaceAll(cleanItem, "'", "")
		cleanItem = strings.TrimSpace(cleanItem)
		quotedItems = append(quotedItems, fmt.Sprintf("\"%s\"", cleanItem))
	}

	return fmt.Sprintf("[%s]", strings.Join(quotedItems, ", "))
}

// CombineTranslatedContent 合并翻译后的内容
func (c *ContentParser) CombineTranslatedContent(frontMatter, body string) string {
	if frontMatter == "" {
		return body
	}
	return frontMatter + "\n\n" + body
}

// AnalyzeArticleContent 分析文章内容统计
func (c *ContentParser) AnalyzeArticleContent(content string) (int, int) {
	_, body := c.ParseArticleContent(content)

	// 统计字数
	wordCount := len(strings.Fields(body))

	// 统计段落数
	paragraphs := c.splitIntoParagraphs(body)
	paragraphCount := len(paragraphs)

	return wordCount, paragraphCount
}

// EstimateTranslationTime 估算翻译时间
func (c *ContentParser) EstimateTranslationTime(paragraphCount int) string {
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

// splitIntoParagraphs 将文本分割成段落
func (c *ContentParser) splitIntoParagraphs(text string) []string {
	preliminaryParagraphs := strings.Split(text, "\n\n")
	var finalParagraphs []string

	for _, p := range preliminaryParagraphs {
		trimmed := strings.TrimSpace(p)
		if trimmed == "" {
			continue
		}

		if strings.Contains(trimmed, "```") {
			finalParagraphs = append(finalParagraphs, trimmed)
		} else {
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
					if c.isSpecialFormatLine(line) && c.translationUtils.ContainsChinese(line) {
						if len(currentParagraph) > 0 {
							finalParagraphs = append(finalParagraphs, strings.Join(currentParagraph, "\n"))
							currentParagraph = nil
						}
						finalParagraphs = append(finalParagraphs, line)
					} else if c.isSpecialFormatLine(line) {
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
func (c *ContentParser) isSpecialFormatLine(line string) bool {
	trimmed := strings.TrimSpace(line)

	// 标题
	if strings.HasPrefix(trimmed, "#") {
		return true
	}

	// 列表
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

// IsMarkdownElement 检查是否为markdown元素行
func (c *ContentParser) IsMarkdownElement(line string) bool {
	trimmed := strings.TrimSpace(line)

	// 空行
	if trimmed == "" {
		return false
	}

	// 代码块标记
	if strings.HasPrefix(trimmed, "```") {
		return true
	}

	// 标题 (# ## ### 等)
	if strings.HasPrefix(trimmed, "#") && (len(trimmed) == 1 || trimmed[1] == '#' || trimmed[1] == ' ') {
		return true
	}

	// 有序列表
	if matched, _ := regexp.MatchString(`^\d+\.\s`, trimmed); matched {
		return true
	}

	// 无序列表
	if matched, _ := regexp.MatchString(`^[-*+]\s`, trimmed); matched {
		return true
	}

	// 引用
	if strings.HasPrefix(trimmed, ">") {
		return true
	}

	// 水平分割线
	if matched, _ := regexp.MatchString(`^(-{3,}|\*{3,}|_{3,})$`, trimmed); matched {
		return true
	}

	// 表格行
	if strings.Contains(trimmed, "|") && (strings.Count(trimmed, "|") >= 2) {
		return true
	}

	// 链接定义
	if matched, _ := regexp.MatchString(`^\[.+\]:\s+`, trimmed); matched {
		return true
	}

	// HTML标签
	if matched, _ := regexp.MatchString(`^<[^>]+>`, trimmed); matched {
		return true
	}

	return false
}

// ExtractMarkdownPrefix 提取markdown前缀
func (c *ContentParser) ExtractMarkdownPrefix(line string) (string, string) {
	trimmed := strings.TrimSpace(line)

	// 标题处理
	if strings.HasPrefix(trimmed, "#") {
		hashCount := 0
		spaceIndex := -1

		// 计算连续的#号数量，并找到第一个空格位置
		for i, r := range trimmed {
			if r == '#' {
				hashCount++
			} else if r == ' ' {
				spaceIndex = i
				break
			} else {
				// 遇到非#非空格字符，说明格式不标准
				break
			}
		}

		// 如果找到了空格，提取前缀和内容
		if spaceIndex > 0 {
			prefix := trimmed[:spaceIndex+1] // 包含空格
			content := strings.TrimSpace(trimmed[spaceIndex+1:])
			return prefix, content
		}

		// 如果只有#号没有空格，补充空格
		if hashCount > 0 && spaceIndex == -1 {
			if hashCount == len(trimmed) {
				// 只有#号，没有内容
				return trimmed + " ", ""
			} else {
				// 有#号和内容但没有空格，插入空格
				prefix := trimmed[:hashCount] + " "
				content := strings.TrimSpace(trimmed[hashCount:])
				return prefix, content
			}
		}
	}

	// 列表项
	if matched, _ := regexp.MatchString(`^(\d+\.\s|[-*+]\s)`, trimmed); matched {
		re := regexp.MustCompile(`^(\d+\.\s|[-*+]\s)`)
		matches := re.FindStringSubmatch(trimmed)
		if len(matches) > 1 {
			prefix := matches[1]
			content := strings.TrimSpace(trimmed[len(prefix):])
			return prefix, content
		}
	}

	// 引用
	if strings.HasPrefix(trimmed, ">") {
		re := regexp.MustCompile(`^(>\s*)`)
		matches := re.FindStringSubmatch(trimmed)
		if len(matches) > 1 {
			prefix := matches[1]
			content := strings.TrimSpace(trimmed[len(prefix):])
			return prefix, content
		}
	}

	return "", trimmed
}

// ReconstructMarkdownLine 重构markdown行
func (c *ContentParser) ReconstructMarkdownLine(prefix, translatedContent string) string {
	if prefix == "" {
		return translatedContent
	}
	return prefix + translatedContent
}

// ParseContentIntoParagraphs 将内容解析为段落
func (c *ContentParser) ParseContentIntoParagraphs(content string) []string {
	if strings.TrimSpace(content) == "" {
		return []string{}
	}

	lines := strings.Split(content, "\n")
	var paragraphs []string
	var currentParagraph []string
	inCodeBlock := false

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// 检测代码块开始/结束
		if strings.HasPrefix(trimmedLine, "```") {
			if !inCodeBlock {
				// 代码块开始
				if len(currentParagraph) > 0 {
					paragraphs = append(paragraphs, strings.Join(currentParagraph, "\n"))
					currentParagraph = nil
				}
				inCodeBlock = true
				currentParagraph = append(currentParagraph, line)
			} else if trimmedLine == "```" || strings.HasPrefix(trimmedLine, "```") {
				// 代码块结束
				currentParagraph = append(currentParagraph, line)
				paragraphs = append(paragraphs, strings.Join(currentParagraph, "\n"))
				currentParagraph = nil
				inCodeBlock = false
			} else {
				currentParagraph = append(currentParagraph, line)
			}
			continue
		}

		// 在代码块内，直接添加行
		if inCodeBlock {
			currentParagraph = append(currentParagraph, line)
			continue
		}

		// 空行表示段落结束
		if trimmedLine == "" {
			if len(currentParagraph) > 0 {
				paragraphs = append(paragraphs, strings.Join(currentParagraph, "\n"))
				currentParagraph = nil
			}
			continue
		}

		// 标题行必须单独成段
		if c.isHeaderLine(line) {
			// 先结束当前段落
			if len(currentParagraph) > 0 {
				paragraphs = append(paragraphs, strings.Join(currentParagraph, "\n"))
				currentParagraph = nil
			}
			// 标题行单独成段
			paragraphs = append(paragraphs, line)
			continue
		}

		// 其他特殊markdown元素单独成段
		if c.isBlockLevelElement(line) {
			// 先结束当前段落
			if len(currentParagraph) > 0 {
				paragraphs = append(paragraphs, strings.Join(currentParagraph, "\n"))
				currentParagraph = nil
			}
			// 单独成段
			paragraphs = append(paragraphs, line)
			continue
		}

		// 普通行添加到当前段落
		currentParagraph = append(currentParagraph, line)
	}

	// 处理最后一个段落
	if len(currentParagraph) > 0 {
		paragraphs = append(paragraphs, strings.Join(currentParagraph, "\n"))
	}

	return c.cleanEmptyParagraphs(paragraphs)
}

// isHeaderLine 检查是否为标题行
func (c *ContentParser) isHeaderLine(line string) bool {
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

// isBlockLevelElement 检查是否为块级元素
func (c *ContentParser) isBlockLevelElement(line string) bool {
	trimmed := strings.TrimSpace(line)

	// 水平分割线
	if matched, _ := regexp.MatchString(`^(-{3,}|\*{3,}|_{3,})$`, trimmed); matched {
		return true
	}

	// HTML块级标签
	if matched, _ := regexp.MatchString(`^<(div|p|h[1-6]|blockquote|pre|table|ul|ol|li)[^>]*>`, trimmed); matched {
		return true
	}

	// 链接定义
	if matched, _ := regexp.MatchString(`^\[.+\]:\s+https?://`, trimmed); matched {
		return true
	}

	return false
}

// ExtractHeaderPrefix 提取标题前缀和内容
func (c *ContentParser) ExtractHeaderPrefix(line string) (string, string) {
	trimmed := strings.TrimSpace(line)

	if !strings.HasPrefix(trimmed, "#") {
		return "", line
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

	// 构建前缀
	prefix := strings.Repeat("#", hashCount)

	// 提取内容
	content := ""
	if len(trimmed) > hashCount {
		if trimmed[hashCount] == ' ' {
			// 有空格，提取空格后的内容
			content = strings.TrimSpace(trimmed[hashCount+1:])
			prefix += " "
		} else {
			// 没有空格，提取#号后的内容
			content = strings.TrimSpace(trimmed[hashCount:])
			prefix += " " // 补充空格
		}
	}

	return prefix, content
}

// ReconstructHeaderLine 重构标题行
func (c *ContentParser) ReconstructHeaderLine(prefix, translatedContent string) string {
	if prefix == "" {
		return translatedContent
	}
	return prefix + translatedContent
}

// cleanEmptyParagraphs 清理空段落
func (c *ContentParser) cleanEmptyParagraphs(paragraphs []string) []string {
	var result []string
	for _, p := range paragraphs {
		if strings.TrimSpace(p) != "" {
			result = append(result, p)
		}
	}
	return result
}

// CountTranslatableParagraphs 统计需要翻译的段落数
func (c *ContentParser) CountTranslatableParagraphs(paragraphs []string) int {
	count := 0
	for _, p := range paragraphs {
		if c.needsTranslation(p) {
			count++
		}
	}
	return count
}

// needsTranslation 检查段落是否需要翻译
func (c *ContentParser) needsTranslation(paragraph string) bool {
	trimmed := strings.TrimSpace(paragraph)

	// 空段落不需要翻译
	if trimmed == "" {
		return false
	}

	// 纯代码块不需要翻译
	if strings.HasPrefix(trimmed, "```") && strings.HasSuffix(trimmed, "```") {
		return false
	}

	// 水平分割线不需要翻译
	if matched, _ := regexp.MatchString(`^(-{3,}|\*{3,}|_{3,})$`, trimmed); matched {
		return false
	}

	// 链接定义不需要翻译
	if matched, _ := regexp.MatchString(`^\[.+\]:\s+https?://`, trimmed); matched {
		return false
	}

	// 检查是否包含中文
	return c.translationUtils.ContainsChinese(paragraph)
}
