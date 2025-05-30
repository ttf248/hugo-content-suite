package generator

import (
	"bufio"
	"fmt"
	"hugo-content-suite/models"
	"hugo-content-suite/scanner"
	"hugo-content-suite/translator"
	"hugo-content-suite/utils"
	"os"
	"regexp"
	"strings"
)

// ArticleSlugGenerator 文章slug生成器
type ArticleSlugGenerator struct {
	contentDir string
	translator *translator.LLMTranslator
}

// ArticleSlugPreview 文章slug预览信息
type ArticleSlugPreview struct {
	FilePath    string
	Title       string
	CurrentSlug string
	NewSlug     string
	Status      string // "missing", "update", "skip"
}

// NewArticleSlugGenerator 创建新的文章slug生成器
func NewArticleSlugGenerator(contentDir string) *ArticleSlugGenerator {
	return &ArticleSlugGenerator{
		contentDir: contentDir,
		translator: translator.NewLLMTranslator(),
	}
}

// PreviewArticleSlugs 预览文章slug生成
func (g *ArticleSlugGenerator) PreviewArticleSlugs() ([]ArticleSlugPreview, error) {
	articles, err := scanner.ScanArticles(g.contentDir)
	if err != nil {
		return nil, fmt.Errorf("扫描文章失败: %v", err)
	}

	var previews []ArticleSlugPreview

	for _, article := range articles {
		if article.Title == "" {
			continue
		}

		currentSlug := g.extractSlugFromFile(article.FilePath)

		var status string
		if currentSlug == "" {
			status = "missing"
		} else {
			status = "exists"
		}

		preview := ArticleSlugPreview{
			FilePath:    article.FilePath,
			Title:       article.Title,
			CurrentSlug: currentSlug,
			NewSlug:     "[需要生成]", // 简化预览
			Status:      status,
		}

		previews = append(previews, preview)
	}

	return previews, nil
}

