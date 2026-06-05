package config

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTestConfig(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	SetConfigPath(configPath)
	// Clear master password for tests
	ClearMasterPassword()
	return configPath
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Version != CurrentVersion {
		t.Errorf("expected version %d, got %d", CurrentVersion, cfg.Version)
	}
	if cfg.SaveCredentials != false {
		t.Errorf("expected save_credentials false, got %v", cfg.SaveCredentials)
	}
	if cfg.Theme != "auto" {
		t.Errorf("expected theme 'auto', got '%s'", cfg.Theme)
	}
	if cfg.MemoryLimit != 256 {
		t.Errorf("expected memory limit 256, got %d", cfg.MemoryLimit)
	}
	if cfg.Concurrency != 3 {
		t.Errorf("expected concurrency 3, got %d", cfg.Concurrency)
	}
	if cfg.PasswordPolicy.MinLength != 8 {
		t.Errorf("expected min length 8, got %d", cfg.PasswordPolicy.MinLength)
	}
}

func TestSaveAndLoad(t *testing.T) {
	setupTestConfig(t)

	cfg := DefaultConfig()
	cfg.CurrentProfile = "test"
	cfg.Profiles["test"] = &Profile{
		AccessKeyID: "test-id",
		Region:      "cn-hangzhou",
	}

	// Save
	if err := Save(cfg); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	// Reset cache
	mu.Lock()
	globalConfig = nil
	mu.Unlock()

	// Load
	loaded, err := Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if loaded.CurrentProfile != "test" {
		t.Errorf("expected current profile 'test', got '%s'", loaded.CurrentProfile)
	}
	if loaded.Profiles["test"].AccessKeyID != "test-id" {
		t.Errorf("expected access key id 'test-id', got '%s'", loaded.Profiles["test"].AccessKeyID)
	}
	if loaded.Profiles["test"].Region != "cn-hangzhou" {
		t.Errorf("expected region 'cn-hangzhou', got '%s'", loaded.Profiles["test"].Region)
	}
}

func TestInitConfig(t *testing.T) {
	setupTestConfig(t)

	// Init config
	if err := InitConfig(false); err != nil {
		t.Fatalf("failed to init config: %v", err)
	}

	// Check file exists
	path, _ := GetConfigPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("config file should exist")
	}

	// Try to init again (should fail)
	if err := InitConfig(false); err == nil {
		t.Fatal("should fail when config already exists")
	}

	// Force init should work
	if err := InitConfig(true); err != nil {
		t.Fatalf("failed to force init config: %v", err)
	}
}

func TestProfileManagement(t *testing.T) {
	setupTestConfig(t)

	// Add profile (without saving credentials)
	if err := AddProfile("prod", &Profile{
		AccessKeyID: "prod-id",
		Region:      "cn-hangzhou",
	}, false); err != nil {
		t.Fatalf("failed to add profile: %v", err)
	}

	// Add another profile
	if err := AddProfile("dev", &Profile{
		AccessKeyID: "dev-id",
		Region:      "cn-shanghai",
	}, false); err != nil {
		t.Fatalf("failed to add profile: %v", err)
	}

	// List profiles
	profiles := ListProfiles()
	if len(profiles) != 2 {
		t.Errorf("expected 2 profiles, got %d", len(profiles))
	}

	// Switch profile
	if err := SwitchProfile("dev"); err != nil {
		t.Fatalf("failed to switch profile: %v", err)
	}

	// Get current profile
	profile, err := GetProfile()
	if err != nil {
		t.Fatalf("failed to get profile: %v", err)
	}
	if profile.AccessKeyID != "dev-id" {
		t.Errorf("expected access key id 'dev-id', got '%s'", profile.AccessKeyID)
	}

	// Remove profile
	if err := RemoveProfile("prod"); err != nil {
		t.Fatalf("failed to remove profile: %v", err)
	}

	// Verify removal
	profiles = ListProfiles()
	if len(profiles) != 1 {
		t.Errorf("expected 1 profile, got %d", len(profiles))
	}

	// Current profile should still be "dev"
	cfg, _ := Load()
	if cfg.CurrentProfile != "dev" {
		t.Errorf("expected current profile 'dev', got '%s'", cfg.CurrentProfile)
	}
}

