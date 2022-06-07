package middlewares

import (
	"Online-Practice/helper"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AdminCheck() func(c *gin.Context) {
	return func(c *gin.Context) {
		//	检查用户是否为管理员 通过解析token获取用户信息
		auth := c.GetHeader("Authorization")
		userClaims, err := helper.AnalyseToken(auth)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "解析用户token失败",
			})
			return
		}
		if userClaims.IsAdmin != 1 {
			c.Abort()
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "非管理员，权限不足",
			})
			return
		}
		c.Next()
	}
}
