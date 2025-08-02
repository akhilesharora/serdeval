// Command serdeval provides a CLI for validating JSON, YAML, XML, and TOML files.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/akhilesharora/serdeval"
)

var (
	green = color.New(color.FgGreen)
	red   = color.New(color.FgRed)
	cyan  = color.New(color.FgCyan)

	// Version is set at build time via -ldflags
	Version = "dev"
)

type ValidationResult struct {
	Valid    bool   `json:"valid"`
	Format   string `json:"format"`
	Error    string `json:"error,omitempty"`
	FileName string `json:"filename,omitempty"`
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "serdeval",
		Short: "Privacy-focused data format validator for JSON, YAML, XML, and TOML",
		Long: `SerdeVal is a completely offline CLI tool that validates common data formats.
		
PRIVACY GUARANTEE:
• No data logging, tracking, or retention
• No network connections
• No clipboard access
• All validation happens locally
• Your data never leaves your machine`,
	}

	var validateCmd = &cobra.Command{
		Use:   "validate [files...]",
		Short: "Validate data format files",
		Args:  cobra.MinimumNArgs(0),
		Run:   validateFiles,
	}

	var webCmd = &cobra.Command{
		Use:   "web",
		Short: "Start web interface",
		Long:  "Start a local web server with a user-friendly interface for validation and formatting",
		Run:   startWebServer,
	}

	var formatFlag string
	var quietFlag bool
	var jsonOutputFlag bool
	var portFlag int

	validateCmd.Flags().StringVarP(&formatFlag, "format", "f", "auto", "Format to validate (json, yaml, xml, toml, auto)")
	validateCmd.Flags().BoolVarP(&quietFlag, "quiet", "q", false, "Only show errors")
	validateCmd.Flags().BoolVarP(&jsonOutputFlag, "json", "j", false, "Output results as JSON")

	webCmd.Flags().IntVarP(&portFlag, "port", "p", 8080, "Port to serve web interface on")

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("SerdeVal %s\n", Version)
		},
	}

	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(webCmd)
	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func validateFiles(cmd *cobra.Command, args []string) {
	format, _ := cmd.Flags().GetString("format")
	quiet, _ := cmd.Flags().GetBool("quiet")
	jsonOutput, _ := cmd.Flags().GetBool("json")

	var results []ValidationResult

	if len(args) == 0 {
		result := validateStdin(format)
		results = append(results, result)
	} else {
		for _, arg := range args {
			fileResults := validatePath(arg, format)
			results = append(results, fileResults...)
		}
	}

	if jsonOutput {
		output, _ := json.MarshalIndent(results, "", "  ")
		fmt.Println(string(output))

		return
	}

	exitCode := 0
	for _, result := range results {
		if !result.Valid {
			exitCode = 1
		}
		printResult(result, quiet)
	}

	os.Exit(exitCode)
}

func validatePath(path, format string) []ValidationResult {
	var results []ValidationResult

	info, err := os.Stat(path)
	if err != nil {
		results = append(results, ValidationResult{
			Valid:    false,
			Format:   "unknown",
			Error:    fmt.Sprintf("Cannot access file: %v", err),
			FileName: path,
		})

		return results
	}

	if info.IsDir() {
		err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && isValidatableFile(filePath, format) {
				result := validateFile(filePath, format)
				results = append(results, result)
			}

			return nil
		})
		if err != nil {
			results = append(results, ValidationResult{
				Valid:    false,
				Format:   "unknown",
				Error:    fmt.Sprintf("Error walking directory: %v", err),
				FileName: path,
			})
		}
	} else {
		result := validateFile(path, format)
		results = append(results, result)
	}

	return results
}

func validateFile(filename, format string) ValidationResult {
	data, err := os.ReadFile(filename) // #nosec G304 - CLI tool needs to read user-specified files
	if err != nil {
		return ValidationResult{
			Valid:    false,
			Format:   "unknown",
			Error:    fmt.Sprintf("Cannot read file: %v", err),
			FileName: filename,
		}
	}

	return validateData(data, filename, format)
}

func validateStdin(format string) ValidationResult {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return ValidationResult{
			Valid:    false,
			Format:   "unknown",
			Error:    fmt.Sprintf("Cannot read stdin: %v", err),
			FileName: "stdin",
		}
	}

	return validateData(data, "stdin", format)
}

func validateData(data []byte, filename, format string) ValidationResult {
	var result serdeval.Result

	const autoFormat = "auto"

	if format == autoFormat {
		// Try filename first, then content
		detectedFormat := serdeval.DetectFormatFromFilename(filename)
		if detectedFormat != serdeval.FormatUnknown {
			v, _ := serdeval.NewValidator(detectedFormat)
			result = v.Validate(data)
		} else {
			result = serdeval.ValidateAuto(data)
		}
	} else {
		// Try to create validator for the specified format
		var formatType serdeval.Format
		switch format {
		case "json":
			formatType = serdeval.FormatJSON
		case "yaml":
			formatType = serdeval.FormatYAML
		case "xml":
			formatType = serdeval.FormatXML
		case "toml":
			formatType = serdeval.FormatTOML
		default:
			return ValidationResult{
				Valid:    false,
				Format:   format,
				Error:    "unsupported format",
				FileName: filename,
			}
		}

		v, err := serdeval.NewValidator(formatType)
		if err != nil {
			return ValidationResult{
				Valid:    false,
				Format:   format,
				Error:    err.Error(),
				FileName: filename,
			}
		}
		result = v.Validate(data)
	}

	return ValidationResult{
		Valid:    result.Valid,
		Format:   string(result.Format),
		Error:    result.Error,
		FileName: filename,
	}
}

func isValidatableFile(filename, format string) bool {
	const autoFormat = "auto"
	if format != autoFormat {
		return true
	}

	ext := strings.ToLower(filepath.Ext(filename))
	validExts := []string{".json", ".yaml", ".yml", ".xml", ".toml"}

	for _, validExt := range validExts {
		if ext == validExt {
			return true
		}
	}

	return false
}

func printResult(result ValidationResult, quiet bool) {
	if result.Valid {
		if !quiet {
			_, _ = green.Printf("✓ %s: Valid %s\n", result.FileName, result.Format)
		}
	} else {
		_, _ = red.Printf("✗ %s: Invalid %s", result.FileName, result.Format)
		if result.Error != "" {
			fmt.Printf(" - %s", result.Error)
		}
		fmt.Println()
	}
}

func startWebServer(cmd *cobra.Command, args []string) {
	port, _ := cmd.Flags().GetInt("port")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/static/index.html")
	})

	http.HandleFunc("/api/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"version": Version})
	})

	_, _ = cyan.Printf("🌐 SerdeVal web interface starting on http://localhost:%d\n", port)
	_, _ = cyan.Printf("🔒 Privacy-first: All validation happens in your browser\n")
	fmt.Printf("Press Ctrl+C to stop\n\n")

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      nil,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		_, _ = red.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}
}
