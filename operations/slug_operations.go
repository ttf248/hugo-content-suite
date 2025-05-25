package operations

import (
	"bufio"
	"fmt"
	"hugo-content-suite/generator"

	"github.com/fatih/color"
)

func (p *Processor) GenerateArticleSlugs(reader *bufio.Reader) {
	color.Cyan("🔍 正在扫描文章...")

	slugGenerator := generator.NewArticleSlugGenerator(p.contentDir)
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
	missingCount, updateCount := p.countSlugOperations(previews)
	p.displaySlugStats(missingCount, updateCount, len(previews))

	if missingCount == 0 && updateCount == 0 {
		color.Green("✅ 所有文章的slug都是最新的")
		return
	}

	// 选择处理模式
	mode := p.selectSlugMode(missingCount, updateCount, reader)
	if mode == "" {
		return
	}

	if !p.confirmExecution(reader, "\n确认执行？(y/n): ") {
		color.Yellow("❌ 已取消生成")
		return
	}

	color.Cyan("🚀 正在生成文章slug...")
	if err := slugGenerator.GenerateArticleSlugsWithMode(mode); err != nil {
		color.Red("❌ 生成失败: %v", err)
	}
}

func (p *Processor) countSlugOperations(previews []generator.ArticleSlugPreview) (int, int) {
	missingCount := 0
	updateCount := 0
	for _, preview := range previews {
		if preview.Status == "missing" {
			missingCount++
		} else if preview.Status == "update" {
			updateCount++
		}
	}
	return missingCount, updateCount
}

func (p *Processor) displaySlugStats(missingCount, updateCount, total int) {
	fmt.Printf("\n📊 统计信息:\n")
	fmt.Printf("   🆕 缺少slug: %d 篇文章\n", missingCount)
	fmt.Printf("   🔄 需要更新: %d 篇文章\n", updateCount)
	fmt.Printf("   📦 总计: %d 篇文章\n", total)
}

func (p *Processor) selectSlugMode(missingCount, updateCount int, reader *bufio.Reader) string {
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

	choice := p.getChoice(reader, "请选择: ")

	switch choice {
	case "1":
		if missingCount == 0 {
			color.Yellow("⚠️  没有缺少slug的文章")
			return ""
		}
		color.Blue("🆕 将为 %d 篇文章新增slug", missingCount)
		return "missing"
	case "2":
		if updateCount == 0 {
			color.Yellow("⚠️  没有需要更新slug的文章")
			return ""
		}
		color.Blue("🔄 将为 %d 篇文章更新slug", updateCount)
		return "update"
	case "3":
		if missingCount == 0 && updateCount == 0 {
			color.Yellow("⚠️  没有需要处理的文章")
			return ""
		}
		color.Blue("📦 将为 %d 篇文章处理slug", missingCount+updateCount)
		return "all"
	case "0":
		color.Yellow("❌ 已取消操作")
		return ""
	default:
		color.Red("⚠️  无效选择")
		return ""
	}
}
