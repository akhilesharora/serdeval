package serdeval

import (
	"testing"
)

func TestNewValidator(t *testing.T) {
	tests := []struct {
		format  Format
		wantErr bool
	}{
		{FormatJSON, false},
		{FormatYAML, false},
		{FormatXML, false},
		{FormatTOML, false},
		{FormatCSV, false},
		{FormatGraphQL, false},
		{FormatINI, false},
		{FormatHCL, false},
		{FormatProtobuf, false},
		{FormatMarkdown, false},
		{FormatJSONL, false},
		{FormatJupyter, false},
		{FormatRequirements, false},
		{FormatDockerfile, false},
		{Format("invalid"), true},
	}

	for _, tt := range tests {
		t.Run(string(tt.format), func(t *testing.T) {
			v, err := NewValidator(tt.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewValidator() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && v.Format() != tt.format {
				t.Errorf("NewValidator() format = %v, want %v", v.Format(), tt.format)
			}
		})
	}
}

func TestJSONValidator(t *testing.T) {
	v := &JSONValidator{baseValidator{format: FormatJSON}}

	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{"valid object", `{"test": true}`, true},
		{"valid array", `[1, 2, 3]`, true},
		{"invalid syntax", `{"test": }`, false},
		{"empty object", `{}`, true},
		{"null", `null`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.ValidateString(tt.input)
			if result.Valid != tt.valid {
				t.Errorf("ValidateString() = %v, want %v", result.Valid, tt.valid)
			}
			if result.Format != FormatJSON {
				t.Errorf("Format = %v, want %v", result.Format, FormatJSON)
			}
		})
	}
}

func TestValidateAuto(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		format Format
		valid  bool
	}{
		{"json object", `{"test": true}`, FormatJSON, true},
		{"yaml", `test: true`, FormatYAML, true},
		{"xml", `<root><test>true</test></root>`, FormatXML, true},
		{"toml", `test = true`, FormatTOML, true},
		{"dockerfile", `FROM ubuntu:20.04
RUN apt-get update
COPY . /app
CMD ["bash"]`, FormatDockerfile, true},
		{"graphql", `query GetUser {
  user(id: "123") {
    name
    email
  }
}`, FormatGraphQL, true},
		{"csv", `name,age,city
John,30,NYC
Jane,25,LA`, FormatCSV, true},
		{"hcl", `resource "aws_instance" "example" {
  ami = "ami-12345"
  instance_type = "t2.micro"
}`, FormatHCL, true},
		{"ini", `[database]
host = localhost
port = 5432`, FormatINI, true},
		{"markdown", `# Title

This is **bold** text.`, FormatMarkdown, true},
		{"requirements", `numpy==1.21.0
pandas>=1.3.0`, FormatRequirements, true},
		{"jsonl", `{"event": "login"}
{"event": "logout"}`, FormatJSONL, true},
		{"jupyter", `{
  "cells": [],
  "metadata": {},
  "nbformat": 4
}`, FormatJupyter, true},
		{"protobuf", `type_url: "type.googleapis.com/example"
value: "test"`, FormatProtobuf, true},
		{"unknown", `random text`, FormatUnknown, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateAuto([]byte(tt.input))
			if result.Valid != tt.valid {
				t.Errorf("ValidateAuto() valid = %v, want %v", result.Valid, tt.valid)
			}
			if result.Format != tt.format {
				t.Errorf("ValidateAuto() format = %v, want %v", result.Format, tt.format)
			}
		})
	}
}

func TestCSVValidator(t *testing.T) {
	v := &CSVValidator{baseValidator{format: FormatCSV}}

	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{"valid csv", "name,age\nJohn,30\nJane,25", true},
		{"valid single column", "name\nJohn\nJane", true},
		{"empty csv", "", true},
		{"inconsistent columns", "name,age\nJohn,30,extra", false},
		{"quoted values", `"name","age"
"John Doe","30"`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.ValidateString(tt.input)
			if result.Valid != tt.valid {
				t.Errorf("ValidateString() = %v, want %v", result.Valid, tt.valid)
			}
			if result.Format != FormatCSV {
				t.Errorf("Format = %v, want %v", result.Format, FormatCSV)
			}
		})
	}
}

