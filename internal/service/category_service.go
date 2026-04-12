package service

import (
	"context"

	dto "blog/internal/model/dto/response"
	"blog/internal/repository"
)

// CategoryService 分类业务
type CategoryService interface {
	ListAll(ctx context.Context) ([]dto.CategoryItem, error)
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

// NewCategoryService 创建分类服务
func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{categoryRepo: categoryRepo}
}

func (s *categoryService) ListAll(ctx context.Context) ([]dto.CategoryItem, error) {
	list, err := s.categoryRepo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]dto.CategoryItem, 0, len(list))
	for _, c := range list {
		out = append(out, dto.CategoryItem{
			ID:   c.ID,
			Name: c.Name,
			Slug: c.Slug,
		})
	}
	return out, nil
}
