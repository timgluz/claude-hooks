# Claude Hooks

Security hooks for [Claude Code](https://claude.ai/code) that protect sensitive files and enhance safety during AI-assisted development sessions.

## What is this?

This repository provides hook tools that integrate with Claude Code's hooks system to add security guardrails. Hooks are shell commands that execute in response to events like tool calls, allowing you to validate or block operations before they happen.

## Getting Started

### Prerequisites

- [Claude Code CLI](https://claude.ai/code) installed
- Go 1.25.6 or later

### Installation

Install the hooks you need:

```bash
# Install file-guard
go install github.com/timgluz/claude-hooks/file-guard@latest
```

### Configuration

Create or edit `.claude/settings.local.json` in your project directory and add the hooks you want to use. See individual hook documentation below for specific configuration examples.

## Available Hooks

### file-guard

Protects sensitive files from being accessed by Claude Code during Read and Grep operations.

**Features:**
- Blocks access to files matching specified patterns (.env, credentials, SSH keys, etc.)
- Uses substring matching for flexible protection
- Configurable via command-line flags
- Exits with error when protected files are accessed

[Read more →](file-guard/README.md)

**Quick setup:**

```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Read|Grep",
        "hooks": [
          {
            "type": "command",
            "command": "file-guard -protect .env,.secret,id_rsa,credentials.json,*.key,*.pem"
          }
        ]
      }
    ]
  }
}
```

## Development

### Running Tests

```bash
# Test all hooks
go test ./...

# Test specific hook
go test ./file-guard -v
```

### Building from Source

```bash
# Build all hooks
go build ./file-guard

# Install locally
go install ./file-guard
```

## How Hooks Work

Claude Code's hooks system allows you to intercept tool calls:

1. **PreToolUse hooks** - Execute before a tool runs, receiving the tool arguments as JSON via stdin
2. **PostToolUse hooks** - Execute after a tool completes
3. **Matchers** - Regex patterns that determine which tools trigger the hook (e.g., `"Read|Grep"`)

Hooks that exit with a non-zero status code block the operation and display an error to Claude.

## Contributing

Feel free to submit issues or pull requests for new hooks or improvements to existing ones.

## License

See individual hook directories for license information.
