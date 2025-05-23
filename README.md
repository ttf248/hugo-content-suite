# 标签扫描器

这个 Go 脚本用于扫描 Hugo 博客中的文章，解析并统计标签信息。

## 功能特性

- 扫描 `content/post` 目录下的所有 Markdown 文件
- 解析文章的 Front Matter，提取标签和分类信息
- 统计每个标签的使用次数
- 显示详细的标签分布信息
- 找出没有标签的文章
- 统计分类信息

## 使用方法

1. 进入脚本目录：
```bash
cd scripts/tag-scanner
```

2. 运行脚本：
```bash
go run main.go
```

3. 或者指定自定义目录：
```bash
go run main.go /path/to/content/post
```

## 输出示例

