# 配置文件说明

[English](configuration_en.md) | 中文

## 配置文件概述

本工具使用YAML格式的配置文件来管理各种设置，支持环境变量覆盖和热重载功能。

## 配置文件位置

默认配置文件路径：`config.yaml`

可通过以下方式指定配置文件：
```bash
# 使用环境变量
export CONFIG_FILE="/path/to/config.yaml"

# 使用命令行参数
go run main.go --config=/path/to/config.yaml
```

## 完整配置示例

```yaml
# config.yaml
# LM Studio AI翻译配置
lm_studio:
  url: "http://localhost:2234/v1/chat/completions"
  model: "gemma-3-12b-it"
  timeout: 30s
  max_retries: 3
  retry_delay: 1s
  api_key: ""  # 如果需要API密钥

# 缓存系统配置
cache:
  directory: "./cache"
  file_name: "tag_translations_cache.json"
  auto_save: true
  save_interval: 5m
  max_entries: 10000
  expire_after: 720h  # 30天过期

# 日志系统配置
logging:
  level: "INFO"  # DEBUG, INFO, WARN, ERROR
  file_path: "./logs/app.log"
  max_size: 100MB
  max_backups: 5
  max_age: 30  # 天数
  console_output: true
  json_format: false
  
# 性能监控配置
performance:
  enable_monitoring: true
  metrics_interval: 10s
  memory_threshold: 500MB
  cpu_threshold: 80
  enable_profiling: false
  profile_port: 6060

# 文章扫描配置
scanner:
  content_dir: "./content/post"
  include_patterns:
    - "*.md"
    - "*.markdown"
  exclude_patterns:
    - "_*"
    - "draft*"
  max_concurrent: 10

# 标签页面生成配置
generator:
  tags_dir: "./content/tags"
  template_file: ""  # 自定义模板文件
  overwrite_existing: false
  backup_existing: true

# 翻译配置
translator:
  fallback_enabled: true
  custom_mappings:
    "人工智能": "artificial-intelligence"
    "机器学习": "machine-learning"
    "前端开发": "frontend-development"
  batch_size: 10
  request_delay: 500ms

# 用户界面配置
ui:
  language: "zh-CN"  # zh-CN, en-US
  color_output: true
  progress_bar: true
  table_max_width: 120
  
# 安全配置
security:
  max_file_size: 10MB
  allowed_extensions:
    - ".md"
    - ".markdown"
  sandbox_mode: false
```

## 配置项详细说明

### LM Studio 配置 (lm_studio)

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| url | string | "http://localhost:2234/v1/chat/completions" | LM Studio API地址 |
| model | string | "gemma-3-12b-it" | 使用的AI模型名称 |
| timeout | duration | 30s | 请求超时时间 |
| max_retries | int | 3 | 最大重试次数 |
| retry_delay | duration | 1s | 重试间隔时间 |
| api_key | string | "" | API密钥（如果需要） |

### 缓存配置 (cache)

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| directory | string | "./cache" | 缓存目录 |
| file_name | string | "tag_translations_cache.json" | 缓存文件名 |
| auto_save | bool | true | 自动保存缓存 |
| save_interval | duration | 5m | 自动保存间隔 |
| max_entries | int | 10000 | 最大缓存条目数 |
| expire_after | duration | 720h | 缓存过期时间 |

### 日志配置 (logging)

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| level | string | "INFO" | 日志级别 |
| file_path | string | "./logs/app.log" | 日志文件路径 |
| max_size | string | 100MB | 单个日志文件最大大小 |
| max_backups | int | 5 | 保留的日志文件数量 |
| max_age | int | 30 | 日志文件保留天数 |
| console_output | bool | true | 是否输出到控制台 |
| json_format | bool | false | 是否使用JSON格式 |

### 性能监控配置 (performance)

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| enable_monitoring | bool | true | 启用性能监控 |
| metrics_interval | duration | 10s | 指标收集间隔 |
| memory_threshold | string | 500MB | 内存使用警告阈值 |
| cpu_threshold | int | 80 | CPU使用警告阈值(%) |
| enable_profiling | bool | false | 启用性能分析 |
| profile_port | int | 6060 | 性能分析端口 |

