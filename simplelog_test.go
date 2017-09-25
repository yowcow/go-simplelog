package simplelog

import (
	"bytes"
	"io"
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_errorLevelString(t *testing.T) {
	type Case struct {
		Input    int
		Expected string
	}
	cases := []Case{
		{Debug, "DEBUG"},
		{Info, "INFO"},
		{Error, "ERROR"},
		{1234, "???"},
	}

	for _, c := range cases {
		assert.Equal(t, c.Expected, errorLevelString(c.Input), "expecting "+c.Expected)
	}
}

func Test_Itoa(t *testing.T) {
	type Case struct {
		Input    int
		Padding  int
		Expected string
	}
	cases := []Case{
		{1, 1, "1"},
		{2, 2, "02"},
		{123, 4, "0123"},
		{1234, 3, "1234"},
	}

	for _, c := range cases {
		assert.Equal(t, c.Expected, Itoa(c.Input, c.Padding), "expecting "+c.Expected)
	}
}

func newLogger(out io.Writer) *Logger {
	return New(out, "[hoge] ", log.Lshortfile, 2)
}

func Test_Write(t *testing.T) {
	type Case struct {
		Flags    int
		Expected string
	}
	cases := []Case{
		{
			Flags:    0,
			Expected: "[hoge] [DEBUG] hogefuga123\n",
		},
		{
			Flags:    log.LUTC | log.Ldate | log.Lmicroseconds | log.Lshortfile,
			Expected: "[hoge] 2017/02/02 23:01:02.123456 file.txt:2345: [DEBUG] hogefuga123\n",
		},
		{
			Flags:    log.Ldate | log.Ltime | log.Llongfile,
			Expected: "[hoge] 2017/02/03 08:01:02 /path/to/file.txt:2345: [DEBUG] hogefuga123\n",
		},
	}

	for _, c := range cases {
		logbuf := new(bytes.Buffer)
		logger := New(logbuf, "[hoge] ", c.Flags, 2)

		loc, _ := time.LoadLocation("Asia/Tokyo")

		ts := time.Date(2017, 2, 3, 8, 1, 2, 123456789, loc)
		logger.Write(Debug, ts, "/path/to/file.txt", 2345, "hoge", "fuga", 1, 2, 3)

		assert.Equal(t, c.Expected, logbuf.String(), "expecting "+c.Expected)
	}
}

func Test_Debug(t *testing.T) {
	type Case struct {
		Input    []interface{}
		Expected *regexp.Regexp
	}
	cases := []Case{
		{
			Input:    []interface{}{"hoge"},
			Expected: regexp.MustCompile(`\[DEBUG\] hoge\n$`),
		},
		{
			Input:    []interface{}{"hoge", 1, 2},
			Expected: regexp.MustCompile(`\[DEBUG\] hoge12\n$`),
		},
	}

	for _, c := range cases {
		logbuf := new(bytes.Buffer)
		logger := newLogger(logbuf)

		logger.Debug(c.Input...)

		assert.True(t, c.Expected.Match(logbuf.Bytes()))
	}
}

func Test_Debugf(t *testing.T) {
	logbuf := new(bytes.Buffer)
	logger := newLogger(logbuf)

	logger.Debugf("hoge %s -- %d", "fuga", 123)
	re := regexp.MustCompile(`\[DEBUG\] hoge fuga -- 123\n$`)

	assert.True(t, re.Match(logbuf.Bytes()))
}

func Test_Debug_for_levels(t *testing.T) {
	type Case struct {
		Level       int
		ShouldWrite bool
	}
	cases := []Case{
		{Debug, true},
		{Info, false},
		{Error, false},
	}

	for _, c := range cases {
		logbuf := new(bytes.Buffer)
		logger := newLogger(logbuf)
		logger.SetLevel(c.Level)

		logger.Debug("hoge")

		assert.True(t, (len(logbuf.String()) > 0) == c.ShouldWrite)
	}
}

func Test_Info_for_levels(t *testing.T) {
	type Case struct {
		Level       int
		ShouldWrite bool
	}
	cases := []Case{
		{Debug, true},
		{Info, true},
		{Error, false},
	}

	for _, c := range cases {
		logbuf := new(bytes.Buffer)
		logger := newLogger(logbuf)
		logger.SetLevel(c.Level)

		logger.Info("hoge")

		assert.True(t, (len(logbuf.String()) > 0) == c.ShouldWrite)
	}
}

func Test_Error_for_levels(t *testing.T) {
	type Case struct {
		Level       int
		ShouldWrite bool
	}
	cases := []Case{
		{Debug, true},
		{Info, true},
		{Error, true},
	}

	for _, c := range cases {
		logbuf := new(bytes.Buffer)
		logger := newLogger(logbuf)
		logger.SetLevel(c.Level)

		logger.Error("hoge")

		assert.True(t, (len(logbuf.String()) > 0) == c.ShouldWrite)
	}
}
