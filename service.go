package common

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/kataras/golog"
	"github.com/kataras/iris"
	"gopkg.in/go-playground/validator.v9"
)

// Service 公共服务
type Service struct {
	validate *validator.Validate

	Logger *golog.Logger
}

// NewService 新建公共服务
func NewService() *Service {
	return &Service{validate: validator.New()}
}

// Get 请求
func (s *Service) Get(url string, p interface{}) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if s.Logger != nil {
		s.Logger.Info("get -> " + string(body))
	}
	return json.Unmarshal(body, p)
}

// Int 路径中的整数
func (s *Service) Int(ctx iris.Context, name, errMsg string) int {
	value, err := ctx.Params().GetInt(name)
	if err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{"msg": errMsg})
		return 0
	}
	return value
}

// String 路径中的字符串
func (s *Service) String(ctx iris.Context, name, errMsg string) string {
	value := ctx.Params().GetString(name)
	if value == "" {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{"msg": errMsg})
		return ""
	}
	return value
}

// URLString query中的字符串
func (s *Service) URLString(ctx iris.Context, name, errMsg string) string {
	str := ctx.URLParam(name)
	if str == "" {
		if errMsg != "" {
			ctx.StatusCode(iris.StatusNotFound)
			ctx.JSON(iris.Map{"msg": errMsg})
		}
		return ""
	}
	return str
}

// URLInt query中的整数
func (s *Service) URLInt(ctx iris.Context, name, errMsg string) int {
	id, err := ctx.URLParamInt(name)
	if err != nil {
		if errMsg != "" {
			ctx.StatusCode(iris.StatusNotFound)
			ctx.JSON(iris.Map{"msg": errMsg})
		}
		return -1
	}
	return id
}

// Bind 对象绑定
func (s *Service) Bind(ctx iris.Context, obj interface{}, errMsg string) bool {
	// 对象解析
	if err := ctx.ReadJSON(obj); err != nil {
		ctx.Application().Logger().Error(err)
		ctx.StatusCode(iris.StatusNotAcceptable)
		ctx.JSON(iris.Map{"msg": errMsg})
		return false
	}
	// 实体校验
	if err := s.validate.Struct(obj); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			ctx.Application().Logger().Error(err)
			ctx.StatusCode(iris.StatusNotAcceptable)
			ctx.JSON(iris.Map{"msg": errMsg})
			return false
		}

		ctx.StatusCode(iris.StatusUnprocessableEntity)

		// for _, err := range  {
		// 	fmt.Println()
		// 	fmt.Println(err.Namespace())
		// 	fmt.Println(err.Field())
		// 	fmt.Println(err.StructNamespace())
		// 	fmt.Println(err.StructField())
		// 	fmt.Println(err.Tag())
		// 	fmt.Println(err.ActualTag())
		// 	fmt.Println(err.Kind())
		// 	fmt.Println(err.Type())
		// 	fmt.Println(err.Value())
		// 	fmt.Println(err.Param())
		// 	fmt.Println()
		// }
		ctx.JSON(iris.Map{
			"msg":       "校验错误",
			"validator": newFieldError(err.(validator.ValidationErrors)),
		})
		return false
	}
	return true
}
