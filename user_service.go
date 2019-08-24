package common

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
)

// UserService 用户服务
type UserService struct {
	DB     *gorm.DB
	Secret []byte
}

// NewUserService 新建用户服务
func NewUserService(db *gorm.DB, secret []byte) *UserService {
	// db.AutoMigrate(&User{}) // 用户
	return &UserService{DB: db, Secret: secret}
}

// find 查找用户
func (s *UserService) find(unionid string) *User {
	var user User
	s.DB.Where("unionid = ?", unionid).First(&user)
	return &user
}

// Token 生成token
func (s *UserService) Token(u *User) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": u.ID,
	})
	u.Password, _ = token.SignedString(s.Secret)
}

// GetUser 获取用户
func (s *UserService) GetUser(ctx iris.Context) *User {
	var user User
	id := s.GetUserID(ctx)
	if id != "" {
		s.DB.Where("id=?", id).First(&user)
	}
	return &user
}

// GetUserID 获取用户ID
func (s *UserService) GetUserID(ctx iris.Context) string {
	return ctx.Values().GetString("userID")
}

// get 获取用户信息
func (s *UserService) get(ctx iris.Context) {
	ctx.JSON(s.GetUser(ctx))
}

// patchNick 修改昵称
func (s *UserService) patchNick(ctx iris.Context) {
	var nick string
	ctx.ReadJSON(&nick)
	user := s.GetUser(ctx)
	user.Nick = nick
	s.DB.Model(user).Update("nick", nick)
	ctx.JSON(user)
}

// patchPassword 修改密码
func (s *UserService) patchPassword(ctx iris.Context) {
	var m map[string]string
	ctx.ReadJSON(&m)
	user := s.GetUser(ctx)
	user.Password = m["old"]
	// user.Cipher = user.encode(user.ID)
	if user.Check() {
		user.Password = m["password"]
		s.DB.Model(user).Update("cipher", user.Encode(user.ID))
		ctx.JSON(user)
		return
	}
	ctx.StatusCode(iris.StatusUnauthorized)
	ctx.JSON(iris.Map{"msg": "密码错误"})
}

// Party 分组
func (s *UserService) Party(p iris.Party) {
	// 查询当前用户信息
	p.Get("/", s.get)
	// 修改昵称
	p.Patch("/nick", s.patchNick)
	// 修改密码
	p.Patch("/password", s.patchPassword)
}
