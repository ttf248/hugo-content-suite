# Performance Monitoring Guide

English | [中文](performance.md)

## Performance Monitoring Overview

This tool has a built-in comprehensive performance monitoring system that can track program runtime status, resource usage, and operational efficiency in real-time, helping users optimize processing performance.

## Monitoring Metrics

### System Resource Metrics

| Metric Type | Metric Name | Unit | Description |
|-------------|-------------|------|-------------|
| Memory | Used Memory | MB | Current memory usage by program |
| Memory | Memory Usage Rate | % | Percentage of total system memory |
| CPU | CPU Usage Rate | % | Program CPU usage |
| Disk | Disk Read Speed | MB/s | File reading speed |
| Disk | Disk Write Speed | MB/s | File writing speed |

### Business Performance Metrics

| Metric Type | Metric Name | Unit | Description |
|-------------|-------------|------|-------------|
| Processing Speed | Article Scan Speed | articles/s | Number of articles scanned per second |
| Processing Speed | Translation Speed | tags/s | Number of tags translated per second |
| Network | API Response Time | ms | LM Studio API average response time |
| Network | API Success Rate | % | API call success rate |
| Cache | Cache Hit Rate | % | Cache hit percentage |
| Cache | Cache Size | KB | Cache file size |

## Performance Monitoring Configuration

### Basic Configuration

```yaml
# config.yaml
performance:
  enable_monitoring: true        # Enable performance monitoring
  metrics_interval: 10s         # Metrics collection interval
  memory_threshold: 500MB       # Memory usage warning threshold
  cpu_threshold: 80             # CPU usage warning threshold (%)
  enable_profiling: false       # Enable Go profiling
  profile_port: 6060           # Profiling port
  metrics_retention: 24h       # Metrics retention time
  export_metrics: true         # Export Prometheus metrics
```

### Advanced Configuration

```yaml
performance:
  # Detailed monitoring configuration
  detailed_monitoring:
    goroutine_count: true      # Monitor goroutine count
    gc_stats: true            # Monitor garbage collection statistics
    memory_breakdown: true    # Detailed memory analysis
    
  # Alert configuration
  alerts:
    memory_critical: 1GB      # Memory critical alert threshold
    response_time_critical: 5s # Response time critical alert
    error_rate_critical: 10   # Error rate critical alert (%)
    
  # Profiling configuration
  profiling:
    cpu_profile_duration: 30s # CPU profiling duration
    memory_profile_interval: 5m # Memory profiling interval
    block_profile: true       # Block profiling
    mutex_profile: true       # Mutex profiling
```

## Real-time Performance Monitoring

### Console Display

The program displays real-time performance information at the bottom of the interface during runtime:

```
=====================================
📊 Real-time Performance Monitoring
=====================================
💾 Memory Usage: 245.2MB / 500MB (49.0%)
⚡ CPU Usage: 15.3%
🏃 Processing Speed: 12.5 articles/s
🌐 API Response: 1.2s (average)
💾 Cache Hit Rate: 87.5%
📊 Goroutines: 8
⏱️  Runtime: 2m 15s
=====================================
```

### Performance Dashboard

Enable web dashboard to view detailed performance metrics:

```bash
go run main.go --dashboard --port=8080
```

Visit `http://localhost:8080/dashboard` to view:
- Real-time performance charts
- Historical trend analysis
- Resource usage details
- Operation statistics reports

## Performance Analysis Tools

### Go pprof Integration

Enable built-in Go performance analysis tools:

```yaml
performance:
  enable_profiling: true
  profile_port: 6060
```

#### CPU Performance Analysis

```bash
# Collect 30 seconds of CPU performance data
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Analyze CPU hotspots
(pprof) top10
(pprof) list function_name
(pprof) web
```

#### Memory Performance Analysis

```bash
# Memory usage analysis
go tool pprof http://localhost:6060/debug/pprof/heap

# Memory allocation analysis
go tool pprof http://localhost:6060/debug/pprof/allocs
```

#### Goroutine Analysis

```bash
# Goroutine status analysis
go tool pprof http://localhost:6060/debug/pprof/goroutine

# Block analysis
go tool pprof http://localhost:6060/debug/pprof/block
```

### Performance Report Generation

#### Automatic Reports

Automatically generate performance reports after program execution:

```
=====================================
📊 Performance Statistics Report
=====================================
⏱️  Total Runtime: 5m 32s
📝 Total Articles Processed: 1,247
🏷️  Total Tags Translated: 156
📊 Average Processing Speed: 3.7 articles/s
🌐 API Calls: 89 times
💾 Cache Hit Rate: 82.3%
📈 Peak Memory Usage: 287.5MB
⚡ Average CPU Usage: 23.1%
🔄 Garbage Collections: 12 times
=====================================
```

