package errors

import (
	"fmt"
	"strings"
)

// ErrorType represents the type of error.
type ErrorType int

const (
	// ErrorTypeGeneric is a generic error.
	ErrorTypeGeneric ErrorType = iota
	// ErrorTypeAuth is an authentication error.
	ErrorTypeAuth
	// ErrorTypeNetwork is a network error.
	ErrorTypeNetwork
	// ErrorTypePermission is a permission error.
	ErrorTypePermission
	// ErrorTypeNotFound is a not found error.
	ErrorTypeNotFound
	// ErrorTypeValidation is a validation error.
	ErrorTypeValidation
	// ErrorTypeTimeout is a timeout error.
	ErrorTypeTimeout
)

// AppError represents an application error with user-friendly message.
type AppError struct {
	// Type is the error type.
	Type ErrorType
	// Message is the original error message.
	Message string
	// UserMessage is the user-friendly message.
	UserMessage string
	// Suggestion is the suggestion for fixing the error.
	Suggestion string
	// Err is the original error.
	Err error
}

// Error implements the error interface.
func (e *AppError) Error() string {
	return e.Message
}

// Unwrap implements the unwrap interface.
func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError.
func New(errType ErrorType, message string, userMessage string, suggestion string) *AppError {
	return &AppError{
		Type:        errType,
		Message:     message,
		UserMessage: userMessage,
		Suggestion:  suggestion,
	}
}

// Wrap wraps an error with user-friendly information.
func Wrap(err error, errType ErrorType) *AppError {
	if err == nil {
		return nil
	}

	// Check if it's already an AppError
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}

	message := err.Error()
	userMessage, suggestion := translateError(message)

	return &AppError{
		Type:        errType,
		Message:     message,
		UserMessage: userMessage,
		Suggestion:  suggestion,
		Err:         err,
	}
}

// translateError translates technical error messages to user-friendly messages.
func translateError(message string) (string, string) {
	messageLower := strings.ToLower(message)

	// Authentication errors
	if strings.Contains(messageLower, "invalidaccesskeyid") ||
		strings.Contains(messageLower, "accesskey") ||
		strings.Contains(messageLower, "authentication") ||
		strings.Contains(messageLower, "signature") {
		return "AccessKey 认证失败",
			"请检查 AccessKey ID 和 Secret 是否正确，或使用 'cloud-manage config add' 重新配置"
	}

	// Permission errors
	if strings.Contains(messageLower, "permission") ||
		strings.Contains(messageLower, "denied") ||
		strings.Contains(messageLower, "forbidden") ||
		strings.Contains(messageLower, "unauthorized") {
		return "权限不足",
			"请检查当前账号是否有访问该资源的权限"
	}

	// Not found errors
	if strings.Contains(messageLower, "notfound") ||
		strings.Contains(messageLower, "not found") ||
		strings.Contains(messageLower, "does not exist") {
		return "资源不存在",
			"请检查资源 ID 或名称是否正确"
	}

	// Network errors
	if strings.Contains(messageLower, "timeout") ||
		strings.Contains(messageLower, "deadline") {
		return "请求超时",
			"请检查网络连接，或稍后重试"
	}

	if strings.Contains(messageLower, "connection") ||
		strings.Contains(messageLower, "network") ||
		strings.Contains(messageLower, "dial") {
		return "网络连接失败",
			"请检查网络连接是否正常"
	}

	// Rate limit errors
	if strings.Contains(messageLower, "throttl") ||
		strings.Contains(messageLower, "rate limit") ||
		strings.Contains(messageLower, "too many requests") {
		return "请求过于频繁",
			"请稍后重试，或减少并发查询数量"
	}

	// Validation errors
	if strings.Contains(messageLower, "invalid") ||
		strings.Contains(messageLower, "validation") {
		return "参数错误",
			"请检查输入参数是否正确"
	}

	// Default
	return "操作失败",
		"请检查输入参数，或查看详细错误信息"
}

// FormatError formats an error for display.
// In CLI mode, it shows technical details.
// In GUI/TUI mode, it shows user-friendly message.
func FormatError(err error, isCLI bool) string {
	if err == nil {
		return ""
	}

	appErr, ok := err.(*AppError)
	if !ok {
		// Wrap generic errors
		appErr = Wrap(err, ErrorTypeGeneric)
	}

	if isCLI {
		// CLI mode: show technical details
		return fmt.Sprintf("Error: %s\n建议: %s", appErr.Message, appErr.Suggestion)
	}

	// GUI/TUI mode: show user-friendly message
	msg := appErr.UserMessage
	if appErr.Suggestion != "" {
		msg += "\n" + appErr.Suggestion
	}
	return msg
}

// FormatErrorCLI formats an error for CLI display.
func FormatErrorCLI(err error) string {
	return FormatError(err, true)
}

// FormatErrorGUI formats an error for GUI/TUI display.
func FormatErrorGUI(err error) string {
	return FormatError(err, false)
}
