package models

import (
	"gorm.io/gorm"
)

type Problem struct {
	gorm.Model
	Identity   string `gorm:"column:identity;type:varchar(36);" json:"identity"`        // 问题的唯一标识
	CategoryID string `gorm:"column:category_id;type:varchar(255);" json:"category_id"` // 分类ID，以逗号分隔
	Title      string `gorm:"column:title;type:varchar(255);" json:"title"`             // 问题的标题
	Content    string `gorm:"column:content;type:text;" json:"content"`                 // 问题的内容
}

func (table *Problem) TableName() string {
	return "problem"
}
func GetProblemList(keyword string) *gorm.DB {

	return DB.Model(new(Problem)).Where("title like ?", "%"+keyword+"%")
	//data := make([]*Problem, 0)
	//DB.Find(&data)
	//for _, v := range data {
	//	fmt.Printf("Problem===>%v \n", v)
	//}
	//return DB
}
