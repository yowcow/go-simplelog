package simplelog

import (
	"bytes"
	"io"
	"log"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newLogger(out io.Writer) *Logger {
	return New(out, "[hoge] ", log.Lshortfile, 2)
}

func Test_Debug_on_1_arg(t *testing.T) {
	logbuf := new(bytes.Buffer)
	logger := newLogger(logbuf)

	logger.Debug("hoge")
	re := regexp.MustCompile(`\[DEBUG\] hoge\n$`)

	assert.True(t, re.Match(logbuf.Bytes()))
}

func Test_Debug_on_multiple_args(t *testing.T) {
	logbuf := new(bytes.Buffer)
	logger := newLogger(logbuf)

	logger.Debug("hoge", 1, 2)
	re := regexp.MustCompile(`\[DEBUG\] hoge12\n$`)

	assert.True(t, re.Match(logbuf.Bytes()))
}

func Test_Debugf(t *testing.T) {
	logbuf := new(bytes.Buffer)
	logger := newLogger(logbuf)

	logger.Debugf("hoge %s %d", "fuga", 123)
	re := regexp.MustCompile(`\[DEBUG\] hoge fuga 123\n$`)

	assert.True(t, re.Match(logbuf.Bytes()))
}

func Test_Info_on_1_arg(t *testing.T) {
	logbuf := new(bytes.Buffer)
	logger := newLogger(logbuf)

	logger.Info("hoge")
	re := regexp.MustCompile(`\[INFO\] hoge\n$`)

	assert.True(t, re.Match(logbuf.Bytes()))
}

func Test_Info_on_multiple_args(t *testing.T) {
	logbuf := new(bytes.Buffer)
	logger := newLogger(logbuf)

	logger.Info("hoge", 1, 2)
	re := regexp.MustCompile(`\[INFO\] hoge12\n$`)

	assert.True(t, re.Match(logbuf.Bytes()))
}

func Test_Infof(t *testing.T) {
	logbuf := new(bytes.Buffer)
	logger := newLogger(logbuf)

	logger.Infof("hoge %s %d", "fuga", 123)
	re := regexp.MustCompile(`\[INFO\] hoge fuga 123\n$`)

	assert.True(t, re.Match(logbuf.Bytes()))
}

func Test_Error_on_1_arg(t *testing.T) {
	logbuf := new(bytes.Buffer)
	logger := newLogger(logbuf)

	logger.Error("hoge")
	re := regexp.MustCompile(`\[ERROR\] hoge\n$`)

	assert.True(t, re.Match(logbuf.Bytes()))
}

func Test_Error_on_multiple_args(t *testing.T) {
	logbuf := new(bytes.Buffer)
	logger := newLogger(logbuf)

	logger.Error("hoge", 1, 2)
	re := regexp.MustCompile(`\[ERROR\] hoge12\n$`)

	assert.True(t, re.Match(logbuf.Bytes()))
}

func Test_Errorf(t *testing.T) {
	logbuf := new(bytes.Buffer)
	logger := newLogger(logbuf)

	logger.Errorf("hoge %s %d", "fuga", 123)
	re := regexp.MustCompile(`\[ERROR\] hoge fuga 123\n$`)

	assert.True(t, re.Match(logbuf.Bytes()))
}

func Test_Debug_on_level_debug(t *testing.T) {
	logbuf := new(bytes.Buffer)
	logger := newLogger(logbuf)
	logger.SetLevel(Debug)

	logger.Debug("hoge")
	re := regexp.MustCompile(`hoge\n$`)

	assert.True(t, re.Match(logbuf.Bytes()))
}

func Test_Debugf_on_level_debug(t *testing.T) {
	logbuf := new(bytes.Buffer)
	logger := newLogger(logbuf)
	logger.SetLevel(Debug)

	logger.Debugf("%s", "hoge")
	re := regexp.MustCompile(`hoge\n$`)

	assert.True(t, re.Match(logbuf.Bytes()))
}

func Test_Debug_on_level_info(t *testing.T) {
	logbuf := new(bytes.Buffer)
	logger := newLogger(logbuf)
	logger.SetLevel(Info)

	logger.Debug("hoge")

	assert.Equal(t, "", logbuf.String())
}

func Test_Debugf_on_level_info(t *testing.T) {
	logbuf := new(bytes.Buffer)
	logger := newLogger(logbuf)
	logger.SetLevel(Info)

	logger.Debugf("%s", "hoge")

	assert.Equal(t, "", logbuf.String())
}

func Test_Info_on_level_info(t *testing.T) {
	logbuf := new(bytes.Buffer)
	logger := newLogger(logbuf)
	logger.SetLevel(Info)

	logger.Info("hoge")
	re := regexp.MustCompile(`hoge\n$`)

	assert.True(t, re.Match(logbuf.Bytes()))
}

func Test_Infof_on_level_info(t *testing.T) {
	logbuf := new(bytes.Buffer)
	logger := newLogger(logbuf)
	logger.SetLevel(Info)

	logger.Infof("%s", "hoge")
	re := regexp.MustCompile(`hoge\n$`)

	assert.True(t, re.Match(logbuf.Bytes()))
}

func Test_Info_on_level_error(t *testing.T) {
	logbuf := new(bytes.Buffer)
	logger := newLogger(logbuf)
	logger.SetLevel(Error)

	logger.Info("hoge")

	assert.Equal(t, "", logbuf.String())
}

func Test_Infof_on_level_error(t *testing.T) {
	logbuf := new(bytes.Buffer)
	logger := newLogger(logbuf)
	logger.SetLevel(Error)

	logger.Infof("%s", "hoge")

	assert.Equal(t, "", logbuf.String())
}

func Test_Error_on_level_error(t *testing.T) {
	logbuf := new(bytes.Buffer)
	logger := newLogger(logbuf)
	logger.SetLevel(Error)

	logger.Error("hoge")
	re := regexp.MustCompile(`hoge\n$`)

	assert.True(t, re.Match(logbuf.Bytes()))
}

func Test_Errorf_on_level_error(t *testing.T) {
	logbuf := new(bytes.Buffer)
	logger := newLogger(logbuf)
	logger.SetLevel(Error)

	logger.Errorf("%s", "hoge")
	re := regexp.MustCompile(`hoge\n$`)

	assert.True(t, re.Match(logbuf.Bytes()))
}
