package operations

import (
	"bufio"
	"fmt"
	"hugo-content-suite/generator"

	"github.com/fatih/color"
)

func (p *Processor) GenerateArticleSlugs(reader *bufio.Reader) {
	if p.contentDir == "" {
		color.Red("❌ 内容目录未设置")
		return
	}

	// 获取文章slug状态统计
	color.Cyan("正在分析文章slug状态...")
	slugGenerator := generator.NewArticleSlugGenerator(p.contentDir)
	previews, createCount, updateCount, err := slugGenerator.PrepareArticleSlugs()
	if err != nil {
		color.Red("❌ 分析失败: %v", err)
		return
	}

	p.displaySlugStats(createCount, updateCount, len(previews))

	if createCount == 0 && updateCount == 0 {
		color.Green("✅ 所有文章slug都是最新的")
		return
	}

	// 选择处理模式
	mode := p.selectPageMode(ArticleSlug, createCount, updateCount, reader)
	if mode == "" {
		return
	}

	// 根据模式筛选预览
	targetPreviews := filterByMode(previews, mode)

	// 显示警告和确认
	p.displaySlugWarning(mode, createCount, updateCount)

	if !p.confirmExecution(reader, "\n确认开始生成？(y/n): ") {
		color.Yellow("❌ 已取消生成")
		return
	}

	color.Cyan("🚀 开始生成文章slug...")
	if err := slugGenerator.GenerateArticleSlugsWithMode(targetPreviews, mode); err != nil {
		color.Red("❌ 生成失败: %v", err)
	}
}

func (p *Processor) displaySlugStats(createCount, updateCount, total int) {
	fmt.Printf("\n📊 Slug统计信息:\n")
	fmt.Printf("   🆕 需要新建slug的文章: %d 篇\n", createCount)
	fmt.Printf("   🔄 需要更新slug的文章: %d 篇\n", updateCount)
	fmt.Printf("   ✅ slug已是最新的文章: %d 篇\n", total-createCount-updateCount)
	fmt.Printf("   📦 文章总数: %d 篇\n", total)

	if createCount > 0 || updateCount > 0 {
		fmt.Printf("\n💡 说明:\n")
		fmt.Printf("   • 需要新建: 文章front matter中缺少slug字段\n")
		fmt.Printf("   • 需要更新: 现有slug与AI生成的slug不匹配\n")
	}
}

func (p *Processor) displaySlugWarning(mode string, createCount, updateCount int) {
	fmt.Println()
	color.Yellow("⚠️  重要提示:")
	fmt.Println("• Slug生成基于AI翻译，可能需要较长时间")
	fmt.Println("• 生成过程中请保持网络连接稳定")
	fmt.Println("• Slug生成会使用缓存加速重复内容的处理")

	switch mode {
	case "create":
		fmt.Printf("• 将为 %d 篇文章新建slug\n", createCount)
	case "update":
		fmt.Printf("• 将更新 %d 篇文章的slug\n", updateCount)
	case "all":
		fmt.Printf("• 将处理 %d 篇文章的slug（包括新建和更新）\n", createCount+updateCount)
	}
}
