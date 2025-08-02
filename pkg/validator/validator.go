// Package validator provides data format validation for JSON, YAML, XML, TOML, CSV, GraphQL, INI, HCL,
// Protobuf text format, Markdown, JSON Lines, Jupyter Notebooks, Requirements.txt, and Dockerfile
package validator

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/yuin/goldmark"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/types/known/anypb"
	"gopkg.in/ini.v1"
	"gopkg.in/yaml.v3"
)

// Format represents a supported data format type.
// It is used to specify which validator to use and identify detected formats.
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
	// FormatCSV represents CSV format
	FormatCSV Format = "csv"
	// FormatGraphQL represents GraphQL query/schema format
	FormatGraphQL Format = "graphql"
	// FormatINI represents INI configuration format
	FormatINI Format = "ini"
	// FormatHCL represents HCL (HashiCorp Configuration Language) format
	FormatHCL Format = "hcl"
	// FormatProtobuf represents Protobuf text format
	FormatProtobuf Format = "protobuf"
	// FormatMarkdown represents Markdown format
	FormatMarkdown Format = "markdown"
	// FormatJSONL represents JSON Lines format (newline-delimited JSON)
	FormatJSONL Format = "jsonl"
	// FormatJupyter represents Jupyter Notebook format
	FormatJupyter Format = "jupyter"
	// FormatRequirements represents Requirements.txt format
	FormatRequirements Format = "requirements"
	// FormatDockerfile represents Dockerfile format
	FormatDockerfile Format = "dockerfile"
	// FormatAuto represents automatic format detection
	FormatAuto Format = "auto"
	// FormatUnknown represents unknown format
	FormatUnknown Format = "unknown"
)

// Result contains the validation result for a data format validation operation.
// It provides information about whether the validation succeeded, the format
// that was validated, and any error messages if validation failed.
//
// Example:
//
//	result := validator.ValidateString(`{"key": "value"}`)
//	if result.Valid {
//		fmt.Printf("Valid %s data\n", result.Format)
//	} else {
//		fmt.Printf("Invalid data: %s\n", result.Error)
//	}
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

// Validator is the main interface for validating data formats.
// Each validator implementation is specific to a single format.
//
// Implementations are thread-safe and can be reused for multiple validations.
//
// Example:
//
//	validator, err := NewValidator(FormatJSON)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Validate bytes
//	result := validator.Validate(jsonBytes)
//
//	// Or validate string
//	result := validator.ValidateString(jsonString)
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

// baseValidator provides common functionality for all validator implementations.
// It is embedded in specific validator types to share the Format() method.
type baseValidator struct {
	format Format
}

// JSONValidator validates JSON data according to RFC 7159.
// It supports validation of both JSON objects and arrays.
//
// Example:
//
//	validator := &JSONValidator{baseValidator{format: FormatJSON}}
//	result := validator.ValidateString(`{"name": "test", "value": 123}`)
type JSONValidator struct {
	baseValidator
}

// YAMLValidator validates YAML data according to YAML 1.2 specification.
// It supports all standard YAML features including anchors, aliases, and multi-document streams.
//
// Example:
//
//	validator := &YAMLValidator{baseValidator{format: FormatYAML}}
//	result := validator.ValidateString("name: test\nvalue: 123")
type YAMLValidator struct {
	baseValidator
}

// XMLValidator validates XML data for well-formedness.
// It checks that the XML is properly structured with matching tags and valid syntax.
//
// Note: This validator checks for well-formedness only, not validity against a schema.
//
// Example:
//
//	validator := &XMLValidator{baseValidator{format: FormatXML}}
//	result := validator.ValidateString(`<root><item>test</item></root>`)
type XMLValidator struct {
	baseValidator
}

// TOMLValidator validates TOML (Tom's Obvious, Minimal Language) data.
// It supports all TOML v1.0.0 features including tables, arrays, and inline tables.
//
// Example:
//
//	validator := &TOMLValidator{baseValidator{format: FormatTOML}}
//	result := validator.ValidateString(`[server]\nhost = "localhost"\nport = 8080`)
type TOMLValidator struct {
	baseValidator
}

// CSVValidator validates CSV (Comma-Separated Values) data.
// It checks that the data can be parsed as valid CSV with consistent column counts.
//
// Example:
//
//	validator := &CSVValidator{baseValidator{format: FormatCSV}}
//	result := validator.ValidateString("name,age\nJohn,30\nJane,25")
type CSVValidator struct {
	baseValidator
}

// GraphQLValidator validates GraphQL queries, mutations, subscriptions, and schema definitions.
// It uses the GraphQL parser to ensure syntactic validity.
//
// Example:
//
//	validator := &GraphQLValidator{baseValidator{format: FormatGraphQL}}
//	result := validator.ValidateString(`query { user(id: "123") { name email } }`)
type GraphQLValidator struct {
	baseValidator
}

// INIValidator validates INI configuration file format.
// It supports sections, key-value pairs, and comments.
//
// Example:
//
//	validator := &INIValidator{baseValidator{format: FormatINI}}
//	result := validator.ValidateString(`[database]\nhost = localhost\nport = 5432`)
type INIValidator struct {
	baseValidator
}

