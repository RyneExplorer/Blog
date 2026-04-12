package repository

import (
	"context"
	"errors"

	"blog/internal/model/entity"

	"gorm.io/gorm"
)

type commentRepository struct {
	db *gorm.DB
}

// NewCommentRepository 创建评论仓储
func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) CountRootsByArticle(ctx context.Context, articleID uint) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&entity.Comment{}).
		Where("article_id = ? AND parent_id IS NULL AND status = ?", articleID, 1).
		Count(&n).Error
	return n, err
}

// ListJoinedByArticlePage 一次查询取出当前页一级评论及其下所有子评论（CTE），并 JOIN 用户表
func (r *commentRepository) ListJoinedByArticlePage(ctx context.Context, articleID uint, limit, offset int) ([]CommentJoinRow, error) {
	const q = `
WITH roots AS (
  SELECT c.id
  FROM comments c
  WHERE c.article_id = ? AND c.parent_id IS NULL AND c.status = 1
  ORDER BY c.created_at DESC
  LIMIT ? OFFSET ?
)
SELECT c.id, c.parent_id, c.root_id, c.content, c.like_count, c.reply_count,
       c.created_at, c.updated_at, c.article_id, c.user_id,
       u.nickname AS user_nickname, u.avatar AS user_avatar
FROM comments c
INNER JOIN users u ON u.id = c.user_id
WHERE c.article_id = ? AND c.status = 1
  AND (c.id IN (SELECT id FROM roots) OR c.root_id IN (SELECT id FROM roots))
ORDER BY c.created_at ASC
`
	var rows []CommentJoinRow
	if err := r.db.WithContext(ctx).Raw(q, articleID, limit, offset, articleID).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *commentRepository) GetByID(ctx context.Context, id uint) (*entity.Comment, error) {
	var c entity.Comment
	err := r.db.WithContext(ctx).First(&c, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *commentRepository) CountChildren(ctx context.Context, parentID uint) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&entity.Comment{}).Where("parent_id = ?", parentID).Count(&n).Error
	return n, err
}

// CreateWithCountersInTx 发表评论/回复：插入评论并维护文章评论数、父评论回复数
func (r *commentRepository) CreateWithCountersInTx(ctx context.Context, c *entity.Comment) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(c).Error; err != nil {
			return err
		}
		if err := tx.Model(&entity.Article{}).Where("id = ?", c.ArticleID).
			UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1)).Error; err != nil {
			return err
		}
		if c.ParentID != nil {
			if err := tx.Model(&entity.Comment{}).Where("id = ?", *c.ParentID).
				UpdateColumn("reply_count", gorm.Expr("reply_count + ?", 1)).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// DeleteWithCountersInTx 删除评论并回退计数
func (r *commentRepository) DeleteWithCountersInTx(ctx context.Context, c *entity.Comment) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Delete(&entity.Comment{}, c.ID)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		if err := tx.Model(&entity.Article{}).Where("id = ?", c.ArticleID).
			UpdateColumn("comment_count", gorm.Expr("GREATEST(comment_count - ?, 0)", 1)).Error; err != nil {
			return err
		}
		if c.ParentID != nil {
			if err := tx.Model(&entity.Comment{}).Where("id = ?", *c.ParentID).
				UpdateColumn("reply_count", gorm.Expr("GREATEST(reply_count - ?, 0)", 1)).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
