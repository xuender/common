package common

import (
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
	"github.com/xuender/common/wx"
)

// AuthService 身份认证
type AuthService struct {
	DB      *gorm.DB
	CS      *Service
	US      *UserService
	TS      *wx.TokenService
	Creater Creater
	WxURL   string
}

// NewAuthService 新建身份认证服务
func NewAuthService(db *gorm.DB, cs *Service, us *UserService, ts *wx.TokenService, creater Creater) *AuthService {
	db.AutoMigrate(&Code{}) // 临时认证
	return &AuthService{
		DB:      db,
		CS:      cs,
		US:      us,
		TS:      ts,
		Creater: creater,
	}
}

func (s *AuthService) login(ctx iris.Context) {
	var form LoginForm
	if !s.CS.Bind(ctx, &form, "登录信息不全") {
		return
	}
	var user User
	s.DB.Where("phone = ?", form.Phone).First(&user)
	if user.ID == "" {
		// TODO 发送提示信息给这个手机号,一天只发1次
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(iris.Map{"msg": "密码或手机号错误"})
		return
	}
	user.Password = form.Password
	if user.Check() {
		user.Password = s.US.Token(user.ID)
		ctx.JSON(user)
		return
	}
	ctx.StatusCode(iris.StatusUnauthorized)
	ctx.JSON(iris.Map{"msg": "密码或手机号错误"})
}
func (s *AuthService) register(ctx iris.Context) {
	// TODO 临时
	// ctx.StatusCode(iris.StatusUnprocessableEntity)
	// ctx.JSON(iris.Map{"msg": "手机号错误"})
	// return
	// ctx.JSON(tags)
	var form LoginForm
	if !s.CS.Bind(ctx, &form, "注册信息不全") {
		return
	}
	var user User
	s.DB.Where("phone = ?", form.Phone).First(&user)
	if user.ID != "" {
		// TODO 发送提示信息给这个手机号,一天只发1次
		ctx.StatusCode(iris.StatusUnprocessableEntity)
		ctx.JSON(iris.Map{"msg": "手机号错误"})
		return
	}
	user.Phone = form.Phone
	user.Password = form.Password
	s.DB.Create(&user)
	s.Creater.Create(&user)
	// var plan Plan
	// plan.Title = "我的大计划"
	// s.PS.create(&plan, &user)
	user.Password = s.US.Token(user.ID)
	ctx.JSON(user)
}

// LoginForm 用户登录
type LoginForm struct {
	Phone    string `validate:"required"` // 手机
	Password string `validate:"required"` // 密码
}

// code 获取JWT
func (s *AuthService) code(ctx iris.Context) {
	code := ctx.URLParam("code")
	var c Code
	s.DB.Where("id = ?", code).First(&c)
	if c.ID == "" {
		ctx.JSON("")
		ctx.Application().Logger().Info("获取用户jwt失败")
	} else {
		ctx.JSON(c.Jwt)
		s.DB.Delete(&c)
		ctx.Application().Logger().Info("获取用户jwt")
	}
}

// getUser 获取用户信息
func (s *AuthService) getUser(token *wx.Token) (*User, error) {
	var user User
	userInfo, err := s.TS.GetUserInfo(token.AccessToken, token.Openid)
	if err != nil {
		return &user, err
	}
	user.Nick = userInfo.Nickname
	user.Openid = userInfo.Openid
	user.Sex = userInfo.Sex
	user.Province = userInfo.Province
	user.Country = userInfo.Country
	user.City = userInfo.City
	user.Headimgurl = userInfo.HeadImgURL
	user.Unionid = userInfo.Unionid
	return &user, err

}

// 微信授权登录
func (s *AuthService) wxLogin(ctx iris.Context) {
	code := ctx.URLParam("code")
	state := ctx.URLParam("state")
	// 获取token
	if token, err := s.TS.SaveToken(code, strings.HasPrefix(state, "qr")); err == nil {
		ctx.Application().Logger().Info("token " + token.Unionid + " " + token.AccessToken)
		// 查找用户
		var user User
		s.US.Find(token.Unionid, &user)
		ctx.Application().Logger().Info("user id " + user.ID)
		//  用户不存在
		if user.ID == "" {
			ctx.Application().Logger().Info("获取用户信息")
			// 获取用户
			u, err := s.getUser(token)
			if err != nil {
				ctx.Application().Logger().Error("获取用户信息失败", err)
				ctx.JSON(iris.Map{"msg": err.Error()})
				return
			}
			user = *u
			ctx.Application().Logger().Info("用户信息", user)
			s.DB.Create(user)
			s.Creater.Create(&user)
			// var plan Plan
			// plan.Title = "我的大计划"
			// s.PS.create(&plan, user)
			ctx.Application().Logger().Info("创建用户", user)
		}
		// 生成JWT
		user.Password = s.US.Token(user.ID)
		var code Code
		code.Jwt = user.Password
		s.DB.Create(&code)
		ctx.Application().Logger().Info("生成JWT", code.ID)
		// 登录
		ctx.Redirect(s.WxURL + code.ID)
	} else {
		ctx.Application().Logger().Error("Token获取失败", err)
		ctx.JSON(iris.Map{"msg": err.Error()})
	}
}

// Party 身份认证分组
func (s *AuthService) Party(p iris.Party) {
	// 查询当前计划标签
	p.Post("/login", s.login)
	// 用户注册
	p.Post("/register", s.register)
	p.Get("/wx", s.wxLogin)
	p.Get("/code", s.code)
}
