package jsonlog

import (
	"os"
	"sync"
)

var (
	once sync.Once
	// logger single instance created here
	logger *Logger
)

// GetLogger returns a singleton instance of Logger
// so it can be reused wheresoever throughout the application
// find ways to test
func GetLggger() *Logger {
	once.Do(func() {
		logger = New(os.Stdout, LevelInfo)
	})
	return logger
}
