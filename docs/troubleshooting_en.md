# Troubleshooting Guide

[中文](troubleshooting.md) | English

## Common Issues

### 1. LM Studio Connection Issues

#### Issue Symptoms
```
Warning: Unable to connect to LM Studio (dial tcp 172.19.192.1:2234: connectex: No connection could be made...)
```

#### Solutions
1. **Check LM Studio Status**
   ```bash
   # Ensure LM Studio is running and has loaded a model
   ```

2. **Check Network Configuration**
   ```go
   // Modify URL in translator/llm_translator.go
   const LMStudioURL = "http://localhost:2234/v1/chat/completions"
   ```

3. **Test Connection**
   ```bash
   curl -X POST http://localhost:2234/v1/chat/completions \
     -H "Content-Type: application/json" \
     -d '{"model":"your-model","messages":[{"role":"user","content":"test"}]}'
   ```

### 2. File Permission Issues

#### Issue Symptoms
```
Failed to create tags directory: permission denied
Failed to write cache file: permission denied
```

#### Solutions
1. **Check Directory Permissions**
   ```bash
   # Linux/macOS
   chmod 755 /path/to/content
   chmod 644 /path/to/content/**/*.md
   
   # Windows
   # Right-click -> Properties -> Security -> Edit permissions
   ```

2. **Run with Administrator Privileges**
   ```bash
   # Linux/macOS
   sudo go run main.go
   
   # Windows
   # Run Command Prompt as Administrator
   ```

### 3. Go Module Issues

#### Issue Symptoms
```
go: module hugo-content-suite: cannot find module providing package
```

#### Solutions
1. **Reinitialize Module**
   ```bash
   go mod init hugo-content-suite
   go mod tidy
   ```

2. **Clean Module Cache**
   ```bash
   go clean -modcache
   go mod download
   ```

3. **Check Go Version**
   ```bash
   go version
   # Ensure version >= 1.21
   ```

### 4. Article Parsing Issues

#### Issue Symptoms
- Tags not being recognized correctly
- Front Matter parsing failure

#### Solutions
1. **Check File Format**
   ```yaml
   ---
   title: "Article Title"
   tags: ["tag1", "tag2"]  # Ensure array format
   categories: ["category"]
   date: 2024-01-01
   ---
   ```

2. **Supported Tag Formats**
   ```yaml
   # Array format (recommended)
   tags: ["tag1", "tag2"]
   
   # YAML list format
   tags:
     - tag1
     - tag2
   ```

### 5. Cache Related Issues

#### Corrupted Cache File
```bash
# Delete corrupted cache file
rm tag_translations_cache.json

# Or use built-in clear function
# Menu selection: 9. Clear Translation Cache
```

#### Cache Path Issues
```go
// Modify cache path in translator/llm_translator.go
cache: NewTranslationCache("/custom/cache/path"),
```

### 6. Translation Quality Issues

#### Poor AI Translation Results
1. **Change Model**
   ```go
   const ModelName = "better-model-name"
   ```

2. **Adjust Prompt**
   ```go
   // Modify prompt in TranslateToSlug method
   prompt := fmt.Sprintf(`More detailed translation instructions...`)
   ```

3. **Add Predefined Mappings**
   ```go
   // Add to fallbackSlug function
   fallbackTranslations := map[string]string{
       "your-tag": "your-tag",
       // Add more mappings
   }
   ```

## Debugging Techniques

### 1. Enable Verbose Logging

```go
// Add debugging info where needed
fmt.Printf("Debug info: %+v\n", variable)
```

### 2. Check File Contents

```bash
# View generated tag pages
find content/tags -name "_index.md" -exec head -10 {} \;

# Check cache file
cat tag_translations_cache.json | jq .
```

### 3. Test Individual Functions

```go
// Create simple test file
func main() {
    translator := translator.NewLLMTranslator()
    result, err := translator.TranslateToSlug("test tag")
    fmt.Printf("Result: %s, Error: %v\n", result, err)
}
```

## Performance Optimization

### 1. Large Article Processing

```bash
# Process large numbers of articles in batches
# First process part of directory for testing
go run main.go ./content/post/2024

# Process all after confirming no issues
go run main.go ./content/post
```

### 2. Network Timeout Adjustment

```go
// Increase timeout duration
client: &http.Client{
    Timeout: 60 * time.Second,  // Adjust to 60 seconds
},
```

### 3. Concurrent Processing

```go
// Current program is serial processing, be careful with cache synchronization if implementing concurrency
```

## Environment-Specific Issues

### Windows Environment

1. **Path Separator Issues**
   ```go
   // Use filepath.Join instead of manual concatenation
   path := filepath.Join("content", "post")
   ```

2. **Character Encoding Issues**
   ```bash
   # Ensure terminal supports UTF-8
   chcp 65001
   ```

### macOS Environment

1. **Homebrew Go Version**
   ```bash
   brew upgrade go
   go version
   ```

### Linux Environment

1. **Package Installation**
   ```bash
   # Ubuntu/Debian
   sudo apt update
   sudo apt install golang-go

   # CentOS/RHEL
   sudo yum install golang
   ```

## Getting Help

### 1. Check Logs
The program outputs detailed operation logs; pay attention to error messages.

### 2. Verify Configuration
Confirm all configuration items (URL, paths, model names) are correct.

### 3. Community Support
- GitHub Issues: Submit problems and suggestions
- Documentation Updates: Report issues promptly

### 4. Contact Information
If the above methods cannot solve the problem, please:
1. Prepare detailed error information
2. Specify operating system and Go version
3. Provide reproduction steps
4. Submit issue through GitHub Issues

## Preventive Measures

### 1. Backup Important Files
```bash
# Backup content directory before processing
cp -r content content_backup
```

### 2. Small-Scale Testing
```bash
# Verify functionality in test directory first
mkdir test_content
cp content/post/sample.md test_content/
go run main.go test_content
```

### 3. Version Control
```bash
# Use Git to track changes
git add .
git commit -m "Backup before processing"
```

This allows quick rollback to previous state if issues occur.
