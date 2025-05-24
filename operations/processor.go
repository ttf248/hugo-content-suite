package operations

import (
	"bufio"
	"fmt"
	"strings"
	"tag-scanner/display"
	"tag-scanner/generator"
	"tag-scanner/models"
	"tag-scanner/translator"
	"tag-scanner/utils"

	"github.com/fatih/color"
)

type Processor struct {
	contentDir string
}

func NewProcessor(contentDir string) *Processor {
	return &Processor{
		contentDir: contentDir,
	}
}

func (p *Processor) QuickProcessAll(tagStats []models.TagStats, reader *bufio.Reader) {
	color.Cyan("=== 🚀 一键快速处理 ===")
	fmt.Println()
	color.Yellow("此操作将按顺序执行以下步骤：")
	fmt.Println("1. 📦 生成全量翻译缓存")
	fmt.Println("2. 🏷️  生成新增标签页面")
	fmt.Println("3. 📝 生成缺失文章Slug")
	fmt.Println()

	// 显示预览统计
	fmt.Println("🔍 正在分析当前状态...")

	// 分析当前状态
	cachePreview, createTagCount, missingSlugCount, err := p.analyzeCurrentState(tagStats)
	if err != nil {
		color.Red("❌ 分析失败: %v", err)
		return
	}

	// 显示统计信息
	totalOperations := len(cachePreview.MissingTranslations) + createTagCount + missingSlugCount
	p.displayProcessStats(cachePreview, createTagCount, missingSlugCount, totalOperations)

	if totalOperations == 0 {
		color.Green("✅ 所有内容都已是最新状态，无需处理")
		return
	}

	// 确认执行
	if !p.confirmExecution(reader, "确认开始一键处理？(y/n): ") {
		color.Yellow("❌ 已取消一键处理")
		return
	}

	// 执行处理流程
	p.executeProcessFlow(cachePreview, createTagCount, missingSlugCount, tagStats)
}

func (p *Processor) analyzeCurrentState(tagStats []models.TagStats) (*display.BulkTranslationPreview, int, int, error) {
	// 分析翻译缓存状态
	cachePreview, err := p.collectTranslationTargets(tagStats)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("分析翻译缓存失败: %v", err)
	}

	// 分析标签页面状态
	pageGenerator := generator.NewTagPageGenerator(p.contentDir)
	tagPreviews := pageGenerator.PreviewTagPages(tagStats)
	createTagCount := 0
	for _, preview := range tagPreviews {
		if preview.Status == "create" {
			createTagCount++
		}
	}

	// 分析文章Slug状态
	slugGenerator := generator.NewArticleSlugGenerator(p.contentDir)
	slugPreviews, err := slugGenerator.PreviewArticleSlugs()
	if err != nil {
		return nil, 0, 0, fmt.Errorf("分析文章Slug失败: %v", err)
	}

	missingSlugCount := 0
	for _, preview := range slugPreviews {
		if preview.Status == "missing" {
			missingSlugCount++
		}
	}

	return cachePreview, createTagCount, missingSlugCount, nil
}

func (p *Processor) displayProcessStats(cachePreview *display.BulkTranslationPreview, createTagCount, missingSlugCount, totalOperations int) {
	fmt.Printf("\n📊 处理统计预览:\n")
	fmt.Printf("   💾 需要翻译: %d 个内容\n", len(cachePreview.MissingTranslations))
	fmt.Printf("   🏷️  需要新建标签页面: %d 个\n", createTagCount)
	fmt.Printf("   📝 需要新增文章Slug: %d 个\n", missingSlugCount)
	fmt.Printf("   📦 预计总操作数: %d 个\n", totalOperations)
	fmt.Println()
	color.Yellow("⚠️  注意：此操作可能需要较长时间，建议在网络稳定时执行")
}

