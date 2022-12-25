package logger

import (
  "fmt"
  "os"
  "path/filepath"
  "time"
  
  "github.com/lestrrat-go/file-rotatelogs"
  "github.com/rifflock/lfshook"
  "github.com/sirupsen/logrus"
  
  "github.com/feelingsray/Ray-Utils-Go/tools"
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
    rotatelogs.WithMaxAge(24*time.Hour),
    rotatelogs.WithRotationTime(time.Hour),
  )
  wfName := fmt.Sprintf("%s.wf", name)
  wfPath := filepath.Join(dir, wfName)
  wfLogs, _ := rotatelogs.New(fmt.Sprintf("%s%s", wfPath, ".%Y%m%d%H"),
    rotatelogs.WithLinkName(wfPath),
    rotatelogs.WithMaxAge(24*time.Hour),
    rotatelogs.WithRotationTime(time.Hour),
  )
  hook := lfshook.NewHook(lfshook.WriterMap{
    logrus.InfoLevel:  logs,
    logrus.DebugLevel: logs,
    logrus.WarnLevel:  wfLogs,
    logrus.ErrorLevel: wfLogs,
  }, &logrus.JSONFormatter{})
  logger := logrus.New()
  logger.AddHook(hook)
  logger.SetLevel(level)
  return logger, nil
}
