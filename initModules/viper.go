package initModules

import (
	"github.com/spf13/viper"
	"log"
)

func initViper() {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fail to parse 'config.yml': %v", err)
	}
	log.Println("Viper Init Success.")
}
