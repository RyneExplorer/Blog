package repository

import "context"

// FavoriteRepository 收藏仓储接口（favorites 表，仅文章）
type FavoriteRepository interface {
	FavoriteArticleInTx(ctx context.Context, userID, articleID uint) error
	UnfavoriteArticleInTx(ctx context.Context, userID, articleID uint) error
}
