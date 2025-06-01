package operations

import (
	"bufio"
	"fmt"
	"hugo-content-suite/display"
	"hugo-content-suite/models"
	"hugo-content-suite/scanner"
	"hugo-content-suite/translator"
	"hugo-content-suite/utils"
	"strings"

	"github.com/fatih/color"
)

func (p *Processor) ShowCacheStatus() {
	color.Cyan("=== 翻译缓存状态 ===")

	translatorInstance := translator.NewLLMTranslator()

	fmt.Println()
	fmt.Println(translatorInstance.GetCacheInfo())
	fmt.Println()

	totalCount := translatorInstance.GetCacheStats()
	fmt.Printf("📊 统计信息:\n")
	fmt.Printf("   总翻译条目: %d 个\n", totalCount)
}

func (p *Processor) ClearTranslationCache(reader *bufio.Reader) {
	color.Cyan("=== 清空翻译缓存 ===")
	fmt.Println("请选择要清空的缓存类型：")
	fmt.Println("1. 清空标签缓存")
	fmt.Println("2. 清空文章缓存")
	fmt.Println("3. 清空所有缓存")
	fmt.Println("0. 取消操作")

	choice := p.getChoice(reader, "请选择 (0-3): ")

	translatorInstance := translator.NewLLMTranslator()

	switch choice {
	case "1":
		if p.confirmExecution(reader, "⚠️ 确认清空标签缓存？(y/n): ") {
			utils.LogOperation("清空标签缓存", map[string]interface{}{
				"operation_type": "cache_clear",
				"cache_type":     "tag",
			})

			if err := translatorInstance.ClearTagCache(); err != nil {
				utils.ErrorWithFields("清空标签缓存失败", map[string]interface{}{
					"error": err.Error(),
				})
				color.Red("❌ 清空标签缓存失败: %v", err)
			} else {
				utils.InfoWithFields("标签缓存清空成功", map[string]interface{}{
					"operation": "cache_clear_tag",
				})
				color.Green("✅ 标签缓存已清空")
			}
		}
	case "2":
		if p.confirmExecution(reader, "⚠️ 确认清空文章缓存？(y/n): ") {
			utils.LogOperation("清空文章缓存", map[string]interface{}{
				"operation_type": "cache_clear",
				"cache_type":     "article",
			})

			if err := translatorInstance.ClearArticleCache(); err != nil {
				utils.ErrorWithFields("清空文章缓存失败", map[string]interface{}{
					"error": err.Error(),
				})
				color.Red("❌ 清空文章缓存失败: %v", err)
			} else {
				utils.InfoWithFields("文章缓存清空成功", map[string]interface{}{
					"operation": "cache_clear_article",
				})
				color.Green("✅ 文章缓存已清空")
			}
		}
	case "3":
		if p.confirmExecution(reader, "⚠️ 确认清空所有缓存？(y/n): ") {
			utils.LogOperation("清空所有缓存", map[string]interface{}{
				"operation_type": "cache_clear",
				"cache_type":     "all",
			})

			if err := translatorInstance.ClearCache(); err != nil {
				utils.ErrorWithFields("清空所有缓存失败", map[string]interface{}{
					"error": err.Error(),
				})
				color.Red("❌ 清空缓存失败: %v", err)
			} else {
				utils.InfoWithFields("所有缓存清空成功", map[string]interface{}{
					"operation": "cache_clear_all",
				})
				color.Green("✅ 所有缓存已清空")
			}
		}
	case "0":
		color.Yellow("❌ 已取消操作")
	default:
		color.Red("⚠️ 无效选择")
	}
}

func (p *Processor) PreviewBulkTranslationCache(tagStats []models.TagStats) *display.BulkTranslationPreview {
	cachePreview, err := p.collectTranslationTargets(tagStats)
	if err != nil {
		color.Red("❌ 收集翻译目标失败: %v", err)
		// 返回空的预览结构而不是nil，避免程序崩溃
		return &display.BulkTranslationPreview{
			TotalTags:           0,
			TotalSlugs:          0,
			CachedCount:         0,
			MissingTranslations: []string{},
			TagsToTranslate:     []display.TranslationItem{},
			SlugsToTranslate:    []display.TranslationItem{},
		}
	}

	return cachePreview
}

