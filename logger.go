/************************************************************************************
 *
 * yada (yet another discord api), A Lightweight Go library for Discord API
 *
 * SPDX-License-Identifier: BSD-3-Clause
 *
 * Copyright 2025 Marouane Souiri
 *
 * Licensed under the BSD 3-Clause License.
 * See the LICENSE file for details.
 *
 ************************************************************************************/

package yada

import (
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"os"
	"sync"
	"time"
)

// Logger defines the logging interface
type Logger interface {
	Info(msg string)
	Debug(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)

	// WithField adds a single field to the logger context
	WithField(key string, value any) Logger
	// WithFields adds multiple fields to the logger context
	WithFields(fields map[string]any) Logger
}

// LogLevel defines the severity level
type LogLevel int

const (
	LogLevel_DebugLevel LogLevel = iota
	LogLevel_InfoLevel
	LogLevel_WarnLevel
	LogLevel_ErrorLevel
	LogLevel_FatalLevel
)

type DefaultLogger struct {
	out    io.Writer
	mu     sync.Mutex
	fields map[string]any
	level  LogLevel
}

func NewDefaultLogger(out io.Writer, level LogLevel) *DefaultLogger {
	if out == nil {
		out = os.Stdout
	}
	return &DefaultLogger{
		out:    out,
		fields: make(map[string]any),
		level:  level,
	}
}

func (l *DefaultLogger) WithField(key string, value any) Logger {
	return l.WithFields(map[string]any{key: value})
}

func (l *DefaultLogger) WithFields(fields map[string]any) Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	newFields := make(map[string]any, len(l.fields)+len(fields))
	maps.Copy(newFields, l.fields)
	maps.Copy(newFields, fields)

	return &DefaultLogger{
		out:    l.out,
		fields: newFields,
		level:  l.level,
	}
}

func (l *DefaultLogger) log(level LogLevel, levelStr, msg string) {
	if level < l.level {
		return
	}

	data := make(map[string]any, len(l.fields)+3)
	maps.Copy(data, l.fields)

	data["level"] = levelStr
	data["time"] = time.Now().Format(time.RFC3339)
	data["msg"] = msg

	l.mu.Lock()
	defer l.mu.Unlock()
	enc, err := json.Marshal(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "logger marshal error: %v\n", err)
		return
	}
	fmt.Fprintln(l.out, string(enc))

	if level == LogLevel_FatalLevel {
		os.Exit(1)
	}
}

func (l *DefaultLogger) Info(msg string) {
	l.log(LogLevel_InfoLevel, "info", msg)
}

func (l *DefaultLogger) Debug(msg string) {
	l.log(LogLevel_DebugLevel, "debug", msg)
}

func (l *DefaultLogger) Warn(msg string) {
	l.log(LogLevel_WarnLevel, "warn", msg)
}

func (l *DefaultLogger) Error(msg string) {
	l.log(LogLevel_ErrorLevel, "error", msg)
}

func (l *DefaultLogger) Fatal(msg string) {
	l.log(LogLevel_FatalLevel, "fatal", msg)
}
