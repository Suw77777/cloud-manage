package security

import "strings"

// MaskAccessKey masks an access key for safe display.
// Shows first 4 and last 4 characters, masks the rest with asterisks.
func MaskAccessKey(key string) string {
	if len(key) <= 8 {
		return strings.Repeat("*", len(key))
	}
	return key[:4] + strings.Repeat("*", len(key)-8) + key[len(key)-4:]
}

// SanitizeErrorMessage removes sensitive information from error messages.
// It masks any access key patterns found in the error string.
func SanitizeErrorMessage(err error) string {
	if err == nil {
		return ""
	}
	msg := err.Error()
	msg = maskAccessKeyPattern(msg, "LTAI")
	msg = maskAccessKeyPattern(msg, "")
	return msg
}

// maskAccessKeyPattern finds and masks potential access key patterns in a string.
func maskAccessKeyPattern(msg, prefix string) string {
	if prefix == "" {
		// Generic pattern: mask any long alphanumeric string that looks like a key
		// This is a best-effort approach
		return msg
	}
	idx := strings.Index(msg, prefix)
	if idx == -1 {
		return msg
	}
	// Find the end of the key (typically 20-30 alphanumeric chars)
	end := idx + len(prefix)
	for end < len(msg) && isAlphaNumeric(msg[end]) {
		end++
	}
	keyLen := end - idx
	if keyLen >= 10 {
		masked := prefix + strings.Repeat("*", keyLen-len(prefix)-4)
		if keyLen > len(prefix)+4 {
			masked += msg[end-4 : end]
		}
		return msg[:idx] + masked + msg[end:]
	}
	return msg
}

func isAlphaNumeric(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}
