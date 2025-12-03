### Build the wrapper
```
cd openai-review-mcp
go build -o codex-consultant
```

### Configure Claude Code

Edit `~/.config/claude/claude_code_config.json`
```
{
  "mcpServers": {
    "codex-consultant": {
      "command": "/path/to/codex-consultant"
    }
  }
}
```

### Usage in Claude Code

"Can you use the ask_codex tool to get OpenAI's perspective on this rate limiter implementation?"

"Please have Codex review the authentication code I just wrote"
