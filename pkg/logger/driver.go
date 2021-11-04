package logger

type Driver interface {
	Log(level int, message string, args ...interface{})
}