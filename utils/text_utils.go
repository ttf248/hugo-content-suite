package utils

import (
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
