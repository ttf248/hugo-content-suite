# 安装与配置

需要 Go 1.22+ 和可访问的模型服务。进入 `hugo-content-suite` 目录后执行 `go run .`。

`config.json` 已随项目提供。重点字段：

```json
{
  "active_model": "minimax-cn-token-plan",
  "models": [{
    "name": "minimax-cn-token-plan",
    "api_type": "anthropic_messages",
    "url": "https://api.minimaxi.com/anthropic/v1/messages",
    "model": "MiniMax-M2.5",
    "api_key": "",
    "api_key_env": "MINIMAX_API_KEY",
    "timeout_seconds": 60
  }],
  "paths": {
    "default_content_dir": "../../content/post",
    "tags_dir": "../tags",
    "runtime_dir": ".hugo-content-suite"
  }
}
```

相对路径一律以 `config.json` 所在目录为基准。`runtime_dir` 保存缓存和日志，建议保持在 Git 忽略范围内。工具不会在缺少配置时自动创建文件，以免误在错误目录写入配置。

`api_key_env` 优先读取系统环境变量；`api_key` 仅预留给本地未跟踪配置，切勿将真实密钥提交到仓库。启动后通过菜单 `5` 切换模型、菜单 `6` 测试当前模型。
