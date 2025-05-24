package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"tag-scanner/config"
	"tag-scanner/display"
	"tag-scanner/generator"
	"tag-scanner/models"
	"tag-scanner/scanner"
	"tag-scanner/stats"
	"tag-scanner/translator"
	"tag-scanner/utils"

	"github.com/fatih/color"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("配置加载失败:", err)
	}

	// 初始化日志
	if err := utils.InitLogger("tag-scanner.log", utils.INFO); err != nil {
		log.Printf("日志初始化失败: %v", err)
	}

	utils.Info("程序启动")
	defer utils.Info("程序退出")

	contentDir := cfg.Paths.DefaultContentDir
	if len(os.Args) > 1 {
		contentDir = os.Args[1]
	}

	articles, err := scanner.ScanArticles(contentDir)
	if err != nil {
		log.Fatal(err)
	}

	if len(articles) == 0 {
		fmt.Println("未找到任何文章")
		return
	}

	// 计算统计数据
	tagStats := stats.CalculateTagStats(articles)
	categoryStats := stats.CalculateCategoryStats(articles)
	noTagArticles := stats.FindNoTagArticles(articles)

	// 显示概览
	display.DisplaySummary(len(articles), tagStats, categoryStats)

	// 显示标签统计（前20个）
	display.DisplayTagStats(tagStats, 20)

	// 显示分类统计
	display.DisplayCategoryStats(categoryStats)

	// 显示无标签文章（前10篇）
	display.DisplayNoTagArticles(noTagArticles, 10)

	// 交互式菜单
	showInteractiveMenu(tagStats, categoryStats, noTagArticles, contentDir)

	// 显示性能统计
	defer func() {
		stats := utils.GetGlobalStats()
		if stats.TranslationCount > 0 || stats.FileOperations > 0 {
			fmt.Println()
			fmt.Println(stats.String())
		}
	}()
}

func showInteractiveMenu(tagStats []models.TagStats, categoryStats []models.CategoryStats, noTagArticles []models.Article, contentDir string) {
	reader := bufio.NewReader(os.Stdin)

	for {
		color.Cyan("\n=== 博客标签管理工具 ===")
		fmt.Println()

		// 快速处理模块
		color.Red("🚀 快速处理")
		fmt.Println("  1. 一键处理全部 (缓存→标签页面→文章Slug)")
		fmt.Println()

		// 数据查看模块
		color.Green("📊 数据查看")
		fmt.Println("  2. 标签统计与分析")
		fmt.Println("  3. 分类统计")
		fmt.Println("  4. 无标签文章")
		fmt.Println()

		// 页面生成模块
		color.Yellow("🏷️  标签页面管理")
		fmt.Println("  5. 预览标签页面")
		fmt.Println("  6. 生成标签页面")
		fmt.Println()

		// 文章管理模块
		color.Blue("📝 文章Slug管理")
		fmt.Println("  7. 预览文章Slug")
		fmt.Println("  8. 生成文章Slug")
		fmt.Println()

		// 缓存管理模块
		color.Magenta("💾 缓存管理")
		fmt.Println("  9. 查看缓存状态")
		fmt.Println(" 10. 预览全量翻译缓存")
		fmt.Println(" 11. 生成全量翻译缓存")
		fmt.Println(" 12. 清空翻译缓存")
		fmt.Println()

		// 系统工具模块
		color.Cyan("🔧 系统工具")
		fmt.Println(" 13. 查看性能统计")
		fmt.Println(" 14. 重置性能统计")
		fmt.Println()

		color.Red("  0. 退出程序")
		fmt.Println()
		fmt.Print("请选择功能 (0-14): ")

		input, _ := reader.ReadString('\n')
		choice := strings.TrimSpace(input)

		switch choice {
		case "1":
			quickProcessAll(tagStats, contentDir, reader)
		case "2":
			showTagAnalysisMenu(tagStats, reader)
		case "3":
			display.DisplayCategoryStats(categoryStats)
		case "4":
			display.DisplayNoTagArticles(noTagArticles, len(noTagArticles))
		case "5":
			previewTagPages(tagStats, contentDir)
		case "6":
			generateTagPages(tagStats, contentDir, reader)
		case "7":
			previewArticleSlugs(contentDir)
		case "8":
			generateArticleSlugs(contentDir, reader)
		case "9":
			showCacheStatus()
		case "10":
			previewBulkTranslationCache(tagStats, contentDir)
		case "11":
			generateBulkTranslationCache(tagStats, contentDir, reader)
		case "12":
			clearTranslationCache(reader)
		case "13":
			showPerformanceStats()
		case "14":
			resetPerformanceStats(reader)
		case "0":
			color.Green("感谢使用！再见！")
			return
		default:
			color.Red("⚠️  无效选择，请重新输入")
		}
	}
}

