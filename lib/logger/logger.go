// Package logger provides logging functions
package logger

import (
	"flag"
	"github.com/maxkulish/dnscrypt-list/lib/fs"
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Global variables against which all our logging occurs.
// api-monitor-server --logEnv=prod --logFile=/var/log/api-monitor.log
var (
	zapLogger *zap.Logger

	loggerEnv  = flag.String("logEnv", "dev", "Environment in which api-logger is working")
	loggerFile = flag.String("logFile", "./api-monitor.log", "Path to the file where you want to save logs")
)

// newLogger return development or production logger
// dev logger writes to file and Stderr
// prod logger optimized for speed and write in one file
func newLogger(env string, logFile string) (*zap.Logger, error) {

	if env == "dev" {
		return DevLogger(logFile)
	}

	return ProdLogger(logFile)
}

// DevLogger creates new logger debug level with 2 outputs: file, stdout
func DevLogger(file string) (*zap.Logger, error) {

	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{
		file,
		"stdout",
	}
	// TODO: change timeformat to more readable
	// TODO: set UTC timezone

	return cfg.Build()
}

// ProdLogger creates performance critical logger
func ProdLogger(file string) (*zap.Logger, error) {

	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"stdout",
	}

	return cfg.Build()
}

// SetLogger creates new development logger
func SetLogger() {

	err := fs.CreateFileIfNotExist(*loggerFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	if err != nil {
		log.Fatal(err.Error())
	}

	zapLogger, err = newLogger(*loggerEnv, *loggerFile)
	if err != nil {
		log.Fatal(err.Error())
	}
}

// Debug outputs a message at debug level
func Debug(msg string, fields ...zapcore.Field) {
	zapLogger.Debug(msg, fields...)
}

// Error outputs a message at error level
func Error(msg string, fields ...zapcore.Field) {
	zapLogger.Error(msg, fields...)
}

// Info outputs a message at information level
func Info(msg string, fields ...zapcore.Field) {
	zapLogger.Info(msg, fields...)
}

// Sync flushes any buffered log entries
func Sync() error {
	return zapLogger.Sync()
}
