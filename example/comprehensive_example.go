package main

import (
	"fmt"
	"log"
	"os"

	"github.com/akhilesharora/serdeval"
)

func main() {
	// Example 1: Basic JSON validation
	fmt.Println("=== Example 1: Basic JSON Validation ===")
	jsonValidator, _ := serdeval.NewValidator(serdeval.FormatJSON)
	result := jsonValidator.ValidateString(`{"name": "John", "age": 30}`)
	fmt.Printf("Valid: %v\n", result.Valid)
	if !result.Valid {
		fmt.Printf("Error: %s\n", result.Error)
	}

	// Example 2: Invalid JSON
	fmt.Println("\n=== Example 2: Invalid JSON ===")
	result = jsonValidator.ValidateString(`{"name": "John", "age": }`)
	fmt.Printf("Valid: %v\n", result.Valid)
	fmt.Printf("Error: %s\n", result.Error)

	// Example 3: Auto-detect format
	fmt.Println("\n=== Example 3: Auto-detect Format ===")
	yamlData := []byte(`
name: John
age: 30
hobbies:
  - reading
  - coding
`)
	result = serdeval.ValidateAuto(yamlData)
	fmt.Printf("Detected format: %s\n", result.Format)
	fmt.Printf("Valid: %v\n", result.Valid)

	// Example 4: Validate from file
	fmt.Println("\n=== Example 4: Validate from File ===")
	if len(os.Args) > 1 {
		filename := os.Args[1]
		data, err := os.ReadFile(filename) // #nosec G304 - This is an example that reads user-provided files
		if err != nil {
			log.Printf("Error reading file: %v", err)
		} else {
			// Detect format from filename
			format := serdeval.DetectFormatFromFilename(filename)
			if format == serdeval.FormatUnknown {
				// Try content-based detection
				format = serdeval.DetectFormat(data)
			}

			if format != serdeval.FormatUnknown {
				v, _ := serdeval.NewValidator(format)
				result := v.Validate(data)
				fmt.Printf("File: %s\n", filename)
				fmt.Printf("Format: %s\n", format)
				fmt.Printf("Valid: %v\n", result.Valid)
				if !result.Valid {
					fmt.Printf("Error: %s\n", result.Error)
				}
			} else {
				fmt.Println("Unable to detect file format")
			}
		}
	}

	// Example 5: Batch validation
	fmt.Println("\n=== Example 5: Batch Validation ===")
	testData := map[string]string{
		"JSON": `{"valid": true}`,
		"YAML": `valid: true`,
		"XML":  `<root><valid>true</valid></root>`,
		"TOML": `valid = true`,
	}

	for name, data := range testData {
		result := serdeval.ValidateAuto([]byte(data))
		fmt.Printf("%s: Format=%s, Valid=%v\n", name, result.Format, result.Valid)
	}

	// Example 6: Custom error handling
	fmt.Println("\n=== Example 6: Custom Error Handling ===")
	validateWithRetry := func(data string, format serdeval.Format) {
		v, err := serdeval.NewValidator(format)
		if err != nil {
			fmt.Printf("Failed to create validator: %v\n", err)

			return
		}

		result := v.ValidateString(data)
		if !result.Valid {
			fmt.Printf("Validation failed for %s: %s\n", format, result.Error)
			// You could implement retry logic, logging, etc. here
		} else {
			fmt.Printf("Successfully validated %s data\n", format)
		}
	}

	validateWithRetry(`{"test": true}`, serdeval.FormatJSON)
	validateWithRetry(`<invalid`, serdeval.FormatXML)
}
