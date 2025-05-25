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
