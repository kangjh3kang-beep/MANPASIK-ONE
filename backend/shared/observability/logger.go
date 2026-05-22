package observability

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"
)

// LogLevel represents the severity of a log entry.
type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

// logEntry is the JSON structure written to stdout for Loki/Promtail ingestion.
type logEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     LogLevel               `json:"level"`
	Service   string                 `json:"service"`
	Msg       string                 `json:"msg"`
	Caller    string                 `json:"caller,omitempty"`
	Error     string                 `json:"error,omitempty"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
}

// StructuredLogger outputs JSON-formatted log lines to stdout,
// designed for collection by Promtail and ingestion into Loki.
type StructuredLogger struct {
	serviceName string
	fields      map[string]interface{}
}

// NewLogger creates a StructuredLogger for the given service.
func NewLogger(serviceName string) *StructuredLogger {
	return &StructuredLogger{
		serviceName: serviceName,
		fields:      make(map[string]interface{}),
	}
}

// With returns a new logger that includes the given key-value pair
// in every subsequent log entry. It does not mutate the original logger.
func (l *StructuredLogger) With(key string, value interface{}) *StructuredLogger {
	newFields := make(map[string]interface{}, len(l.fields)+1)
	for k, v := range l.fields {
		newFields[k] = v
	}
	newFields[key] = value
	return &StructuredLogger{
		serviceName: l.serviceName,
		fields:      newFields,
	}
}

// Info logs an informational message.
func (l *StructuredLogger) Info(msg string, fields ...map[string]interface{}) {
	l.log(LogLevelInfo, msg, nil, fields...)
}

// Error logs an error message with the associated error value.
func (l *StructuredLogger) Error(msg string, err error, fields ...map[string]interface{}) {
	l.log(LogLevelError, msg, err, fields...)
}

// Warn logs a warning message.
func (l *StructuredLogger) Warn(msg string, fields ...map[string]interface{}) {
	l.log(LogLevelWarn, msg, nil, fields...)
}

// Debug logs a debug message.
func (l *StructuredLogger) Debug(msg string, fields ...map[string]interface{}) {
	l.log(LogLevelDebug, msg, nil, fields...)
}

// log writes a single JSON log line to stdout.
func (l *StructuredLogger) log(level LogLevel, msg string, err error, extraFields ...map[string]interface{}) {
	entry := logEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		Level:     level,
		Service:   l.serviceName,
		Msg:       msg,
		Caller:    getCaller(3),
	}

	if err != nil {
		entry.Error = err.Error()
	}

	// Merge base fields and extra fields
	merged := make(map[string]interface{}, len(l.fields))
	for k, v := range l.fields {
		merged[k] = v
	}
	for _, f := range extraFields {
		for k, v := range f {
			merged[k] = v
		}
	}
	if len(merged) > 0 {
		entry.Fields = merged
	}

	data, jsonErr := json.Marshal(entry)
	if jsonErr != nil {
		// Fallback: write a plain text message if JSON marshalling fails
		fmt.Fprintf(os.Stderr, "logger: json marshal error: %v (original msg: %s)\n", jsonErr, msg)
		return
	}

	// Write to stdout; Promtail picks up container stdout
	os.Stdout.Write(append(data, '\n'))
}

// getCaller returns the file:line of the caller at the given skip depth.
func getCaller(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown"
	}
	// Use short file path (last component)
	short := file
	for i := len(file) - 1; i >= 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	return fmt.Sprintf("%s:%d", short, line)
}
