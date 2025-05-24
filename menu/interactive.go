package menu

import (
	"bufio"
	"fmt"
	"strings"
	"tag-scanner/display"
	"tag-scanner/models"
	"tag-scanner/operations"
	"tag-scanner/stats"
	"tag-scanner/utils"

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

func (m *InteractiveMenu) Show(tagStats []models.TagStats, categoryStats []models.CategoryStats, noTagArticles []models.Article) {
	for {
		m.displayMainMenu()
		choice := m.getChoice("请选择功能 (0-14): ")

		switch choice {
		case "1":
			m.processor.QuickProcessAll(tagStats, m.reader)
		case "2":
			m.showTagAnalysisMenu(tagStats)
		case "3":
			display.DisplayCategoryStats(categoryStats)
		case "4":
			display.DisplayNoTagArticles(noTagArticles, len(noTagArticles))
		case "5":
			m.processor.PreviewTagPages(tagStats)
		case "6":
			m.processor.GenerateTagPages(tagStats, m.reader)
		case "7":
			m.processor.PreviewArticleSlugs()
		case "8":
			m.processor.GenerateArticleSlugs(m.reader)
		case "9":
			m.processor.ShowCacheStatus()
		case "10":
			m.processor.PreviewBulkTranslationCache(tagStats)
		case "11":
			m.processor.GenerateBulkTranslationCache(tagStats, m.reader)
		case "12":
			m.processor.ClearTranslationCache(m.reader)
		case "13":
			m.showPerformanceStats()
		case "14":
			m.resetPerformanceStats()
		case "0":
			color.Green("感谢使用！再见！")
			return
		default:
			color.Red("⚠️  无效选择，请重新输入")
		}
	}
}

func (m *InteractiveMenu) displayMainMenu() {
	color.Cyan("\n=== 博客标签管理工具 ===")
	fmt.Println()

	// 快速处理模块
	color.Red("🚀 快速处理")
	fmt.Println("  1. 一键处理全部 (缓存→标签页面→文章Slug)")
	fmt.Println()

	// 数据查看模块
	color.Green("📊 数据查看")
	fmt.Println("  2. 标签统计与分析")
	fmt.Println("  3. 分类统计")
	fmt.Println("  4. 无标签文章")
	fmt.Println()

	// 页面生成模块
	color.Yellow("🏷️  标签页面管理")
	fmt.Println("  5. 预览标签页面")
	fmt.Println("  6. 生成标签页面")
	fmt.Println()

	// 文章管理模块
	color.Blue("📝 文章Slug管理")
	fmt.Println("  7. 预览文章Slug")
	fmt.Println("  8. 生成文章Slug")
	fmt.Println()

	// 缓存管理模块
	color.Magenta("💾 缓存管理")
	fmt.Println("  9. 查看缓存状态")
	fmt.Println(" 10. 预览全量翻译缓存")
	fmt.Println(" 11. 生成全量翻译缓存")
	fmt.Println(" 12. 清空翻译缓存")
	fmt.Println()

	// 系统工具模块
	color.Cyan("🔧 系统工具")
	fmt.Println(" 13. 查看性能统计")
	fmt.Println(" 14. 重置性能统计")
	fmt.Println()

	color.Red("  0. 退出程序")
	fmt.Println()
}

func (m *InteractiveMenu) showTagAnalysisMenu(tagStats []models.TagStats) {
	for {
		color.Cyan("\n=== 标签统计与分析 ===")
		fmt.Println("1. 查看所有标签")
		fmt.Println("2. 查看特定标签详情")
		fmt.Println("3. 按频率分组查看")
		fmt.Println("4. 返回主菜单")

		choice := m.getChoice("请选择 (1-4): ")

		switch choice {
		case "1":
			display.DisplayTagStats(tagStats, len(tagStats))
		case "2":
			tagName := m.getChoice("请输入标签名: ")
			if tagName != "" {
				display.DisplayTagDetails(tagStats, tagName)
			} else {
				color.Red("标签名不能为空")
			}
		case "3":
			m.showTagFrequencyGroups(tagStats)
		case "4":
			return
		default:
			color.Red("⚠️  无效选择，请重新输入")
		}
	}
}

func (m *InteractiveMenu) showTagFrequencyGroups(tagStats []models.TagStats) {
	high, medium, low := stats.GroupTagsByFrequency(tagStats)

	color.Green("=== 高频标签 (≥5篇) ===")
	if len(high) > 0 {
		display.DisplayTagStats(high, len(high))
	} else {
		fmt.Println("没有高频标签")
	}

	color.Yellow("=== 中频标签 (2-4篇) ===")
	if len(medium) > 0 {
		display.DisplayTagStats(medium, len(medium))
	} else {
		fmt.Println("没有中频标签")
	}

	color.Blue("=== 低频标签 (1篇) ===")
	if len(low) > 0 {
		fmt.Printf("共有 %d 个低频标签，显示前20个：\n", len(low))
		limit := 20
		if len(low) < 20 {
			limit = len(low)
		}
		display.DisplayTagStats(low, limit)
	} else {
		fmt.Println("没有低频标签")
	}
}

func (m *InteractiveMenu) showPerformanceStats() {
	color.Cyan("=== 系统性能统计 ===")
	perfStats := utils.GetGlobalStats()
	fmt.Println()
	fmt.Println(perfStats.String())
	fmt.Println()
}

func (m *InteractiveMenu) resetPerformanceStats() {
	color.Yellow("⚠️  警告：此操作将重置所有性能统计数据")
	confirm := m.getChoice("确认重置？(y/n): ")

	if strings.TrimSpace(strings.ToLower(confirm)) != "y" {
		color.Yellow("❌ 已取消重置")
		return
	}

	utils.ResetGlobalStats()
	color.Green("✅ 性能统计已重置")
}

func (m *InteractiveMenu) getChoice(prompt string) string {
	fmt.Print(prompt)
	input, _ := m.reader.ReadString('\n')
	return strings.TrimSpace(input)
}
