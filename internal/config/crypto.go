package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
	"unicode"

	"golang.org/x/crypto/pbkdf2"
)

const (
	// SaltSize is the size of the salt in bytes.
	SaltSize = 16
	// NonceSize is the size of the nonce in bytes.
	NonceSize = 12
	// Iterations is the number of PBKDF2 iterations.
	Iterations = 100000
	// KeySize is the size of the AES key in bytes.
	KeySize = 32
)

// EncryptedPrefix is the prefix for encrypted secrets.
const EncryptedPrefix = "encrypted:"

// Encrypt encrypts a plaintext string using AES-256-GCM with a password.
func Encrypt(plaintext string, password string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	// Generate random salt
	salt := make([]byte, SaltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Derive key from password
	key := pbkdf2.Key([]byte(password), salt, Iterations, KeySize, sha256.New)

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate random nonce
	nonce := make([]byte, NonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt
	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)

	// Combine salt + nonce + ciphertext
	combined := make([]byte, 0, SaltSize+NonceSize+len(ciphertext))
	combined = append(combined, salt...)
	combined = append(combined, nonce...)
	combined = append(combined, ciphertext...)

	// Encode to base64
	encoded := base64.StdEncoding.EncodeToString(combined)

	return EncryptedPrefix + encoded, nil
}

// Decrypt decrypts an encrypted string using AES-256-GCM with a password.
func Decrypt(encrypted string, password string) (string, error) {
	if encrypted == "" {
		return "", nil
	}

	// Check prefix
	if !strings.HasPrefix(encrypted, EncryptedPrefix) {
		return "", fmt.Errorf("invalid encrypted format: missing prefix")
	}

	// Remove prefix
	encoded := strings.TrimPrefix(encrypted, EncryptedPrefix)

	// Decode from base64
	combined, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	// Check minimum length
	minLen := SaltSize + NonceSize
	if len(combined) < minLen {
		return "", fmt.Errorf("invalid encrypted data: too short")
	}

	// Extract salt, nonce, ciphertext
	salt := combined[:SaltSize]
	nonce := combined[SaltSize:SaltSize+NonceSize]
	ciphertext := combined[SaltSize+NonceSize:]

	// Derive key from password
	key := pbkdf2.Key([]byte(password), salt, Iterations, KeySize, sha256.New)

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}

// IsEncrypted checks if a string is encrypted.
func IsEncrypted(s string) bool {
	return strings.HasPrefix(s, EncryptedPrefix)
}

// ValidatePassword validates a password against the password policy.
func ValidatePassword(password string, policy PasswordPolicy) error {
	if len(password) < policy.MinLength {
		return fmt.Errorf("密码长度至少需要 %d 位", policy.MinLength)
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
		}
	}

	if policy.RequireUppercase && !hasUpper {
		return fmt.Errorf("密码需要包含大写字母")
	}
	if policy.RequireLowercase && !hasLower {
		return fmt.Errorf("密码需要包含小写字母")
	}
	if policy.RequireDigit && !hasDigit {
		return fmt.Errorf("密码需要包含数字")
	}
	if policy.RequireSpecial && !hasSpecial {
		return fmt.Errorf("密码需要包含特殊字符")
	}

	return nil
}

// EncryptProfile encrypts the access key secret in a profile.
func EncryptProfile(profile *Profile, password string) (*Profile, error) {
	if profile.AccessKeySecret == "" || IsEncrypted(profile.AccessKeySecret) {
		return profile, nil
	}

	encrypted, err := Encrypt(profile.AccessKeySecret, password)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt access key secret: %w", err)
	}

	// Create a copy to avoid modifying the original
	encryptedProfile := *profile
	encryptedProfile.AccessKeySecret = encrypted

	return &encryptedProfile, nil
}

// DecryptProfile decrypts the access key secret in a profile.
func DecryptProfile(profile *Profile, password string) (*Profile, error) {
	if profile.AccessKeySecret == "" || !IsEncrypted(profile.AccessKeySecret) {
		return profile, nil
	}

	decrypted, err := Decrypt(profile.AccessKeySecret, password)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt access key secret: %w", err)
	}

	// Create a copy to avoid modifying the original
	decryptedProfile := *profile
	decryptedProfile.AccessKeySecret = decrypted

	return &decryptedProfile, nil
}

// EncryptAllSecrets encrypts all profile secrets in the config.
func EncryptAllSecrets(cfg *Config, password string) error {
	for name, profile := range cfg.Profiles {
		encrypted, err := EncryptProfile(profile, password)
		if err != nil {
			return fmt.Errorf("failed to encrypt profile %s: %w", name, err)
		}
		cfg.Profiles[name] = encrypted
	}
	return nil
}

// DecryptAllSecrets decrypts all profile secrets in the config.
func DecryptAllSecrets(cfg *Config, password string) error {
	for name, profile := range cfg.Profiles {
		decrypted, err := DecryptProfile(profile, password)
		if err != nil {
			return fmt.Errorf("failed to decrypt profile %s: %w", name, err)
		}
		cfg.Profiles[name] = decrypted
	}
	return nil
}
