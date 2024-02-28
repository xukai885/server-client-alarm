package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server-client-alarm/controllers"
)

func SetUp() *gin.Engine {
	r := gin.New()

	v1 := r.Group("/api")

	v1.GET("/health", controllers.Ping)
	v1.POST("/register", controllers.Register)
	v1.POST("/alive", controllers.Alive)
	v1.POST("/delete", controllers.Delete)
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})

	return r
}
