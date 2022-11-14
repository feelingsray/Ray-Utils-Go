package logger

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/weekface/mgorus"
	"gopkg.in/mgo.v2"

	"github.com/feelingsray/Ray-Utils-Go/tools"
)

const (
	DebugLevel = logrus.DebugLevel
	InfoLevel  = logrus.InfoLevel
	ErrorLevel = logrus.ErrorLevel
	WarnLevels = logrus.WarnLevel
)

func LoggerConsoleHandle(level logrus.Level) *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(level)
	return logger
}

// 文本输出
func LoggerFileHandle(dir string, name string, level logrus.Level) (*logrus.Logger, error) {
	// 判断文件夹是否存在,
	isExist, err := tools.PathExists(dir)
	if err != nil {
		return nil, err
	}
	if !isExist {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	logger := logrus.New()
	hook := NewHook(filepath.Join(dir, name+".log"))
	logger.AddHook(hook)
	logger.SetLevel(level)
	return logger, nil

}

// MongoDB
func LoggerSimpleMongoHandle(host string, port int, dbname string, collection string, username string, password string, level logrus.Level) (*logrus.Logger, error) {
	logger := logrus.New()
	dialInfo := mgo.DialInfo{}
	dialInfo.Addrs = []string{fmt.Sprintf("%s:%d", host, port)}
	dialInfo.Direct = false
	dialInfo.Username = username
	dialInfo.Password = password
	dialInfo.Source = "admin"
	session, err := mgo.DialWithInfo(&dialInfo)
	if err != nil {
		return nil, err
	}
	c := session.DB(dbname).C(collection)
	hooker := mgorus.NewHookerFromCollection(c)
	logger.AddHook(hooker)
	logger.SetLevel(level)
	return logger, nil
}

func LoggerMongoHandle(session *mgo.Session, dbname string, collection string, level logrus.Level) (*logrus.Logger, error) {
	if session == nil {
		return nil, errors.New("mongo session nil")
	}
	c := session.DB(dbname).C(collection)
	hooker := mgorus.NewHookerFromCollection(c)
	logger := logrus.New()
	logger.AddHook(hooker)
	logger.SetLevel(level)
	return logger, nil
}
