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

// Find 查找用户,参数是用户指针
func (s *UserService) Find(unionid string, user interface{}) {
	s.DB.Where("unionid = ?", unionid).First(user)
}

// Token 生成token
func (s *UserService) Token(userID string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": userID,
	})
	jwt, _ := token.SignedString(s.Secret)
	return jwt
}

// GetUser 获取用户
func (s *UserService) GetUser(ctx iris.Context, user interface{}) {
	id := s.GetUserID(ctx)
	if id != "" {
		s.DB.Where("id=?", id).First(user)
	}
}

// GetUserID 获取用户ID
func (s *UserService) GetUserID(ctx iris.Context) string {
	return ctx.Values().GetString("userID")
}

// get 获取用户信息
func (s *UserService) get(ctx iris.Context) {
	var user User
	s.GetUser(ctx, &user)
	ctx.JSON(user)
}

// patchNick 修改昵称
func (s *UserService) patchNick(ctx iris.Context) {
	var nick string
	ctx.ReadJSON(&nick)
	var user User
	s.GetUser(ctx, &user)
	user.Nick = nick
	s.DB.Model(user).Update("nick", nick)
	ctx.JSON(user)
}

// patchPassword 修改密码
func (s *UserService) patchPassword(ctx iris.Context) {
	var m map[string]string
	ctx.ReadJSON(&m)
	var user User
	s.GetUser(ctx, &user)
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
