package models

type Article struct {
	FilePath string
	Title    string
	Tags     []string
	Category string
	Date     string
	// 新增：内容信息
	FrontMatter string // 原始前置信息
	BodyContent string // 正文内容
	FullContent string // 完整文件内容
	CharCount   int    // 正文字符数
}

type TagStats struct {
	Name  string
	Count int
	Files []string
}

type CategoryStats struct {
	Name  string
	Count int
}