func TestGraphQLValidator(t *testing.T) {
	v := &GraphQLValidator{baseValidator{format: FormatGraphQL}}

	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{"valid query", `query { user { name } }`, true},
		{"valid mutation", `mutation { createUser(name: "John") { id } }`, true},
		{"valid schema", `type User { id: ID! name: String }`, true},
		{"invalid syntax", `query { user {`, false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.ValidateString(tt.input)
			if result.Valid != tt.valid {
				t.Errorf("ValidateString() = %v, want %v", result.Valid, tt.valid)
			}
			if result.Format != FormatGraphQL {
				t.Errorf("Format = %v, want %v", result.Format, FormatGraphQL)
			}
		})
	}
}

func TestINIValidator(t *testing.T) {
	v := &INIValidator{baseValidator{format: FormatINI}}

	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{"valid ini", "[section]\nkey = value", true},
		{"no section", "key = value", true},
		{"multiple sections", "[section1]\nkey1 = value1\n[section2]\nkey2 = value2", true},
		{"empty", "", true},
		{"comments", "; comment\n[section]\nkey = value", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.ValidateString(tt.input)
			if result.Valid != tt.valid {
				t.Errorf("ValidateString() = %v, want %v", result.Valid, tt.valid)
			}
			if result.Format != FormatINI {
				t.Errorf("Format = %v, want %v", result.Format, FormatINI)
			}
		})
	}
}

func TestHCLValidator(t *testing.T) {
	v := &HCLValidator{baseValidator{format: FormatHCL}}

	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{"valid hcl", `variable "region" { default = "us-west-2" }`, true},
		{"terraform config", `resource "aws_instance" "example" { ami = "ami-12345" }`, true},
		{"invalid syntax", `resource "test" {`, false},
		{"empty", "", true},
		{"nested blocks", `provider "aws" { region = var.region }`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.ValidateString(tt.input)
			if result.Valid != tt.valid {
				t.Errorf("ValidateString() = %v, want %v, error: %v", result.Valid, tt.valid, result.Error)
			}
			if result.Format != FormatHCL {
				t.Errorf("Format = %v, want %v", result.Format, FormatHCL)
			}
		})
	}
}

func TestProtobufValidator(t *testing.T) {
	v := &ProtobufValidator{baseValidator{format: FormatProtobuf}}

	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{"valid protobuf text", `type_url: "type.googleapis.com/google.protobuf.StringValue"
value: "\n\x05hello"`, true},
		{"empty", "", true}, // Empty is valid protobuf text format
		{"invalid syntax", `type_url: "missing quote`, false},
		{"simple message", `value: "test"`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.ValidateString(tt.input)
			if result.Valid != tt.valid {
				t.Errorf("ValidateString() = %v, want %v, error: %v", result.Valid, tt.valid, result.Error)
			}
			if result.Format != FormatProtobuf {
				t.Errorf("Format = %v, want %v", result.Format, FormatProtobuf)
			}
		})
	}
}

func TestMarkdownValidator(t *testing.T) {
	v := &MarkdownValidator{baseValidator{format: FormatMarkdown}}

	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{"valid markdown", "# Hello\n\nThis is **bold** text.", true},
		{"empty", "", true},
		{"plain text", "Just plain text", true},
		{"code block", "```go\nfunc main() {}\n```", true},
		{"list", "- Item 1\n- Item 2\n  - Nested", true},
		{"link", "[Link](https://example.com)", true},
		{"table", "| Col1 | Col2 |\n|------|------|\n| A    | B    |", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.ValidateString(tt.input)
			if result.Valid != tt.valid {
				t.Errorf("ValidateString() = %v, want %v, error: %v", result.Valid, tt.valid, result.Error)
			}
			if result.Format != FormatMarkdown {
				t.Errorf("Format = %v, want %v", result.Format, FormatMarkdown)
			}
		})
	}
}

