package common

// Issue 问题
type Issue struct {
	NumModel
	UserID      string  `gorm:"primary_key;size:22" json:"userID"`                  // 用户主键
	Title       string  `gorm:"size:100;not null" json:"title" validate:"required"` // 标题
	Description string  `gorm:"size:250" json:"description"`                        // 说明
	Replies     []Reply `gorm:"foreignkey:IssueID" json:"-"`                        // 回复
	Echo        string  `gorm:"size:250" json:"echo"`                               // 回音
}