// HCLValidator validates HCL (HashiCorp Configuration Language) data.
// It supports HCL2 syntax used in Terraform, Packer, and other HashiCorp tools.
//
// Example:
//
//	validator := &HCLValidator{baseValidator{format: FormatHCL}}
//	result := validator.ValidateString(`resource "aws_instance" "example" { ami = "ami-123" }`)
type HCLValidator struct {
	baseValidator
}

// ProtobufValidator validates Protocol Buffers text format data.
// It checks that the data can be parsed as valid protobuf text format.
//
// Example:
//
//	validator := &ProtobufValidator{baseValidator{format: FormatProtobuf}}
//	result := validator.ValidateString(`type_url: "type.googleapis.com/Example" value: "\x08\x01"`)
type ProtobufValidator struct {
	baseValidator
}

// MarkdownValidator validates Markdown formatted text.
// It uses the CommonMark specification to parse and validate the content.
//
// Example:
//
//	validator := &MarkdownValidator{baseValidator{format: FormatMarkdown}}
//	result := validator.ValidateString("# Title\n\nThis is **bold** text.")
type MarkdownValidator struct {
	baseValidator
}

// JSONLValidator validates JSON Lines (newline-delimited JSON) data.
// Each line must be a valid JSON object or array.
//
// Example:
//
//	validator := &JSONLValidator{baseValidator{format: FormatJSONL}}
//	result := validator.ValidateString(`{"name": "John"}\n{"name": "Jane"}`)
type JSONLValidator struct {
	baseValidator
}

// JupyterValidator validates Jupyter Notebook (.ipynb) files.
// It checks for the required notebook structure including cells, metadata, and nbformat.
//
// Example:
//
//	validator := &JupyterValidator{baseValidator{format: FormatJupyter}}
//	result := validator.Validate(jupyterNotebookBytes)
type JupyterValidator struct {
	baseValidator
}

// RequirementsValidator validates Python requirements.txt file format.
// It supports package names with version specifiers (==, >=, <=, ~=) and comments.
//
// Example:
//
//	validator := &RequirementsValidator{baseValidator{format: FormatRequirements}}
//	result := validator.ValidateString("flask==2.0.1\nrequests>=2.25.0")
type RequirementsValidator struct {
	baseValidator
}

// DockerfileValidator validates Dockerfile syntax.
// It checks for valid Docker instructions and ensures at least one FROM instruction exists.
//
// Example:
//
//	validator := &DockerfileValidator{baseValidator{format: FormatDockerfile}}
//	result := validator.ValidateString("FROM golang:1.19\nWORKDIR /app\nCOPY . .")
type DockerfileValidator struct {
	baseValidator
}

