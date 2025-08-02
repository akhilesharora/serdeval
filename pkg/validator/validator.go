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

// CSVValidator validates CSV data
type CSVValidator struct {
	baseValidator
}

// GraphQLValidator validates GraphQL queries and schemas
type GraphQLValidator struct {
	baseValidator
}

// INIValidator validates INI configuration data
type INIValidator struct {
	baseValidator
}

// HCLValidator validates HCL configuration data
type HCLValidator struct {
	baseValidator
}

// ProtobufValidator validates Protobuf text format data
type ProtobufValidator struct {
	baseValidator
}

// MarkdownValidator validates Markdown data
type MarkdownValidator struct {
	baseValidator
}

// JSONLValidator validates JSON Lines data
type JSONLValidator struct {
	baseValidator
}

// JupyterValidator validates Jupyter Notebook data
type JupyterValidator struct {
	baseValidator
}

// RequirementsValidator validates Requirements.txt data
type RequirementsValidator struct {
	baseValidator
}

// DockerfileValidator validates Dockerfile data
type DockerfileValidator struct {
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
// Supported formats: FormatJSON, FormatYAML, FormatXML, FormatTOML, FormatCSV, FormatGraphQL,
// FormatINI, FormatHCL, FormatProtobuf, FormatMarkdown, FormatJSONL, FormatJupyter,
// FormatRequirements, FormatDockerfile
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
	case FormatCSV:
		return &CSVValidator{baseValidator{format: FormatCSV}}, nil
	case FormatGraphQL:
		return &GraphQLValidator{baseValidator{format: FormatGraphQL}}, nil
	case FormatINI:
		return &INIValidator{baseValidator{format: FormatINI}}, nil
	case FormatHCL:
		return &HCLValidator{baseValidator{format: FormatHCL}}, nil
	case FormatProtobuf:
		return &ProtobufValidator{baseValidator{format: FormatProtobuf}}, nil
	case FormatMarkdown:
		return &MarkdownValidator{baseValidator{format: FormatMarkdown}}, nil
	case FormatJSONL:
		return &JSONLValidator{baseValidator{format: FormatJSONL}}, nil
	case FormatJupyter:
		return &JupyterValidator{baseValidator{format: FormatJupyter}}, nil
	case FormatRequirements:
		return &RequirementsValidator{baseValidator{format: FormatRequirements}}, nil
	case FormatDockerfile:
		return &DockerfileValidator{baseValidator{format: FormatDockerfile}}, nil
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

// Validate validates CSV data
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

// ValidateString validates CSV string
func (v *CSVValidator) ValidateString(data string) Result {
	return v.Validate([]byte(data))
}

// Validate validates GraphQL queries and schemas
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

// ValidateString validates GraphQL string
func (v *GraphQLValidator) ValidateString(data string) Result {
	return v.Validate([]byte(data))
}

// Validate validates INI configuration data
func (v *INIValidator) Validate(data []byte) Result {
	_, err := ini.Load(data)

	return Result{
		Valid:  err == nil,
		Format: v.format,
		Error:  errorString(err),
	}
}

// ValidateString validates INI string
func (v *INIValidator) ValidateString(data string) Result {
	return v.Validate([]byte(data))
}

// Validate validates HCL configuration data
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

// ValidateString validates HCL string
func (v *HCLValidator) ValidateString(data string) Result {
	return v.Validate([]byte(data))
}

// Validate validates Protobuf text format data
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

// ValidateString validates Protobuf text format string
func (v *ProtobufValidator) ValidateString(data string) Result {
	return v.Validate([]byte(data))
}

// Validate validates Markdown data
func (v *MarkdownValidator) Validate(data []byte) Result {
	md := goldmark.New()
	err := md.Convert(data, io.Discard)

	return Result{
		Valid:  err == nil,
		Format: v.format,
		Error:  errorString(err),
	}
}

// ValidateString validates Markdown string
func (v *MarkdownValidator) ValidateString(data string) Result {
	return v.Validate([]byte(data))
}

// Validate validates JSON Lines data (each line must be valid JSON)
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

// ValidateString validates JSON Lines string
func (v *JSONLValidator) ValidateString(data string) Result {
	return v.Validate([]byte(data))
}

// Validate validates Jupyter Notebook data (must be valid JSON with notebook structure)
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

// ValidateString validates Jupyter Notebook string
func (v *JupyterValidator) ValidateString(data string) Result {
	return v.Validate([]byte(data))
}

// Validate validates Requirements.txt data
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

// ValidateString validates Requirements.txt string
func (v *RequirementsValidator) ValidateString(data string) Result {
	return v.Validate([]byte(data))
}

// Validate validates Dockerfile data
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

// ValidateString validates Dockerfile string
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
	// Check for Dockerfile without extension first
	baseName := strings.ToLower(filename[strings.LastIndex(filename, "/")+1:])
	if baseName == "dockerfile" || strings.HasPrefix(baseName, "dockerfile.") {
		return FormatDockerfile
	}

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
	case "csv":
		return FormatCSV
	case "graphql", "gql":
		return FormatGraphQL
	case "ini", "cfg", "conf":
		return FormatINI
	case "hcl", "tf", "tfvars":
		return FormatHCL
	case "pb", "proto", "textproto", "pbtxt":
		return FormatProtobuf
	case "md", "markdown", "mkd", "mdwn", "mdown", "mdtxt", "mdtext":
		return FormatMarkdown
	case "jsonl", "ndjson", "jsonlines":
		return FormatJSONL
	case "ipynb":
		return FormatJupyter
	case "txt":
		// Check if it might be requirements.txt
		if strings.Contains(strings.ToLower(filename), "requirements") {
			return FormatRequirements
		}
		return FormatUnknown
	case "dockerfile", "containerfile":
		return FormatDockerfile
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
