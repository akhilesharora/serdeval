# SerdeVal

![SerdeVal Banner](https://img.shields.io/badge/SerdeVal-Privacy%20First-brightgreen)
[![Go Report Card](https://goreportcard.com/badge/github.com/akhilesharora/serdeval)](https://goreportcard.com/report/github.com/akhilesharora/serdeval)
[![Go Reference](https://pkg.go.dev/badge/github.com/akhilesharora/serdeval.svg)](https://pkg.go.dev/github.com/akhilesharora/serdeval)
[![CI](https://github.com/akhilesharora/serdeval/actions/workflows/ci.yml/badge.svg)](https://github.com/akhilesharora/serdeval/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/akhilesharora/serdeval/graph/badge.svg?token=KWAMELIEPS)](https://codecov.io/gh/akhilesharora/serdeval)
[![Go Version](https://img.shields.io/badge/go%20version-%3E=1.22-61CFDD.svg?style=flat-square)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A privacy-focused, blazingly fast data format validator supporting 14+ formats including JSON, YAML, XML, TOML, CSV, GraphQL, Markdown, and more.

**Privacy-focused**: All validation happens locally on your machine.

## 🌐 Live Demo

**[freedatavalidator.xyz](https://freedatavalidator.xyz)** - Try it online (client-side validation)

## 🚀 Features

- **🔒 Privacy-First**: No logging, tracking, or data retention
- **⚡ Blazingly Fast**: Zero-dependency Go implementation
- **🎯 Auto-Detection**: Automatically detects data formats
- **📱 Multiple Interfaces**: CLI, Go library, and web interface
- **🌐 Cross-Platform**: Windows, macOS, and Linux support
- **🧠 Smart Formatting**: Beautifies and validates in one step
- **📊 Developer-Friendly**: JSON output for CI/CD pipelines

### Supported Formats

| Format | Extensions | Auto-Detection | Validation | Use Case |
|--------|------------|----------------|------------|----------|
| JSON   | `.json`    | ✅             | ✅         | APIs, Config files |
| YAML   | `.yaml`, `.yml` | ✅       | ✅         | Kubernetes, CI/CD |
| XML    | `.xml`     | ✅             | ✅         | Enterprise, SOAP |
| TOML   | `.toml`    | ✅             | ✅         | Config files |
| CSV    | `.csv`     | ✅             | ✅         | Data exchange |
| GraphQL| `.graphql`, `.gql` | ✅    | ✅         | API schemas |
| INI    | `.ini`, `.cfg`, `.conf` | ✅ | ✅      | Config files |
| HCL    | `.hcl`, `.tf`, `.tfvars` | ✅ | ✅    | Terraform |
| Protobuf| `.proto`, `.textproto` | ✅ | ✅      | Protocol Buffers |
| Markdown| `.md`, `.markdown` | ✅   | ✅         | Documentation |
| JSON Lines| `.jsonl`, `.ndjson` | ✅ | ✅       | Streaming data |
| Jupyter | `.ipynb`  | ✅             | ✅         | Data science |
| Requirements.txt | `.txt` | ✅     | ✅         | Python deps |
| Dockerfile | `Dockerfile*` | ✅     | ✅         | Containers |

## 📦 Installation

### Using Go

```bash
# As a CLI tool
go install github.com/akhilesharora/serdeval/cmd/serdeval@latest

# As a library
go get github.com/akhilesharora/serdeval/pkg/validator
```

### Pre-built Binaries

Download the latest binary for your platform from the [releases page](https://github.com/akhilesharora/serdeval/releases).

### From Source

#### Linux/macOS
```bash
git clone https://github.com/akhilesharora/serdeval
cd serdeval
make build
```

#### Windows
```powershell
git clone https://github.com/akhilesharora/serdeval
cd serdeval
go build -o serdeval.exe ./cmd/serdeval
```

### Development Setup

For contributors, set up pre-commit hooks to ensure code quality:

```bash
# Install pre-commit hooks
make pre-commit

# Or manually
./scripts/setup-hooks.sh
```

This will install hooks that:
- Format code automatically  
- Run tests before push
- Check for linting issues
- Prevent commits with formatting errors

### Deployment

For server deployment, copy `deploy.example.sh` to `deploy.sh` and customize it with your specific configuration:

```bash
cp deploy.example.sh deploy.sh
# Edit deploy.sh with your server details
```

## 🖥️ Usage

### Command Line Interface

```bash
# Validate a single file
serdeval validate config.json

# Validate multiple files
serdeval validate config.json data.yaml settings.toml

# Validate from stdin
echo '{"name": "John", "age": 30}' | serdeval validate

# Specify format explicitly
serdeval validate --format json config.txt

# Output as JSON for CI/CD pipelines
serdeval validate --json config.json

# Start web interface
serdeval web --port 8080

# Check version
serdeval version

# Windows examples
serdeval.exe validate config.json
type config.json | serdeval.exe validate
```

### Go Library

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/akhilesharora/serdeval/pkg/validator"
)

func main() {
    // Auto-detect format
    data := []byte(`{"name": "John", "age": 30}`)
    result := validator.ValidateAuto(data)
    
    if result.Valid {
        fmt.Printf("Valid %s data\n", result.Format)
    } else {
        log.Fatalf("Invalid data: %s", result.Error)
    }
    
    // Validate specific format
    v, _ := validator.NewValidator(validator.FormatJSON)
    result = v.ValidateString(`{"test": true}`)
    fmt.Printf("Valid: %v\n", result.Valid)
}
```

### Web Interface

Start the built-in web server:

```bash
serdeval web --port 8080
```

Then visit http://localhost:8080 for a user-friendly interface with:
- Real-time validation as you type
- Auto-format with beautification
- Copy-to-clipboard functionality
- Format auto-detection
- **100% client-side processing** (your data never leaves your browser)

## 🛠️ Development

### Prerequisites

- Go 1.22 or higher
- Make (optional, for convenience commands)

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Run linters
make lint

# Format code
make fmt

# Run benchmarks
make bench
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...

# Run benchmarks
go test -bench=. ./...
```

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](.github/CONTRIBUTING.md) for details.

### Quick Start for Contributors

1. Fork the repository
2. Set up development environment: `make pre-commit`
3. Make your changes
4. Run tests: `make test`
5. Submit a pull request

Please read our [Code of Conduct](.github/CODE_OF_CONDUCT.md).

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with ❤️ for the developer community
- Inspired by the need for privacy-focused validation tools
- Thanks to all contributors who help improve this project


---

**[Visit SerdeVal](https://freedatavalidator.xyz)** | **[Documentation](https://pkg.go.dev/github.com/akhilesharora/serdeval)** | **[Report Issues](https://github.com/akhilesharora/serdeval/issues)**