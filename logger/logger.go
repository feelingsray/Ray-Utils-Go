package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"

	"github.com/feelingsray/ray-utils-go/v2/tools"
)

const (
	DebugLevel = logrus.DebugLevel
	InfoLevel  = logrus.InfoLevel
	WarnLevel  = logrus.WarnLevel
	ErrorLevel = logrus.ErrorLevel
)

func ConsoleHandle(level logrus.Level) *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(level)
	return logger
}

func MultiFileHandle(dir, name string, level logrus.Level) (*logrus.Logger, error) {
	isExist, err := tools.PathExists(dir)
	if err != nil {
		return nil, err
	}
	if !isExist {
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			return nil, err
		}
	}
	path := filepath.Join(dir, name)
	logs, _ := rotatelogs.New(fmt.Sprintf("%s%s", path, ".%Y%m%d%H"),
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(12*time.Hour),
		rotatelogs.WithRotationTime(time.Hour),
	)
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(level)
	logger.SetOutput(logs)
	return logger, nil
}
