package operations

import (
	"bufio"
	"fmt"
	"hugo-content-suite/generator"
	"hugo-content-suite/models"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fatih/color"
)

type Processor struct {
	contentDir string
}

func NewProcessor(contentDir string) *Processor {
	return &Processor{
		contentDir: contentDir,
	}
}

func (p *Processor) QuickProcessAll(tagStats []models.TagStats, reader *bufio.Reader) {
	color.Cyan("=== 一键处理全部 ===")
	fmt.Println("正在自动执行完整的处理流程：")
	fmt.Println("1. 生成全量翻译缓存")
	fmt.Println("2. 生成新增标签页面")
	fmt.Println("3. 生成缺失文章Slug")
	fmt.Println("4. 翻译新增文章为英文")
	fmt.Println()

	// 预览批量翻译缓存
	cachePreview := p.PreviewBulkTranslationCache(tagStats)

	// 预览标签页面
	tagGenerator := generator.NewTagPageGenerator(p.contentDir)
	tagPreviews := tagGenerator.PreviewTagPages(tagStats)
	createTagCount := 0
	for _, preview := range tagPreviews {
		if preview.Status == "create" {
			createTagCount++
		}
	}

	// 预览文章Slug
	slugGenerator := generator.NewArticleSlugGenerator(p.contentDir)
	slugPreviews, err := slugGenerator.PreviewArticleSlugs()
	missingSlugCount := 0
	if err == nil {
		for _, preview := range slugPreviews {
			if preview.Status == "missing" {
				missingSlugCount++
			}
		}
	}

	// 预览文章翻译
	articleTranslator := generator.NewArticleTranslator(p.contentDir)
	translationStatus, err := articleTranslator.GetTranslationStatus()
	missingTranslationCount := 0
	if err == nil {
		missingTranslationCount = translationStatus.MissingArticles
	}

	// 显示总体预览
	fmt.Printf("📊 检测到需要处理的内容:\n")
	fmt.Printf("   🔄 需要翻译: %d 个项目\n", len(cachePreview.MissingTranslations))
	fmt.Printf("   🏷️  需要创建标签页面: %d 个\n", createTagCount)
	fmt.Printf("   📝 需要添加文章Slug: %d 个\n", missingSlugCount)
	fmt.Printf("   🌐 需要翻译文章: %d 篇\n", missingTranslationCount)

	totalTasks := 0
	if len(cachePreview.MissingTranslations) > 0 {
		totalTasks++
	}
	if createTagCount > 0 {
		totalTasks++
	}
	if missingSlugCount > 0 {
		totalTasks++
	}
	if missingTranslationCount > 0 {
		totalTasks++
	}

	if totalTasks == 0 {
		color.Green("✅ 所有内容都已是最新状态，无需处理")
		return
	}

	// 直接执行处理流程，无需确认
	currentStep := 1
	color.Cyan("🚀 开始自动执行处理流程...")

	// 步骤1: 生成翻译缓存
	if len(cachePreview.MissingTranslations) > 0 {
		fmt.Printf("\n步骤 %d/%d: 生成翻译缓存\n", currentStep, totalTasks)
		p.GenerateBulkTranslationCache(tagStats, reader)
		currentStep++
	}

	// 步骤2: 生成标签页面
	if createTagCount > 0 {
		fmt.Printf("\n步骤 %d/%d: 生成标签页面\n", currentStep, totalTasks)
		err := tagGenerator.GenerateTagPagesWithMode(tagStats, "create")
		if err != nil {
			color.Red("❌ 生成标签页面失败: %v", err)
		} else {
			color.Green("✅ 标签页面生成完成")
		}
		currentStep++
	}

	// 步骤3: 生成文章Slug
	if missingSlugCount > 0 {
		fmt.Printf("\n步骤 %d/%d: 生成文章Slug\n", currentStep, totalTasks)
		err := slugGenerator.GenerateArticleSlugsWithMode("missing")
		if err != nil {
			color.Red("❌ 生成文章Slug失败: %v", err)
		} else {
			color.Green("✅ 文章Slug生成完成")
		}
		currentStep++
	}

	// 步骤4: 翻译文章
	if missingTranslationCount > 0 {
		fmt.Printf("\n步骤 %d/%d: 翻译文章\n", currentStep, totalTasks)
		err := articleTranslator.TranslateArticles("missing")
		if err != nil {
			color.Red("❌ 文章翻译失败: %v", err)
		} else {
			color.Green("✅ 文章翻译完成")
		}
	}

	color.Green("\n🎉 一键处理全部完成！")
}

func (p *Processor) confirmExecution(reader *bufio.Reader, prompt string) bool {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(strings.ToLower(input)) == "y"
}

func (p *Processor) ScanLanguages() ([]string, error) {
	// 扫描 contentDir 下所有文章，收集所有语言（假设文件名格式为 xxx.{lang}.md）
	langSet := make(map[string]struct{})
	err := filepath.Walk(p.contentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		base := filepath.Base(path)
		parts := strings.Split(base, ".")
		// 只处理以 .md 结尾的文件
		if len(parts) >= 3 && parts[len(parts)-1] == "md" {
			lang := parts[len(parts)-2]
			langSet[lang] = struct{}{}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	langs := make([]string, 0, len(langSet))
	for lang := range langSet {
		langs = append(langs, lang)
	}
	sort.Strings(langs)
	return langs, nil
}

func (p *Processor) DeleteArticlesByLanguage(lang string) error {
	// 删除 contentDir 下所有指定语言的文章
	return filepath.Walk(p.contentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		base := filepath.Base(path)
		parts := strings.Split(base, ".")
		if len(parts) >= 3 && parts[len(parts)-2] == lang {
			return os.Remove(path)
		}
		return nil
	})
}
