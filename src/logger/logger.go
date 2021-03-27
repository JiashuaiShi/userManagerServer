package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"zego.com/userManageServer/src/models"
)

var logger *logrus.Logger

type Field = logrus.Fields

// 初始化日志配置
func Init(conf models.LogConfig) (err error) {
	logger = logrus.New()

	// 设置级别
	logger.SetLevel(models.LogLevelMap[conf.Level])

	// 设置格式
	logger.SetFormatter(&logrus.JSONFormatter{})

	// 设置路径
	var file *os.File
	// path为空，输出到控制台上
	if len(conf.Path) == 0 {
		logger.SetOutput(os.Stdout)
		fmt.Printf("log output stdout!")
		return err
	} else {
		//path不为空，输出到对应路径文件
		if file, err = os.OpenFile(conf.Path, os.O_CREATE|os.O_WRONLY, 0666); err != nil {
			fmt.Printf("log file open failed! program cannnot continue!")
			return err
		}
		logger.Out = file
	}

	return err
}

func Debug(fields logrus.Fields, msg string, args ...interface{}) {
	if fields != nil {
		logger.WithFields(fields).Debug(msg, args)
	} else {
		logger.Debugf(msg, args...)
	}
}

func Info(fields Field, msg string, args ...interface{}) {
	if fields != nil {
		logger.WithFields(fields).Infof(msg, args...)
	} else {
		logger.Infof(msg, args...)
	}
}

func Error(fields Field, msg string, args ...interface{}) {
	if fields != nil {
		logger.WithFields(fields).Errorf(msg, args)
	} else {
		logger.Errorf(msg)
	}
}
