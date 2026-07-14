package main

import (
	"bufio"
	"fmt"
	"hugo-content-suite/config"
	"hugo-content-suite/operations"
	"hugo-content-suite/translator"
	"hugo-content-suite/utils"
	"log"
	"os"

	"github.com/fatih/color"
)

type InteractiveMenu struct {
	reader    *bufio.Reader
	processor *operations.Processor
	cfg       *config.Config
}

const processNewFlag = "--process-new"

// parseStartupMode 保持既有的“位置参数为内容目录”行为，同时提供无需交互的增量处理入口。
func parseStartupMode(args []string) (runProcessNew bool, contentDirOverride string, err error) {
	if len(args) == 0 {
		return false, "", nil
	}
	if args[0] != processNewFlag {
		if len(args) > 1 {
			return false, "", fmt.Errorf("内容目录只能指定一次")
		}
		return false, args[0], nil
	}
	if len(args) > 2 {
		return false, "", fmt.Errorf("%s 最多接受一个内容目录参数", processNewFlag)
	}
	if len(args) == 2 {
		return true, args[1], nil
	}
	return true, "", nil
}

func NewInteractiveMenu(reader *bufio.Reader, contentDir string, cfg *config.Config) *InteractiveMenu {
	return &InteractiveMenu{
		reader:    reader,
		processor: operations.NewProcessor(contentDir),
		cfg:       cfg,
	}
}

func (m *InteractiveMenu) Show() {
	for {
		m.displayMainMenu()
		choice := utils.GetChoice(m.reader, "请选择功能 (0-6): ")

		switch choice {
		case ".":
			m.processor.ProcessAllContent(m.reader)
		case "1":
			m.processor.GenerateTagPages(m.reader)
		case "2":
			m.processor.GenerateArticleSlugs(m.reader)
		case "3":
			m.processor.TranslateArticles(m.reader)
		case "4":
			m.processor.DeleteArticles(m.reader)
		case "5":
			m.selectModel()
		case "6":
			if err := translator.NewTranslationUtils().TestConnection(); err != nil {
				color.Red("模型连接失败: %v", err)
			} else {
				color.Green("模型连接成功: %s", m.cfg.ActiveModel)
			}

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
	fmt.Println("  .  一键处理所有内容（仅新增）")
	fmt.Println()

	// 内容管理模块
	color.Green("📝 内容管理")
	fmt.Println("  1. 生成标签页面")
	fmt.Println("  2. 生成文章Slug")
	fmt.Println("  3. 翻译文章为多语言版本")
	fmt.Println("  4. 删除指定语言的文章")
	fmt.Println("  5. 选择翻译模型")
	fmt.Println("  6. 测试当前翻译模型")
	fmt.Println()

	fmt.Println()

	color.Red("  0. 退出程序")
	fmt.Println()
}

func (m *InteractiveMenu) selectModel() {
	fmt.Println("\n可用翻译模型：")
	for i, model := range m.cfg.Models {
		fmt.Printf("  %d. %s%s\n", i+1, model.Name, map[bool]string{true: "（当前）", false: ""}[model.Name == m.cfg.ActiveModel])
	}
	choice := utils.GetChoice(m.reader, "选择模型编号（0 取消）: ")
	var index int
	if _, err := fmt.Sscanf(choice, "%d", &index); err != nil || index < 1 || index > len(m.cfg.Models) {
		return
	}
	if err := m.cfg.SelectModel(m.cfg.Models[index-1].Name); err != nil {
		color.Red("切换模型失败: %v", err)
		return
	}
	color.Green("当前翻译模型: %s", m.cfg.ActiveModel)
}

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("config.local.json")
	if err != nil {
		log.Fatal("配置加载失败:", err)
	}
	config.SetGlobalConfig(cfg)

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

	runProcessNew, contentDirOverride, err := parseStartupMode(os.Args[1:])
	if err != nil {
		log.Fatal("命令行参数错误:", err)
	}

	contentDir := cfg.Paths.DefaultContentDir
	if contentDirOverride != "" {
		contentDir = contentDirOverride
		utils.InfoWithFields("使用命令行参数指定目录", map[string]interface{}{
			"content_dir": contentDir,
		})
	}

	fmt.Printf("📂 内容目录: %s\n", contentDir)
	if runProcessNew {
		operations.NewProcessor(contentDir).ProcessAllContent(bufio.NewReader(os.Stdin))
		return
	}

	// 未传入 CLI 标志时维持原有交互菜单。
	reader := bufio.NewReader(os.Stdin)
	interactiveMenu := NewInteractiveMenu(reader, contentDir, cfg)
	interactiveMenu.Show()
}
