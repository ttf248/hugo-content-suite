# Contributing Guide

[中文](contributing.md) | English

Welcome to contribute to the Hugo Tag Auto Management Tool!

## How to Contribute

### 🐛 Report Bugs
1. Search [GitHub Issues](https://github.com/ttf248/hugo-slug-auto/issues) for existing similar issues
2. If none exists, create a new Issue with:
   - Detailed problem description
   - Steps to reproduce
   - Expected behavior
   - Actual behavior
   - Environment information (OS, Go version, etc.)

### 💡 Feature Suggestions
1. Describe new feature requirements in detail in Issues
2. Explain why this feature is needed
3. Provide possible implementation approaches

### 🔧 Code Contributions

#### Development Environment Setup
```bash
# 1. Fork the project to your GitHub account
# 2. Clone your fork
git clone https://github.com/your-username/hugo-slug-auto.git
cd hugo-slug-auto

# 3. Add upstream repository
git remote add upstream https://github.com/ttf248/hugo-slug-auto.git

# 4. Install dependencies
go mod tidy

# 5. Run tests
go test ./...
```

#### Development Workflow
1. **Create Feature Branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Write Code**
   - Follow existing code style
   - Add necessary comments
   - Write unit tests

3. **Test Code**
   ```bash
   go test ./...
   go run main.go  # Manual testing
   ```

4. **Commit Changes**
   ```bash
   git add .
   git commit -m "feat: add new feature description"
   ```

5. **Sync Upstream**
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

6. **Push and Create PR**
   ```bash
   git push origin feature/your-feature-name
   ```

## Code Standards

### Go Code Style
- Use `gofmt` to format code
- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use meaningful variable and function names
- Add necessary error handling

### Commit Message Format
Use [Conventional Commits](https://www.conventionalcommits.org/) format:

```
<type>[optional scope]: <description>

[optional body]

[optional footer]
```

Types include:
- `feat`: New features
- `fix`: Bug fixes
- `docs`: Documentation updates
- `style`: Code formatting changes
- `refactor`: Code refactoring
- `test`: Adding tests
- `chore`: Build process or auxiliary tool changes

Example:
```
feat(translator): add new translation engine support

- Support OpenAI GPT API
- Add configuration options
- Update documentation

Closes #123
```

### Code Comments
```go
// Package translator provides AI translation functionality
package translator

// TranslateToSlug translates Chinese tags to English slugs
// Parameter tag: Chinese tag to translate
// Returns: English slug and possible error
func (t *LLMTranslator) TranslateToSlug(tag string) (string, error) {
    // Implementation logic...
}
```

## Project Structure

### Directory Description
```
hugo-slug-auto/
├── main.go              # Program entry point
├── models/              # Data model definitions
├── scanner/             # File scanning and parsing
├── stats/               # Statistical analysis functions
├── translator/          # AI translation module
├── generator/           # Content generators
├── display/             # User interface display
├── docs/               # Project documentation
└── tests/              # Test files
```

### Module Design Principles
- **Single Responsibility**: Each module focuses on specific functionality
- **Loose Coupling**: Modules interact through interfaces
- **High Cohesion**: Related functions organized in the same module

## Testing Guide

### Unit Tests
```go
// translator_test.go
func TestTranslateToSlug(t *testing.T) {
    translator := NewLLMTranslator()
    
    tests := []struct {
        input    string
        expected string
    }{
        {"人工智能", "artificial-intelligence"},
        {"机器学习", "machine-learning"},
    }
    
    for _, test := range tests {
        result, err := translator.TranslateToSlug(test.input)
        assert.NoError(t, err)
        assert.Equal(t, test.expected, result)
    }
}
```

### Integration Tests
```bash
# Create test data
mkdir test-content
echo '---\ntitle: "Test Article"\ntags: ["test"]\n---\nTest content' > test-content/test.md

# Run program test
go run main.go test-content
```

## New Feature Development

### Adding New Translators
1. Create new file in `translator/` directory
2. Implement `Translator` interface
3. Add configuration options
4. Write tests
5. Update documentation

### Adding New Generators
1. Create new file in `generator/` directory
2. Implement generation logic
3. Integrate into main menu
4. Add preview functionality
5. Write tests

### Adding New Display Features
1. Add new functions in `display/` directory
2. Use unified colors and formatting
3. Support table and list display
4. Consider pagination for large datasets

## Documentation Contributions

### Documentation Types
- **README**: Quick start guide
- **Installation Docs**: Detailed installation and configuration instructions
- **Usage Docs**: Feature usage instructions
- **API Docs**: Interface documentation
- **Troubleshooting**: Common problem solutions

### Documentation Writing Standards
- Use clear heading structure
- Provide specific code examples
- Include necessary screenshots or diagrams
- Keep content up-to-date

## PR Review Process

### Pre-submission Checklist
- [ ] Code passes all tests
- [ ] Code conforms to project style
- [ ] Added necessary documentation
- [ ] Updated CHANGELOG (if applicable)

### Review Criteria
1. **Functional Correctness**: Does implementation meet requirements
2. **Code Quality**: Readability and maintainability
3. **Performance Impact**: Does it affect existing functionality performance
4. **Backward Compatibility**: Does it break existing APIs
5. **Security**: Does it introduce security risks

### Review Process
1. Automated checks (CI/CD)
2. Code review
3. Functional testing
4. Documentation review
5. Merge to main branch

## Release Process

### Version Numbering Rules
Follow [Semantic Versioning](https://semver.org/):
- `MAJOR.MINOR.PATCH`
- MAJOR: Incompatible API changes
- MINOR: Backward-compatible feature additions
- PATCH: Backward-compatible bug fixes

### Release Steps
1. Update version number
2. Update CHANGELOG
3. Create Release Tag
4. Write Release Notes
5. Publish to GitHub Releases

## Community Communication

### Discussion Platforms
- GitHub Issues: Bug reports and feature discussions
- GitHub Discussions: General discussions and Q&A

### Code of Conduct
- Stay friendly and respectful
- Focus on technical discussions
- Welcome newcomer questions
- Share knowledge and experience

## Acknowledgments

Thanks to all developers who contribute to the project! Your contributions make this project better.

### Contributors List
- Maintained in [CONTRIBUTORS.md](CONTRIBUTORS.md) file

---

Thank you again for your attention and contribution to the project! 🎉
