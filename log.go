package main

import (
	"fmt"
	"log"
)

type LoggerMode byte

const (
	ModeDebug = iota
	ModeInfo
	ModeError
	ModeWarn
)

type Logger interface {
	Debug(string, ...any)
	Info(string, ...any)
	Error(string, ...any)
	Warn(string, ...any)
}

type SimpleLogger struct {
	mode LoggerMode
}

func NewSimpleLogger(mode LoggerMode) *SimpleLogger {
	return &SimpleLogger{mode: mode}
}

func (l *SimpleLogger) Debug(format string, args ...any) {
	if ModeDebug&l.mode != 0 {
		log.Printf(fmt.Sprintf("debug: %s", format), args)
	}
}

func (l *SimpleLogger) Info(format string, args ...any) {
	if ModeInfo&l.mode != 0 {
		log.Printf(fmt.Sprintf("info: %s", format), args)
	}
}

func (l *SimpleLogger) Warn(format string, args ...any) {
	if ModeWarn&l.mode != 0 {
		log.Printf(fmt.Sprintf("warning: %s", format), args)
	}
}

func (l *SimpleLogger) Error(format string, args ...any) {
	log.Printf(fmt.Sprintf("error: %s", format), args)
}
