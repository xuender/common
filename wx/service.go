package wx

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/kataras/iris"
	"github.com/lithammer/shortuuid"
	"github.com/xuender/toolkit"
)

const (
	tokenKey     = "TOKEN"
	ticketKey    = "TICKET"
	getTicketURL = "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi"
	getTokenURL  = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
)

// Service js服务
type Service struct {
	Cache  *toolkit.Cache
	Config *Config
	TS     *TokenService
}

// NewService 新建js服务
func NewService(ts *TokenService, config *Config) *Service {
	return &Service{
		Cache:  toolkit.NewCache(time.Second * 7200),
		Config: config,
		TS:     ts,
	}
}

// GetToken 获取授权
func (s *Service) getToken() (*Token, error) {
	url := fmt.Sprintf(getTokenURL, s.Config.Appid, s.Config.AppSecret)
	fmt.Println(url)
	var token Token
	err := s.TS.get(url, &token)
	return &token, err
}

// GetToken 获取Token
func (s *Service) GetToken() (string, error) {
	if token, ok := s.Cache.Get(tokenKey); ok {
		return token.(string), nil
	}
	token, err := s.getToken()
	if err != nil {
		return "", err
	}
	if token.ErrCode > 0 {
		return "", errors.New(token.ErrMsg)
	}
	s.Cache.Set(tokenKey, token.AccessToken)
	return token.AccessToken, nil
}

// GetTicket 获取JS-SDK使用权限
func (s *Service) GetTicket() (string, error) {
	// 读缓存
	if ticket, ok := s.Cache.Get(ticketKey); ok {
		return ticket.(string), nil
	}
	// 获取Token
	token, err := s.GetToken()
	if err != nil {
		return "", err
	}
	// 获取ticket
	url := fmt.Sprintf(getTicketURL, token)
	var ticket Ticket
	err = s.TS.get(url, &ticket)
	if err != nil {
		return "", err
	}
	if ticket.ErrCode > 0 {
		return "", errors.New(ticket.ErrMsg)
	}
	// 保存缓存
	s.Cache.Set(ticketKey, ticket.Ticket)
	return ticket.Ticket, err
}

// GetConfig 获取jssdk配置
func (s *Service) GetConfig(uri string) (*JsConfig, error) {
	ticket, err := s.GetTicket()
	if err != nil {
		return nil, err
	}
	nonceStr := shortuuid.New()
	timestamp := time.Now().Unix()
	sigStr := Signature(fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s", ticket, nonceStr, timestamp, uri))
	config := &JsConfig{
		Debug:     false,
		Appid:     s.Config.Appid,
		NonceStr:  nonceStr,
		Timestamp: timestamp,
		Signature: sigStr,
		// TODO 缺少付费接口
		// JsAPIList: []string{"checkJsApi", "updateAppMessageShareData", "updateTimelineShareData", "onMenuShareWeibo", "onMenuShareQZone"},
		JsAPIList: []string{"checkJsApi", "updateAppMessageShareData", "updateTimelineShareData"},
	}
	return config, nil
}

// Signature sha1签名
func Signature(params ...string) string {
	sort.Strings(params)
	h := sha1.New()
	for _, s := range params {
		io.WriteString(h, s)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Party 分组
func (s *Service) Party(p iris.Party) {
	p.Get("/config", func(ctx iris.Context) {
		url := ctx.URLParam("url")
		config, err := s.GetConfig(url)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Text(err.Error())
			return
		}
		ctx.JSON(config)
	})
}
