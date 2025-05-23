package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"tag-scanner/display"
	"tag-scanner/generator"
	"tag-scanner/models"
	"tag-scanner/scanner"
	"tag-scanner/stats"

	"github.com/fatih/color"
)

func main() {
	contentDir := "../../content/post"
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
}

func showInteractiveMenu(tagStats []models.TagStats, categoryStats []models.CategoryStats, noTagArticles []models.Article, contentDir string) {
	reader := bufio.NewReader(os.Stdin)

	for {
		color.Cyan("\n=== 博客标签管理工具 ===")
		fmt.Println()

		// 数据查看模块
		color.Green("📊 数据查看")
		fmt.Println("  1. 标签统计与分析")
		fmt.Println("  2. 分类统计")
		fmt.Println("  3. 无标签文章")
		fmt.Println()

		// 页面生成模块
		color.Yellow("🏷️  标签页面管理")
		fmt.Println("  4. 预览标签页面")
		fmt.Println("  5. 生成标签页面")
		fmt.Println()

		// 文章管理模块
		color.Blue("📝 文章Slug管理")
		fmt.Println("  6. 预览文章Slug")
		fmt.Println("  7. 生成文章Slug")
		fmt.Println()

		color.Red("  0. 退出程序")
		fmt.Println()
		fmt.Print("请选择功能 (0-7): ")

		input, _ := reader.ReadString('\n')
		choice := strings.TrimSpace(input)

		switch choice {
		case "1":
			showTagAnalysisMenu(tagStats, reader)
		case "2":
			display.DisplayCategoryStats(categoryStats)
		case "3":
			display.DisplayNoTagArticles(noTagArticles, len(noTagArticles))
		case "4":
			previewTagPages(tagStats, contentDir)
		case "5":
			generateTagPages(tagStats, contentDir, reader)
		case "6":
			previewArticleSlugs(contentDir)
		case "7":
			generateArticleSlugs(contentDir, reader)
		case "0":
			color.Green("感谢使用！再见！")
			return
		default:
			color.Red("⚠️  无效选择，请重新输入")
		}
	}
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
