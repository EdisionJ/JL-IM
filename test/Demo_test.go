package test

import (
	"fmt"
	"github.com/GUAIK-ORG/go-snowflake/snowflake"
	"log"
	"testing"
)

func TestDemo(t *testing.T) {
	s, err := snowflake.NewSnowflake(int64(0), int64(0))
	if err != nil {
		log.Fatal("IN generator init fail! ")
		return
	}

	for true {
		fmt.Println(s.NextVal())
	}
}