// validatorMap maps formats to their validator constructors
var validatorMap = map[Format]func() Validator{
	FormatJSON:         func() Validator { return &JSONValidator{baseValidator{format: FormatJSON}} },
	FormatYAML:         func() Validator { return &YAMLValidator{baseValidator{format: FormatYAML}} },
	FormatXML:          func() Validator { return &XMLValidator{baseValidator{format: FormatXML}} },
	FormatTOML:         func() Validator { return &TOMLValidator{baseValidator{format: FormatTOML}} },
	FormatCSV:          func() Validator { return &CSVValidator{baseValidator{format: FormatCSV}} },
	FormatGraphQL:      func() Validator { return &GraphQLValidator{baseValidator{format: FormatGraphQL}} },
	FormatINI:          func() Validator { return &INIValidator{baseValidator{format: FormatINI}} },
	FormatHCL:          func() Validator { return &HCLValidator{baseValidator{format: FormatHCL}} },
	FormatProtobuf:     func() Validator { return &ProtobufValidator{baseValidator{format: FormatProtobuf}} },
	FormatMarkdown:     func() Validator { return &MarkdownValidator{baseValidator{format: FormatMarkdown}} },
	FormatJSONL:        func() Validator { return &JSONLValidator{baseValidator{format: FormatJSONL}} },
	FormatJupyter:      func() Validator { return &JupyterValidator{baseValidator{format: FormatJupyter}} },
	FormatRequirements: func() Validator { return &RequirementsValidator{baseValidator{format: FormatRequirements}} },
	FormatDockerfile:   func() Validator { return &DockerfileValidator{baseValidator{format: FormatDockerfile}} },
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
// Supported formats: FormatJSON, FormatYAML, FormatXML, FormatTOML, FormatCSV, FormatGraphQL,
// FormatINI, FormatHCL, FormatProtobuf, FormatMarkdown, FormatJSONL, FormatJupyter,
// FormatRequirements, FormatDockerfile
// Returns an error if an unsupported format is specified.
func NewValidator(format Format) (Validator, error) {
	constructor, ok := validatorMap[format]
	if !ok {
		return nil, fmt.Errorf("unsupported format: %s", format)
	}

	return constructor(), nil
}

// Format returns the data format type associated with this validator.
// This method is available on all validator implementations.
//
// Example:
//
//	validator, _ := NewValidator(FormatJSON)
//	fmt.Println(validator.Format()) // Output: json
func (v baseValidator) Format() Format {
	return v.format
}

// Validate checks if the provided byte slice contains valid JSON data.
// It attempts to unmarshal the data and returns a Result indicating success or failure.
//
// The validation checks for proper JSON syntax including:
//   - Matching braces and brackets
//   - Proper string escaping
//   - Valid number formats
//   - Correct use of null, true, and false
//
// Example:
//
//	validator := &JSONValidator{baseValidator{format: FormatJSON}}
//	result := validator.Validate([]byte(`{"valid": true}`))
//	if result.Valid {
//		fmt.Println("Valid JSON!")
//	}
func (v *JSONValidator) Validate(data []byte) Result {
	var jsonData interface{}
	err := json.Unmarshal(data, &jsonData)

	return Result{
		Valid:  err == nil,
		Format: v.format,
		Error:  errorString(err),
	}
}

// ValidateString is a convenience method that validates a JSON string.
// It converts the string to bytes and calls Validate.
//
// Example:
//
//	validator := &JSONValidator{baseValidator{format: FormatJSON}}
//	result := validator.ValidateString(`{"name": "test"}`)
func (v *JSONValidator) ValidateString(data string) Result {
	return v.Validate([]byte(data))
}

// Validate checks if the provided byte slice contains valid YAML data.
// It supports all YAML 1.2 features including multi-document streams.
//
// Example:
//
//	validator := &YAMLValidator{baseValidator{format: FormatYAML}}
//	result := validator.Validate([]byte("key: value\nlist:\n  - item1\n  - item2"))
func (v *YAMLValidator) Validate(data []byte) Result {
	var yamlData interface{}
	err := yaml.Unmarshal(data, &yamlData)

	return Result{
		Valid:  err == nil,
		Format: v.format,
		Error:  errorString(err),
	}
}

// ValidateString is a convenience method that validates a YAML string.
// It converts the string to bytes and calls Validate.
//
// Example:
//
//	validator := &YAMLValidator{baseValidator{format: FormatYAML}}
//	result := validator.ValidateString("name: test\nvalue: 123")
func (v *YAMLValidator) ValidateString(data string) Result {
	return v.Validate([]byte(data))
}

// Validate checks if the provided byte slice contains well-formed XML data.
// Note that this only checks for well-formedness, not validity against a schema.
//
// Example:
//
//	validator := &XMLValidator{baseValidator{format: FormatXML}}
//	result := validator.Validate([]byte(`<?xml version="1.0"?><root></root>`))
func (v *XMLValidator) Validate(data []byte) Result {
	var xmlData interface{}
	err := xml.Unmarshal(data, &xmlData)

	return Result{
		Valid:  err == nil,
		Format: v.format,
		Error:  errorString(err),
	}
}

// ValidateString is a convenience method that validates an XML string.
// It converts the string to bytes and calls Validate.
//
// Example:
//
//	validator := &XMLValidator{baseValidator{format: FormatXML}}
//	result := validator.ValidateString(`<root><item>test</item></root>`)
func (v *XMLValidator) ValidateString(data string) Result {
	return v.Validate([]byte(data))
}

// Validate checks if the provided byte slice contains valid TOML data.
// It supports all TOML v1.0.0 features.
//
// Example:
//
//	validator := &TOMLValidator{baseValidator{format: FormatTOML}}
//	result := validator.Validate([]byte(`[server]\nport = 8080`))
func (v *TOMLValidator) Validate(data []byte) Result {
	var tomlData interface{}
	err := toml.Unmarshal(data, &tomlData)

	return Result{
		Valid:  err == nil,
		Format: v.format,
		Error:  errorString(err),
	}
}

// ValidateString is a convenience method that validates a TOML string.
// It converts the string to bytes and calls Validate.
//
// Example:
//
//	validator := &TOMLValidator{baseValidator{format: FormatTOML}}
//	result := validator.ValidateString(`title = "TOML Example"`)
func (v *TOMLValidator) ValidateString(data string) Result {
	return v.Validate([]byte(data))
}

// Validate checks if the provided byte slice contains valid CSV data.
// It reads all records to ensure consistent column counts and proper formatting.
//
// Example:
//
//	validator := &CSVValidator{baseValidator{format: FormatCSV}}
//	result := validator.Validate([]byte("name,age\nJohn,30"))
func (v *CSVValidator) Validate(data []byte) Result {
	r := csv.NewReader(strings.NewReader(string(data)))
	// Read all records to validate
	_, err := r.ReadAll()

	return Result{
		Valid:  err == nil,
		Format: v.format,
		Error:  errorString(err),
	}
}

// ValidateString is a convenience method that validates a CSV string.
// It converts the string to bytes and calls Validate.
//
// Example:
//
//	validator := &CSVValidator{baseValidator{format: FormatCSV}}
//	result := validator.ValidateString("header1,header2\nvalue1,value2")
func (v *CSVValidator) ValidateString(data string) Result {
	return v.Validate([]byte(data))
}

// Validate checks if the provided byte slice contains valid GraphQL syntax.
// It can validate queries, mutations, subscriptions, and schema definitions.
// Empty content is considered invalid.
//
// Example:
//
//	validator := &GraphQLValidator{baseValidator{format: FormatGraphQL}}
//	result := validator.Validate([]byte(`query { user { name } }`))
func (v *GraphQLValidator) Validate(data []byte) Result {
	// GraphQL requires non-empty content
	if len(data) == 0 {
		return Result{
			Valid:  false,
			Format: v.format,
			Error:  "empty GraphQL content",
		}
	}
	s := source.NewSource(&source.Source{
		Body: data,
		Name: "GraphQL",
	})
	_, err := parser.Parse(parser.ParseParams{Source: s})

	return Result{
		Valid:  err == nil,
		Format: v.format,
		Error:  errorString(err),
	}
}

// ValidateString is a convenience method that validates a GraphQL string.
// It converts the string to bytes and calls Validate.
//
// Example:
//
//	validator := &GraphQLValidator{baseValidator{format: FormatGraphQL}}
//	result := validator.ValidateString(`mutation { createUser(name: "John") { id } }`)
func (v *GraphQLValidator) ValidateString(data string) Result {
	return v.Validate([]byte(data))
}

// Validate checks if the provided byte slice contains valid INI format data.
// It supports sections, key-value pairs, and comments.
//
// Example:
//
//	validator := &INIValidator{baseValidator{format: FormatINI}}
//	result := validator.Validate([]byte(`[section]\nkey = value`))
func (v *INIValidator) Validate(data []byte) Result {
	_, err := ini.Load(data)

	return Result{
		Valid:  err == nil,
		Format: v.format,
		Error:  errorString(err),
	}
}

// ValidateString is a convenience method that validates an INI string.
// It converts the string to bytes and calls Validate.
//
// Example:
//
//	validator := &INIValidator{baseValidator{format: FormatINI}}
//	result := validator.ValidateString("[database]\nhost = localhost")
func (v *INIValidator) ValidateString(data string) Result {
	return v.Validate([]byte(data))
}

// Validate checks if the provided byte slice contains valid HCL2 syntax.
// It uses the HashiCorp HCL parser to validate the configuration.
//
// Example:
//
//	validator := &HCLValidator{baseValidator{format: FormatHCL}}
//	result := validator.Validate([]byte(`variable "region" { default = "us-west-2" }`))
func (v *HCLValidator) Validate(data []byte) Result {
	_, diags := hclsyntax.ParseConfig(data, "hcl", hcl.InitialPos)
	var errStr string
	if diags.HasErrors() {
		errStr = diags.Error()
	}

	return Result{
		Valid:  !diags.HasErrors(),
		Format: v.format,
		Error:  errStr,
	}
}

// ValidateString is a convenience method that validates an HCL string.
// It converts the string to bytes and calls Validate.
//
// Example:
//
//	validator := &HCLValidator{baseValidator{format: FormatHCL}}
//	result := validator.ValidateString(`resource "aws_instance" "web" { ami = "ami-123" }`)
func (v *HCLValidator) ValidateString(data string) Result {
	return v.Validate([]byte(data))
}

// Validate checks if the provided byte slice contains valid Protocol Buffers text format.
// It attempts to unmarshal the data as protobuf text format.
//
// Example:
//
//	validator := &ProtobufValidator{baseValidator{format: FormatProtobuf}}
//	result := validator.Validate([]byte(`type_url: "example.com/Type"`))
func (v *ProtobufValidator) Validate(data []byte) Result {
	// Try to unmarshal as protobuf text format into Any message
	msg := &anypb.Any{}
	err := prototext.Unmarshal(data, msg)

	return Result{
		Valid:  err == nil,
		Format: v.format,
		Error:  errorString(err),
	}
}

// ValidateString is a convenience method that validates a Protobuf text format string.
// It converts the string to bytes and calls Validate.
//
// Example:
//
//	validator := &ProtobufValidator{baseValidator{format: FormatProtobuf}}
//	result := validator.ValidateString(`type_url: "type.googleapis.com/Example"`)
func (v *ProtobufValidator) ValidateString(data string) Result {
	return v.Validate([]byte(data))
}

// Validate checks if the provided byte slice contains valid Markdown.
// It uses the CommonMark specification to parse the content.
//
// Example:
//
//	validator := &MarkdownValidator{baseValidator{format: FormatMarkdown}}
//	result := validator.Validate([]byte("# Title\n\nParagraph with **bold** text."))
func (v *MarkdownValidator) Validate(data []byte) Result {
	md := goldmark.New()
	err := md.Convert(data, io.Discard)

	return Result{
		Valid:  err == nil,
		Format: v.format,
		Error:  errorString(err),
	}
}

// ValidateString is a convenience method that validates a Markdown string.
// It converts the string to bytes and calls Validate.
//
// Example:
//
//	validator := &MarkdownValidator{baseValidator{format: FormatMarkdown}}
//	result := validator.ValidateString("## Heading\n\n- List item 1\n- List item 2")
func (v *MarkdownValidator) ValidateString(data string) Result {
	return v.Validate([]byte(data))
}

// Validate checks if the provided byte slice contains valid JSON Lines format.
// Each non-empty line must be a valid JSON object or array.
// Empty lines are allowed and ignored.
//
// Example:
//
//	validator := &JSONLValidator{baseValidator{format: FormatJSONL}}
//	result := validator.Validate([]byte(`{"id":1}\n{"id":2}`))
func (v *JSONLValidator) Validate(data []byte) Result {
	if len(data) == 0 {
		return Result{
			Valid:  true,
			Format: v.format,
			Error:  "",
		}
	}

	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		// Skip empty lines
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Each line must be valid JSON
		var jsonData interface{}
		if err := json.Unmarshal([]byte(line), &jsonData); err != nil {
			return Result{
				Valid:  false,
				Format: v.format,
				Error:  fmt.Sprintf("invalid JSON on line %d: %s", i+1, err.Error()),
			}
		}
	}

	return Result{
		Valid:  true,
		Format: v.format,
		Error:  "",
	}
}

