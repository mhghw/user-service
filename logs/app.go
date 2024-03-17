package logs

import "github.com/sirupsen/logrus"

func ApplicationLogger() *logrus.Logger {
	return logrus.StandardLogger()
}
