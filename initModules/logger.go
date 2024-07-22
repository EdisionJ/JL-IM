package initModules

import (
	"IM/globle"
	"IM/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log"
	"os"
)

func initLogger() {
	logger := logrus.New()                                                     //新建一个实例
	logger.SetOutput(os.Stdout)                                                //设置输出类型
	logger.SetReportCaller(viper.GetBool("logger.show_line"))                  //开启返回函数名和行号
	logger.SetFormatter(&utils.LogFormatter{viper.GetString("logger.prefix")}) //设置自己定义的Formatter
	level, err := logrus.ParseLevel(viper.GetString("logger.level"))
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)
	globle.Logger = logger
	log.Println("Logger Init Success")
}
