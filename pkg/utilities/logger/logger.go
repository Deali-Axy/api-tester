package logger

import (
	strext "api-tester/pkg/utilities/goext/str-ext"
	"fmt"
	"github.com/nleeper/goment"
	"github.com/spf13/afero"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
)

type Options struct {
	WriteToFile bool
	Folder      string
}

func BuildLogger(options *Options) (*zap.SugaredLogger, error) {
	var (
		core   zapcore.Core
		config = zap.NewProductionEncoderConfig()
	)

	config.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	if options.WriteToFile {
		fs := afero.NewOsFs()
		g, _ := goment.New()
		timeStr := g.Format("YYYY-MM-DD_HH-mm-ss-x")

		if strext.IsNullOrWhiteSpace(options.Folder) {
			options.Folder = "logs"
		}
		folder, err := filepath.Abs(options.Folder)
		if err != nil {
			return nil, fmt.Errorf("get fullpath error: %v", err)
		}
		if exists, err := afero.Exists(fs, folder); !exists {
			if err != nil {
				return nil, fmt.Errorf("check path error: %v", err)
			}
			err := os.MkdirAll(folder, os.ModeDir)
			if err != nil {
				return nil, fmt.Errorf("create folder for saving log files failed. err: %v", err)
			}
		}

		path := filepath.Join(folder, fmt.Sprintf("apitester_%s.log", timeStr))
		logFile, _ := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
		fileEncoder := zapcore.NewJSONEncoder(config)
		core = zapcore.NewTee(
			zapcore.NewCore(fileEncoder, zapcore.AddSync(logFile), zap.DebugLevel),
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.DebugLevel),
		)
	} else {
		core = zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.DebugLevel),
		)
	}

	var zapLogging = zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	var logger = zapLogging.Sugar()
	return logger, nil
}
