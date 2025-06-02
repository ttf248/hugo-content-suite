package scanner

import (
	"fmt"
	"hugo-content-suite/config"
	"hugo-content-suite/models"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/tmc/langchaingo/textsplitter"
	"gopkg.in/yaml.v3"
)

// Article 类型别名，方便引用
type Article = models.Article

// 默认只扫描 index.md，不读取内容详情
func ScanArticles(dir string) ([]Article, error) {
	return scanArticlesInternal(dir, false, false)
}

// 支持 allLangs 参数，不读取内容详情
func ScanArticlesWithLangs(dir string, allLangs bool) ([]Article, error) {
	return scanArticlesInternal(dir, allLangs, false)
}

// 用于翻译模块：扫描并读取完整内容信息
func ScanArticlesForTranslation(dir string) ([]Article, error) {
	return scanArticlesInternal(dir, false, true)
}

// 内部统一扫描函数
func scanArticlesInternal(dir string, allLangs bool, withContent bool) ([]Article, error) {
	var articles []Article

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		base := filepath.Base(path)
		if allLangs {
			// 匹配 index.md, index.en.md, index.zh-cn.md 等
			if !strings.HasPrefix(base, "index.") || !strings.HasSuffix(base, ".md") {
				return nil
			}
		} else {
			// 只扫描 index.md
			if base != "index.md" {
				return nil
			}
		}

		article, err := parseMarkdownFile(path, withContent)
		if err != nil {
			return nil
		}

		if article != nil {
			articles = append(articles, *article)
		}

		return nil
	})

	return articles, err
}

func parseMarkdownFile(filePath string, withContent bool) (*Article, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 读取整个文件内容
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	contentStr := string(content)

	// 查找前置数据边界
	lines := strings.Split(contentStr, "\n")
	var frontMatterStart, frontMatterEnd int = -1, -1

	for i, line := range lines {
		if strings.TrimSpace(line) == "---" {
			if frontMatterStart == -1 {
				frontMatterStart = i
			} else {
				frontMatterEnd = i
				break
			}
		}
	}

	article := &Article{
		FilePath: filePath,
	}

	// 解析前置数据
	var frontMatterContent string
	var bodyContent string

	if frontMatterStart >= 0 && frontMatterEnd > frontMatterStart {
		// 提取前置数据
		frontMatterLines := lines[frontMatterStart+1 : frontMatterEnd]
		frontMatterContent = strings.Join(frontMatterLines, "\n")

		// 定义前置数据结构体
		type FrontMatter struct {
			Title      string   `yaml:"title"`
			Subtitle   string   `yaml:"subtitle"`
			Summary    string   `yaml:"summary"`
			Tags       []string `yaml:"tags"`
			Categories []string `yaml:"categories"`
			Date       string   `yaml:"date"`
			LastMod    string   `yaml:"lastmod"`
			Featured   bool     `yaml:"featured"`
			Draft      bool     `yaml:"draft"`
			Slug       string   `yaml:"slug"`
		}

		// 解析 YAML 前置数据
		var frontMatter FrontMatter
		if err := yaml.Unmarshal([]byte(frontMatterContent), &frontMatter); err != nil {
			fmt.Printf("❌ YAML解析错误: %s\n", err)
			fmt.Printf("📄 文章路径: %s\n", filePath)
			os.Exit(1)
		}
		article.Title = frontMatter.Title
		article.Subtitle = frontMatter.Subtitle
		article.Summary = frontMatter.Summary
		article.Tags = frontMatter.Tags
		article.Categories = frontMatter.Categories
		article.Date = frontMatter.Date
		article.LastMod = frontMatter.LastMod
		article.Featured = frontMatter.Featured
		article.Draft = frontMatter.Draft
		article.Slug = frontMatter.Slug

		// 提取正文内容
		if frontMatterEnd+1 < len(lines) {
			bodyLines := lines[frontMatterEnd+1:]
			bodyContent = strings.Join(bodyLines, "\n")
		}
	} else {
		// 没有前置数据，整个内容都是正文
		bodyContent = contentStr
	}

	// 如果需要内容信息，则填充相关字段
	if withContent {
		// 构建前置信息
		if frontMatterContent != "" {
			article.FrontMatter = frontMatterContent
		}

		// 解析正文为段落
		article.BodyContent = splitTextIntoParagraphs(bodyContent)

		// 计算正文字符数
		article.CharCount = len([]rune(bodyContent))
	}

	return article, nil
}

// splitTextIntoParagraphs 将文本分割成段落，使用 langchaingo 的 MarkdownTextSplitter
func splitTextIntoParagraphs(text string) []string {
	cfg := config.GetGlobalConfig()

	// 创建 MarkdownTextSplitter，设置较大的 chunk 大小以保持段落完整
	splitter := textsplitter.NewMarkdownTextSplitter(
		textsplitter.WithChunkSize(cfg.Paragraph.MaxLength), // 设置较大的 chunk 大小
		textsplitter.WithChunkOverlap(0),                    // 不需要重叠
		textsplitter.WithCodeBlocks(true),
	)

	// 使用 markdown splitter 分割文本
	chunks, err := splitter.SplitText(text)
	if err != nil {
		// 如果分割失败，回退到简单的段落分割
		paragraphs := strings.Split(strings.TrimSpace(text), "\n\n")
		var result []string
		for _, paragraph := range paragraphs {
			paragraph = strings.TrimSpace(paragraph)
			if paragraph != "" {
				result = append(result, paragraph)
			}
		}
		return result
	}

	return chunks
}

// extractTagsFromYAML 从 YAML 数据中提取标签
func extractTagsFromYAML(tags interface{}) []string {
	var result []string

	switch v := tags.(type) {
	case []interface{}:
		for _, tag := range v {
			if tagStr, ok := tag.(string); ok {
				result = append(result, tagStr)
			}
		}
	case string:
		// 处理单个字符串标签
		result = append(result, v)
	}

	return result
}

// extractCategoriesFromYAML 从 YAML 数据中提取分类
func extractCategoriesFromYAML(categories interface{}) []string {
	var result []string

	switch v := categories.(type) {
	case []interface{}:
		for _, category := range v {
			if categoryStr, ok := category.(string); ok {
				result = append(result, categoryStr)
			}
		}
	case string:
		// 处理单个字符串分类
		result = append(result, v)
	}

	return result
}
