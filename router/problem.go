package router

import (
	"Online-Practice/helper"
	"Online-Practice/models"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	//根据keyword关键词查找问题
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

// ProblemCreate
// @Tags 管理员私有接口
// @Summary 创建问题
//@Param authorization header string true "authorization"
// @Param title formData string true "title"
// @Param content formData string true "content"
// @Param test_cases formData array false "test_cases"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /admin/problem-create [post]
func ProblemCreate(c *gin.Context) {
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
				"msg":  "测试用例输入格式错误",
			})
			return
		}
		if _, ok := caseMap["output"]; !ok {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "测试用例输出格式错误",
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

// ProblemModify
// @Tags 管理员私有接口
// @Summary 问题修改
//@Param authorization header string true "authorization"
//@Param identity formData string true "identity"
// @Param title formData string true "title"
// @Param content formData string true "content"
// @Param test_cases formData array false "test_cases"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /admin/problem-modify [put]
func ProblemModify(c *gin.Context) {
	identity := c.PostForm("identity")
	title := c.PostForm("title")
	content := c.PostForm("content")
	testCases := c.PostFormArray("test_cases")
	if identity == "" || title == "" || content == "" || len(testCases) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不能为空",
		})
		return
	}
	if err := models.DB.Transaction(func(tx *gorm.DB) error {
		//	问题基本信息保存到Problem
		problem := &models.Problem{
			Title:   title,
			Content: content,
		}
		err := tx.Where("identity = ?", identity).Updates(problem).Error
		if err != nil {
			return err
		}
		//关联测试用例的更新
		//1.删除已存在的关联测试用例
		err = tx.Where("problem_identity = ?", identity).Delete(new(models.TestCase)).Error
		if err != nil {
			return err
		}
		//2.插入新的关联测试用例
		testCase := make([]*models.TestCase, 0)
		for _, v := range testCases {
			caseMap := make(map[string]string)
			err := json.Unmarshal([]byte(v), &caseMap)
			if err != nil {
				return errors.New("测试用例格式错误")
			}
			if _, ok := caseMap["input"]; !ok {
				return errors.New("测试用例input格式错误")
			}
			if _, ok := caseMap["output"]; !ok {
				return errors.New("测试用例output格式错误")
			}
			newTestcase := &models.TestCase{
				Identity:        helper.GetUUID(),
				ProblemIdentity: identity,
				Input:           caseMap["input"],
				Output:          caseMap["output"],
			}
			testCase = append(testCase, newTestcase)
		}
		err = tx.Create(testCase).Error
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "问题修改失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "问题修改成功",
	})
}

// ProblemDelete
// @Tags 管理员私有接口
// @Summary 问题删除
//@Param authorization header string true "authorization"
//@Param identity query string true "identity"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /admin/problem-delete [delete]
func ProblemDelete(c *gin.Context) {
	identity := c.Query("identity")
	if identity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不能为空",
		})
		return
	}
	err := models.DB.Where("identity= ?", identity).Delete(&models.Problem{}).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "问题删除出错",
		})
		log.Println(err)
		return
	}
	//删除已存在的关联测试用例
	err = models.DB.Where("problem_identity = ?", identity).Delete(&models.TestCase{}).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "关联测试用例删除出错" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "问题删除成功",
	})
}
