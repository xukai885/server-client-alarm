package controllers

import (
	"github.com/gin-gonic/gin"
	"log"
	"server-client-alarm/dao/mysql"
	"server-client-alarm/modules"
	"server-client-alarm/watchClient"
	"time"
)

// Register 注册接口
func Register(c *gin.Context) {
	var clientBody modules.Client
	clientBody.Ip = c.ClientIP()
	if err := c.ShouldBindJSON(&clientBody); err != nil {
		ResponseError(c)
		return
	}
	if clientBody.Id == "" {
		ResponseError(c)
		log.Println("客户端ID为空")
		return
	}
	// 写入到mysql中
	if err := mysql.Register(&clientBody); err != nil {
		if err.Error() == "数据已存在" {
			ResponseSuccess(c, "🏅️")
			return
		}
		ResponseError(c)
		return
	}
	watchClient.ClientListSum = append(watchClient.ClientListSum, &clientBody)
	watchClient.InitClientListSum()
	log.Println("添加监控端", clientBody.Id)
	ResponseSuccess(c, clientBody.Count)
}

// Alive 暴露接受client活着信息
func Alive(c *gin.Context) {
	var clientBody modules.Client
	if err := c.ShouldBindJSON(&clientBody); err != nil {
		ResponseError(c)
		return
	}

	for i, clientSum := range watchClient.ClientListSum {
		if clientSum.Id == clientBody.Id {
			watchClient.ClientListSum[i].Time = time.Now()
			ResponseSuccess(c, "🏅️")
			return
		}
	}

	ResponseError(c)

}
