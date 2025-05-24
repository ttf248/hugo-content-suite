package operations

import (
	"bufio"
	"fmt"
	"tag-scanner/display"
	"tag-scanner/models"
	"tag-scanner/scanner"
	"tag-scanner/translator"

	"github.com/fatih/color"
)

func (p *Processor) ShowCacheStatus() {
	color.Cyan("=== 翻译缓存状态 ===")

	translatorInstance := translator.NewLLMTranslator()

	fmt.Println()
	fmt.Println(translatorInstance.GetCacheInfo())
	fmt.Println()

	totalCount, expiredCount := translatorInstance.GetCacheStats()
	fmt.Printf("📊 统计信息:\n")
	fmt.Printf("   总翻译条目: %d 个\n", totalCount)
	fmt.Printf("   过期条目: %d 个\n", expiredCount)
	fmt.Printf("   有效条目: %d 个\n", totalCount-expiredCount)
}

func (p *Processor) ClearTranslationCache(reader *bufio.Reader) {
	color.Yellow("⚠️  警告：此操作将清空所有翻译缓存")
	if !p.confirmExecution(reader, "确认清空缓存？(y/n): ") {
		color.Yellow("❌ 已取消清空操作")
		return
	}

	translatorInstance := translator.NewLLMTranslator()
	if err := translatorInstance.ClearCache(); err != nil {
		color.Red("❌ 清空缓存失败: %v", err)
		return
	}

	color.Green("✅ 翻译缓存已清空")
}

func (p *Processor) PreviewBulkTranslationCache(tagStats []models.TagStats) {
	color.Cyan("=== 全量翻译缓存预览 ===")

	cachePreview, err := p.collectTranslationTargets(tagStats)
	if err != nil {
		color.Red("❌ 收集翻译目标失败: %v", err)
		return
	}

	display.DisplayBulkTranslationPreview(cachePreview, 20)
}

func (p *Processor) GenerateBulkTranslationCache(tagStats []models.TagStats, reader *bufio.Reader) {
	color.Cyan("🔍 正在收集翻译目标...")

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

	_, err = translatorInstance.BatchTranslate(cachePreview.MissingTranslations)
	if err != nil {
		color.Red("❌ 批量翻译失败: %v", err)
		return
	}

	color.Green("✅ 全量翻译缓存生成完成！")
}

func (p *Processor) displayCacheStats(cachePreview *display.BulkTranslationPreview) {
	fmt.Printf("\n📊 翻译缓存统计:\n")
	fmt.Printf("   🏷️  标签总数: %d 个\n", cachePreview.TotalTags)
	fmt.Printf("   📝 文章总数: %d 篇\n", cachePreview.TotalArticles)
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

	// 合并所有需要翻译的文本
	allTexts := append(tagNames, articleTitles...)

	// 检查缓存状态
	missingTexts, cachedCount := translatorInstance.PrepareBulkTranslation(allTexts)

	// 分离标签和文章的缺失项
	var tagsToTranslate []display.TranslationItem
	var articlesToTranslate []display.TranslationItem

	for _, text := range missingTexts {
		// 检查是否为标签
		isTag := false
		for _, stat := range tagStats {
			if stat.Name == text {
				tagsToTranslate = append(tagsToTranslate, display.TranslationItem{
					Type:     "标签",
					Original: text,
					Count:    stat.Count,
				})
				isTag = true
				break
			}
		}

		// 如果不是标签，则为文章标题
		if !isTag {
			articlesToTranslate = append(articlesToTranslate, display.TranslationItem{
				Type:     "文章",
				Original: text,
				Count:    1,
			})
		}
	}

	return &display.BulkTranslationPreview{
		TotalTags:           len(tagStats),
		TotalArticles:       len(articleTitles),
		CachedCount:         cachedCount,
		MissingTranslations: missingTexts,
		TagsToTranslate:     tagsToTranslate,
		ArticlesToTranslate: articlesToTranslate,
	}, nil
}
