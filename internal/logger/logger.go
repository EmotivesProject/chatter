package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	myLogger      *zap.Logger
	mySugarLogger *zap.SugaredLogger
	serviceName   string
)

func InitLogger(name string) {
	atom := zap.NewAtomicLevel()

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	myLogger = zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))

	defer myLogger.Sync()

	mySugarLogger = myLogger.Sugar()

	atom.SetLevel(zap.InfoLevel)

	serviceName = name

	Info("Created Logger")
}

func Info(message string) {
	myLogger.Info(message, zap.String("service", serviceName))
}

func Infof(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	myLogger.Info(message, zap.String("service", serviceName))
}

func Error(err error) {
	myLogger.Error(err.Error(), zap.String("service", serviceName))
}

func Fatal(err error) {
	myLogger.Fatal(err.Error(), zap.String("service", serviceName))
}
