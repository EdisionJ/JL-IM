package test

import (
	"IM/utils"
	"fmt"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	str1, _ := utils.GenToken(1)
	fmt.Println(str1)
	token1, _ := utils.ParseToken(str1)
	fmt.Println(token1)

	str2, _ := utils.GenToken(1)
	fmt.Println(str2)
	time.Sleep(time.Second * 2)
	token2, err := utils.ParseToken(str2)
	fmt.Println(token2)
	fmt.Println(err)

}
