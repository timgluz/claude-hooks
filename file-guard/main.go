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

type HookPayload struct {
	SessionID      string    `json:"session_id"`
	TranscriptPath string    `json:"transcript_path"`
	Cwd            string    `json:"cwd"`
	PermissionMode string    `json:"permission_mode"`
	HookEventName  string    `json:"hook_event_name"`
	ToolName       string    `json:"tool_name"`
	ToolInput      ToolInput `json:"tool_input"`
	ToolUseID      string    `json:"tool_use_id"`
}

type ToolInput struct {
	FilePath string `json:"file_path"`
	Path     string `json:"path"`
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

	var input HookPayload
	if err := json.Unmarshal(data, &input); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSON input: %v\n", err)
		os.Exit(1)
	}

	// Check if Claude is trying to read the protected file
	if input.ToolInput.FilePath != "" {
		if isProtected, pattern := isProtectedFile(config.ProtectedFiles, input.ToolInput.FilePath); isProtected {
			fmt.Fprintf(os.Stderr, "BLOCKED: File '%s' matches protected pattern '%s'.\n", input.ToolInput.FilePath, pattern)
			fmt.Fprintf(os.Stderr, "To allow access, remove '%s' from the -protect flag in .claude/settings.local.json\n", pattern)
			os.Exit(2)
		}
	}

	if input.ToolInput.Path != "" {
		if isProtected, pattern := isProtectedFile(config.ProtectedFiles, input.ToolInput.Path); isProtected {
			fmt.Fprintf(os.Stderr, "BLOCKED: Path '%s' matches protected pattern '%s'.\n", input.ToolInput.Path, pattern)
			fmt.Fprintf(os.Stderr, "To allow access, remove '%s' from the -protect flag in .claude/settings.local.json\n", pattern)
			os.Exit(2)
		}
	}

	// Proceed with normal processing
	fmt.Println("Input is valid and does not access protected files.")
}

func isProtectedFile(protectedFiles []string, filePath string) (bool, string) {
	for _, protectedFile := range protectedFiles {
		if strings.Contains(filePath, protectedFile) {
			return true, protectedFile
		}
	}

	return false, ""
}

func splitAndTrim(csv, separator string) []string {
	var result []string
	for _, item := range strings.Split(csv, string(separator)) {
		result = append(result, strings.TrimSpace(item))
	}

	return result
}
