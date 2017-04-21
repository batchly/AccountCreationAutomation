package logger

import (
	"os"

	"github.com/Sirupsen/logrus"
)

var (
	logger *logrus.Logger
	f      *os.File
	err    error
)

func Set() {
	logger = logrus.New()
	logger.Formatter = new(logrus.JSONFormatter)

	f, err = os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err == nil {
		logger.Out = f
	} else {
		logger.Fatal("Unable to write logs to file. Printing to console.", err)
	}
}

func Get() *logrus.Logger {
	return logger
}

func GetFileHandle() *os.File {
	return f
}
