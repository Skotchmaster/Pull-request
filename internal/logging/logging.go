package logging

import "log"

type Logger interface {
	Infof(format string, args ...any)
	Errorf(format string, args ...any)
}

type stdLogger struct{}

func NewStdLogger() Logger {
	return &stdLogger{}
}

func (l *stdLogger) Infof(format string, args ...any) {
	log.Printf("[INFO] "+format, args...)
}

func (l *stdLogger) Errorf(format string, args ...any) {
	log.Printf("[ERROR] "+format, args...)
}
