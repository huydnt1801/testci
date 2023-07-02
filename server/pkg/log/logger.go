package log

import (
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

func ZapLogger() logr.Logger {
	zapLog := zap.Must(zap.NewDevelopment())
	log := zapr.NewLogger(zapLog)
	return log
}
