package profile

import (
	"bufio"
	"bytes"
	"os"
	"strings"
)

// UpdateEnvFile updates or creates a .env file, removing existing COPILOT_ vars and appending new ones.
func UpdateEnvFile(filePath, newContent string) error {
	content, err := os.ReadFile(filePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var newFileContent bytes.Buffer
	if !os.IsNotExist(err) {
		scanner := bufio.NewScanner(bytes.NewReader(content))
		for scanner.Scan() {
			line := scanner.Text()
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "COPILOT_") {
				continue // Skip existing copilot vars to ensure clean replacement
			}
			newFileContent.WriteString(line + "\n")
		}
	}

	if newContent != "" {
		// Ensure trailing newline before appending new content
		if newFileContent.Len() > 0 {
			lastChar := newFileContent.Bytes()[newFileContent.Len()-1]
			if lastChar != '\n' {
				newFileContent.WriteString("\n")
			}
		}
		newFileContent.WriteString(strings.TrimSpace(newContent) + "\n")
	}

	return os.WriteFile(filePath, newFileContent.Bytes(), 0600)
}
