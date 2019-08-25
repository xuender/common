package wx

// CommonError 错误代码
type CommonError struct {
	ErrCode int    `gorm:"-" json:"errcode"`
	ErrMsg  string `gorm:"-" json:"errmsg"`
}
