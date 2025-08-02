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

#### Basic Usage

```bash
# Validate a single file (auto-detects format)
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
```

#### Validate Each Format

```bash
# JSON files
serdeval validate package.json
serdeval validate data.json

# YAML files
serdeval validate docker-compose.yaml
serdeval validate config.yml

# XML files
serdeval validate pom.xml
serdeval validate web.xml

# TOML files
serdeval validate Cargo.toml
serdeval validate pyproject.toml

# CSV files
serdeval validate data.csv
serdeval validate report.csv

# GraphQL files
serdeval validate schema.graphql
serdeval validate query.gql

# INI/Config files
serdeval validate config.ini
serdeval validate settings.cfg
serdeval validate app.conf

# HCL/Terraform files
serdeval validate main.tf
serdeval validate variables.tfvars
serdeval validate config.hcl

# Protobuf text format
serdeval validate message.textproto
serdeval validate data.pbtxt

# Markdown files
serdeval validate README.md
serdeval validate docs.markdown

# JSON Lines files
serdeval validate logs.jsonl
serdeval validate events.ndjson

# Jupyter notebooks
serdeval validate analysis.ipynb

# Python requirements
serdeval validate requirements.txt
serdeval validate requirements-dev.txt

# Dockerfiles
serdeval validate Dockerfile
serdeval validate Dockerfile.prod
```

### Go Library

#### Basic Usage

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
}
```

#### Examples for Each Format

```go
// JSON Validation
jsonValidator, _ := validator.NewValidator(validator.FormatJSON)
result := jsonValidator.ValidateString(`{"name": "test", "value": 123}`)
fmt.Printf("JSON valid: %v\n", result.Valid)

// YAML Validation
yamlValidator, _ := validator.NewValidator(validator.FormatYAML)
result = yamlValidator.ValidateString(`
name: test
value: 123
items:
  - one
  - two
`)
fmt.Printf("YAML valid: %v\n", result.Valid)

// XML Validation
xmlValidator, _ := validator.NewValidator(validator.FormatXML)
result = xmlValidator.ValidateString(`<?xml version="1.0"?>
<root>
  <name>test</name>
  <value>123</value>
</root>`)
fmt.Printf("XML valid: %v\n", result.Valid)

// TOML Validation
tomlValidator, _ := validator.NewValidator(validator.FormatTOML)
result = tomlValidator.ValidateString(`
[server]
host = "localhost"
port = 8080
`)
fmt.Printf("TOML valid: %v\n", result.Valid)

// CSV Validation
csvValidator, _ := validator.NewValidator(validator.FormatCSV)
result = csvValidator.ValidateString(`name,age,city
John,30,NYC
Jane,25,LA`)
fmt.Printf("CSV valid: %v\n", result.Valid)

// GraphQL Validation
graphqlValidator, _ := validator.NewValidator(validator.FormatGraphQL)
result = graphqlValidator.ValidateString(`
query GetUser {
  user(id: "123") {
    name
    email
  }
}`)
fmt.Printf("GraphQL valid: %v\n", result.Valid)

// INI Validation
iniValidator, _ := validator.NewValidator(validator.FormatINI)
result = iniValidator.ValidateString(`
[database]
host = localhost
port = 5432
`)
fmt.Printf("INI valid: %v\n", result.Valid)

// HCL Validation
hclValidator, _ := validator.NewValidator(validator.FormatHCL)
result = hclValidator.ValidateString(`
resource "aws_instance" "example" {
  ami           = "ami-12345"
  instance_type = "t2.micro"
}`)
fmt.Printf("HCL valid: %v\n", result.Valid)

// Protobuf Text Validation
protoValidator, _ := validator.NewValidator(validator.FormatProtobuf)
result = protoValidator.ValidateString(`
type_url: "type.googleapis.com/example"
value: "\n\x05hello"
`)
fmt.Printf("Protobuf valid: %v\n", result.Valid)

// Markdown Validation
mdValidator, _ := validator.NewValidator(validator.FormatMarkdown)
result = mdValidator.ValidateString(`# Title

This is **bold** text.

- Item 1
- Item 2
`)
fmt.Printf("Markdown valid: %v\n", result.Valid)

// JSON Lines Validation
jsonlValidator, _ := validator.NewValidator(validator.FormatJSONL)
result = jsonlValidator.ValidateString(`{"event": "login", "user": "john"}
{"event": "logout", "user": "john"}
{"event": "login", "user": "jane"}`)
fmt.Printf("JSONL valid: %v\n", result.Valid)

// Jupyter Notebook Validation
jupyterValidator, _ := validator.NewValidator(validator.FormatJupyter)
result = jupyterValidator.ValidateString(`{
  "cells": [],
  "metadata": {"kernelspec": {"name": "python3"}},
  "nbformat": 4,
  "nbformat_minor": 2
}`)
fmt.Printf("Jupyter valid: %v\n", result.Valid)

// Requirements.txt Validation
reqValidator, _ := validator.NewValidator(validator.FormatRequirements)
result = reqValidator.ValidateString(`numpy==1.21.0
pandas>=1.3.0
scikit-learn~=1.0.0`)
fmt.Printf("Requirements valid: %v\n", result.Valid)

// Dockerfile Validation
dockerValidator, _ := validator.NewValidator(validator.FormatDockerfile)
result = dockerValidator.ValidateString(`FROM python:3.9
WORKDIR /app
COPY . .
RUN pip install -r requirements.txt
CMD ["python", "app.py"]`)
fmt.Printf("Dockerfile valid: %v\n", result.Valid)
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