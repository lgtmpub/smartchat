package logger

// Logger is the interface that wraps the basic logging methods. slog.Logger implements this interface.
type Logger interface {
	// Debug logs a debug message.
	Debug(msg string, fields ...any)
	// Info logs an info message.
	Info(msg string, fields ...any)
	// Warn logs a warning message.
	Warn(msg string, fields ...any)
	// Error logs an error message.
	Error(msg string, fields ...any)
}
