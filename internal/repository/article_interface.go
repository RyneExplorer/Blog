package repository

import (
	"context"
	"database/sql"
	"time"

	"blog/internal/model/entity"
)

// ArticleRepository 文章仓储接口
type ArticleRepository interface {
	ListPublishedWithJoin(ctx context.Context, offset, limit int, categoryID *uint, sort string) ([]ArticleListJoinRow, error)
	CountPublished(ctx context.Context, categoryID *uint) (int64, error)
	GetPublishedDetailJoin(ctx context.Context, id uint) (*ArticleDetailJoinRow, error)
	IncrementViewInTx(ctx context.Context, id uint) error
	ExistsPublished(ctx context.Context, id uint) (bool, error)

	// 用户文章模块
	ListByUserWithJoin(ctx context.Context, userID uint, offset, limit int, categoryID *uint, sort string) ([]MyArticleListJoinRow, error)
	CountByUser(ctx context.Context, userID uint, categoryID *uint) (int64, error)
	GetByIDWithCategories(ctx context.Context, id uint) (*entity.Article, error)
	CreateWithCategoriesInTx(ctx context.Context, article *entity.Article, categoryIDs []uint) error
	UpdateByAuthorWithCategoriesInTx(ctx context.Context, id uint, userID uint, updates map[string]interface{}, categoryIDs []uint) (bool, error)
	UpdateStatusByAuthor(ctx context.Context, id uint, userID uint, status int) (bool, error)
	DeleteByAuthorInTx(ctx context.Context, id uint, userID uint) (bool, error)

	// 用户收藏模块
	ListFavoritesWithJoin(ctx context.Context, userID uint, offset, limit int, categoryID *uint, sort string) ([]MyArticleListJoinRow, error)
	CountFavorites(ctx context.Context, userID uint, categoryID *uint) (int64, error)
}

// ArticleListJoinRow 列表连表查询结果
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

// ArticleDetailJoinRow 详情连表查询结果
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

// MyArticleListJoinRow “我的文章列表”连表查询结果（返回文章 + 一个分类）
type MyArticleListJoinRow struct {
	ID            uint           `gorm:"column:id"`
	Title         string         `gorm:"column:title"`
	Content       string         `gorm:"column:content"`
	Summary       sql.NullString `gorm:"column:summary"`
	CoverImage    string         `gorm:"column:cover_image"`
	Status        int            `gorm:"column:status"`
	ViewCount     int            `gorm:"column:view_count"`
	LikeCount     int64          `gorm:"column:like_count"`
	FavoriteCount int64          `gorm:"column:favorite_count"`
	CommentCount  int64          `gorm:"column:comment_count"`
	CreatedAt     time.Time      `gorm:"column:created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at"`
	CategoryRefID sql.NullInt64  `gorm:"column:category_ref_id"`
	CategoryName  string         `gorm:"column:category_name"`
	CategorySlug  string         `gorm:"column:category_slug"`
}
