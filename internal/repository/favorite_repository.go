package repository

import (
	"context"
	"errors"

	"blog/internal/model/entity"

	"gorm.io/gorm"
)

// ErrUnfavoriteMissing 取消收藏时未找到记录
var ErrUnfavoriteMissing = errors.New("not_favorited")

type favoriteRepository struct {
	db *gorm.DB
}

// NewFavoriteRepository 创建收藏仓储
func NewFavoriteRepository(db *gorm.DB) FavoriteRepository {
	return &favoriteRepository{db: db}
}

// FavoriteArticleInTx 收藏文章：写入 favorites + 冗余 favorite_count
func (r *favoriteRepository) FavoriteArticleInTx(ctx context.Context, userID, articleID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var exists int64
		if err := tx.Model(&entity.Article{}).Where("id = ? AND status = ?", articleID, 1).Count(&exists).Error; err != nil {
			return err
		}
		if exists == 0 {
			return gorm.ErrRecordNotFound
		}
		rec := &entity.Favorite{UserID: userID, ArticleID: articleID}
		if err := tx.Create(rec).Error; err != nil {
			return err
		}
		return tx.Model(&entity.Article{}).Where("id = ?", articleID).
			UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error
	})
}

func (r *favoriteRepository) UnfavoriteArticleInTx(ctx context.Context, userID, articleID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Where("user_id = ? AND article_id = ?", userID, articleID).Delete(&entity.Favorite{})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return ErrUnfavoriteMissing
		}
		return tx.Model(&entity.Article{}).Where("id = ?", articleID).
			UpdateColumn("favorite_count", gorm.Expr("GREATEST(favorite_count - ?, 0)", 1)).Error
	})
}
