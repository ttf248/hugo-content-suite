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
- **LM Studio**（可选）：用于 AI 翻译功能
  - 支持的模型：gemma-3-12b-it 或其他兼容模型
  - 默认地址：`http://172.19.192.1:2234`

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

### 3. 配置 LM Studio（可选）
如果要使用 AI 翻译功能，需要配置 LM Studio：

```go
// 在 translator/llm_translator.go 中修改配置
const (
    LMStudioURL = "http://your-lm-studio-url:port/v1/chat/completions"
    ModelName   = "your-model-name"
)
```

### 4. 验证安装
```bash
go run main.go --help
```

## 🚀 快速开始

### 基础用法
```bash
# 使用默认目录 (../../content/post)
go run main.go

# 指定自定义目录
go run main.go /path/to/your/hugo/content/post
```

### 典型工作流程

1. **扫描分析**
   ```bash
   go run main.go
   ```
   查看项目统计信息和标签分布

2. **预览标签页面生成**
   - 选择菜单选项 `5. 预览标签页面生成`
   - 查看将要生成的标签页面列表

3. **生成标签页面**
   - 选择菜单选项 `6. 生成标签页面文件`
   - 确认后自动生成所有标签页面

4. **预览文章 Slug**
   - 选择菜单选项 `7. 预览文章Slug生成`
   - 查看将要添加/更新的 slug

5. **生成文章 Slug**
   - 选择菜单选项 `8. 生成文章Slug`
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
=== 标签使用统计 ===
┌──────┬─────────────┬──────────┬──────────┐
│ 排名 │   标签名    │ 使用次数 │ 频率级别 │
├──────┼─────────────┼──────────┼──────────┤
│  1   │ JavaScript  │    15    │   高频   │
│  2   │   前端开发  │    12    │   高频   │
│  3   │    React    │    10    │   高频   │
└──────┴─────────────┴──────────┴──────────┘
```

#### 标签页面预览
```
=== 标签页面生成预览 ===
┌─────────────┬───────────┬─────────────────────┬────────┬──────┐
│   标签名    │ 当前Slug  │       新Slug        │ 文章数 │ 状态 │
├─────────────┼───────────┼─────────────────────┼────────┼──────┤
│ JavaScript  │    无     │     javascript      │   15   │ 新建 │
│  前端开发   │    无     │ frontend-development│   12   │ 新建 │
│    React    │ react-js  │       react         │   10   │ 更新 │
└─────────────┴───────────┴─────────────────────┴────────┴──────┘
```

## 🏗️ 项目架构

### 目录结构
```
hugo-slug-auto/
├── main.go                    # 主程序入口
├── go.mod                     # Go 模块配置
├── README.md                  # 项目文档
├── .vscode/                   # VS Code 配置
│   ├── settings.json
│   └── launch.json
├── models/                    # 数据模型
│   └── article.go
├── scanner/                   # 文章扫描器
│   └── parser.go
├── stats/                     # 统计计算
│   └── calculator.go
├── display/                   # 显示层
│   └── tables.go
├── translator/                # 翻译服务
│   └── llm_translator.go
└── generator/                 # 页面生成器
    ├── page_generator.go
    └── article_slug_generator.go
```

### 核心组件

#### 1. 文章扫描器 (scanner)
- **功能**：递归扫描目录，解析 Markdown 文件
- **支持格式**：YAML Front Matter
- **提取内容**：标题、标签、分类、日期等元数据

#### 2. 统计计算器 (stats)
- **标签统计**：计算使用频率，生成排序列表
- **分类统计**：分析文章分类分布
- **文章分析**：识别无标签文章，频率分组

#### 3. 翻译服务 (translator)
- **AI 翻译**：集成 LM Studio 进行智能翻译
- **备用翻译**：内置常用标签映射表
- **Slug 规范化**：生成符合 URL 标准的 slug

#### 4. 页面生成器 (generator)
- **标签页面生成**：创建 Hugo 标签页面文件
- **文章 Slug 生成**：为文章添加/更新 slug
- **批量处理**：支持大量文件的批量操作

#### 5. 显示层 (display)
- **表格显示**：使用 tablewriter 生成美观的表格
- **颜色输出**：使用 fatih/color 提供彩色输出
- **交互界面**：实现用户友好的命令行界面

## 🔧 配置选项

### LM Studio 配置
```go
// translator/llm_translator.go
const (
    LMStudioURL = "http://172.19.192.1:2234/v1/chat/completions"
    ModelName   = "gemma-3-12b-it"
)
```

### 备用翻译映射
```go
// translator/llm_translator.go
fallbackTranslations := map[string]string{
    "人工智能": "artificial-intelligence",
    "机器学习": "machine-learning",
    "深度学习": "deep-learning",
    "前端开发": "frontend-development",
    "后端开发": "backend-development",
    // ... 更多映射
}
```

### 文件路径配置
```go
// main.go
contentDir := "../../content/post"  // 默认文章目录
tagsDir := "content/tags"           // 标签页面目录
```

## 🐛 故障排除

### 常见问题

**1. LM Studio 连接失败**
```
错误：无法连接到LM Studio
解决：检查 LM Studio 是否运行，确认 URL 和端口正确
```

**2. 文件权限错误**
```
错误：写入文件失败
解决：确保程序对目标目录有写入权限
```

**3. Front Matter 解析失败**
```
错误：无法解析文章元数据
解决：检查 Markdown 文件的 YAML 格式是否正确
```

**4. 目录不存在**
```
错误：找不到指定目录
解决：确认文章目录路径正确，或使用绝对路径
```

### 调试模式

使用 VS Code 调试：
1. 打开项目在 VS Code 中
2. 按 F5 启动调试
3. 在代码中设置断点进行调试

命令行调试：
```bash
go run -race main.go  # 检测竞态条件
go run main.go 2>&1 | tee output.log  # 记录输出日志
```

## 🤝 贡献指南

### 开发环境设置
1. Fork 项目
2. 创建功能分支
3. 提交更改
4. 创建 Pull Request

### 代码规范
- 遵循 Go 官方代码规范
- 添加适当的注释和文档
- 编写单元测试
- 确保代码通过 `go vet` 和 `golint` 检查

### 测试
```bash
go test ./...                    # 运行所有测试
go test -v ./scanner            # 运行特定包的测试
go test -cover ./...            # 生成测试覆盖率报告
```

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🔗 相关链接

- [Hugo 官方文档](https://gohugo.io/documentation/)
- [LM Studio](https://lmstudio.ai/)
- [Go 官方文档](https://golang.org/doc/)

## 📞 支持与反馈

如果您遇到问题或有改进建议，请：

1. 查看 [常见问题](#故障排除)
2. 搜索现有的 [Issues](issues)
3. 创建新的 Issue 描述问题
4. 提交 Pull Request 贡献代码

---

⭐ 如果这个项目对您有帮助，请给它一个星标！

