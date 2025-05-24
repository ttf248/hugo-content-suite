# Configuration Guide

English | [中文](configuration.md)

## Configuration Overview

This tool uses YAML format configuration files to manage various settings, with support for environment variable overrides and hot reload functionality.

## Configuration File Location

Default configuration file path: `config.yaml`

You can specify the configuration file in the following ways:
```bash
# Using environment variable
export CONFIG_FILE="/path/to/config.yaml"

# Using command line argument
go run main.go --config=/path/to/config.yaml
```

## Complete Configuration Example

```yaml
# config.yaml
# LM Studio AI Translation Configuration
lm_studio:
  url: "http://localhost:2234/v1/chat/completions"
  model: "gemma-3-12b-it"
  timeout: 30s
  max_retries: 3
  retry_delay: 1s
  api_key: ""  # If API key is required

# Cache System Configuration
cache:
  directory: "./cache"
  file_name: "tag_translations_cache.json"
  auto_save: true
  save_interval: 5m
  max_entries: 10000
  expire_after: 720h  # 30 days expiration

# Logging System Configuration
logging:
  level: "INFO"  # DEBUG, INFO, WARN, ERROR
  file_path: "./logs/app.log"
  max_size: 100MB
  max_backups: 5
  max_age: 30  # days
  console_output: true
  json_format: false
  
# Performance Monitoring Configuration
performance:
  enable_monitoring: true
  metrics_interval: 10s
  memory_threshold: 500MB
  cpu_threshold: 80
  enable_profiling: false
  profile_port: 6060

# Article Scanning Configuration
scanner:
  content_dir: "./content/post"
  include_patterns:
    - "*.md"
    - "*.markdown"
  exclude_patterns:
    - "_*"
    - "draft*"
  max_concurrent: 10

# Tag Page Generation Configuration
generator:
  tags_dir: "./content/tags"
  template_file: ""  # Custom template file
  overwrite_existing: false
  backup_existing: true

# Translation Configuration
translator:
  fallback_enabled: true
  custom_mappings:
    "人工智能": "artificial-intelligence"
    "机器学习": "machine-learning"
    "前端开发": "frontend-development"
  batch_size: 10
  request_delay: 500ms

# User Interface Configuration
ui:
  language: "en-US"  # zh-CN, en-US
  color_output: true
  progress_bar: true
  table_max_width: 120
  
# Security Configuration
security:
  max_file_size: 10MB
  allowed_extensions:
    - ".md"
    - ".markdown"
  sandbox_mode: false
```

## Detailed Configuration Items

### LM Studio Configuration (lm_studio)

| Item | Type | Default | Description |
|------|------|---------|-------------|
| url | string | "http://localhost:2234/v1/chat/completions" | LM Studio API address |
| model | string | "gemma-3-12b-it" | AI model name to use |
| timeout | duration | 30s | Request timeout |
| max_retries | int | 3 | Maximum retry attempts |
| retry_delay | duration | 1s | Retry interval |
| api_key | string | "" | API key (if required) |

### Cache Configuration (cache)

| Item | Type | Default | Description |
|------|------|---------|-------------|
| directory | string | "./cache" | Cache directory |
| file_name | string | "tag_translations_cache.json" | Cache file name |
| auto_save | bool | true | Auto-save cache |
| save_interval | duration | 5m | Auto-save interval |
| max_entries | int | 10000 | Maximum cache entries |
| expire_after | duration | 720h | Cache expiration time |

### Logging Configuration (logging)

| Item | Type | Default | Description |
|------|------|---------|-------------|
| level | string | "INFO" | Log level |
| file_path | string | "./logs/app.log" | Log file path |
| max_size | string | 100MB | Maximum single log file size |
| max_backups | int | 5 | Number of log files to retain |
| max_age | int | 30 | Log file retention days |
| console_output | bool | true | Output to console |
| json_format | bool | false | Use JSON format |

### Performance Monitoring Configuration (performance)

