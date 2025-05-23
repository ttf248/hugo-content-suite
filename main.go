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
		color.Cyan("\n=== 交互式菜单 ===")
		fmt.Println("1. 查看所有标签")
		fmt.Println("2. 查看特定标签详情")
		fmt.Println("3. 查看所有无标签文章")
		fmt.Println("4. 查看标签频率分组")
		fmt.Println("5. 预览标签页面生成")
		fmt.Println("6. 生成标签页面文件")
		fmt.Println("7. 预览文章Slug生成")
		fmt.Println("8. 生成文章Slug")
		fmt.Println("0. 退出")
		fmt.Print("请选择操作 (0-8): ")

		input, _ := reader.ReadString('\n')
		choice := strings.TrimSpace(input)

		switch choice {
		case "1":
			display.DisplayTagStats(tagStats, len(tagStats))
		case "2":
			fmt.Print("请输入标签名: ")
			tagName, _ := reader.ReadString('\n')
			tagName = strings.TrimSpace(tagName)
			display.DisplayTagDetails(tagStats, tagName)
		case "3":
			display.DisplayNoTagArticles(noTagArticles, len(noTagArticles))
		case "4":
			showTagFrequencyGroups(tagStats)
		case "5":
			previewTagPages(tagStats, contentDir)
		case "6":
			generateTagPages(tagStats, contentDir, reader)
		case "7":
			previewArticleSlugs(contentDir)
		case "8":
			generateArticleSlugs(contentDir, reader)
		case "0":
			fmt.Println("再见！")
			return
		default:
			fmt.Println("无效选择，请重新输入")
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
		fmt.Println("没有找到任何标签，无法生成页面")
		return
	}

	fmt.Printf("将为 %d 个标签生成页面文件\n", len(tagStats))
	fmt.Print("确认生成？(y/n): ")

	input, _ := reader.ReadString('\n')
	if strings.TrimSpace(strings.ToLower(input)) != "y" {
		fmt.Println("已取消生成")
		return
	}

	fmt.Println("正在生成标签页面...")
	pageGenerator := generator.NewTagPageGenerator(contentDir)
	if err := pageGenerator.GenerateTagPages(tagStats); err != nil {
		fmt.Printf("生成失败: %v\n", err)
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
	fmt.Println("正在扫描文章...")

	slugGenerator := generator.NewArticleSlugGenerator(contentDir)
	previews, err := slugGenerator.PreviewArticleSlugs()
	if err != nil {
		fmt.Printf("扫描失败: %v\n", err)
		return
	}

	if len(previews) == 0 {
		fmt.Println("没有找到需要处理的文章")
		return
	}

	fmt.Printf("将为 %d 篇文章生成/更新 slug\n", len(previews))
	fmt.Print("确认生成？(y/n): ")

	input, _ := reader.ReadString('\n')
	if strings.TrimSpace(strings.ToLower(input)) != "y" {
		fmt.Println("已取消生成")
		return
	}

	fmt.Println("正在生成文章slug...")
	if err := slugGenerator.GenerateArticleSlugs(); err != nil {
		fmt.Printf("生成失败: %v\n", err)
	}
}