// ValidateString is a convenience method that validates a JSON Lines string.
// It converts the string to bytes and calls Validate.
//
// Example:
//
//	validator := &JSONLValidator{baseValidator{format: FormatJSONL}}
//	result := validator.ValidateString(`{"event":"start"}\n{"event":"end"}`)
func (v *JSONLValidator) ValidateString(data string) Result {
	return v.Validate([]byte(data))
}

// Validate checks if the provided byte slice contains a valid Jupyter Notebook.
// It verifies that the data is valid JSON and contains required notebook fields:
// cells, metadata, and nbformat.
//
// Example:
//
//	validator := &JupyterValidator{baseValidator{format: FormatJupyter}}
//	notebookData, _ := os.ReadFile("notebook.ipynb")
//	result := validator.Validate(notebookData)
func (v *JupyterValidator) Validate(data []byte) Result {
	var notebook map[string]interface{}
	if err := json.Unmarshal(data, &notebook); err != nil {
		return Result{
			Valid:  false,
			Format: v.format,
			Error:  "invalid JSON: " + err.Error(),
		}
	}

	// Check for required notebook fields
	if _, ok := notebook["cells"]; !ok {
		return Result{
			Valid:  false,
			Format: v.format,
			Error:  "missing required field: cells",
		}
	}
	if _, ok := notebook["metadata"]; !ok {
		return Result{
			Valid:  false,
			Format: v.format,
			Error:  "missing required field: metadata",
		}
	}
	if _, ok := notebook["nbformat"]; !ok {
		return Result{
			Valid:  false,
			Format: v.format,
			Error:  "missing required field: nbformat",
		}
	}

	return Result{
		Valid:  true,
		Format: v.format,
		Error:  "",
	}
}

