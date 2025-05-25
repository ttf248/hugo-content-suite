package generator

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FileUtils 文件操作工具
type FileUtils struct{}

// NewFileUtils 创建文件工具实例
func NewFileUtils() *FileUtils {
	return &FileUtils{}
}

// EnsureDir 确保目录存在
func (f *FileUtils) EnsureDir(dirPath string) error {
	return os.MkdirAll(dirPath, 0755)
}

// FileExists 检查文件是否存在
func (f *FileUtils) FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

// ReadFileContent 读取文件内容
func (f *FileUtils) ReadFileContent(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %v", err)
	}
	return string(content), nil
}

// WriteFileContent 写入文件内容
func (f *FileUtils) WriteFileContent(filePath, content string) error {
	// 确保目录存在
	if err := f.EnsureDir(filepath.Dir(filePath)); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	// 写入文件
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	return nil
}

// ExtractSlugFromFile 从标签页面文件中提取现有的slug
func (f *FileUtils) ExtractSlugFromFile(filePath string) string {
	if !f.FileExists(filePath) {
		return ""
	}

	file, err := os.Open(filePath)
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	inFrontMatter := false

	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == "---" {
			if !inFrontMatter {
				inFrontMatter = true
				continue
			} else {
				break
			}
		}

		if inFrontMatter && strings.HasPrefix(strings.TrimSpace(line), "slug:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				slug := strings.TrimSpace(parts[1])
				slug = strings.Trim(slug, "\"'")
				return slug
			}
		}
	}

	return ""
}

// GenerateTagContent 生成标签页面内容
func (f *FileUtils) GenerateTagContent(tagName, slug string) string {
	return fmt.Sprintf(`---
title: %s
slug: "%s"
---
`, tagName, slug)
}

// BuildTargetFilePath 根据语言构建目标文件路径
func (f *FileUtils) BuildTargetFilePath(originalPath, targetLang string) string {
	dir := filepath.Dir(originalPath)
	baseName := filepath.Base(originalPath)

	if !strings.HasSuffix(baseName, ".md") {
		return ""
	}

	switch targetLang {
	case "ja":
		return filepath.Join(dir, "index.ja.md")
	case "ko":
		return filepath.Join(dir, "index.ko.md")
	default: // "en" 或其他
		return filepath.Join(dir, "index.en.md")
	}
}
