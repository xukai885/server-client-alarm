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

func Delete(c *gin.Context) {
	var clientBody modules.Client
	if err := c.ShouldBindJSON(&clientBody); err != nil {
		ResponseError(c)
		return
	}
	if err := mysql.Delete(clientBody.Id); err != nil {
		ResponseError(c)
		return
	}

	// 遍历切片以查找并删除具有指定属性的对象
	for i := 0; i < len(watchClient.ClientListSum); i++ {
		if watchClient.ClientListSum[i].Id == clientBody.Id {
			watchClient.ClientListSum = append(watchClient.ClientListSum[:i], watchClient.ClientListSum[i+1:]...)
			i-- // 调整循环变量，以考虑已移除的元素
		}
	}
	for i := 0; i < len(watchClient.LastClientListSum); i++ {
		if watchClient.LastClientListSum[i].Id == clientBody.Id {
			watchClient.LastClientListSum = append(watchClient.LastClientListSum[:i], watchClient.LastClientListSum[i+1:]...)
			i-- // 调整循环变量，以考虑已移除的元素
		}
	}

	ResponseSuccess(c, "🏅")
}
