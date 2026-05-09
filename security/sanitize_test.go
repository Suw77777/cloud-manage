package security

import (
	"errors"
	"testing"
)

func TestMaskAccessKey(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"short key", "abc", "***"},
		{"8 char key", "12345678", "********"},
		{"normal key", "LTAI5tTest1234567890", "LTAI************7890"},
		{"empty string", "", ""},
		{"20 char key", "ABCDEFGHIJKLMNOPQRST", "ABCD************QRST"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaskAccessKey(tt.input)
			if result != tt.expected {
				t.Errorf("MaskAccessKey(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSanitizeErrorMessage(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		result := SanitizeErrorMessage(nil)
		if result != "" {
			t.Errorf("expected empty string for nil error, got %q", result)
		}
	})

	t.Run("error without key", func(t *testing.T) {
		err := errors.New("some connection error")
		result := SanitizeErrorMessage(err)
		if result != "some connection error" {
			t.Errorf("expected original message, got %q", result)
		}
	})

	t.Run("error with LTAI key", func(t *testing.T) {
		err := errors.New("auth failed for key LTAI5tFakeKey1234567890 expired")
		result := SanitizeErrorMessage(err)
		if result == err.Error() {
			t.Error("expected key to be masked")
		}
	})
}