func (p *Processor) executeProcessFlow(cachePreview *display.BulkTranslationPreview, createTagCount, missingSlugCount int, tagStats []models.TagStats) {
	fmt.Println()
	color.Cyan("🚀 开始一键处理流程...")
	utils.Info("开始一键处理流程")

	// 步骤1: 生成全量翻译缓存
	if len(cachePreview.MissingTranslations) > 0 {
		if !p.processTranslationCache(cachePreview) {
			return
		}
	} else {
		color.Green("\n✅ 步骤1/3: 翻译缓存已是最新")
	}

	// 步骤2: 生成新增标签页面
	if createTagCount > 0 {
		if !p.processTagPages(tagStats, createTagCount) {
			return
		}
	} else {
		color.Green("\n✅ 步骤2/3: 标签页面已是最新")
	}

	// 步骤3: 生成缺失文章Slug
	if missingSlugCount > 0 {
		if !p.processArticleSlugs(missingSlugCount) {
			return
		}
	} else {
		color.Green("\n✅ 步骤3/3: 文章Slug已是最新")
	}

	// 显示最终统计
	p.displayFinalStats()
}

func (p *Processor) processTranslationCache(cachePreview *display.BulkTranslationPreview) bool {
	color.Blue("\n📦 步骤1/3: 生成全量翻译缓存")
	fmt.Printf("需要翻译 %d 个内容...\n", len(cachePreview.MissingTranslations))

	translatorInstance := translator.NewLLMTranslator()

	// 分别处理标签和文章翻译
	if len(cachePreview.TagsToTranslate) > 0 {
		fmt.Printf("  🏷️ 翻译 %d 个标签...\n", len(cachePreview.TagsToTranslate))
		tagNames := make([]string, len(cachePreview.TagsToTranslate))
		for i, item := range cachePreview.TagsToTranslate {
			tagNames[i] = item.Original
		}
		_, err := translatorInstance.BatchTranslateTags(tagNames)
		if err != nil {
			color.Red("❌ 标签翻译失败: %v", err)
			return false
		}
	}

	if len(cachePreview.ArticlesToTranslate) > 0 {
		fmt.Printf("  📝 翻译 %d 个文章标题...\n", len(cachePreview.ArticlesToTranslate))
		articleTitles := make([]string, len(cachePreview.ArticlesToTranslate))
		for i, item := range cachePreview.ArticlesToTranslate {
			articleTitles[i] = item.Original
		}
		_, err := translatorInstance.BatchTranslateArticles(articleTitles)
		if err != nil {
			color.Red("❌ 文章翻译失败: %v", err)
			return false
		}
	}

	color.Green("✅ 翻译缓存生成完成")
	return true
}

func (p *Processor) processTagPages(tagStats []models.TagStats, createTagCount int) bool {
	color.Blue("\n🏷️  步骤2/3: 生成新增标签页面")
	fmt.Printf("需要创建 %d 个标签页面...\n", createTagCount)

	pageGenerator := generator.NewTagPageGenerator(p.contentDir)
	err := pageGenerator.GenerateTagPagesWithMode(tagStats, "create")
	if err != nil {
		color.Red("❌ 标签页面生成失败: %v", err)
		return false
	}
	color.Green("✅ 标签页面生成完成")
	return true
}

func (p *Processor) processArticleSlugs(missingSlugCount int) bool {
	color.Blue("\n📝 步骤3/3: 生成缺失文章Slug")
	fmt.Printf("需要添加 %d 个文章Slug...\n", missingSlugCount)

	slugGenerator := generator.NewArticleSlugGenerator(p.contentDir)
	err := slugGenerator.GenerateArticleSlugsWithMode("missing")
	if err != nil {
		color.Red("❌ 文章Slug生成失败: %v", err)
		return false
	}
	color.Green("✅ 文章Slug生成完成")
	return true
}

func (p *Processor) displayFinalStats() {
	fmt.Println()
	color.Green("🎉 一键处理完成！")

	// 显示性能统计
	perfStats := utils.GetGlobalStats()
	if perfStats.TranslationCount > 0 || perfStats.FileOperations > 0 {
		fmt.Println()
		color.Cyan("📊 本次处理统计:")
		fmt.Println(perfStats.String())
	}

	utils.Info("一键处理流程完成")
	fmt.Println()
}

func (p *Processor) confirmExecution(reader *bufio.Reader, prompt string) bool {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(strings.ToLower(input)) == "y"
}
