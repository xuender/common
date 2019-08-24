package common

// Router 路由器
type Router interface {
	// Auth 身份认证
	Auth() Partyer
	// API 应用接口
	API() Partyer
}
