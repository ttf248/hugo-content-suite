# Hugo 标签与 Slug 自动化工具

一个强大的 Go 工具，用于自动化管理 Hugo 博客的标签系统和 URL slug。支持批量扫描、智能翻译、标签页面生成和文章 slug 优化。

## 🚀 功能特性

### 📊 文章分析与统计
- **智能扫描**：自动扫描 Hugo 博客目录下的所有 Markdown 文件
- **Front Matter 解析**：准确解析 YAML 格式的文章元数据
- **标签统计**：统计标签使用频率，按高频(≥5篇)、中频(2-4篇)、低频(1篇)分组
- **分类统计**：分析文章分类分布情况
- **无标签检测**：识别缺少标签的文章

### 🌐 智能翻译系统
- **AI 翻译**：集成 LM Studio 本地大语言模型进行中文标签翻译
- **备用方案**：内置常用标签映射表，确保翻译服务不可用时的可靠性
- **Slug 规范化**：自动生成符合 URL 标准的英文 slug
- **批量处理**：支持大量标签的批量翻译

### 📄 自动化页面生成
- **标签页面**：为每个标签自动生成 Hugo 标签页面文件
- **目录结构**：按 Hugo 规范创建 `content/tags/[标签名]/_index.md` 结构
- **元数据完整**：包含标题、slug、描述等完整的页面元数据
- **增量更新**：智能检测现有页面，支持新建和更新操作

### 🔗 文章 Slug 管理
- **Slug 生成**：为文章标题自动生成英文 slug
- **批量更新**：批量为所有文章添加或更新 slug 字段
- **规范化处理**：确保 slug 符合 URL 标准（小写、连字符分隔）
- **冲突避免**：智能处理重复和特殊字符

### 🎯 交互式操作
- **菜单驱动**：直观的交互式命令行界面
- **实时预览**：操作前预览将要进行的更改
- **确认机制**：重要操作前提供确认提示
- **进度显示**：批量操作时显示详细进度信息

## 📋 系统要求

- **Go 版本**：Go 1.21 或更高版本
- **操作系统**：Windows、macOS、Linux
- **依赖库**：
  - `github.com/fatih/color` - 彩色终端输出
  - `github.com/olekukonko/tablewriter` - 表格显示
- **LM Studio**（可选）：用于 AI 翻译功能

## 🛠️ 安装与配置

### 1. 克隆项目
```bash
git clone <repository-url>
cd hugo-slug-auto
```

### 2. 安装依赖
```bash
go mod tidy
```

### 3. 构建项目
```bash
go build -o tag-scanner main.go
```

### 4. 验证安装
```bash
./tag-scanner
```

## 🚀 快速开始

### 基础用法
```bash
# 使用默认目录 (../../content/post)
./tag-scanner

# 指定自定义目录
./tag-scanner /path/to/your/hugo/content/post
```

### 典型工作流程

1. **扫描分析**
   ```bash
   ./tag-scanner
   ```
   查看项目统计信息和标签分布

2. **预览标签页面生成**
   - 选择菜单选项 `5. 预览标签页面生成`
   - 查看将要生成的标签页面列表

3. **生成标签页面**
   - 选择菜单选项 `6. 生成标签页面文件`
   - 选择处理模式（新增/更新/全部）
   - 确认后自动生成所有标签页面

4. **预览文章 Slug**
   - 选择菜单选项 `7. 预览文章Slug生成`
   - 查看将要添加/更新的 slug

5. **生成文章 Slug**
   - 选择菜单选项 `8. 生成文章Slug`
   - 选择处理模式（缺失/更新/全部）
   - 确认后批量处理所有文章

## 📖 详细使用说明

### 交互式菜单功能

| 选项 | 功能 | 说明 |
|------|------|------|
| 1 | 查看所有标签 | 显示完整的标签使用统计 |
| 2 | 查看特定标签详情 | 查看指定标签的使用文章列表 |
| 3 | 查看所有无标签文章 | 列出所有缺少标签的文章 |
| 4 | 查看标签频率分组 | 按使用频率分组显示标签 |
| 5 | 预览标签页面生成 | 预览即将生成的标签页面 |
| 6 | 生成标签页面文件 | 实际生成标签页面文件 |
| 7 | 预览文章Slug生成 | 预览文章 slug 生成结果 |
| 8 | 生成文章Slug | 批量生成/更新文章 slug |
| 0 | 退出 | 退出程序 |

### 处理模式说明

#### 标签页面生成模式
- **仅新增**：只生成不存在的标签页面
- **仅更新**：只更新已存在的标签页面
- **全部处理**：新建和更新所有标签页面

#### 文章 Slug 生成模式
- **仅新增**：只为缺少 slug 的文章添加 slug
- **仅更新**：只更新已有但不匹配的 slug
- **全部处理**：处理所有需要添加或更新 slug 的文章

### 输出示例