func quickProcessAll(tagStats []models.TagStats, contentDir string, reader *bufio.Reader) {
	color.Cyan("=== 🚀 一键快速处理 ===")
	fmt.Println()
	color.Yellow("此操作将按顺序执行以下步骤：")
	fmt.Println("1. 📦 生成全量翻译缓存")
	fmt.Println("2. 🏷️  生成新增标签页面")
	fmt.Println("3. 📝 生成缺失文章Slug")
	fmt.Println()

	// 显示预览统计
	fmt.Println("🔍 正在分析当前状态...")

	// 步骤1: 分析翻译缓存状态
	cachePreview, err := collectTranslationTargets(tagStats, contentDir)
	if err != nil {
		color.Red("❌ 分析翻译缓存失败: %v", err)
		return
	}

	// 步骤2: 分析标签页面状态
	pageGenerator := generator.NewTagPageGenerator(contentDir)
	tagPreviews := pageGenerator.PreviewTagPages(tagStats)
	createTagCount := 0
	for _, preview := range tagPreviews {
		if preview.Status == "create" {
			createTagCount++
		}
	}

	// 步骤3: 分析文章Slug状态
	slugGenerator := generator.NewArticleSlugGenerator(contentDir)
	slugPreviews, err := slugGenerator.PreviewArticleSlugs()
	if err != nil {
		color.Red("❌ 分析文章Slug失败: %v", err)
		return
	}

	missingSlugCount := 0
	for _, preview := range slugPreviews {
		if preview.Status == "missing" {
			missingSlugCount++
		}
	}

	// 显示统计信息
	fmt.Printf("\n📊 处理统计预览:\n")
	fmt.Printf("   💾 需要翻译: %d 个内容\n", len(cachePreview.MissingTranslations))
	fmt.Printf("   🏷️  需要新建标签页面: %d 个\n", createTagCount)
	fmt.Printf("   📝 需要新增文章Slug: %d 个\n", missingSlugCount)

	totalOperations := len(cachePreview.MissingTranslations) + createTagCount + missingSlugCount
	if totalOperations == 0 {
		color.Green("✅ 所有内容都已是最新状态，无需处理")
		return
	}

	fmt.Printf("   📦 预计总操作数: %d 个\n", totalOperations)
	fmt.Println()

	color.Yellow("⚠️  注意：此操作可能需要较长时间，建议在网络稳定时执行")
	fmt.Print("确认开始一键处理？(y/n): ")

	input, _ := reader.ReadString('\n')
	if strings.TrimSpace(strings.ToLower(input)) != "y" {
		color.Yellow("❌ 已取消一键处理")
		return
	}

	fmt.Println()
	color.Cyan("🚀 开始一键处理流程...")
	utils.Info("开始一键处理流程")

	// 步骤1: 生成全量翻译缓存
	if len(cachePreview.MissingTranslations) > 0 {
		color.Blue("\n📦 步骤1/3: 生成全量翻译缓存")
		fmt.Printf("需要翻译 %d 个内容...\n", len(cachePreview.MissingTranslations))

		translatorInstance := translator.NewLLMTranslator()
		_, err = translatorInstance.BatchTranslate(cachePreview.MissingTranslations)
		if err != nil {
			color.Red("❌ 翻译缓存生成失败: %v", err)
			return
		}
		color.Green("✅ 翻译缓存生成完成")
	} else {
		color.Green("\n✅ 步骤1/3: 翻译缓存已是最新")
	}

	// 步骤2: 生成新增标签页面
	if createTagCount > 0 {
		color.Blue("\n🏷️  步骤2/3: 生成新增标签页面")
		fmt.Printf("需要创建 %d 个标签页面...\n", createTagCount)

		err = pageGenerator.GenerateTagPagesWithMode(tagStats, "create")
		if err != nil {
			color.Red("❌ 标签页面生成失败: %v", err)
			return
		}
		color.Green("✅ 标签页面生成完成")
	} else {
		color.Green("\n✅ 步骤2/3: 标签页面已是最新")
	}

	// 步骤3: 生成缺失文章Slug
	if missingSlugCount > 0 {
		color.Blue("\n📝 步骤3/3: 生成缺失文章Slug")
		fmt.Printf("需要添加 %d 个文章Slug...\n", missingSlugCount)

		err = slugGenerator.GenerateArticleSlugsWithMode("missing")
		if err != nil {
			color.Red("❌ 文章Slug生成失败: %v", err)
			return
		}
		color.Green("✅ 文章Slug生成完成")
	} else {
		color.Green("\n✅ 步骤3/3: 文章Slug已是最新")
	}

	// 显示最终统计
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

func showTagAnalysisMenu(tagStats []models.TagStats, reader *bufio.Reader) {
	for {
		color.Cyan("\n=== 标签统计与分析 ===")
		fmt.Println("1. 查看所有标签")
		fmt.Println("2. 查看特定标签详情")
		fmt.Println("3. 按频率分组查看")
		fmt.Println("4. 返回主菜单")
		fmt.Print("请选择 (1-4): ")

		input, _ := reader.ReadString('\n')
		choice := strings.TrimSpace(input)

		switch choice {
		case "1":
			display.DisplayTagStats(tagStats, len(tagStats))
		case "2":
			fmt.Print("请输入标签名: ")
			tagName, _ := reader.ReadString('\n')
			tagName = strings.TrimSpace(tagName)
			if tagName != "" {
				display.DisplayTagDetails(tagStats, tagName)
			} else {
				color.Red("标签名不能为空")
			}
		case "3":
			showTagFrequencyGroups(tagStats)
		case "4":
			return
		default:
			color.Red("⚠️  无效选择，请重新输入")
		}
	}
}

func showTagFrequencyGroups(tagStats []models.TagStats) {
	high, medium, low := stats.GroupTagsByFrequency(tagStats)

	color.Green("=== 高频标签 (≥5篇) ===")
	if len(high) > 0 {
		display.DisplayTagStats(high, len(high))
	} else {
		fmt.Println("没有高频标签")
	}

	color.Yellow("=== 中频标签 (2-4篇) ===")
	if len(medium) > 0 {
		display.DisplayTagStats(medium, len(medium))
	} else {
		fmt.Println("没有中频标签")
	}

	color.Blue("=== 低频标签 (1篇) ===")
	if len(low) > 0 {
		fmt.Printf("共有 %d 个低频标签，显示前20个：\n", len(low))
		limit := 20
		if len(low) < 20 {
			limit = len(low)
		}
		display.DisplayTagStats(low, limit)
	} else {
		fmt.Println("没有低频标签")
	}
}

func previewTagPages(tagStats []models.TagStats, contentDir string) {
	if len(tagStats) == 0 {
		fmt.Println("没有找到任何标签，无法预览")
		return
	}

	pageGenerator := generator.NewTagPageGenerator(contentDir)
	fmt.Printf("即将为 %d 个标签生成页面预览...\n", len(tagStats))

	previews := pageGenerator.PreviewTagPages(tagStats)
	display.DisplayTagPagePreview(previews, 20)
}

func generateTagPages(tagStats []models.TagStats, contentDir string, reader *bufio.Reader) {
	if len(tagStats) == 0 {
		color.Yellow("⚠️  没有找到任何标签，无法生成页面")
		return
	}

	// 先预览以获取统计信息
	color.Cyan("正在分析标签页面状态...")
	pageGenerator := generator.NewTagPageGenerator(contentDir)
	previews := pageGenerator.PreviewTagPages(tagStats)

	createCount := 0
	updateCount := 0
	for _, preview := range previews {
		if preview.Status == "create" {
			createCount++
		} else if preview.Status == "update" {
			updateCount++
		}
	}

	fmt.Printf("\n📊 统计信息:\n")
	fmt.Printf("   🆕 需要新建: %d 个标签页面\n", createCount)
	fmt.Printf("   🔄 需要更新: %d 个标签页面\n", updateCount)
	fmt.Printf("   📦 总计: %d 个标签页面\n", len(previews))

	if createCount == 0 && updateCount == 0 {
		color.Green("✅ 所有标签页面都是最新的")
		return
	}

	// 选择处理模式
	fmt.Println("\n🔧 请选择处理模式:")
	options := []string{}
	if createCount > 0 {
		options = append(options, fmt.Sprintf("1. 仅新增 (%d 个)", createCount))
	}
	if updateCount > 0 {
		options = append(options, fmt.Sprintf("2. 仅更新 (%d 个)", updateCount))
	}
	if createCount > 0 && updateCount > 0 {
		options = append(options, fmt.Sprintf("3. 全部处理 (%d 个)", createCount+updateCount))
	}

	for _, option := range options {
		fmt.Printf("   %s\n", option)
	}
	fmt.Println("   0. 取消操作")
	fmt.Print("请选择: ")

	input, _ := reader.ReadString('\n')
	choice := strings.TrimSpace(input)

	var mode string
	switch choice {
	case "1":
		if createCount == 0 {
			color.Yellow("⚠️  没有需要新增的标签页面")
			return
		}
		mode = "create"
		color.Blue("🆕 将新增 %d 个标签页面", createCount)
	case "2":
		if updateCount == 0 {
			color.Yellow("⚠️  没有需要更新的标签页面")
			return
		}
		mode = "update"
		color.Blue("🔄 将更新 %d 个标签页面", updateCount)
	case "3":
		if createCount == 0 && updateCount == 0 {
			color.Yellow("⚠️  没有需要处理的标签页面")
			return
		}
		mode = "all"
		color.Blue("📦 将处理 %d 个标签页面", createCount+updateCount)
	case "0":
		color.Yellow("❌ 已取消操作")
		return
	default:
		color.Red("⚠️  无效选择")
		return
	}

	fmt.Print("\n确认执行？(y/n): ")
	input, _ = reader.ReadString('\n')
	if strings.TrimSpace(strings.ToLower(input)) != "y" {
		color.Yellow("❌ 已取消生成")
		return
	}

	color.Cyan("🚀 正在生成标签页面...")
	if err := pageGenerator.GenerateTagPagesWithMode(tagStats, mode); err != nil {
		color.Red("❌ 生成失败: %v", err)
	}
}

func previewArticleSlugs(contentDir string) {
	fmt.Println("正在扫描文章并生成Slug预览...")

	slugGenerator := generator.NewArticleSlugGenerator(contentDir)
	previews, err := slugGenerator.PreviewArticleSlugs()
	if err != nil {
		fmt.Printf("预览失败: %v\n", err)
		return
	}

	if len(previews) == 0 {
		fmt.Println("没有找到需要处理的文章")
		return
	}

	display.DisplayArticleSlugPreview(previews, 20)
}

func generateArticleSlugs(contentDir string, reader *bufio.Reader) {
	color.Cyan("🔍 正在扫描文章...")

	slugGenerator := generator.NewArticleSlugGenerator(contentDir)
	previews, err := slugGenerator.PreviewArticleSlugs()
	if err != nil {
		color.Red("❌ 扫描失败: %v", err)
		return
	}

	if len(previews) == 0 {
		color.Green("✅ 没有找到需要处理的文章")
		return
	}

	// 统计信息
	missingCount := 0
	updateCount := 0
	for _, preview := range previews {
		if preview.Status == "missing" {
			missingCount++
		} else if preview.Status == "update" {
			updateCount++
		}
	}

	fmt.Printf("\n📊 统计信息:\n")
	fmt.Printf("   🆕 缺少slug: %d 篇文章\n", missingCount)
	fmt.Printf("   🔄 需要更新: %d 篇文章\n", updateCount)
	fmt.Printf("   📦 总计: %d 篇文章\n", len(previews))

	if missingCount == 0 && updateCount == 0 {
		color.Green("✅ 所有文章的slug都是最新的")
		return
	}

	// 选择处理模式
	fmt.Println("\n🔧 请选择处理模式:")
	options := []string{}
	if missingCount > 0 {
		options = append(options, fmt.Sprintf("1. 仅新增 (%d 篇)", missingCount))
	}
	if updateCount > 0 {
		options = append(options, fmt.Sprintf("2. 仅更新 (%d 篇)", updateCount))
	}
	if missingCount > 0 && updateCount > 0 {
		options = append(options, fmt.Sprintf("3. 全部处理 (%d 篇)", missingCount+updateCount))
	}

	for _, option := range options {
		fmt.Printf("   %s\n", option)
	}
	fmt.Println("   0. 取消操作")
	fmt.Print("请选择: ")

	input, _ := reader.ReadString('\n')
	choice := strings.TrimSpace(input)

	var mode string
	switch choice {
	case "1":
		if missingCount == 0 {
			color.Yellow("⚠️  没有缺少slug的文章")
			return
		}
		mode = "missing"
		color.Blue("🆕 将为 %d 篇文章新增slug", missingCount)
	case "2":
		if updateCount == 0 {
			color.Yellow("⚠️  没有需要更新slug的文章")
			return
		}
		mode = "update"
		color.Blue("🔄 将为 %d 篇文章更新slug", updateCount)
	case "3":
		if missingCount == 0 && updateCount == 0 {
			color.Yellow("⚠️  没有需要处理的文章")
			return
		}
		mode = "all"
		color.Blue("📦 将为 %d 篇文章处理slug", missingCount+updateCount)
	case "0":
		color.Yellow("❌ 已取消操作")
		return
	default:
		color.Red("⚠️  无效选择")
		return
	}

	fmt.Print("\n确认执行？(y/n): ")
	input, _ = reader.ReadString('\n')
	if strings.TrimSpace(strings.ToLower(input)) != "y" {
		color.Yellow("❌ 已取消生成")
		return
	}

	color.Cyan("🚀 正在生成文章slug...")
	if err := slugGenerator.GenerateArticleSlugsWithMode(mode); err != nil {
		color.Red("❌ 生成失败: %v", err)
	}
}

func showCacheStatus() {
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

func clearTranslationCache(reader *bufio.Reader) {
	color.Yellow("⚠️  警告：此操作将清空所有翻译缓存")
	fmt.Print("确认清空缓存？(y/n): ")

	input, _ := reader.ReadString('\n')
	if strings.TrimSpace(strings.ToLower(input)) != "y" {
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

func previewBulkTranslationCache(tagStats []models.TagStats, contentDir string) {
	color.Cyan("=== 全量翻译缓存预览 ===")

	// 收集所有需要翻译的文本
	cachePreview, err := collectTranslationTargets(tagStats, contentDir)
	if err != nil {
		color.Red("❌ 收集翻译目标失败: %v", err)
		return
	}

	display.DisplayBulkTranslationPreview(cachePreview, 20)
}

func generateBulkTranslationCache(tagStats []models.TagStats, contentDir string, reader *bufio.Reader) {
	color.Cyan("🔍 正在收集翻译目标...")

	// 收集所有需要翻译的文本
	cachePreview, err := collectTranslationTargets(tagStats, contentDir)
	if err != nil {
		color.Red("❌ 收集翻译目标失败: %v", err)
		return
	}

	if len(cachePreview.MissingTranslations) == 0 {
		color.Green("✅ 所有内容都已有翻译缓存")
		return
	}

	fmt.Printf("\n📊 翻译缓存统计:\n")
	fmt.Printf("   🏷️  标签总数: %d 个\n", cachePreview.TotalTags)
	fmt.Printf("   📝 文章总数: %d 篇\n", cachePreview.TotalArticles)
	fmt.Printf("   ✅ 已缓存: %d 个\n", cachePreview.CachedCount)
	fmt.Printf("   🔄 需翻译: %d 个\n", len(cachePreview.MissingTranslations))

	if len(cachePreview.MissingTranslations) == 0 {
		color.Green("✅ 所有翻译都已缓存")
		return
	}

	fmt.Print("\n确认生成全量翻译缓存？(y/n): ")
	input, _ := reader.ReadString('\n')
	if strings.TrimSpace(strings.ToLower(input)) != "y" {
		color.Yellow("❌ 已取消生成")
		return
	}

	color.Cyan("🚀 正在生成全量翻译缓存...")
	translatorInstance := translator.NewLLMTranslator()

	// 批量翻译所有缺失的内容
	_, err = translatorInstance.BatchTranslate(cachePreview.MissingTranslations)
	if err != nil {
		color.Red("❌ 批量翻译失败: %v", err)
		return
	}

	color.Green("✅ 全量翻译缓存生成完成！")
}

func collectTranslationTargets(tagStats []models.TagStats, contentDir string) (*display.BulkTranslationPreview, error) {
	translatorInstance := translator.NewLLMTranslator()

	// 收集所有标签
	tagNames := make([]string, len(tagStats))
	for i, stat := range tagStats {
		tagNames[i] = stat.Name
	}

	// 收集所有文章标题
	articles, err := scanner.ScanArticles(contentDir)
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

func showPerformanceStats() {
	color.Cyan("=== 系统性能统计 ===")
	perfStats := utils.GetGlobalStats()
	fmt.Println()
	fmt.Println(perfStats.String())
	fmt.Println()
}

func resetPerformanceStats(reader *bufio.Reader) {
	color.Yellow("⚠️  警告：此操作将重置所有性能统计数据")
	fmt.Print("确认重置？(y/n): ")

	input, _ := reader.ReadString('\n')
	if strings.TrimSpace(strings.ToLower(input)) != "y" {
		color.Yellow("❌ 已取消重置")
		return
	}

	utils.ResetGlobalStats()
	color.Green("✅ 性能统计已重置")
}
