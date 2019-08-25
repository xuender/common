package common

// Config 配置信息
type Config struct {
	Name      string // 应用名称
	Address   string // 监听端口
	Secret    []byte // JWT密钥
	MockID    string // 模拟用户ID
	Develop   bool   // 开发模式
	Redirect  string // 重定向目录
	LogDir    string // 日志目录
	StaticDir string // 静态资源目录
	WxURL     string // 微信首页
}
