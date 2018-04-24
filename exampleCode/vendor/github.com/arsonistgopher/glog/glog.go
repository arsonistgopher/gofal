package glog

import (
	"fmt"
	"runtime"
)

// Logger concrete type which embeds the Logging interface for exported 'things'
type Logger struct {
	LoggingBase Logging
	Name        string
}

// Logging interface exported with four methods. These are variadic and need unwinding in each decorator function.
type Logging interface {
	Info(v ...interface{})
	Error(v ...interface{})
	// Critical(v ...interface{})
	Debug(v ...interface{})
}

// Info method provides users the opportunity to decorate their Info logs.
func (l Logger) Info(v ...interface{}) {
	for _, msg := range v {
		l.LoggingBase.Info(msg)
	}
}

// Debug method provides users the opportunity to decorate their Debug logs.
// This method shows functions and packages that call the Debug() func.
func (l Logger) Debug(v ...interface{}) {
	for _, msg := range v {
		pc, _, _, ok := runtime.Caller(1)
		details := runtime.FuncForPC(pc)
		if ok && details != nil {
			l.LoggingBase.Debug(fmt.Sprintf("%s: %v", details.Name(), msg))
		}
	}
}

/*
// Critical method provides users the opportunity to decorate their Critical logs.
func (l Logger) Critical(v ...interface{}) {
	for _, msg := range v {
		l.LoggingBase.Critical(msg)
	}
}
*/

// Error method provides users the opportunity to decorate their Error logs.
func (l Logger) Error(v ...interface{}) {
	for _, msg := range v {
		l.LoggingBase.Error(msg)
	}
}
