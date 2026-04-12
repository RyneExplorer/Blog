package entity

import "time"

// ArticleCategory 文章与分类的中间表（与 Article.Categories / Category.Articles 的 many2many 对应）
type ArticleCategory struct {
	ArticleID  uint      `gorm:"primaryKey;column:article_id;comment:文章ID"`
	CategoryID uint      `gorm:"primaryKey;column:category_id;comment:分类ID"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (ArticleCategory) TableName() string {
	return "article_categories"
}
