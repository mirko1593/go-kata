package logger

import "go.uber.org/zap"

// DefaultLogger ...
var DefaultLogger *zap.Logger

// NewDevelopment ...
func NewDevelopment() (*zap.Logger, error) {
	l, err := zap.NewDevelopment()
	DefaultLogger = l
	return l, err
}

// Sync ...
func Sync() {
	if DefaultLogger != nil {
		DefaultLogger.Sync()
	}
}
