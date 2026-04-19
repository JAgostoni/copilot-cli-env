package envdetect

import (
	"testing"
)

func TestNewDetector(t *testing.T) {
	detector := NewDetector()
	if detector == nil {
		t.Fatal("Expected NewDetector to return a Detector instance")
	}

	env, err := detector.Detect()
	if err != nil {
		t.Fatalf("Detect returned an error: %v", err)
	}

	if env.OS == "" {
		t.Error("Expected OS to be populated")
	}
	if env.Shell == "" {
		t.Error("Expected Shell to be populated")
	}
}
