package logger

import (
	"fmt"
	"io"
	"reflect"
	"sync"
	
	"github.com/sirupsen/logrus"
)

type WriterMap map[logrus.Level]io.Writer

type LevelHook struct {
	writers   WriterMap
	levels    []logrus.Level
	lock      *sync.Mutex
	formatter logrus.Formatter
}

func NewHook(output any, formatter logrus.Formatter) *LevelHook {
	hook := &LevelHook{
		lock: new(sync.Mutex),
	}
	hook.SetFormatter(formatter)
	switch output.(type) {
	case WriterMap:
		hook.writers = output.(WriterMap)
		for level := range output.(WriterMap) {
			hook.levels = append(hook.levels, level)
		}
	default:
		panic(fmt.Sprintf("unsupported level map type: %v", reflect.TypeOf(output)))
	}
	return hook
}

func (hook *LevelHook) SetFormatter(formatter logrus.Formatter) {
	hook.lock.Lock()
	defer hook.lock.Unlock()
	if formatter == nil {
		formatter = &logrus.JSONFormatter{}
	}
	hook.formatter = formatter
}

func (hook *LevelHook) Fire(entry *logrus.Entry) error {
	hook.lock.Lock()
	defer hook.lock.Unlock()
	if hook.writers != nil {
		return hook.ioWrite(entry)
	}
	return nil
}

func (hook *LevelHook) ioWrite(entry *logrus.Entry) error {
	writer, ok := hook.writers[entry.Level]
	if !ok {
		panic("no writer for the level")
	}
	msg, err := hook.formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = writer.Write(msg)
	return err
}

func (hook *LevelHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
