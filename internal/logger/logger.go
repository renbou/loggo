package logger

import (
	"fmt"
	"net/url"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// newline is used to split logs by line in our custom log sink and route them to the storage
	newline = zapcore.DefaultLineEnding

	schemeLoggo = "loggo"
	outputLoggo = "loggo://loggo"
)

// LogConsumer is supposed to consume log messages. The default usage routes logs to the storage.
type LogConsumer func(t time.Time, message []byte) error

var (
	globalSink   consumerSink
	globalLogger *zap.SugaredLogger
)

// InitGlobal initializes the global logger using the specified consumer. Note that no synchronization
// is done here, which is why this must be done during the "init" step of the program.
func InitGlobal(consumer LogConsumer, outputs ...string) error {
	globalSink.consumer = consumer

	// rebuild logger with proper output
	logger, err := newLogger(append(outputs, outputLoggo)...)
	if err != nil {
		return fmt.Errorf("creating global logger: %w", err)
	}

	globalLogger = logger
	return nil
}

// Infow logs an informational message to the global logger.
func Infow(msg string, kvs ...any) {
	globalLogger.Infow(msg, kvs...)
}

// Errorw logs an error message to the global logger.
func Errorw(msg string, kvs ...any) {
	globalLogger.Errorw(msg, kvs...)
}

// Sync runs the global logger Sync routine. Meant to be run at the end of the app's execution.
func Sync() error {
	return globalLogger.Sync()
}

func newLogger(outputs ...string) (*zap.SugaredLogger, error) {
	cfg := zap.NewProductionConfig()
	cfg.DisableCaller = true
	cfg.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.LineEnding = newline
	cfg.OutputPaths = outputs

	logger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("building global logger: %w", err)
	}

	return logger.Sugar(), nil
}

func init() {
	globalSink = consumerSink{lineEnding: []byte(newline)}

	if err := zap.RegisterSink(schemeLoggo, func(u *url.URL) (zap.Sink, error) {
		return &globalSink, nil
	}); err != nil {
		panic(fmt.Sprintf("registering custom log sink: %s", err))
	}

	// Since the consumer won't be set at this point, this is basically a noop logger.
	logger, err := newLogger(outputLoggo)
	if err != nil {
		panic(fmt.Sprintf("initializing default global logger: %s", err))
	}

	globalLogger = logger
}
