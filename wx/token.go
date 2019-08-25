package wx

import "time"

// Token 授权
type Token struct {
	CommonError
	Openid       string    `gorm:"primary_key;size:36" json:"openid"`
	AccessToken  string    `gorm:"size:512" json:"access_token"`
	ExpiresIn    int       `json:"expires_in"`
	RefreshToken string    `gorm:"size:512" json:"refresh_token"`
	Scope        string    `gorm:"size:20" json:"scope"`
	Unionid      string    `gorm:"size:36" json:"unionid"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
