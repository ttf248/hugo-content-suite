package main

import (
	"bufio"
	"fmt"
	"hugo-content-suite/config"
	"hugo-content-suite/display"
	"hugo-content-suite/menu"
	"hugo-content-suite/scanner"
	"hugo-content-suite/stats"
	"hugo-content-suite/utils"
	"log"
	"os"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("配置加载失败:", err)
	}

	// 从配置读取日志等级并初始化日志
	logLevel := utils.INFO // 默认等级
	switch cfg.Logging.Level {
	case "DEBUG":
		logLevel = utils.DEBUG
	case "INFO":
		logLevel = utils.INFO
	case "WARN":
		logLevel = utils.WARN
	case "ERROR":
		logLevel = utils.ERROR
	}

	if err := utils.InitLogger(cfg.Logging.File, logLevel); err != nil {
		log.Printf("日志初始化失败: %v", err)
	}

	utils.Info("程序启动，日志等级: %s", cfg.Logging.Level)
	defer utils.Info("程序退出")
	defer utils.Close() // 确保程序退出时关闭日志文件

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

	// 显示简化概览
	display.DisplaySummary(len(articles), tagStats, categoryStats)

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
