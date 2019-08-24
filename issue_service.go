package common

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
)

// IssueService 问题服务
type IssueService struct {
	DB *gorm.DB
	CS *Service
	US *UserService
}

// NewIssueService 新建问题服务
func NewIssueService(db *gorm.DB, cs *Service, us *UserService) *IssueService {
	db.AutoMigrate(&Issue{}) // 问题
	return &IssueService{DB: db, CS: cs, US: us}
}

func (s *IssueService) getIssue(id int, userID string) *Issue {
	var i Issue
	s.DB.Where("id=? and user_id=?", id, userID).First(&i)
	return &i
}

// create 创建里问题
func (s *IssueService) create(i *Issue, user *User) {
	user.IssuesNum++
	i.ID = user.IssuesNum
	i.UserID = user.ID
	s.DB.Create(i)
	s.DB.Model(user).Update("issues_num", gorm.Expr("issues_num + 1"))
}

// update 修改问题
func (s *IssueService) update(i *Issue, user *User) (*Issue, bool) {
	// 存在判断
	if old := s.getIssue(i.ID, user.ID); old.ID > 0 {
		old.Title = i.Title
		old.Description = i.Description
		s.DB.Omit("created_at", "state").Save(old)
		return old, true
	}
	return nil, false
}

func (s *IssueService) checkIssue(ctx iris.Context) *Issue {
	id, err := ctx.URLParamInt("issue")
	if err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{"msg": "问题ID错误"})
		return nil
	}
	i := s.getIssue(id, s.US.GetUserID(ctx))
	if i.ID == 0 {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{"msg": fmt.Sprintf("问题未找到:%d", id)})
		return nil
	}
	return i
}
func (s *IssueService) list(ctx iris.Context) {
	var is []Issue
	userID := s.US.GetUserID(ctx)
	s.DB.Order("id desc").Where("user_id=?", userID).Find(&is)
	ctx.JSON(is)
}
func (s *IssueService) post(ctx iris.Context) {
	var input Issue
	if s.CS.Bind(ctx, &input, "问题数据错误") {
		var user User
		s.US.GetUser(ctx, &user)
		s.create(&input, &user)
		ctx.JSON(input)
	}
}
func (s *IssueService) put(ctx iris.Context) {
	var input Issue
	if !s.CS.Bind(ctx, &input, "问题数据库错误") {
		return
	}
	input.ID = s.CS.Int(ctx, "id", "问题ID错误")
	if input.ID == 0 {
		return
	}
	var user User
	s.US.GetUser(ctx, &user)
	if i, ok := s.update(&input, &user); ok {
		ctx.JSON(i)
		return
	}
	ctx.StatusCode(iris.StatusUnauthorized)
	ctx.JSON(iris.Map{"msg": "问题未找到"})
}

func (s *IssueService) get(ctx iris.Context) {
	id := s.CS.Int(ctx, "id", "问题ID错误")
	if id == 0 {
		return
	}
	i := s.getIssue(id, s.US.GetUserID(ctx))
	if i.ID == 0 {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{"msg": fmt.Sprintf("问题未找到: %d", id)})
		return
	}
	ctx.JSON(i)
}

func (s *IssueService) patchState(ctx iris.Context) {
	id := s.CS.Int(ctx, "id", "问题ID错误")
	if id == 0 {
		return
	}
	i := s.getIssue(id, s.US.GetUserID(ctx))
	if i.ID == 0 {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(iris.Map{"msg": "问题未找到"})
		return
	}
	var state int
	ctx.ReadJSON(&state)
	s.DB.Model(i).Update("state", state)
	ctx.JSON(i)
}

// Party 分组
func (s *IssueService) Party(p iris.Party) {
	p.Get("/", s.list)
	p.Get("/{id:int}", s.get)
	p.Post("/", s.post)
	p.Put("/{id:int}", s.put)
	p.Patch("/{id:int}/state", s.patchState)
}
