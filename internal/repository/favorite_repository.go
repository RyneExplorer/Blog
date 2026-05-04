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
		if err := tx.Model(&entity.Article{}).Where("id = ? AND status = ?", articleID, 2).Count(&exists).Error; err != nil {
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

// ListFavoritedArticleIDs 批量查询当前用户已收藏的文章 ID
// 1. 未登录或文章 ID 为空时返回空 map，公开列表无需强制登录。
// 2. 只查 favorites 关系表，不依赖前端 localStorage 推断收藏状态。
// 3. 转成 map 方便服务层快速填充每篇文章的 favorited 状态。
func (r *favoriteRepository) ListFavoritedArticleIDs(ctx context.Context, userID uint, articleIDs []uint) (map[uint]bool, error) {
	result := make(map[uint]bool)
	if userID == 0 || len(articleIDs) == 0 {
		return result, nil
	}

	var ids []uint
	if err := r.db.WithContext(ctx).Model(&entity.Favorite{}).
		Where("user_id = ? AND article_id IN ?", userID, articleIDs).
		Pluck("article_id", &ids).Error; err != nil {
		return nil, err
	}

	for _, id := range ids {
		result[id] = true
	}
	return result, nil
}
