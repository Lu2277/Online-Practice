# 基于Gin、Gorm实现的在线练习系统

## 数据表

> 1.问题表problem
>
>2.用户表user
>
> 3.提交表submit
>

## 整合 Swagger

参考文档： https://github.com/swaggo/gin-swagger
接口访问地址：http://localhost:8080/swagger/index.html

```text
// GetProblemList
// @Tags 公共方法
// @Summary 问题列表
// @Param page query int false "page"
// @Param size query int false "size"
// @Success 200 {string} json "{"code":"200","msg","","data":""}"
// @Router /problem-list [get]
```

## 安装 jwt
