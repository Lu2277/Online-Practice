package router

import (
	"Online-Practice/helper"
	"Online-Practice/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

// Login
// @Tags 公共接口
// @Summary 用户登录
// @Param username formData string false "用户名称username"
// @Param password formData string false "用户密码password"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /login [post]
func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户名、密码不能为空",
		})
		return
	}
	//获取密码的md5格式
	password = helper.GetMd5(password)
	//print(username, password)
	data := new(models.User)
	err := models.DB.Where("name= ? AND password = ?", username, password).First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "用户名或者密码错误",
			})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "Get User Error:" + err.Error(),
			})
			return
		}
	}
	token, err := helper.GenerateToken(data.Identity, data.Name)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "生成token错误：" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": token,
		"msg":  "ok",
	})
}
