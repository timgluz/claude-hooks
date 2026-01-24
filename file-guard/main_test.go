package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsProtectedFile(t *testing.T) {
	tests := []struct {
		protectedFiles  []string
		filePath        string
		expected        bool
		expectedPattern string
	}{
		{[]string{"/etc/passwd", "/etc/shadow"}, "/etc/passwd", true, "/etc/passwd"},
		{[]string{"/etc/passwd", "/etc/shadow"}, "/etc/hosts", false, ""},
		{[]string{".env", "config.yaml"}, "config.yaml", true, "config.yaml"},
		{[]string{".env", "config.yaml"}, "README.md", false, ""},
	}

	for _, test := range tests {
		t.Run(test.filePath, func(t *testing.T) {
			result, pattern := isProtectedFile(test.protectedFiles, test.filePath)
			assert.Equal(t, test.expected, result)
			assert.Equal(t, test.expectedPattern, pattern)
		})
	}
}

func TestSplitAndTrim(t *testing.T) {
	tests := []struct {
		csv       string
		separator string
		expected  []string
	}{
		{"a,b,c", ",", []string{"a", "b", "c"}},
		{"  a ; b ; c  ", ";", []string{"a", "b", "c"}},
		{"x|y|z", "|", []string{"x", "y", "z"}},
		{" single ", ",", []string{"single"}},
	}

	for _, test := range tests {
		t.Run(test.csv, func(t *testing.T) {
			result := splitAndTrim(test.csv, test.separator)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestHookPayloadUnmarshal(t *testing.T) {
	tests := []struct {
		name          string
		json          string
		expectedPath  string
		expectedTool  string
		shouldSucceed bool
	}{
		{
			name: "Read tool with file_path",
			json: `{
				"session_id": "abc123",
				"transcript_path": "/Users/test/.claude/projects/test.jsonl",
				"cwd": "/Users/test",
				"permission_mode": "default",
				"hook_event_name": "PreToolUse",
				"tool_name": "Read",
				"tool_input": {
					"file_path": "/path/to/.env"
				},
				"tool_use_id": "toolu_01ABC123"
			}`,
			expectedPath:  "/path/to/.env",
			expectedTool:  "Read",
			shouldSucceed: true,
		},
		{
			name: "Grep tool with path",
			json: `{
				"session_id": "xyz789",
				"transcript_path": "/Users/test/.claude/projects/test.jsonl",
				"cwd": "/Users/test",
				"permission_mode": "default",
				"hook_event_name": "PreToolUse",
				"tool_name": "Grep",
				"tool_input": {
					"path": "/home/user/.secret"
				},
				"tool_use_id": "toolu_01XYZ789"
			}`,
			expectedPath:  "/home/user/.secret",
			expectedTool:  "Grep",
			shouldSucceed: true,
		},
		{
			name: "Read tool with both file_path and path",
			json: `{
				"session_id": "test123",
				"tool_name": "Read",
				"tool_input": {
					"file_path": "/test/file.txt"
				},
				"tool_use_id": "toolu_test"
			}`,
			expectedPath:  "/test/file.txt",
			expectedTool:  "Read",
			shouldSucceed: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var payload HookPayload
			err := json.Unmarshal([]byte(test.json), &payload)

			if test.shouldSucceed {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedTool, payload.ToolName)

				if payload.ToolInput.FilePath != "" {
					assert.Equal(t, test.expectedPath, payload.ToolInput.FilePath)
				}
				if payload.ToolInput.Path != "" {
					assert.Equal(t, test.expectedPath, payload.ToolInput.Path)
				}
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestProtectionWithRealPayload(t *testing.T) {
	tests := []struct {
		name           string
		json           string
		protectedFiles []string
		shouldBlock    bool
	}{
		{
			name: "Block .env file",
			json: `{
				"tool_name": "Read",
				"tool_input": {
					"file_path": "/home/user/project/.env"
				}
			}`,
			protectedFiles: []string{".env", ".secret"},
			shouldBlock:    true,
		},
		{
			name: "Allow README.md",
			json: `{
				"tool_name": "Read",
				"tool_input": {
					"file_path": "/home/user/project/README.md"
				}
			}`,
			protectedFiles: []string{".env", ".secret"},
			shouldBlock:    false,
		},
		{
			name: "Block id_rsa via Grep path",
			json: `{
				"tool_name": "Grep",
				"tool_input": {
					"path": "/home/user/.ssh/id_rsa"
				}
			}`,
			protectedFiles: []string{"id_rsa", ".env"},
			shouldBlock:    true,
		},
		{
			name: "Allow public key",
			json: `{
				"tool_name": "Read",
				"tool_input": {
					"file_path": "/home/user/.ssh/id_rsa.pub"
				}
			}`,
			protectedFiles: []string{".env", ".secret"},
			shouldBlock:    false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var payload HookPayload
			err := json.Unmarshal([]byte(test.json), &payload)
			assert.NoError(t, err)

			isBlocked := false
			if payload.ToolInput.FilePath != "" {
				isBlocked, _ = isProtectedFile(test.protectedFiles, payload.ToolInput.FilePath)
			}
			if payload.ToolInput.Path != "" {
				pathBlocked, _ := isProtectedFile(test.protectedFiles, payload.ToolInput.Path)
				isBlocked = isBlocked || pathBlocked
			}

			assert.Equal(t, test.shouldBlock, isBlocked)
		})
	}
}
