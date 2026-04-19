# Session 6: Safe Shell Profile Persistence

## Goal
Implement the logic to safely write configurations directly to shell startup files, using idempotent managed blocks.

## Tooling & Library Versions
- **Go Version:** `1.26.x`
- **Libraries:** Standard library `os`, `path/filepath`, `bufio`.

## Setup Instructions
```text
internal/
  profile/
    manager.go
    manager_test.go
```

## Tasks

### 1. Define the Managed Block Markers
In `internal/profile/manager.go`, define the start and end markers:
```go
const (
	BlockStart = "# --- BEGIN COPILOT ENV MANAGER ---"
	BlockEnd   = "# --- END COPILOT ENV MANAGER ---"
)
```

### 2. Profile File Resolution
Create logic to determine the correct file based on OS and Shell:
- Bash (Linux): `~/.bashrc`
- Zsh (macOS): `~/.zshrc` or `~/.zprofile`
- PowerShell: Read `$PROFILE` via a quick subprocess call (`pwsh -NoProfile -Command "Write-Host $PROFILE"`), or construct the default path `~\Documents\PowerShell\Microsoft.PowerShell_profile.ps1`.

### 3. Implement Idempotent Write Logic
Write a function `UpdateProfile(filePath string, newContent string) error`:
1. Read the file contents.
2. Look for `BlockStart` and `BlockEnd`.
3. If they exist, replace everything between them (and the markers themselves) with the `newContent` (wrapped in markers).
4. If they do not exist, append the markers and `newContent` to the end of the file.
5. Save the file with safe permissions (e.g., `0600`).

### 4. Implement Preview Logic
Write a function `GeneratePreview(...)` that returns a string showing what the file *will* look like, or just the diff of what is being injected.

## Unit Testing
Write unit tests in `manager_test.go`. Use Go's `t.TempDir()` to create temporary mock shell profile files safely.
- Test injecting a block into an empty file.
- Test appending to a file without an existing block.
- Test updating a file that already has a block.
- Test updating a file where the block is in the middle of other user configurations to ensure it doesn't corrupt surrounding lines.
