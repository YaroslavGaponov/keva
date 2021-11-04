package console

import (
	"fmt"
	"log"
)

type Console struct {
	levels int
}

func ConsoleLoggerNew(levels int) *Console {
	return &Console{levels: levels}
}

func (c *Console) Log(level int, message string, v ...interface{}) {
	if c.levels&level != 0 {
		log.Println("["+levelToText(level)+"]", fmt.Sprintf(message, v...))
	}
}

func levelToText(level int) string {
	switch level {
	case 1 << 0:
		return "INFO"
	case 1 << 1:
		return "WARN"
	case 1 << 2:
		return "ERROR"
	case 1 << 3:
		return "DEBUG"
	case 1 << 4:
		return "TRACE"
	case 1 << 5:
		return "FATAL"
	}
	return "UNKNOWN"
}
