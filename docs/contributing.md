# 贡献指南

欢迎为 Hugo 标签自动化管理工具做出贡献！

## 如何贡献

### 🐛 报告Bug
1. 在 [GitHub Issues](https://github.com/ttf248/hugo-slug-auto/issues) 中搜索是否已有相同问题
2. 如果没有，创建新的Issue并提供：
   - 详细的问题描述
   - 复现步骤
   - 期待的行为
   - 实际的行为
   - 环境信息（操作系统、Go版本等）

### 💡 功能建议
1. 在Issues中详细描述新功能的需求
2. 说明为什么需要这个功能
3. 提供可能的实现思路

### 🔧 代码贡献

#### 开发环境设置
```bash
# 1. Fork 项目到你的GitHub账户
# 2. 克隆你的Fork
git clone https://github.com/your-username/hugo-slug-auto.git
cd hugo-slug-auto

# 3. 添加上游仓库
git remote add upstream https://github.com/ttf248/hugo-slug-auto.git

# 4. 安装依赖
go mod tidy

# 5. 运行测试
go test ./...
```

#### 开发流程
1. **创建功能分支**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **编写代码**
   - 遵循现有的代码风格
   - 添加必要的注释
   - 编写单元测试

3. **测试代码**
   ```bash
   go test ./...
   go run main.go  # 手动测试
   ```

4. **提交更改**
   ```bash
   git add .
   git commit -m "feat: 添加新功能描述"
   ```

5. **同步上游**
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

6. **推送并创建PR**
   ```bash
   git push origin feature/your-feature-name
   ```

## 代码规范

### Go 代码风格
- 使用 `gofmt` 格式化代码
- 遵循 [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- 使用有意义的变量和函数名
- 添加必要的错误处理

### 提交信息格式
使用 [Conventional Commits](https://www.conventionalcommits.org/) 格式：

```
<type>[optional scope]: <description>

[optional body]

[optional footer]
```

类型包括：
- `feat`: 新功能
- `fix`: 修复bug
- `docs`: 文档更新
- `style`: 代码格式修改
- `refactor`: 代码重构
- `test`: 添加测试
- `chore`: 构建过程或辅助工具的变动

示例：
```
feat(translator): 添加新的翻译引擎支持

- 支持OpenAI GPT API
- 添加配置选项
- 更新文档

Closes #123
```

### 代码注释
```go
// Package translator 提供AI翻译功能
package translator

// TranslateToSlug 将中文标签翻译为英文slug
// 参数 tag: 要翻译的中文标签
// 返回值: 英文slug和可能的错误
func (t *LLMTranslator) TranslateToSlug(tag string) (string, error) {
    // 实现逻辑...
}
```

## 项目结构

### 目录说明
```
hugo-slug-auto/
├── main.go              # 程序入口
├── models/              # 数据模型定义
├── scanner/             # 文件扫描和解析
├── stats/               # 统计分析功能
├── translator/          # AI翻译模块
├── generator/           # 内容生成器
├── display/             # 用户界面展示
├── docs/               # 项目文档
└── tests/              # 测试文件
```

### 模块设计原则
- **单一职责**: 每个模块专注于特定功能
- **松耦合**: 模块间通过接口交互
- **高内聚**: 相关功能组织在同一模块内

## 测试指南

### 单元测试
```go
// translator_test.go
func TestTranslateToSlug(t *testing.T) {
    translator := NewLLMTranslator()
    
    tests := []struct {
        input    string
        expected string
    }{
        {"人工智能", "artificial-intelligence"},
        {"机器学习", "machine-learning"},
    }
    
    for _, test := range tests {
        result, err := translator.TranslateToSlug(test.input)
        assert.NoError(t, err)
        assert.Equal(t, test.expected, result)
    }
}
```

### 集成测试
```bash
# 创建测试数据
mkdir test-content
echo '---\ntitle: "测试文章"\ntags: ["测试"]\n---\n测试内容' > test-content/test.md

# 运行程序测试
go run main.go test-content
```

## 新功能开发

### 添加新的翻译器
1. 在 `translator/` 目录下创建新文件
2. 实现 `Translator` 接口
3. 添加配置选项
4. 编写测试
5. 更新文档

### 添加新的生成器
1. 在 `generator/` 目录下创建新文件
2. 实现生成逻辑
3. 集成到主菜单
4. 添加预览功能
5. 编写测试

### 添加新的显示功能
1. 在 `display/` 目录下添加新函数
2. 使用统一的颜色和格式
3. 支持表格和列表显示
4. 考虑大数据量的分页显示

## 文档贡献

### 文档类型
- **README**: 快速入门指南
- **安装文档**: 详细的安装配置说明
- **使用文档**: 功能使用说明
- **API文档**: 接口说明
- **故障排除**: 常见问题解决方案

### 文档写作规范
- 使用清晰的标题结构
- 提供具体的代码示例
- 包含必要的截图或图表
- 保持内容的时效性

## PR 审查流程

### 提交PR前检查
- [ ] 代码通过所有测试
- [ ] 代码符合项目风格
- [ ] 添加了必要的文档
- [ ] 更新了CHANGELOG（如适用）

### 审查标准
1. **功能正确性**: 实现是否符合需求
2. **代码质量**: 可读性、可维护性
3. **性能影响**: 是否影响现有功能性能
4. **向后兼容**: 是否破坏现有API
5. **安全性**: 是否引入安全风险

### 审查流程
1. 自动化检查（CI/CD）
2. 代码审查
3. 功能测试
4. 文档审查
5. 合并到主分支

## 发布流程

### 版本号规则
遵循 [Semantic Versioning](https://semver.org/)：
- `MAJOR.MINOR.PATCH`
- MAJOR: 不兼容的API修改
- MINOR: 向后兼容的功能性新增
- PATCH: 向后兼容的问题修正

### 发布步骤
1. 更新版本号
2. 更新CHANGELOG
3. 创建Release Tag
4. 编写Release Notes
5. 发布到GitHub Releases

## 社区交流

### 讨论平台
- GitHub Issues: 问题报告和功能讨论
- GitHub Discussions: 一般性讨论和问答

### 行为准则
- 保持友善和尊重
- 专注于技术讨论
- 欢迎新手提问
- 分享知识和经验

## 致谢

感谢所有为项目做出贡献的开发者！你们的贡献让这个项目变得更好。

### 贡献者名单
- 维护在 [CONTRIBUTORS.md](CONTRIBUTORS.md) 文件中

---

再次感谢你对项目的关注和贡献！🎉
