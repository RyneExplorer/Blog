package entity

// Comment 评论（楼中楼：一级 parent_id 为空；子评论带 root_id 指向一级评论）
type Comment struct {
	BaseEntity
	ArticleID  uint   `gorm:"index;not null;comment:文章ID" json:"article_id"`
	UserID     uint   `gorm:"index;not null;comment:评论用户ID" json:"user_id"`
	ParentID   *uint  `gorm:"index;comment:父评论ID，一级为空" json:"parent_id"`
	RootID     *uint  `gorm:"index;comment:一级评论ID，一级为空" json:"root_id"`
	Content    string `gorm:"type:text;not null;comment:内容" json:"content"`
	Status     int    `gorm:"type:tinyint;not null;default:1;index;comment:0待审核,1已发布,2已拒绝" json:"status"`
	LikeCount  int    `gorm:"type:int unsigned;default:0;comment:点赞数" json:"like_count"`
	ReplyCount int64  `gorm:"default:0;comment:直接回复数" json:"reply_count"`
}

func (Comment) TableName() string {
	return "comments"
}
