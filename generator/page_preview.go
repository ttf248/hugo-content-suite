package generator

import (
	"fmt"
	"hugo-content-suite/models"
	"hugo-content-suite/translator"
	"path/filepath"
	"time"
)

// PreviewTagPages 预览即将生成的标签页面
func (g *TagPageGenerator) PreviewTagPages(tagStats []models.TagStats) []TagPagePreview {
	var previews []TagPagePreview

	// 测试LM Studio连接
	fmt.Print("🔗 测试LM Studio连接... ")
	useAI := true
	if err := g.translationUtils.TestConnection(); err != nil {
		fmt.Printf("❌ 失败 (%v)\n", err)
		fmt.Println("⚠️  将使用备用翻译")
		useAI = false
	} else {
		fmt.Println("✅ 成功")
	}

	// 收集所有标签名
	tagNames := make([]string, len(tagStats))
	for i, stat := range tagStats {
		tagNames[i] = stat.Name
	}

	fmt.Printf("🌐 正在生成 %d 个标签的slug...\n", len(tagNames))

	// 生成slug映射
	var slugMap map[string]string
	var err error

	if useAI {
		fmt.Println("🤖 使用AI翻译...")
		// 使用批量翻译（带缓存）
		slugMap, err = g.translationUtils.BatchTranslateWithCache(tagNames, "en", translator.TagCache)
		if err != nil {
			fmt.Printf("⚠️ 批量翻译失败: %v，使用备用方案\n", err)
			useAI = false
		}
	}

	if !useAI {
		fmt.Println("🔄 使用备用翻译...")
		slugMap = make(map[string]string)
		for i, tag := range tagNames {
			fmt.Printf("  [%d/%d] %s -> ", i+1, len(tagNames), tag)
			slugMap[tag] = g.translationUtils.FallbackSlug(tag)
			fmt.Printf("%s\n", slugMap[tag])
			time.Sleep(50 * time.Millisecond) // 短暂延迟让用户看到进度
		}
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

		if g.fileUtils.FileExists(indexFile) {
			status = "update"
			updateCount++
			fmt.Printf(" 🔄 需要更新\n")
		} else {
			status = "create"
			createCount++
			fmt.Printf(" ✨ 需要新建\n")
		}

		preview := TagPagePreview{
			TagName:       stat.Name,
			Slug:          slugMap[stat.Name],
			ArticleCount:  stat.Count,
			DirectoryPath: fmt.Sprintf("tags/%s/", stat.Name),
			FilePath:      fmt.Sprintf("tags/%s/_index.md", stat.Name),
			Status:        status,
			ExistingSlug:  g.fileUtils.ExtractSlugFromFile(indexFile),
		}
		previews = append(previews, preview)

		time.Sleep(10 * time.Millisecond) // 短暂延迟
	}

	fmt.Printf("\n📈 统计结果:\n")
	fmt.Printf("   ✨ 需要新建: %d 个\n", createCount)
	fmt.Printf("   🔄 需要更新: %d 个\n", updateCount)
	fmt.Printf("   📦 总计: %d 个\n", len(previews))

	return previews
}
