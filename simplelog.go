package simplelog

import (
	"fmt"
	"io"
	"path/filepath"
	"runtime"
)

const (
	Debug = iota + 1
	Info
	Error
)

type Logger struct {
	out        io.Writer
	prefix     string
	callerSkip int
	errorLevel int
}

func New(out io.Writer, prefix string, callerSkip int) *Logger {
	return &Logger{
		out:        out,
		prefix:     prefix,
		callerSkip: callerSkip,
		errorLevel: Debug,
	}
}

func (l *Logger) SetLevel(errorLevel int) {
	l.errorLevel = errorLevel
}

func (l *Logger) write(level string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(l.callerSkip)
	row := []interface{}{
		fmt.Sprintf("%s%s:%d [%s]", l.prefix, filepath.Base(file), line, level),
	}
	fmt.Fprintln(
		l.out,
		append(row, v...)...,
	)
}

func (l *Logger) writef(level string, f string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(l.callerSkip)
	fmt.Fprintf(
		l.out,
		fmt.Sprintf("%s%s:%d [%s] %s\n", l.prefix, filepath.Base(file), line, level, f),
		v...,
	)
}

func (l *Logger) Debug(v ...interface{}) {
	if l.errorLevel <= Debug {
		l.write("DEBUG", v...)
	}
}

func (l *Logger) Debugf(f string, v ...interface{}) {
	if l.errorLevel <= Debug {
		l.writef("DEBUG", f, v...)
	}
}

func (l *Logger) Info(v ...interface{}) {
	if l.errorLevel <= Info {
		l.write("INFO", v...)
	}
}

func (l *Logger) Infof(f string, v ...interface{}) {
	if l.errorLevel <= Info {
		l.writef("INFO", f, v...)
	}
}

func (l *Logger) Error(v ...interface{}) {
	if l.errorLevel <= Error {
		l.write("ERROR", v...)
	}
}

func (l *Logger) Errorf(f string, v ...interface{}) {
	if l.errorLevel <= Error {
		l.writef("ERROR", f, v...)
	}
}
