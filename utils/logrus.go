package utils

import (
	"IM/config"
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path"
)

const (
	red    = 31
	yellow = 33
	blue   = 36
	gray   = 37
)

type LogFormatter struct {
	Config config.LogConf
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
		fmt.Fprintf(b, "%s [%s] \x1b[%dm[%s]\x1b[0m %s %s %s\n", t.Config.Prefix, timestamp, level_color, entry.Level, fileVal, funcVal, entry.Message)
	} else {
		fmt.Fprintf(b, "%s [%s] \x1b[%dm[%s]\x1b[0m %s\n", t.Config.Prefix, timestamp, level_color, entry.Level, entry.Message)
	}
	return b.Bytes(), nil
}

func InitLogger(c config.LogConf) *logrus.Logger {
	logger := logrus.New()                //新建一个实例
	logger.SetOutput(os.Stdout)           //设置输出类型
	logger.SetReportCaller(c.ShowLine)    //开启返回函数名和行号
	logger.SetFormatter(&LogFormatter{c}) //设置自己定义的Formatter
	level, err := logrus.ParseLevel(c.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level) //设置最低的Level
	return logger
}
