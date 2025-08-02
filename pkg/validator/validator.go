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

	// Check for Jupyter Notebook first (it's a specific type of JSON)
	if strings.HasPrefix(trimmed, "{") && strings.Contains(trimmed, "\"cells\"") &&
		strings.Contains(trimmed, "\"metadata\"") && strings.Contains(trimmed, "\"nbformat\"") {
		return FormatJupyter
	}

	// Check for JSON Lines before regular JSON
	if len(lines) > 1 {
		validJSONLines := 0
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" {
				if (strings.HasPrefix(line, "{") && strings.HasSuffix(line, "}")) ||
					(strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]")) {
					validJSONLines++
				} else {
					validJSONLines = 0

					break
				}
			}
		}
		if validJSONLines > 1 {
			return FormatJSONL
		}
	}

	// Check JSON (after Jupyter and JSONL)
	if trimmed[0] == '{' || trimmed[0] == '[' {
		// Verify it's valid JSON structure
		if strings.HasPrefix(trimmed, "{") && strings.HasSuffix(trimmed, "}") {
			return FormatJSON
		}
		if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
			return FormatJSON
		}
	}

	// Check XML
	if strings.HasPrefix(trimmed, "<?xml") || (strings.HasPrefix(trimmed, "<") && strings.Contains(trimmed, ">")) {
		return FormatXML
	}

	// Check Dockerfile - look for common Docker instructions
	upperTrimmed := strings.ToUpper(trimmed)
	dockerKeywords := []string{
		"FROM ", "RUN ", "CMD ", "EXPOSE ", "ENV ", "ADD ", "COPY ",
		"ENTRYPOINT ", "VOLUME ", "USER ", "WORKDIR ", "ARG ",
	}
	dockerScore := 0
	for _, keyword := range dockerKeywords {
		if strings.Contains(upperTrimmed, keyword) {
			dockerScore++
		}
	}
	// If we find FROM or multiple Docker keywords, it's likely a Dockerfile
	if strings.Contains(upperTrimmed, "FROM ") || dockerScore >= 3 {
		return FormatDockerfile
	}

	// Check HCL/Terraform before GraphQL (HCL is more specific)
	hclPatterns := []string{
		"resource ", "variable ", "provider ", "module ",
		"output ", "locals ", "terraform ", "data ",
	}
	hclScore := 0
	for _, pattern := range hclPatterns {
		if strings.Contains(trimmed, pattern) {
			hclScore++
		}
	}
	// HCL uses quotes and braces with equals signs
	if hclScore > 0 && strings.Contains(trimmed, "=") &&
		strings.Contains(trimmed, "\"") && strings.Contains(trimmed, "{") {
		return FormatHCL
	}

	// Check GraphQL - look for GraphQL-specific patterns
	graphqlPatterns := []string{
		"query ", "mutation ", "subscription ", "fragment ", "type ",
		"interface ", "enum ", "input ", "scalar ", "schema ",
	}
	graphqlScore := 0
	for _, pattern := range graphqlPatterns {
		if strings.Contains(trimmed, pattern) {
			graphqlScore++
		}
	}
	// Check for GraphQL-specific syntax
	if graphqlScore > 0 && strings.Contains(trimmed, "{") && strings.Contains(trimmed, "}") {
		// Additional check for field selection syntax
		if strings.Contains(trimmed, "{\n") || strings.Contains(trimmed, "{ ") || graphqlScore >= 2 {
			return FormatGraphQL
		}
	}

	// Check Markdown
	if len(lines) > 0 {
		// Check for common markdown patterns
		if strings.HasPrefix(lines[0], "#") || strings.Contains(trimmed, "```") ||
			strings.Contains(trimmed, "**") || strings.Contains(trimmed, "~~") ||
			(strings.Contains(trimmed, "[") && strings.Contains(trimmed, "](")) {
			return FormatMarkdown
		}
	}

	// Check Requirements.txt
	if strings.Contains(trimmed, "==") || strings.Contains(trimmed, ">=") ||
		strings.Contains(trimmed, "<=") || strings.Contains(trimmed, "~=") {
		// Check if it looks like Python packages
		hasPackagePattern := false
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" && !strings.HasPrefix(line, "#") {
				// Simple check for package-like names
				if strings.ContainsAny(line, "abcdefghijklmnopqrstuvwxyz") {
					hasPackagePattern = true

					break
				}
			}
		}
		if hasPackagePattern {
			return FormatRequirements
		}
	}

	// Check CSV
	if strings.Contains(trimmed, ",") {
		// Simple CSV detection - check if lines have consistent comma counts
		if len(lines) > 1 {
			firstLineCommas := strings.Count(lines[0], ",")
			isCSV := true
			for i := 1; i < len(lines) && i < 5; i++ {
				if lines[i] != "" && strings.Count(lines[i], ",") != firstLineCommas {
					isCSV = false

					break
				}
			}
			if isCSV && firstLineCommas > 0 {
				return FormatCSV
			}
		}
	}

	// Check INI (after CSV to avoid confusion)
	if strings.Contains(trimmed, "[") && strings.Contains(trimmed, "]") {
		// Check for INI section pattern
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
				return FormatINI
			}
		}
	}

	// Check for Protobuf text format before YAML (more specific)
	if (strings.Contains(trimmed, "type_url:") || strings.Contains(trimmed, "value:")) &&
		strings.Contains(trimmed, "\"") {
		return FormatProtobuf
	}

	// Check YAML
	if strings.Contains(trimmed, "---") || (strings.Contains(trimmed, ":") && !strings.Contains(trimmed, "://")) {
		// Additional check to ensure it's not a URL or Protobuf
		if !strings.Contains(trimmed, "type_url:") && !strings.Contains(trimmed, "value:") {
			if strings.Contains(trimmed, ": ") || strings.HasSuffix(trimmed, ":") {
				return FormatYAML
			}
		}
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

// errorString returns empty string if error is nil
func errorString(err error) string {
	if err == nil {
		return ""
	}

	return err.Error()
}
