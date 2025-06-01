package generator

import (
	"bufio"
	"fmt"
	"hugo-content-suite/models"
	"hugo-content-suite/translator"
	"hugo-content-suite/utils"
	"os"
	"path/filepath"
	"strings"
	"time"
)

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

// 实现 StatusLike 接口
func (t TagPagePreview) GetStatus() string {
	return t.Status
}

// TagPageGenerator 标签页面生成器
type TagPageGenerator struct {
	contentDir       string
	translationUtils *TranslationUtils
	slugCache        map[string]string
}

// NewTagPageGenerator 创建新的标签页面生成器
func NewTagPageGenerator(contentDir string) *TagPageGenerator {
	return &TagPageGenerator{
		contentDir:       contentDir,
		translationUtils: NewTranslationUtils(),
		slugCache:        make(map[string]string),
	}
}

// GenerateTagPagesWithMode 根据模式生成标签页面文件
func (g *TagPageGenerator) GenerateTagPagesWithMode(targetPreviews []TagPagePreview, mode string) error {
	fmt.Println("\n🏷️  标签页面生成器 (模式选择)")
	fmt.Println("===============================")

	if len(targetPreviews) == 0 {
		fmt.Printf("ℹ️  根据选择的模式 '%s'，没有需要处理的标签\n", mode)
		return nil
	}

	fmt.Printf("📊 将处理 %d 个标签 (模式: %s)\n", len(targetPreviews), mode)

	tagsDir := filepath.Join(g.contentDir, "..", "tags")
	if err := utils.EnsureDir(tagsDir); err != nil {
		return fmt.Errorf("❌ 创建tags目录失败: %v", err)
	}

	return g.processTargetPreviews(targetPreviews, tagsDir)
}

// processTargetPreviews 处理目标预览
func (g *TagPageGenerator) processTargetPreviews(targetPreviews []TagPagePreview, tagsDir string) error {
	createdCount := 0
	updatedCount := 0
	errorCount := 0

	fmt.Printf("\n📝 正在生成标签页面...\n")
	fmt.Println("========================")

	for i, preview := range targetPreviews {
		fmt.Printf("  [%d/%d] %s", i+1, len(targetPreviews), preview.TagName)

		tagDir := filepath.Join(tagsDir, preview.TagName)
		indexFile := filepath.Join(tagDir, "_index.md")
		content := g.GenerateTagContent(preview.TagName, preview.Slug)

		if err := utils.WriteFileContent(indexFile, content); err != nil {
			fmt.Printf(" ❌ 失败\n")
			fmt.Printf("     错误: %v\n", err)
			errorCount++
			continue
		}

		if preview.Status == "create" {
			fmt.Printf(" ✨ 新建\n")
			fmt.Printf("     slug: %s\n", preview.Slug)
			createdCount++
		} else {
			fmt.Printf(" 🔄 更新\n")
			fmt.Printf("     slug: %s\n", preview.Slug)
			updatedCount++
		}
	}

	fmt.Printf("\n🎉 标签页面生成完成！\n")
	fmt.Printf("   ✨ 新建: %d 个\n", createdCount)
	fmt.Printf("   🔄 更新: %d 个\n", updatedCount)
	if errorCount > 0 {
		fmt.Printf("   ❌ 失败: %d 个\n", errorCount)
	}
	fmt.Printf("   📦 总计: %d 个\n", len(targetPreviews))

	return nil
}

// PreviewTagPages 预览即将生成的标签页面
func (g *TagPageGenerator) PrepareTagPages(tagStats []models.TagStats) ([]TagPagePreview, int, int) {
	var previews []TagPagePreview

	// 测试LM Studio连接
	fmt.Print("🔗 测试LM Studio连接... ")
	if err := g.translationUtils.TestConnection(); err != nil {
		fmt.Printf("❌ 失败 (%v)\n", err)
		fmt.Println("⚠️  无法连接AI翻译，终止操作")
		return previews, 0, 0
	} else {
		fmt.Println("✅ 成功")
	}

	// 收集所有标签名
	tagNames := make([]string, len(tagStats))
	for i, stat := range tagStats {
		tagNames[i] = stat.Name
	}

	fmt.Printf("🌐 正在生成 %d 个标签的slug...\n", len(tagNames))

	// 只使用AI批量翻译（带缓存）
	fmt.Println("🤖 使用AI翻译...")
	slugMap, err := g.translationUtils.BatchTranslateWithCache(tagNames, "en", translator.TagCache)
	if err != nil {
		fmt.Printf("⚠️ 批量翻译失败: %v\n", err)
		return previews, 0, 0
	}

	// 格式化所有slug
	for tag, slug := range slugMap {
		slugMap[tag] = g.translationUtils.FormatSlugField(slug)
	}

	fmt.Printf("\n📊 正在分析标签状态...\n")
	createCount := 0
	updateCount := 0

	for i, stat := range tagStats {
		fmt.Printf("  [%d/%d] 检查: %s", i+1, len(tagStats), stat.Name)

		var status string

		// 检查标签目录是否已存在
		tagsDir := filepath.Join(g.contentDir, "..", "tags")
		tagDir := filepath.Join(tagsDir, stat.Name)
		indexFile := filepath.Join(tagDir, "_index.md")

		if utils.FileExists(indexFile) {
			status = "update"
			updateCount++
			fmt.Printf(" 🔄 需要更新\n")
		} else {
			status = "create"
			createCount++
			fmt.Printf(" ✨ 需要新建\n")
		}

		// 生成slug（从映射中获取）
		slug := slugMap[stat.Name]

		preview := TagPagePreview{
			TagName:       stat.Name,
			Slug:          slug,
			ArticleCount:  stat.Count,
			DirectoryPath: fmt.Sprintf("tags/%s/", stat.Name),
			FilePath:      fmt.Sprintf("tags/%s/_index.md", stat.Name),
			Status:        status,
			ExistingSlug:  g.ExtractSlugFromFile(indexFile),
		}
		previews = append(previews, preview)

		time.Sleep(10 * time.Millisecond) // 短暂延迟
	}

	fmt.Printf("\n📈 统计结果:\n")
	fmt.Printf("   ✨ 需要新建: %d 个\n", createCount)
	fmt.Printf("   🔄 需要更新: %d 个\n", updateCount)
	fmt.Printf("   📦 总计: %d 个\n", len(previews))

	return previews, createCount, updateCount
}

// GenerateTagContent 生成标签页面内容
func (g *TagPageGenerator) GenerateTagContent(tagName, slug string) string {
	return fmt.Sprintf(`---
title: %s
slug: "%s"
---
`, tagName, slug)
}

// ExtractSlugFromFile 从标签页面文件中提取现有的slug
func (g *TagPageGenerator) ExtractSlugFromFile(filePath string) string {
	if !utils.FileExists(filePath) {
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
