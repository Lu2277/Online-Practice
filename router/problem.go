package router

import (
	"Online-Practice/helper"
	"Online-Practice/models"
	"encoding/json"
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
	//Offset 从哪一页开始、默认从第一页开始；Limit 限制每页显示的记录数，默认为20条
	err = tx.Count(&count).Offset(page).Limit(size).Find(&data).Error
	if err != nil {
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"问题总数": count,
			"问题列表": data,
		},
		"msg": "ok",
	})
}

// AddProblem
// @Tags 管理员私有接口
// @Summary 创建问题
//@Param authorization header string true "authorization"
// @Param title formData string true "title"
// @Param content formData string true "content"
// @Param test_cases formData array true "test_cases"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /problem-add [post]
func AddProblem(c *gin.Context) {
	title := c.PostForm("title")
	content := c.PostForm("content")
	testCases := c.PostFormArray("test_cases")
	if title == "" || content == "" || len(testCases) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不能为空",
		})
		return
	}
	problemIdentity := helper.GetUUID()
	data := &models.Problem{
		Identity: problemIdentity,
		Title:    title,
		Content:  content,
	}
	//处理测试用例
	testCase := make([]*models.TestCase, 0)
	for _, v := range testCases {
		caseMap := make(map[string]string)
		err := json.Unmarshal([]byte(v), &caseMap)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "测试用例格式错误",
			})
			return
		}
		if _, ok := caseMap["input"]; !ok {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "测试用例格式错误",
			})
			return
		}
		if _, ok := caseMap["output"]; !ok {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "测试用例格式错误",
			})
			return
		}
		newTestcase := &models.TestCase{
			Identity:        helper.GetUUID(),
			ProblemIdentity: problemIdentity,
			Input:           caseMap["input"],
			Output:          caseMap["output"],
		}
		testCase = append(testCase, newTestcase)
	}
	data.TestCases = testCase
	//数据插入、创建问题
	err := models.DB.Create(data).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "问题创建失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":     200,
		"msg":      "问题创建成功",
		"identity": data.Identity,
	})
}
