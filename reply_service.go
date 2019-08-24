package common

import (
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
)

// ReplyService 回复服务
type ReplyService struct {
	DB *gorm.DB
	CS *Service
	IS *IssueService
}

// NewReplyService 新建问题服务
func NewReplyService(db *gorm.DB, cs *Service, is *IssueService) *ReplyService {
	db.AutoMigrate(&Reply{}) // 回复
	return &ReplyService{DB: db, CS: cs, IS: is}
}

func (s *ReplyService) list(ctx iris.Context) {
	i := s.IS.checkIssue(ctx)
	if i == nil {
		return
	}
	var rs []Reply
	s.DB.Model(i).Order("created_at").Where("user_id=?", i.UserID).Related(&rs)
	ctx.JSON(rs)
}

func (s *ReplyService) create(r *Reply, i *Issue) {
	r.UserID = i.UserID
	r.IssueID = i.ID
	s.DB.Create(r)
}
func (s *ReplyService) post(ctx iris.Context) {
	i := s.IS.checkIssue(ctx)
	if i == nil {
		return
	}
	var r Reply
	if s.CS.Bind(ctx, &r, "回复数据错误") {
		s.create(&r, i)
		ctx.JSON(r)
	}
}

func (s *ReplyService) getReply(id string) *Reply {
	var r Reply
	s.DB.Where("id=?", id).First(&r)
	return &r
}
func (s *ReplyService) update(r *Reply) (*Reply, bool) {
	// 存在判断
	if old := s.getReply(r.ID); old.ID != "" {
		old.Content = r.Content
		s.DB.Omit("created_at").Save(old)
		return old, true
	}
	return nil, false
}
func (s *ReplyService) put(ctx iris.Context) {
	var input Reply
	if !s.CS.Bind(ctx, &input, "回复数据错误") {
		return
	}
	input.ID = s.CS.String(ctx, "id", "回复ID错误")
	if input.ID == "" {
		return
	}
	if i, ok := s.update(&input); ok {
		ctx.JSON(i)
		return
	}
	ctx.StatusCode(iris.StatusUnauthorized)
	ctx.JSON(iris.Map{"msg": "回复未找到"})
}

// Party 分组
func (s *ReplyService) Party(p iris.Party) {
	p.Get("/", s.list)
	p.Post("/", s.post)
	p.Put("/{id:string}", s.put)
}
