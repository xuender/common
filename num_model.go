package common

import "time"

// NumModel 主键int模型
type NumModel struct {
	ID        int       `gorm:"primary_key;auto_increment:false" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