// ValidateString is a convenience method that validates a Jupyter Notebook string.
// It converts the string to bytes and calls Validate.
//
// Example:
//
//	validator := &JupyterValidator{baseValidator{format: FormatJupyter}}
//	result := validator.ValidateString(notebookJSONString)
func (v *JupyterValidator) ValidateString(data string) Result {
	return v.Validate([]byte(data))
}

// Validate checks if the provided byte slice contains valid requirements.txt format.
// It supports package names with version specifiers and comments.
// Empty lines and lines starting with # are allowed.
//
// Example:
//
//	validator := &RequirementsValidator{baseValidator{format: FormatRequirements}}
//	result := validator.Validate([]byte("django==3.2\nrequests>=2.25.0"))
func (v *RequirementsValidator) Validate(data []byte) Result {
	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Basic validation: check if line contains package name
		// Valid formats: package, package==version, package>=version, etc.
		if !strings.ContainsAny(line, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-") {
			return Result{
				Valid:  false,
				Format: v.format,
				Error:  fmt.Sprintf("invalid requirement on line %d: %s", i+1, line),
			}
		}
	}

	return Result{
		Valid:  true,
		Format: v.format,
		Error:  "",
	}
}

// ValidateString is a convenience method that validates a requirements.txt string.
// It converts the string to bytes and calls Validate.
//
// Example:
//
//	validator := &RequirementsValidator{baseValidator{format: FormatRequirements}}
//	result := validator.ValidateString("numpy>=1.19.0\npandas==1.3.0")
func (v *RequirementsValidator) ValidateString(data string) Result {
	return v.Validate([]byte(data))
}

// Validate checks if the provided byte slice contains valid Dockerfile syntax.
// It verifies that valid Docker instructions are used and at least one FROM instruction exists.
// Line continuations with backslash are supported.
//
// Example:
//
//	validator := &DockerfileValidator{baseValidator{format: FormatDockerfile}}
//	result := validator.Validate([]byte("FROM alpine:latest\nRUN apk add --no-cache curl"))
func (v *DockerfileValidator) Validate(data []byte) Result {
	lines := strings.Split(string(data), "\n")
	hasFrom := false

	for i, line := range lines {
		line = strings.TrimSpace(line)
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Check for FROM instruction (required)
		if strings.HasPrefix(strings.ToUpper(line), "FROM ") {
			hasFrom = true
		}

		// Basic validation: check if line starts with valid instruction
		upperLine := strings.ToUpper(line)
		validInstructions := []string{"FROM", "RUN", "CMD", "LABEL", "EXPOSE", "ENV", "ADD", "COPY",
			"ENTRYPOINT", "VOLUME", "USER", "WORKDIR", "ARG", "ONBUILD", "STOPSIGNAL", "HEALTHCHECK", "SHELL"}

		hasValidInstruction := false
		for _, instruction := range validInstructions {
			if strings.HasPrefix(upperLine, instruction+" ") || upperLine == instruction {
				hasValidInstruction = true

				break
			}
		}

		// Allow line continuations
		if i > 0 && strings.HasSuffix(lines[i-1], "\\") {
			continue
		}

		if !hasValidInstruction {
			return Result{
				Valid:  false,
				Format: v.format,
				Error:  fmt.Sprintf("invalid instruction on line %d: %s", i+1, line),
			}
		}
	}

	if !hasFrom {
		return Result{
			Valid:  false,
			Format: v.format,
			Error:  "missing required FROM instruction",
		}
	}

	return Result{
		Valid:  true,
		Format: v.format,
		Error:  "",
	}
}

