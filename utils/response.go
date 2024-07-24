package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	Fail = iota
	Success
)

type RepData struct {
	Code    int    `json:"code"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func DefaultRsp(c *gin.Context, code int, ok bool, msg string) {
	if ok {
		c.JSON(code, RepData{
			Code:    Success,
			Success: ok,
			Message: msg,
		})
	} else {
		c.JSON(http.StatusOK, RepData{
			Code:    Fail,
			Success: ok,
			Message: msg,
		})
	}
}

func RspWithData(c *gin.Context, code int, ok bool, msg string, data any) {
	if ok {
		c.JSON(code, RepData{
			Code:    Success,
			Success: ok,
			Message: msg,
			Data:    data,
		})
	} else {
		c.JSON(code, RepData{
			Code:    Fail,
			Success: ok,
			Message: msg,
			Data:    data,
		})
	}
}
