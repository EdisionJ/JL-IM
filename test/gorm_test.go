package test

import (
	"IM/db/model"
	"IM/db/query"
	"IM/service/RR"
	"IM/utils"
	"context"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestGorm(t *testing.T) {
	dsn := "root:123123@tcp(127.0.0.1:3306)/im?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic("failed to connect database")
	}

	newUser := RR.UserSingUp{
		Name:        "jiangshiwei",
		Email:       "jiangshiwei76@163.com",
		PhoneNumber: "17602832214",
		PassWD:      "123123",
		RePassWD:    "123123",
	}

	UserQ := query.Use(db).User
	_, err = UserQ.WithContext(context.Background()).Where(UserQ.PhoneNumber.Eq(newUser.PhoneNumber)).Or(UserQ.Name.Eq(newUser.Name)).First()
	if err != nil && err.Error() != "record not found" {
		println(err)
	}
	var NewUser model.User
	NewUser.Name = newUser.Name
	NewUser.PassWd = utils.EncodeWithSHA256(newUser.PassWD)
	NewUser.PhoneNumber = newUser.PhoneNumber
	if err := UserQ.WithContext(context.Background()).Create(&NewUser); err != nil {
		println(err)
	}
}
