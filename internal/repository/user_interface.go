package repository

import (
	"blog/internal/model/entity"
)

// UserListFilter 用户列表筛选条件
type UserListFilter struct {
	Username string
	Nickname string
	Status   *int
}

// UserRepository 用户仓储接口
type UserRepository interface {
	FindByID(id uint) (*entity.User, error)
	FindByUsername(username string) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	Create(user *entity.User) error
	Update(id uint, updates map[string]interface{}) error
	Delete(id uint) error
	AdminList(offset, limit int, filter *UserListFilter) ([]*entity.User, int64, error)
	ExistsByUsername(username string) (bool, error)
	ExistsByEmail(email string) (bool, error)
}
