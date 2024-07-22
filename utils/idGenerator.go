package utils

import (
	"github.com/GUAIK-ORG/go-snowflake/snowflake"
	"log"
)

var s *snowflake.Snowflake

func init() {
	var err error
	s, err = snowflake.NewSnowflake(int64(0), int64(0))
	if err != nil {
		log.Fatal("IN generator init fail! ")
		return
	}
}

func GenID() int64 {
	return s.NextVal()
}
