package common

// Reply 问题
type Reply struct {
	IDModel
	IssueID int    `json:"issueID"`                 // 问题主键
	UserID  string `gorm:"size:22" json:"userID"`   // 用户主键
	Content string `gorm:"size:250" json:"content"` // 内容
	Echo    string `gorm:"size:250" json:"echo"`    // 回音
}
