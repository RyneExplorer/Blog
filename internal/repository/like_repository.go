package repository

import (
	"context"
	"errors"

	"blog/internal/model/entity"

	"gorm.io/gorm"
)

// 与业务层约定：取消点赞时未找到记录
var (
	ErrArticleUnlikeMissing = errors.New("not_liked")
	ErrCommentUnlikeMissing = errors.New("comment_not_liked")
)

type likeRepository struct {
	db *gorm.DB
}

// NewLikeRepository 创建点赞仓储
func NewLikeRepository(db *gorm.DB) LikeRepository {
	return &likeRepository{db: db}
}

// LikeArticleInTx 文章点赞：写入 likes + 冗余 like_count
func (r *likeRepository) LikeArticleInTx(ctx context.Context, userID, articleID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var exists int64
		if err := tx.Model(&entity.Article{}).Where("id = ? AND status = ?", articleID, 1).Count(&exists).Error; err != nil {
			return err
		}
		if exists == 0 {
			return gorm.ErrRecordNotFound
		}
		rec := &entity.Like{
			UserID:     userID,
			TargetType: entity.LikeTargetArticle,
			TargetID:   articleID,
		}
		if err := tx.Create(rec).Error; err != nil {
			return err
		}
		return tx.Model(&entity.Article{}).Where("id = ?", articleID).
			UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error
	})
}

func (r *likeRepository) UnlikeArticleInTx(ctx context.Context, userID, articleID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Where("user_id = ? AND target_type = ? AND target_id = ?", userID, entity.LikeTargetArticle, articleID).
			Delete(&entity.Like{})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return ErrArticleUnlikeMissing
		}
		return tx.Model(&entity.Article{}).Where("id = ?", articleID).
			UpdateColumn("like_count", gorm.Expr("GREATEST(like_count - ?, 0)", 1)).Error
	})
}

// LikeCommentInTx 评论点赞：写入 likes + 冗余 like_count
func (r *likeRepository) LikeCommentInTx(ctx context.Context, userID, commentID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var exists int64
		if err := tx.Model(&entity.Comment{}).Where("id = ? AND status = ?", commentID, 1).Count(&exists).Error; err != nil {
			return err
		}
		if exists == 0 {
			return gorm.ErrRecordNotFound
		}
		rec := &entity.Like{
			UserID:     userID,
			TargetType: entity.LikeTargetComment,
			TargetID:   commentID,
		}
		if err := tx.Create(rec).Error; err != nil {
			return err
		}
		return tx.Model(&entity.Comment{}).Where("id = ?", commentID).
			UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error
	})
}

func (r *likeRepository) UnlikeCommentInTx(ctx context.Context, userID, commentID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Where("user_id = ? AND target_type = ? AND target_id = ?", userID, entity.LikeTargetComment, commentID).
			Delete(&entity.Like{})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return ErrCommentUnlikeMissing
		}
		return tx.Model(&entity.Comment{}).Where("id = ?", commentID).
			UpdateColumn("like_count", gorm.Expr("GREATEST(like_count - ?, 0)", 1)).Error
	})
}
