package operations

import (
	"bufio"
	"fmt"
	"hugo-content-suite/generator"
	"hugo-content-suite/models"
	"hugo-content-suite/utils"

	"github.com/fatih/color"
)

func (p *Processor) GenerateTagPages(tagStats []models.TagStats, reader *bufio.Reader) {
	if len(tagStats) == 0 {
		color.Yellow("⚠️  没有找到任何标签，无法生成页面")
		return
	}

	// 先预览以获取统计信息
	color.Cyan("正在分析标签页面状态...")
	pageGenerator := generator.NewTagPageGenerator(p.contentDir)
	previews := pageGenerator.PreviewTagPages(tagStats)

	createCount, updateCount := pageGenerator.CountPageOperations(previews)

	if createCount == 0 && updateCount == 0 {
		color.Green("✅ 所有标签页面都是最新的")
		return
	}

	// 选择处理模式
	mode := p.selectPageMode(createCount, updateCount, reader)
	if mode == "" {
		return
	}

	// 根据模式筛选预览
	var targetPreviews []generator.TagPagePreview
	for _, preview := range previews {
		switch mode {
		case "create":
			if preview.Status == "create" {
				targetPreviews = append(targetPreviews, preview)
			}
		case "update":
			if preview.Status == "update" {
				targetPreviews = append(targetPreviews, preview)
			}
		case "all":
			targetPreviews = append(targetPreviews, preview)
		}
	}

	if !p.confirmExecution(reader, "\n确认执行？(y/n): ") {
		color.Yellow("❌ 已取消生成")
		return
	}

	color.Cyan("🚀 正在生成标签页面...")
	if err := pageGenerator.GenerateTagPagesWithMode(targetPreviews, mode); err != nil {
		color.Red("❌ 生成失败: %v", err)
	}
}

func (p *Processor) selectPageMode(createCount, updateCount int, reader *bufio.Reader) string {
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
			color.Yellow("⚠️  没有需要新增的标签页面")
			return ""
		}
		color.Blue("🆕 将新增 %d 个标签页面", createCount)
		return "create"
	case "2":
		if updateCount == 0 {
			color.Yellow("⚠️  没有需要更新的标签页面")
			return ""
		}
		color.Blue("🔄 将更新 %d 个标签页面", updateCount)
		return "update"
	case "3":
		if createCount == 0 && updateCount == 0 {
			color.Yellow("⚠️  没有需要处理的标签页面")
			return ""
		}
		color.Blue("📦 将处理 %d 个标签页面", createCount+updateCount)
		return "all"
	case "0":
		color.Yellow("❌ 已取消操作")
		return ""
	default:
		color.Red("⚠️  无效选择")
		return ""
	}
}
