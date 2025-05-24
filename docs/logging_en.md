# Logging Guide

English | [中文](logging.md)

## Logging System Overview

This tool uses a hierarchical logging system that supports multiple output formats and automatic rotation features, helping users track program execution status and troubleshoot issues.

## Log Levels

### Level Description

| Level | Value | Purpose | Examples |
|-------|-------|---------|----------|
| DEBUG | 0 | Detailed debugging info | Function calls, variable values, execution flow |
| INFO | 1 | General information | Operation start/completion, config loading, statistics |
| WARN | 2 | Warning messages | Config issues, performance problems, compatibility issues |
| ERROR | 3 | Error messages | Operation failures, network errors, file I/O errors |

### Level Configuration

```yaml
# config.yaml
logging:
  level: "INFO"  # Only output INFO level and above
```

Environment variable override:
```bash
export HUGO_SLUG_LOGGING_LEVEL="DEBUG"
```

## Log Output Configuration

### File Output

```yaml
logging:
  file_path: "./logs/app.log"
  max_size: 100MB      # Maximum single file size
  max_backups: 5       # Number of historical files to retain
  max_age: 30          # File retention days
```

### Console Output

```yaml
logging:
  console_output: true  # Also output to console
  json_format: false    # Use plain text format for console
```

### Dual Output

The program supports simultaneous output to both file and console, convenient for development debugging and production monitoring.

## Log Formats

### Text Format (Default)

```
2024-01-15 14:30:25 [INFO] [scanner] Starting article directory scan: ./content/post
2024-01-15 14:30:25 [DEBUG] [scanner] Found file: article1.md
2024-01-15 14:30:25 [INFO] [translator] Starting batch translation of 5 tags
2024-01-15 14:30:26 [WARN] [cache] Cache file does not exist, will create new file
2024-01-15 14:30:27 [ERROR] [translator] LM Studio connection failed: connection refused
```

### JSON Format

```yaml
logging:
  json_format: true
```

JSON output example:
```json
{
  "timestamp": "2024-01-15T14:30:25.123Z",
  "level": "INFO",
  "module": "scanner",
  "message": "Starting article directory scan",
  "data": {
    "directory": "./content/post",
    "file_count": 0
  },
  "correlation_id": "req_123456"
}
```

## Log Modules

### Core Module Identifiers

| Module | Identifier | Function |
|--------|------------|----------|
| Scanner | scanner | Article scanning and parsing |
| Translator | translator | AI translation functionality |
| Generator | generator | Content generation |
| Cache | cache | Cache management |
| Performance | performance | Performance monitoring |
| Config | config | Configuration management |
| Main | main | Main program flow |

## Log File Management

### Automatic Rotation

When log files reach the specified size, they automatically rotate:

```
logs/
├── app.log           # Current log file
├── app.log.1         # Most recent backup
├── app.log.2         # Second most recent backup
├── app.log.3
├── app.log.4
└── app.log.5         # Oldest backup
```

### Automatic Cleanup

- Files exceeding `max_backups` count will be deleted
- Files older than `max_age` days will be deleted
- Cleanup operations execute during each rotation

### Compressed Storage

Historical log files can be automatically compressed (requires configuration):

```yaml
logging:
  compress_backups: true  # Compress historical logs
```

Compressed file format:
```
logs/
├── app.log
├── app.log.1.gz
├── app.log.2.gz
└── ...
```

## Log Viewing and Analysis

### Real-time Viewing

```bash
# View real-time logs
tail -f logs/app.log

# View last 100 lines
tail -n 100 logs/app.log

# Real-time viewing with error highlighting
tail -f logs/app.log | grep --color=always -E "(ERROR|WARN|$)"
```

### Filter by Level

```bash
# View error logs
grep "ERROR" logs/app.log

# View warnings and errors
grep -E "(WARN|ERROR)" logs/app.log

# View specific module logs
grep "[translator]" logs/app.log
```

### Filter by Time

