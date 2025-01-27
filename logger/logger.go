package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type Logger struct {
	instance *zap.Logger
}

// NewLogger initializes a new logger with the given service name, environment, log level, and log file.
// Return a new logger instance and an error if the logger could not be created.
func NewLogger(svcName, env, logLevel, logFile string) (*Logger, error) {
	// Set the encoder (JSON for production, console for development)
	var encoder zapcore.Encoder

	if env == "prod" {
		encoder = zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		})
	} else {
		encoder = zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			TimeKey:      "timestamp",
			LevelKey:     "level",
			CallerKey:    "caller",
			MessageKey:   "message",
			EncodeLevel:  zapcore.CapitalColorLevelEncoder,
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		})
	}

	// Determine the log level
	var atomicLevel zap.AtomicLevel
	switch logLevel {
	case "debug":
		atomicLevel = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		atomicLevel = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		atomicLevel = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		atomicLevel = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	default:
		atomicLevel = zap.NewAtomicLevelAt(zapcore.InfoLevel) // Default to info
	}

	// Set log output destination
	var writeSyncer zapcore.WriteSyncer
	if logFile != "" {
		file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return nil, err
		}
		writeSyncer = zapcore.AddSync(file)
	} else {
		writeSyncer = zapcore.AddSync(os.Stdout) // Default to stdout
	}

	// Create the core
	core := zapcore.NewCore(encoder, writeSyncer, atomicLevel)

	// Build the logger with service-specific fields
	logger := zap.New(core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.Fields(zap.String("service", svcName)),
	)

	return &Logger{instance: logger}, nil
}

// Info logs an info-level message.
func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.instance.Info(msg, fields...)
}

// Error logs an error-level message.
func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.instance.Error(msg, fields...)
}

// Debug logs a debug-level message.
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.instance.Debug(msg, fields...)
}

// Sync flushes the logger's buffered logs.
func (l *Logger) Sync() error {
	return l.instance.Sync()
}
