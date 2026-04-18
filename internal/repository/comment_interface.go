package repository

import (
	"context"
	"database/sql"
	"time"

	"blog/internal/model/entity"
)

// CommentJoinRow 评论列表连表查询结果
type CommentJoinRow struct {
	ID           uint          `gorm:"column:id"`
	ParentID     sql.NullInt64 `gorm:"column:parent_id"`
	RootID       sql.NullInt64 `gorm:"column:root_id"`
	Content      string        `gorm:"column:content"`
	LikeCount    int           `gorm:"column:like_count"`
	ReplyCount   int64         `gorm:"column:reply_count"`
	CreatedAt    time.Time     `gorm:"column:created_at"`
	UpdatedAt    time.Time     `gorm:"column:updated_at"`
	ArticleID    uint          `gorm:"column:article_id"`
	UserID       uint          `gorm:"column:user_id"`
	UserNickname string        `gorm:"column:user_nickname"`
	UserAvatar   string        `gorm:"column:user_avatar"`
}

// CommentRepository 评论仓储接口
type CommentRepository interface {
	CountRootsByArticle(ctx context.Context, articleID uint) (int64, error)
	ListJoinedByArticlePage(ctx context.Context, articleID uint, limit, offset int) ([]CommentJoinRow, error)
	GetByID(ctx context.Context, id uint) (*entity.Comment, error)
	CountChildren(ctx context.Context, parentID uint) (int64, error)
	CreateWithCountersInTx(ctx context.Context, c *entity.Comment) error
	DeleteWithCountersInTx(ctx context.Context, c *entity.Comment) error
}
