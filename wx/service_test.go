package wx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignature(t *testing.T) {
	assert.Equal(t, "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3", Signature("test"), "摘要")
}

func TestNewService(t *testing.T) {
	// ts := TokenService{}
	// config := Config{Appid: "appid", AppSecret: "appsecret"}
	// js := NewService(&ts, &config)
	// token, err := js.GetToken()
	// assert.Nil(t, err, "JsService GetToken")
	// assert.NotEmpty(t, token, "JsService token")
	// ticket, err := js.GetTicket()
	// assert.Nil(t, err, "JsService GetTicket")
	// assert.NotEmpty(t, ticket, "JsService ticket")
	// jsConfig, err := js.GetConfig("https://anicca.cn/api/wx/config")
	// assert.Nil(t, err, "JsService GetConfig")
	// assert.Equal(t, jsConfig.Appid, config.Appid, "Appid")
	// assert.NotEmpty(t, jsConfig.Signature, "JsService signature")
}
