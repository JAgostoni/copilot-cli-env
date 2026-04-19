package profile

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	BlockStart = "# --- BEGIN COPILOT ENV MANAGER ---"
	BlockEnd   = "# --- END COPILOT ENV MANAGER ---"
)

// ResolveProfilePath attempts to find the correct profile path based on OS and shell
func ResolveProfilePath(osName, shell string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	switch shell {
	case "bash":
		if osName == "darwin" {
			// On macOS, bash uses .bash_profile by default for login shells, but .bashrc is standard.
			// .bashrc is generally safer if properly sourced.
			return filepath.Join(home, ".bashrc"), nil
		}
		return filepath.Join(home, ".bashrc"), nil
	case "gitbash":
		return filepath.Join(home, ".bash_profile"), nil
	case "zsh":
		return filepath.Join(home, ".zshrc"), nil
	case "powershell", "pwsh":
		// Attempt to run pwsh/powershell to get the true $PROFILE path
		cmdName := "powershell"
		if shell == "pwsh" {
			cmdName = "pwsh"
		}
		cmd := exec.Command(cmdName, "-NoProfile", "-Command", "Write-Host -NoNewline $PROFILE")
		out, err := cmd.Output()
		if err == nil && len(out) > 0 {
			return string(out), nil
		}
		// Fallback for Windows
		if osName == "windows" {
			return filepath.Join(home, "Documents", "PowerShell", "Microsoft.PowerShell_profile.ps1"), nil
		}
		// Fallback for Linux/macOS pwsh
		return filepath.Join(home, ".config", "powershell", "Microsoft.PowerShell_profile.ps1"), nil
	}

	return "", fmt.Errorf("unsupported shell for profile management: %s", shell)
}

// UpdateProfile injects or updates the managed block in the given file
func UpdateProfile(filePath, newContent string) error {
	content, err := os.ReadFile(filePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read profile file: %w", err)
	}

	var newFileContent bytes.Buffer
	var inBlock bool
	blockReplaced := false

	if !os.IsNotExist(err) {
		scanner := bufio.NewScanner(bytes.NewReader(content))
		for scanner.Scan() {
			line := scanner.Text()

			if strings.TrimSpace(line) == BlockStart {
				inBlock = true
				blockReplaced = true
				
				if strings.TrimSpace(newContent) != "" {
					newFileContent.WriteString(BlockStart + "\n")
					newFileContent.WriteString(strings.TrimSpace(newContent) + "\n")
					newFileContent.WriteString(BlockEnd + "\n")
				}
				continue
			}

			if strings.TrimSpace(line) == BlockEnd {
				inBlock = false
				continue
			}

			if !inBlock {
				newFileContent.WriteString(line + "\n")
			}
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("error parsing profile file: %w", err)
		}
	}

	// If block was never found, append it
	if !blockReplaced && strings.TrimSpace(newContent) != "" {
		// Ensure there's a trailing newline before appending if the file isn't empty and doesn't end in one
		if newFileContent.Len() > 0 {
			lastChar := newFileContent.Bytes()[newFileContent.Len()-1]
			if lastChar != '\n' {
				newFileContent.WriteString("\n")
			}
		}
		newFileContent.WriteString(BlockStart + "\n")
		newFileContent.WriteString(strings.TrimSpace(newContent) + "\n")
		newFileContent.WriteString(BlockEnd + "\n")
	}

	// Ensure the directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create profile directory: %w", err)
	}

	// Write the file safely
	if err := os.WriteFile(filePath, newFileContent.Bytes(), 0600); err != nil {
		return fmt.Errorf("failed to write profile file: %w", err)
	}

	return nil
}

// GeneratePreview simulates the update and returns the new full file content
func GeneratePreview(filePath, newContent string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return BlockStart + "\n" + strings.TrimSpace(newContent) + "\n" + BlockEnd + "\n", nil
		}
		return "", fmt.Errorf("failed to read profile file: %w", err)
	}

	var newFileContent bytes.Buffer
	var inBlock bool
	blockReplaced := false

	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == BlockStart {
			inBlock = true
			blockReplaced = true
			
			newFileContent.WriteString(BlockStart + "\n")
			newFileContent.WriteString(strings.TrimSpace(newContent) + "\n")
			newFileContent.WriteString(BlockEnd + "\n")
			continue
		}

		if strings.TrimSpace(line) == BlockEnd {
			inBlock = false
			continue
		}

		if !inBlock {
			newFileContent.WriteString(line + "\n")
		}
	}

	if !blockReplaced {
		if newFileContent.Len() > 0 {
			lastChar := newFileContent.Bytes()[newFileContent.Len()-1]
			if lastChar != '\n' {
				newFileContent.WriteString("\n")
			}
		}
		newFileContent.WriteString(BlockStart + "\n")
		newFileContent.WriteString(strings.TrimSpace(newContent) + "\n")
		newFileContent.WriteString(BlockEnd + "\n")
	}

	return newFileContent.String(), nil
}
