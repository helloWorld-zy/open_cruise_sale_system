package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap.Logger with additional methods
type Logger struct {
	*zap.SugaredLogger
}

// Field represents a log field
type Field = zap.Field

// New creates a new logger
func New(level string) *Logger {
	config := zap.NewDevelopmentConfig()

	switch level {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	zapLogger, err := config.Build()
	if err != nil {
		panic(err)
	}

	return &Logger{
		SugaredLogger: zapLogger.Sugar(),
	}
}

// Error creates an error field
func Error(err error) Field {
	return zap.Error(err)
}
