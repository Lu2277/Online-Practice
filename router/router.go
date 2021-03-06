package router

import (
	_ "Online-Practice/docs"
	"Online-Practice/middlewares"
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
	//公共接口
	//swagger配置
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//问题列表
	r.GET("/problem-list", GetProblem)
	//提交列表
	r.GET("/submit-list", GetSubmit)
	//用户列表
	r.GET("/user-list", GetUser)
	//用户注册
	r.POST("/register", Register)
	//发送验证码
	r.POST("/send-code", SendCode)
	//用户登陆
	r.POST("/login", Login)
	//用户排行榜
	r.GET("/rank-list", GetRankList)

	//管理员私有接口
	//问题创建 先检查是否为管理员，再创建/修改/删除问题
	admin := r.Group("/admin", middlewares.AdminCheck())
	{
		admin.POST("/problem-create", ProblemCreate)
		admin.PUT("/problem-modify", ProblemModify)
		admin.DELETE("/problem-delete", ProblemDelete)
	}
	//r.POST("/problem-create", middlewares.AdminCheck(), ProblemCreate)
	//r.PUT("/problem-modify", middlewares.AdminCheck(), ProblemModify)
	//r.DELETE("/problem-delete", middlewares.AdminCheck(), ProblemDelete)

	return r
}
