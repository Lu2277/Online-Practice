package main

import (
	"Online-Practice/router"
)

// @title 在线练习系统
// @version 1.0
// @description 基于Gin、Gorm框架的在线练习系统
// @contact.name LL
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8000
func main() {
	r := router.Router()
	r.Run(":8000")
	//var err error
	//dsn := "root:123456@tcp(127.0.0.1:3306)/online-practice?charset=utf8mb4&parseTime=True&loc=Local"
	//db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//if err != nil {
	//	panic(err)
	//}

}
