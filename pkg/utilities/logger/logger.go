package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func BuildLogger(logFilepath string) *zap.SugaredLogger {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(config)
	fileEncoder := zapcore.NewJSONEncoder(config)
	logFile, _ := os.OpenFile(logFilepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 06666)

	tee := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(logFile), zap.DebugLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.DebugLevel),
	)
	var zapLogging = zap.New(
		tee,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	var logger = zapLogging.Sugar()
	return logger
}