## 环境变量覆盖

所有配置项都可以通过环境变量覆盖，变量名格式：`HUGO_SLUG_<SECTION>_<KEY>`

### 常用环境变量

```bash
# LM Studio配置
export HUGO_SLUG_LM_STUDIO_URL="http://192.168.1.100:2234/v1/chat/completions"
export HUGO_SLUG_LM_STUDIO_MODEL="gpt-4"
export HUGO_SLUG_LM_STUDIO_TIMEOUT="60s"

# 日志配置
export HUGO_SLUG_LOGGING_LEVEL="DEBUG"
export HUGO_SLUG_LOGGING_FILE_PATH="./logs/debug.log"

# 缓存配置
export HUGO_SLUG_CACHE_DIRECTORY="./custom_cache"
export HUGO_SLUG_CACHE_MAX_ENTRIES="20000"

# 性能配置
export HUGO_SLUG_PERFORMANCE_ENABLE_MONITORING="false"
```

### 嵌套配置覆盖

对于嵌套的配置项，使用下划线分隔：

```bash
# 覆盖 translator.custom_mappings
export HUGO_SLUG_TRANSLATOR_CUSTOM_MAPPINGS_AI="artificial-intelligence"

# 覆盖 ui.language
export HUGO_SLUG_UI_LANGUAGE="en-US"
```

## 配置验证

程序启动时会自动验证配置文件：

### 验证规则
- 必需字段检查
- 数据类型验证
- 取值范围检查
- 文件路径有效性
- 网络地址格式验证

### 验证错误示例
```
配置验证失败:
- lm_studio.timeout: 必须是有效的时间格式 (如: 30s, 1m, 2h)
- logging.level: 必须是 DEBUG, INFO, WARN, ERROR 中的一个
- cache.max_entries: 必须是正整数
- performance.memory_threshold: 必须是有效的字节大小 (如: 100MB, 1GB)
```

## 配置热重载

支持在运行时重新加载配置文件：

### 触发方式
1. **信号触发**（Linux/macOS）：
   ```bash
   kill -SIGHUP <进程ID>
   ```

2. **文件监控**：
   程序会自动监控配置文件变化并重新加载

3. **API接口**：
   ```bash
   curl -X POST http://localhost:8080/api/reload-config
   ```

### 热重载限制
- 某些配置项需要重启程序才能生效（如日志文件路径）
- 重载过程中会暂停相关功能
- 配置验证失败时会保持原有配置

## 配置模板生成

生成默认配置文件：

```bash
# 生成默认配置
go run main.go --generate-config

# 生成带注释的详细配置
go run main.go --generate-config --with-comments

# 指定输出文件
go run main.go --generate-config --output=my-config.yaml
```

## 最佳实践

### 开发环境
```yaml
logging:
  level: "DEBUG"
  console_output: true
  
performance:
  enable_monitoring: true
  enable_profiling: true
```

### 生产环境
```yaml
logging:
  level: "INFO"
  console_output: false
  max_backups: 10
  
performance:
  enable_monitoring: true
  enable_profiling: false
  
security:
  sandbox_mode: true
```

### 性能优化
```yaml
cache:
  max_entries: 50000
  save_interval: 1m
  
translator:
  batch_size: 20
  request_delay: 200ms
  
scanner:
  max_concurrent: 20
```

## 故障排除

### 常见配置问题

1. **配置文件不存在**
   ```
   解决方案：使用 --generate-config 生成默认配置
   ```

2. **YAML格式错误**
   ```
   解决方案：检查缩进和语法，使用在线YAML验证器
   ```

3. **环境变量未生效**
   ```
   解决方案：检查变量名格式，确保程序重启后生效
   ```

4. **权限问题**
   ```
   解决方案：检查配置文件和目录的读写权限
   ```

## 相关文档

- [安装配置指南](installation.md)
- [日志系统指南](logging.md)
- [性能监控指南](performance.md)
- [故障排除](troubleshooting.md)
