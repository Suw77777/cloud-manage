package config

import (
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	password := "TestPassword123"
	plaintext := "my-secret-access-key"

	// Encrypt
	encrypted, err := Encrypt(plaintext, password)
	if err != nil {
		t.Fatalf("failed to encrypt: %v", err)
	}

	// Check prefix
	if !IsEncrypted(encrypted) {
		t.Fatal("encrypted string should have prefix")
	}

	// Decrypt
	decrypted, err := Decrypt(encrypted, password)
	if err != nil {
		t.Fatalf("failed to decrypt: %v", err)
	}

	// Check result
	if decrypted != plaintext {
		t.Errorf("expected '%s', got '%s'", plaintext, decrypted)
	}
}

func TestEncryptEmptyString(t *testing.T) {
	password := "TestPassword123"

	// Encrypt empty string
	encrypted, err := Encrypt("", password)
	if err != nil {
		t.Fatalf("failed to encrypt empty string: %v", err)
	}

	// Should return empty
	if encrypted != "" {
		t.Errorf("expected empty string, got '%s'", encrypted)
	}
}

func TestDecryptEmptyString(t *testing.T) {
	password := "TestPassword123"

	// Decrypt empty string
	decrypted, err := Decrypt("", password)
	if err != nil {
		t.Fatalf("failed to decrypt empty string: %v", err)
	}

	// Should return empty
	if decrypted != "" {
		t.Errorf("expected empty string, got '%s'", decrypted)
	}
}

func TestDecryptWrongPassword(t *testing.T) {
	password := "TestPassword123"
	wrongPassword := "WrongPassword456"
	plaintext := "my-secret-access-key"

	// Encrypt
	encrypted, err := Encrypt(plaintext, password)
	if err != nil {
		t.Fatalf("failed to encrypt: %v", err)
	}

	// Decrypt with wrong password
	_, err = Decrypt(encrypted, wrongPassword)
	if err == nil {
		t.Fatal("should fail with wrong password")
	}
}

func TestDecryptInvalidFormat(t *testing.T) {
	password := "TestPassword123"

	// Try to decrypt invalid string
	_, err := Decrypt("invalid-format", password)
	if err == nil {
		t.Fatal("should fail with invalid format")
	}
}

func TestValidatePassword(t *testing.T) {
	policy := PasswordPolicy{
		MinLength:        8,
		RequireUppercase: true,
		RequireLowercase: true,
		RequireDigit:     true,
		RequireSpecial:   false,
	}

	// Valid password
	err := ValidatePassword("TestPass123", policy)
	if err != nil {
		t.Errorf("valid password should pass: %v", err)
	}

	// Too short
	err = ValidatePassword("Test1", policy)
	if err == nil {
		t.Error("short password should fail")
	}

	// No uppercase
	err = ValidatePassword("testpass123", policy)
	if err == nil {
		t.Error("password without uppercase should fail")
	}

	// No lowercase
	err = ValidatePassword("TESTPASS123", policy)
	if err == nil {
		t.Error("password without lowercase should fail")
	}

	// No digit
	err = ValidatePassword("TestPassword", policy)
	if err == nil {
		t.Error("password without digit should fail")
	}

	// With special character required
	policy.RequireSpecial = true
	err = ValidatePassword("TestPass123!", policy)
	if err != nil {
		t.Errorf("password with special should pass: %v", err)
	}

	err = ValidatePassword("TestPass123", policy)
	if err == nil {
		t.Error("password without special should fail")
	}
}

func TestEncryptDecryptProfile(t *testing.T) {
	password := "TestPassword123"

	profile := &Profile{
		AccessKeyID:     "test-id",
		AccessKeySecret: "test-secret",
		Region:          "cn-hangzhou",
	}

	// Encrypt
	encrypted, err := EncryptProfile(profile, password)
	if err != nil {
		t.Fatalf("failed to encrypt profile: %v", err)
	}

	// Check that secret is encrypted
	if !IsEncrypted(encrypted.AccessKeySecret) {
		t.Error("access key secret should be encrypted")
	}

	// Check that other fields are unchanged
	if encrypted.AccessKeyID != profile.AccessKeyID {
		t.Error("access key id should be unchanged")
	}
	if encrypted.Region != profile.Region {
		t.Error("region should be unchanged")
	}

	// Decrypt
	decrypted, err := DecryptProfile(encrypted, password)
	if err != nil {
		t.Fatalf("failed to decrypt profile: %v", err)
	}

	// Check that secret is decrypted
	if decrypted.AccessKeySecret != profile.AccessKeySecret {
		t.Errorf("expected '%s', got '%s'", profile.AccessKeySecret, decrypted.AccessKeySecret)
	}
}

func TestEncryptDecryptAlreadyEncrypted(t *testing.T) {
	password := "TestPassword123"

	profile := &Profile{
		AccessKeyID:     "test-id",
		AccessKeySecret: "encrypted:already-encrypted",
		Region:          "cn-hangzhou",
	}

	// Encrypt should return same profile
	encrypted, err := EncryptProfile(profile, password)
	if err != nil {
		t.Fatalf("failed to encrypt profile: %v", err)
	}

	if encrypted.AccessKeySecret != profile.AccessKeySecret {
		t.Error("already encrypted secret should not be re-encrypted")
	}
}

func TestEncryptDecryptAllSecrets(t *testing.T) {
	password := "TestPassword123"

	cfg := &Config{
		Profiles: map[string]*Profile{
			"prod": {
				AccessKeyID:     "prod-id",
				AccessKeySecret: "prod-secret",
				Region:          "cn-hangzhou",
			},
			"dev": {
				AccessKeyID:     "dev-id",
				AccessKeySecret: "dev-secret",
				Region:          "cn-shanghai",
			},
		},
	}

	// Encrypt all
	err := EncryptAllSecrets(cfg, password)
	if err != nil {
		t.Fatalf("failed to encrypt all secrets: %v", err)
	}

	// Check that all secrets are encrypted
	for name, profile := range cfg.Profiles {
		if !IsEncrypted(profile.AccessKeySecret) {
			t.Errorf("profile %s secret should be encrypted", name)
		}
	}

	// Decrypt all
	err = DecryptAllSecrets(cfg, password)
	if err != nil {
		t.Fatalf("failed to decrypt all secrets: %v", err)
	}

	// Check that all secrets are decrypted
	if cfg.Profiles["prod"].AccessKeySecret != "prod-secret" {
		t.Error("prod secret should be decrypted")
	}
	if cfg.Profiles["dev"].AccessKeySecret != "dev-secret" {
		t.Error("dev secret should be decrypted")
	}
}
