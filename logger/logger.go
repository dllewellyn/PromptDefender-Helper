package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func init() {
	if Log != nil {
		return
	}

	var err error
	Log, err = zap.NewProduction()

	if err != nil {
		panic(err)
	}
}
