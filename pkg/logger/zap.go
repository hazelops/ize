package logger

import (
	"errors"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const maxLogLen = 1024

type StandartLogger struct {
	*zap.SugaredLogger
}

func (l *StandartLogger) Write(p []byte) (n int, err error) {
	if len(p) > maxLogLen {
		l.Debugw(string(p[:maxLogLen]), "truncated", true)
		return maxLogLen, errors.New("log line too long")
	}
	l.Debug(string(p))
	return len(p), nil
}

func (l *StandartLogger) Infof(template string, args ...interface{}) {
	l.SugaredLogger.Infof(template, args...)
}

func NewSugaredLogger(logLevel zapcore.Level) StandartLogger {
	config := NewStandardZapConfig(logLevel)

	zapLogger, err := config.Build()
	if err != nil {
		panic(err)
	}

	return StandartLogger{
		SugaredLogger: zapLogger.Sugar(),
	}
}

// NewStandardLogger creates a new zap.Logger based on common configuration
//
// This is intended to be used with zap.ReplaceGlobals() in an application's
// main.go.
func NewStandardLogger(logLevel zapcore.Level) (l *zap.Logger, err error) {
	config := NewStandardZapConfig(logLevel)
	return config.Build()
}

// NewStandardZapConfig returns a sensible [config](https://godoc.org/go.uber.org/zap#Config) for a Zap logger.
func NewStandardZapConfig(logLevel zapcore.Level) zap.Config {
	if logLevel.String() == "info" || logLevel.String() == "debug" {
		return zap.Config{
			Level:       zap.NewAtomicLevelAt(logLevel),
			Development: false,
			Sampling: &zap.SamplingConfig{
				Initial:    100,
				Thereafter: 100,
			},
			Encoding: "console",
			EncoderConfig: zapcore.EncoderConfig{
				TimeKey:        "timestamp",
				LevelKey:       "level",
				NameKey:        "logger",
				CallerKey:      "caller",
				MessageKey:     "message",
				StacktraceKey:  "stacktrace",
				EncodeLevel:    zapcore.LowercaseLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.SecondsDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			},
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		}
	}

	return zap.Config{
		Level:       zap.NewAtomicLevelAt(logLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "console",
	}
}
