package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:123123@tcp(127.0.0.1:3306)/im?charset=utf8&parseTime=True&loc=Local"
	gormdb, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic("failed to connect database")
	}

	g := gen.NewGenerator(gen.Config{
		OutPath: "db/query",
		Mode:    gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	g.UseDB(gormdb) // reuse your gorm db

	g.ApplyBasic(g.GenerateAllTable()...)

	// Generate the code
	g.Execute()
}
