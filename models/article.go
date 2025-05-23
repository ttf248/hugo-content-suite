package models

type Article struct {
	FilePath string
	Title    string
	Tags     []string
	Category string
	Date     string
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