func TestJSONLValidator(t *testing.T) {
	v := &JSONLValidator{baseValidator{format: FormatJSONL}}

	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{"valid jsonl", `{"name": "Alice", "age": 30}
{"name": "Bob", "age": 25}
{"name": "Charlie", "age": 35}`, true},
		{"empty", "", true},
		{"single line", `{"key": "value"}`, true},
		{"empty lines between", `{"a": 1}

{"b": 2}`, true},
		{"invalid json on line", `{"valid": true}
{invalid json}
{"valid": true}`, false},
		{"mixed types", `{"type": "object"}
[1, 2, 3]
"string"
123
true
null`, true},
		{"trailing newline", `{"a": 1}
{"b": 2}
`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.ValidateString(tt.input)
			if result.Valid != tt.valid {
				t.Errorf("ValidateString() = %v, want %v, error: %v", result.Valid, tt.valid, result.Error)
			}
			if result.Format != FormatJSONL {
				t.Errorf("Format = %v, want %v", result.Format, FormatJSONL)
			}
		})
	}
}

func TestJupyterValidator(t *testing.T) {
	v := &JupyterValidator{baseValidator{format: FormatJupyter}}

	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{"valid notebook", `{
			"cells": [],
			"metadata": {},
			"nbformat": 4,
			"nbformat_minor": 2
		}`, true},
		{"invalid json", `{invalid}`, false},
		{"missing cells", `{"metadata": {}, "nbformat": 4}`, false},
		{"missing metadata", `{"cells": [], "nbformat": 4}`, false},
		{"missing nbformat", `{"cells": [], "metadata": {}}`, false},
		{"empty", `{}`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.ValidateString(tt.input)
			if result.Valid != tt.valid {
				t.Errorf("ValidateString() = %v, want %v, error: %v", result.Valid, tt.valid, result.Error)
			}
		})
	}
}

func TestRequirementsValidator(t *testing.T) {
	v := &RequirementsValidator{baseValidator{format: FormatRequirements}}

	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{"valid requirements", "numpy==1.21.0\npandas>=1.3.0\nscikit-learn", true},
		{"with comments", "# ML dependencies\nnumpy==1.21.0\n# Data processing\npandas", true},
		{"empty", "", true},
		{"empty lines", "numpy\n\npandas\n\n", true},
		{"complex versions", "torch>=1.9.0,<2.0.0\ntensorflow~=2.8.0", true},
		{"git urls", "git+https://github.com/user/repo.git", true},
		{"invalid line", "===", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.ValidateString(tt.input)
			if result.Valid != tt.valid {
				t.Errorf("ValidateString() = %v, want %v, error: %v", result.Valid, tt.valid, result.Error)
			}
		})
	}
}

func TestDockerfileValidator(t *testing.T) {
	v := &DockerfileValidator{baseValidator{format: FormatDockerfile}}

	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{"valid dockerfile", "FROM ubuntu:20.04\nRUN apt-get update\nCMD [\"bash\"]", true},
		{"multi-stage", "FROM node:14 as builder\nRUN npm install\nFROM node:14-alpine\nCOPY --from=builder /app /app", true},
		{"with comments", "# Base image\nFROM python:3.9\n# Install deps\nRUN pip install numpy", true},
		{"missing FROM", "RUN apt-get update\nCMD [\"bash\"]", false},
		{"invalid instruction", "FROM ubuntu\nINVALID instruction here", false},
		{"empty", "", false},
		{"line continuation", "FROM ubuntu\nRUN apt-get update && \\\n    apt-get install -y python3", true},
		{"all instructions", `FROM ubuntu
ARG VERSION=latest
ENV PATH=/usr/local/bin:$PATH
LABEL maintainer="test@example.com"
RUN apt-get update
COPY . /app
ADD file.tar.gz /tmp/
WORKDIR /app
EXPOSE 8080
VOLUME /data
USER nobody
ENTRYPOINT ["python"]
CMD ["app.py"]
HEALTHCHECK CMD curl -f http://localhost/ || exit 1
STOPSIGNAL SIGTERM
SHELL ["/bin/bash", "-c"]
ONBUILD RUN echo "trigger"`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.ValidateString(tt.input)
			if result.Valid != tt.valid {
				t.Errorf("ValidateString() = %v, want %v, error: %v", result.Valid, tt.valid, result.Error)
			}
		})
	}
}

