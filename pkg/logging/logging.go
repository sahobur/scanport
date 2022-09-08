package logging

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

type writeHook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

func (hook *writeHook) WriteLog(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	for _, w := range hook.Writer {
		w.Write([]byte(line))
	}
	return err
}

func (hook writeHook) Levels() []logrus.Level {
	return hook.LogLevels
}

func Init() {
	l := logrus.New()
	l.SetReportCaller(true)
	l.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s%d", filename, f.Line)
		},
		DisableColors: false,
		FullTimestamp: true,
	}
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err := os.MkdirAll("logs", 0644)
		if err != nil {
			panic(err)
		}

	}
	logf, err := os.OpenFile("logs/logs.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		panic(err)
	}
	l.SetOutput(io.Discard)
	l.AddHook(&writeHook{
		Writer: []io.Writer{logf, os.Stdout},
		LogLevels: logrus.AllLevels
	})
}
