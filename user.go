package common

import (
	"bytes"

	"github.com/jinzhu/gorm"
	"github.com/lithammer/shortuuid"
	"github.com/spaolacci/murmur3"
)

// User 用户
type User struct {
	IDModel
	IssuesNum  int    `gorm:"default:0" json:"-"`            // 问题数量
	Nick       string `gorm:"size:60" json:"nickname"`       // 昵称
	Openid     string `gorm:"size:36" json:"openid"`         // 微信ID
	Unionid    string `gorm:"unique;size:36" json:"unionid"` // 微信唯一标识
	Sex        int    `json:"sex"`                           // 性别
	Province   string `gorm:"size:36" json:"province"`       // 省份
	City       string `gorm:"size:36" json:"city"`           // 城市
	Country    string `gorm:"size:16" json:"country"`        // 国家
	Headimgurl string `gorm:"size:250" json:"headimgurl"`    // 头像
	// 用户可修改
	Phone    string `gorm:"size:20" json:"phone"` // 手机
	Password string `gorm:"-" json:"token"`       // 密码
	Cipher   []byte `gorm:"size:16" json:"-"`     // 密文
}

// Encode 加密
func (u *User) Encode(id string) []byte {
	h := murmur3.New128()
	bs := []byte(id)
	h.Write(bs)
	h.Write([]byte(u.Password))
	h.Write(bs)
	return h.Sum(nil)
}

// Check 密码检查
func (u *User) Check() bool {
	return bytes.Equal(u.Cipher, u.Encode(u.ID))
}

// BeforeCreate 创建用户前加密密码
func (u *User) BeforeCreate(scope *gorm.Scope) error {
	id := shortuuid.New()
	scope.SetColumn("ID", id)
	scope.SetColumn("Cipher", u.Encode(id))
	return nil
}