```bash
# View today's logs
grep "$(date +%Y-%m-%d)" logs/app.log

# View specific time period
grep "2024-01-15 14:" logs/app.log

# View errors from last hour
grep "ERROR" logs/app.log | grep "$(date -d '1 hour ago' +%Y-%m-%d\ %H):"
```

### Statistical Analysis

```bash
# Count logs by level
grep -c "ERROR" logs/app.log
grep -c "WARN" logs/app.log
grep -c "INFO" logs/app.log

# Count errors by hour
grep "ERROR" logs/app.log | cut -d' ' -f2 | cut -d':' -f1 | sort | uniq -c

# Find most frequent errors
grep "ERROR" logs/app.log | cut -d']' -f3- | sort | uniq -c | sort -nr | head -10
```

## Log Monitoring and Alerting

### Log Monitoring Script

```bash
#!/bin/bash
# monitor_logs.sh

LOG_FILE="./logs/app.log"
ERROR_THRESHOLD=10

# Count errors in last 5 minutes
ERROR_COUNT=$(tail -n 1000 "$LOG_FILE" | grep "$(date -d '5 minutes ago' +%Y-%m-%d\ %H:%M)" | grep -c "ERROR")

if [ "$ERROR_COUNT" -gt "$ERROR_THRESHOLD" ]; then
    echo "Warning: $ERROR_COUNT errors in last 5 minutes, exceeding threshold $ERROR_THRESHOLD"
    # Send alert notification
fi
```

### Monitoring System Integration

#### Prometheus + Grafana

```yaml
# prometheus.yml
- job_name: 'hugo-slug-auto'
  static_configs:
    - targets: ['localhost:8080']
  metrics_path: '/metrics'
```

#### ELK Stack

```yaml
# filebeat.yml
filebeat.inputs:
- type: log
  enabled: true
  paths:
    - ./logs/*.log
  fields:
    service: hugo-slug-auto
```

## Debugging Techniques

### Enable Verbose Logging

```yaml
logging:
  level: "DEBUG"
  console_output: true
```

### Track Specific Operations

Add correlation IDs in code:
```go
correlationID := uuid.New().String()
logger.WithField("correlation_id", correlationID).Info("Starting request processing")
```

### Performance Debugging

Enable performance logging:
```yaml
logging:
  performance_logging: true
```

Performance log example:
```
2024-01-15 14:30:25 [PERF] [translator] Translation time: 1.23s, tag count: 10
2024-01-15 14:30:26 [PERF] [cache] Cache hit rate: 85.2%
```

## Best Practices

### Development Environment

```yaml
logging:
  level: "DEBUG"
  console_output: true
  file_path: "./logs/dev.log"
  json_format: false
  max_size: 10MB
  max_backups: 3
```

### Production Environment

```yaml
logging:
  level: "INFO"
  console_output: false
  file_path: "/var/log/hugo-slug-auto/app.log"
  json_format: true
  max_size: 100MB
  max_backups: 10
  max_age: 30
  compress_backups: true
```

### Performance Considerations

1. **Avoid Excessive Logging**
   - Use INFO level in production
   - Avoid DEBUG logs in loops

2. **Asynchronous Log Writing**
   ```yaml
   logging:
     async_write: true
     buffer_size: 1000
   ```

3. **Log Rotation Strategy**
   - Set reasonable file sizes based on disk space
   - Balance historical retention and storage costs

## Troubleshooting

### Common Issues

1. **Cannot create log file**
   ```
   Cause: Insufficient directory permissions
   Solution: Check and modify directory permissions
   ```

2. **Log rotation failure**
   ```
   Cause: Insufficient disk space
   Solution: Free up disk space or adjust log policy
   ```

3. **JSON format error**
   ```
   Cause: Log content contains special characters
   Solution: Enable character escaping or use text format
   ```

4. **Performance impact**
   ```
   Cause: Too many DEBUG level logs
   Solution: Adjust log level or enable asynchronous writing
   ```

## Related Documentation

- [Configuration Guide](configuration_en.md)
- [Performance Guide](performance_en.md)
- [Troubleshooting](troubleshooting_en.md)
- [API Documentation](api_en.md)
