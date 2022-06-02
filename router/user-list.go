package router

import (
	"Online-Practice/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetUser
// @Tags 公共接口
// @Summary 用户详情
// @Param identity query string false "用户标识identity"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /user-list [get]
func GetUser(c *gin.Context) {
	identity := c.Query("identity")
	if identity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户标识不能为空",
		})
		return
	}
	data := make([]*models.User, 0)
	err := models.DB.Omit("password").Where("identity =?", identity).Find(&data).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get UserList by identity:" + identity + " Error" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data,
		"msg":  "ok",
	})
}
