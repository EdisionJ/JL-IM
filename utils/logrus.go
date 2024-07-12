package utils

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"path"
)

const (
	red    = 31
	yellow = 33
	blue   = 36
	gray   = 37
)

type LogFormatter struct {
	Prefix string
}

func (t LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var level_color int
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		level_color = gray
	case logrus.WarnLevel:
		level_color = yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		level_color = red
	default:
		level_color = blue
	}

	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	//自定义日期格式
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	if entry.HasCaller() {
		//自定义文件路径
		funcVal := entry.Caller.Function
		fileVal := fmt.Sprintf("%s:%d", path.Base(entry.Caller.File), entry.Caller.Line)
		//自定义输出格式
		fmt.Fprintf(b, "%s [%s] \x1b[%dm[%s]\x1b[0m %s %s %s\n", t.Prefix, timestamp, level_color, entry.Level, fileVal, funcVal, entry.Message)
	} else {
		fmt.Fprintf(b, "%s [%s] \x1b[%dm[%s]\x1b[0m %s\n", t.Prefix, timestamp, level_color, entry.Level, entry.Message)
	}
	return b.Bytes(), nil
}
