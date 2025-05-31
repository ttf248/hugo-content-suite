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

	utils.InfoWithFields("程序启动", map[string]interface{}{
		"log_level":    cfg.Logging.Level,
		"version":      "3.0.0", // 更新版本号以反映重构
		"config":       "loaded",
		"architecture": "refactored",
	})

	defer func() {
		utils.InfoWithFields("程序退出", map[string]interface{}{
			"exit_reason": "normal",
		})
		utils.Close()
	}()

	contentDir := cfg.Paths.DefaultContentDir
	if len(os.Args) > 1 {
		contentDir = os.Args[1]
		utils.InfoWithFields("使用命令行参数指定目录", map[string]interface{}{
			"content_dir": contentDir,
		})
	}

	// 扫描文章
	absContentDir, err := utils.GetAbsolutePath(contentDir)
	if err != nil {
		log.Fatal("无法转换为绝对路径:", err)
	}
	fmt.Printf("📂 扫描目录: %s\n", absContentDir)

	articles, err := scanner.ScanArticles(absContentDir)
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
}