// ValidateString is a convenience method that validates a Dockerfile string.
// It converts the string to bytes and calls Validate.
//
// Example:
//
//	validator := &DockerfileValidator{baseValidator{format: FormatDockerfile}}
//	result := validator.ValidateString("FROM node:16\nWORKDIR /app\nCOPY . .")
func (v *DockerfileValidator) ValidateString(data string) Result {
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

// isJupyterNotebook checks if the content appears to be a Jupyter notebook.
// It looks for the required JSON structure with cells, metadata, and nbformat fields.
func isJupyterNotebook(trimmed string) bool {
	return strings.HasPrefix(trimmed, "{") &&
		strings.Contains(trimmed, "\"cells\"") &&
		strings.Contains(trimmed, "\"metadata\"") &&
		strings.Contains(trimmed, "\"nbformat\"")
}

// isJSONLines checks if the content appears to be JSON Lines format.
// Returns true if there are multiple lines, each containing valid JSON.
func isJSONLines(lines []string) bool {
	if len(lines) <= 1 {
		return false
	}

	validLines := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if !isJSONLine(line) {
			return false
		}
		validLines++
	}

	return validLines > 1
}

// isJSONLine checks if a single line appears to be valid JSON.
// It performs a simple syntax check for JSON objects or arrays.
func isJSONLine(line string) bool {
	return (strings.HasPrefix(line, "{") && strings.HasSuffix(line, "}")) ||
		(strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]"))
}

// isJSON checks if the content appears to be regular JSON format.
// It looks for opening and closing braces or brackets.
func isJSON(trimmed string) bool {
	if len(trimmed) == 0 {
		return false
	}

	return (strings.HasPrefix(trimmed, "{") && strings.HasSuffix(trimmed, "}")) ||
		(strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]"))
}

// detectJSONFamily attempts to detect JSON-based formats.
// It checks for Jupyter notebooks first, then JSON Lines, then regular JSON.
// Returns FormatUnknown if no JSON-based format is detected.
func detectJSONFamily(trimmed string, lines []string) Format {
	// Check for Jupyter Notebook first (it's a specific type of JSON)
	if isJupyterNotebook(trimmed) {
		return FormatJupyter
	}

	// Check for JSON Lines before regular JSON
	if isJSONLines(lines) {
		return FormatJSONL
	}

	// Check JSON (after Jupyter and JSONL)
	if isJSON(trimmed) {
		return FormatJSON
	}

	return FormatUnknown
}

// countPatterns counts how many of the provided patterns are found in the text.
// Used for heuristic format detection based on keyword presence.
func countPatterns(text string, patterns []string) int {
	count := 0
	for _, pattern := range patterns {
		if strings.Contains(text, pattern) {
			count++
		}
	}

	return count
}

// isDockerfile checks if the content appears to be a Dockerfile.
// It looks for Docker instructions like FROM, RUN, CMD, etc.
// The input should be uppercase for case-insensitive matching.
func isDockerfile(upperTrimmed string) bool {
	dockerKeywords := []string{
		"FROM ", "RUN ", "CMD ", "EXPOSE ", "ENV ", "ADD ", "COPY ",
		"ENTRYPOINT ", "VOLUME ", "USER ", "WORKDIR ", "ARG ",
	}
	dockerScore := countPatterns(upperTrimmed, dockerKeywords)

	return strings.Contains(upperTrimmed, "FROM ") || dockerScore >= 3
}

// isHCL checks if the content appears to be HCL (HashiCorp Configuration Language).
// It looks for HCL-specific keywords like resource, variable, provider, etc.
func isHCL(trimmed string) bool {
	hclPatterns := []string{
		"resource ", "variable ", "provider ", "module ",
		"output ", "locals ", "terraform ", "data ",
	}
	hclScore := countPatterns(trimmed, hclPatterns)

	return hclScore > 0 &&
		strings.Contains(trimmed, "=") &&
		strings.Contains(trimmed, "\"") &&
		strings.Contains(trimmed, "{")
}

// isGraphQL checks if the content appears to be GraphQL format.
// It looks for GraphQL keywords like query, mutation, type, schema, etc.
func isGraphQL(trimmed string) bool {
	graphqlPatterns := []string{
		"query ", "mutation ", "subscription ", "fragment ", "type ",
		"interface ", "enum ", "input ", "scalar ", "schema ",
	}
	graphqlScore := countPatterns(trimmed, graphqlPatterns)
	if graphqlScore == 0 || !strings.Contains(trimmed, "{") || !strings.Contains(trimmed, "}") {
		return false
	}

	return strings.Contains(trimmed, "{\n") || strings.Contains(trimmed, "{ ") || graphqlScore >= 2
}

