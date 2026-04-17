package repository

import (
	"context"
	"errors"

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
		Where("articles.status = ?", 2)
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
	q := r.db.WithContext(ctx).Model(&entity.Article{}).Where("status = ?", 2)
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
		Where("articles.id = ? AND articles.status = ?", id, 2).
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
	err := r.db.WithContext(ctx).Model(&entity.Article{}).Where("id = ? AND status = ?", id, 2).Count(&n).Error
	return n > 0, err
}

// IncrementViewInTx 使用事务更新浏览量（便于后续扩展流水等逻辑）
func (r *articleRepository) IncrementViewInTx(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&entity.Article{}).Where("id = ? AND status = ?", id, 2).
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

func (r *articleRepository) baseMyListQuery(ctx context.Context, userID uint, categoryID *uint) *gorm.DB {
	q := r.db.WithContext(ctx).Table("articles").
		Select(`articles.id, articles.title, articles.content, articles.summary, articles.cover_image, articles.status,
			articles.view_count, articles.like_count, articles.favorite_count, articles.comment_count,
			articles.created_at, articles.updated_at,
			categories.id AS category_ref_id, categories.name AS category_name, categories.slug AS category_slug`).
		Joins(`LEFT JOIN article_categories ac ON ac.article_id = articles.id AND ac.category_id = (
			SELECT MIN(ac2.category_id) FROM article_categories ac2 WHERE ac2.article_id = articles.id
		)`).
		Joins("LEFT JOIN categories ON categories.id = ac.category_id").
		Where("articles.user_id = ?", userID)
	if categoryID != nil {
		q = q.Where("articles.id IN (SELECT article_id FROM article_categories WHERE category_id = ?)", *categoryID)
	}
	return q
}

func (r *articleRepository) ListByUserWithJoin(ctx context.Context, userID uint, offset, limit int, categoryID *uint, sort string) ([]MyArticleListJoinRow, error) {
	q := r.baseMyListQuery(ctx, userID, categoryID).Session(&gorm.Session{})
	switch sort {
	case "hottest":
		q = q.Order("articles.view_count DESC, articles.created_at DESC")
	default:
		q = q.Order("articles.created_at DESC")
	}
	var rows []MyArticleListJoinRow
	if err := q.Offset(offset).Limit(limit).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *articleRepository) CountByUser(ctx context.Context, userID uint, categoryID *uint) (int64, error) {
	q := r.db.WithContext(ctx).Model(&entity.Article{}).Where("user_id = ?", userID)
	if categoryID != nil {
		q = q.Where("id IN (SELECT article_id FROM article_categories WHERE category_id = ?)", *categoryID)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (r *articleRepository) baseFavoritesQuery(ctx context.Context, userID uint, categoryID *uint) *gorm.DB {
	q := r.db.WithContext(ctx).Table("articles").
		Select(`articles.id, articles.title, articles.content, articles.summary, articles.cover_image, articles.status,
			articles.view_count, articles.like_count, articles.favorite_count, articles.comment_count,
			articles.created_at, articles.updated_at,
			categories.id AS category_ref_id, categories.name AS category_name, categories.slug AS category_slug`).
		Joins(`INNER JOIN favorites f ON f.article_id = articles.id AND f.user_id = ?`, userID).
		Joins(`LEFT JOIN article_categories ac ON ac.article_id = articles.id AND ac.category_id = (
			SELECT MIN(ac2.category_id) FROM article_categories ac2 WHERE ac2.article_id = articles.id
		)`).
		Joins("LEFT JOIN categories ON categories.id = ac.category_id")
	if categoryID != nil {
		q = q.Where("articles.id IN (SELECT article_id FROM article_categories WHERE category_id = ?)", *categoryID)
	}
	return q
}

func (r *articleRepository) ListFavoritesWithJoin(ctx context.Context, userID uint, offset, limit int, categoryID *uint, sort string) ([]MyArticleListJoinRow, error) {
	q := r.baseFavoritesQuery(ctx, userID, categoryID).Session(&gorm.Session{})
	switch sort {
	case "hottest":
		q = q.Order("articles.view_count DESC, articles.created_at DESC")
	default:
		q = q.Order("articles.created_at DESC")
	}
	var rows []MyArticleListJoinRow
	if err := q.Offset(offset).Limit(limit).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *articleRepository) CountFavorites(ctx context.Context, userID uint, categoryID *uint) (int64, error) {
	q := r.baseFavoritesQuery(ctx, userID, categoryID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (r *articleRepository) GetByIDWithCategories(ctx context.Context, id uint) (*entity.Article, error) {
	var a entity.Article
	err := r.db.WithContext(ctx).Model(&entity.Article{}).
		Preload("Categories").
		Where("id = ?", id).
		First(&a).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}

func uniqUint(in []uint) []uint {
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

func (r *articleRepository) CreateWithCategoriesInTx(ctx context.Context, article *entity.Article, categoryIDs []uint) error {
	categoryIDs = uniqUint(categoryIDs)
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(article).Error; err != nil {
			return err
		}
		if len(categoryIDs) == 0 {
			return nil
		}
		links := make([]entity.ArticleCategory, 0, len(categoryIDs))
		for _, cid := range categoryIDs {
			links = append(links, entity.ArticleCategory{
				ArticleID:  article.ID,
				CategoryID: cid,
			})
		}
		return tx.Model(&entity.ArticleCategory{}).Create(&links).Error
	})
}

func (r *articleRepository) UpdateByAuthorWithCategoriesInTx(ctx context.Context, id uint, userID uint, updates map[string]interface{}, categoryIDs []uint) (bool, error) {
	categoryIDs = uniqUint(categoryIDs)
	returned := false
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&entity.Article{}).Where("id = ? AND user_id = ?", id, userID).Updates(updates)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return nil
		}
		returned = true
		if err := tx.Where("article_id = ?", id).Delete(&entity.ArticleCategory{}).Error; err != nil {
			return err
		}
		if len(categoryIDs) == 0 {
			return nil
		}
		links := make([]entity.ArticleCategory, 0, len(categoryIDs))
		for _, cid := range categoryIDs {
			links = append(links, entity.ArticleCategory{
				ArticleID:  id,
				CategoryID: cid,
			})
		}
		return tx.Model(&entity.ArticleCategory{}).Create(&links).Error
	})
	return returned, err
}

func (r *articleRepository) UpdateStatusByAuthor(ctx context.Context, id uint, userID uint, status int) (bool, error) {
	res := r.db.WithContext(ctx).Model(&entity.Article{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("status", status)
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}

func (r *articleRepository) DeleteByAuthorInTx(ctx context.Context, id uint, userID uint) (bool, error) {
	deleted := false
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Where("id = ? AND user_id = ?", id, userID).Delete(&entity.Article{})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return nil
		}
		deleted = true
		return tx.Where("article_id = ?", id).Delete(&entity.ArticleCategory{}).Error
	})
	return deleted, err
}
