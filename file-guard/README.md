# File Guard

Protects sensitive files from being accessed by Claude Code.

## Installation

```bash
go install github.com/timgluz/claude-hooks/file-guard@latest
```

## Configuration

Copy the example `settings.local.json` to your project's `.claude/` directory:

```bash
cp settings.local.json /path/to/your/project/.claude/
```

Or add this to your existing `.claude/settings.local.json`:

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

## Usage

Customize the `-protect` flag with comma-separated patterns:

- `.env` - blocks any file containing ".env"
- `*.key` - blocks files ending in .key
- `credentials.json` - blocks this specific filename
