package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
定义返回信息
*/

type ResponseData struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"`
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: CodeSuccess,
		Msg:  CodeSuccess.Msg(),
		Data: data,
	})
}
func ResponseError(c *gin.Context) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: CodeError,
		Msg:  CodeError.Msg(),
		Data: nil,
	})
}
