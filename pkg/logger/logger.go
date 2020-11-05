package logger

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/sirupsen/logrus"
)

type formatter struct {
	logrus.TextFormatter
}

// Format .
func (f *formatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	level := entry.Level
	fmt.Fprintf(b, "\x1b[%dm", f.getColorByLevel(level))
	b.WriteString("[AUTHENTICATOR][")
	b.WriteString(strings.ToUpper(level.String()))
	b.WriteString("] ")
	b.WriteString(entry.Time.Format(f.TimestampFormat))
	b.WriteString(" ")
	if level == logrus.ErrorLevel {
		s := strings.Split(entry.Caller.Function, ".")
		funcName := s[len(s)-1]
		b.WriteString(fmt.Sprintf("%s:%d (%s) ", path.Base(entry.Caller.File), entry.Caller.Line, funcName))
	}
	b.WriteString(entry.Message)
	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *formatter) getColorByLevel(level logrus.Level) int {
	switch level {
	case logrus.DebugLevel:
		return 37
	case logrus.WarnLevel:
		return 33
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return 31
	default:
		return 36
	}
}

func init() {
	f := &formatter{}
	f.ForceColors = true
	f.DisableLevelTruncation = true
	f.TimestampFormat = "2006-01-02 15:04:05"
	f.FullTimestamp = true
	logrus.SetReportCaller(true)
	logrus.SetFormatter(f)
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetOutput(os.Stdout)
	// file, err := os.OpenFile("access.log", os.O_CREATE|os.O_WRONLY, 0666)
	// if err != nil {
	// log.SetOutput(os.Stdout)
	// } else {
	// log.SetOutput(file)
	// }
}

// SetLevel .
func SetLevel(l int) {
	level := logrus.InfoLevel
	switch l {
	case 0:
		level = logrus.PanicLevel
	case 1:
		level = logrus.FatalLevel
	case 2:
		level = logrus.ErrorLevel
	case 3:
		level = logrus.WarnLevel
	case 4:
		level = logrus.InfoLevel
	case 5:
		level = logrus.DebugLevel
	case 6:
		level = logrus.TraceLevel
	}
	logrus.SetLevel(level)
}

// Debug .
func Debug(v ...interface{}) {
	logrus.Debug(v...)
}

// Debugf .
func Debugf(format string, v ...interface{}) {
	logrus.Debugf(format, v...)
}

// Error .
func Error(v ...interface{}) {
	logrus.Error(v...)
}

// Errorln .
func Errorln(v ...interface{}) {
	logrus.Errorln(v)
}

// Errorf .
func Errorf(format string, v ...interface{}) {
	logrus.Errorf(format, v...)
}

// Info .
func Info(v ...interface{}) {
	logrus.Info(v...)
}

// Infof .
func Infof(format string, v ...interface{}) {
	logrus.Infof(format, v...)
}

// Warn .
func Warn(v ...interface{}) {
	logrus.Warn(v...)
}

// Warnf .
func Warnf(format string, v ...interface{}) {
	logrus.Warnf(format, v...)
}

// Panic .
func Panic(v ...interface{}) {
	logrus.Panic(v)
}

// Panicf .
func Panicf(format string, v ...interface{}) {
	logrus.Panicf(format, v)
}

// Fatal .
func Fatal(v ...interface{}) {
	logrus.Fatal(v)
}

// Fatalln .
func Fatalln(v ...interface{}) {
	logrus.Fatalln(v)
}
