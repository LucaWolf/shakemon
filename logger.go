package main

import (
	"bytes"
	"fmt"
	"os"
)

// type Logger interface {
// 	Errorf(format string, v ...interface{})
// 	Warnf(format string, v ...interface{})
// 	Debugf(format string, v ...interface{})
// }

type Logger struct {
	data bytes.Buffer
}

func NewLogger() *Logger {
	return &Logger{
		data: bytes.Buffer{},
	}
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.data.WriteString(fmt.Sprintf("error: "+format+"\n", v...))
	l.data.WriteTo(os.Stdout)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.data.WriteString(fmt.Sprintf("warn: "+format+"\n", v...))
	l.data.WriteTo(os.Stdout)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.data.WriteString(fmt.Sprintf("debug: "+format+"\n", v...))
	l.data.WriteTo(os.Stdout)
}
