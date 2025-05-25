package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *logrus.Logger

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// 自定义格式化器，添加源码信息
type CustomFormatter struct {
	logrus.TextFormatter
}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 获取调用者信息
	if pc, file, line, ok := runtime.Caller(7); ok {
		funcName := runtime.FuncForPC(pc).Name()
		// 简化包名和函数名
		if idx := strings.LastIndex(funcName, "."); idx != -1 {
			funcName = funcName[idx+1:]
		}
		fileName := filepath.Base(file)
		entry.Data["source"] = fmt.Sprintf("%s:%d:%s", fileName, line, funcName)
	}

	return f.TextFormatter.Format(entry)
}

func InitLogger(filename string, level LogLevel) error {
	logger = logrus.New()

	// 设置日志级别
	switch level {
	case DEBUG:
		logger.SetLevel(logrus.DebugLevel)
	case INFO:
		logger.SetLevel(logrus.InfoLevel)
	case WARN:
		logger.SetLevel(logrus.WarnLevel)
	case ERROR:
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	// 创建日志目录
	logDir := filepath.Dir(filename)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %v", err)
	}

	// 配置文件输出（带轮转）
	fileWriter := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    100, // MB
		MaxBackups: 5,
		MaxAge:     30, // days
		Compress:   true,
	}

	// 配置多输出：文件和控制台
	logger.SetOutput(fileWriter)

	// 设置自定义格式化器
	logger.SetFormatter(&CustomFormatter{
		TextFormatter: logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			PadLevelText:    true,
		},
	})

	// 添加钩子以同时输出到控制台
	logger.AddHook(&ConsoleHook{})

	Info("日志系统初始化完成")
	return nil
}

// 控制台输出钩子
type ConsoleHook struct{}

func (hook *ConsoleHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}

	// 根据日志级别使用不同颜色输出到控制台
	switch entry.Level {
	case logrus.ErrorLevel:
		fmt.Fprint(os.Stderr, "\033[31m"+line+"\033[0m") // 红色
	case logrus.WarnLevel:
		fmt.Fprint(os.Stdout, "\033[33m"+line+"\033[0m") // 黄色
	case logrus.InfoLevel:
		fmt.Fprint(os.Stdout, "\033[32m"+line+"\033[0m") // 绿色
	case logrus.DebugLevel:
		fmt.Fprint(os.Stdout, "\033[36m"+line+"\033[0m") // 青色
	default:
		fmt.Fprint(os.Stdout, line)
	}
	return nil
}

func (hook *ConsoleHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// 结构化日志方法
func WithFields(fields map[string]interface{}) *logrus.Entry {
	if logger == nil {
		// 如果logger未初始化，使用默认配置
		InitLogger("logs/app.log", INFO)
	}
	return logger.WithFields(fields)
}

func WithError(err error) *logrus.Entry {
	if logger == nil {
		InitLogger("logs/app.log", INFO)
	}
	return logger.WithError(err)
}

// 传统日志方法（保持向后兼容）
func Debug(format string, args ...interface{}) {
	if logger == nil {
		InitLogger("logs/app.log", INFO)
	}
	logger.Debugf(format, args...)
}

func Info(format string, args ...interface{}) {
	if logger == nil {
		InitLogger("logs/app.log", INFO)
	}
	logger.Infof(format, args...)
}

func Warn(format string, args ...interface{}) {
	if logger == nil {
		InitLogger("logs/app.log", INFO)
	}
	logger.Warnf(format, args...)
}

func Error(format string, args ...interface{}) {
	if logger == nil {
		InitLogger("logs/app.log", INFO)
	}
	logger.Errorf(format, args...)
}

// 新增结构化日志方法
func DebugWithFields(msg string, fields map[string]interface{}) {
	WithFields(fields).Debug(msg)
}

func InfoWithFields(msg string, fields map[string]interface{}) {
	WithFields(fields).Info(msg)
}

func WarnWithFields(msg string, fields map[string]interface{}) {
	WithFields(fields).Warn(msg)
}

func ErrorWithFields(msg string, fields map[string]interface{}) {
	WithFields(fields).Error(msg)
}

// 操作追踪日志
func LogOperation(operation string, details map[string]interface{}) {
	fields := map[string]interface{}{
		"operation": operation,
		"module":    "operation",
	}
	for k, v := range details {
		fields[k] = v
	}
	InfoWithFields("执行操作", fields)
}

// 性能日志
func LogPerformance(operation string, duration interface{}, details map[string]interface{}) {
	fields := map[string]interface{}{
		"operation": operation,
		"duration":  duration,
		"module":    "performance",
	}
	for k, v := range details {
		fields[k] = v
	}
	InfoWithFields("性能统计", fields)
}

// 翻译日志
func LogTranslation(source, target, result string, success bool) {
	fields := map[string]interface{}{
		"source":  source,
		"target":  target,
		"result":  result,
		"success": success,
		"module":  "translation",
	}
	if success {
		InfoWithFields("翻译完成", fields)
	} else {
		ErrorWithFields("翻译失败", fields)
	}
}

// 缓存日志
func LogCache(operation string, key string, hit bool, details map[string]interface{}) {
	fields := map[string]interface{}{
		"operation": operation,
		"key":       key,
		"hit":       hit,
		"module":    "cache",
	}
	for k, v := range details {
		fields[k] = v
	}
	InfoWithFields("缓存操作", fields)
}

func Close() {
	if logger != nil {
		// logrus不需要显式关闭
		Info("日志系统关闭")
	}
}
