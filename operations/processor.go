package operations

import (
	"bufio"
	"fmt"
	"hugo-content-suite/generator"
	"hugo-content-suite/utils"
	"strings"
	"time"

	"github.com/fatih/color"
)

const (
	TagPageLabel    = "标签页面"
	ArticleCategory = "文章分类"
	ArticleSlug     = "文章Slug"

	ModeUpdate = "update"
	ModeCreate = "create"
	ModeAll    = "all"
)

type StatusLike interface {
	GetStatus() string
}

type Processor struct {
	contentDir string
}

func NewProcessor(contentDir string) *Processor {
	return &Processor{
		contentDir: contentDir,
	}
}

// 新增GenerateArticleSlugs方法声明（在article_operations.go中实现）

// 通用筛选函数
func filterByMode[T StatusLike](items []T, mode string) []T {
	var result []T
	for _, item := range items {
		switch mode {
		case ModeCreate:
			if item.GetStatus() == ModeCreate {
				result = append(result, item)
			}
		case ModeUpdate:
			if item.GetStatus() == ModeUpdate {
				result = append(result, item)
			}
		case ModeAll:
			if item.GetStatus() == ModeCreate || item.GetStatus() == ModeUpdate {
				result = append(result, item)
			}
		}
	}
	return result
}

func (p *Processor) confirmExecution(reader *bufio.Reader, prompt string) bool {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(strings.ToLower(input)) == "y"
}

func (p *Processor) selectPageMode(info string, createCount, updateCount int, reader *bufio.Reader) string {
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

	choice := utils.GetChoice(reader, "请选择: ")

	switch choice {
	case "1":
		if createCount == 0 {
			color.Yellow(fmt.Sprintf("⚠️  没有需要新增的 %s", info))
			return ""
		}
		color.Blue("🆕 将新增 %d 个 %s", createCount, info)
		return "create"
	case "2":
		if updateCount == 0 {
			color.Yellow(fmt.Sprintf("⚠️  没有需要更新的 %s", info))
			return ""
		}
		color.Blue("🔄 将更新 %d 个 %s", updateCount, info)
		return "update"
	case "3":
		if createCount == 0 && updateCount == 0 {
			color.Yellow(fmt.Sprintf("⚠️  没有需要处理的 %s", info))
			return ""
		}
		color.Blue("📦 将处理 %d 个 %s", createCount+updateCount, info)
		return "all"
	case "0":
		color.Yellow("❌ 已取消操作")
		return ""
	default:
		color.Red("⚠️  无效选择")
		return ""
	}
}

// ProcessAllContent 一键处理所有内容（仅新增数据）
func (p *Processor) ProcessAllContent(reader *bufio.Reader) {
	if p.contentDir == "" {
		color.Red("❌ 内容目录未设置")
		return
	}

	color.Cyan("🚀 一键处理所有内容")
	color.Cyan("=================")
	fmt.Println("将依次执行以下操作（仅处理新增内容）：")
	fmt.Println("  1. 生成标签页面")
	fmt.Println("  2. 生成文章Slug")
	fmt.Println("  3. 翻译文章为多语言版本")
	fmt.Println()

	startTime := time.Now()
	var totalErrors int

	// 步骤1：生成标签页面
	color.Cyan("\n📖 步骤 1/3: 生成标签页面")
	color.Cyan("=======================")
	if err := p.processTagPagesAutomatically(); err != nil {
		color.Red("❌ 标签页面生成失败: %v", err)
		totalErrors++
	} else {
		color.Green("✅ 标签页面生成完成")
	}

	// 步骤2：生成文章Slug
	color.Cyan("\n📝 步骤 2/3: 生成文章Slug")
	color.Cyan("========================")
	if err := p.processArticleSlugsAutomatically(); err != nil {
		color.Red("❌ 文章Slug生成失败: %v", err)
		totalErrors++
	} else {
		color.Green("✅ 文章Slug生成完成")
	}

	// 步骤3：翻译文章
	color.Cyan("\n🌐 步骤 3/3: 翻译文章")
	color.Cyan("==================")
	if err := p.processArticleTranslationAutomatically(); err != nil {
		color.Red("❌ 文章翻译失败: %v", err)
		totalErrors++
	} else {
		color.Green("✅ 文章翻译完成")
	}

	// 总结
	duration := time.Since(startTime)
	color.Cyan("\n🎉 一键处理完成!")
	color.Cyan("===============")
	fmt.Printf("⏱️  总用时: %v\n", duration.Round(time.Second))
	if totalErrors > 0 {
		color.Yellow("⚠️  完成时遇到 %d 个错误，请检查日志", totalErrors)
	} else {
		color.Green("✅ 所有操作成功完成")
	}
}

// processTagPagesAutomatically 自动处理标签页面生成
func (p *Processor) processTagPagesAutomatically() error {
	pageGenerator := generator.NewTagPageGenerator(p.contentDir)
	previews, createCount, _ := pageGenerator.PrepareTagPages()

	if createCount == 0 {
		color.Green("✅ 所有标签页面都是最新的")
		return nil
	}

	// 只处理新增的标签页面
	targetPreviews := filterByMode(previews, "create")
	if len(targetPreviews) == 0 {
		color.Green("✅ 没有需要新建的标签页面")
		return nil
	}

	color.Cyan("🚀 自动生成新标签页面...")
	return pageGenerator.GenerateTagPagesWithMode(targetPreviews, "create")
}

// processArticleSlugsAutomatically 自动处理文章Slug生成
func (p *Processor) processArticleSlugsAutomatically() error {
	slugGenerator := generator.NewArticleSlugGenerator(p.contentDir)
	previews, createCount, _, err := slugGenerator.PrepareArticleSlugs()
	if err != nil {
		return fmt.Errorf("分析文章slug失败: %v", err)
	}

	if createCount == 0 {
		color.Green("✅ 所有文章slug都是最新的")
		return nil
	}

	// 只处理缺失的slug
	targetPreviews := filterByMode(previews, "create")
	if len(targetPreviews) == 0 {
		color.Green("✅ 没有需要新建的文章slug")
		return nil
	}

	color.Cyan("🚀 自动生成新文章slug...")
	return slugGenerator.GenerateArticleSlugsWithMode(targetPreviews, "create")
}

// processArticleTranslationAutomatically 自动处理文章翻译
func (p *Processor) processArticleTranslationAutomatically() error {
	articleTranslator := generator.NewArticleTranslator(p.contentDir)
	previews, createCount, _, err := articleTranslator.PrepareArticleTranslations()
	if err != nil {
		return fmt.Errorf("分析文章翻译失败: %v", err)
	}

	if createCount == 0 {
		color.Green("✅ 所有文章都已完全翻译")
		return nil
	}

	// 只处理缺失的翻译
	targetPreviews := filterTranslationsByMode(previews, "create")
	if len(targetPreviews) == 0 {
		color.Green("✅ 没有需要新建的文章翻译")
		return nil
	}

	color.Cyan("🚀 自动翻译缺失的文章...")
	return articleTranslator.TranslateArticlesWithMode(targetPreviews, "create")
}
