package router

import (
	"Online-Practice/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// GetProblem
// @Tags 公共接口
// @Summary 问题列表
// @Param page query int false "page"
// @Param size query int false "size"
// @Param keyword query string false "keyword"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /problem-list [get]
func GetProblem(c *gin.Context) {
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
	//获取关键字
	keyword := c.Query("keyword")
	data := make([]*models.Problem, 0)
	tx := models.GetProblemList(keyword)
	err = tx.Count(&count).Offset(page).Limit(size).Find(&data).Error
	if err != nil {
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"问题列表": data,
			"问题总数": count,
		},
		"msg": "ok",
	})
}
