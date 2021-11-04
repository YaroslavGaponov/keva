package logger

type Logger struct {
	driver Driver
}

func NewLogger(driver Driver) *Logger {
	return &Logger{driver: driver}
}

func (l *Logger) Info(message string, args ...interface{}) {
	l.driver.Log(INFO, message, args...)
}

func (l *Logger) Warn(message string, args ...interface{}) {
	l.driver.Log(WARN, message, args...)
}

func (l *Logger) Error(message string, args ...interface{}) {
	l.driver.Log(ERROR, message, args...)
}

func (l *Logger) Debug(message string, args ...interface{}) {
	l.driver.Log(DEBUG, message, args...)
}

func (l *Logger) Trace(message string, args ...interface{}) {
	l.driver.Log(TRACE, message, args...)
}

func (l *Logger) Fatal(message string, args ...interface{}) {
	l.driver.Log(FATAL, message, args...)
}
