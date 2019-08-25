package wx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jinzhu/gorm"
)

const (
	accessTokenURL        = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	refreshAccessTokenURL = "https://api.weixin.qq.com/sns/oauth2/refresh_token?appid=%s&grant_type=refresh_token&refresh_token=%s"
	userInfoURL           = "https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN"
	checkAccessTokenURL   = "https://api.weixin.qq.com/sns/auth?access_token=%s&openid=%s"
)

// TokenService 授权服务
type TokenService struct {
	DB     *gorm.DB
	Config *Config
}

// NewTokenService 新建授权服务
func NewTokenService(db *gorm.DB, config *Config) *TokenService {
	db.AutoMigrate(&Token{}) // 授权
	return &TokenService{DB: db, Config: config}
}

func (s *TokenService) get(url string, p interface{}) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, p)
}

// FindToken 查找Token
func (s *TokenService) FindToken(openid string) (*Token, error) {
	// 查询数据库
	var token Token
	s.DB.Where("openid = ?", openid).First(&token)
	if token.Openid == "" {
		return nil, errors.New("未保存过Token")
	}
	// 检查
	if ok, err := s.CheckAccessToken(token.AccessToken, token.Openid); ok && err == nil {
		return &token, err
	}
	// 刷新
	return s.RefreshAccessToken(token.Openid, token.RefreshToken)
}

// SaveToken 获取授权
func (s *TokenService) SaveToken(code string, web bool) (*Token, error) {
	url := fmt.Sprintf(accessTokenURL, s.Config.Appid, s.Config.AppSecret, code)
	if web {
		url = fmt.Sprintf(accessTokenURL, s.Config.WebAppid, s.Config.WebAppSecret, code)
	}
	var token Token
	err := s.get(url, &token)
	if err == nil && token.Openid != "" {
		var old Token
		s.DB.Where("openid = ?", token.Openid).First(&old)
		if old.Openid == "" {
			s.DB.Create(&token)
		} else {
			s.DB.Update(&token)
		}
	}
	return &token, err
}

// RefreshAccessToken 更新授权
func (s *TokenService) RefreshAccessToken(openid, refreshToken string) (*Token, error) {
	url := fmt.Sprintf(refreshAccessTokenURL, openid, refreshToken)
	var token Token
	err := s.get(url, &token)
	if err == nil && token.Openid != "" {
		s.DB.Save(&token)
	}
	return &token, err
}

// CheckAccessToken 检验access_token是否有效
func (s *TokenService) CheckAccessToken(accessToken, openid string) (bool, error) {
	url := fmt.Sprintf(checkAccessTokenURL, accessToken, openid)
	var token Token
	err := s.get(url, &token)
	if err == nil && token.ErrCode == 0 {
		return true, nil
	}
	return false, err
}

//GetUserInfo 如果scope为 snsapi_userinfo 则可以通过此方法获取到用户基本信息
func (s *TokenService) GetUserInfo(accessToken, openid string) (*UserInfo, error) {
	url := fmt.Sprintf(userInfoURL, accessToken, openid)
	var user UserInfo
	err := s.get(url, &user)
	if err != nil {
		return nil, err
	}
	if user.ErrCode != 0 {
		return nil, fmt.Errorf("GetUserInfo error : errcode=%v , errmsg=%v", user.ErrCode, user.ErrMsg)
	}
	return &user, nil
}