// isProtobuf checks if the content appears to be Protocol Buffers text format.
// It looks for protobuf-specific fields like type_url and value.
func isProtobuf(trimmed string) bool {
	return (strings.Contains(trimmed, "type_url:") || strings.Contains(trimmed, "value:")) &&
		strings.Contains(trimmed, "\"")
}

// detectDeveloperFormats attempts to detect developer tool formats.
// It checks for Dockerfile, HCL, GraphQL, and Protobuf formats in order of specificity.
// Returns FormatUnknown if no developer format is detected.
func detectDeveloperFormats(trimmed string, lines []string) Format {
	upperTrimmed := strings.ToUpper(trimmed)

	// Check Dockerfile - look for common Docker instructions
	if isDockerfile(upperTrimmed) {
		return FormatDockerfile
	}

	// Check HCL/Terraform before GraphQL (HCL is more specific)
	if isHCL(trimmed) {
		return FormatHCL
	}

	// Check GraphQL - look for GraphQL-specific patterns
	if isGraphQL(trimmed) {
		return FormatGraphQL
	}

	// Check for Protobuf text format before YAML (more specific)
	if isProtobuf(trimmed) {
		return FormatProtobuf
	}

	return FormatUnknown
}

// detectCSV checks if the content appears to be CSV format.
// It verifies that the content has commas and consistent column counts across rows.
func detectCSV(trimmed string, lines []string) bool {
	if !strings.Contains(trimmed, ",") || len(lines) <= 1 {
		return false
	}

	firstLineCommas := strings.Count(lines[0], ",")
	if firstLineCommas == 0 {
		return false
	}

	for i := 1; i < len(lines) && i < 5; i++ {
		if lines[i] != "" && strings.Count(lines[i], ",") != firstLineCommas {
			return false
		}
	}

	return true
}

// detectMarkdown checks if the content appears to be Markdown format.
// It looks for common Markdown syntax like headers (#), code blocks (```),
// bold text (**), and links []().
func detectMarkdown(trimmed string, lines []string) bool {
	if len(lines) == 0 {
		return false
	}

	// Check for common markdown patterns
	return strings.HasPrefix(lines[0], "#") ||
		strings.Contains(trimmed, "```") ||
		strings.Contains(trimmed, "**") ||
		strings.Contains(trimmed, "~~") ||
		(strings.Contains(trimmed, "[") && strings.Contains(trimmed, "]("))
}

// detectRequirements checks if the content appears to be a Python requirements.txt file.
// It looks for package names with version specifiers like ==, >=, <=, or ~=.
func detectRequirements(trimmed string, lines []string) bool {
	if !strings.Contains(trimmed, "==") && !strings.Contains(trimmed, ">=") &&
		!strings.Contains(trimmed, "<=") && !strings.Contains(trimmed, "~=") {
		return false
	}

	// Check if it looks like Python packages
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			if strings.ContainsAny(line, "abcdefghijklmnopqrstuvwxyz") {
				return true
			}
		}
	}

	return false
}

// detectDataFormats attempts to detect data exchange and documentation formats.
// It checks for CSV, Markdown, and Requirements.txt formats.
// Returns FormatUnknown if no data format is detected.
func detectDataFormats(trimmed string, lines []string) Format {
	// Check CSV
	if detectCSV(trimmed, lines) {
		return FormatCSV
	}

	// Check Markdown
	if detectMarkdown(trimmed, lines) {
		return FormatMarkdown
	}

	// Check Requirements.txt
	if detectRequirements(trimmed, lines) {
		return FormatRequirements
	}

	return FormatUnknown
}

// isXML checks if the content appears to be XML format.
// It looks for XML declaration (<?xml) or angle bracket tags.
func isXML(trimmed string) bool {
	return strings.HasPrefix(trimmed, "<?xml") ||
		(strings.HasPrefix(trimmed, "<") && strings.Contains(trimmed, ">"))
}

// isINI checks if the content appears to be INI configuration format.
// It looks for section headers in square brackets [section].
func isINI(trimmed string, lines []string) bool {
	if !strings.Contains(trimmed, "[") || !strings.Contains(trimmed, "]") {
		return false
	}

	// Check for INI section pattern
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			return true
		}
	}

	return false
}

// isYAML checks if the content appears to be YAML format.
// It looks for YAML document markers (---) or key: value patterns.
// Excludes content that looks like URLs or Protobuf.
func isYAML(trimmed string) bool {
	if strings.Contains(trimmed, "---") {
		return true
	}

	if !strings.Contains(trimmed, ":") || strings.Contains(trimmed, "://") {
		return false
	}

	// Additional check to ensure it's not a URL or Protobuf
	if strings.Contains(trimmed, "type_url:") || strings.Contains(trimmed, "value:") {
		return false
	}

	return strings.Contains(trimmed, ": ") || strings.HasSuffix(trimmed, ":")
}

