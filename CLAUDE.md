# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This repository provides security hooks for Claude Code, specifically the `file-guard` tool that protects sensitive files from being accessed during Claude Code sessions.

## Architecture

### file-guard

A Go-based CLI tool that acts as a PreToolUse hook for Claude Code's Read and Grep operations. It intercepts tool calls and validates that the requested file paths don't match protected patterns.

**How it works:**
1. Claude Code hooks system pipes tool arguments as JSON to stdin
2. `file-guard` parses the JSON to extract `file_path` or `path` fields
3. Checks if the path matches any protected patterns (substring matching)
4. Exits with status 1 (blocking the tool call) if protected, or 0 (allowing it) if safe

**Core components:**
- `main.go:12-19` - Defines the JSON structure for Claude Code tool inputs
- `main.go:75-83` - `isProtectedFile()` uses substring matching to detect protected patterns
- `main.go:85-92` - `splitAndTrim()` parses comma-separated protection patterns

## Development Commands

### Building
```bash
go build -o file-guard ./file-guard
```

### Installing
```bash
go install github.com/timgluz/claude-hooks/file-guard@latest
```

### Testing
```bash
# Run all tests
go test ./file-guard

# Run with verbose output
go test -v ./file-guard

# Run specific test
go test -v ./file-guard -run TestIsProtectedFile
```

### Manual testing
```bash
# Test protected file blocking
echo '{"tool_input":{"file_path":".env"}}' | ./file-guard -protect .env,.secret

# Test allowed file
echo '{"tool_input":{"file_path":"README.md"}}' | ./file-guard -protect .env,.secret
```

## Hook Integration

The tool is designed to integrate with Claude Code's hooks system via `.claude/settings.local.json`:

```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Read|Grep",
        "hooks": [
          {
            "type": "command",
            "command": "file-guard -protect .env,.secret,id_rsa"
          }
        ]
      }
    ]
  }
}
```

The `matcher` field uses regex to target specific Claude Code tools (Read, Grep, etc.).

## Pattern Matching Behavior

Protection patterns use substring matching (not glob or regex):
- `.env` blocks any path containing ".env" (e.g., `.env`, `config/.env.local`, `test.env`)
- `id_rsa` blocks paths containing "id_rsa" (e.g., `~/.ssh/id_rsa`, `backup/id_rsa.pub`)
- Use specific patterns to avoid over-blocking (e.g., `credentials.json` instead of `json`)

## Error Messages

When a file is blocked, file-guard provides clear feedback:
```
BLOCKED: File '/path/to/.env' matches protected pattern '.env'.
To allow access, remove '.env' from the -protect flag in .claude/settings.local.json
```

The error message includes:
- Which file was blocked
- Which pattern triggered the block
- How to modify the configuration to allow access
