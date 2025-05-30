package generator

import (
	"fmt"
	"hugo-content-suite/models"
	"hugo-content-suite/translator"
	"path/filepath"
	"time"
)

// TagPageGenerator 标签页面生成器
type TagPageGenerator struct {
	contentDir       string
	translationUtils *TranslationUtils
	fileUtils        *FileUtils
	slugCache        map[string]string
}

// NewTagPageGenerator 创建新的标签页面生成器
func NewTagPageGenerator(contentDir string) *TagPageGenerator {
	return &TagPageGenerator{
		contentDir:       contentDir,
		translationUtils: NewTranslationUtils(),
		fileUtils:        NewFileUtils(),
		slugCache:        make(map[string]string),
	}
}

// GenerateTagPages 生成标签页面文件
func (g *TagPageGenerator) GenerateTagPages(tagStats []models.TagStats) error {
	fmt.Println("\n🏷️  标签页面生成器")
	fmt.Println("==================")

	fmt.Print("🔗 测试LM Studio连接... ")
	useAI := true
	if err := g.translationUtils.TestConnection(); err != nil {
		fmt.Printf("❌ 失败 (%v)\n", err)
		fmt.Println("⚠️  将使用备用翻译方案")
		useAI = false
	} else {
		fmt.Println("✅ 成功")
	}

	tagsDir := filepath.Join(g.contentDir, "..", "tags")
	fmt.Printf("📁 确保目录存在: %s\n", tagsDir)
	if err := g.fileUtils.EnsureDir(tagsDir); err != nil {
		return fmt.Errorf("❌ 创建tags目录失败: %v", err)
	}

	// 批量翻译所有标签
	fmt.Printf("\n🌐 正在翻译 %d 个标签...\n", len(tagStats))
	tagNames := make([]string, len(tagStats))
	for i, stat := range tagStats {
		tagNames[i] = stat.Name
	}

	var slugMap map[string]string
	var err error

	if useAI {
		// 使用带缓存的批量翻译
		slugMap, err = g.translationUtils.BatchTranslateWithCache(tagNames, "en", translator.TagCache)
		if err != nil {
			fmt.Printf("⚠️ 翻译失败: %v，使用备用方案\n", err)
			useAI = false
		}
	}

	if !useAI {
		// 使用原文作为备用方案
		fmt.Println("🔄 使用原文作为slug...")
		slugMap = make(map[string]string)
		for i, tag := range tagNames {
			fmt.Printf("  [%d/%d] %s -> ", i+1, len(tagNames), tag)
			slug := g.translationUtils.FormatSlugField(tag)
			slugMap[tag] = slug
			fmt.Printf("%s\n", slug)
			time.Sleep(10 * time.Millisecond)
		}
	}

	// 格式化所有slug
	for tag, slug := range slugMap {
		slugMap[tag] = g.translationUtils.FormatSlugField(slug)
	}

	g.slugCache = slugMap

	return g.generateTagFiles(tagStats, tagsDir)
}

// GenerateTagPagesWithMode 根据模式生成标签页面文件
func (g *TagPageGenerator) GenerateTagPagesWithMode(tagStats []models.TagStats, mode string) error {
	fmt.Println("\n🏷️  标签页面生成器 (模式选择)")
	fmt.Println("===============================")

	fmt.Print("🔍 生成预览信息... ")
	previews := g.PreviewTagPages(tagStats)
	fmt.Printf("完成 (%d 个标签)\n", len(previews))

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
		fmt.Printf("ℹ️  根据选择的模式 '%s'，没有需要处理的标签\n", mode)
		return nil
	}

	fmt.Printf("📊 将处理 %d 个标签 (模式: %s)\n", len(targetPreviews), mode)

	tagsDir := filepath.Join(g.contentDir, "..", "tags")
	if err := g.fileUtils.EnsureDir(tagsDir); err != nil {
		return fmt.Errorf("❌ 创建tags目录失败: %v", err)
	}

	return g.processTargetPreviews(targetPreviews, tagsDir)
}

// generateTagFiles 生成标签文件
func (g *TagPageGenerator) generateTagFiles(tagStats []models.TagStats, tagsDir string) error {
	createdCount := 0
	updatedCount := 0

	fmt.Printf("\n📝 正在生成标签页面文件...\n")
	fmt.Println("================================")

	for i, stat := range tagStats {
		fmt.Printf("  [%d/%d] 处理标签: %s", i+1, len(tagStats), stat.Name)

		tagDir := filepath.Join(tagsDir, stat.Name)
		indexFile := filepath.Join(tagDir, "_index.md")

		exists := g.fileUtils.FileExists(indexFile)
		slug := g.slugCache[stat.Name]
		content := g.fileUtils.GenerateTagContent(stat.Name, slug)

		if err := g.fileUtils.WriteFileContent(indexFile, content); err != nil {
			fmt.Printf(" ❌ 失败\n")
			fmt.Printf("     错误: %v\n", err)
			return fmt.Errorf("写入文件 %s 失败: %v", indexFile, err)
		}

		if exists {
			fmt.Printf(" 🔄 更新\n")
			updatedCount++
		} else {
			fmt.Printf(" ✨ 新建\n")
			createdCount++
		}

		// 显示slug信息
		fmt.Printf("     slug: %s\n", slug)
	}

	fmt.Printf("\n🎉 标签页面生成完成！\n")
	fmt.Printf("   ✨ 新建: %d 个\n", createdCount)
	fmt.Printf("   🔄 更新: %d 个\n", updatedCount)
	fmt.Printf("   📦 总计: %d 个\n", len(tagStats))

	return nil
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
		content := g.fileUtils.GenerateTagContent(preview.TagName, preview.Slug)

		if err := g.fileUtils.WriteFileContent(indexFile, content); err != nil {
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
