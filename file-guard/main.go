package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var version = "0.0.1"

type ToolInput struct {
	FilePath string `json:"file_path"`
	Path     string `json:"path"`
}

type ToolArgs struct {
	ToolInput ToolInput `json:"tool_input"`
}

type AppConfig struct {
	ProtectedFiles []string
}

var DefaultProtectedFilesCSV = "/etc/passwd,/etc/shadow,.env"

func NewAppConfig() *AppConfig {
	return &AppConfig{
		ProtectedFiles: make([]string, 0),
	}
}

func main() {
	config := NewAppConfig()

	protectedFilesCSV := ""
	separator := ","
	showVersion := false
	flag.StringVar(&protectedFilesCSV, "protect", "/etc/passwd,/etc/shadow,.env", "List of protected file paths")
	flag.StringVar(&separator, "separator", ",", "Separator for protected file paths")
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.Parse()

	if showVersion {
		fmt.Printf("FileGuard Version: %s\n", version)
		fmt.Println("  Claude Code Hook to protect sensitive files from being accessed.")
		os.Exit(0)
	}

	config.ProtectedFiles = append(config.ProtectedFiles, splitAndTrim(protectedFilesCSV, separator)...)
	if len(config.ProtectedFiles) == 0 {
		config.ProtectedFiles = append(config.ProtectedFiles, splitAndTrim(DefaultProtectedFilesCSV, separator)...)
	}

	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
		os.Exit(1)
	}

	var args ToolArgs
	if err := json.Unmarshal(data, &args); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSON input: %v\n", err)
		os.Exit(1)
	}

	// Check if Claude is trying to read the protected file
	if args.ToolInput.FilePath != "" {
		if isProtectedFile(config.ProtectedFiles, args.ToolInput.FilePath) {
			fmt.Fprintf(os.Stderr, "Access to protected file '%s' is denied.\n", args.ToolInput.FilePath)
			os.Exit(1)
		}
	}

	if args.ToolInput.Path != "" {
		if isProtectedFile(config.ProtectedFiles, args.ToolInput.Path) {
			fmt.Fprintf(os.Stderr, "Access to protected path '%s' is denied.\n", args.ToolInput.Path)
			os.Exit(1)
		}
	}

	// Proceed with normal processing
	fmt.Println("Input is valid and does not access protected files.")
}

func isProtectedFile(protectedFiles []string, filePath string) bool {
	for _, protectedFile := range protectedFiles {
		if strings.Contains(filePath, protectedFile) {
			return true
		}
	}

	return false
}

func splitAndTrim(csv, separator string) []string {
	var result []string
	for _, item := range strings.Split(csv, string(separator)) {
		result = append(result, strings.TrimSpace(item))
	}

	return result
}
