package common

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lithammer/shortuuid"
)

// IDModel 主键UUID模型
type IDModel struct {
	ID        string    `gorm:"primary_key;size:22" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// BeforeCreate 创建实体前生成UUID
func (m *IDModel) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", shortuuid.New())
	return nil
}