func TestResetConfig(t *testing.T) {
	setupTestConfig(t)

	// Create config
	if err := InitConfig(false); err != nil {
		t.Fatalf("failed to init config: %v", err)
	}

	// Reset config
	if err := ResetConfig(); err != nil {
		t.Fatalf("failed to reset config: %v", err)
	}

	// File should not exist
	path, _ := GetConfigPath()
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("config file should not exist after reset")
	}

	// HasConfig should return false
	if HasConfig() {
		t.Fatal("HasConfig should return false after reset")
	}
}

func TestHasConfig(t *testing.T) {
	setupTestConfig(t)

	// Should not have config initially
	if HasConfig() {
		t.Fatal("should not have config initially")
	}

	// Init config
	if err := InitConfig(false); err != nil {
		t.Fatalf("failed to init config: %v", err)
	}

	// Should have config now
	if !HasConfig() {
		t.Fatal("should have config after init")
	}
}

func TestSaveCredentials(t *testing.T) {
	setupTestConfig(t)

	// Init config
	if err := InitConfig(false); err != nil {
		t.Fatalf("failed to init config: %v", err)
	}

	// Load config
	cfg, err := Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	// Default should be false
	if cfg.SaveCredentials {
		t.Fatal("save_credentials should be false by default")
	}

	// Add profile without saving credentials
	if err := AddProfile("test", &Profile{
		AccessKeyID:     "test-id",
		AccessKeySecret: "test-secret",
		Region:          "cn-hangzhou",
	}, false); err != nil {
		t.Fatalf("failed to add profile: %v", err)
	}

	// Reload config
	mu.Lock()
	globalConfig = nil
	mu.Unlock()

	cfg, err = Load()
	if err != nil {
		t.Fatalf("failed to reload config: %v", err)
	}

	// Secret should not be saved
	if cfg.Profiles["test"].AccessKeySecret != "" {
		t.Fatal("access_key_secret should be empty when save_credentials is false")
	}
}

func TestUpdateTheme(t *testing.T) {
	setupTestConfig(t)

	// Init config
	if err := InitConfig(false); err != nil {
		t.Fatalf("failed to init config: %v", err)
	}

	// Update theme
	if err := UpdateTheme("dark"); err != nil {
		t.Fatalf("failed to update theme: %v", err)
	}

	// Reload config
	mu.Lock()
	globalConfig = nil
	mu.Unlock()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("failed to reload config: %v", err)
	}

	if cfg.Theme != "dark" {
		t.Errorf("expected theme 'dark', got '%s'", cfg.Theme)
	}

	// Invalid theme should fail
	if err := UpdateTheme("invalid"); err == nil {
		t.Fatal("should fail with invalid theme")
	}
}

func TestMigration(t *testing.T) {
	setupTestConfig(t)

	// Create a v1 config with secret
	cfg := &Config{
		Version:     1,
		MemoryLimit: 256,
		Concurrency: 3,
		Profiles: map[string]*Profile{
			"prod": {
				AccessKeyID:     "prod-id",
				AccessKeySecret: "some-secret",
				Region:          "cn-hangzhou",
			},
		},
	}

	// Save it
	if err := Save(cfg); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	// Reset cache
	mu.Lock()
	globalConfig = nil
	mu.Unlock()

	// Load should migrate
	loaded, err := Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	// Should be migrated to v2
	if loaded.Version != 2 {
		t.Errorf("expected version 2, got %d", loaded.Version)
	}

	// save_credentials should be true because profile had secret
	if !loaded.SaveCredentials {
		t.Fatal("save_credentials should be true after migration")
	}
}