// GenerateArticleSlugs 生成文章slug
func (g *ArticleSlugGenerator) GenerateArticleSlugs() error {
	utils.LogOperation("开始生成文章Slug", map[string]interface{}{
		"content_dir": g.contentDir,
	})

	articles, err := scanner.ScanArticles(g.contentDir)
	if err != nil {
		utils.ErrorWithFields("扫描文章失败", map[string]interface{}{
			"content_dir": g.contentDir,
			"error":       err.Error(),
		})
		return fmt.Errorf("扫描文章失败: %v", err)
	}

	// 测试LM Studio连接
	fmt.Println("正在测试与LM Studio的连接...")
	if err := g.translator.TestConnection(); err != nil {
		utils.WarnWithFields("LM Studio连接失败", map[string]interface{}{
			"error": err.Error(),
		})
		fmt.Printf("警告：无法连接到LM Studio (%v)，将使用备用翻译方案\n", err)
	} else {
		utils.InfoWithFields("LM Studio连接成功", map[string]interface{}{
			"status": "connected",
		})
		fmt.Println("LM Studio连接成功！")
	}

	processedCount := 0
	updatedCount := 0
	createdCount := 0
	errorCount := 0

	for i, article := range articles {
		if article.Title == "" {
			continue
		}

		utils.DebugWithFields("处理文章", map[string]interface{}{
			"article_index": i + 1,
			"total_count":   len(articles),
			"title":         article.Title,
			"file_path":     article.FilePath,
		})

		fmt.Printf("处理文章 (%d/%d): %s\n", i+1, len(articles), article.Title)

		// 生成新的slug
		newSlug, err := g.translator.TranslateToSlug(article.Title)
		if err != nil {
			utils.WarnWithFields("翻译失败", map[string]interface{}{
				"title": article.Title,
				"error": err.Error(),
			})
			fmt.Printf("  翻译失败: %v，跳过此文章\n", err)
			errorCount++
			continue
		}

		// 更新文件
		currentSlug := g.extractSlugFromFile(article.FilePath)
		if currentSlug == "" {
			// 添加slug
			if err := g.addSlugToFile(article.FilePath, newSlug); err != nil {
				utils.ErrorWithFields("添加slug失败", map[string]interface{}{
					"file_path": article.FilePath,
					"new_slug":  newSlug,
					"error":     err.Error(),
				})
				fmt.Printf("  添加slug失败: %v\n", err)
				errorCount++
				continue
			}
			utils.InfoWithFields("添加slug成功", map[string]interface{}{
				"file_path": article.FilePath,
				"slug":      newSlug,
			})
			fmt.Printf("  ✓ 添加slug: %s\n", newSlug)
			createdCount++
		} else if currentSlug != newSlug {
			// 更新slug
			if err := g.updateSlugInFile(article.FilePath, currentSlug, newSlug); err != nil {
				utils.ErrorWithFields("更新slug失败", map[string]interface{}{
					"file_path": article.FilePath,
					"old_slug":  currentSlug,
					"new_slug":  newSlug,
					"error":     err.Error(),
				})
				fmt.Printf("  更新slug失败: %v\n", err)
				errorCount++
				continue
			}
			utils.InfoWithFields("更新slug成功", map[string]interface{}{
				"file_path": article.FilePath,
				"old_slug":  currentSlug,
				"new_slug":  newSlug,
			})
			fmt.Printf("  ✓ 更新slug: %s -> %s\n", currentSlug, newSlug)
			updatedCount++
		} else {
			fmt.Printf("  - slug已是最新: %s\n", currentSlug)
		}

		processedCount++
	}

	utils.LogOperation("文章Slug生成完成", map[string]interface{}{
		"processed_count": processedCount,
		"created_count":   createdCount,
		"updated_count":   updatedCount,
		"error_count":     errorCount,
	})

	fmt.Printf("\n文章slug生成完成！\n")
	fmt.Printf("- 处理文章: %d 篇\n", processedCount)
	fmt.Printf("- 新增slug: %d 个\n", createdCount)
	fmt.Printf("- 更新slug: %d 个\n", updatedCount)
	fmt.Printf("- 处理失败: %d 个\n", errorCount)

	return nil
}

