package wx

// JsConfig js-api配置
type JsConfig struct {
	Debug     bool     `json:"debug"`     // 开启调试模式,调用的所有api的返回值会在客户端alert出来，若要查看传入的参数，可以在pc端打开，参数信息会通过log打出，仅在pc端时才会打印。
	Appid     string   `json:"appId"`     // 必填，公众号的唯一标识
	Timestamp int64    `json:"timestamp"` // 必填，生成签名的时间戳
	NonceStr  string   `json:"nonceStr"`  // 必填，生成签名的随机串
	Signature string   `json:"signature"` // 必填，签名
	JsAPIList []string `json:"jsApiList"` // 必填，需要使用的JS接口列表
}
