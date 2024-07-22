package initModules

import (
	"fmt"
)

func init() {
	fmt.Println("*************************** INIT MODULES ***************************")
	initViper()
	initGorm()
	initLogger()
	initRedis()
	initRocketMq()
	fmt.Println("*********************************************************************")
}
