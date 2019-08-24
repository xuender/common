package common

// APIer 路由器
type APIer interface {
	// API 应用接口
	API() Partyer
}
