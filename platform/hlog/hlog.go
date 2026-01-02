package hlog

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

// Logger is the interface for the logger.
type Logger interface {
	// GetZapLogger returns the underlying zap logger.
	GetZapLogger() *zap.Logger

	// Named returns a new named logger.
	Named(s string) *logger

	// With returns a new logger with the given fields.
	With(fields ...zap.Field) *logger

	// Debug logs a debug message.
	Debug(ctx context.Context, msg string, fields ...zap.Field)

	// Info logs an info message.
	Info(ctx context.Context, msg string, fields ...zap.Field)

	// Warn logs a warning message.
	Warn(ctx context.Context, msg string, fields ...zap.Field)

	// Error logs an error message with stack trace.
	Error(ctx context.Context, msg string, fields ...zap.Field)

	// Panic logs a panic message and panics.
	Panic(ctx context.Context, msg string, fields ...zap.Field)

	// Fatal logs a fatal message and exits with os.Exit(1).
	Fatal(ctx context.Context, msg string, fields ...zap.Field)

	// Log is an implementation for pgx logger.
	Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{})

	extract(ctx context.Context) []zap.Field
}

// ExtractField is a struct to define the transformation of a field in the context.
type ExtractField struct {
	KeyInContext any
	Func         func(any) zap.Field
}

// Options is a struct to define the options for the logger.
type Options struct {
	ExtractFields              []ExtractField
	IgnoreDefaultExtractFields bool
	IgnoreDefaultTimeField     bool
}

type logger struct {
	logger  *zap.Logger
	sentry  *sentry.Client // Add Sentry client reference
	options Options
}

// New initializes a new logger with Zap and Sentry integration.
func New(l *zap.Logger, options Options, sentryClient *sentry.Client) Logger {
	if !options.IgnoreDefaultExtractFields {
		options.ExtractFields = append(options.ExtractFields, []ExtractField{
			{
				KeyInContext: "x-request-id",
				Func: func(v any) zap.Field {
					if vString, ok := v.(string); ok {
						return zap.String("x-request-id", vString)
					}
					return zap.Skip()
				},
			},
			{
				KeyInContext: "x-user-id",
				Func: func(v any) zap.Field {
					if vString, ok := v.(string); ok {
						return zap.String("x-user-id", vString)
					}
					return zap.Skip()
				},
			},
			{
				KeyInContext: "request-start-time",
				Func: func(v any) zap.Field {
					if vTime, ok := v.(time.Time); ok {
						return zap.Float64("time-since-request", float64(time.Since(vTime).Milliseconds()))
					}
					return zap.Skip()
				},
			},
			{
				KeyInContext: "x-ws-request-id",
				Func: func(v any) zap.Field {
					if vString, ok := v.(string); ok {
						return zap.String("x-ws-request-id", vString)
					}
					return zap.Skip()
				},
			},
		}...)
	}
	if sentryClient == nil {
		sentryClient = sentry.CurrentHub().Client()
		if sentryClient == nil {
			log.Fatal("sentry client is not initialized")
		}
	}

	return &logger{
		logger:  l,
		sentry:  sentryClient,
		options: options,
	}
}

// GetZapLogger returns the underlying zap logger.
func (l *logger) GetZapLogger() *zap.Logger {
	return l.logger
}

// Named returns a new named logger.
func (l *logger) Named(s string) *logger {
	l2 := l.logger.Named(s)
	return &logger{
		logger:  l2,
		sentry:  l.sentry,
		options: l.options,
	}
}

// With returns a new logger with the given fields.
func (l *logger) With(fields ...zap.Field) *logger {
	l2 := l.logger.With(fields...)
	return &logger{
		logger:  l2,
		sentry:  l.sentry,
		options: l.options,
	}
}

// Debug logs a debug message to Zap only (Sentry typically doesn’t need debug logs).
func (l *logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.With(l.extract(ctx)...).Debug(msg, fields...)
}

// Info logs an info message to both Zap and Sentry.
func (l *logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	combinedFields := append(l.extract(ctx), fields...)
	l.logger.Info(msg, combinedFields...)
	// l.sendToSentry(ctx, sentry.LevelInfo, msg, combinedFields)
}

// Warn logs a warning message to both Zap and Sentry.
func (l *logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	combinedFields := append(l.extract(ctx), fields...)
	l.logger.Warn(msg, combinedFields...)
	l.sendToSentry(ctx, sentry.LevelWarning, msg, combinedFields)
}

// Error logs an error message with stack trace to both Zap and Sentry.
func (l *logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	combinedFields := append(l.extract(ctx), fields...)
	l.logger.Error(msg, combinedFields...)
	l.sendToSentry(ctx, sentry.LevelError, msg, combinedFields)
}

// Panic logs a panic message to both Zap and Sentry, then panics.
func (l *logger) Panic(ctx context.Context, msg string, fields ...zap.Field) {
	combinedFields := append(l.extract(ctx), fields...)
	l.logger.Panic(msg, combinedFields...)
	l.sendToSentry(ctx, sentry.LevelFatal, msg, combinedFields) // Will still send before panic
	panic(msg)
}

// Fatal logs a fatal message to both Zap and Sentry, then exits.
func (l *logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	combinedFields := append(l.extract(ctx), fields...)
	l.logger.Fatal(msg, combinedFields...)
	l.sendToSentry(ctx, sentry.LevelFatal, msg, combinedFields) // Will still send before exit
	os.Exit(1)
}

// sendToSentry sends the log to Sentry if the client is available.
func (l *logger) sendToSentry(ctx context.Context, level sentry.Level, msg string, fields []zap.Field) {
	if l.sentry == nil {
		return // Sentry not initialized yet, skip
	}

	hub := sentry.CurrentHub().Clone()
	hub.WithScope(func(scope *sentry.Scope) {
		extras := make(map[string]interface{})
		for _, field := range fields {
			extras[field.Key] = field.Interface
		}
		scope.SetExtras(extras)
		scope.SetLevel(level)
		hub.CaptureMessage(msg)
	})
}

func (l *logger) extract(ctx context.Context) []zap.Field {
	var fields []zap.Field

	if !l.options.IgnoreDefaultTimeField {
		fields = append(fields, zap.String("time", time.Now().Format(time.RFC3339)))
	}

	if ctx != nil {
		for _, field := range l.options.ExtractFields {
			if v := ctx.Value(field.KeyInContext); v != nil {
				fields = append(fields, field.Func(v))
			}
		}
	}

	return fields
}

// Printf is the kafka logger function implementation.
func (l *logger) Printf(msg string, fields ...interface{}) {
	l.Info(context.Background(), fmt.Sprintf(msg, fields...))
}

// Log is an implementation for pgx logger.
func (l *logger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	fields := make([]zap.Field, 0, len(data))

	data["pgx_time"] = data["time"]
	delete(data, "time")

	for k, v := range data {
		if k == "args" {
			if args, ok := v.([]interface{}); ok {
				var argsStr []string
				for _, arg := range args {
					if argByte, ok := arg.(pgtype.JSON); ok {
						arg = string(argByte.Bytes)
					}
					argsStr = append(argsStr, fmt.Sprintf("%v", arg))
				}
				v = argsStr
			}
		}
		fields = append(fields, zap.Any(k, v))
	}

	switch level {
	case pgx.LogLevelInfo:
		l.Info(ctx, msg, fields...)
	case pgx.LogLevelWarn:
		l.Warn(ctx, msg, fields...)
	case pgx.LogLevelError:
		l.Error(ctx, msg, fields...)
	default:
		l.Debug(ctx, msg, fields...)
	}
}
