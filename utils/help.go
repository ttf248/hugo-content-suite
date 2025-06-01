package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// IsEnglishOnly 检查字符串是否只包含英文字符
func IsEnglishOnly(s string) bool {
	for _, r := range s {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '-' || r == '_' || r == ' ') {
			return false
		}
	}
	return true
}

// NormalizeSlug 标准化slug格式
func NormalizeSlug(s string) string {
	// 转为小写
	s = strings.ToLower(s)

	// 移除引号和其他特殊字符
	s = strings.Trim(s, "\"'`")

	// 替换空格为连字符
	s = strings.ReplaceAll(s, " ", "-")

	// 移除非法字符，只保留字母、数字和连字符
	reg := regexp.MustCompile(`[^a-z0-9\-]`)
	s = reg.ReplaceAllString(s, "")

	// 移除多个连续的连字符
	reg = regexp.MustCompile(`-+`)
	s = reg.ReplaceAllString(s, "-")

	// 移除开头和结尾的连字符
	s = strings.Trim(s, "-")

	return s
}

// ContainsChinese 检查文本是否包含中文
func ContainsChinese(text string) bool {
	for _, r := range text {
		if r >= 0x4e00 && r <= 0x9fff {
			return true
		}
	}
	return false
}

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

// EnsureDir 确保目录存在
func EnsureDir(dirPath string) error {
	return os.MkdirAll(dirPath, 0755)
}

// FileExists 检查文件是否存在
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

// WriteFileContent 写入文件内容
func WriteFileContent(filePath, content string) error {
	// 确保目录存在
	if err := EnsureDir(filepath.Dir(filePath)); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	// 写入文件
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	return nil
}

// ReadFileContent 读取文件内容
func ReadFileContent(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %v", err)
	}
	return string(content), nil
}

// BuildTargetFilePath 根据语言构建目标文件路径
func BuildTargetFilePath(originalPath, targetLang string) string {
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
	case "fr":
		return filepath.Join(dir, "index.fr.md")
	case "ru":
		return filepath.Join(dir, "index.ru.md")
	case "hi":
		return filepath.Join(dir, "index.hi.md")
	default: // "en" 或其他
		return filepath.Join(dir, "index.en.md")
	}
}
