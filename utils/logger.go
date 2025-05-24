package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

type Logger struct {
	level LogLevel
	file  *os.File
}

var globalLogger *Logger

func InitLogger(logFile string, level LogLevel) error {
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("创建日志文件失败: %v", err)
	}

	globalLogger = &Logger{
		level: level,
		file:  file,
	}

	return nil
}

func (l *Logger) log(level LogLevel, msg string, args ...interface{}) {
	if level < l.level {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMsg := fmt.Sprintf("[%s] %s: %s", timestamp, levelToString(level), fmt.Sprintf(msg, args...))

	// 写入文件
	if l.file != nil {
		l.file.WriteString(logMsg + "\n")
	}

	// 控制台输出（带颜色）
	switch level {
	case DEBUG:
		color.HiBlack(logMsg)
	case INFO:
		color.White(logMsg)
	case WARN:
		color.Yellow(logMsg)
	case ERROR:
		color.Red(logMsg)
	}
}

func levelToString(level LogLevel) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// 全局日志函数
func Debug(msg string, args ...interface{}) {
	if globalLogger != nil {
		globalLogger.log(DEBUG, msg, args...)
	}
}

func Info(msg string, args ...interface{}) {
	if globalLogger != nil {
		globalLogger.log(INFO, msg, args...)
	}
}

func Warn(msg string, args ...interface{}) {
	if globalLogger != nil {
		globalLogger.log(WARN, msg, args...)
	}
}

func Error(msg string, args ...interface{}) {
	if globalLogger != nil {
		globalLogger.log(ERROR, msg, args...)
	}
}
