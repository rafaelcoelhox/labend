package logger

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger - interface principal do logger
type Logger interface {
	// Métodos básicos
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)

	// Métodos auxiliares para contextos específicos
	HTTP(method, path string, statusCode int, duration time.Duration, fields ...zap.Field)
	Database(operation, table string, duration time.Duration, fields ...zap.Field)
	Event(eventType, source string, fields ...zap.Field)
	Performance(operation string, duration time.Duration, fields ...zap.Field)

	// Métodos com contexto
	WithFields(fields ...zap.Field) Logger
	WithRequestID(requestID string) Logger
	WithUserID(userID string) Logger

	// Sync flush dos logs
	Sync() error

	// GetZapLogger retorna o *zap.Logger subjacente
	GetZapLogger() *zap.Logger
}

// Config - configuração do logger
type Config struct {
	Level            string
	Environment      string
	EnableCaller     bool
	EnableStacktrace bool
	OutputPaths      []string
	ErrorOutputPaths []string
}

// zapLogger - implementação com zap
type zapLogger struct {
	zap    *zap.Logger
	config Config
}

// New - cria logger para produção
func New() (Logger, error) {
	config := Config{
		Level:            "info",
		Environment:      "production",
		EnableCaller:     true,
		EnableStacktrace: true,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	return NewWithConfig(config)
}

// NewDevelopment - cria logger para desenvolvimento com cores
func NewDevelopment() (Logger, error) {
	config := Config{
		Level:            "debug",
		Environment:      "development",
		EnableCaller:     true,
		EnableStacktrace: false,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	return NewWithConfig(config)
}

// NewWithConfig - cria logger com configuração customizada
func NewWithConfig(config Config) (Logger, error) {
	var zapConfig zap.Config

	if config.Environment == "production" {
		zapConfig = zap.NewProductionConfig()
		zapConfig.EncoderConfig.TimeKey = "timestamp"
		zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		zapConfig.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	} else {
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig = getColoredEncoderConfig()
		zapConfig.Development = true
	}

	// Configurar nível
	level, err := zapcore.ParseLevel(config.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}
	zapConfig.Level = zap.NewAtomicLevelAt(level)

	// Configurar outputs
	zapConfig.OutputPaths = config.OutputPaths
	zapConfig.ErrorOutputPaths = config.ErrorOutputPaths

	// Configurar caller e stacktrace
	zapConfig.DisableCaller = !config.EnableCaller
	zapConfig.DisableStacktrace = !config.EnableStacktrace

	zapLog, err := zapConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	return &zapLogger{
		zap:    zapLog,
		config: config,
	}, nil
}

// getColoredEncoderConfig - configuração com cores para desenvolvimento
func getColoredEncoderConfig() zapcore.EncoderConfig {
	config := zap.NewDevelopmentEncoderConfig()

	// Cores para níveis
	config.EncodeLevel = func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		var coloredLevel string
		switch level {
		case zapcore.DebugLevel:
			coloredLevel = "\033[36mDEBUG\033[0m" // Cyan
		case zapcore.InfoLevel:
			coloredLevel = "\033[32mINFO\033[0m" // Green
		case zapcore.WarnLevel:
			coloredLevel = "\033[33mWARN\033[0m" // Yellow
		case zapcore.ErrorLevel:
			coloredLevel = "\033[31mERROR\033[0m" // Red
		case zapcore.FatalLevel:
			coloredLevel = "\033[35mFATAL\033[0m" // Magenta
		default:
			coloredLevel = level.CapitalString()
		}
		enc.AppendString(coloredLevel)
	}

	// Timestamp colorido
	config.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		timestamp := fmt.Sprintf("\033[90m%s\033[0m", t.Format("15:04:05.000"))
		enc.AppendString(timestamp)
	}

	// Caller colorido
	config.EncodeCaller = func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		coloredCaller := fmt.Sprintf("\033[90m%s\033[0m", caller.TrimmedPath())
		enc.AppendString(coloredCaller)
	}

	// Configurar separadores e formato
	config.ConsoleSeparator = " "
	config.TimeKey = "time"
	config.CallerKey = "caller"
	config.MessageKey = "msg"

	return config
}

// Métodos básicos
func (l *zapLogger) Info(msg string, fields ...zap.Field) {
	l.zap.Info(msg, fields...)
}

func (l *zapLogger) Error(msg string, fields ...zap.Field) {
	l.zap.Error(msg, fields...)
}

func (l *zapLogger) Debug(msg string, fields ...zap.Field) {
	l.zap.Debug(msg, fields...)
}

func (l *zapLogger) Warn(msg string, fields ...zap.Field) {
	l.zap.Warn(msg, fields...)
}

func (l *zapLogger) Fatal(msg string, fields ...zap.Field) {
	l.zap.Fatal(msg, fields...)
}

