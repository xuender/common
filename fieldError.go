package common

import "gopkg.in/go-playground/validator.v9"

// FieldError 字段错误信息
type FieldError struct {
	Field string `json:"field"` // 字段
	Tag   string `json:"tag"`   // 标签
	Param string `json:"param"` // 参数
}

func newFieldError(fes []validator.FieldError) []FieldError {
	var ret []FieldError
	for _, err := range fes {
		ret = append(ret, FieldError{Field: err.Field(), Tag: err.Tag(), Param: err.Param()})
	}
	return ret
}
