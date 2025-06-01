package operations

import (
	"bufio"
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
	mode := utils.SelectPageMode(utils.TagPageLabel, createCount, updateCount, reader)
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
