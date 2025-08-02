// Package validator provides data format validation for JSON, YAML, XML, and TOML
package validator

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

// Format represents a data format type
type Format string

const (
	// FormatJSON represents JSON format
	FormatJSON Format = "json"
	// FormatYAML represents YAML format
	FormatYAML Format = "yaml"
	// FormatXML represents XML format
	FormatXML Format = "xml"
	// FormatTOML represents TOML format
	FormatTOML Format = "toml"
	// FormatAuto represents automatic format detection
	FormatAuto Format = "auto"
	// FormatUnknown represents unknown format
	FormatUnknown Format = "unknown"
)

// Result contains the validation result for a data format validation operation.
type Result struct {
	// Valid indicates whether the data is valid for the detected/specified format
	Valid bool `json:"valid"`
	// Format indicates the data format that was validated
	Format Format `json:"format"`
	// Error contains the validation error message if Valid is false
	Error string `json:"error,omitempty"`
	// FileName is an optional field to track which file was validated
	FileName string `json:"filename,omitempty"`
}

// Validator is the main validator interface for validating data formats.
// Each validator is specific to a single format (JSON, YAML, XML, or TOML).
type Validator interface {
	// Validate checks if the provided byte slice is valid for this validator's format.
	// Returns a Result containing validation status and any error messages.
	Validate(data []byte) Result

	// ValidateString is a convenience method that validates a string.
	// Internally converts the string to []byte and calls Validate.
	ValidateString(data string) Result

	// Format returns the data format this validator is configured for.
	Format() Format
}

// baseValidator implements common validation logic
type baseValidator struct {
	format Format
}

// JSONValidator validates JSON data
type JSONValidator struct {
	baseValidator
}

// YAMLValidator validates YAML data
type YAMLValidator struct {
	baseValidator
}

// XMLValidator validates XML data
type XMLValidator struct {
	baseValidator
}

// TOMLValidator validates TOML data
type TOMLValidator struct {
	baseValidator
}

