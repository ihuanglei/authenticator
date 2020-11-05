package logger

import (
	xlog "xorm.io/xorm/log"
)

// XormLog xorm日志
type XormLog struct {
	xlog.DiscardLogger
	showSQL bool
	level   xlog.LogLevel
}

// Debug .
func (l *XormLog) Debug(v ...interface{}) {
	if l.Level() <= xlog.LOG_DEBUG {
		Debug(v...)
	}
}

// Debugf .
func (l *XormLog) Debugf(format string, v ...interface{}) {
	if l.Level() <= xlog.LOG_DEBUG {
		Debugf(format, v...)
	}
}

// Error .
func (l *XormLog) Error(v ...interface{}) {
	if l.Level() <= xlog.LOG_ERR {
		Error(v...)
	}
}

// Errorf .
func (l *XormLog) Errorf(format string, v ...interface{}) {
	if l.Level() <= xlog.LOG_ERR {
		Errorf(format, v...)
	}
}

// Info .
func (l *XormLog) Info(v ...interface{}) {
	if l.Level() <= xlog.LOG_INFO {
		Debug(v...)
	}
}

// Infof .
func (l *XormLog) Infof(format string, v ...interface{}) {
	if l.Level() <= xlog.LOG_INFO {
		Debugf(format, v...)
	}
}

// Warn .
func (l *XormLog) Warn(v ...interface{}) {
	if l.Level() <= xlog.LOG_WARNING {
		Warn(v...)
	}
}

// Warnf .
func (l *XormLog) Warnf(format string, v ...interface{}) {
	if l.Level() <= xlog.LOG_WARNING {
		Warnf(format, v...)
	}
}

// IsShowSQL .
func (l *XormLog) IsShowSQL() bool {
	return l.showSQL
}

// Level empty implementation
func (l *XormLog) Level() xlog.LogLevel {
	return l.level
}

// NewXormLogger .
func NewXormLogger(level int, show bool) *XormLog {
	return &XormLog{level: xlog.LogLevel(level), showSQL: show}
}
