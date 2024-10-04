package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
)

var validTypes = []string{"feat", "fix", "test", "chore", "refactor", "docs", "style", "perf", "build", "ci", "revert"}

func main() {
	// Use promptui to select commit type
	prompt := promptui.Select{
		Label: "🔧 Choose commit type",
		Items: validTypes,
	}

	_, commitType, err := prompt.Run()
	if err != nil {
		fmt.Printf("❌ Prompt failed: %v\n", err)
		return
	}

	// Use promptui to get commit scope (optional)
	scopePrompt := promptui.Prompt{
		Label:     "📦 Enter scope (optional)",
		AllowEdit: true, // Allowing empty input as scope is optional
	}
	scope, err := scopePrompt.Run()
	if err != nil {
		fmt.Printf("❌ Prompt failed: %v\n", err)
		return
	}

	// Use promptui to get commit subject with validation in a loop until valid
	subjectPrompt := promptui.Prompt{
		Label: "🖍️  Enter commit subject",
		Validate: func(input string) error {
			if strings.TrimSpace(input) == "" {
				return errors.New("❌ Subject cannot be empty")
			}
			return nil
		},
	}

	subject, err := subjectPrompt.Run()
	if err != nil {
		fmt.Printf("❌ Prompt failed: %v\n", err)
		return
	}

	// Use promptui to get commit desc (optional)
	descPrompt := promptui.Prompt{
		Label:     "📋 Enter description (optional)",
		AllowEdit: true, // Allowing empty input as desc is optional
	}
	desc, err := descPrompt.Run()
	if err != nil {
		fmt.Printf("❌ Prompt failed: %v\n", err)
		return
	}

	// Use promptui to get commit breakingChange (optional)
	breakingChangePrompt := promptui.Prompt{
		Label:     "🔥 Enter Breaking Change (optional)",
		AllowEdit: true, // Allowing empty input as breakingChange is optional
	}
	breakingChange, err := breakingChangePrompt.Run()
	if err != nil {
		fmt.Printf("❌ Prompt failed: %v\n", err)
		return
	}

	hookPrompt := promptui.Select{
		Label: "🚀 Should run git hooks?",
		Items: []string{"No", "Yes"},
	}
	_, hook, _ := hookPrompt.Run()
	var shouldUseHook bool

	if hook != "Yes" && hook != "No" {
		return
	}

	if hook == "Yes" {
		shouldUseHook = true
	}

	// Format the commit message
	commitMessage := formatCommitMessage(commitType, scope, subject, desc, breakingChange)

	runGitCommit(commitMessage, shouldUseHook)
}

// Format the commit message according to Angular convention
func formatCommitMessage(commitType, scope, subject, desc, breakingChange string) string {
	var s string

	if breakingChange != "" {
		commitType = commitType + "!"
	}

	if strings.TrimSpace(scope) != "" {
		s = fmt.Sprintf("%s(%s): %s", commitType, scope, subject)
	} else {
		s = fmt.Sprintf("%s: %s", commitType, subject)
	}

	if desc != "" {
		s += fmt.Sprintf("\n\n%s", desc)
	}

	if breakingChange != "" {
		s += fmt.Sprintf("\n\nBREAKING CHANGE: %s", breakingChange)
	}

	return s
}

func runGitCommit(commitMessage string, shouldUseHook bool) {
	var cmd *exec.Cmd

	if !shouldUseHook {
		cmd = exec.Command("git", "commit", "-m", commitMessage, "--no-verify")
		defer fmt.Printf("git commit -m \"%s\" --no-verify\n", commitMessage)
	} else {
		cmd = exec.Command("git", "commit", "-m", commitMessage)
		defer fmt.Printf("git commit -m \"%s\"\n", commitMessage)
	}

	// Set the standard output and error to the console
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the git commit command
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to execute git commit: %v", err)
	}
}
