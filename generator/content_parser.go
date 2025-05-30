package generator

import (
	"fmt"
	"hugo-content-suite/config"
	"regexp"
	"strings"
)

// ParagraphMapping 段落映射信息，用于追踪拆分的段落关系
type ParagraphMapping struct {
	OriginalIndex int      // 原始段落索引
	SplitParts    []string // 拆分后的部分
	IsOriginal    bool     // 是否为原始段落（未拆分）
}

// SplitResult 段落拆分结果
type SplitResult struct {
	Paragraphs []string           // 拆分后的段落列表
	Mappings   []ParagraphMapping // 段落映射关系
}

// ContentParser 内容解析器
type ContentParser struct {
	translationUtils *TranslationUtils
	config           *config.Config
}

// NewContentParser 创建内容解析器
func NewContentParser() *ContentParser {
	return &ContentParser{
		translationUtils: NewTranslationUtils(),
		config:           config.GetGlobalConfig(),
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
	cleanedParagraphs := c.cleanEmptyParagraphs(paragraphs)

	// 应用段落拆分
	splitParagraphs := c.applySplittingToParagraphs(cleanedParagraphs)

	// 生成并记录统计信息
	if c.config.Paragraph.EnableSplitting {
		stats := c.GetParagraphSplitStats(cleanedParagraphs, splitParagraphs)
		c.LogParagraphSplitInfo(stats)
	}

	return splitParagraphs
}

// ParseContentIntoParagraphsWithMapping 将内容解析为段落并保留映射关系
func (c *ContentParser) ParseContentIntoParagraphsWithMapping(content string) (*SplitResult, error) {
	if strings.TrimSpace(content) == "" {
		return &SplitResult{
			Paragraphs: []string{},
			Mappings:   []ParagraphMapping{},
		}, nil
	}

	// 首先使用现有的解析逻辑获取基础段落
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

	cleanedParagraphs := c.cleanEmptyParagraphs(paragraphs)

	// 应用段落拆分并生成映射关系
	return c.applySplittingWithMapping(cleanedParagraphs), nil
}

// applySplittingWithMapping 对段落列表应用拆分并生成映射关系
func (c *ContentParser) applySplittingWithMapping(paragraphs []string) *SplitResult {
	var resultParagraphs []string
	var mappings []ParagraphMapping

	for originalIndex, paragraph := range paragraphs {
		// 检查是否为特殊格式（代码块、标题等），这些不需要拆分
		if c.shouldSkipSplitting(paragraph) || !c.config.Paragraph.EnableSplitting {
			resultParagraphs = append(resultParagraphs, paragraph)
			mappings = append(mappings, ParagraphMapping{
				OriginalIndex: originalIndex,
				SplitParts:    []string{paragraph},
				IsOriginal:    true,
			})
			continue
		}

		// 对普通段落应用拆分
		splitParagraphs := c.splitLongParagraph(paragraph)
		resultParagraphs = append(resultParagraphs, splitParagraphs...)

		// 记录映射关系
		mappings = append(mappings, ParagraphMapping{
			OriginalIndex: originalIndex,
			SplitParts:    splitParagraphs,
			IsOriginal:    len(splitParagraphs) == 1 && splitParagraphs[0] == paragraph,
		})
	}

	result := &SplitResult{
		Paragraphs: resultParagraphs,
		Mappings:   mappings,
	}

	// 生成并记录统计信息
	if c.config.Paragraph.EnableSplitting {
		stats := c.GetParagraphSplitStats(paragraphs, resultParagraphs)
		c.LogParagraphSplitInfo(stats)
	}

	return result
}

// MergeTranslatedParagraphs 合并翻译后的拆分段落
func (c *ContentParser) MergeTranslatedParagraphs(translatedParagraphs []string, mappings []ParagraphMapping) ([]string, error) {
	if !c.config.Paragraph.MergeAfterTranslation {
		// 如果配置为不合并，直接返回翻译后的段落
		return translatedParagraphs, nil
	}

	var mergedParagraphs []string
	var currentIndex int

	for _, mapping := range mappings {
		if mapping.IsOriginal {
			// 原始段落未被拆分，直接添加
			if currentIndex < len(translatedParagraphs) {
				mergedParagraphs = append(mergedParagraphs, translatedParagraphs[currentIndex])
				currentIndex++
			}
		} else {
			// 段落被拆分了，需要合并翻译后的片段
			var parts []string
			for i := 0; i < len(mapping.SplitParts); i++ {
				if currentIndex < len(translatedParagraphs) {
					parts = append(parts, strings.TrimSpace(translatedParagraphs[currentIndex]))
					currentIndex++
				}
			}

			if len(parts) > 0 {
				// 合并拆分的段落为单个段落
				merged := strings.Join(parts, " ")
				mergedParagraphs = append(mergedParagraphs, merged)
			}
		}
	}

	return mergedParagraphs, nil
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

// applySplittingToParagraphs 对段落列表应用拆分
func (c *ContentParser) applySplittingToParagraphs(paragraphs []string) []string {
	if !c.config.Paragraph.EnableSplitting {
		return paragraphs
	}

	var result []string

	for _, paragraph := range paragraphs {
		// 检查是否为特殊格式（代码块、标题等），这些不需要拆分
		if c.shouldSkipSplitting(paragraph) {
			result = append(result, paragraph)
			continue
		}

		// 对普通段落应用拆分
		splitParagraphs := c.splitLongParagraph(paragraph)
		result = append(result, splitParagraphs...)
	}

	return result
}

// shouldSkipSplitting 检查段落是否应该跳过拆分
func (c *ContentParser) shouldSkipSplitting(paragraph string) bool {
	trimmed := strings.TrimSpace(paragraph)

	// 空段落
	if trimmed == "" {
		return true
	}

	// 代码块
	if strings.HasPrefix(trimmed, "```") || strings.Contains(paragraph, "```") {
		return true
	}

	// 标题行
	if c.isHeaderLine(paragraph) {
		return true
	}

	// 其他块级元素
	if c.isBlockLevelElement(paragraph) {
		return true
	}

	// 特殊格式行
	lines := strings.Split(paragraph, "\n")
	if len(lines) == 1 && c.IsMarkdownElement(lines[0]) {
		return true
	}

	// 表格
	if strings.Contains(trimmed, "|") && strings.Count(trimmed, "|") >= 2 {
		return true
	}

	return false
}

// splitLongParagraph 拆分过长的段落
func (c *ContentParser) splitLongParagraph(paragraph string) []string {
	// 如果未启用拆分或段落不超长，直接返回
	if !c.config.Paragraph.EnableSplitting || len(paragraph) <= c.config.Paragraph.MaxLength {
		return []string{paragraph}
	}

	var result []string

	// 如果启用了在句子边界拆分
	if c.config.Paragraph.SplitAtSentences {
		result = c.splitAtSentenceBoundaries(paragraph)
	} else {
		result = c.splitAtCharacterLimit(paragraph)
	}

	// 过滤掉过短的段落片段
	var filteredResult []string
	for _, part := range result {
		if strings.TrimSpace(part) != "" && len(strings.TrimSpace(part)) >= c.config.Paragraph.MinSplitLength {
			filteredResult = append(filteredResult, part)
		}
	}

	// 如果过滤后没有结果，返回原段落
	if len(filteredResult) == 0 {
		return []string{paragraph}
	}

	return filteredResult
}

// splitAtSentenceBoundaries 在句子边界拆分段落
func (c *ContentParser) splitAtSentenceBoundaries(paragraph string) []string {
	var result []string
	var currentSegment strings.Builder

	// 定义句子结束标记的正则表达式
	sentenceEndRegex := regexp.MustCompile(`[.!?。！？]\s*`)

	// 查找所有句子结束位置
	matches := sentenceEndRegex.FindAllStringIndex(paragraph, -1)

	if len(matches) == 0 {
		// 没有找到句子边界，按字符限制拆分
		return c.splitAtCharacterLimit(paragraph)
	}

	lastEnd := 0

	for _, match := range matches {
		sentenceEnd := match[1]
		sentence := paragraph[lastEnd:sentenceEnd]

		// 检查当前段落加上这个句子是否超过长度限制
		if currentSegment.Len()+len(sentence) > c.config.Paragraph.MaxLength && currentSegment.Len() > 0 {
			// 保存当前段落
			if segment := strings.TrimSpace(currentSegment.String()); segment != "" {
				result = append(result, segment)
			}
			currentSegment.Reset()
		}

		currentSegment.WriteString(sentence)
		lastEnd = sentenceEnd
	}

	// 处理剩余部分
	if lastEnd < len(paragraph) {
		remaining := paragraph[lastEnd:]
		if currentSegment.Len()+len(remaining) > c.config.Paragraph.MaxLength && currentSegment.Len() > 0 {
			if segment := strings.TrimSpace(currentSegment.String()); segment != "" {
				result = append(result, segment)
			}
			if remaining := strings.TrimSpace(remaining); remaining != "" {
				result = append(result, remaining)
			}
		} else {
			currentSegment.WriteString(remaining)
		}
	}

	// 添加最后的段落
	if segment := strings.TrimSpace(currentSegment.String()); segment != "" {
		result = append(result, segment)
	}

	return result
}

// splitAtCharacterLimit 按字符限制拆分段落
func (c *ContentParser) splitAtCharacterLimit(paragraph string) []string {
	var result []string
	maxLength := c.config.Paragraph.MaxLength

	// 如果段落长度小于等于限制，直接返回
	if len(paragraph) <= maxLength {
		return []string{paragraph}
	}

	// 按空白字符分割为单词
	words := strings.Fields(paragraph)
	var currentSegment strings.Builder

	for _, word := range words {
		// 检查添加这个单词是否会超过长度限制
		testLength := currentSegment.Len()
		if testLength > 0 {
			testLength += 1 // 加上空格
		}
		testLength += len(word)

		if testLength > maxLength && currentSegment.Len() > 0 {
			// 保存当前段落并开始新的段落
			if segment := strings.TrimSpace(currentSegment.String()); segment != "" {
				result = append(result, segment)
			}
			currentSegment.Reset()
			currentSegment.WriteString(word)
		} else {
			if currentSegment.Len() > 0 {
				currentSegment.WriteString(" ")
			}
			currentSegment.WriteString(word)
		}
	}

	// 添加最后的段落
	if segment := strings.TrimSpace(currentSegment.String()); segment != "" {
		result = append(result, segment)
	}

	return result
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

// GetParagraphSplitStats 获取段落拆分统计信息
func (c *ContentParser) GetParagraphSplitStats(originalParagraphs, splitParagraphs []string) map[string]interface{} {
	stats := make(map[string]interface{})

	originalCount := len(originalParagraphs)
	splitCount := len(splitParagraphs)

	// 计算原始段落中超长的数量
	longParagraphsCount := 0
	totalOriginalLength := 0
	totalSplitLength := 0

	for _, p := range originalParagraphs {
		length := len(strings.TrimSpace(p))
		totalOriginalLength += length
		if length > c.config.Paragraph.MaxLength {
			longParagraphsCount++
		}
	}

	for _, p := range splitParagraphs {
		totalSplitLength += len(strings.TrimSpace(p))
	}

	stats["original_paragraph_count"] = originalCount
	stats["split_paragraph_count"] = splitCount
	stats["long_paragraphs_count"] = longParagraphsCount
	stats["paragraphs_added"] = splitCount - originalCount
	stats["average_original_length"] = 0
	stats["average_split_length"] = 0

	if originalCount > 0 {
		stats["average_original_length"] = totalOriginalLength / originalCount
	}

	if splitCount > 0 {
		stats["average_split_length"] = totalSplitLength / splitCount
	}

	stats["splitting_enabled"] = c.config.Paragraph.EnableSplitting
	stats["max_length_config"] = c.config.Paragraph.MaxLength
	stats["min_split_length_config"] = c.config.Paragraph.MinSplitLength
	stats["split_at_sentences"] = c.config.Paragraph.SplitAtSentences

	return stats
}

// LogParagraphSplitInfo 记录段落拆分信息
func (c *ContentParser) LogParagraphSplitInfo(stats map[string]interface{}) {
	if !stats["splitting_enabled"].(bool) {
		fmt.Println("📝 段落拆分功能已禁用")
		return
	}

	originalCount := stats["original_paragraph_count"].(int)
	splitCount := stats["split_paragraph_count"].(int)
	longCount := stats["long_paragraphs_count"].(int)
	added := stats["paragraphs_added"].(int)

	if added > 0 {
		fmt.Printf("✂️ 段落拆分完成: %d个段落 → %d个段落 (新增%d个)\n",
			originalCount, splitCount, added)
		fmt.Printf("📊 发现%d个超长段落已被拆分\n", longCount)
		fmt.Printf("📏 平均长度: %d字符 → %d字符\n",
			stats["average_original_length"].(int),
			stats["average_split_length"].(int))
	} else {
		fmt.Printf("📝 段落分析完成: %d个段落，无需拆分\n", originalCount)
	}
}
