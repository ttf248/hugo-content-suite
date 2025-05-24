package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"tag-scanner/config"
	"tag-scanner/display"
	"tag-scanner/menu"
	"tag-scanner/scanner"
	"tag-scanner/stats"
	"tag-scanner/utils"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("配置加载失败:", err)
	}

	// 初始化日志
	if err := utils.InitLogger("tag-scanner.log", utils.DEBUG); err != nil {
		log.Printf("日志初始化失败: %v", err)
	}

	utils.Info("程序启动")
	defer utils.Info("程序退出")

	contentDir := cfg.Paths.DefaultContentDir
	if len(os.Args) > 1 {
		contentDir = os.Args[1]
	}

	// 扫描文章
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
	display.DisplayTagStats(tagStats, 20)
	display.DisplayCategoryStats(categoryStats)
	display.DisplayNoTagArticles(noTagArticles, 10)

	// 启动交互式菜单
	reader := bufio.NewReader(os.Stdin)
	interactiveMenu := menu.NewInteractiveMenu(reader, contentDir)
	interactiveMenu.Show(tagStats, categoryStats, noTagArticles)

	// 显示性能统计
	perfStats := utils.GetGlobalStats()
	if perfStats.TranslationCount > 0 || perfStats.FileOperations > 0 {
		fmt.Println()
		fmt.Println(perfStats.String())
	}
}
