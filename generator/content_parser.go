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

	frontMatter := strings.Join(lines[0:frontMatterEnd+1], "\n")
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
		cleanItem := c.translationUtils.RemoveQuotes(item)
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
