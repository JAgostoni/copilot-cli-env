package profile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUpdateProfile(t *testing.T) {
	tmpDir := t.TempDir()
	profilePath := filepath.Join(tmpDir, ".testrc")

	// Test 1: Inject block into an empty/non-existent file
	newContent := "export FOO=bar"
	err := UpdateProfile(profilePath, newContent)
	if err != nil {
		t.Fatalf("Failed to update empty profile: %v", err)
	}

	content, err := os.ReadFile(profilePath)
	if err != nil {
		t.Fatalf("Failed to read profile: %v", err)
	}
	contentStr := string(content)
	if !strings.Contains(contentStr, BlockStart) || !strings.Contains(contentStr, BlockEnd) || !strings.Contains(contentStr, newContent) {
		t.Errorf("Block not injected correctly. Got:\n%s", contentStr)
	}

	// Test 2: Update existing block
	updatedContent := "export FOO=baz"
	err = UpdateProfile(profilePath, updatedContent)
	if err != nil {
		t.Fatalf("Failed to update existing profile block: %v", err)
	}

	content, err = os.ReadFile(profilePath)
	if err != nil {
		t.Fatalf("Failed to read profile: %v", err)
	}
	contentStr = string(content)
	if strings.Contains(contentStr, newContent) {
		t.Errorf("Old content still present. Got:\n%s", contentStr)
	}
	if !strings.Contains(contentStr, updatedContent) {
		t.Errorf("New content not injected correctly. Got:\n%s", contentStr)
	}

	// Test 3: Append to a file without an existing block
	otherConfig := "alias ll='ls -al'\n"
	profilePath2 := filepath.Join(tmpDir, ".testrc2")
	os.WriteFile(profilePath2, []byte(otherConfig), 0644)

	err = UpdateProfile(profilePath2, newContent)
	if err != nil {
		t.Fatalf("Failed to append to profile: %v", err)
	}

	content2, err := os.ReadFile(profilePath2)
	if err != nil {
		t.Fatalf("Failed to read profile 2: %v", err)
	}
	contentStr2 := string(content2)
	if !strings.HasPrefix(contentStr2, otherConfig) {
		t.Errorf("Original content lost or corrupted. Got:\n%s", contentStr2)
	}
	if !strings.Contains(contentStr2, BlockStart) {
		t.Errorf("Block not appended correctly. Got:\n%s", contentStr2)
	}

	// Test 4: Update a file where the block is in the middle
	mixedContent := "alias ll='ls -al'\n" + BlockStart + "\nexport OLD=config\n" + BlockEnd + "\necho 'hello'\n"
	profilePath3 := filepath.Join(tmpDir, ".testrc3")
	os.WriteFile(profilePath3, []byte(mixedContent), 0644)

	err = UpdateProfile(profilePath3, newContent)
	if err != nil {
		t.Fatalf("Failed to update middle block: %v", err)
	}

	content3, err := os.ReadFile(profilePath3)
	if err != nil {
		t.Fatalf("Failed to read profile 3: %v", err)
	}
	contentStr3 := string(content3)
	if !strings.Contains(contentStr3, "alias ll='ls -al'") || !strings.Contains(contentStr3, "echo 'hello'") {
		t.Errorf("Surrounding lines corrupted. Got:\n%s", contentStr3)
	}
	if !strings.Contains(contentStr3, newContent) {
		t.Errorf("New content not injected correctly. Got:\n%s", contentStr3)
	}
}

func TestGeneratePreview(t *testing.T) {
	tmpDir := t.TempDir()
	profilePath := filepath.Join(tmpDir, ".testrc")

	// Test preview for new file
	newContent := "export PREVIEW=1"
	preview, err := GeneratePreview(profilePath, newContent)
	if err != nil {
		t.Fatalf("GeneratePreview failed: %v", err)
	}
	if !strings.Contains(preview, BlockStart) || !strings.Contains(preview, newContent) {
		t.Errorf("Preview incorrect for new file. Got:\n%s", preview)
	}
}
