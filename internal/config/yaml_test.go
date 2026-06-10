package config

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestLoadYAMLWithComments(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.yaml")

	content := `# This is a comment
version: 1
# Another comment
theme: auto
`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	doc, err := LoadYAMLWithComments(path)
	if err != nil {
		t.Fatalf("failed to load YAML: %v", err)
	}

	if doc.Kind != yaml.DocumentNode {
		t.Errorf("expected document node, got %d", doc.Kind)
	}
}

func TestFindMappingValue(t *testing.T) {
	content := `version: 1
theme: auto
memory_limit: 256
`

	var doc yaml.Node
	if err := yaml.Unmarshal([]byte(content), &doc); err != nil {
		t.Fatalf("failed to parse YAML: %v", err)
	}

	// Find version
	versionNode := FindMappingValue(doc.Content[0], "version")
	if versionNode == nil {
		t.Fatal("version node not found")
	}
	if versionNode.Value != "1" {
		t.Errorf("expected version '1', got '%s'", versionNode.Value)
	}

	// Find theme
	themeNode := FindMappingValue(doc.Content[0], "theme")
	if themeNode == nil {
		t.Fatal("theme node not found")
	}
	if themeNode.Value != "auto" {
		t.Errorf("expected theme 'auto', got '%s'", themeNode.Value)
	}

	// Find non-existent key
	missingNode := FindMappingValue(doc.Content[0], "missing")
	if missingNode != nil {
		t.Error("missing node should be nil")
	}
}

func TestSetMappingValue(t *testing.T) {
	content := `# Comment
version: 1
theme: auto
`

	var doc yaml.Node
	if err := yaml.Unmarshal([]byte(content), &doc); err != nil {
		t.Fatalf("failed to parse YAML: %v", err)
	}

	// Update existing value
	if err := SetMappingValue(doc.Content[0], "theme", "dark"); err != nil {
		t.Fatalf("failed to set value: %v", err)
	}

	themeNode := FindMappingValue(doc.Content[0], "theme")
	if themeNode == nil {
		t.Fatal("theme node not found")
	}
	if themeNode.Value != "dark" {
		t.Errorf("expected theme 'dark', got '%s'", themeNode.Value)
	}

	// Add new value
	if err := SetMappingValue(doc.Content[0], "new_key", "new_value"); err != nil {
		t.Fatalf("failed to add value: %v", err)
	}

	newNode := FindMappingValue(doc.Content[0], "new_key")
	if newNode == nil {
		t.Fatal("new_key node not found")
	}
	if newNode.Value != "new_value" {
		t.Errorf("expected new_value, got '%s", newNode.Value)
	}

	// Marshal and check comments are preserved
	data, err := yaml.Marshal(&doc)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	output := string(data)
	if !containsString(output, "# Comment") {
		t.Error("comment should be preserved")
	}
	if !containsString(output, "theme: dark") {
		t.Error("theme should be 'dark'")
	}
	if !containsString(output, "new_key: new_value") {
		t.Error("new_key should be present")
	}
}

func TestSaveYAMLWithComments(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.yaml")

	content := `# This is a comment
version: 1
theme: auto
`

	var doc yaml.Node
	if err := yaml.Unmarshal([]byte(content), &doc); err != nil {
		t.Fatalf("failed to parse YAML: %v", err)
	}

	// Modify value
	SetMappingValue(doc.Content[0], "theme", "dark")

	// Save
	if err := SaveYAMLWithComments(path, &doc); err != nil {
		t.Fatalf("failed to save: %v", err)
	}

	// Read back
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read: %v", err)
	}

	output := string(data)
	if !containsString(output, "# This is a comment") {
		t.Error("comment should be preserved")
	}
	if !containsString(output, "theme: dark") {
		t.Error("theme should be 'dark'")
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstr(s, substr))
}

func containsSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
