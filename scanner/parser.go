package scanner

import (
	"bufio"
	"hugo-content-suite/models"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Article 类型别名，方便引用
type Article = models.Article

// 默认只扫描 index.md
func ScanArticles(dir string) ([]Article, error) {
	return ScanArticlesWithLangs(dir, false)
}

// 支持 allLangs 参数
func ScanArticlesWithLangs(dir string, allLangs bool) ([]Article, error) {
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

		article, err := parseMarkdownFile(path)
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

func parseMarkdownFile(filePath string) (*Article, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	inFrontMatter := false
	frontMatterLines := []string{}

	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == "---" {
			if !inFrontMatter {
				inFrontMatter = true
				continue
			} else {
				break
			}
		}

		if inFrontMatter {
			frontMatterLines = append(frontMatterLines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	article := &Article{
		FilePath: filePath,
	}

	for _, line := range frontMatterLines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "title:") {
			article.Title = extractValue(line)
		} else if strings.HasPrefix(line, "tags:") {
			article.Tags = extractTags(line, frontMatterLines)
		} else if strings.HasPrefix(line, "categories:") {
			categories := extractCategories(line, frontMatterLines)
			if len(categories) > 0 {
				article.Category = categories[0]
			}
		} else if strings.HasPrefix(line, "date:") {
			article.Date = extractValue(line)
		}
	}

	return article, nil
}

func extractValue(line string) string {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return ""
	}
	value := strings.TrimSpace(parts[1])
	if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
		value = value[1 : len(value)-1]
	}
	return value
}

func extractTags(line string, frontMatterLines []string) []string {
	if strings.Contains(line, "[") {
		re := regexp.MustCompile(`\[(.*?)\]`)
		matches := re.FindStringSubmatch(line)
		if len(matches) > 1 {
			tagsStr := matches[1]
			tags := strings.Split(tagsStr, ",")
			var result []string
			for _, tag := range tags {
				tag = strings.TrimSpace(tag)
				tag = strings.Trim(tag, "\"")
				if tag != "" {
					result = append(result, tag)
				}
			}
			return result
		}
	}

	var tags []string
	inTagsArray := false

	for _, fmLine := range frontMatterLines {
		fmLine = strings.TrimSpace(fmLine)

		if strings.HasPrefix(fmLine, "tags:") {
			inTagsArray = true
			continue
		}

		if inTagsArray {
			if strings.HasPrefix(fmLine, "-") {
				tag := strings.TrimSpace(strings.TrimPrefix(fmLine, "-"))
				tag = strings.Trim(tag, "\"")
				if tag != "" {
					tags = append(tags, tag)
				}
			} else if !strings.HasPrefix(fmLine, " ") && fmLine != "" {
				break
			}
		}
	}

	return tags
}

func extractCategories(line string, frontMatterLines []string) []string {
	if strings.Contains(line, "[") {
		re := regexp.MustCompile(`\[(.*?)\]`)
		matches := re.FindStringSubmatch(line)
		if len(matches) > 1 {
			categoriesStr := matches[1]
			categories := strings.Split(categoriesStr, ",")
			var result []string
			for _, category := range categories {
				category = strings.TrimSpace(category)
				category = strings.Trim(category, "\"")
				if category != "" {
					result = append(result, category)
				}
			}
			return result
		}
	}

	var categories []string
	inCategoriesArray := false

	for _, fmLine := range frontMatterLines {
		fmLine = strings.TrimSpace(fmLine)

		if strings.HasPrefix(fmLine, "categories:") {
			inCategoriesArray = true
			continue
		}

		if inCategoriesArray {
			if strings.HasPrefix(fmLine, "-") {
				category := strings.TrimSpace(strings.TrimPrefix(fmLine, "-"))
				category = strings.Trim(category, "\"")
				if category != "" {
					categories = append(categories, category)
				}
			} else if !strings.HasPrefix(fmLine, " ") && fmLine != "" {
				break
			}
		}
	}

	return categories
}
