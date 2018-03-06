package sailgo

import (
	"log"
	"os"
)

const (
	LevelTrace = iota
	LevelDebug
	LevelInfo
	LevelWarning
	LevelError
	LevelCritical
)

var SailLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
var Level = LevelTrace

func SetLogger(logger *log.Logger) {
	SailLogger = logger
}

func SetLevel(l int) {
	Level = l
}

func Trace(v ...interface{}) {
	if Level <= LevelTrace {
		SailLogger.Printf("[T] %v\n", v)
	}
}

func Debug(v ...interface{}) {
	if Level <= LevelDebug {
		SailLogger.Printf("[D] %v\n", v)
	}
}

func Info(v ...interface{}) {
	if Level <= LevelInfo {
		SailLogger.Printf("[I] %v\n", v)
	}
}

func Warning(v ...interface{}) {
	if Level <= LevelWarning {
		SailLogger.Printf("[W] %v\n", v)
	}
}

func Error(v ...interface{}) {
	if Level <= LevelError {
		SailLogger.Printf("[E] %v\n", v)
	}
}

func Critical(v ...interface{}) {
	if Level <= LevelCritical {
		SailLogger.Printf("[C] %v\n", v)
	}
}