| Item | Type | Default | Description |
|------|------|---------|-------------|
| enable_monitoring | bool | true | Enable performance monitoring |
| metrics_interval | duration | 10s | Metrics collection interval |
| memory_threshold | string | 500MB | Memory usage warning threshold |
| cpu_threshold | int | 80 | CPU usage warning threshold (%) |
| enable_profiling | bool | false | Enable profiling |
| profile_port | int | 6060 | Profiling port |

## Environment Variable Override

All configuration items can be overridden with environment variables using the format: `HUGO_SLUG_<SECTION>_<KEY>`

### Common Environment Variables

```bash
# LM Studio Configuration
export HUGO_SLUG_LM_STUDIO_URL="http://192.168.1.100:2234/v1/chat/completions"
export HUGO_SLUG_LM_STUDIO_MODEL="gpt-4"
export HUGO_SLUG_LM_STUDIO_TIMEOUT="60s"

# Logging Configuration
export HUGO_SLUG_LOGGING_LEVEL="DEBUG"
export HUGO_SLUG_LOGGING_FILE_PATH="./logs/debug.log"

# Cache Configuration
export HUGO_SLUG_CACHE_DIRECTORY="./custom_cache"
export HUGO_SLUG_CACHE_MAX_ENTRIES="20000"

# Performance Configuration
export HUGO_SLUG_PERFORMANCE_ENABLE_MONITORING="false"
```

### Nested Configuration Override

For nested configuration items, use underscores as separators:

```bash
# Override translator.custom_mappings
export HUGO_SLUG_TRANSLATOR_CUSTOM_MAPPINGS_AI="artificial-intelligence"

# Override ui.language
export HUGO_SLUG_UI_LANGUAGE="en-US"
```

## Configuration Validation

The program automatically validates the configuration file on startup:

### Validation Rules
- Required field checks
- Data type validation
- Value range checks
- File path validity
- Network address format validation

### Validation Error Examples
```
Configuration validation failed:
- lm_studio.timeout: Must be a valid time format (e.g., 30s, 1m, 2h)
- logging.level: Must be one of DEBUG, INFO, WARN, ERROR
- cache.max_entries: Must be a positive integer
- performance.memory_threshold: Must be a valid byte size (e.g., 100MB, 1GB)
```

## Configuration Hot Reload

Supports reloading configuration files at runtime:

### Trigger Methods
1. **Signal Trigger** (Linux/macOS):
   ```bash
   kill -SIGHUP <process_id>
   ```

2. **File Monitoring**:
   The program automatically monitors configuration file changes and reloads

3. **API Interface**:
   ```bash
   curl -X POST http://localhost:8080/api/reload-config
   ```

### Hot Reload Limitations
- Some configuration items require program restart to take effect (like log file path)
- Related functions are paused during reload
- Original configuration is maintained if validation fails

## Configuration Template Generation

Generate default configuration file:

```bash
# Generate default configuration
go run main.go --generate-config

# Generate detailed configuration with comments
go run main.go --generate-config --with-comments

# Specify output file
go run main.go --generate-config --output=my-config.yaml
```

## Best Practices

### Development Environment
```yaml
logging:
  level: "DEBUG"
  console_output: true
  
performance:
  enable_monitoring: true
  enable_profiling: true
```

### Production Environment
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

### Performance Optimization
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

## Troubleshooting

### Common Configuration Issues

1. **Configuration file does not exist**
   ```
   Solution: Use --generate-config to generate default configuration
   ```

2. **YAML format error**
   ```
   Solution: Check indentation and syntax, use online YAML validators
   ```

3. **Environment variables not taking effect**
   ```
   Solution: Check variable name format, ensure program restart after changes
   ```

4. **Permission issues**
   ```
   Solution: Check read/write permissions for configuration files and directories
   ```

## Related Documentation

- [Installation Guide](installation_en.md)
- [Logging Guide](logging_en.md)
- [Performance Guide](performance_en.md)
- [Troubleshooting](troubleshooting_en.md)
