package models

import "gorm.io/gorm"

type Submit struct {
	Identity        string   `gorm:"column:identity;type:varchar(36);" json:"identity"`                 // 提交的标识
	ProblemIdentity string   `gorm:"column:problem_identity;type:varchar(36);" json:"problem_identity"` // 题目的唯一标识
	ProblemLink     *Problem `gorm:"foreignKey:identity;references:problem_identity"`                   //关联问题表
	UserIdentity    string   `gorm:"column:user_identity;type:varchar(36);" json:"user_identity"`       // 用户的唯一标识
	UserLink        *User    `gorm:"foreignKey:identity;references:user_identity"`                      //关联用户表
	Path            string   `gorm:"column:path;type:varchar(255);" json:"path"`                        // 提交的路径
	Status          int      `gorm:"column:status;type:tinyint(1);" json:"status"`                      // 提交的状态
}

func (table Submit) TableName() string {
	return "submit"
}

// GetSubmitList 获取提交列表
func GetSubmitList(problemIdentity, userIdentity string, status int) *gorm.DB {
	//预加载问题表、用户表
	tx := DB.Model(&Submit{}).Preload("ProblemLink").Preload("UserLink")
	if problemIdentity != "" {
		tx.Where("problem_identity=?", problemIdentity)
	}
	if userIdentity != "" {
		tx.Where("user_identity=?", userIdentity)
	}
	if status != 0 {
		tx.Where("status=?", status)
	}
	return tx
}
