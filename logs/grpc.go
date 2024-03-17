package logs

import (
	"github.com/mhghw/user-service/config"

	"github.com/sirupsen/logrus"
)

func GRPCLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "time",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
	})
	logger.SetLevel(logrus.WarnLevel)

	if config.IsDevelopment() {
		logger.SetLevel(logrus.DebugLevel)
	}

	return logger
}
