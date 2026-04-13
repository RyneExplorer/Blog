package repository

import (
	"context"
	"strings"

	"blog/internal/model/entity"

	"gorm.io/gorm"
)

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) baseAdminQuery(ctx context.Context, f *AdminListFilter) *gorm.DB {
	q := r.db.WithContext(ctx).Table("articles").
		Select(`articles.id, articles.title, articles.content, articles.summary, articles.cover_image, articles.reject_reason, articles.status,
			articles.view_count, articles.like_count, articles.favorite_count, articles.comment_count,
			articles.created_at, articles.updated_at, articles.user_id,
			u.username AS username, u.nickname AS nickname, u.avatar AS avatar, u.bio AS bio,
			categories.id AS category_ref_id, categories.name AS category_name, categories.slug AS category_slug`).
		Joins("INNER JOIN users u ON u.id = articles.user_id").
		Joins(`LEFT JOIN article_categories ac ON ac.article_id = articles.id AND ac.category_id = (
			SELECT MIN(ac2.category_id) FROM article_categories ac2 WHERE ac2.article_id = articles.id
		)`).
		Joins("LEFT JOIN categories ON categories.id = ac.category_id")

	if f != nil {
		if f.Status != nil {
			q = q.Where("articles.status = ?", *f.Status)
		}
		if f.CategoryID != nil {
			q = q.Where("articles.id IN (SELECT article_id FROM article_categories WHERE category_id = ?)", *f.CategoryID)
		}
		if strings.TrimSpace(f.Username) != "" {
			q = q.Where("u.username LIKE ?", "%"+strings.TrimSpace(f.Username)+"%")
		}
	}
	return q
}

func (r *reviewRepository) ListForAdmin(ctx context.Context, req *AdminListFilter, offset, limit int) ([]AdminArticleJoinRow, error) {
	q := r.baseAdminQuery(ctx, req).Order("articles.created_at DESC")
	var rows []AdminArticleJoinRow
	if err := q.Offset(offset).Limit(limit).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *reviewRepository) CountForAdmin(ctx context.Context, req *AdminListFilter) (int64, error) {
	q := r.baseAdminQuery(ctx, req).Select("COUNT(DISTINCT articles.id)")
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (r *reviewRepository) GetDetail(ctx context.Context, id uint) (*AdminArticleDetailRow, error) {
	var row AdminArticleDetailRow
	err := r.db.WithContext(ctx).Table("articles").
		Select(`articles.id, articles.user_id, articles.title, articles.content, articles.summary, articles.cover_image, articles.reject_reason, articles.status,
			articles.view_count, articles.like_count, articles.favorite_count, articles.comment_count,
			articles.created_at, articles.updated_at,
			u.username AS username, u.nickname AS nickname, u.avatar AS avatar, u.bio AS bio`).
		Joins("INNER JOIN users u ON u.id = articles.user_id").
		Where("articles.id = ?", id).
		Scan(&row).Error
	if err != nil {
		return nil, err
	}
	if row.ID == 0 {
		return nil, nil
	}
	return &row, nil
}

func uniqUintReview(in []uint) []uint {
	if len(in) == 0 {
		return nil
	}
	seen := make(map[uint]struct{}, len(in))
	out := make([]uint, 0, len(in))
	for _, v := range in {
		if v == 0 {
			continue
		}
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}

func (r *reviewRepository) ApproveInTx(ctx context.Context, articleID uint, adminID uint) (bool, error) {
	ok := false
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&entity.Article{}).Where("id = ? AND status = ?", articleID, 1).Updates(map[string]interface{}{
			"status":        2,
			"reject_reason": "",
		})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return nil
		}
		ok = true
		return tx.Create(&entity.ReviewLog{
			ArticleID: articleID,
			AdminID:   adminID,
			Action:    "approve",
			Reason:    "",
		}).Error
	})
	return ok, err
}

func (r *reviewRepository) RejectInTx(ctx context.Context, articleID uint, adminID uint, reason string) (bool, error) {
	ok := false
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&entity.Article{}).Where("id = ? AND status = ?", articleID, 1).Updates(map[string]interface{}{
			"status":        3,
			"reject_reason": strings.TrimSpace(reason),
		})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return nil
		}
		ok = true
		return tx.Create(&entity.ReviewLog{
			ArticleID: articleID,
			AdminID:   adminID,
			Action:    "reject",
			Reason:    strings.TrimSpace(reason),
		}).Error
	})
	return ok, err
}

func (r *reviewRepository) BanInTx(ctx context.Context, articleID uint, adminID uint, reason string) (bool, error) {
	ok := false
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&entity.Article{}).Where("id = ?", articleID).Update("status", 4)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return nil
		}
		ok = true
		return tx.Create(&entity.ReviewLog{
			ArticleID: articleID,
			AdminID:   adminID,
			Action:    "ban",
			Reason:    strings.TrimSpace(reason),
		}).Error
	})
	return ok, err
}

func (r *reviewRepository) UpdateCategoriesInTx(ctx context.Context, articleID uint, categoryIDs []uint, adminID uint) (bool, error) {
	categoryIDs = uniqUintReview(categoryIDs)
	ok := false
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var cnt int64
		if err := tx.Model(&entity.Article{}).Where("id = ?", articleID).Count(&cnt).Error; err != nil {
			return err
		}
		if cnt == 0 {
			return nil
		}
		ok = true
		if err := tx.Where("article_id = ?", articleID).Delete(&entity.ArticleCategory{}).Error; err != nil {
			return err
		}
		if len(categoryIDs) > 0 {
			links := make([]entity.ArticleCategory, 0, len(categoryIDs))
			for _, cid := range categoryIDs {
				links = append(links, entity.ArticleCategory{ArticleID: articleID, CategoryID: cid})
			}
			if err := tx.Create(&links).Error; err != nil {
				return err
			}
		}
		return tx.Create(&entity.ReviewLog{
			ArticleID: articleID,
			AdminID:   adminID,
			Action:    "category",
			Reason:    "",
		}).Error
	})
	return ok, err
}
