# Hugo Content Suite

面向 Hugo 内容目录的交互式管理工具。它支持 OpenAI Chat 与 Anthropic Messages 协议，生成标签 slug、文章 slug 和多语言译文。

## 启动

在本目录执行：

```powershell
go run .
```

也可将内容目录作为唯一的位置参数传入：

```powershell
go run . ..\..\content\post
```

保留默认交互菜单的同时，可用以下 CLI 模式直接执行“仅新增的一键处理”：

```powershell
go run . --process-new
go run . --process-new ..\..\content\post
```

VS Code 可选择 `Hugo Content Suite: 增量一键处理` 调试配置，以当前本地配置启动同一流程。

程序读取同目录的 `config.local.json`，并以受跟踪的 `config.example.json` 为基础合并本地覆盖项。所有相对路径均相对于本地配置文件解析；缓存和日志默认写入 `paths.runtime_dir`，该运行目录不纳入版本控制。

`models` 可配置多组模型，`active_model` 指定启动默认项。菜单 `5` 可在本次运行中切换，菜单 `6` 测试当前选择。MiniMax 优先读取本地配置的 `api_key`，并兼容 `api_key_env` 指向的环境变量回退；`config.local.json` 已被 Git 忽略，绝不提交真实密钥。

## 菜单能力

- `.`：仅处理新增内容，依次生成标签页、缺失 slug 和缺失译文。
- `1`：预览并生成或更新标签页。
- `2`：预览并生成或更新文章 slug。
- `3`：预览并翻译文章。
- `4`：删除一种语言的译文。必须输入语言编号，再输入完整语言代码确认；永不删除 `index.md` 源文。
- `5`：选择翻译模型。
- `6`：测试当前翻译模型的连通性。
- `0`：退出。

翻译、标签和 slug 操作要求当前模型服务可访问；扫描与删除不依赖该服务。

## 验证

```powershell
go test ./...
go vet ./...
```

详细的配置与操作步骤见 [安装说明](docs/installation.md) 和 [使用说明](docs/usage.md)。