#### Detailed Report Export

```bash
# Export detailed performance report
go run main.go --export-performance-report=./reports/performance.json

# Export Prometheus format metrics
go run main.go --export-metrics=./metrics/metrics.txt
```

## Performance Optimization Recommendations

### Memory Optimization

#### Monitor Memory Leaks

```bash
# Periodically check memory usage trends
watch -n 5 'curl -s http://localhost:6060/debug/pprof/heap | head -20'
```

#### Optimization Configuration

```yaml
# Configuration to reduce memory usage
cache:
  max_entries: 5000      # Reduce cache entry count
  expire_after: 24h      # Shorter cache expiration time

scanner:
  max_concurrent: 5      # Reduce concurrent processing count

translator:
  batch_size: 5          # Reduce batch size
```

### CPU Optimization

#### Concurrent Processing Optimization

```yaml
scanner:
  max_concurrent: 16     # Adjust based on CPU core count

performance:
  enable_profiling: true # Enable profiling to find hotspots
```

#### I/O Optimization

```yaml
cache:
  auto_save: false       # Manual control of cache save timing
  save_interval: 10m     # Increase save interval
```

### Network Optimization

#### API Call Optimization

```yaml
lm_studio:
  timeout: 10s           # Reduce timeout
  max_retries: 1         # Reduce retry count

translator:
  batch_size: 20         # Increase batch size
  request_delay: 200ms   # Reduce request interval
```

## Performance Monitoring Alerts

### Threshold Alerts

```yaml
performance:
  alerts:
    memory_warning: 400MB    # Memory warning threshold
    memory_critical: 800MB   # Memory critical threshold
    cpu_warning: 70         # CPU warning threshold
    cpu_critical: 90        # CPU critical threshold
    response_time_warning: 3s # Response time warning
    cache_hit_rate_low: 60  # Low cache hit rate threshold
```

### Alert Notifications

#### Email Notifications

```yaml
alerts:
  email:
    enabled: true
    smtp_server: "smtp.gmail.com:587"
    username: "your-email@gmail.com"
    password: "your-password"
    recipients:
      - "admin@example.com"
```

#### Webhook Notifications

```yaml
alerts:
  webhook:
    enabled: true
    url: "https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK"
    timeout: 5s
```

## Performance Benchmarking

### Benchmark Script

```bash
#!/bin/bash
# benchmark.sh

echo "Starting performance benchmark..."

# Prepare test data
mkdir -p test_data
for i in {1..1000}; do
    echo "---
title: \"Test Article $i\"
tags: [\"tag1\", \"tag2\", \"tag3\"]
date: 2024-01-01
---
Test content" > "test_data/article_$i.md"
done

# Run benchmark
echo "Test start time: $(date)"
time go run main.go test_data --benchmark
echo "Test end time: $(date)"

# Cleanup test data
rm -rf test_data
```

### Benchmark Results

```
=====================================
📊 Benchmark Results
=====================================
📝 Test Article Count: 1,000
🏷️  Test Tag Count: 50
⏱️  Total Processing Time: 2m 15s
📊 Processing Speed: 7.4 articles/s
🌐 Translation Speed: 3.2 tags/s
💾 Peak Memory: 156.7MB
⚡ Peak CPU: 45.2%
💾 Cache Hit Rate: 94.1%
=====================================
```

## Troubleshooting

### Performance Issue Diagnosis

#### High Memory Usage

```bash
# Check memory allocation hotspots
go tool pprof http://localhost:6060/debug/pprof/heap
(pprof) top10
(pprof) list function_name
```

#### High CPU Usage

```bash
# Check CPU hotspots
go tool pprof http://localhost:6060/debug/pprof/profile
(pprof) top10
(pprof) web
```

#### Slow Response Diagnosis

```bash
# Check blocking points
go tool pprof http://localhost:6060/debug/pprof/block
```

### Common Performance Issues

1. **Continuous Memory Growth**
   - Check for memory leaks
   - Reduce cache size
   - Enable periodic garbage collection

2. **Slow Processing Speed**
   - Increase concurrent processing count
   - Optimize file I/O
   - Use SSD storage

3. **Slow API Response**
   - Check network connection
   - Adjust timeout settings
   - Use local models

## Related Documentation

- [Configuration Guide](configuration_en.md)
- [Logging Guide](logging_en.md)
- [Troubleshooting](troubleshooting_en.md)
- [API Documentation](api_en.md)
