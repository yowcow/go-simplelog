package simplelog

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"
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
	flags      int
	errorLevel int
}

func New(out io.Writer, prefix string, flags int, calldepth int) *Logger {
	return &Logger{
		mu:         &sync.Mutex{},
		out:        out,
		prefix:     prefix,
		calldepth:  calldepth,
		flags:      flags,
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

func Itoa(i int, pad int) string {
	a := strconv.Itoa(i)

	if fill := pad - len(a); fill > 0 {
		buf := new(bytes.Buffer)
		for ; fill > 0; fill-- {
			buf.WriteByte('0')
		}
		buf.WriteString(a)
		a = buf.String()
	}

	return a
}

var bufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func (l *Logger) Write(errorLevel int, now time.Time, file string, line int, msgs ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	buf := bufPool.Get().(*bytes.Buffer)
	buf.WriteString(l.prefix)

	if l.flags&(log.Ldate|log.Ltime|log.Lmicroseconds) != 0 {
		if l.flags&log.LUTC != 0 {
			now = now.UTC()
		}
		if l.flags&log.Ldate != 0 {
			year, month, day := now.Date()
			buf.WriteString(Itoa(year, 4))
			buf.WriteByte('/')
			buf.WriteString(Itoa(int(month), 2))
			buf.WriteByte('/')
			buf.WriteString(Itoa(day, 2))
			buf.WriteByte(' ')
		}
		if l.flags&(log.Ltime|log.Lmicroseconds) != 0 {
			hour, min, sec := now.Clock()
			buf.WriteString(Itoa(hour, 2))
			buf.WriteByte(':')
			buf.WriteString(Itoa(min, 2))
			buf.WriteByte(':')
			buf.WriteString(Itoa(sec, 2))

			if l.flags&log.Lmicroseconds != 0 {
				buf.WriteByte('.')
				buf.WriteString(Itoa(now.Nanosecond()/1e3, 6))
			}

			buf.WriteByte(' ')
		}
	}

	if l.flags&(log.Lshortfile|log.Llongfile) != 0 {
		if l.flags&log.Lshortfile != 0 {
			buf.WriteString(filepath.Base(file))
		} else if l.flags&log.Llongfile != 0 {
			buf.WriteString(file)
		}

		buf.WriteByte(':')
		buf.WriteString(strconv.Itoa(line))
		buf.WriteByte(':')
		buf.WriteByte(' ')
	}

	buf.WriteByte('[')
	buf.WriteString(errorLevelString(errorLevel))
	buf.WriteString("] ")

	for _, msg := range msgs {
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

	var line int
	_, file, line, ok := runtime.Caller(l.calldepth)

	if !ok {
		file = "???"
		line = 0
	}

	l.Write(errorLevel, time.Now(), file, line, v...)
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
