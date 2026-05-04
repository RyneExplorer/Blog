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
		if err := tx.Model(&entity.Article{}).Where("id = ? AND status = ?", articleID, 2).Count(&exists).Error; err != nil {
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

// ListLikedArticleIDs 批量查询当前用户已点赞的文章 ID
// 1. 未登录或文章 ID 为空时直接返回空结果，避免无意义查询。
// 2. 只查询 likes 表中 target_type=article 的记录，保证文章和评论点赞互不干扰。
// 3. 转成 map 方便服务层 O(1) 填充每篇文章的 liked 状态。
func (r *likeRepository) ListLikedArticleIDs(ctx context.Context, userID uint, articleIDs []uint) (map[uint]bool, error) {
	result := make(map[uint]bool)
	if userID == 0 || len(articleIDs) == 0 {
		return result, nil
	}

	var ids []uint
	if err := r.db.WithContext(ctx).Model(&entity.Like{}).
		Where("user_id = ? AND target_type = ? AND target_id IN ?", userID, entity.LikeTargetArticle, articleIDs).
		Pluck("target_id", &ids).Error; err != nil {
		return nil, err
	}

	for _, id := range ids {
		result[id] = true
	}
	return result, nil
}

// ListLikedCommentIDs 批量查询当前用户已点赞的评论 ID
// 1. 未登录或评论 ID 为空时直接返回空结果，公开评论列表无需强制登录。
// 2. 只查询 likes 表中 target_type=comment 的记录，避免误用文章点赞记录。
// 3. 转成 map 供评论树组装时快速标记 liked 字段。
func (r *likeRepository) ListLikedCommentIDs(ctx context.Context, userID uint, commentIDs []uint) (map[uint]bool, error) {
	result := make(map[uint]bool)
	if userID == 0 || len(commentIDs) == 0 {
		return result, nil
	}

	var ids []uint
	if err := r.db.WithContext(ctx).Model(&entity.Like{}).
		Where("user_id = ? AND target_type = ? AND target_id IN ?", userID, entity.LikeTargetComment, commentIDs).
		Pluck("target_id", &ids).Error; err != nil {
		return nil, err
	}

	for _, id := range ids {
		result[id] = true
	}
	return result, nil
}
