package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
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

var Default = New(os.Stdout)

func New(out io.Writer) *Logger {
	return &Logger{out: out}
}

func (l *Logger) log(level Level, msg string, context map[string]interface{}) {
	_, file, line, _ := runtime.Caller(2)
	// pega só o caminho relativo
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}

	e := entry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     level,
		Message:   msg,
		File:      fmt.Sprintf("%s:%d", short, line),
		Context:   context,
	}

	b, _ := json.Marshal(e)
	fmt.Fprintln(l.out, string(b))
}

// Info — equivalente ao Log::info() do Laravel
func (l *Logger) Info(msg string, ctx ...map[string]interface{}) {
	l.log(INFO, msg, mergeContext(ctx))
}

// Warning — equivalente ao Log::warning() do Laravel
func (l *Logger) Warning(msg string, ctx ...map[string]interface{}) {
	l.log(WARNING, msg, mergeContext(ctx))
}

// Error — equivalente ao Log::error() do Laravel
func (l *Logger) Error(msg string, ctx ...map[string]interface{}) {
	l.log(ERROR, msg, mergeContext(ctx))
}

// Debug — equivalente ao Log::debug() do Laravel
func (l *Logger) Debug(msg string, ctx ...map[string]interface{}) {
	l.log(DEBUG, msg, mergeContext(ctx))
}

// Atalhos globais — uso: logger.Info("msg", logger.With("key", val))
func Info(msg string, ctx ...map[string]interface{})    { Default.Info(msg, ctx...) }
func Warning(msg string, ctx ...map[string]interface{}) { Default.Warning(msg, ctx...) }
func Error(msg string, ctx ...map[string]interface{})   { Default.Error(msg, ctx...) }
func Debug(msg string, ctx ...map[string]interface{})   { Default.Debug(msg, ctx...) }

// With — helper para contexto: logger.Info("msg", logger.With("user_id", 1))
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
