package utils

import (
	"IM/globle"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"time"
)

type Claim struct {
	UID int64
	jwt.RegisteredClaims
}

var secret = []byte(globle.Project)

func GenToken(uid int64) (tokenStr string, err error) {
	c := Claim{
		UID: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    globle.Project,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(viper.GetInt("jwt.valid_time")))),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	tokenStr, err = token.SignedString(secret)
	return
}

func ParseToken(tokenStr string) (*Claim, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claim{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claim); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
