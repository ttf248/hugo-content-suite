# Performance Monitoring Guide

English | [中文](performance.md)

## Performance Monitoring Overview

Built-in performance monitoring system tracks program runtime status and operational efficiency.

## Key Metrics

### System Resources
- Memory usage and CPU utilization
- File I/O operations
- Network request latency

### Business Performance
- Translation speed and cache hit rate
- Processing throughput
- Error count and success rate

## Configuration

### Basic Setup

```yaml
# config.yaml
performance:
  enable_monitoring: true
  metrics_interval: 10s
  memory_threshold: 500MB
```

## Performance Analysis

### Real-time Monitoring
View performance statistics through the menu system (option 8).

### Performance Reports
```
📊 Performance Statistics:
🔄 Translation Count: 156
⚡ Cache Hit Rate: 87.5%
⏱️ Average Translation Time: 1.2s
📁 File Operations: 89
```

## Optimization Tips

### Memory Optimization
- Reduce cache size if memory usage is high
- Use batch processing for large datasets

### Speed Optimization
- Ensure good network connection to LM Studio
- Utilize translation cache effectively
- Process in smaller batches if needed

## Troubleshooting

### Common Issues
1. **High Memory Usage**: Reduce cache entries or batch size
2. **Slow Processing**: Check network connection and API response times
3. **Cache Misses**: Generate bulk cache before processing

## Related Documentation
- [Configuration Guide](configuration_en.md)
- [Troubleshooting](troubleshooting_en.md)
