package repository

import (
	"context"
	"database/sql"
	"time"
)

// ReviewRepository 管理员审核仓储
type ReviewRepository interface {
	ListForAdmin(ctx context.Context, req *AdminListFilter, offset, limit int) ([]AdminArticleJoinRow, error)
	CountForAdmin(ctx context.Context, req *AdminListFilter) (int64, error)
	GetDetail(ctx context.Context, id uint) (*AdminArticleDetailRow, error)

	ApproveInTx(ctx context.Context, articleID uint, adminID uint) (bool, error)
	RejectInTx(ctx context.Context, articleID uint, adminID uint, reason string) (bool, error)
	BanInTx(ctx context.Context, articleID uint, adminID uint, reason string) (bool, error)
	UpdateCategoriesInTx(ctx context.Context, articleID uint, categoryIDs []uint, adminID uint) (bool, error)
}

// AdminListFilter 管理员列表筛选条件
type AdminListFilter struct {
	CategoryID *uint
	Username   string
	Status     *int
}

// AdminArticleJoinRow 管理员列表连表查询结果
type AdminArticleJoinRow struct {
	ID            uint           `gorm:"column:id"`
	Title         string         `gorm:"column:title"`
	Content       string         `gorm:"column:content"`
	Summary       sql.NullString `gorm:"column:summary"`
	CoverImage    string         `gorm:"column:cover_image"`
	RejectReason  sql.NullString `gorm:"column:reject_reason"`
	Status        int            `gorm:"column:status"`
	ViewCount     int            `gorm:"column:view_count"`
	LikeCount     int64          `gorm:"column:like_count"`
	FavoriteCount int64          `gorm:"column:favorite_count"`
	CommentCount  int64          `gorm:"column:comment_count"`
	CreatedAt     time.Time      `gorm:"column:created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at"`
	UserID        uint           `gorm:"column:user_id"`
	Username      string         `gorm:"column:username"`
	Nickname      string         `gorm:"column:nickname"`
	Avatar        string         `gorm:"column:avatar"`
	Bio           string         `gorm:"column:bio"`
	CategoryRefID sql.NullInt64  `gorm:"column:category_ref_id"`
	CategoryName  string         `gorm:"column:category_name"`
	CategorySlug  string         `gorm:"column:category_slug"`
}

// AdminArticleDetailRow 管理员详情行（包含内容）
type AdminArticleDetailRow struct {
	ID            uint           `gorm:"column:id"`
	UserID        uint           `gorm:"column:user_id"`
	Title         string         `gorm:"column:title"`
	Content       string         `gorm:"column:content"`
	Summary       sql.NullString `gorm:"column:summary"`
	CoverImage    string         `gorm:"column:cover_image"`
	RejectReason  sql.NullString `gorm:"column:reject_reason"`
	Status        int            `gorm:"column:status"`
	ViewCount     int            `gorm:"column:view_count"`
	LikeCount     int64          `gorm:"column:like_count"`
	FavoriteCount int64          `gorm:"column:favorite_count"`
	CommentCount  int64          `gorm:"column:comment_count"`
	CreatedAt     time.Time      `gorm:"column:created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at"`
	Username      string         `gorm:"column:username"`
	Nickname      string         `gorm:"column:nickname"`
	Avatar        string         `gorm:"column:avatar"`
	Bio           string         `gorm:"column:bio"`
}
