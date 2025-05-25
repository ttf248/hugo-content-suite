# 故障排除指南

## 常见问题

### 1. LM Studio 连接问题

#### 问题现象
```
警告：无法连接到LM Studio (dial tcp 172.19.192.1:2234: connectex: No connection could be made...)
```

#### 解决方法
1. **检查LM Studio状态**
   ```bash
   # 确保LM Studio正在运行并已加载模型
   ```

2. **检查网络配置**
   ```go
   // 修改 translator/llm_translator.go 中的URL
   const LMStudioURL = "http://localhost:2234/v1/chat/completions"
   ```

3. **测试连接**
   ```bash
   curl -X POST http://localhost:2234/v1/chat/completions \
     -H "Content-Type: application/json" \
     -d '{"model":"your-model","messages":[{"role":"user","content":"test"}]}'
   ```

### 2. 文件权限问题

#### 问题现象
```
创建tags目录失败: permission denied
写入缓存文件失败: permission denied
```

#### 解决方法
1. **检查目录权限**
   ```bash
   # Linux/macOS
   chmod 755 /path/to/content
   chmod 644 /path/to/content/**/*.md
   
   # Windows
   # 右键 -> 属性 -> 安全 -> 编辑权限
   ```

2. **使用管理员权限运行**
   ```bash
   # Linux/macOS
   sudo go run main.go
   
   # Windows
   # 以管理员身份运行命令提示符
   ```

### 3. Go 模块问题

#### 问题现象
```
go: module hugo-content-suite: cannot find module providing package
```

#### 解决方法
1. **重新初始化模块**
   ```bash
   go mod init hugo-content-suite
   go mod tidy
   ```

2. **清理模块缓存**
   ```bash
   go clean -modcache
   go mod download
   ```

3. **检查Go版本**
   ```bash
   go version
   # 确保版本 >= 1.21
   ```

### 4. 文章解析问题

#### 问题现象
- 标签没有被正确识别
- Front Matter解析失败

#### 解决方法
1. **检查文件格式**
   ```yaml
   ---
   title: "文章标题"
   tags: ["标签1", "标签2"]  # 确保使用数组格式
   categories: ["分类"]
   date: 2024-01-01
   ---
   ```

2. **支持的标签格式**
   ```yaml
   # 数组格式（推荐）
   tags: ["tag1", "tag2"]
   
   # YAML列表格式
   tags:
     - tag1
     - tag2
   ```

### 5. 缓存相关问题

#### 缓存文件损坏
```bash
# 删除损坏的缓存文件
rm tag_translations_cache.json

# 或使用程序内置的清空功能
# 菜单选择: 9. 清空翻译缓存
```

#### 缓存路径问题
```go
// 修改缓存路径在 translator/llm_translator.go
cache: NewTranslationCache("/custom/cache/path"),
```

### 6. 翻译质量问题

#### AI翻译结果不理想
1. **更换模型**
   ```go
   const ModelName = "better-model-name"
   ```

2. **调整提示词**
   ```go
   // 在 TranslateToSlug 方法中修改 prompt
   prompt := fmt.Sprintf(`更详细的翻译指令...`)
   ```

3. **添加预定义映射**
   ```go
   // 在 fallbackSlug 函数中添加
   fallbackTranslations := map[string]string{
       "你的标签": "your-tag",
       // 添加更多映射
   }
   ```

## 调试技巧

### 1. 启用详细日志

```go
// 在需要调试的地方添加
fmt.Printf("调试信息: %+v\n", variable)
```

### 2. 检查文件内容

```bash
# 查看生成的标签页面
find content/tags -name "_index.md" -exec head -10 {} \;

# 检查缓存文件
cat tag_translations_cache.json | jq .
```

### 3. 测试单个功能

```go
// 创建简单的测试文件
func main() {
    translator := translator.NewLLMTranslator()
    result, err := translator.TranslateToSlug("测试标签")
    fmt.Printf("结果: %s, 错误: %v\n", result, err)
}
```

## 性能优化

### 1. 大量文章处理

```bash
# 分批处理大量文章
# 先处理部分目录测试效果
go run main.go ./content/post/2024

# 确认无误后处理全部
go run main.go ./content/post
```

### 2. 网络超时调整

```go
// 增加超时时间
client: &http.Client{
    Timeout: 60 * time.Second,  // 调整为60秒
},
```

### 3. 并发处理

```go
// 目前程序是串行处理，如需并发请小心处理缓存同步
```

## 环境特定问题

### Windows 环境

1. **路径分隔符问题**
   ```go
   // 使用 filepath.Join 而不是手动拼接
   path := filepath.Join("content", "post")
   ```

2. **字符编码问题**
   ```bash
   # 确保终端支持UTF-8
   chcp 65001
   ```

### macOS 环境

1. **Homebrew Go版本**
   ```bash
   brew upgrade go
   go version
   ```

### Linux 环境

1. **依赖包安装**
   ```bash
   # Ubuntu/Debian
   sudo apt update
   sudo apt install golang-go

   # CentOS/RHEL
   sudo yum install golang
   ```

## 获取帮助

### 1. 查看日志
程序会输出详细的操作日志，注意查看错误信息。

### 2. 检查配置
确认所有配置项（URL、路径、模型名称）都正确。

### 3. 社区支持
- GitHub Issues: 提交问题和建议
- 文档更新: 发现问题请及时反馈

### 4. 联系方式
如果以上方法都无法解决问题，请：
1. 准备详细的错误信息
2. 说明操作系统和Go版本
3. 提供复现步骤
4. 通过GitHub Issues提交问题

## 预防措施

### 1. 备份重要文件
```bash
# 处理前备份content目录
cp -r content content_backup
```

### 2. 小范围测试
```bash
# 先在测试目录验证功能
mkdir test_content
cp content/post/sample.md test_content/
go run main.go test_content
```

### 3. 版本控制
```bash
# 使用Git跟踪变更
git add .
git commit -m "处理前的备份"
```

这样可以在出现问题时快速回滚到之前的状态。
