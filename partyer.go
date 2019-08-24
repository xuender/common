package common

import "github.com/kataras/iris"

// Partyer 服务注册
type Partyer interface {
	Party(p iris.Party)
}
