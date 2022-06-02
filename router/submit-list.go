package router

import (
	"Online-Practice/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// GetSubmit
// @Tags 公共接口
// @Summary 提交列表
// @Param page query int false "page"
// @Param size query int false "size"
// @Param problem_identity query string false "problem_identity"
// @Param user_identity query string false "user_identity"
// @Param status query int false "status"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /submit-list [get]
func GetSubmit(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		log.Println(err)
		return
	}
	size, err := strconv.Atoi(c.DefaultQuery("size", "20"))
	if err != nil {
		log.Println(err)
		return
	}
	page = (page - 1) * size
	var count int64
	data := make([]*models.Submit, 0)
	//获取关键字
	problemIdentity := c.Query("problem_identity")
	userIdentity := c.Query("user_identity")
	status, _ := strconv.Atoi(c.Query("status"))
	tx := models.GetSubmitList(problemIdentity, userIdentity, status)
	err = tx.Count(&count).Offset(page).Limit(size).Find(&data).Error
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get Submit List error:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"提交列表": data,
			"提交总数": count,
		},
		"msg": "ok",
	})
}
