package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/setarek/pym-particle-microservice/config"
)

const DevelopmentMode = "development"

type Logger interface {
	InitLogger()
	Info(args ...interface{})
	Debug(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
}

type logger struct {
	config     *config.Config
	sugarLogger *zap.SugaredLogger
}

func NewLogger(config *config.Config) *logger {
	return &logger{config: config}
}

func (l *logger) InitLogger() {
	var err error
	config := zap.NewProductionConfig()
	enccoderConfig := zap.NewProductionEncoderConfig()
	zapcore.TimeEncoderOfLayout("Jan _2 15:04:05.000000000")
	config.EncoderConfig = enccoderConfig
	if l.config.GetString("run_mode") == DevelopmentMode {
		config.Development = true
	}

	log, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
	l.sugarLogger = log.Sugar()
}

func (l logger) Info(args ...interface{}) {
	l.sugarLogger.Info(args)
}

func (l logger) Debug(args ...interface{}) {
	l.sugarLogger.Debug(args)
}

func (l logger) Error(args ...interface{}) {
	l.sugarLogger.Error(args)
}

func (l logger) Fatal(args ...interface{}) {
	l.sugarLogger.Fatal(args)
}
