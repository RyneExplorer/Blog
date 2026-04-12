package repository

import (
	"context"

	"blog/internal/model/entity"

	"gorm.io/gorm"
)

type articleRepository struct {
	db *gorm.DB
}

// NewArticleRepository 创建文章仓储
func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{db: db}
}

// baseListQuery 文章列表基础查询（分类仅通过 article_categories，文章表无 category_id）
func (r *articleRepository) baseListQuery(ctx context.Context, categoryID *uint) *gorm.DB {
	q := r.db.WithContext(ctx).Table("articles").
		Select(`articles.id, articles.title, articles.content, articles.summary, articles.cover_image, articles.status,
			articles.view_count, articles.like_count, articles.favorite_count, articles.comment_count,
			articles.created_at, articles.updated_at, articles.user_id,
			users.nickname AS author_nickname, users.avatar AS author_avatar, users.bio AS author_bio,
			categories.id AS category_ref_id, categories.name AS category_name, categories.slug AS category_slug`).
		Joins("INNER JOIN users ON users.id = articles.user_id").
		Joins(`LEFT JOIN article_categories ac ON ac.article_id = articles.id AND ac.category_id = (
			SELECT MIN(ac2.category_id) FROM article_categories ac2 WHERE ac2.article_id = articles.id
		)`).
		Joins("LEFT JOIN categories ON categories.id = ac.category_id").
		Where("articles.status = ?", 1)
	if categoryID != nil {
		q = q.Where("articles.id IN (SELECT article_id FROM article_categories WHERE category_id = ?)", *categoryID)
	}
	return q
}

func (r *articleRepository) ListPublishedWithJoin(ctx context.Context, offset, limit int, categoryID *uint, sort string) ([]ArticleListJoinRow, error) {
	q := r.baseListQuery(ctx, categoryID).Session(&gorm.Session{})
	switch sort {
	case "hottest":
		q = q.Order("(articles.like_count + articles.view_count) DESC, articles.created_at DESC")
	default:
		q = q.Order("articles.created_at DESC")
	}
	var rows []ArticleListJoinRow
	if err := q.Offset(offset).Limit(limit).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *articleRepository) CountPublished(ctx context.Context, categoryID *uint) (int64, error) {
	q := r.db.WithContext(ctx).Model(&entity.Article{}).Where("status = ?", 1)
	if categoryID != nil {
		q = q.Where("id IN (SELECT article_id FROM article_categories WHERE category_id = ?)", *categoryID)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (r *articleRepository) GetPublishedDetailJoin(ctx context.Context, id uint) (*ArticleDetailJoinRow, error) {
	var row ArticleDetailJoinRow
	err := r.db.WithContext(ctx).Table("articles").
		Select(`articles.id, articles.title, articles.content, articles.cover_image, articles.status,
			articles.view_count, articles.like_count, articles.favorite_count, articles.comment_count,
			articles.created_at, articles.updated_at,
			(SELECT c.name FROM article_categories ac
			 INNER JOIN categories c ON c.id = ac.category_id
			 WHERE ac.article_id = articles.id ORDER BY ac.category_id ASC LIMIT 1) AS category_name,
			users.nickname AS nickname, users.bio AS bio, users.avatar AS avatar`).
		Joins("INNER JOIN users ON users.id = articles.user_id").
		Where("articles.id = ? AND articles.status = ?", id, 1).
		Scan(&row).Error
	if err != nil {
		return nil, err
	}
	if row.ID == 0 {
		return nil, nil
	}
	return &row, nil
}

func (r *articleRepository) ExistsPublished(ctx context.Context, id uint) (bool, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&entity.Article{}).Where("id = ? AND status = ?", id, 1).Count(&n).Error
	return n > 0, err
}

// IncrementViewInTx 使用事务更新浏览量（便于后续扩展流水等逻辑）
func (r *articleRepository) IncrementViewInTx(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&entity.Article{}).Where("id = ? AND status = ?", id, 1).
			UpdateColumn("view_count", gorm.Expr("view_count + ?", 1))
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}
