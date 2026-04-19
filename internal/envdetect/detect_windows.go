//go:build windows
// +build windows

package envdetect

import (
	"os"
	"strings"

	"github.com/mitchellh/go-ps"
)

func detectShell() string {
	pid := os.Getppid()
	// Walk up the process tree (limit to 10 to prevent infinite loops)
	for i := 0; i < 10; i++ {
		p, err := ps.FindProcess(pid)
		if err != nil || p == nil {
			break
		}

		name := strings.ToLower(p.Executable())
		switch name {
		case "pwsh.exe", "powershell.exe":
			return "powershell"
		case "cmd.exe":
			return "cmd"
		case "bash.exe", "git-bash.exe":
			return "gitbash"
		}

		// Move to parent process
		parentPid := p.PPid()
		if parentPid == 0 || parentPid == pid {
			break
		}
		pid = parentPid
	}

	// Fallback heuristics
	if os.Getenv("PSModulePath") != "" && os.Getenv("PROMPT") == "" {
		// PROMPT is often set in CMD, so we can use it to distinguish slightly if needed.
		// However, PSModulePath is globally set on modern Windows, making it a tricky fallback.
		return "powershell"
	}

	return "cmd" // safe fallback
}
