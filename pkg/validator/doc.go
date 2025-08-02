/*
Package validator provides comprehensive data format validation for Go applications.

The validator package supports validation of multiple data formats including JSON, YAML, XML, TOML,
CSV, GraphQL, INI, HCL, Protobuf text format, Markdown, JSON Lines, Jupyter Notebooks,
Python requirements.txt, and Dockerfiles. It features automatic format detection and a unified API
for all supported formats.

# Features

  - Support for 15+ data formats
  - Automatic format detection from content or filename
  - Zero dependencies for core validation logic
  - Thread-safe validator implementations
  - Simple, consistent API across all formats
  - Privacy-focused: no logging, network calls, or data retention
  - Comprehensive error messages for validation failures

# Basic Usage

Create a validator for a specific format:

	import "github.com/akhilesharora/datavalidator/pkg/validator"

	// Create a JSON validator
	v, err := validator.NewValidator(validator.FormatJSON)
	if err != nil {
		log.Fatal(err)
	}

	// Validate a string
	result := v.ValidateString(`{"name": "test", "value": 123}`)
	if result.Valid {
		fmt.Println("Valid JSON!")
	} else {
		fmt.Printf("Invalid JSON: %s\n", result.Error)
	}

# Automatic Format Detection

Use ValidateAuto for automatic format detection:

	data := []byte(`{"key": "value"}`)
	result := validator.ValidateAuto(data)
	fmt.Printf("Detected format: %s\n", result.Format)
	fmt.Printf("Valid: %v\n", result.Valid)

Detect format from filename:

	format := validator.DetectFormatFromFilename("config.yaml")
	// format == validator.FormatYAML

# Supported Formats

The package supports the following formats:

  - JSON (FormatJSON): RFC 7159 compliant JSON validation
  - YAML (FormatYAML): YAML 1.2 specification
  - XML (FormatXML): Well-formed XML validation
  - TOML (FormatTOML): TOML v1.0.0 format
  - CSV (FormatCSV): Comma-separated values with consistent columns
  - GraphQL (FormatGraphQL): GraphQL queries, mutations, and schemas
  - INI (FormatINI): INI configuration files with sections
  - HCL (FormatHCL): HashiCorp Configuration Language (HCL2)
  - Protobuf (FormatProtobuf): Protocol Buffers text format
  - Markdown (FormatMarkdown): CommonMark specification
  - JSON Lines (FormatJSONL): Newline-delimited JSON
  - Jupyter (FormatJupyter): Jupyter Notebook .ipynb files
  - Requirements (FormatRequirements): Python requirements.txt
  - Dockerfile (FormatDockerfile): Docker container definitions

# Advanced Usage

Validate multiple files with different formats:

	files := []string{"config.json", "data.yaml", "schema.graphql"}

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			log.Printf("Error reading %s: %v", file, err)
			continue
		}

		// Detect format from filename
		format := validator.DetectFormatFromFilename(file)
		if format == validator.FormatUnknown {
			// Fall back to content detection
			format = validator.DetectFormat(data)
		}

		// Create appropriate validator
		v, err := validator.NewValidator(format)
		if err != nil {
			log.Printf("Unsupported format for %s", file)
			continue
		}

		// Validate the file
		result := v.Validate(data)
		result.FileName = file

		if result.Valid {
			fmt.Printf("✓ %s: Valid %s\n", file, result.Format)
		} else {
			fmt.Printf("✗ %s: %s\n", file, result.Error)
		}
	}

# Format Detection Heuristics

The package uses intelligent heuristics for format detection:

  - File extension matching for reliable format identification
  - Content-based detection using format-specific patterns
  - Precedence given to more specific formats (e.g., Jupyter over JSON)
  - Support for files without extensions (e.g., Dockerfile)

# Privacy and Security

This package is designed with privacy and security in mind:

  - All validation is performed in-memory
  - No network connections are made
  - No temporary files are created
  - No data is logged or retained
  - Input data is never modified

# Thread Safety

All validator implementations are thread-safe and can be reused across multiple goroutines:

	validator, _ := validator.NewValidator(validator.FormatJSON)

	// Safe to use concurrently
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(data string) {
			defer wg.Done()
			result := validator.ValidateString(data)
			// Process result
		}(jsonData[i])
	}
	wg.Wait()

# Error Handling

Validation errors include detailed information about what went wrong:

	result := validator.ValidateString(invalidJSON)
	if !result.Valid {
		// result.Error contains specific error message
		// e.g., "unexpected end of JSON input"
		fmt.Printf("Validation failed: %s\n", result.Error)
	}
*/
package validator
