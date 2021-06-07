// Package slog is a structured logging library which is for use by shellrepeater
// or any other servers, it should not be used by commands that are run directly
// by users (e.g. the earthly binary)
package slog

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type field struct {
	key   string
	value string
}

// Logger represents a logger with some structured metadata associated.
type Logger struct {
	fields []field
}

// With adds metadata to the logger.
func (l Logger) With(key string, value interface{}) Logger {
	valueStr := ""
	switch v := value.(type) {
	case string:
		valueStr = v
	case error:
		valueStr = v.Error()
	default:
		valueStr = fmt.Sprintf("%v", v)
	}
	copy := l.clone()
	copy.fields = append(copy.fields, field{
		key:   key,
		value: valueStr,
	})
	return copy
}

// Debug logs debug message.
func (l Logger) Debug(msg string) {
	l.entry().Debug(msg)
}

// Info logs info message.
func (l Logger) Info(msg string) {
	l.entry().Info(msg)
}

// Warning logs warning message.
func (l Logger) Warning(msg string) {
	l.entry().Warning(msg)
}

// Error logs error message.
func (l Logger) Error(err error) {
	l.entry().Error(err.Error())
}

// Fatal logs fatal message and calls os.Exit(1).
func (l Logger) Fatal(msg string) {
	l.entry().Fatal(msg)
}

// Panic logs panic message and calls panic.
func (l Logger) Panic(msg string) {
	l.entry().Panic(msg)
}

func (l Logger) clone() Logger {
	fieldsCopy := make([]field, len(l.fields))
	copy(fieldsCopy, l.fields)
	return Logger{
		fields: fieldsCopy,
	}
}

func (l Logger) entry() *logrus.Entry {
	data := logrus.Fields{}
	for _, field := range l.fields {
		data[field.key] = field.value
	}
	return logrus.WithFields(data)
}
