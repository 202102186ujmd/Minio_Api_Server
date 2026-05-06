package logger

import (
	"os"
	"strings"

	"github.com/202102186ujmd/Minio_Api_Server/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(cfg config.Config) (*zap.Logger, error) {
	level := zapcore.InfoLevel
	if err := level.UnmarshalText([]byte(strings.ToLower(cfg.LogLevel))); err != nil {
		level = zapcore.InfoLevel
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.CallerKey = "caller"

	var encoder zapcore.Encoder
	if strings.ToLower(cfg.LogFormat) == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	var sink zapcore.WriteSyncer
	if strings.ToLower(cfg.LogOutput) == "stdout" || cfg.LogOutput == "" {
		sink = zapcore.AddSync(os.Stdout)
	} else {
		file, err := os.OpenFile(cfg.LogOutput, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return nil, err
		}
		sink = zapcore.AddSync(file)
	}

	core := zapcore.NewCore(encoder, sink, level)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return logger, nil
}

func String(key, value string) zap.Field { return zap.String(key, value) }
func Int(key string, value int) zap.Field { return zap.Int(key, value) }
func Int64(key string, value int64) zap.Field { return zap.Int64(key, value) }
func Bool(key string, value bool) zap.Field { return zap.Bool(key, value) }
func Error(err error) zap.Field { return zap.Error(err) }
