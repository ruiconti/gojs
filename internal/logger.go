package internal

import (
	"fmt"
	"os"
)

type LoggerMode byte

const (
	ModeDebug = 1
	ModeInfo  = iota
	ModeError
	ModeWarn
)

type Logger interface {
	Debug(string, ...any)
	Info(string, ...any)
	Error(string, ...any)
	Warn(string, ...any)
	DumpLogs()
}

type SimpleLogger struct {
	mode   LoggerMode
	writer *LogWriter
}

type LogWriter struct {
	buff [][]byte
}

func (w *LogWriter) WriteString(s string) (n int, err error) {
	line := make([]byte, len(s)+1)
	line = append(line, s...)
	line = append(line, '\n')

	w.buff = append(w.buff, line)
	return len(s), nil
}

func (w *LogWriter) WriteToStdout() {
	for _, line := range w.buff {
		os.Stdout.Write(line)
	}
}

func NewSimpleLogger(mode LoggerMode) *SimpleLogger {
	writer := &LogWriter{buff: [][]byte{}}
	return &SimpleLogger{mode: mode, writer: writer}
}

func (l *SimpleLogger) Debug(format string, args ...any) {
	if ModeDebug&l.mode > 0 {
		l.writer.WriteString(fmt.Sprintf(format, args...))
	}
}

func (l *SimpleLogger) Info(format string, args ...any) {
	if ModeInfo&l.mode > 0 {
		l.writer.WriteString(fmt.Sprintf(format, args...))
	}
}

func (l *SimpleLogger) Warn(format string, args ...any) {
	if ModeWarn&l.mode > 0 {
		l.writer.WriteString(fmt.Sprintf(format, args...))
	}
}

func (l *SimpleLogger) Error(format string, args ...any) {
	s := fmt.Sprintf(format, args...)
	l.writer.WriteString(fmt.Sprintf("ERROR:%s", s))
}

func (l *SimpleLogger) DumpLogs() {
	l.writer.WriteToStdout()
}