// GenerateArticleSlugsWithMode 根据模式生成文章slug
func (g *ArticleSlugGenerator) GenerateArticleSlugsWithMode(mode string) error {
	articles, err := scanner.ScanArticles(g.contentDir)
	if err != nil {
		return fmt.Errorf("扫描文章失败: %v", err)
	}

	// 测试LM Studio连接
	fmt.Println("正在测试与LM Studio的连接...")
	if err := g.translator.TestConnection(); err != nil {
		fmt.Printf("警告：无法连接到LM Studio (%v)，将使用备用翻译方案\n", err)
	} else {
		fmt.Println("LM Studio连接成功！")
	}

	// 根据模式过滤需要处理的文章
	var targetArticles []models.Article
	for _, article := range articles {
		if article.Title == "" {
			continue
		}

		currentSlug := g.extractSlugFromFile(article.FilePath)

		switch mode {
		case "missing":
			if currentSlug == "" {
				targetArticles = append(targetArticles, article)
			}
		case "update":
			if currentSlug != "" {
				targetArticles = append(targetArticles, article)
			}
		case "all":
			targetArticles = append(targetArticles, article)
		}
	}

	if len(targetArticles) == 0 {
		fmt.Println("根据选择的模式，没有需要处理的文章")
		return nil
	}

	// 收集所有需要翻译的标题
	var titlesToTranslate []string
	for _, article := range targetArticles {
		titlesToTranslate = append(titlesToTranslate, article.Title)
	}

	// 批量翻译所有标题（使用文章专用接口）
	fmt.Printf("正在批量翻译 %d 个文章标题...\n", len(titlesToTranslate))
	translationMap, err := g.translator.BatchTranslateArticles(titlesToTranslate)
	if err != nil {
		fmt.Printf("⚠️ 批量翻译失败: %v，将逐个翻译\n", err)
		translationMap = make(map[string]string)
	}

	processedCount := 0
	updatedCount := 0
	createdCount := 0
	skippedCount := 0
	errorCount := 0

	for i, article := range targetArticles {
		fmt.Printf("处理文章 (%d/%d): %s\n", i+1, len(targetArticles), article.Title)

		// 获取翻译结果
		var newSlug string
		if slug, exists := translationMap[article.Title]; exists {
			newSlug = slug
		} else {
			// 如果批量翻译失败，尝试单独翻译（使用文章专用接口）
			slug, err := g.translator.TranslateToArticleSlug(article.Title)
			if err != nil {
				fmt.Printf("  翻译失败: %v，跳过此文章\n", err)
				errorCount++
				continue
			} else {
				newSlug = slug
			}
		}

		// 检查当前slug
		currentSlug := g.extractSlugFromFile(article.FilePath)

		if currentSlug == "" {
			// 添加slug
			if err := g.addSlugToFile(article.FilePath, newSlug); err != nil {
				fmt.Printf("  添加slug失败: %v\n", err)
				errorCount++
				continue
			}
			fmt.Printf("  ✓ 添加slug: %s\n", newSlug)
			createdCount++
		} else if currentSlug != newSlug && (mode == "update" || mode == "all") {
			// 更新slug
			if err := g.updateSlugInFile(article.FilePath, currentSlug, newSlug); err != nil {
				fmt.Printf("  更新slug失败: %v\n", err)
				errorCount++
				continue
			}
			fmt.Printf("  ✓ 更新slug: %s -> %s\n", currentSlug, newSlug)
			updatedCount++
		} else {
			fmt.Printf("  - 跳过: slug已是最新 (%s)\n", currentSlug)
			skippedCount++
		}

		processedCount++
	}

	fmt.Printf("\n文章slug生成完成！\n")
	fmt.Printf("- 处理文章: %d 篇\n", processedCount)
	fmt.Printf("- 新增slug: %d 个\n", createdCount)
	fmt.Printf("- 更新slug: %d 个\n", updatedCount)
	fmt.Printf("- 跳过: %d 个\n", skippedCount)
	fmt.Printf("- 处理失败: %d 个\n", errorCount)

	return nil
}

// extractSlugFromFile 从文件中提取现有的slug
func (g *ArticleSlugGenerator) extractSlugFromFile(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	inFrontMatter := false

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

		if inFrontMatter && strings.HasPrefix(strings.TrimSpace(line), "slug:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				slug := strings.TrimSpace(parts[1])
				slug = strings.Trim(slug, "\"'")
				return slug
			}
		}
	}

	return ""
}

// addSlugToFile 向文件添加slug字段
func (g *ArticleSlugGenerator) addSlugToFile(filePath, slug string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string
	inFrontMatter := false
	frontMatterEnd := -1

	for i, line := range lines {
		if strings.TrimSpace(line) == "---" {
			if !inFrontMatter {
				inFrontMatter = true
			} else {
				frontMatterEnd = i
				break
			}
		}
	}

	if frontMatterEnd == -1 {
		return fmt.Errorf("找不到front matter结束标记")
	}

	// 在front matter结束前添加slug
	for i, line := range lines {
		newLines = append(newLines, line)
		if i == frontMatterEnd-1 {
			newLines = append(newLines, fmt.Sprintf("slug: \"%s\"", slug))
		}
	}

	return os.WriteFile(filePath, []byte(strings.Join(newLines, "\n")), 0644)
}

// updateSlugInFile 更新文件中的slug字段
func (g *ArticleSlugGenerator) updateSlugInFile(filePath, oldSlug, newSlug string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// 使用正则表达式替换slug
	slugPattern := regexp.MustCompile(`slug:\s*["']?` + regexp.QuoteMeta(oldSlug) + `["']?`)
	newContent := slugPattern.ReplaceAllString(string(content), fmt.Sprintf("slug: \"%s\"", newSlug))

	return os.WriteFile(filePath, []byte(newContent), 0644)
}
