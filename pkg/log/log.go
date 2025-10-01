/*
Copyright (C) 2025 Keith Chu <cqroot@outlook.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package log

import (
	"io"
	"sync"

	"github.com/rs/zerolog"
)

// Logger is a wrapped zerolog logger with thread-safe capabilities
type Logger struct {
	zerolog.Logger
	mu sync.Mutex
}

var once sync.Once

// New creates a new logger with the given output writer
func New(w io.Writer) *Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := zerolog.New(w).
		With().
		Timestamp().
		Caller().
		Logger()
	return &Logger{Logger: logger}
}

func NewFromLogger(logger zerolog.Logger) *Logger {
	return &Logger{
		Logger: logger,
	}
}

// With returns a zerolog.Context for adding structured fields
func (l *Logger) With() zerolog.Context {
	return l.Logger.With()
}

// Level sets the logging level (thread-safe)
func (l *Logger) Level(level zerolog.Level) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Logger = l.Logger.Level(level)
	return l
}

// Output sets the output writer for the logger (thread-safe)
func (l *Logger) Output(w io.Writer) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Logger = l.Logger.Output(w)
	return l
}

// Debug creates a debug level log event
func (l *Logger) Debug() *zerolog.Event {
	return l.Logger.Debug()
}

// Info creates an info level log event
func (l *Logger) Info() *zerolog.Event {
	return l.Logger.Info()
}

// Warn creates a warning level log event
func (l *Logger) Warn() *zerolog.Event {
	return l.Logger.Warn()
}

// Error creates an error level log event
func (l *Logger) Error() *zerolog.Event {
	return l.Logger.Error()
}

// Fatal creates a fatal level log event and exits
func (l *Logger) Fatal() *zerolog.Event {
	return l.Logger.Fatal()
}

// Panic creates a panic level log event and panics
func (l *Logger) Panic() *zerolog.Event {
	return l.Logger.Panic()
}

// Log creates a log event without specific level
func (l *Logger) Log() *zerolog.Event {
	return l.Logger.Log()
}

// Print provides compatibility with standard log package (no level)
func (l *Logger) Print(v ...any) {
	l.Logger.Print(v...)
}

// Printf provides compatibility with standard log package (formatted, no level)
func (l *Logger) Printf(format string, v ...any) {
	l.Logger.Printf(format, v...)
}

// Sample creates a sampled logger
func (l *Logger) Sample(s zerolog.Sampler) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Logger = l.Logger.Sample(s)
	return l
}

// Hook adds a hook to the logger
func (l *Logger) Hook(h zerolog.Hook) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Logger = l.Logger.Hook(h)
	return l
}
