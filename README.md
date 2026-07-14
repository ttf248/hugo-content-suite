# Hugo Content Suite

面向 Hugo 内容目录的交互式管理工具。它使用 LM Studio 的 OpenAI 兼容接口生成标签 slug、文章 slug 和多语言译文。

## 启动

在本目录执行：

```powershell
go run .
```

也可将内容目录作为唯一的位置参数传入：

```powershell
go run . ..\..\content\post
```

程序读取同目录的 `config.json`。所有相对路径均相对此配置文件解析；缓存和日志默认写入 `paths.runtime_dir`，该运行目录不纳入版本控制。

## 菜单能力

- `.`：仅处理新增内容，依次生成标签页、缺失 slug 和缺失译文。
- `1`：预览并生成或更新标签页。
- `2`：预览并生成或更新文章 slug。
- `3`：预览并翻译文章。
- `4`：删除一种语言的译文。必须输入语言编号，再输入完整语言代码确认；永不删除 `index.md` 源文。
- `0`：退出。

翻译、标签和 slug 操作要求配置的 LM Studio 服务可访问；扫描与删除不依赖该服务。

## 验证

```powershell
go test ./...
go vet ./...
```

详细的配置与操作步骤见 [安装说明](docs/installation.md) 和 [使用说明](docs/usage.md)。
