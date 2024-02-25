package controllers

import (
	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	ResponseSuccess(c, "OK")
}
