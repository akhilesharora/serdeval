/*
Package validator provides a privacy-focused data format validation library for Go.

It supports validation of JSON, YAML, XML, and TOML formats with automatic format detection.

Features:
  - Zero dependencies for core validation logic
  - No logging or data retention
  - Format auto-detection
  - Simple, clean API
  - Thread-safe validators

Basic usage:

	import "github.com/akhilesharora/serdeval/pkg/validator"

	// Validate with explicit format
	v, _ := validator.NewValidator(validator.FormatJSON)
	result := v.ValidateString(`{"test": true}`)
	if result.Valid {
		fmt.Println("Valid JSON!")
	}

	// Validate with auto-detection
	result := validator.ValidateAuto([]byte(data))
	fmt.Printf("Format: %s, Valid: %v\n", result.Format, result.Valid)

Advanced usage:

	// Create reusable validators
	jsonValidator, _ := validator.NewValidator(validator.FormatJSON)
	yamlValidator, _ := validator.NewValidator(validator.FormatYAML)

	// Validate multiple files
	for _, file := range files {
		data, _ := os.ReadFile(file)
		format := validator.DetectFormatFromFilename(file)

		v, _ := validator.NewValidator(format)
		result := v.Validate(data)

		if !result.Valid {
			log.Printf("%s: %s", file, result.Error)
		}
	}

Privacy guarantee:
  - No network connections
  - No temporary files
  - No logging of validated data
  - All validation happens in-memory
*/
package validator
