package log

import (
	"fibera/internal/config"
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// SetupLogger 初始化日志
func SetupLogger(conf *config.Config) error {
	// 配置 zap logger
	zapConfig := zap.NewProductionConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(getZapLogLevel(conf.Log.Level))
	zapConfig.EncoderConfig.TimeKey = "timestamp"
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapConfig.EncoderConfig.ConsoleSeparator = " | "

	// 根据 format 设置编码器类型
	var encoderConfig zapcore.Encoder
	if conf.Log.Format == "json" {
		encoderConfig = zapcore.NewJSONEncoder(zapConfig.EncoderConfig)
	} else {
		encoderConfig = zapcore.NewConsoleEncoder(zapConfig.EncoderConfig)
	}

	outputDir := "../" + conf.Log.OutputDir
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create log directory: %v", err)
	}

	fileWriter := &lumberjack.Logger{
		Filename:   outputDir + "/app.log",
		MaxSize:    conf.Log.MaxSize,
		MaxBackups: conf.Log.MaxBackups,
		MaxAge:     conf.Log.MaxAge,
		Compress:   true, // 开启压缩
	}

	var cores []zapcore.Core

	// 根据 outputMethod 添加不同输出
	switch conf.Log.OutputMethod {
	case "console":
		cores = append(cores, zapcore.NewCore(
			encoderConfig,
			zapcore.AddSync(os.Stdout),
			zapConfig.Level,
		))
	case "file":
		cores = append(cores, zapcore.NewCore(
			encoderConfig,
			zapcore.AddSync(fileWriter),
			zapConfig.Level,
		))
	default:
		cores = append(cores, zapcore.NewCore(
			encoderConfig,
			zapcore.AddSync(os.Stdout),
			zapConfig.Level,
		))
		cores = append(cores, zapcore.NewCore(
			encoderConfig,
			zapcore.AddSync(fileWriter),
			zapConfig.Level,
		))
	}

	// 创建 logger 实例
	logger := zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))

	// 设置全局 logger
	zap.ReplaceGlobals(logger)
	return nil
}

// getZapLogLevel 获取 zap 日志级别
func getZapLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
