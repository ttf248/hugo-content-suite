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
	contentDir       string
	translationUtils *translator.TranslationUtils
}

// ArticleSlugPreview 文章slug预览信息
type ArticleSlugPreview struct {
	FilePath    string
	Title       string
	CurrentSlug string
	NewSlug     string
	Status      string // "missing", "update", "skip"
}

// 实现 StatusLike 接口
func (a ArticleSlugPreview) GetStatus() string {
	if a.Status == "missing" {
		return "create"
	}
	return "update"
}

// NewArticleSlugGenerator 创建新的文章slug生成器
func NewArticleSlugGenerator(contentDir string) *ArticleSlugGenerator {
	return &ArticleSlugGenerator{
		contentDir:       contentDir,
		translationUtils: translator.NewTranslationUtils(),
	}
}

// PrepareArticleSlugs 预处理文章slug生成
func (g *ArticleSlugGenerator) PrepareArticleSlugs() ([]ArticleSlugPreview, int, int, error) {
	var previews []ArticleSlugPreview

	// 扫描文章 - 使用基础扫描函数，不需要内容详情
	articles, err := scanner.ScanArticles(g.contentDir)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("扫描文章失败: %v", err)
	}

	// 测试LM Studio连接
	fmt.Print("🔗 测试LM Studio连接... ")
	if err := g.translationUtils.TestConnection(); err != nil {
		fmt.Printf("❌ 失败 (%v)\n", err)
		fmt.Println("⚠️ 无法连接AI翻译，终止操作")
		return nil, 0, 0, fmt.Errorf("AI翻译连接失败: %v", err)
	} else {
		fmt.Println("✅ 成功")
	}

	// 收集需要处理的文章标题
	var validArticles []models.Article
	var titleList []string
	for _, article := range articles {
		if article.Title != "" {
			validArticles = append(validArticles, article)
			titleList = append(titleList, article.Title)
		}
	}

	if len(titleList) == 0 {
		return previews, 0, 0, nil
	}

	fmt.Printf("🌐 正在生成 %d 个文章的slug...\n", len(titleList))

	// 使用AI批量翻译slug
	slugMap, err := g.translationUtils.TranslateArticlesSlugs(titleList)
	if err != nil {
		fmt.Printf("⚠️ 批量翻译失败: %v\n", err)
		return nil, 0, 0, fmt.Errorf("批量翻译失败: %v", err)
	}

	// 格式化所有slug
	for title, slug := range slugMap {
		slugMap[title] = utils.FormatSlugField(slug)
	}

	fmt.Printf("\n📊 正在分析文章slug状态...\n")
	createCount := 0
	updateCount := 0

	for i, article := range validArticles {
		fmt.Printf("  [%d/%d] 检查: %s", i+1, len(validArticles), article.Title)

		currentSlug := g.extractSlugFromFile(article.FilePath)
		newSlug := slugMap[article.Title]

		var status string
		if currentSlug == "" {
			status = "missing"
			createCount++
			fmt.Printf(" ✨ 需要新建\n")
		} else if currentSlug != newSlug {
			status = "update"
			updateCount++
			fmt.Printf(" 🔄 需要更新\n")
		} else {
			status = "skip"
			fmt.Printf(" ✅ 已是最新\n")
		}

		preview := ArticleSlugPreview{
			FilePath:    article.FilePath,
			Title:       article.Title,
			CurrentSlug: currentSlug,
			NewSlug:     newSlug,
			Status:      status,
		}
		previews = append(previews, preview)
	}

	fmt.Printf("\n📈 统计结果:\n")
	fmt.Printf("   ✨ 需要新建: %d 个\n", createCount)
	fmt.Printf("   🔄 需要更新: %d 个\n", updateCount)
	fmt.Printf("   📦 总计: %d 个\n", len(previews))

	return previews, createCount, updateCount, nil
}

// GenerateArticleSlugsWithMode 根据模式生成文章slug
func (g *ArticleSlugGenerator) GenerateArticleSlugsWithMode(targetPreviews []ArticleSlugPreview, mode string) error {
	fmt.Println("\n📝 文章Slug生成器 (模式选择)")
	fmt.Println("===============================")

	if len(targetPreviews) == 0 {
		fmt.Printf("ℹ️  根据选择的模式 '%s'，没有需要处理的文章\n", mode)
		return nil
	}

	fmt.Printf("📊 将处理 %d 篇文章 (模式: %s)\n", len(targetPreviews), mode)

	return g.processTargetPreviews(targetPreviews)
}

// processTargetPreviews 处理目标预览
func (g *ArticleSlugGenerator) processTargetPreviews(targetPreviews []ArticleSlugPreview) error {
	createdCount := 0
	updatedCount := 0
	errorCount := 0

	fmt.Printf("\n📝 正在生成文章slug...\n")
	fmt.Println("========================")

	for i, preview := range targetPreviews {
		fmt.Printf("  [%d/%d] %s", i+1, len(targetPreviews), preview.Title)

		var err error
		if preview.Status == "missing" {
			err = g.addSlugToFile(preview.FilePath, preview.NewSlug)
			if err == nil {
				fmt.Printf(" ✨ 新建\n")
				fmt.Printf("     slug: %s\n", preview.NewSlug)
				createdCount++
			}
		} else if preview.Status == "update" {
			err = g.updateSlugInFile(preview.FilePath, preview.CurrentSlug, preview.NewSlug)
			if err == nil {
				fmt.Printf(" 🔄 更新\n")
				fmt.Printf("     slug: %s -> %s\n", preview.CurrentSlug, preview.NewSlug)
				updatedCount++
			}
		}

		if err != nil {
			fmt.Printf(" ❌ 失败\n")
			fmt.Printf("     错误: %v\n", err)
			errorCount++
		}
	}

	fmt.Printf("\n🎉 文章slug生成完成！\n")
	fmt.Printf("   ✨ 新建: %d 个\n", createdCount)
	fmt.Printf("   🔄 更新: %d 个\n", updatedCount)
	if errorCount > 0 {
		fmt.Printf("   ❌ 失败: %d 个\n", errorCount)
	}
	fmt.Printf("   📦 总计: %d 个\n", len(targetPreviews))

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
