package router

import (
	_ "Online-Practice/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	//swagger配置
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//问题列表
	r.GET("/problem-list", GetProblem)
	//提交列表
	r.GET("/submit-list", GetSubmit)
	//用户列表
	r.GET("/user-list", GetUser)
	//用户登陆
	r.GET("/login", GetUser)

	return r
}
