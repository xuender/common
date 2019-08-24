package common

// Code 用户登录身份识别
type Code struct {
	IDModel
	Jwt string `gorm:"size:200" json:"jwt"` // 用户jwt
}
