package foundation

import (
	"log"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger() *zap.Logger {
	config := zap.NewProductionConfig()
	level := viper.GetInt("logger.level")

	// Validate the level to ensure it fits within the int8 range
	if level < int(zapcore.DebugLevel) || level > int(zapcore.FatalLevel) {
		log.Fatalf("invalid logger level: %d", level)
	}

	config.Level = zap.NewAtomicLevelAt(zapcore.Level(level))

	logger, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	return logger
}
