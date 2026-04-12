package entity

import "time"

// Favorite 文章收藏表（复合主键 user_id + article_id）
type Favorite struct {
	UserID    uint      `gorm:"primaryKey;column:user_id"`
	ArticleID uint      `gorm:"primaryKey;column:article_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (Favorite) TableName() string {
	return "favorites"
}
