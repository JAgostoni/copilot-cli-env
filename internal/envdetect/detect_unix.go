//go:build !windows
// +build !windows

package envdetect

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-ps"
)

func detectShell() string {
	pid := os.Getppid()
	for i := 0; i < 10; i++ {
		p, err := ps.FindProcess(pid)
		if err != nil || p == nil {
			break
		}

		name := p.Executable()
		// Clean up the name (e.g. /bin/zsh -> zsh, -zsh -> zsh)
		name = filepath.Base(name)
		name = strings.TrimPrefix(name, "-")
		if isShell(name) {
			return name
		}

		parentPid := p.PPid()
		if parentPid == 0 || parentPid == pid {
			break
		}
		pid = parentPid
	}

	// Fallback
	shellEnv := os.Getenv("SHELL")
	if shellEnv != "" {
		return filepath.Base(shellEnv)
	}
	return "bash" // safe fallback
}

func isShell(name string) bool {
	switch name {
	case "bash", "zsh", "fish", "sh", "tcsh", "ksh", "dash":
		return true
	}
	return false
}
