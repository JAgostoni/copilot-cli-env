package envdetect

import (
	"os"
	"runtime"
)

type Environment struct {
	OS          string // "windows", "darwin", "linux"
	Shell       string // "bash", "zsh", "fish", "powershell", "cmd", "gitbash"
	Terminal    string // "windows-terminal", "alacritty", "vscode", etc.
}

type Detector interface {
	Detect() (Environment, error)
}

type defaultDetector struct{}

// NewDetector creates a new environment detector
func NewDetector() Detector {
	return &defaultDetector{}
}

func (d *defaultDetector) Detect() (Environment, error) {
	env := Environment{
		OS: runtime.GOOS,
	}

	env.Shell = detectShell()

	termProg := os.Getenv("TERM_PROGRAM")
	if termProg == "vscode" {
		env.Terminal = "vscode"
	} else if os.Getenv("WT_SESSION") != "" {
		env.Terminal = "windows-terminal"
	} else if os.Getenv("ALACRITTY_LOG") != "" || os.Getenv("ALACRITTY_WINDOW_ID") != "" {
		env.Terminal = "alacritty"
	}

	return env, nil
}
