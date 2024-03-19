package utils

import (
	"log"

	"go.uber.org/zap"
)

var MySugarLogger *zap.SugaredLogger

func init() {

	logger, err := zap.NewProduction()
	defer logger.Sync()
	if err != nil {
		log.Fatalf(err.Error())
	}
	MySugarLogger = logger.Sugar()
}
