# Codex Consultant MCP Server

An MCP (Model Context Protocol) server that enables Claude Code to consult with OpenAI's Codex for second opinions on code, implementations, and reviews.

## What is This?

This is a bridge between Claude Code and OpenAI's Codex CLI. It provides two powerful tools that Claude Code can use:

- **`ask_codex`** - Get a second opinion from OpenAI Codex on code, plans, or implementations
- **`codex_review`** - Have OpenAI Codex review code changes, files, or snippets

Think of it as getting a peer review from another AI assistant while working with Claude Code.

## Features

- 🔍 **Smart Context Handling** - Automatically reads file contents when file paths are provided
- 📊 **Git Integration** - Review current git changes with "current changes" or "git diff"
- 📁 **File Support** - Pass file paths directly for review or context
- 🎯 **Flexible Reviews** - Review specific files, code snippets, or git changes
- ✅ **CLI Validation** - Checks that Codex CLI is available on startup
- 🔧 **Model Selection** - Choose between different Codex models (gpt-5-codex, gpt-5, etc.)

## Prerequisites

1. **OpenAI Codex CLI** - You must have the `codex` command-line tool installed and available in your PATH
   - Install from: [OpenAI Codex CLI](https://github.com/openai/openai-cli) (or your specific Codex CLI source)
   - Verify installation: `codex --version`

2. **Go 1.21+** - Required to build the server

## Installation

### 1. Build the MCP Server

```bash
cd codex-consultant
go mod tidy
go build -o codex-consultant
```

This creates an executable named `codex-consultant` in your current directory.

### 2. Configure Claude Code

Add the MCP server to your Claude Code configuration file at `~/.config/claude/claude_code_config.json`:

```json
{
  "mcpServers": {
    "codex-consultant": {
      "command": "/absolute/path/to/codex-consultant"
    }
  }
}
```

**Important:** Use the full absolute path to the `codex-consultant` executable.

### 3. Restart Claude Code

After updating the configuration, restart Claude Code to load the new MCP server.

## Usage

Once configured, Claude Code will automatically have access to the Codex tools. You can use them naturally in conversation:

### Ask Codex for Opinions

```
"Can you use ask_codex to get OpenAI's perspective on this rate limiter implementation?"

"Ask Codex what it thinks about using channels vs mutexes here"

"Get a second opinion from Codex on my error handling approach"
```

### Review Code with Codex

```
"Use codex_review to review the current changes"

"Have Codex review the authentication.go file for security issues"

"Ask Codex to review this code snippet for performance problems"
```

## Tool Reference

### `ask_codex`

Get a second opinion from OpenAI Codex on code, plans, or implementations.

**Parameters:**
- `prompt` (required) - The question or code to ask Codex about
- `context` (optional) - Additional context or file path to include
- `model` (optional) - Model to use (default: "gpt-5-codex")

**Examples:**
```
ask_codex(
  prompt="Is this the best way to implement a connection pool?",
  context="database.go"
)

ask_codex(
  prompt="Review this algorithm",
  context="func quickSort(...) { ... }",
  model="gpt-5"
)
```

### `codex_review`

Have OpenAI Codex review code changes, files, or implementation plans.

**Parameters:**
- `target` (required) - What to review:
  - `"current changes"` or `"git diff"` - Reviews uncommitted git changes
  - File path (e.g., `"main.go"`) - Reviews the entire file
  - Code snippet - Reviews the provided code directly
- `focus` (optional) - Specific areas to focus on (default: "code quality, bugs, and best practices")

**Examples:**
```
codex_review(
  target="current changes",
  focus="security vulnerabilities"
)

codex_review(
  target="src/auth/middleware.go",
  focus="performance and error handling"
)

codex_review(
  target="func processPayment() { ... }",
  focus="edge cases and validation"
)
```

## How It Works

1. Claude Code calls the MCP tool (e.g., `ask_codex`)
2. The MCP server processes the request:
   - Validates parameters
   - Reads files if paths are provided
   - Gets git diffs if "current changes" is specified
   - Formats the prompt for Codex
3. Executes `codex exec` with the formatted prompt
4. Returns Codex's response back to Claude Code

## Troubleshooting

### "Codex CLI validation failed"
- Ensure `codex` is installed: `which codex`
- Verify it's executable: `codex --version`
- Check that it's in your PATH

### "Failed to get git diff"
- Ensure you're in a git repository
- Check that git is installed: `git --version`

### "Failed to read file"
- Verify the file path is correct
- Ensure the file exists and is readable
- Use absolute paths or paths relative to Claude Code's working directory

## Development

### Project Structure
```
codex-consultant/
├── codex-consultant.go  # Main server implementation
├── go.mod              # Go module definition
├── go.sum              # Go dependencies
└── README.md           # This file
```

### Building for Development
```bash
go build -o codex-consultant
./codex-consultant  # Test directly (will run as MCP stdio server)
```

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.
