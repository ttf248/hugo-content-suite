package generator

import (
	"bufio"
	"fmt"
	"hugo-content-suite/models"
	"hugo-content-suite/translator"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// TagPageGenerator 标签页面生成器
type TagPageGenerator struct {
	contentDir string
	translator *translator.LLMTranslator
	slugCache  map[string]string
}

// NewTagPageGenerator 创建新的标签页面生成器
func NewTagPageGenerator(contentDir string) *TagPageGenerator {
	return &TagPageGenerator{
		contentDir: contentDir,
		translator: translator.NewLLMTranslator(),
		slugCache:  make(map[string]string),
	}
}

// GenerateTagPages 生成标签页面文件
func (g *TagPageGenerator) GenerateTagPages(tagStats []models.TagStats) error {
	// 首先测试连接
	fmt.Println("正在测试与LM Studio的连接...")
	if err := g.translator.TestConnection(); err != nil {
		fmt.Printf("警告：无法连接到LM Studio (%v)，将使用备用翻译方案\n", err)
	} else {
		fmt.Println("LM Studio连接成功！")
	}

	// 确定tags目录路径
	tagsDir := filepath.Join(g.contentDir, "..", "tags")

	// 确保tags目录存在
	if err := os.MkdirAll(tagsDir, 0755); err != nil {
		return fmt.Errorf("创建tags目录失败: %v", err)
	}

	// 批量翻译所有标签
	fmt.Println("正在翻译标签...")
	tagNames := make([]string, len(tagStats))
	for i, stat := range tagStats {
		tagNames[i] = stat.Name
	}

	slugMap, err := g.translator.BatchTranslate(tagNames)
	if err != nil {
		return fmt.Errorf("批量翻译失败: %v", err)
	}

	g.slugCache = slugMap

	createdCount := 0
	updatedCount := 0

	fmt.Println("\n正在生成标签页面文件...")
	for _, stat := range tagStats {
		tagDir := filepath.Join(tagsDir, stat.Name)
		indexFile := filepath.Join(tagDir, "_index.md")

		// 检查文件是否已存在
		exists := false
		if _, err := os.Stat(indexFile); err == nil {
			exists = true
		}

		// 确保标签目录存在
		if err := os.MkdirAll(tagDir, 0755); err != nil {
			return fmt.Errorf("创建标签目录 %s 失败: %v", tagDir, err)
		}

		// 获取翻译后的slug
		slug := g.slugCache[stat.Name]

		// 生成文件内容
		content := fmt.Sprintf(`---
title: %s
slug: "%s"
---
`, stat.Name, slug)

		// 写入文件
		if err := os.WriteFile(indexFile, []byte(content), 0644); err != nil {
			return fmt.Errorf("写入文件 %s 失败: %v", indexFile, err)
		}

		if exists {
			updatedCount++
		} else {
			createdCount++
		}
	}

	fmt.Printf("标签页面生成完成！\n")
	fmt.Printf("- 新建: %d 个\n", createdCount)
	fmt.Printf("- 更新: %d 个\n", updatedCount)
	fmt.Printf("- 总计: %d 个\n", len(tagStats))

	return nil
}

// GenerateTagPagesWithMode 根据模式生成标签页面文件
func (g *TagPageGenerator) GenerateTagPagesWithMode(tagStats []models.TagStats, mode string) error {
	// 首先获取预览信息以确定状态
	previews := g.PreviewTagPages(tagStats)

	// 根据模式过滤需要处理的标签
	var targetPreviews []TagPagePreview
	for _, preview := range previews {
		switch mode {
		case "create":
			if preview.Status == "create" {
				targetPreviews = append(targetPreviews, preview)
			}
		case "update":
			if preview.Status == "update" {
				targetPreviews = append(targetPreviews, preview)
			}
		case "all":
			targetPreviews = append(targetPreviews, preview)
		}
	}

	if len(targetPreviews) == 0 {
		fmt.Println("根据选择的模式，没有需要处理的标签")
		return nil
	}

	// 确定tags目录路径
	tagsDir := filepath.Join(g.contentDir, "..", "tags")
	if err := os.MkdirAll(tagsDir, 0755); err != nil {
		return fmt.Errorf("创建tags目录失败: %v", err)
	}

	createdCount := 0
	updatedCount := 0
	errorCount := 0

	fmt.Println("正在生成标签页面...")

	for i, preview := range targetPreviews {
		fmt.Printf("处理标签 (%d/%d): %s\n", i+1, len(targetPreviews), preview.TagName)

		tagDir := filepath.Join(tagsDir, preview.TagName)
		indexFile := filepath.Join(tagDir, "_index.md")

		// 确保标签目录存在
		if err := os.MkdirAll(tagDir, 0755); err != nil {
			fmt.Printf("  创建目录失败: %v\n", err)
			errorCount++
			continue
		}

		// 生成文件内容
		content := fmt.Sprintf(`---
title: %s
slug: "%s"
---
`, preview.TagName, preview.Slug)

		// 写入文件
		if err := os.WriteFile(indexFile, []byte(content), 0644); err != nil {
			fmt.Printf("  写入文件失败: %v\n", err)
			errorCount++
			continue
		}

		if preview.Status == "create" {
			fmt.Printf("  ✓ 新建: %s\n", preview.Slug)
			createdCount++
		} else {
			fmt.Printf("  ✓ 更新: %s\n", preview.Slug)
			updatedCount++
		}
	}

	fmt.Printf("\n标签页面生成完成！\n")
	fmt.Printf("- 新建: %d 个\n", createdCount)
	fmt.Printf("- 更新: %d 个\n", updatedCount)
	if errorCount > 0 {
		fmt.Printf("- 失败: %d 个\n", errorCount)
	}

	return nil
}

// PreviewTagPages 预览即将生成的标签页面
func (g *TagPageGenerator) PreviewTagPages(tagStats []models.TagStats) []TagPagePreview {
	var previews []TagPagePreview

	fmt.Println("正在生成标签页面预览...")

	// 测试LM Studio连接
	fmt.Print("测试LM Studio连接... ")
	useAI := true
	if err := g.translator.TestConnection(); err != nil {
		fmt.Printf("失败 (%v)，将使用备用翻译\n", err)
		useAI = false
	} else {
		fmt.Println("成功")
	}

	// 收集所有标签名
	tagNames := make([]string, len(tagStats))
	for i, stat := range tagStats {
		tagNames[i] = stat.Name
	}

	// 批量翻译标签（利用缓存）
	var slugMap map[string]string
	var err error

	if useAI {
		slugMap, err = g.translator.BatchTranslateTags(tagNames)
		if err != nil {
			fmt.Printf("⚠️ 批量翻译失败: %v，使用备用方案\n", err)
			useAI = false
		}
	}

	if !useAI {
		// 使用备用翻译
		slugMap = make(map[string]string)
		for _, tag := range tagNames {
			slugMap[tag] = g.fallbackSlug(tag)
		}
	}

	for _, stat := range tagStats {
		var status string

		// 检查标签目录是否已存在
		tagsDir := filepath.Join(g.contentDir, "..", "tags")
		tagDir := filepath.Join(tagsDir, stat.Name)
		indexFile := filepath.Join(tagDir, "_index.md")

		if _, err := os.Stat(indexFile); err == nil {
			status = "update"
		} else {
			status = "create"
		}

		preview := TagPagePreview{
			TagName:       stat.Name,
			Slug:          slugMap[stat.Name],
			ArticleCount:  stat.Count,
			DirectoryPath: fmt.Sprintf("tags/%s/", stat.Name),
			FilePath:      fmt.Sprintf("tags/%s/_index.md", stat.Name),
			Status:        status,
			ExistingSlug:  g.extractSlugFromTagFile(indexFile),
		}
		previews = append(previews, preview)
	}

	return previews
}

// TagPagePreview 标签页面预览信息
type TagPagePreview struct {
	TagName       string
	Slug          string
	ArticleCount  int
	DirectoryPath string
	FilePath      string
	Status        string // "create", "update"
	ExistingSlug  string
}

// extractSlugFromTagFile 从标签页面文件中提取现有的slug
func (g *TagPageGenerator) extractSlugFromTagFile(filePath string) string {
	if _, err := os.Stat(filePath); err != nil {
		return ""
	}

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

// fallbackSlug 备用slug生成方案
func (g *TagPageGenerator) fallbackSlug(tag string) string {
	// 预定义的映射表
	fallbackTranslations := map[string]string{
		"人工智能":       "artificial-intelligence",
		"机器学习":       "machine-learning",
		"深度学习":       "deep-learning",
		"前端开发":       "frontend-development",
		"后端开发":       "backend-development",
		"JavaScript": "javascript",
		"Python":     "python",
		"Go":         "golang",
		"技术":         "technology",
		"教程":         "tutorial",
		"编程":         "programming",
		"开发":         "development",
	}

	if slug, exists := fallbackTranslations[tag]; exists {
		return slug
	}

	// 简单处理
	slug := strings.ToLower(tag)
	slug = strings.ReplaceAll(slug, " ", "-")
	// 移除特殊字符
	reg := regexp.MustCompile(`[^\w\x{4e00}-\x{9fff}\-]`)
	slug = reg.ReplaceAllString(slug, "")
	slug = strings.Trim(slug, "-")

	return slug
}
