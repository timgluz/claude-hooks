# File Guard

A security hook for Claude Code that protects sensitive files from being accessed during AI-assisted development sessions.

## What It Does

File Guard intercepts Claude Code's Read and Grep operations and blocks access to files matching protected patterns. This prevents accidentally exposing credentials, private keys, and other sensitive data during development sessions.

## Installation

```bash
go install github.com/timgluz/claude-hooks/file-guard@latest
```

## Configuration

Add this to your project's `.claude/settings.local.json`:

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

## Pattern Matching

File Guard uses **substring matching** (not glob or regex):

- `.env` → blocks any path containing ".env"
  - Matches: `.env`, `config/.env.local`, `test.env`
- `id_rsa` → blocks paths containing "id_rsa"
  - Matches: `~/.ssh/id_rsa`, `backup/id_rsa.pub`
- `credentials.json` → blocks paths containing "credentials.json"
  - Matches: `credentials.json`, `config/credentials.json`

**Tip:** Use specific patterns to avoid over-blocking (e.g., use `credentials.json` instead of `json`)

## Common Protection Patterns

Recommended patterns for different types of sensitive files:

```bash
# Environment files
file-guard -protect .env,.env.local,.env.production

# Credentials and secrets
file-guard -protect .secret,credentials,secrets.yaml,config.json

# SSH and PGP keys
file-guard -protect id_rsa,id_ed25519,.ssh/,*.key,*.pem

# API keys and tokens
file-guard -protect .aws/credentials,.gcp/,.azure/,token

# Database credentials
file-guard -protect database.yml,.pgpass,.my.cnf

# Comprehensive protection
file-guard -protect .env,.secret,id_rsa,credentials,*.key,*.pem,token,.aws/
```

## Error Messages

When a file is blocked, File Guard provides clear, actionable feedback:

```
BLOCKED: File '/home/user/project/.env' matches protected pattern '.env'.
To allow access, remove '.env' from the -protect flag in .claude/settings.local.json
```

The error message includes:
- Which file was blocked
- Which pattern triggered the block
- How to modify the configuration to allow access

## Testing

### Manual Testing

Test that protection works:

```bash
# Should block (exit code 2)
echo '{"tool_input":{"file_path":".env"}}' | file-guard -protect .env,.secret

# Should allow (exit code 0)
echo '{"tool_input":{"file_path":"README.md"}}' | file-guard -protect .env,.secret
```

### Running Unit Tests

```bash
# Run all tests
go test ./file-guard -v

# Run specific test
go test ./file-guard -v -run TestIsProtectedFile
```

## How It Works

1. Claude Code pipes tool arguments as JSON to stdin when a Read or Grep operation occurs
2. File Guard parses the JSON to extract the `file_path` or `path` field
3. Checks if the path contains any protected patterns (substring matching)
4. Exits with status code 2 if protected (blocking the operation)
5. Exits with status code 0 if allowed (operation proceeds)

## Command-Line Options

```bash
file-guard [options]

Options:
  -protect string
        Comma-separated list of protected file patterns (default "/etc/passwd,/etc/shadow,.env")
  -separator string
        Separator for protected file patterns (default ",")
  -version
        Show version information
```

## Examples

### Protect environment files only
```json
{
  "type": "command",
  "command": "file-guard -protect .env,.env.local"
}
```

### Protect SSH keys and credentials
```json
{
  "type": "command",
  "command": "file-guard -protect id_rsa,id_ed25519,credentials.json"
}
```

### Comprehensive protection
```json
{
  "type": "command",
  "command": "file-guard -protect .env,.secret,id_rsa,*.key,*.pem,credentials,token,.aws/,.gcp/"
}
```

## Troubleshooting

**Problem:** File Guard blocks a file I need access to

**Solution:** Remove the matching pattern from the `-protect` flag:
```json
{
  "command": "file-guard -protect .secret,id_rsa"
}
```

**Problem:** Hook isn't triggering

**Solution:**
1. Verify the hook is in `.claude/settings.local.json`
2. Check that `file-guard` is in your PATH: `which file-guard`
3. Test manually: `echo '{"tool_input":{"file_path":".env"}}' | file-guard -protect .env`

**Problem:** Want to use different separators

**Solution:** Use the `-separator` flag:
```json
{
  "command": "file-guard -protect .env;.secret;id_rsa -separator ;"
}
```
