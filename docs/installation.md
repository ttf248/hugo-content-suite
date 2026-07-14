# 安装与配置

需要 Go 1.22+ 和可访问的 LM Studio OpenAI 兼容接口。进入 `hugo-content-suite` 目录后执行 `go run .`。

`config.json` 已随项目提供。重点字段：

```json
{
  "lm_studio": { "url": "http://localhost:1234/v1/chat/completions", "model": "your-model", "timeout_seconds": 30 },
  "paths": {
    "default_content_dir": "../../content/post",
    "tags_dir": "../tags",
    "runtime_dir": ".hugo-content-suite"
  }
}
```

相对路径一律以 `config.json` 所在目录为基准。`runtime_dir` 保存缓存和日志，建议保持在 Git 忽略范围内。工具不会在缺少配置时自动创建文件，以免误在错误目录写入配置。
