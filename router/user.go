package router

import (
	"Online-Practice/helper"
	"Online-Practice/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
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
		"msg":  "登录成功",
	})
}

// SendCode
// @Tags 公共接口
// @Summary 发送验证码
// @Param email formData string true "email"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /send-code [post]
func SendCode(c *gin.Context) {
	email := c.PostForm("email")
	code := helper.CreateRandCode()
	//设置验证码有效期5分钟
	models.RDB.Set(c, email, code, time.Minute*5)
	err := helper.SendCode(email, code)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "发送验证码错误：" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "验证码发送成功",
	})
}

// Register
// @Tags 公共接口
// @Summary 用户注册
// @Param email formData string true "邮箱"
// @Param code formData string true "验证码"
// @Param name formData string true "名字"
// @Param password formData string true "密码"
// @Param phone formData string false "手机号"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /register [post]
func Register(c *gin.Context) {
	email := c.PostForm("email")
	code := c.PostForm("code")
	name := c.PostForm("name")
	password := c.PostForm("password")
	phone := c.PostForm("phone")
	if email == "" || code == "" || name == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "邮箱、验证码、名字、密码不能为空",
		})
	}
	//判断验证码是否正确
	trueCode, err := models.RDB.Get(c, email).Result()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "获取验证码错误",
		})
		return
	}
	if trueCode != code {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "验证码错误",
		})
		return
	}
	//判断邮箱是否已经注册
	var num int64
	err = models.DB.Where("email= ?", email).Model(&models.User{}).Count(&num).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "获取用户数量错误：" + err.Error(),
		})
		return
	}
	if num > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "该邮箱已经注册了",
		})
		return
	}
	//数据插入
	data := &models.User{
		Identity: helper.GetUUID(), //唯一标识码uuid
		Name:     name,
		Password: helper.GetMd5(password), //密码加密为md5格式
		Phone:    phone,
		Email:    email,
	}
	//插入到数据库
	err = models.DB.Create(data).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户创建失败：" + err.Error(),
		})
		return
	}
	//生成token
	token, err := helper.GenerateToken(data.Identity, data.Name)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "生成token错误：" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":  200,
		"token": token,
		"msg":   "ok",
	})
}

// GetRankList 获取排行榜
// @Tags 公共接口
// @Summary 用户排行榜
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /rank-list [get]
func GetRankList(c *gin.Context) {
	var count int64
	data := make([]*models.User, 0)
	err := models.DB.Model(new(models.User)).Order("right_num DESC,submit_num ASC").Count(&count).Find(&data).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "获取排行榜失败：" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":  200,
		"count": count,
		"data":  data,
		"msg":   "ok",
	})
}
