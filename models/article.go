package models

type Article struct {
	FilePath   string
	Title      string
	Subtitle   string
	Summary    string
	Tags       []string
	Categories []string
	Date       string
	LastMod    string
	Featured   bool
	Draft      bool
	Slug       string
	// 新增：内容信息
	FrontMatter string   // 原始前置信息
	BodyContent []string // 分段后的正文内容
	CharCount   int      // 正文字符数
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