// NewValidator creates a new validator for the specified format.
//
// Example:
//
//	validator, err := NewValidator(FormatJSON)
//	if err != nil {
//		log.Fatal(err)
//	}
//	result := validator.ValidateString(`{"key": "value"}`)
//	if result.Valid {
//		fmt.Println("Valid JSON!")
//	}
//
// Supported formats: FormatJSON, FormatYAML, FormatXML, FormatTOML
// Returns an error if an unsupported format is specified.
func NewValidator(format Format) (Validator, error) {
	switch format {
	case FormatJSON:
		return &JSONValidator{baseValidator{format: FormatJSON}}, nil
	case FormatYAML:
		return &YAMLValidator{baseValidator{format: FormatYAML}}, nil
	case FormatXML:
		return &XMLValidator{baseValidator{format: FormatXML}}, nil
	case FormatTOML:
		return &TOMLValidator{baseValidator{format: FormatTOML}}, nil
	case FormatAuto, FormatUnknown:
		return nil, fmt.Errorf("unsupported format: %s", format)
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// Format returns the validator's format
func (v baseValidator) Format() Format {
	return v.format
}

// Validate validates JSON data
func (v *JSONValidator) Validate(data []byte) Result {
	var jsonData interface{}
	err := json.Unmarshal(data, &jsonData)

	return Result{
		Valid:  err == nil,
		Format: v.format,
		Error:  errorString(err),
	}
}

// ValidateString validates JSON string
func (v *JSONValidator) ValidateString(data string) Result {
	return v.Validate([]byte(data))
}

// Validate validates YAML data
func (v *YAMLValidator) Validate(data []byte) Result {
	var yamlData interface{}
	err := yaml.Unmarshal(data, &yamlData)

	return Result{
		Valid:  err == nil,
		Format: v.format,
		Error:  errorString(err),
	}
}

// ValidateString validates YAML string
func (v *YAMLValidator) ValidateString(data string) Result {
	return v.Validate([]byte(data))
}

// Validate validates XML data
func (v *XMLValidator) Validate(data []byte) Result {
	var xmlData interface{}
	err := xml.Unmarshal(data, &xmlData)

	return Result{
		Valid:  err == nil,
		Format: v.format,
		Error:  errorString(err),
	}
}

// ValidateString validates XML string
func (v *XMLValidator) ValidateString(data string) Result {
	return v.Validate([]byte(data))
}

// Validate validates TOML data
func (v *TOMLValidator) Validate(data []byte) Result {
	var tomlData interface{}
	err := toml.Unmarshal(data, &tomlData)

	return Result{
		Valid:  err == nil,
		Format: v.format,
		Error:  errorString(err),
	}
}

// ValidateString validates TOML string
func (v *TOMLValidator) ValidateString(data string) Result {
	return v.Validate([]byte(data))
}

// ValidateAuto validates data with automatic format detection.
// It first attempts to detect the format, then validates using the appropriate validator.
//
// Example:
//
//	data := []byte(`{"name": "test"}`)
//	result := ValidateAuto(data)
//	fmt.Printf("Format: %s, Valid: %v\n", result.Format, result.Valid)
//	// Output: Format: json, Valid: true
//
// Returns a Result with Format=FormatUnknown if the format cannot be detected.
func ValidateAuto(data []byte) Result {
	format := DetectFormat(data)
	if format == FormatUnknown {
		return Result{
			Valid:  false,
			Format: FormatUnknown,
			Error:  "unable to detect format",
		}
	}

	validator, err := NewValidator(format)
	if err != nil {
		return Result{
			Valid:  false,
			Format: format,
			Error:  err.Error(),
		}
	}

	result := validator.Validate(data)

	return result
}

// DetectFormat attempts to detect the data format by analyzing the content.
// Uses simple heuristics to identify JSON, XML, YAML, or TOML formats.
//
// Detection rules:
//   - JSON: Starts with '{' or '['
//   - XML: Starts with '<?xml' or '<'
//   - YAML: Contains '---' or has key:value pattern
//   - TOML: Contains '[section]' pattern with key=value pairs
//
// Returns FormatUnknown if the format cannot be determined.
func DetectFormat(data []byte) Format {
	trimmed := strings.TrimSpace(string(data))
	if len(trimmed) == 0 {
		return FormatUnknown
	}

	// Check JSON
	if trimmed[0] == '{' || trimmed[0] == '[' {
		return FormatJSON
	}

	// Check XML
	if strings.HasPrefix(trimmed, "<?xml") || strings.HasPrefix(trimmed, "<") {
		return FormatXML
	}

	// Check YAML
	if strings.Contains(trimmed, "---") || strings.Contains(trimmed, ":") {
		return FormatYAML
	}

	// Check TOML - simple key=value pattern
	if strings.Contains(trimmed, "=") && !strings.Contains(trimmed, ":") {
		// Make sure it's not JSON or XML
		if !strings.HasPrefix(trimmed, "{") && !strings.HasPrefix(trimmed, "[") && !strings.HasPrefix(trimmed, "<") {
			return FormatTOML
		}
	}

	return FormatUnknown
}

// DetectFormatFromFilename attempts to detect format from filename extension.
//
// Supported extensions:
//   - .json → FormatJSON
//   - .yaml, .yml → FormatYAML
//   - .xml → FormatXML
//   - .toml → FormatTOML
//
// Example:
//
//	format := DetectFormatFromFilename("config.json")
//	// format == FormatJSON
//
// Returns FormatUnknown if the extension is not recognized.
func DetectFormatFromFilename(filename string) Format {
	lastDot := strings.LastIndex(filename, ".")
	if lastDot == -1 {
		return FormatUnknown
	}
	ext := strings.ToLower(strings.TrimPrefix(filename[lastDot:], "."))

	switch ext {
	case "json":
		return FormatJSON
	case "yaml", "yml":
		return FormatYAML
	case "xml":
		return FormatXML
	case "toml":
		return FormatTOML
	default:
		return FormatUnknown
	}
}

// errorString returns empty string if error is nil
func errorString(err error) string {
	if err == nil {
		return ""
	}

	return err.Error()
}