func (p *Processor) GenerateBulkTranslationCache(tagStats []models.TagStats, reader *bufio.Reader) {
	color.Cyan("🔍 正在分析翻译需求...")

	cachePreview, err := p.collectTranslationTargets(tagStats)
	if err != nil {
		color.Red("❌ 收集翻译目标失败: %v", err)
		return
	}

	if len(cachePreview.MissingTranslations) == 0 {
		color.Green("✅ 所有内容都已有翻译缓存")
		return
	}

	p.displayCacheStats(cachePreview)

	if !p.confirmExecution(reader, "\n确认生成全量翻译缓存？(y/n): ") {
		color.Yellow("❌ 已取消生成")
		return
	}

	color.Cyan("🚀 正在生成全量翻译缓存...")
	translatorInstance := translator.NewLLMTranslator()

	// 分别批量翻译标签和文章
	if len(cachePreview.TagsToTranslate) > 0 {
		tagNames := make([]string, len(cachePreview.TagsToTranslate))
		for i, item := range cachePreview.TagsToTranslate {
			tagNames[i] = item.Original
		}
		_, err = translatorInstance.BatchTranslateTags(tagNames)
		if err != nil {
			color.Red("❌ 标签批量翻译失败: %v", err)
			return
		}
	}

	if len(cachePreview.SlugsToTranslate) > 0 {
		articleTitles := make([]string, len(cachePreview.SlugsToTranslate))
		for i, item := range cachePreview.SlugsToTranslate {
			articleTitles[i] = item.Original
		}
		_, err = translatorInstance.BatchTranslateSlugs(articleTitles)
		if err != nil {
			color.Red("❌ Slug批量翻译失败: %v", err)
			return
		}
	}

	color.Green("✅ 全量翻译缓存生成完成！")
}

func (p *Processor) displayCacheStats(cachePreview *display.BulkTranslationPreview) {
	fmt.Printf("\n📊 翻译缓存统计:\n")
	fmt.Printf("   🏷️  标签总数: %d 个\n", cachePreview.TotalTags)
	fmt.Printf("   📝 Slug总数: %d 篇\n", cachePreview.TotalSlugs)
	fmt.Printf("   ✅ 已缓存: %d 个\n", cachePreview.CachedCount)
	fmt.Printf("   🔄 需翻译: %d 个\n", len(cachePreview.MissingTranslations))
}

func (p *Processor) collectTranslationTargets(tagStats []models.TagStats) (*display.BulkTranslationPreview, error) {
	translatorInstance := translator.NewLLMTranslator()

	// 收集所有标签
	tagNames := make([]string, len(tagStats))
	for i, stat := range tagStats {
		tagNames[i] = stat.Name
	}

	// 收集所有文章标题
	articles, err := scanner.ScanArticles(p.contentDir)
	if err != nil {
		return nil, fmt.Errorf("扫描文章失败: %v", err)
	}

	var articleTitles []string
	for _, article := range articles {
		if article.Title != "" {
			articleTitles = append(articleTitles, article.Title)
		}
	}

	// 分别检查标签和文章的缓存状态
	missingTags := translatorInstance.GetMissingTags(tagNames)
	missingArticles := translatorInstance.GetMissingArticles(articleTitles)

	// 合并所有缺失的文本
	allMissingTexts := append(missingTags, missingArticles...)
	cachedCount := len(tagNames) + len(articleTitles) - len(allMissingTexts)

	// 分离标签和文章的缺失项
	var tagsToTranslate []display.TranslationItem
	var articlesToTranslate []display.TranslationItem

	for _, tag := range missingTags {
		for _, stat := range tagStats {
			if stat.Name == tag {
				tagsToTranslate = append(tagsToTranslate, display.TranslationItem{
					Type:     "标签",
					Original: tag,
					Count:    stat.Count,
				})
				break
			}
		}
	}

	for _, title := range missingArticles {
		articlesToTranslate = append(articlesToTranslate, display.TranslationItem{
			Type:     "文章",
			Original: title,
			Count:    1,
		})
	}

	return &display.BulkTranslationPreview{
		TotalTags:           len(tagStats),
		TotalSlugs:          len(articleTitles),
		CachedCount:         cachedCount,
		MissingTranslations: allMissingTexts,
		TagsToTranslate:     tagsToTranslate,
		SlugsToTranslate:    articlesToTranslate,
	}, nil
}

func (p *Processor) getChoice(reader *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
