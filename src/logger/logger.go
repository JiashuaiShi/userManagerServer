/**
 * @Author: ryder
 * @Description:
 * @File: log
 * @Version: 1.0.0
 * @Date: 2021/3/24 10:00 下午
 */

package logger

import (
	"fmt"
	//_ "errors"
	"github.com/sirupsen/logrus"
	"os"
)

var logger *logrus.Logger

type Field = logrus.Fields

func Init() {
	logger = logrus.New()

	logger.WithFields(logrus.Fields{
		"animal": "walrus",
		"size":   10,
	}).Info("A group of walrus emerges from the ocean")

	file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		logger.SetFormatter(&logrus.JSONFormatter{})
		logger.Out = file
	} else {
		fmt.Println("log init failed!")
		return
	}
}

func Info(fields Field, msg string, args ...interface{}) {
	logger.WithFields(fields).Infof(msg, args)
}

func Error(fields Field, msg string, args ...interface{}) {
	logger.WithFields(fields).Error(msg, args)
}

func Debug(fields logrus.Fields, msg string, args ...interface{}) {
	logger.WithFields(fields).Debug(msg, args)
}