#### 统计概览
```
=== 博客文章统计概览 ===
┌─────────────────┬──────┐
│     统计项      │ 数量 │
├─────────────────┼──────┤
│    总文章数     │ 156  │
│    标签总数     │  89  │
│    分类总数     │  12  │
│ 高频标签 (≥5篇) │   8  │
│ 中频标签 (2-4篇) │  25  │
│ 低频标签 (1篇)  │  56  │
└─────────────────┴──────┘
```

#### 标签统计
```
=== 标签使用统计（前20个）===
┌──────┬─────────────┬──────────┬──────────┐
│ 排名 │   标签名    │ 使用次数 │ 频率级别 │
├──────┼─────────────┼──────────┼──────────┤
│  1   │ JavaScript  │    15    │   高频   │
│  2   │   前端开发  │    12    │   高频   │
│  3   │    React    │    10    │   高频   │
└──────┴─────────────┴──────────┴──────────┘
```

## 🏗️ 项目架构

### 目录结构
```
hugo-slug-auto/
├── main.go                    # 主程序入口和交互式菜单
├── go.mod                     # Go 模块配置
├── README.md                  # 项目文档
├── models/                    # 数据模型定义
│   └── article.go            # Article、TagStats、CategoryStats 结构
├── scanner/                   # 文章扫描功能（需要实现）
├── stats/                     # 统计计算功能（需要实现）
├── display/                   # 显示层功能（需要实现）
├── generator/                 # 页面生成功能（需要实现）
└── translator/                # 翻译服务功能（需要实现）
```

### 核心数据结构

#### Article 模型
```go
type Article struct {
    FilePath string   // 文件路径
    Title    string   // 文章标题
    Tags     []string // 标签列表
    Category string   // 文章分类
    Date     string   // 发布日期
}
```

#### 统计模型
```go
type TagStats struct {
    Name  string   // 标签名称
    Count int      // 使用次数
    Files []string // 使用该标签的文件列表
}

type CategoryStats struct {
    Name  string // 分类名称
    Count int    // 文章数量
}
```

### 主要功能模块

#### 1. 文章扫描器 (scanner)
- 递归扫描指定目录下的 Markdown 文件
- 解析 YAML Front Matter 提取元数据
- 返回 Article 结构体数组

#### 2. 统计计算器 (stats)
- `CalculateTagStats()` - 计算标签使用统计
- `CalculateCategoryStats()` - 计算分类统计
- `FindNoTagArticles()` - 查找无标签文章
- `GroupTagsByFrequency()` - 按频率分组标签

#### 3. 显示层 (display)
- `DisplaySummary()` - 显示统计概览
- `DisplayTagStats()` - 显示标签统计表格
- `DisplayCategoryStats()` - 显示分类统计
- `DisplayNoTagArticles()` - 显示无标签文章列表
- `DisplayTagDetails()` - 显示特定标签详情

#### 4. 页面生成器 (generator)
- `TagPageGenerator` - 标签页面生成器
- `ArticleSlugGenerator` - 文章 Slug 生成器
- 支持预览和实际生成功能

#### 5. 翻译服务 (translator)

##### LM Studio 集成
工具集成了 LM Studio 本地大语言模型，提供智能翻译服务：

- **本地运行**：无需互联网连接，保护数据隐私
- **智能翻译**：理解上下文，生成准确的技术术语翻译
- **自动降级**：连接失败时自动切换到备用翻译方案

##### 备用翻译映射
内置常用标签翻译映射表，确保翻译功能的可靠性：

```go
fallbackTranslations := map[string]string{
    "人工智能": "artificial-intelligence",
    "机器学习": "machine-learning",
    "深度学习": "deep-learning",
    "前端开发": "frontend-development",
    "后端开发": "backend-development",
    "JavaScript": "javascript",
    "Python": "python",
    "Go": "golang",
    // ...更多映射
}
```

## ⚙️ 配置选项

### 默认路径配置
```go
// 默认文章目录
contentDir := "../../content/post"

// 可通过命令行参数覆盖
// ./tag-scanner /your/custom/path
```

### 标签频率分组阈值
- 高频标签：≥5 篇文章
- 中频标签：2-4 篇文章  
- 低频标签：1 篇文章

## 🔧 开发指南

### 环境要求
- Go 1.21+
- 支持的终端（用于彩色输出）

### 构建和测试
```bash
# 构建
go build -o tag-scanner main.go

# 运行测试（当实现测试后）
go test ./...

# 代码格式化
go fmt ./...

# 静态检查
go vet ./...
```

### VS Code 配置
项目包含 VS Code 配置，支持：
- 调试启动配置
- Go 语言设置
- 推荐的扩展

## 📄 许可证

本项目采用 MIT 许可证。

## 🤝 贡献

欢迎提交 Issue 和 Pull Request 来改进这个项目！

---

⭐ 如果这个项目对您有帮助，请给它一个星标！

