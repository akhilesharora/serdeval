package validator

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
