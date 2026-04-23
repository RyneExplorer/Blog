package entity

// Article 文章实体
// Categories 与 Category 为多对多关系，中间表为 article_categories。
type Article struct {
	BaseEntity
	UserID        uint       `gorm:"index;not null;comment:作者用户ID" json:"user_id"`
	User          User       `gorm:"foreignKey:UserID" json:"user"`
	Title         string     `gorm:"type:varchar(255);not null;comment:标题" json:"title"`
	Content       string     `gorm:"type:longtext;not null;comment:正文" json:"content"`
	Summary       string     `gorm:"type:varchar(255);comment:摘要" json:"summary"`
	CoverImage    string     `gorm:"type:varchar(500);default:'';comment:封面图" json:"cover_image"`
	RejectReason  string     `gorm:"type:varchar(255);default:'';comment:驳回原因" json:"reject_reason"`
	Status        int        `gorm:"type:tinyint;default:0;index;comment:状态 0草稿,1待审核,2已发布,3已驳回,4已封禁" json:"status"`
	ViewCount     int        `gorm:"type:int unsigned;default:0;comment:浏览量" json:"view_count"`
	LikeCount     int64      `gorm:"default:0;comment:点赞数" json:"like_count"`
	FavoriteCount int64      `gorm:"default:0;comment:收藏数" json:"favorite_count"`
	CommentCount  int64      `gorm:"default:0;comment:评论数" json:"comment_count"`
	Categories    []Category `gorm:"many2many:article_categories" json:"categories"`
}

func (Article) TableName() string {
	return "articles"
}