// isTOML checks if the content appears to be TOML format.
// It looks for key = value patterns while excluding JSON and XML.
func isTOML(trimmed string) bool {
	if !strings.Contains(trimmed, "=") || strings.Contains(trimmed, ":") {
		return false
	}

	// Make sure it's not JSON or XML
	return !strings.HasPrefix(trimmed, "{") &&
		!strings.HasPrefix(trimmed, "[") &&
		!strings.HasPrefix(trimmed, "<")
}

// detectConfigFormats attempts to detect configuration file formats.
// It checks for XML, INI, YAML, and TOML formats in order of specificity.
// Returns FormatUnknown if no config format is detected.
func detectConfigFormats(trimmed string, lines []string) Format {
	// Check XML
	if isXML(trimmed) {
		return FormatXML
	}

	// Check INI (after CSV to avoid confusion)
	if isINI(trimmed, lines) {
		return FormatINI
	}

	// Check YAML
	if isYAML(trimmed) {
		return FormatYAML
	}

	// Check TOML - simple key=value pattern
	if isTOML(trimmed) {
		return FormatTOML
	}

	return FormatUnknown
}

// DetectFormat attempts to detect the data format by analyzing the content.
// Uses simple heuristics to identify various data formats.
//
// Detection rules:
//   - JSON: Starts with '{' or '['
//   - XML: Starts with '<?xml' or '<'
//   - YAML: Contains '---' or has key:value pattern
//   - TOML: Contains '[section]' pattern with key=value pairs
//   - CSV: Has comma-separated values with consistent columns
//   - GraphQL: Contains query/mutation/type/schema keywords
//   - INI: Has [section] headers or key=value pairs
//   - Dockerfile: Starts with FROM instruction
//   - Markdown: Contains markdown syntax like #, *, -, ```
//   - Requirements.txt: Contains package names with version specifiers
//
// Returns FormatUnknown if the format cannot be determined.
func DetectFormat(data []byte) Format {
	trimmed := strings.TrimSpace(string(data))
	if len(trimmed) == 0 {
		return FormatUnknown
	}

	// Split into lines for multi-line format detection
	lines := strings.Split(trimmed, "\n")

	// Sequential detection for now (parallel overhead not worth it for simple string checks)
	// Try detection in order of specificity
	// Check JSON family first as they have distinct patterns
	if format := detectJSONFamily(trimmed, lines); format != FormatUnknown {
		return format
	}

	// Check developer formats
	if format := detectDeveloperFormats(trimmed, lines); format != FormatUnknown {
		return format
	}

	// Check data formats
	if format := detectDataFormats(trimmed, lines); format != FormatUnknown {
		return format
	}

	// Check config formats last as they're more general
	if format := detectConfigFormats(trimmed, lines); format != FormatUnknown {
		return format
	}

	return FormatUnknown
}

// extensionMap maps file extensions to formats
var extensionMap = map[string]Format{
	"json":          FormatJSON,
	"yaml":          FormatYAML,
	"yml":           FormatYAML,
	"xml":           FormatXML,
	"toml":          FormatTOML,
	"csv":           FormatCSV,
	"graphql":       FormatGraphQL,
	"gql":           FormatGraphQL,
	"ini":           FormatINI,
	"cfg":           FormatINI,
	"conf":          FormatINI,
	"hcl":           FormatHCL,
	"tf":            FormatHCL,
	"tfvars":        FormatHCL,
	"pb":            FormatProtobuf,
	"proto":         FormatProtobuf,
	"textproto":     FormatProtobuf,
	"pbtxt":         FormatProtobuf,
	"md":            FormatMarkdown,
	"markdown":      FormatMarkdown,
	"mkd":           FormatMarkdown,
	"mdwn":          FormatMarkdown,
	"mdown":         FormatMarkdown,
	"mdtxt":         FormatMarkdown,
	"mdtext":        FormatMarkdown,
	"jsonl":         FormatJSONL,
	"ndjson":        FormatJSONL,
	"jsonlines":     FormatJSONL,
	"ipynb":         FormatJupyter,
	"dockerfile":    FormatDockerfile,
	"containerfile": FormatDockerfile,
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
	// Check for Dockerfile without extension first
	baseName := strings.ToLower(filename[strings.LastIndex(filename, "/")+1:])
	const dockerfileName = "dockerfile"
	if baseName == dockerfileName || strings.HasPrefix(baseName, dockerfileName+".") {
		return FormatDockerfile
	}

	lastDot := strings.LastIndex(filename, ".")
	if lastDot == -1 {
		return FormatUnknown
	}
	ext := strings.ToLower(strings.TrimPrefix(filename[lastDot:], "."))

	// Special case for txt files
	if ext == "txt" && strings.Contains(strings.ToLower(filename), "requirements") {
		return FormatRequirements
	}

	if format, ok := extensionMap[ext]; ok {
		return format
	}

	return FormatUnknown
}

// errorString is a helper function that safely converts an error to string.
// Returns an empty string if the error is nil, otherwise returns err.Error().
//
// This is used internally to populate the Error field in Result structs.
func errorString(err error) string {
	if err == nil {
		return ""
	}

	return err.Error()
}
