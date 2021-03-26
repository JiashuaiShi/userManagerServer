/**
 * @Author: ryder
 * @Description:
 * @File: log
 * @Version: 1.0.0
 * @Date: 2021/3/24 10:00 下午
 */

package logger

import (
	"errors"
	"github.com/sirupsen/logrus"
	"os"
	"zego.com/userManageServer/src/models"
)

var logger *logrus.Logger

type Field = logrus.Fields

func Init(conf models.LogConfig) (err error) {
	logger = logrus.New()
	// 设置级别
	logger.SetLevel(models.LogLevelMap[conf.Level])
	// 设置格式
	logger.SetFormatter(&logrus.JSONFormatter{})
	// 设置路径
	if len(conf.Path) == 0{
		logger.SetOutput(os.Stdout)
		return err
	}else {
		if file, err := os.OpenFile(conf.Path, os.O_CREATE|os.O_WRONLY, 0666); err == nil {
			logger.Out = file
			return err
		} else {
			return errors.New("log conf failed")
		}
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

func Debug(fields logrus.Fields, msg string, args ...interface{}) {
	if fields != nil {
		logger.WithFields(fields).Debug(msg, args)
	} else {
		logger.Debugf(msg, args...)
	}
}
