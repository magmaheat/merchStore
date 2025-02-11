package app

import (
	"github.com/sirupsen/logrus"
	"os"
)

func setupLogger(level string) {
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logLevel)
	}

	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logrus.SetOutput(os.Stdout)
}
