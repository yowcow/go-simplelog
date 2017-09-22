package simplelog

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
)

const (
	Debug = iota + 1
	Info
	Error
)

type Logger struct {
	mu         *sync.Mutex
	out        io.Writer
	prefix     string
	calldepth  int
	errorLevel int
}

func New(out io.Writer, prefix string, calldepth int) *Logger {
	return &Logger{
		mu:         &sync.Mutex{},
		out:        out,
		prefix:     prefix,
		calldepth:  calldepth,
		errorLevel: Debug,
	}
}

func (l *Logger) SetLevel(errorLevel int) {
	l.errorLevel = errorLevel
}

func errorLevelString(errorLevel int) string {
	switch errorLevel {
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	case Error:
		return "ERROR"
	default:
		return "???"
	}
}

var bufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func (l *Logger) write(errorLevel int, file string, line int, msgs ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	buf := bufPool.Get().(*bytes.Buffer)
	buf.WriteString(l.prefix)
	buf.WriteString(file)
	buf.WriteByte(':')
	buf.WriteString(strconv.Itoa(line))
	buf.WriteString(" [")
	buf.WriteString(errorLevelString(errorLevel))
	buf.WriteString("]")

	for _, msg := range msgs {
		buf.WriteByte(' ')

		switch msg.(type) {
		case string:
			buf.WriteString(msg.(string))
		case int:
			buf.WriteString(strconv.Itoa(msg.(int)))
		}
	}

	buf.WriteString("\n")
	l.out.Write(buf.Bytes())

	buf.Reset()
	bufPool.Put(buf)
}

func (l *Logger) Output(errorLevel int, v ...interface{}) {
	if errorLevel < l.errorLevel {
		return
	}

	var filename string
	var line int
	_, file, line, ok := runtime.Caller(l.calldepth)

	if !ok {
		filename = "???"
		line = 0
	} else {
		filename = filepath.Base(file)
	}

	l.write(errorLevel, filename, line, v...)
}

func (l *Logger) Debug(v ...interface{}) {
	l.Output(Debug, v...)
}

func (l *Logger) Debugf(f string, v ...interface{}) {
	l.Output(Debug, fmt.Sprintf(f, v...))
}

func (l *Logger) Info(v ...interface{}) {
	l.Output(Info, v...)
}

func (l *Logger) Infof(f string, v ...interface{}) {
	l.Output(Info, fmt.Sprintf(f, v...))
}

func (l *Logger) Error(v ...interface{}) {
	l.Output(Error, v...)
}

func (l *Logger) Errorf(f string, v ...interface{}) {
	l.Output(Error, fmt.Sprintf(f, v...))
}
