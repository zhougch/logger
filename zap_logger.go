package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

type LogConfig struct {
	Develop	bool	`json:"develop"`
	Level string	`json:"level"`
	Structured bool `json:"structured"`
	Path string 	`json:"path"`
	ErrorPath string	`json:"errorPath"`
	MaxFileSize int	`json:"maxFileSize"`  // megabytes
	MaxBackups int	`json:"maxBackups"`
}

func NewRotateLogger(logConfig LogConfig) *zap.Logger {
	var zapLogger *zap.Logger

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	if logConfig.Develop {
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	var encoder zapcore.Encoder
	if logConfig.Structured {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	}

	hook := lumberjack.Logger{
		Filename:   logConfig.Path,
		MaxSize:    logConfig.MaxFileSize,
		MaxBackups: logConfig.MaxBackups,
	}
	errorHook := lumberjack.Logger{
		Filename:   logConfig.ErrorPath,
		MaxSize:    logConfig.MaxFileSize,
		MaxBackups: logConfig.MaxBackups,
	}

	fileWriter := zapcore.AddSync(&hook)
	errorFileWriter := zapcore.AddSync(&errorHook)
	consoleWriter := zapcore.Lock(os.Stdout)

	zapLevel := new(zapcore.Level)
	zapLevel.Set(logConfig.Level)
	levelEnable := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= *zapLevel
	})

	errLevelEnable := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	var core zapcore.Core
	if logConfig.Develop {
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, consoleWriter, levelEnable),
			zapcore.NewCore(encoder, fileWriter, levelEnable),
			zapcore.NewCore(encoder, errorFileWriter, errLevelEnable),
		)
	} else {
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, fileWriter, levelEnable),
			zapcore.NewCore(encoder, errorFileWriter, errLevelEnable),
		)
	}

	zapLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.DPanicLevel))

	return zapLogger
}