func TestDetectFormat(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		format Format
	}{
		// Basic formats
		{"json object", `{"test": true}`, FormatJSON},
		{"json array", `[1, 2, 3]`, FormatJSON},
		{"yaml", `key: value`, FormatYAML},
		{"xml", `<?xml version="1.0"?><root></root>`, FormatXML},
		{"xml simple", `<root><item>test</item></root>`, FormatXML},
		{"toml", `key = "value"`, FormatTOML},

		// Dockerfile with various patterns
		{"dockerfile basic", `FROM ubuntu:20.04`, FormatDockerfile},
		{"dockerfile with comments", `# My Dockerfile
FROM alpine:latest
RUN apk add curl`, FormatDockerfile},
		{"dockerfile multi instruction", `WORKDIR /app
COPY . .
RUN npm install
CMD ["node", "app.js"]`, FormatDockerfile},

		// GraphQL
		{"graphql query", `query { users { name } }`, FormatGraphQL},
		{"graphql mutation", `mutation CreateUser {
  createUser(name: "John") {
    id
  }
}`, FormatGraphQL},
		{"graphql schema", `type User {
  id: ID!
  name: String
}`, FormatGraphQL},

		// HCL/Terraform
		{"hcl resource", `resource "aws_instance" "example" {
  ami = "ami-12345"
}`, FormatHCL},
		{"hcl variable", `variable "region" {
  default = "us-west-2"
}`, FormatHCL},

		// CSV
		{"csv simple", `name,age
John,30
Jane,25`, FormatCSV},

		// INI
		{"ini", `[section]
key = value`, FormatINI},

		// Markdown
		{"markdown heading", `# Hello World`, FormatMarkdown},
		{"markdown code block", "```go\nfunc main() {}\n```", FormatMarkdown},
		{"markdown bold", `This is **bold** text`, FormatMarkdown},

		// Requirements.txt
		{"requirements", `numpy==1.21.0
pandas>=1.3.0`, FormatRequirements},

		// JSON Lines
		{"jsonl", `{"a": 1}
{"b": 2}
{"c": 3}`, FormatJSONL},

		// Jupyter
		{"jupyter", `{"cells": [], "metadata": {}, "nbformat": 4}`, FormatJupyter},

		// Protobuf text
		{"protobuf", `type_url: "example.com/Type"
value: "data"`, FormatProtobuf},

		// Unknown
		{"plain text", `just some random text`, FormatUnknown},
		{"empty", ``, FormatUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectFormat([]byte(tt.input))
			if result != tt.format {
				t.Errorf("DetectFormat() = %v, want %v", result, tt.format)
			}
		})
	}
}