// Métodos auxiliares para contextos específicos
func (l *zapLogger) HTTP(method, path string, statusCode int, duration time.Duration, fields ...zap.Field) {
	var colorPrefix string
	var level zapcore.Level

	switch {
	case statusCode >= 200 && statusCode < 300:
		colorPrefix = "\033[32m" // Green
		level = zapcore.InfoLevel
	case statusCode >= 300 && statusCode < 400:
		colorPrefix = "\033[33m" // Yellow
		level = zapcore.InfoLevel
	case statusCode >= 400 && statusCode < 500:
		colorPrefix = "\033[93m" // Bright Yellow
		level = zapcore.WarnLevel
	case statusCode >= 500:
		colorPrefix = "\033[31m" // Red
		level = zapcore.ErrorLevel
	default:
		colorPrefix = "\033[37m" // White
		level = zapcore.InfoLevel
	}

	msg := fmt.Sprintf("%sHTTP %s %s\033[0m", colorPrefix, method, path)
	combinedFields := append([]zap.Field{
		zap.String("method", method),
		zap.String("path", path),
		zap.Int("status_code", statusCode),
		zap.Duration("duration", duration),
		zap.String("duration_ms", fmt.Sprintf("%.2fms", float64(duration.Nanoseconds())/1000000)),
	}, fields...)

	l.zap.Log(level, msg, combinedFields...)
}

func (l *zapLogger) Database(operation, table string, duration time.Duration, fields ...zap.Field) {
	msg := fmt.Sprintf("\033[34mDatabase %s on %s\033[0m", operation, table)
	combinedFields := append([]zap.Field{
		zap.String("db_operation", operation),
		zap.String("db_table", table),
		zap.Duration("duration", duration),
		zap.String("duration_ms", fmt.Sprintf("%.2fms", float64(duration.Nanoseconds())/1000000)),
	}, fields...)

	if duration > 100*time.Millisecond {
		l.zap.Warn(msg+" \033[31m(SLOW QUERY)\033[0m", combinedFields...)
	} else {
		l.zap.Debug(msg, combinedFields...)
	}
}

func (l *zapLogger) Event(eventType, source string, fields ...zap.Field) {
	msg := fmt.Sprintf("\033[35mEvent %s from %s\033[0m", eventType, source)
	combinedFields := append([]zap.Field{
		zap.String("event_type", eventType),
		zap.String("event_source", source),
		zap.Time("event_time", time.Now()),
	}, fields...)

	l.zap.Info(msg, combinedFields...)
}

func (l *zapLogger) Performance(operation string, duration time.Duration, fields ...zap.Field) {
	var colorPrefix string
	var level zapcore.Level

	switch {
	case duration < 10*time.Millisecond:
		colorPrefix = "\033[36m" // Cyan
		level = zapcore.DebugLevel
	case duration < 100*time.Millisecond:
		colorPrefix = "\033[32m" // Green
		level = zapcore.InfoLevel
	case duration < 1*time.Second:
		colorPrefix = "\033[33m" // Yellow
		level = zapcore.WarnLevel
	default:
		colorPrefix = "\033[31m" // Red
		level = zapcore.ErrorLevel
	}

	msg := fmt.Sprintf("%sPerformance %s\033[0m", colorPrefix, operation)
	combinedFields := append([]zap.Field{
		zap.String("operation", operation),
		zap.Duration("duration", duration),
		zap.String("duration_ms", fmt.Sprintf("%.2fms", float64(duration.Nanoseconds())/1000000)),
	}, fields...)

	l.zap.Log(level, msg, combinedFields...)
}

// Métodos com contexto
func (l *zapLogger) WithFields(fields ...zap.Field) Logger {
	return &zapLogger{
		zap:    l.zap.With(fields...),
		config: l.config,
	}
}

func (l *zapLogger) WithRequestID(requestID string) Logger {
	return l.WithFields(zap.String("request_id", requestID))
}

func (l *zapLogger) WithUserID(userID string) Logger {
	return l.WithFields(zap.String("user_id", userID))
}

func (l *zapLogger) Sync() error {
	return l.zap.Sync()
}

// GetZapLogger retorna o *zap.Logger subjacente
func (l *zapLogger) GetZapLogger() *zap.Logger {
	return l.zap
}

// Helper functions para criar fields comuns
func String(key, value string) zap.Field {
	return zap.String(key, value)
}

func Int(key string, value int) zap.Field {
	return zap.Int(key, value)
}

func Error(err error) zap.Field {
	return zap.Error(err)
}

func Duration(key string, duration time.Duration) zap.Field {
	return zap.Duration(key, duration)
}

func Any(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}

// GetLogLevel - helper para pegar nível de log do ambiente
func GetLogLevel() string {
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		return level
	}
	return "info"
}
