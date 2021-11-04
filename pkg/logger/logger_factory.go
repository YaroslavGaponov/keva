package logger

import (
	"sync"

	"github.com/YaroslavGaponov/keva/pkg/logger/driver/console"
	"github.com/YaroslavGaponov/keva/pkg/utils"
)

var (
	lock     sync.Mutex
	instance *Logger
)

func CreateLogger() *Logger {
	lock.Lock()
	defer lock.Unlock()
	if instance == nil {
		alias := utils.GetEnvVariableOrDefult("LOG_LEVEL", "info")
		driver := console.ConsoleLoggerNew(getLogLevel(alias))
		instance = NewLogger(driver)
		instance.Info("Log level is %s", alias)
	}
	return instance
}

func getLogLevel(s string) int {
	switch s {
	case "all":
		return INFO | DEBUG | WARN | ERROR | FATAL | TRACE
	case "none":
		return 0
	case "debug":
		return INFO | DEBUG | TRACE
	case "info":
		return INFO
	case "warn":
		return INFO | WARN
	case "error":
		return INFO | ERROR
	case "fatal":
		return INFO | FATAL
	default:
		return INFO | ERROR | FATAL
	}
}
