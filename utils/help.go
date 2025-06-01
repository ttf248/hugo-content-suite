package utils

import (
	"bufio"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

const (
	TagPageLabel    = "标签页面"
	ArticleCategory = "文章分类"
	ArticleSlug     = "文章Slug"
)

// GetAbsolutePath converts a relative path to an absolute path.
func GetAbsolutePath(relativePath string) (string, error) {
	absPath, err := filepath.Abs(relativePath)
	if err != nil {
		return "", err
	}
	return absPath, nil
}

func GetChoice(reader *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func SelectPageMode(info string, createCount, updateCount int, reader *bufio.Reader) string {
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

	choice := GetChoice(reader, "请选择: ")

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
