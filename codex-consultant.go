package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Validate that codex CLI is available
	if err := validateCodexCLI(); err != nil {
		fmt.Fprintf(os.Stderr, "Codex CLI validation failed: %v\n", err)
		fmt.Fprintf(os.Stderr, "Please ensure the 'codex' command is installed and available in PATH\n")
		os.Exit(1)
	}

	s := server.NewMCPServer(
		"Codex Consultant",
		"0.1.0",
		server.WithToolCapabilities(true),
	)

	// Tool: ask_codex
	askCodexTool := mcp.NewTool("ask_codex",
		mcp.WithDescription("Get a second opinion from OpenAI Codex on code, plans, or implementations"),
		mcp.WithString("prompt",
			mcp.Required(),
			mcp.Description("The question or code to ask Codex about"),
		),
		mcp.WithString("context",
			mcp.Description("Additional context or files to include"),
		),
		mcp.WithString("model",
			mcp.Description("Model to use (e.g., gpt-5-codex, gpt-5). Default: gpt-5-codex"),
		),
	)

	s.AddTool(askCodexTool, askCodexHandler)

	// Tool: codex_review
	reviewTool := mcp.NewTool("codex_review",
		mcp.WithDescription("Have OpenAI Codex review code changes or implementation plans"),
		mcp.WithString("target",
			mcp.Required(),
			mcp.Description("What to review (file path, code snippet, or 'current changes')"),
		),
		mcp.WithString("focus",
			mcp.Description("Specific areas to focus on (security, performance, bugs, etc.)"),
		),
	)

	s.AddTool(reviewTool, codexReviewHandler)

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}

func askCodexHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	prompt, err := request.RequireString("prompt")
	if err != nil {
		return mcp.NewToolResultError("prompt is required"), nil
	}

	model := request.GetString("model", "gpt-5-codex")
	contextStr := request.GetString("context", "")

	// Build the full prompt
	fullPrompt := prompt
	if contextStr != "" {
		// Try to read context as a file path
		expandedContext := expandContext(contextStr)
		fullPrompt = fmt.Sprintf("Context: %s\n\nQuestion: %s", expandedContext, prompt)
	}

	// Execute Codex in non-interactive mode
	cmd := exec.CommandContext(ctx, "codex", "exec",
		"--model", model,
		"--",
		fullPrompt,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Codex execution failed: %v\nOutput: %s", err, string(output))), nil
	}

	return mcp.NewToolResultText(string(output)), nil
}

func codexReviewHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	target, err := request.RequireString("target")
	if err != nil {
		return mcp.NewToolResultError("target is required"), nil
	}

	focus := request.GetString("focus", "code quality, bugs, and best practices")

	// Handle different target types
	var reviewContent string
	targetLower := strings.ToLower(strings.TrimSpace(target))

	if targetLower == "current changes" || targetLower == "git diff" {
		// Get current git changes
		cmd := exec.CommandContext(ctx, "git", "diff", "HEAD")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get git diff: %v", err)), nil
		}
		if len(output) == 0 {
			// Try staged changes if no unstaged changes
			cmd = exec.CommandContext(ctx, "git", "diff", "--staged")
			output, err = cmd.CombinedOutput()
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to get staged git diff: %v", err)), nil
			}
		}
		reviewContent = string(output)
		if len(reviewContent) == 0 {
			return mcp.NewToolResultError("No git changes found to review"), nil
		}
	} else if fileExists(target) {
		// Target is a file path
		content, err := os.ReadFile(target)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to read file %s: %v", target, err)), nil
		}
		reviewContent = fmt.Sprintf("File: %s\n\n%s", target, string(content))
	} else {
		// Treat as code snippet
		reviewContent = target
	}

	// Build review prompt
	reviewPrompt := fmt.Sprintf("Please review the following code with focus on: %s. Provide specific, actionable feedback.\n\n%s", focus, reviewContent)

	// Execute Codex with /review command in exec mode
	cmd := exec.CommandContext(ctx, "codex", "exec",
		"--model", "gpt-5-codex",
		"--",
		reviewPrompt,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Codex review failed: %v\nOutput: %s", err, string(output))), nil
	}

	return mcp.NewToolResultText(string(output)), nil
}

// expandContext attempts to read context as a file path or returns it as-is
func expandContext(context string) string {
	// Check if it's a file path
	if fileExists(context) {
		content, err := os.ReadFile(context)
		if err == nil {
			return fmt.Sprintf("File: %s\n\n%s", context, string(content))
		}
	}
	// Return as-is if not a file
	return context
}

// fileExists checks if a file exists and is not a directory
func fileExists(path string) bool {
	// Clean and expand the path
	path = filepath.Clean(path)
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// validateCodexCLI checks if the codex CLI is available
func validateCodexCLI() error {
	cmd := exec.Command("codex", "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("codex command not found or not executable: %w", err)
	}
	return nil
}
