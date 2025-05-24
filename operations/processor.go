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
	color.Cyan("=== 一键处理全部 ===")
	fmt.Println("这将自动执行以下操作：")
	fmt.Println("1. 生成全量翻译缓存")
	fmt.Println("2. 生成新增标签页面")
	fmt.Println("3. 生成缺失文章Slug")
	fmt.Println("4. 翻译新增文章为英文")
	fmt.Println()

	// 预览批量翻译缓存
	cachePreview := p.PreviewBulkTranslationCache(tagStats)

	// 预览标签页面
	tagGenerator := generator.NewTagPageGenerator(p.contentDir)
	tagPreviews := tagGenerator.PreviewTagPages(tagStats)
	createTagCount := 0
	for _, preview := range tagPreviews {
		if preview.Status == "create" {
			createTagCount++
		}
	}

	// 预览文章Slug
	slugGenerator := generator.NewArticleSlugGenerator(p.contentDir)
	slugPreviews, err := slugGenerator.PreviewArticleSlugs()
	missingSlugCount := 0
	if err == nil {
		for _, preview := range slugPreviews {
			if preview.Status == "missing" {
				missingSlugCount++
			}
		}
	}

	// 预览文章翻译
	articleTranslator := generator.NewArticleTranslator(p.contentDir)
	translationPreviews, err := articleTranslator.PreviewArticleTranslations()
	missingTranslationCount := 0
	if err == nil {
		for _, preview := range translationPreviews {
			if preview.Status == "missing" {
				missingTranslationCount++
			}
		}
	}

	// 显示总体预览
	fmt.Printf("📊 总体预览:\n")
	fmt.Printf("   🔄 需要翻译: %d 个项目\n", len(cachePreview.MissingTranslations))
	fmt.Printf("   🏷️  需要创建标签页面: %d 个\n", createTagCount)
	fmt.Printf("   📝 需要添加文章Slug: %d 个\n", missingSlugCount)
	fmt.Printf("   🌐 需要翻译文章: %d 篇\n", missingTranslationCount)

	totalTasks := 0
	if len(cachePreview.MissingTranslations) > 0 {
		totalTasks++
	}
	if createTagCount > 0 {
		totalTasks++
	}
	if missingSlugCount > 0 {
		totalTasks++
	}
	if missingTranslationCount > 0 {
		totalTasks++
	}

	if totalTasks == 0 {
		color.Green("✅ 所有内容都已是最新状态，无需处理")
		return
	}

	fmt.Printf("\n需要执行 %d 个步骤\n", totalTasks)

	if !p.confirmExecution(reader, "\n⚠️ 确认开始一键处理？(y/n): ") {
		color.Yellow("⏹️ 操作已取消")
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

	// 获取文章翻译预览信息
	articleTranslator := generator.NewArticleTranslator(p.contentDir)
	translationPreviews, err := articleTranslator.PreviewArticleTranslations()
	if err != nil {
		color.Red("❌ 获取文章翻译预览失败: %v", err)
		utils.Error("获取文章翻译预览失败: %v", err)
		return
	}

	// 统计需要翻译的文章数量
	missingTranslationCount := 0
	for _, preview := range translationPreviews {
		if preview.Status == "missing" {
			missingTranslationCount++
		}
	}

	// 步骤1: 生成全量翻译缓存
	if len(cachePreview.MissingTranslations) > 0 {
		if !p.processTranslationCache(cachePreview) {
			return
		}
	} else {
		color.Green("\n✅ 步骤1/4: 翻译缓存已是最新")
	}

	// 步骤2: 生成新增标签页面
	if createTagCount > 0 {
		if !p.processTagPages(tagStats, createTagCount) {
			return
		}
	} else {
		color.Green("\n✅ 步骤2/4: 标签页面已是最新")
	}

	// 步骤3: 生成缺失文章Slug
	if missingSlugCount > 0 {
		if !p.processArticleSlugs(missingSlugCount) {
			return
		}
	} else {
		color.Green("\n✅ 步骤3/4: 文章Slug已是最新")
	}

	// 步骤4: 翻译新增文章为英文
	if missingTranslationCount > 0 {
		if !p.processArticleTranslations(missingTranslationCount) {
			return
		}
	} else {
		color.Green("\n✅ 步骤4/4: 文章翻译已是最新")
	}

	// 显示最终统计
	p.displayFinalStats()
}

func (p *Processor) processTranslationCache(cachePreview *display.BulkTranslationPreview) bool {
	color.Blue("\n📦 步骤1/4: 生成全量翻译缓存")
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
	color.Blue("\n🏷️  步骤2/4: 生成新增标签页面")
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
	color.Blue("\n📝 步骤3/4: 生成缺失文章Slug")
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

func (p *Processor) processArticleTranslations(missingCount int) bool {
	color.Yellow("\n🔄 步骤4/4: 翻译新增文章为英文")
	fmt.Printf("需要翻译 %d 篇文章\n", missingCount)

	utils.Info("开始处理文章翻译，缺失数量: %d", missingCount)

	articleTranslator := generator.NewArticleTranslator(p.contentDir)

	fmt.Print("正在翻译文章...")
	if err := articleTranslator.TranslateArticles("missing"); err != nil {
		fmt.Println()
		color.Red("❌ 文章翻译失败: %v", err)
		utils.Error("文章翻译失败: %v", err)
		return false
	}

	color.Green("✅ 步骤4/4: 文章翻译完成")
	utils.Info("文章翻译处理完成")
	return true
}

func (p *Processor) displayFinalStats() {
	color.Green("\n🎉 一键处理流程完成！")

	fmt.Println("\n📊 处理结果总结:")
	fmt.Println("✅ 翻译缓存已更新")
	fmt.Println("✅ 标签页面已生成")
	fmt.Println("✅ 文章Slug已完善")
	fmt.Println("✅ 文章翻译已完成")

	fmt.Println("\n💡 提示:")
	fmt.Println("   - 所有缓存已更新，后续操作将更加快速")
	fmt.Println("   - 标签页面已生成到 content/tags/ 目录")
	fmt.Println("   - 文章Slug已添加到各文章的front matter")
	fmt.Println("   - 英文版本文章已生成到对应目录")
	fmt.Println("   - 可以使用其他菜单选项进行具体查看和管理")

	utils.Info("一键处理流程全部完成")
}

func (p *Processor) confirmExecution(reader *bufio.Reader, prompt string) bool {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(strings.ToLower(input)) == "y"
}
