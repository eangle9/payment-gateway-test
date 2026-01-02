package middleware

import (
	"context"
	"pg/internal/constant"
	"pg/platform/hlog"
	"time"

	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Logger(logger hlog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			req := c.Request()
			path := req.URL.Path
			query := req.URL.RawQuery
			id := uuid.New().String()

			// Add x-request-id and request-start-time to context
			ctx := context.WithValue(req.Context(), constant.ContextKey("x-request-id"), id)
			ctx = context.WithValue(ctx, constant.ContextKey("request-start-time"), start)
			c.SetRequest(req.WithContext(ctx))

			err := next(c)

			end := time.Now()
			latency := end.Sub(start)

			status := c.Response().Status
			if err != nil {
				// If error handler hasn't been called yet or we want to log the error status
				if he, ok := err.(*echo.HTTPError); ok {
					status = he.Code
				} else {
					status = 500 // Default to 500 if unknown error
				}
			}

			fields := []zapcore.Field{
				zap.Int("status", status),
				zap.String("method", req.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", c.RealIP()),
				zap.String("user-agent", req.UserAgent()),
				zap.Int64("request-latency", latency.Milliseconds()),
			}

			// Log via hlog
			if status >= 500 || err != nil {
				logger.Error(ctx, "Request failed", fields...)
			} else {
				logger.Info(ctx, "Request completed", fields...)
			}

			// Send to Sentry
			if hub := sentryecho.GetHubFromContext(c); hub != nil {
				hub.WithScope(func(scope *sentry.Scope) {
					extras := map[string]interface{}{
						"status":          status,
						"method":          req.Method,
						"path":            path,
						"query":           query,
						"ip":              c.RealIP(),
						"user-agent":      req.UserAgent(),
						"request-latency": latency.Milliseconds(),
						"request-id":      id,
					}
					scope.SetExtras(extras)
					if status >= 500 || err != nil {
						scope.SetLevel(sentry.LevelError)
						hub.CaptureMessage("Request failed")
					} else {
						scope.SetLevel(sentry.LevelInfo)
						hub.CaptureMessage("Request completed")
					}
				})
			}
			return err
		}
	}
}
