package repository

import (
	"context"

	"blog/internal/model/entity"

	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository 创建分类仓储
func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) ListAll(ctx context.Context) ([]entity.Category, error) {
	var list []entity.Category
	err := r.db.WithContext(ctx).Model(&entity.Category{}).Order("id ASC").Find(&list).Error
	return list, err
}
