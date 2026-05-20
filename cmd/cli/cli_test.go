package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestCLI_Version(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Expected no error, got: %v\nOutput: %s", err, output)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "v0.0.12") {
		t.Error("Expected version output to contain version number")
	}
}

func TestCLI_CMSProducts(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "cms", "products")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Expected no error, got: %v\nOutput: %s", err, output)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "ecs") {
		t.Error("Expected products output to contain 'ecs'")
	}
}

func TestCLI_MissingCredentials(t *testing.T) {
	os.Unsetenv("CLOUD_ACCESS_KEY_ID")
	os.Unsetenv("CLOUD_ACCESS_KEY_SECRET")

	cmd := exec.Command("go", "run", ".", "ecs", "list")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("Expected error for missing credentials")
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "AccessKey") {
		t.Error("Expected error message about AccessKey")
	}
}
