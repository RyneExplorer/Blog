package repository

import (
	"context"

	"blog/internal/model/entity"
)

// CategoryRepository 分类仓储接口
type CategoryRepository interface {
	ListAll(ctx context.Context) ([]entity.Category, error)
}
