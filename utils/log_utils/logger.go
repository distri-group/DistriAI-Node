package logs

import (
	"fmt"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	yellow = "\033[33m"
	red    = "\033[31m"
	green  = "\033[32m"
	reset  = "\033[0m"
)

var Logger *zap.Logger

func init() {
	t := time.Now()
	formattedTime := t.Format("2006-01-02_15:04:05")
	logFileName := fmt.Sprintf("./logs/log_%s.txt", formattedTime)
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFileName,
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     28,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(lumberjackLogger),
		zap.InfoLevel,
	)

	Logger = zap.New(core)
}

func Normal(message string) {
	fmt.Println(message)
}

func Warning(message string) {
	fmt.Println(yellow + message + reset)
	Logger.Warn(message)
}

func Error(message string) {
	fmt.Println(red + message + reset)
	Logger.Error(message)
}

func Result(message string) {
	fmt.Println(green + message + reset)
	Logger.Info(message)
}
