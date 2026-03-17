package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type Level string

const (
	DEBUG   Level = "DEBUG"
	INFO    Level = "INFO"
	WARNING Level = "WARNING"
	ERROR   Level = "ERROR"
)

type entry struct {
	Timestamp string                 `json:"timestamp"`
	Level     Level                  `json:"level"`
	Message   string                 `json:"message"`
	File      string                 `json:"file,omitempty"`
	Context   map[string]interface{} `json:"context,omitempty"`
}

type Logger struct {
	out io.Writer
}

// Default - Default logger instance writing to Stdout
var Default = New(os.Stdout)

func New(out io.Writer) *Logger {
	return &Logger{out: out}
}

func (l *Logger) log(level Level, msg string, context map[string]interface{}) {
	_, file, line, _ := runtime.Caller(2)

	// Uses filepath.Base to correctly get the filename on both Linux and Windows
	short := fmt.Sprintf("%s:%d", filepath.Base(file), line)

	e := entry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     level,
		Message:   msg,
		File:      short,
		Context:   context,
	}

	b, _ := json.Marshal(e)
	fmt.Fprintln(l.out, string(b))
}

// Info — Equivalent to Laravel's Log::info()
func (l *Logger) Info(msg string, ctx ...map[string]interface{}) {
	l.log(INFO, msg, mergeContext(ctx))
}

// Warning — Equivalent to Laravel's Log::warning()
func (l *Logger) Warning(msg string, ctx ...map[string]interface{}) {
	l.log(WARNING, msg, mergeContext(ctx))
}

// Error — Equivalent to Laravel's Log::error()
func (l *Logger) Error(msg string, ctx ...map[string]interface{}) {
	l.log(ERROR, msg, mergeContext(ctx))
}

// Debug — Equivalent to Laravel's Log::debug()
func (l *Logger) Debug(msg string, ctx ...map[string]interface{}) {
	l.log(DEBUG, msg, mergeContext(ctx))
}

// Global helpers — Usage: logger.Info("msg", logger.With("key", val))
func Info(msg string, ctx ...map[string]interface{})    { Default.Info(msg, ctx...) }
func Warning(msg string, ctx ...map[string]interface{}) { Default.Warning(msg, ctx...) }
func Error(msg string, ctx ...map[string]interface{})   { Default.Error(msg, ctx...) }
func Debug(msg string, ctx ...map[string]interface{})   { Default.Debug(msg, ctx...) }

// With — Context helper: logger.Info("msg", logger.With("user_id", 1))
func With(key string, value interface{}) map[string]interface{} {
	return map[string]interface{}{key: value}
}

func mergeContext(ctxSlice []map[string]interface{}) map[string]interface{} {
	if len(ctxSlice) == 0 {
		return nil
	}
	merged := make(map[string]interface{})
	for _, c := range ctxSlice {
		for k, v := range c {
			merged[k] = v
		}
	}
	return merged
}
