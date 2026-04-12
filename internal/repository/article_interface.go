package repository

import (
	"context"
	"database/sql"
	"time"
)

// ArticleRepository 文章仓储接口
type ArticleRepository interface {
	ListPublishedWithJoin(ctx context.Context, offset, limit int, categoryID *uint, sort string) ([]ArticleListJoinRow, error)
	CountPublished(ctx context.Context, categoryID *uint) (int64, error)
	GetPublishedDetailJoin(ctx context.Context, id uint) (*ArticleDetailJoinRow, error)
	IncrementViewInTx(ctx context.Context, id uint) error
	ExistsPublished(ctx context.Context, id uint) (bool, error)
}

// ArticleListJoinRow 列表 JOIN 扫描行
type ArticleListJoinRow struct {
	ID             uint           `gorm:"column:id"`
	Title          string         `gorm:"column:title"`
	Content        string         `gorm:"column:content"`
	Summary        sql.NullString `gorm:"column:summary"`
	CoverImage     string         `gorm:"column:cover_image"`
	Status         int            `gorm:"column:status"`
	ViewCount      int            `gorm:"column:view_count"`
	LikeCount      int64          `gorm:"column:like_count"`
	FavoriteCount  int64          `gorm:"column:favorite_count"`
	CommentCount   int64          `gorm:"column:comment_count"`
	CreatedAt      time.Time      `gorm:"column:created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at"`
	UserID         uint           `gorm:"column:user_id"`
	AuthorNickname string         `gorm:"column:author_nickname"`
	AuthorAvatar   string         `gorm:"column:author_avatar"`
	AuthorBio      string         `gorm:"column:author_bio"`
	CategoryRefID  sql.NullInt64  `gorm:"column:category_ref_id"`
	CategoryName   string         `gorm:"column:category_name"`
	CategorySlug   string         `gorm:"column:category_slug"`
}

// ArticleDetailJoinRow 详情 JOIN 扫描行
type ArticleDetailJoinRow struct {
	ID            uint      `gorm:"column:id"`
	Title         string    `gorm:"column:title"`
	Content       string    `gorm:"column:content"`
	CoverImage    string    `gorm:"column:cover_image"`
	Status        int       `gorm:"column:status"`
	ViewCount     int       `gorm:"column:view_count"`
	LikeCount     int64     `gorm:"column:like_count"`
	FavoriteCount int64     `gorm:"column:favorite_count"`
	CommentCount  int64     `gorm:"column:comment_count"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
	CategoryName  string    `gorm:"column:category_name"`
	Nickname      string    `gorm:"column:nickname"`
	Bio           string    `gorm:"column:bio"`
	Avatar        string    `gorm:"column:avatar"`
}