func TestDetectFormatFromFilename(t *testing.T) {
	tests := []struct {
		filename string
		format   Format
	}{
		{"test.json", FormatJSON},
		{"test.JSON", FormatJSON},
		{"test.yaml", FormatYAML},
		{"test.yml", FormatYAML},
		{"test.xml", FormatXML},
		{"test.toml", FormatTOML},
		{"test.csv", FormatCSV},
		{"test.CSV", FormatCSV},
		{"test.graphql", FormatGraphQL},
		{"test.gql", FormatGraphQL},
		{"test.ini", FormatINI},
		{"test.cfg", FormatINI},
		{"test.hcl", FormatHCL},
		{"test.tf", FormatHCL},
		{"test.proto", FormatProtobuf},
		{"test.pb", FormatProtobuf},
		{"test.textproto", FormatProtobuf},
		{"test.pbtxt", FormatProtobuf},
		{"test.md", FormatMarkdown},
		{"test.markdown", FormatMarkdown},
		{"test.jsonl", FormatJSONL},
		{"test.ndjson", FormatJSONL},
		{"test.ipynb", FormatJupyter},
		{"notebook.ipynb", FormatJupyter},
		{"requirements.txt", FormatRequirements},
		{"requirements-dev.txt", FormatRequirements},
		{"Dockerfile", FormatDockerfile},
		{"dockerfile", FormatDockerfile},
		{"Dockerfile.prod", FormatDockerfile},
		{"my.dockerfile", FormatDockerfile},
		{"test.txt", FormatUnknown},
		{"test", FormatUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			format := DetectFormatFromFilename(tt.filename)
			if format != tt.format {
				t.Errorf("DetectFormatFromFilename(%s) = %v, want %v", tt.filename, format, tt.format)
			}
		})
	}
}

func TestHelperFunctions(t *testing.T) {
	t.Run("isJSON edge cases", func(t *testing.T) {
		// Test empty string
		if isJSON("") {
			t.Error("isJSON(\"\") should return false")
		}
		// Test string that starts with { but doesn't end with }
		if isJSON("{incomplete") {
			t.Error("isJSON(\"{incomplete\") should return false")
		}
		// Test string that starts with [ but doesn't end with ]
		if isJSON("[incomplete") {
			t.Error("isJSON(\"[incomplete\") should return false")
		}
	})

	t.Run("detectCSV edge cases", func(t *testing.T) {
		// Test no commas
		if detectCSV("no commas here", []string{"no commas here"}) {
			t.Error("detectCSV should return false for text without commas")
		}
		// Test single line
		if detectCSV("a,b,c", []string{"a,b,c"}) {
			t.Error("detectCSV should return false for single line")
		}
		// Test empty first line commas
		if detectCSV("test", []string{"test", "a,b"}) {
			t.Error("detectCSV should return false when first line has no commas")
		}
	})

	t.Run("detectMarkdown edge cases", func(t *testing.T) {
		// Test empty lines
		if detectMarkdown("", []string{}) {
			t.Error("detectMarkdown should return false for empty content")
		}
		// Test no markdown patterns
		if detectMarkdown("plain text", []string{"plain text"}) {
			t.Error("detectMarkdown should return false for plain text")
		}
	})

	t.Run("isINI edge cases", func(t *testing.T) {
		// Test with brackets but not INI format
		if isINI("array[0]", []string{"array[0]"}) {
			t.Error("isINI should return false for non-INI brackets")
		}
		// Test missing closing bracket
		if isINI("[section", []string{"[section"}) {
			t.Error("isINI should return false for unclosed section")
		}
	})

	t.Run("isYAML edge cases", func(t *testing.T) {
		// Test URL with colon
		if isYAML("http://example.com") {
			t.Error("isYAML should return false for URLs")
		}
		// Test protobuf-like content
		if isYAML("type_url: value") {
			t.Error("isYAML should return false for protobuf content")
		}
		// Test no colon
		if isYAML("no colon here") {
			t.Error("isYAML should return false for text without colon")
		}
	})

	t.Run("detectRequirements edge cases", func(t *testing.T) {
		// Test with version specifiers but no package names
		if detectRequirements("==1.0.0", []string{"==1.0.0"}) {
			t.Error("detectRequirements should return false for version without package")
		}
		// Test comment only
		if detectRequirements("# comment", []string{"# comment"}) {
			t.Error("detectRequirements should return false for comments only")
		}
	})

	t.Run("isJSONLines edge cases", func(t *testing.T) {
		// Test mixed valid/invalid lines
		lines := []string{"{\"valid\": true}", "not json", "{\"another\": true}"}
		if isJSONLines(lines) {
			t.Error("isJSONLines should return false when any line is invalid JSON")
		}
	})
}
