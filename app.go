package common

import (
	"fmt"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/iris-contrib/middleware/cors"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/jinzhu/gorm"
	"github.com/kataras/golog"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
)

// App 应用
type App struct {
	DB    *gorm.DB
	CS    *Service
	AS    *AuthService
	APIer APIer

	Log *golog.Logger
}

// NewApp 新建应用
func NewApp(db *gorm.DB, cs *Service, as *AuthService, apier APIer) *App {
	return &App{
		DB:    db,
		CS:    cs,
		AS:    as,
		APIer: apier,
	}
}

// Run 运行
func (a *App) Run(c *Config) {
	// 应用
	app := iris.Default()
	a.Log = app.Logger()
	// 压缩
	// app.Use(iris.Gzip)
	// 真实 IP
	app.Configure(iris.WithConfiguration(iris.Configuration{
		RemoteAddrHeaders: map[string]bool{
			"X-Real-Ip":       true,
			"X-Forwarded-For": true,
		},
	}))
	// 开发模式
	if c.Develop {
		app.Logger().SetLevel("debug")
		a.DB.LogMode(true)
	} else {
		c.MockID = ""
	}
	// 404 返回首页
	tmpl := iris.HTML(c.StaticDir, ".html")
	tmpl.Reload(true)
	app.RegisterView(tmpl)
	app.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
		// ctx.Redirect(c.Redirect)
		ctx.View("index.html")
	})
	// 日志配置
	if c.LogDir != "" {
		f, err := os.OpenFile(
			fmt.Sprintf("%s/big-plan.log", c.LogDir),
			os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		app.Logger().SetOutput(f)
		a.DB.SetLogger(app.Logger())
		a.CS.Logger = app.Logger()
	}
	// 静态资源
	app.Logger().Info("资源目录: ", c.StaticDir)
	app.StaticWeb("/", c.StaticDir)
	// JWT 安全认证
	jwtHandler := jwtmiddleware.New(jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return c.Secret, nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
	v1 := app.Party("/api", cors.AllowAll()).AllowMethods(iris.MethodOptions)
	if len(c.MockID) < 4 {
		v1.Use(jwtHandler.Serve)
	}
	v1.Use(func(ctx context.Context) {
		if ctx.Values().Get("jwt") != nil {
			j := ctx.Values().Get("jwt").(*jwt.Token)
			m := j.Claims.(jwt.MapClaims)
			ctx.Values().Set("userID", m["id"].(string))
		} else {
			ctx.Values().Set("userID", c.MockID)
		}
		ctx.Next()
	})
	// AUTH
	a.AS.Party(app.Party("/auth", cors.AllowAll()).AllowMethods(iris.MethodOptions))
	// API
	a.APIer.API().Party(v1)
	// 首页
	a.AS.WxURL = c.WxURL
	// 启动
	app.Run(iris.Addr(c.Address))
}
