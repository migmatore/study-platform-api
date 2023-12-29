package logger

import (
	"errors"
	"github.com/migmatore/study-platform-api/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"syscall"
)

type Logger interface {
	InitLogger()
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	DPanic(args ...interface{})
	DPanicf(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
}

type logger struct {
	cfg         *config.Config
	sugarLogger *zap.SugaredLogger
}

var loggerLevelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func (l *logger) getLoggerLevel(cfg *config.Config) zapcore.Level {
	level, exist := loggerLevelMap[cfg.Logger.Level]
	if !exist {
		return zapcore.DebugLevel
	}

	return level
}

func NewLogger(config *config.Config) *logger {
	return &logger{cfg: config}
}

func (l *logger) InitLogger() {
	logLevel := l.getLoggerLevel(l.cfg)

	logWriter := zapcore.AddSync(os.Stdout)

	var encoderCfg zapcore.EncoderConfig
	if l.cfg.Server.Mode == "dev" {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
	}

	encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	encoder := zapcore.NewConsoleEncoder(encoderCfg)

	core := zapcore.NewCore(encoder, logWriter, zap.NewAtomicLevelAt(logLevel))

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	l.sugarLogger = logger.Sugar()
	if err := l.sugarLogger.Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
		l.sugarLogger.Error(err)
	}
}

func (l *logger) Debug(args ...interface{}) {
	l.sugarLogger.Debug(args...)
}

func (l *logger) Debugf(template string, args ...interface{}) {
	l.sugarLogger.Debugf(template, args...)
}

func (l *logger) Info(args ...interface{}) {
	l.sugarLogger.Info(args...)
}

func (l *logger) Infof(template string, args ...interface{}) {
	l.sugarLogger.Infof(template, args...)
}

func (l *logger) Warn(args ...interface{}) {
	l.sugarLogger.Warn(args...)
}

func (l *logger) Warnf(template string, args ...interface{}) {
	l.sugarLogger.Warnf(template, args...)
}

func (l *logger) Error(args ...interface{}) {
	l.sugarLogger.Error(args...)
}

func (l *logger) Errorf(template string, args ...interface{}) {
	l.sugarLogger.Errorf(template, args...)
}

func (l *logger) DPanic(args ...interface{}) {
	l.sugarLogger.DPanic(args...)
}

func (l *logger) DPanicf(template string, args ...interface{}) {
	l.sugarLogger.DPanicf(template, args...)
}

func (l *logger) Panic(args ...interface{}) {
	l.sugarLogger.Panic(args...)
}

func (l *logger) Panicf(template string, args ...interface{}) {
	l.sugarLogger.Panicf(template, args...)
}

func (l *logger) Fatal(args ...interface{}) {
	l.sugarLogger.Fatal(args...)
}

func (l *logger) Fatalf(template string, args ...interface{}) {
	l.sugarLogger.Fatalf(template, args...)
}
