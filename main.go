package main

import (
	"bufio"
	"fmt"
	"hugo-content-suite/config"
	"hugo-content-suite/operations"
	"hugo-content-suite/utils"
	"log"
	"os"

	"github.com/fatih/color"
)

type InteractiveMenu struct {
	reader    *bufio.Reader
	processor *operations.Processor
}

func NewInteractiveMenu(reader *bufio.Reader, contentDir string) *InteractiveMenu {
	return &InteractiveMenu{
		reader:    reader,
		processor: operations.NewProcessor(contentDir),
	}
}

func (m *InteractiveMenu) Show() {
	for {
		m.displayMainMenu()
		choice := utils.GetChoice(m.reader, "请选择功能 (0-8): ")

		switch choice {
		case "`":
			m.processor.ProcessAllContent(m.reader)
		case "1":
			m.processor.GenerateTagPages(m.reader)
		case "2":
			m.processor.GenerateArticleSlugs(m.reader)
		case "3":
			m.processor.TranslateArticles(m.reader)
		case "4":
			m.processor.DeleteArticles(m.reader)

		case "0":
			color.Green("感谢使用！再见！")
			return
		default:
			color.Red("⚠️  无效选择，请重新输入")
		}
	}
}

func (m *InteractiveMenu) displayMainMenu() {
	color.Cyan("\n=== Hugo 博客管理工具 ===")
	fmt.Println()

	color.Yellow("⚡ 一键操作")
	fmt.Println("  `. 一键处理所有内容（仅新增）")
	fmt.Println()

	// 内容管理模块
	color.Green("📝 内容管理")
	fmt.Println("  1. 生成标签页面")
	fmt.Println("  2. 生成文章Slug")
	fmt.Println("  3. 翻译文章为多语言版本")
	fmt.Println("  4. 删除指定语言的文章")
	fmt.Println()

	fmt.Println()

	color.Red("  0. 退出程序")
	fmt.Println()
}

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("config.json")
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
	fmt.Printf("📂 内容目录: %s\n", contentDir)

	// 启动交互式菜单
	reader := bufio.NewReader(os.Stdin)
	interactiveMenu := NewInteractiveMenu(reader, contentDir)
	interactiveMenu.Show()
}
